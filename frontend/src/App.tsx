import { useState, useEffect, useRef } from 'react'
import { Shield, Search, RefreshCw, Download, AlertTriangle, CheckCircle, XCircle, HelpCircle, FolderOpen, Edit3, Save, BadgeCheck, ArrowRightLeft, Monitor, Terminal, Play, ChevronDown } from 'lucide-react'
import { 
  ValidateExtension, 
  GetInstalledExtensions, 
  AuditAllExtensions,
  GetDefaultExtensionsPath,
  SelectDirectory,
  DownloadOfficialExtension,
  SearchMarketplaceExtension,
  SearchMarketplace,
  GetEditorProfiles,
  GetAllEditorStatuses,
  GetEditorExtensions,
  GetCLIStatus,
  InstallExtensionViaCLI,
  SyncExtensions,
  DetectSyncConflicts,
  GetCLIInstallStatus,
  InstallCLI,
  UninstallCLI,
} from './wailsjs/go/main/App'

interface ExtensionMetadata {
  id: string
  publisher: string
  publisherDomain?: string
  isVerifiedPublisher: boolean
  name: string
  version: string
  displayName: string
  description: string
  repositoryUrl?: string
  downloadUrl?: string
  source: string
}

interface Extension {
  id: string
  path: string
  publisher: string
  name: string
  version: string
  isEnabled: boolean
  lastModified: string
}

interface ValidationResult {
  extensionId: string
  trustLevel: string
  marketplaceData?: any
  openvsxData?: any
  differences?: string[]
  recommendation: string
  validationTime: string
  error?: string
}

interface AuditReport {
  totalExtensions: number
  legitimateCount: number
  suspiciousCount: number
  maliciousCount: number
  unknownCount: number
  results: ValidationResult[]
  auditTime: string
}

interface EditorProfile {
  id: string
  name: string
  extensionsDir: string
  indexFile?: string
  cliCommand?: string
  isVSCodeFamily?: boolean
  isCustom?: boolean
}

interface EditorStatus {
  editor: EditorProfile
  dirExists?: boolean
  indexFileExists?: boolean
  cliAvailable?: boolean
  cliPath?: string
  extensionCount?: number
  disabledReason?: string
  isAvailable?: boolean
}

interface CLIStatus {
  vscodeAvailable?: boolean
  vscodePath?: string
  insidersAvailable?: boolean
  insidersPath?: string
  codiumAvailable?: boolean
  codiumPath?: string
  anyAvailable?: boolean
  preferredCli?: string
}

interface SyncResult {
  targetEditor: string
  success?: boolean
  copiedCount?: number
  skippedCount?: number
  overwrittenCount?: number
  indexUpdated?: boolean
  conflicts?: string[]
  errors?: string[]
}

interface SyncReport {
  sourceEditor: string
  results?: SyncResult[]
  totalCopied?: number
  totalSkipped?: number
  totalErrors?: number
}

function App() {
  const [extensions, setExtensions] = useState<Extension[]>([])
  // Cache audit reports per editor so user can switch between editors
  const [auditReportsCache, setAuditReportsCache] = useState<Record<string, AuditReport>>({})
  const [auditReport, setAuditReport] = useState<AuditReport | null>(null)
  const [selectedExtension, setSelectedExtension] = useState<Extension | null>(null)
  const [validationResult, setValidationResult] = useState<ValidationResult | null>(null)
  const [searchQuery, setSearchQuery] = useState('')
  const [loading, setLoading] = useState(false)
  const [extensionsPath, setExtensionsPath] = useState('')
  const [view, setView] = useState<'list' | 'audit' | 'search' | 'sync' | 'settings'>('list')
  const [error, setError] = useState<string | null>(null)
  
  // Editor-related state
  const [editorProfiles, setEditorProfiles] = useState<EditorProfile[]>([])
  const [editorStatuses, setEditorStatuses] = useState<EditorStatus[]>([])
  const [selectedEditor, setSelectedEditor] = useState<string>('vscode')
  const [cliStatus, setCLIStatus] = useState<CLIStatus | null>(null)
  
  // Sync-related state
  const [syncSourceEditor, setSyncSourceEditor] = useState<string>('vscode')
  const [syncTargetEditors, setSyncTargetEditors] = useState<string[]>([])
  const [syncSelectedExtensions, setSyncSelectedExtensions] = useState<string[]>([])
  const [syncReport, setSyncReport] = useState<SyncReport | null>(null)
  const [syncConflicts, setSyncConflicts] = useState<string[]>([])
  const [showConflictDialog, setShowConflictDialog] = useState(false)
  const [targetEditorExtensions, setTargetEditorExtensions] = useState<Record<string, string[]>>({})
  const [syncFilterMode, setSyncFilterMode] = useState<'all' | 'missing' | 'present'>('all')
  const [syncSearchFilter, setSyncSearchFilter] = useState('')
  
  // Install-related state  
  const [installTargetEditor, setInstallTargetEditor] = useState<string>('vscode')
  const [installing, setInstalling] = useState(false)
  const [isEditingPath, setIsEditingPath] = useState(false)
  const [tempPath, setTempPath] = useState('')
  const [navCollapsed, setNavCollapsed] = useState(false)
  const [_vsynxCliStatus, setVsynxCliStatus] = useState<any>(null)
  const [cliInstalling, setCliInstalling] = useState(false)
  const [marketplaceSearchQuery, setMarketplaceSearchQuery] = useState('')
  const [marketplaceSearchResults, setMarketplaceSearchResults] = useState<ExtensionMetadata[]>([])
  const [selectedSearchResult, setSelectedSearchResult] = useState<ExtensionMetadata | null>(null)
  const [marketplaceSearchResult, setMarketplaceSearchResult] = useState<ValidationResult | null>(null)
  const [detailsLoading, setDetailsLoading] = useState(false)
  const [suggestions, setSuggestions] = useState<ExtensionMetadata[]>([])
  const [showSuggestions, setShowSuggestions] = useState(false)
  const [suggestionIndex, setSuggestionIndex] = useState(-1)
  
  // Install+Sync dialog state
  const [showInstallSyncDialog, setShowInstallSyncDialog] = useState(false)
  const [installSyncTargets, setInstallSyncTargets] = useState<string[]>([])
  const [installSyncInProgress, setInstallSyncInProgress] = useState(false)
  
  // Audit-specific state
  const [auditLoading, setAuditLoading] = useState(false)
  const auditCancelRef = useRef(false)

  useEffect(() => {
    loadDefaultPath()
    loadEditorData()
  }, [])

  const loadEditorData = async () => {
    try {
      const profiles = await GetEditorProfiles()
      setEditorProfiles(profiles || [])
      
      const statuses = await GetAllEditorStatuses()
      setEditorStatuses(statuses || [])
      
      const cli = await GetCLIStatus()
      setCLIStatus(cli)
      
      // Set install target to first available CLI
      if (cli?.vscodeAvailable) {
        setInstallTargetEditor('vscode')
      } else if (cli?.insidersAvailable) {
        setInstallTargetEditor('vscode-insiders')
      } else if (cli?.codiumAvailable) {
        setInstallTargetEditor('vscodium')
      }
    } catch (error) {
      console.error('[Frontend] Failed to load editor data:', error)
    }
  }

  const loadDefaultPath = async () => {
    console.log('[Frontend] Loading default extensions path...')
    try {
      const path = await GetDefaultExtensionsPath()
      console.log('[Frontend] Default path:', path)
      setExtensionsPath(path)
      loadExtensions(path)
    } catch (error) {
      console.error('[Frontend] Failed to get default path:', error)
      const errorMessage = String(error)
      if (errorMessage.includes('could not find VS Code extensions directory')) {
        setError('VS Code extensions directory not found. Please click "Change Path" to select your extensions folder manually.')
      } else {
        setError(`Failed to get default path: ${error}`)
      }
    }
  }

  const loadExtensions = async (path?: string) => {
    const targetPath = path || extensionsPath
    console.log('[Frontend] Loading extensions from:', targetPath)
    setLoading(true)
    setError(null)
    try {
      const exts = await GetInstalledExtensions(targetPath)
      console.log('[Frontend] Loaded extensions:', exts?.length || 0)
      setExtensions(exts || [])
      if (!exts || exts.length === 0) {
        setError('No extensions found in the specified directory')
      }
    } catch (error) {
      console.error('[Frontend] Failed to load extensions:', error)
      setError(`Failed to load extensions: ${error}`)
    } finally {
      setLoading(false)
    }
  }

  const handleAudit = async () => {
    console.log('[Frontend] Starting audit for path:', extensionsPath)
    auditCancelRef.current = false
    setAuditLoading(true)
    setError(null)
    try {
      const report = await AuditAllExtensions(extensionsPath)
      if (auditCancelRef.current) return
      console.log('[Frontend] Audit complete:', report)
      setAuditReport(report)
      // Cache the report for this editor
      setAuditReportsCache(prev => ({ ...prev, [selectedEditor]: report }))
    } catch (error) {
      if (auditCancelRef.current) return
      console.error('[Frontend] Failed to audit extensions:', error)
      setError(`Failed to audit extensions: ${error}`)
    } finally {
      if (!auditCancelRef.current) setAuditLoading(false)
    }
  }

  const handleCancelAudit = () => {
    auditCancelRef.current = true
    setAuditLoading(false)
  }

  const handleValidate = async (extensionId: string) => {
    console.log('[Frontend] Validating extension:', extensionId)
    setLoading(true)
    setError(null)
    try {
      const result = await ValidateExtension(extensionId)
      console.log('[Frontend] Validation result:', result)
      setValidationResult(result)
    } catch (error) {
      console.error('[Frontend] Failed to validate extension:', error)
      setError(`Failed to validate extension: ${error}`)
    } finally {
      setLoading(false)
    }
  }

  const handleMarketplaceSearch = async () => {
    if (!marketplaceSearchQuery.trim()) {
      setError('Please enter a search term (e.g., python, prettier, eslint)')
      return
    }
    console.log('[Frontend] Searching marketplace for:', marketplaceSearchQuery)
    // Suppress any pending or future suggestion displays
    suppressSuggestionsRef.current = true
    // Clear debounce timeout to prevent race condition with autocomplete
    if (debounceTimeoutRef.current) {
      clearTimeout(debounceTimeoutRef.current)
      debounceTimeoutRef.current = null
    }
    setShowSuggestions(false)
    setSuggestions([])
    setLoading(true)
    setError(null)
    setMarketplaceSearchResult(null)
    setSelectedSearchResult(null)
    try {
      const results = await SearchMarketplace(marketplaceSearchQuery.trim())
      console.log('[Frontend] Marketplace search results:', results)
      setMarketplaceSearchResults(results || [])
      setView('search')
      if (!results || results.length === 0) {
        setError(`No extensions found matching "${marketplaceSearchQuery}"`)
      }
    } catch (error) {
      console.error('[Frontend] Failed to search marketplace:', error)
      setError(`Failed to search marketplace: ${error}`)
    } finally {
      setLoading(false)
    }
  }

  const handleValidateSearchResult = async (extensionId: string) => {
    console.log('[Frontend] Validating search result:', extensionId)
    setDetailsLoading(true)
    setMarketplaceSearchResult(null)
    setError(null)
    try {
      const result = await SearchMarketplaceExtension(extensionId)
      console.log('[Frontend] Validation result:', result)
      setMarketplaceSearchResult(result)
    } catch (error) {
      console.error('[Frontend] Failed to validate:', error)
      setError(`Failed to validate: ${error}`)
    } finally {
      setDetailsLoading(false)
    }
  }

  // Debounced search for autocomplete suggestions
  const debounceTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null)
  const suppressSuggestionsRef = useRef(false)
  
  const handleSearchInputChange = (value: string) => {
    setMarketplaceSearchQuery(value)
    setSuggestionIndex(-1)
    
    // Clear previous timeout
    if (debounceTimeoutRef.current) {
      clearTimeout(debounceTimeoutRef.current)
    }
    
    // Don't search if query is too short
    if (value.trim().length < 2) {
      setSuggestions([])
      setShowSuggestions(false)
      return
    }
    
    // Reset suppress flag when user types
    suppressSuggestionsRef.current = false
    
    // Debounce search
    debounceTimeoutRef.current = setTimeout(async () => {
      try {
        const results = await SearchMarketplace(value.trim())
        // Only show suggestions if not suppressed (e.g., user clicked search button)
        if (!suppressSuggestionsRef.current) {
          setSuggestions(results?.slice(0, 8) || [])
          setShowSuggestions(true)
        }
      } catch (error) {
        console.error('[Frontend] Autocomplete search failed:', error)
        setSuggestions([])
      }
    }, 300)
  }

  const handleSuggestionSelect = (result: ExtensionMetadata) => {
    setMarketplaceSearchQuery(result.id)
    setShowSuggestions(false)
    setSuggestions([])
    setSelectedSearchResult(result)
    setMarketplaceSearchResult(null)
    setView('search')
    handleValidateSearchResult(result.id)
  }

  const handleSearchKeyDown = (e: React.KeyboardEvent) => {
    if (!showSuggestions || suggestions.length === 0) {
      if (e.key === 'Enter') {
        handleMarketplaceSearch()
      }
      return
    }
    
    switch (e.key) {
      case 'ArrowDown':
        e.preventDefault()
        setSuggestionIndex(prev => Math.min(prev + 1, suggestions.length - 1))
        break
      case 'ArrowUp':
        e.preventDefault()
        setSuggestionIndex(prev => Math.max(prev - 1, -1))
        break
      case 'Enter':
        e.preventDefault()
        if (suggestionIndex >= 0 && suggestionIndex < suggestions.length) {
          handleSuggestionSelect(suggestions[suggestionIndex])
        } else {
          setShowSuggestions(false)
          handleMarketplaceSearch()
        }
        break
      case 'Escape':
        setShowSuggestions(false)
        setSuggestionIndex(-1)
        break
    }
  }

  const handleDownload = async (extensionId: string) => {
    console.log('[Frontend] Downloading extension:', extensionId)
    setLoading(true)
    setError(null)
    try {
      await DownloadOfficialExtension(extensionId)
      alert(`Downloaded official version of ${extensionId}`)
    } catch (error) {
      console.error('[Frontend] Failed to download extension:', error)
      const errorMsg = `Failed to download: ${error}`
      setError(errorMsg)
      alert(errorMsg)
    } finally {
      setLoading(false)
    }
  }

  const handleSelectDirectory = async () => {
    try {
      const path = await SelectDirectory()
      if (path) {
        setExtensionsPath(path)
        setTempPath(path)
        setIsEditingPath(false)
        loadExtensions(path)
      }
    } catch (error) {
      console.error('Failed to select directory:', error)
    }
  }

  const handleEditPath = () => {
    setTempPath(extensionsPath)
    setIsEditingPath(true)
  }

  const handleSavePath = () => {
    if (tempPath.trim()) {
      setExtensionsPath(tempPath.trim())
      setIsEditingPath(false)
      loadExtensions(tempPath.trim())
    }
  }

  const handleCancelEdit = () => {
    setTempPath(extensionsPath)
    setIsEditingPath(false)
  }

  // Editor selection handler
  const handleEditorChange = async (editorId: string) => {
    console.log('[Frontend] Switching to editor:', editorId)
    setSelectedEditor(editorId)
    setLoading(true)
    setError(null)
    // Restore cached audit report for this editor (if any) or clear it
    setAuditReport(auditReportsCache[editorId] || null)
    setAuditLoading(false)
    auditCancelRef.current = true
    try {
      const exts = await GetEditorExtensions(editorId)
      console.log('[Frontend] Loaded extensions for editor:', exts?.length || 0)
      setExtensions(exts || [])
      
      // Update path display from editor status
      const status = editorStatuses.find(s => s.editor.id === editorId)
      if (status) {
        setExtensionsPath(status.editor.extensionsDir)
      }
    } catch (error) {
      console.error('[Frontend] Failed to load extensions for editor:', error)
      setError(`Failed to load extensions: ${error}`)
      setExtensions([])
    } finally {
      setLoading(false)
    }
  }

  // Install extension via CLI handler
  const handleInstallViaCLI = async (extensionId: string) => {
    if (!cliStatus?.anyAvailable) {
      setError('No VS Code CLI available. Please install VS Code and ensure "code" command is in PATH.')
      return
    }
    
    const cliCommand = getCLICommand(installTargetEditor)
    if (!cliCommand) {
      setError(`CLI not available for ${installTargetEditor}`)
      return
    }
    
    console.log('[Frontend] Installing extension via CLI:', extensionId, 'using', cliCommand)
    setInstalling(true)
    setError(null)
    try {
      await InstallExtensionViaCLI(cliCommand, extensionId)
      alert(`Successfully installed ${extensionId} to ${installTargetEditor}`)
    } catch (error) {
      console.error('[Frontend] Failed to install extension:', error)
      setError(`Failed to install: ${error}`)
    } finally {
      setInstalling(false)
    }
  }

  // Get CLI command for an editor
  const getCLICommand = (editorId: string): string | null => {
    if (!cliStatus) return null
    switch (editorId) {
      case 'vscode': return cliStatus.vscodeAvailable ? 'code' : null
      case 'vscode-insiders': return cliStatus.insidersAvailable ? 'code-insiders' : null
      case 'vscodium': return cliStatus.codiumAvailable ? 'codium' : null
      default: return null
    }
  }

  // Check if an editor has CLI available
  const isEditorCLIAvailable = (editorId: string): boolean => {
    return getCLICommand(editorId) !== null
  }

  // Open Install+Sync dialog
  const handleOpenInstallSyncDialog = () => {
    setInstallSyncTargets([])
    setShowInstallSyncDialog(true)
  }

  // Toggle target in Install+Sync dialog
  const handleToggleInstallSyncTarget = (editorId: string) => {
    setInstallSyncTargets(prev => 
      prev.includes(editorId) 
        ? prev.filter(id => id !== editorId)
        : [...prev, editorId]
    )
  }

  // Execute Install + Sync workflow
  const handleInstallAndSync = async () => {
    if (!selectedSearchResult) return
    if (!cliStatus?.anyAvailable) {
      setError('No VS Code CLI available')
      return
    }
    if (installSyncTargets.length === 0) {
      setError('Please select at least one target editor')
      return
    }

    const cliCommand = getCLICommand(installTargetEditor)
    if (!cliCommand) {
      setError(`CLI not available for ${installTargetEditor}`)
      return
    }

    setInstallSyncInProgress(true)
    setError(null)

    try {
      // Step 1: Install via CLI
      console.log('[Frontend] Install+Sync: Installing', selectedSearchResult.id, 'via', cliCommand)
      await InstallExtensionViaCLI(cliCommand, selectedSearchResult.id)

      // Step 2: Sync to selected targets
      console.log('[Frontend] Install+Sync: Syncing to', installSyncTargets)
      const report = await SyncExtensions(
        installTargetEditor,
        installSyncTargets,
        [selectedSearchResult.id],
        true // overwrite conflicts for this shortcut workflow
      )

      setShowInstallSyncDialog(false)
      
      if (report) {
        const totalCopied = report.totalCopied || 0
        const totalErrors = report.totalErrors || 0
        if (totalErrors > 0) {
          setError(`Installed and synced with ${totalErrors} errors`)
        } else {
          alert(`Successfully installed ${selectedSearchResult.id} and synced to ${totalCopied} target(s)`)
        }
      }
    } catch (error) {
      console.error('[Frontend] Install+Sync failed:', error)
      setError(`Install+Sync failed: ${error}`)
    } finally {
      setInstallSyncInProgress(false)
    }
  }

  // Sync handlers
  const handleSyncToggleExtension = (extId: string) => {
    setSyncSelectedExtensions(prev => 
      prev.includes(extId) 
        ? prev.filter(id => id !== extId)
        : [...prev, extId]
    )
  }

  const handleSyncSelectAll = (filteredExts: Extension[]) => {
    const visibleIds = filteredExts.map(e => e.id)
    setSyncSelectedExtensions(prev => [...new Set([...prev, ...visibleIds])])
  }

  const handleSyncClearAll = () => {
    setSyncSelectedExtensions([])
  }

  // Select all extensions that are missing from all target editors
  const handleSelectMissing = (filteredExts: Extension[]) => {
    const missingIds = filteredExts
      .filter((ext: Extension) => {
        const extIdLower = ext.id.toLowerCase()
        return !syncTargetEditors.some((editorId: string) => 
          targetEditorExtensions[editorId]?.includes(extIdLower)
        )
      })
      .map((ext: Extension) => ext.id)
    setSyncSelectedExtensions(prev => [...new Set([...prev, ...missingIds])])
  }

  // Select all extensions that exist in at least one target editor
  const handleSelectPresent = (filteredExts: Extension[]) => {
    const presentIds = filteredExts
      .filter((ext: Extension) => {
        const extIdLower = ext.id.toLowerCase()
        return syncTargetEditors.some((editorId: string) => 
          targetEditorExtensions[editorId]?.includes(extIdLower)
        )
      })
      .map((ext: Extension) => ext.id)
    setSyncSelectedExtensions(prev => [...new Set([...prev, ...presentIds])])
  }

  const handleSyncToggleTarget = async (editorId: string) => {
    const isRemoving = syncTargetEditors.includes(editorId)
    setSyncTargetEditors(prev => 
      isRemoving 
        ? prev.filter(id => id !== editorId)
        : [...prev, editorId]
    )
    
    // Load extensions from target editor if adding
    if (!isRemoving && !targetEditorExtensions[editorId]) {
      try {
        const exts = await GetEditorExtensions(editorId)
        const extIds = (exts || []).map((e: Extension) => e.id.toLowerCase())
        setTargetEditorExtensions(prev => ({ ...prev, [editorId]: extIds }))
      } catch (error) {
        console.error('[Frontend] Failed to load target editor extensions:', error)
      }
    }
  }

  const handleStartSync = async () => {
    if (syncSelectedExtensions.length === 0) {
      setError('Please select at least one extension to sync')
      return
    }
    if (syncTargetEditors.length === 0) {
      setError('Please select at least one target editor')
      return
    }

    // Check for conflicts first
    let allConflicts: string[] = []
    for (const target of syncTargetEditors) {
      try {
        const conflicts = await DetectSyncConflicts(syncSourceEditor, target, syncSelectedExtensions)
        allConflicts = [...new Set([...allConflicts, ...(conflicts || [])])]
      } catch (error) {
        console.error('[Frontend] Failed to detect conflicts:', error)
      }
    }

    if (allConflicts.length > 0) {
      setSyncConflicts(allConflicts)
      setShowConflictDialog(true)
    } else {
      await executeSync(false)
    }
  }

  const executeSync = async (overwriteConflicts: boolean) => {
    setShowConflictDialog(false)
    setLoading(true)
    setError(null)
    try {
      const report = await SyncExtensions(syncSourceEditor, syncTargetEditors, syncSelectedExtensions, overwriteConflicts)
      console.log('[Frontend] Sync complete:', report)
      setSyncReport(report)
      
      // Show summary first
      const totalCopied = report?.totalCopied || 0
      const totalErrors = report?.totalErrors || 0

      // Display the report briefly
      setSyncReport(report)

      if (totalErrors > 0) {
        setError(`Sync completed with ${totalErrors} errors. ${totalCopied} extensions copied.`)
      } else {
        alert(`Sync complete! ${totalCopied} extensions copied to ${syncTargetEditors.length} editor(s).`)
      }

      // Reset sync UI to initial state after showing the report
      setTimeout(() => {
        setTargetEditorExtensions({})
        setSyncTargetEditors([])
        setSyncSelectedExtensions([])
        setSyncConflicts([])
        setSyncReport(null)
      }, 100)

    } catch (error) {
      console.error('[Frontend] Failed to sync extensions:', error)
      setError(`Failed to sync: ${error}`)
    } finally {
      setLoading(false)
    }
  }

  const filteredExtensions = extensions.filter(ext => 
    ext.id.toLowerCase().includes(searchQuery.toLowerCase()) ||
    ext.name.toLowerCase().includes(searchQuery.toLowerCase())
  )

  const getTrustIcon = (trustLevel: string) => {
    switch (trustLevel) {
      case 'Legitimate':
        return <CheckCircle className="w-5 h-5 text-trust-legitimate" />
      case 'Suspicious':
        return <AlertTriangle className="w-5 h-5 text-trust-suspicious" />
      case 'Malicious':
        return <XCircle className="w-5 h-5 text-trust-malicious" />
      default:
        return <HelpCircle className="w-5 h-5 text-trust-unknown" />
    }
  }

  const getTrustColor = (trustLevel: string) => {
    switch (trustLevel) {
      case 'Legitimate':
        return 'bg-green-100 border-green-500 text-green-800'
      case 'Suspicious':
        return 'bg-yellow-100 border-yellow-500 text-yellow-800'
      case 'Malicious':
        return 'bg-red-100 border-red-500 text-red-800'
      default:
        return 'bg-gray-100 border-gray-500 text-gray-800'
    }
  }

  return (
    <div className="h-screen flex bg-gray-50">
      {/* Sidebar Navigation */}
      <aside className={`bg-gradient-to-b from-blue-600 to-blue-800 text-white flex flex-col transition-all duration-300 ${navCollapsed ? 'w-16' : 'w-56'}`}>
        {/* Logo */}
        <div className="p-4 border-b border-blue-500">
          <div className="flex items-center space-x-3">
            <Shield className="w-8 h-8 flex-shrink-0" />
            {!navCollapsed && (
              <div className="overflow-hidden">
                <h1 className="text-lg font-bold truncate">Vsynx</h1>
                <p className="text-xs text-blue-200 truncate">Extension Manager</p>
              </div>
            )}
          </div>
        </div>
        
        {/* Navigation Items */}
        <nav className="flex-1 py-4">
          <button
            onClick={() => setView('list')}
            className={`w-full flex items-center space-x-3 px-4 py-3 transition ${view === 'list' ? 'bg-blue-700 border-r-4 border-white' : 'hover:bg-blue-700'}`}
            title="Extensions"
          >
            <Monitor className="w-5 h-5 flex-shrink-0" />
            {!navCollapsed && <span>Extensions</span>}
          </button>
          <button
            onClick={() => setView('search')}
            className={`w-full flex items-center space-x-3 px-4 py-3 transition ${view === 'search' ? 'bg-blue-700 border-r-4 border-white' : 'hover:bg-blue-700'}`}
            title="Marketplace"
          >
            <Search className="w-5 h-5 flex-shrink-0" />
            {!navCollapsed && <span>Marketplace</span>}
          </button>
          <button
            onClick={() => setView('audit')}
            className={`w-full flex items-center space-x-3 px-4 py-3 transition ${view === 'audit' ? 'bg-blue-700 border-r-4 border-white' : 'hover:bg-blue-700'}`}
            title="Audit"
          >
            <Shield className="w-5 h-5 flex-shrink-0" />
            {!navCollapsed && <span>Audit</span>}
          </button>
          <button
            onClick={() => setView('sync')}
            className={`w-full flex items-center space-x-3 px-4 py-3 transition ${view === 'sync' ? 'bg-blue-700 border-r-4 border-white' : 'hover:bg-blue-700'}`}
            title="Sync"
          >
            <ArrowRightLeft className="w-5 h-5 flex-shrink-0" />
            {!navCollapsed && <span>Sync</span>}
          </button>
          <button
            onClick={() => setView('settings')}
            className={`w-full flex items-center space-x-3 px-4 py-3 transition ${view === 'settings' ? 'bg-blue-700 border-r-4 border-white' : 'hover:bg-blue-700'}`}
            title="Settings"
          >
            <HelpCircle className="w-5 h-5 flex-shrink-0" />
            {!navCollapsed && <span>Settings</span>}
          </button>
        </nav>
        
        {/* Collapse Toggle */}
        <button
          onClick={() => setNavCollapsed(!navCollapsed)}
          className="p-4 border-t border-blue-500 hover:bg-blue-700 transition flex items-center justify-center"
          title={navCollapsed ? 'Expand sidebar' : 'Collapse sidebar'}
        >
          <ChevronDown className={`w-5 h-5 transition-transform ${navCollapsed ? '-rotate-90' : 'rotate-90'}`} />
        </button>
      </aside>

      {/* Main Content Area */}
      <div className="flex-1 flex flex-col overflow-hidden">
        {/* Header */}
        <header className="bg-white border-b border-gray-200 px-6 py-3 shadow-sm">
          <div className="flex items-center justify-between">
            <h2 className="text-xl font-semibold text-gray-800">
              {view === 'list' ? 'Extensions' : view === 'search' ? 'Marketplace Search' : view === 'audit' ? 'Security Audit' : view === 'sync' ? 'Sync Extensions' : 'Settings'}
            </h2>
          </div>
        </header>

        {/* Error Banner */}
        {error && (
          <div className="bg-red-50 border-b border-red-200 px-6 py-2">
            <div className="flex items-center justify-between">
              <div className="flex items-center space-x-2">
                <AlertTriangle className="w-4 h-4 text-red-600" />
                <span className="text-sm text-red-800">{error}</span>
              </div>
              <button onClick={() => setError(null)} className="text-red-600 hover:text-red-800">
                <XCircle className="w-4 h-4" />
              </button>
            </div>
          </div>
        )}

        {/* Toolbar */}
        <div className="bg-gray-50 border-b border-gray-200 px-6 py-3">
          <div className="flex flex-wrap items-center gap-4">
            {/* Editor Selector + Path */}
            <div className="flex flex-col gap-2">
              <div className="flex items-center space-x-2">
                <Monitor className="w-4 h-4 text-gray-500" />
                <select
                  value={selectedEditor}
                  onChange={(e) => handleEditorChange(e.target.value)}
                  className="px-3 py-1.5 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 bg-white"
                >
                  {editorStatuses.map((status: EditorStatus) => (
                    <option 
                      key={status.editor.id} 
                      value={status.editor.id}
                      disabled={!status.isAvailable}
                    >
                      {status.editor.name} {status.isAvailable ? `(${status.extensionCount})` : '(not found)'}
                    </option>
                  ))}
                </select>
              </div>
              {/* Path controls - now under editor */}
              <div className="flex items-center space-x-2">
                {isEditingPath ? (
                  <div className="flex items-center space-x-2">
                    <input
                      type="text"
                      value={tempPath}
                      onChange={(e) => setTempPath(e.target.value)}
                      onKeyDown={(e) => {
                        if (e.key === 'Enter') handleSavePath()
                        if (e.key === 'Escape') handleCancelEdit()
                      }}
                      placeholder="Enter extensions path..."
                      className="px-2 py-1 text-sm border border-blue-500 rounded focus:ring-1 focus:ring-blue-500 w-64"
                      autoFocus
                    />
                    <button
                      onClick={handleSavePath}
                      className="px-2 py-1 text-sm bg-green-600 text-white rounded hover:bg-green-700"
                    >
                      <Save className="w-3.5 h-3.5" />
                    </button>
                    <button
                      onClick={handleCancelEdit}
                      className="px-2 py-1 text-sm bg-gray-400 text-white rounded hover:bg-gray-500"
                    >
                      <XCircle className="w-3.5 h-3.5" />
                    </button>
                  </div>
                ) : (
                  <>
                    <span className="text-xs text-gray-500 truncate max-w-[200px]" title={extensionsPath}>
                      {extensionsPath || 'Path not set'}
                    </span>
                    <button
                      onClick={handleEditPath}
                      className="p-1 text-gray-500 hover:text-gray-700 hover:bg-gray-200 rounded"
                      title="Edit path"
                    >
                      <Edit3 className="w-3.5 h-3.5" />
                    </button>
                    <button
                      onClick={handleSelectDirectory}
                      className="p-1 text-blue-600 hover:text-blue-700 hover:bg-blue-100 rounded"
                      title="Browse"
                    >
                      <FolderOpen className="w-3.5 h-3.5" />
                    </button>
                  </>
                )}
              </div>
            </div>

            {/* Search */}
            <div className="relative flex-1 min-w-[200px] max-w-md">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-4 h-4" />
              <input
                type="text"
                placeholder="Filter extensions..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="w-full pl-9 pr-4 py-1.5 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
              />
            </div>

            {/* Refresh */}
            <button
              onClick={() => loadExtensions()}
              disabled={loading}
              className="flex items-center space-x-1.5 px-3 py-1.5 text-sm bg-gray-100 hover:bg-gray-200 rounded-lg transition"
            >
              <RefreshCw className={`w-4 h-4 ${loading ? 'animate-spin' : ''}`} />
              <span>Refresh</span>
            </button>
          </div>
        </div>

        {/* Main Content */}
        <main className="flex-1 overflow-hidden">
        {view === 'list' ? (
          <ExtensionsList
            extensions={filteredExtensions}
            selectedExtension={selectedExtension}
            validationResult={validationResult}
            loading={loading}
            onSelectExtension={setSelectedExtension}
            onValidate={handleValidate}
            onDownload={handleDownload}
            getTrustIcon={getTrustIcon}
            getTrustColor={getTrustColor}
          />
        ) : view === 'search' ? (
          <MarketplaceSearchView
            searchQuery={marketplaceSearchQuery}
            setSearchQuery={handleSearchInputChange}
            searchResults={marketplaceSearchResults}
            selectedResult={selectedSearchResult}
            validationResult={marketplaceSearchResult}
            loading={loading}
            detailsLoading={detailsLoading}
            installing={installing}
            cliStatus={cliStatus}
            installTargetEditor={installTargetEditor}
            setInstallTargetEditor={setInstallTargetEditor}
            suggestions={suggestions}
            showSuggestions={showSuggestions}
            suggestionIndex={suggestionIndex}
            editorProfiles={editorProfiles}
            editorStatuses={editorStatuses}
            showInstallSyncDialog={showInstallSyncDialog}
            installSyncTargets={installSyncTargets}
            installSyncInProgress={installSyncInProgress}
            onSearch={handleMarketplaceSearch}
            onSearchKeyDown={handleSearchKeyDown}
            onSuggestionSelect={handleSuggestionSelect}
            onCloseSuggestions={() => setShowSuggestions(false)}
            onSelectResult={setSelectedSearchResult}
            onClearSelection={() => { setSelectedSearchResult(null); setMarketplaceSearchResult(null); }}
            onValidate={handleValidateSearchResult}
            onInstall={handleInstallViaCLI}
            onOpenInstallSyncDialog={handleOpenInstallSyncDialog}
            onToggleInstallSyncTarget={handleToggleInstallSyncTarget}
            onInstallAndSync={handleInstallAndSync}
            onCloseInstallSyncDialog={() => setShowInstallSyncDialog(false)}
            isEditorCLIAvailable={isEditorCLIAvailable}
            getTrustIcon={getTrustIcon}
            getTrustColor={getTrustColor}
          />
        ) : view === 'sync' ? (
          <SyncView
            editorProfiles={editorProfiles}
            editorStatuses={editorStatuses}
            extensions={extensions}
            syncSourceEditor={syncSourceEditor}
            setSyncSourceEditor={setSyncSourceEditor}
            syncTargetEditors={syncTargetEditors}
            syncSelectedExtensions={syncSelectedExtensions}
            syncReport={syncReport}
            syncConflicts={syncConflicts}
            showConflictDialog={showConflictDialog}
            loading={loading}
            targetEditorExtensions={targetEditorExtensions}
            syncFilterMode={syncFilterMode}
            setSyncFilterMode={setSyncFilterMode}
            syncSearchFilter={syncSearchFilter}
            setSyncSearchFilter={setSyncSearchFilter}
            onToggleExtension={handleSyncToggleExtension}
            onSelectAll={handleSyncSelectAll}
            onClearAll={handleSyncClearAll}
            onSelectMissing={handleSelectMissing}
            onSelectPresent={handleSelectPresent}
            onToggleTarget={handleSyncToggleTarget}
            onStartSync={handleStartSync}
            onConfirmOverwrite={() => executeSync(true)}
            onCancelConflict={() => executeSync(false)}
            onEditorChange={handleEditorChange}
          />
        ) : view === 'audit' ? (
          <AuditView
            report={auditReport}
            loading={auditLoading}
            onStartAudit={handleAudit}
            onCancelAudit={handleCancelAudit}
            getTrustIcon={getTrustIcon}
            getTrustColor={getTrustColor}
          />
        ) : view === 'settings' ? (
          <SettingsView
            setVsynxCliStatus={setVsynxCliStatus}
            cliInstalling={cliInstalling}
            setCliInstalling={setCliInstalling}
          />
        ) : null}
        </main>
      </div>
    </div>
  )
}

// Extension List Component
function ExtensionsList({ 
  extensions, 
  selectedExtension, 
  validationResult, 
  loading,
  onSelectExtension, 
  onValidate, 
  onDownload,
  getTrustIcon,
  getTrustColor,
}: any) {
  return (
    <div className="h-full flex">
      {/* Extensions List */}
      <div className="w-1/2 border-r border-gray-200 overflow-y-auto">
        <div className="p-6">
          <h2 className="text-xl font-semibold mb-4">
            Installed Extensions ({extensions.length})
          </h2>
          <div className="space-y-2">
            {extensions.map((ext: Extension, idx: number) => (
              <div
                key={`${ext.id}-${ext.version}-${idx}`}
                onClick={() => onSelectExtension(ext)}
                className={`p-4 border rounded-lg cursor-pointer transition ${
                  selectedExtension?.id === ext.id
                    ? 'border-blue-500 bg-blue-50'
                    : 'border-gray-200 hover:border-gray-300 bg-white'
                }`}
              >
                <div className="flex items-center justify-between">
                  <div className="flex-1">
                    <h3 className="font-semibold text-gray-900">{ext.id}</h3>
                    <p className="text-sm text-gray-600">Version {ext.version}</p>
                  </div>
                  {validationResult?.extensionId === ext.id && (
                    <div className="ml-4">
                      {getTrustIcon(validationResult.trustLevel)}
                    </div>
                  )}
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Details Panel */}
      <div className="w-1/2 overflow-y-auto bg-gray-50">
        <div className="p-6">
          {selectedExtension ? (
            <>
              <div className="bg-white rounded-lg shadow p-6 mb-6">
                <h2 className="text-2xl font-bold mb-4">{selectedExtension.id}</h2>
                <div className="space-y-2 text-sm">
                  <p><span className="font-semibold">Publisher:</span> {selectedExtension.publisher}</p>
                  <p><span className="font-semibold">Name:</span> {selectedExtension.name}</p>
                  <p><span className="font-semibold">Version:</span> {selectedExtension.version}</p>
                  <p><span className="font-semibold">Status:</span> {selectedExtension.isEnabled ? 'Enabled' : 'Disabled'}</p>
                  <p><span className="font-semibold">Path:</span> <code className="text-xs bg-gray-100 px-2 py-1 rounded">{selectedExtension.path}</code></p>
                </div>
              </div>

              <div className="flex space-x-3 mb-6">
                <button
                  onClick={() => onValidate(selectedExtension.id)}
                  disabled={loading}
                  className="flex-1 flex items-center justify-center space-x-2 px-4 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition disabled:opacity-50"
                >
                  <Shield className="w-5 h-5" />
                  <span>Validate</span>
                </button>
                <button
                  onClick={() => onDownload(selectedExtension.id)}
                  disabled={loading}
                  className="flex-1 flex items-center justify-center space-x-2 px-4 py-3 bg-green-600 text-white rounded-lg hover:bg-green-700 transition disabled:opacity-50"
                >
                  <Download className="w-5 h-5" />
                  <span>Download Official</span>
                </button>
              </div>

              {validationResult?.extensionId === selectedExtension.id && (
                <div className={`border-l-4 p-6 rounded-lg ${getTrustColor(validationResult.trustLevel)}`}>
                  <div className="flex items-center space-x-3 mb-4">
                    {getTrustIcon(validationResult.trustLevel)}
                    <h3 className="text-lg font-bold">Trust Level: {validationResult.trustLevel}</h3>
                  </div>
                  
                  <p className="mb-4">{validationResult.recommendation}</p>

                  {validationResult.differences && validationResult.differences.length > 0 && (
                    <div className="mt-4">
                      <h4 className="font-semibold mb-2">Differences Found:</h4>
                      <ul className="list-disc list-inside space-y-1 text-sm">
                        {validationResult.differences.map((diff: string, idx: number) => (
                          <li key={idx}>{diff}</li>
                        ))}
                      </ul>
                    </div>
                  )}

                  {validationResult.error && (
                    <div className="mt-4 p-3 bg-red-50 border border-red-200 rounded">
                      <p className="text-sm text-red-800">{validationResult.error}</p>
                    </div>
                  )}
                </div>
              )}
            </>
          ) : (
            <div className="flex items-center justify-center h-full text-gray-400">
              <p>Select an extension to view details</p>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

// Audit View Component
function AuditView({ report, loading, onStartAudit, onCancelAudit, getTrustIcon, getTrustColor }: any) {
  // Empty state - no report yet, not loading
  if (!report && !loading) {
    return (
      <div className="h-full flex flex-col items-center justify-center text-gray-500">
        <Shield className="w-16 h-16 text-gray-300 mb-4" />
        <h3 className="text-xl font-semibold text-gray-700 mb-2">Security Audit</h3>
        <p className="text-gray-500 mb-6">Scan all installed extensions for security issues</p>
        <button
          onClick={onStartAudit}
          className="flex items-center space-x-2 px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition font-semibold"
        >
          <Shield className="w-5 h-5" />
          <span>Start Audit</span>
        </button>
      </div>
    )
  }

  // Loading state
  if (loading) {
    return (
      <div className="h-full flex flex-col items-center justify-center">
        <RefreshCw className="w-12 h-12 animate-spin text-blue-600 mb-4" />
        <p className="text-gray-600 mb-4">Auditing extensions...</p>
        <button
          onClick={onCancelAudit}
          className="px-4 py-2 text-gray-600 bg-gray-100 rounded-lg hover:bg-gray-200 transition"
        >
          Cancel
        </button>
      </div>
    )
  }

  return (
    <div className="h-full overflow-y-auto p-6">
      <div className="container mx-auto max-w-6xl">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-2xl font-bold">Extension Audit Report</h2>
          <button
            onClick={onStartAudit}
            className="flex items-center space-x-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition"
          >
            <RefreshCw className="w-4 h-4" />
            <span>Re-audit</span>
          </button>
        </div>

        {/* Summary Cards */}
        <div className="grid grid-cols-4 gap-4 mb-8">
          <div className="bg-white rounded-lg shadow p-6">
            <p className="text-sm text-gray-600 mb-1">Total Extensions</p>
            <p className="text-3xl font-bold text-gray-900">{report.totalExtensions}</p>
          </div>
          <div className="bg-white rounded-lg shadow p-6 border-l-4 border-green-500">
            <p className="text-sm text-gray-600 mb-1">Legitimate</p>
            <p className="text-3xl font-bold text-green-600">{report.legitimateCount}</p>
          </div>
          <div className="bg-white rounded-lg shadow p-6 border-l-4 border-yellow-500">
            <p className="text-sm text-gray-600 mb-1">Suspicious</p>
            <p className="text-3xl font-bold text-yellow-600">{report.suspiciousCount}</p>
          </div>
          <div className="bg-white rounded-lg shadow p-6 border-l-4 border-red-500">
            <p className="text-sm text-gray-600 mb-1">Malicious</p>
            <p className="text-3xl font-bold text-red-600">{report.maliciousCount}</p>
          </div>
        </div>

        {/* Results */}
        <div className="bg-white rounded-lg shadow overflow-hidden">
          <div className="px-6 py-4 bg-gray-50 border-b border-gray-200">
            <h3 className="text-lg font-semibold">Detailed Results</h3>
          </div>
          <div className="divide-y divide-gray-200">
            {report.results.map((result: ValidationResult, idx: number) => (
              <div key={idx} className="p-6">
                <div className="flex items-start justify-between">
                  <div className="flex-1">
                    <div className="flex items-center space-x-3 mb-2">
                      {getTrustIcon(result.trustLevel)}
                      <h4 className="text-lg font-semibold">{result.extensionId}</h4>
                      <span className={`px-3 py-1 rounded-full text-xs font-semibold ${getTrustColor(result.trustLevel)}`}>
                        {result.trustLevel}
                      </span>
                    </div>
                    <p className="text-sm text-gray-700 mb-2">{result.recommendation}</p>
                    {result.differences && result.differences.length > 0 && (
                      <ul className="list-disc list-inside text-sm text-gray-600 space-y-1">
                        {result.differences.map((diff: string, i: number) => (
                          <li key={i}>{diff}</li>
                        ))}
                      </ul>
                    )}
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}

// Verified Publisher Badge Component
function VerifiedBadge({ publisher, domain }: { publisher: string, domain?: string }) {
  return (
    <div className="inline-flex items-center space-x-2 px-3 py-1.5 bg-green-50 border-2 border-green-500 rounded-lg">
      <BadgeCheck className="w-5 h-5 text-green-600" />
      <div className="flex flex-col">
        <span className="text-sm font-semibold text-green-800">{publisher}</span>
        {domain && <span className="text-xs text-green-600">{domain}</span>}
      </div>
    </div>
  )
}

// Marketplace Search View Component
function MarketplaceSearchView({ 
  searchQuery, 
  setSearchQuery, 
  searchResults, 
  selectedResult, 
  validationResult, 
  loading,
  detailsLoading,
  installing,
  cliStatus,
  installTargetEditor,
  setInstallTargetEditor,
  suggestions,
  showSuggestions,
  suggestionIndex,
  editorProfiles,
  editorStatuses,
  showInstallSyncDialog,
  installSyncTargets,
  installSyncInProgress,
  onSearch,
  onSearchKeyDown,
  onSuggestionSelect,
  onCloseSuggestions,
  onSelectResult,
  onClearSelection,
  onValidate, 
  onInstall,
  onOpenInstallSyncDialog,
  onToggleInstallSyncTarget,
  onInstallAndSync,
  onCloseInstallSyncDialog,
  isEditorCLIAvailable,
  getTrustIcon, 
  getTrustColor 
}: any) {
  const detailsRef = useRef<HTMLDivElement>(null)
  
  // Auto-scroll to details panel when a result is selected
  useEffect(() => {
    if (selectedResult && detailsRef.current) {
      detailsRef.current.scrollIntoView({ behavior: 'smooth', block: 'start' })
    }
  }, [selectedResult])

  return (
    <div className="h-full overflow-y-auto p-6 bg-gray-50">
      <div className="container mx-auto max-w-6xl">
        <div className="bg-white rounded-lg shadow-lg p-6 mb-6">
          <h2 className="text-2xl font-bold mb-4 flex items-center">
            <Search className="w-6 h-6 mr-2 text-blue-600" />
            Search Marketplace Extensions
          </h2>
          <p className="text-gray-600 mb-6">
            Search extensions by keywords, publisher name, or extension name. Click any result to validate and compare against both registries.
          </p>

          <div className="flex space-x-3 relative">
            <div className="relative flex-1">
              <input
                type="text"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                onKeyDown={onSearchKeyDown}
                onBlur={() => setTimeout(onCloseSuggestions, 150)}
                placeholder="Search by keyword (e.g., python, formatter, linter)"
                className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
              {/* Autocomplete Suggestions Dropdown */}
              {showSuggestions && suggestions.length > 0 && (
                <div className="absolute z-50 w-full mt-1 bg-white border border-gray-300 rounded-lg shadow-lg max-h-80 overflow-y-auto">
                  {suggestions.map((suggestion: ExtensionMetadata, idx: number) => (
                    <div
                      key={suggestion.id}
                      className={`px-4 py-3 cursor-pointer border-b border-gray-100 last:border-b-0 ${
                        idx === suggestionIndex ? 'bg-blue-50' : 'hover:bg-gray-50'
                      }`}
                      onMouseDown={() => onSuggestionSelect(suggestion)}
                    >
                      <div className="flex items-center justify-between">
                        <div className="flex-1 min-w-0">
                          <div className="flex items-center gap-2">
                            <span className="font-medium text-gray-900 truncate">
                              {suggestion.displayName || suggestion.name}
                            </span>
                            {suggestion.isVerifiedPublisher && (
                              <BadgeCheck className="w-4 h-4 text-green-600 flex-shrink-0" />
                            )}
                          </div>
                          <p className="text-xs text-gray-500 truncate">{suggestion.id}</p>
                          <p className="text-xs text-gray-400">{suggestion.publisher}</p>
                        </div>
                        <span className="text-xs text-gray-400 ml-2">v{suggestion.version}</span>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
            <button
              onClick={onSearch}
              disabled={loading}
              className="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition disabled:opacity-50 flex items-center space-x-2"
            >
              {loading ? (
                <>
                  <RefreshCw className="w-5 h-5 animate-spin" />
                  <span>Searching...</span>
                </>
              ) : (
                <>
                  <Search className="w-5 h-5" />
                  <span>Search</span>
                </>
              )}
            </button>
          </div>

          <div className="mt-4 text-sm text-gray-500">
            <p className="font-semibold mb-2">Quick searches:</p>
            <div className="flex flex-wrap gap-2">
              <button onClick={() => { setSearchQuery('python'); onSearch(); }} className="px-3 py-1 bg-gray-100 rounded hover:bg-gray-200">python</button>
              <button onClick={() => { setSearchQuery('prettier'); onSearch(); }} className="px-3 py-1 bg-gray-100 rounded hover:bg-gray-200">prettier</button>
              <button onClick={() => { setSearchQuery('eslint'); onSearch(); }} className="px-3 py-1 bg-gray-100 rounded hover:bg-gray-200">eslint</button>
              <button onClick={() => { setSearchQuery('copilot'); onSearch(); }} className="px-3 py-1 bg-gray-100 rounded hover:bg-gray-200">copilot</button>
            </div>
          </div>
        </div>

        {/* Search Results */}
        {searchResults && searchResults.length > 0 && (
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-4 mb-6">
            {searchResults.slice(0, 20).map((result: ExtensionMetadata) => (
              <div
                key={result.id}
                className={`bg-white rounded-lg shadow p-4 cursor-pointer hover:shadow-lg transition ${
                  selectedResult?.id === result.id ? 'ring-2 ring-blue-500' : ''
                }`}
                onClick={() => {
                  onSelectResult(result);
                  onValidate(result.id);
                }}
              >
                <div className="flex items-start justify-between mb-2">
                  <div className="flex-1">
                    <h3 className="font-bold text-lg">{result.displayName || result.name}</h3>
                    <p className="text-sm text-gray-600">{result.id}</p>
                  </div>
                  {result.isVerifiedPublisher && (
                    <BadgeCheck className="w-6 h-6 text-green-600 flex-shrink-0" />
                  )}
                </div>
                <p className="text-sm text-gray-700 mb-2 line-clamp-2">{result.description}</p>
                <div className="flex items-center justify-between text-xs text-gray-500">
                  <span>{result.publisher}</span>
                  <span>v{result.version}</span>
                </div>
              </div>
            ))}
          </div>
        )}

        {/* Extension Details Panel - shows when selectedResult exists */}
        {selectedResult && (
          <div ref={detailsRef} className="bg-white rounded-lg shadow-lg p-6">
            {/* Header with selected extension info */}
            <div className="flex items-center justify-between mb-6 pb-4 border-b">
              <div className="flex items-center space-x-3">
                <div>
                  <h3 className="text-xl font-bold">{selectedResult.displayName || selectedResult.name}</h3>
                  <p className="text-sm text-gray-500">{selectedResult.id}</p>
                  <p className="text-sm text-gray-400">by {selectedResult.publisher} • v{selectedResult.version}</p>
                </div>
              </div>
              {selectedResult.isVerifiedPublisher && (
                <VerifiedBadge 
                  publisher={selectedResult.publisher} 
                  domain={selectedResult.publisherDomain} 
                />
              )}
            </div>

            {selectedResult.description && (
              <p className="text-gray-700 mb-4">{selectedResult.description}</p>
            )}

            {/* Open in Marketplace Link */}
            <div className="flex gap-2 mb-6">
              <a
                href={`https://marketplace.visualstudio.com/items?itemName=${selectedResult.id}`}
                target="_blank"
                rel="noopener noreferrer"
                className="inline-flex items-center px-3 py-1.5 text-sm text-blue-600 bg-blue-50 rounded-lg hover:bg-blue-100 transition"
              >
                <span className="mr-2">📦</span>
                Open in Microsoft Marketplace
              </a>
              {selectedResult.repositoryUrl && (
                <a
                  href={selectedResult.repositoryUrl}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="inline-flex items-center px-3 py-1.5 text-sm text-gray-600 bg-gray-100 rounded-lg hover:bg-gray-200 transition"
                >
                  <span className="mr-2">📁</span>
                  Repository
                </a>
              )}
            </div>

            {/* Loading state while validating */}
            {detailsLoading && (
              <div className="flex items-center justify-center py-8 mb-6 bg-gray-50 rounded-lg">
                <RefreshCw className="w-6 h-6 animate-spin text-blue-600 mr-3" />
                <span className="text-gray-600">Validating extension...</span>
              </div>
            )}

            {/* Validation Result - shows after loading completes */}
            {validationResult && !detailsLoading && (
              <div className={`border-l-4 p-6 rounded-lg mb-6 ${getTrustColor(validationResult.trustLevel)}`}>
                <div className="flex items-center space-x-3 mb-4">
                  {getTrustIcon(validationResult.trustLevel)}
                  <div>
                    <p className="font-semibold">Trust Level: {validationResult.trustLevel}</p>
                  </div>
                </div>

                <p className="mb-4">{validationResult.recommendation}</p>

                {validationResult.differences && validationResult.differences.length > 0 && (
                  <div className="mt-4">
                    <h4 className="font-semibold mb-2">Validation Details:</h4>
                    <ul className="space-y-2">
                      {validationResult.differences.map((diff: string, idx: number) => (
                        <li key={idx} className="text-sm flex items-start">
                          <span className="mr-2">{diff.includes('✓') ? '✓' : diff.includes('⚠') ? '⚠️' : '•'}</span>
                          <span>{diff.replace(/^[✓⚠]\s*/, '')}</span>
                        </li>
                      ))}
                    </ul>
                  </div>
                )}

                {validationResult.shaMismatchDetails && (
                  <div className="mt-4 p-4 bg-red-50 border border-red-200 rounded">
                    <h4 className="font-semibold text-red-800 mb-2">⚠️ SHA256 Mismatch Detected</h4>
                    <p className="text-sm text-red-700">{validationResult.shaMismatchDetails}</p>
                  </div>
                )}

                {validationResult.error && (
                  <div className="mt-4 p-3 bg-yellow-50 border border-yellow-200 rounded">
                    <p className="text-sm text-yellow-800">{validationResult.error}</p>
                  </div>
                )}
              </div>
            )}

            {/* Install to Editor Section - always visible when selectedResult exists */}
            <div className="bg-gray-50 rounded-lg p-4 mb-6">
              <h4 className="font-semibold mb-3 flex items-center">
                <Download className="w-5 h-5 mr-2 text-blue-600" />
                Install Extension
              </h4>
              
              {cliStatus?.anyAvailable ? (
                <div className="space-y-3">
                  <div className="flex flex-wrap items-center gap-3">
                    <div className="flex items-center gap-2">
                      <span className="text-sm text-gray-600">Target:</span>
                      <select
                        value={installTargetEditor}
                        onChange={(e) => setInstallTargetEditor(e.target.value)}
                        className="px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                      >
                        {cliStatus.vscodeAvailable && <option value="vscode">VS Code</option>}
                        {cliStatus.insidersAvailable && <option value="vscode-insiders">VS Code Insiders</option>}
                        {cliStatus.codiumAvailable && <option value="vscodium">VSCodium</option>}
                      </select>
                    </div>
                    <button
                      onClick={() => onInstall(selectedResult.id)}
                      disabled={installing || !isEditorCLIAvailable(installTargetEditor)}
                      className="flex items-center space-x-2 px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                      {installing ? (
                        <>
                          <RefreshCw className="w-4 h-4 animate-spin" />
                          <span>Installing...</span>
                        </>
                      ) : (
                        <>
                          <Download className="w-4 h-4" />
                          <span>Install to {installTargetEditor === 'vscode' ? 'VS Code' : installTargetEditor === 'vscode-insiders' ? 'Insiders' : 'VSCodium'}</span>
                        </>
                      )}
                    </button>
                    <button
                      onClick={onOpenInstallSyncDialog}
                      disabled={installing}
                      className="flex items-center space-x-2 px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition disabled:opacity-50 disabled:cursor-not-allowed"
                      title="Install to VS Code and sync to clone editors in one step"
                    >
                      <ArrowRightLeft className="w-4 h-4" />
                      <span>Install + Sync...</span>
                    </button>
                  </div>
                  <p className="text-xs text-gray-500">
                    Use "Install + Sync" to install and copy to clone editors (Windsurf, Cursor, Kiro) in one step.
                  </p>
                </div>
              ) : (
                <div className="flex items-center gap-2 text-gray-500">
                  <AlertTriangle className="w-5 h-5 text-yellow-500" />
                  <span className="text-sm">
                    VS Code CLI not found. Install VS Code and run "Shell Command: Install 'code' command in PATH" from the command palette.
                  </span>
                </div>
              )}
            </div>

            {/* Source Comparison - only show when validation is complete */}
            {validationResult && !detailsLoading && (
              <div className="grid grid-cols-2 gap-4">
                {validationResult.marketplaceData && (
                <div className="border rounded-lg p-4">
                  <h4 className="font-semibold mb-2 flex items-center">
                    <span className="mr-2">📦</span>
                    Microsoft Marketplace
                  </h4>
                  <dl className="text-sm space-y-1">
                    <div><dt className="inline font-medium">Version:</dt> <dd className="inline">{validationResult.marketplaceData.version}</dd></div>
                    <div><dt className="inline font-medium">Publisher:</dt> <dd className="inline">{validationResult.marketplaceData.publisher}</dd></div>
                    {validationResult.marketplaceData.displayName && (
                      <div><dt className="inline font-medium">Name:</dt> <dd className="inline">{validationResult.marketplaceData.displayName}</dd></div>
                    )}
                  </dl>
                </div>
              )}

              {validationResult.openvsxData && (
                <div className="border rounded-lg p-4">
                  <h4 className="font-semibold mb-2 flex items-center">
                    <span className="mr-2">🔓</span>
                    OpenVSX Registry
                  </h4>
                  <dl className="text-sm space-y-1">
                    <div><dt className="inline font-medium">Version:</dt> <dd className="inline">{validationResult.openvsxData.version}</dd></div>
                    <div><dt className="inline font-medium">Publisher:</dt> <dd className="inline">{validationResult.openvsxData.publisher}</dd></div>
                    {validationResult.openvsxData.displayName && (
                      <div><dt className="inline font-medium">Name:</dt> <dd className="inline">{validationResult.openvsxData.displayName}</dd></div>
                    )}
                  </dl>
                </div>
              )}
              </div>
            )}
          </div>
        )}

        {/* Empty State when no search yet */}
        {!selectedResult && searchResults.length === 0 && !loading && (
          <div className="bg-white rounded-lg shadow-lg p-12 text-center">
            <Search className="w-16 h-16 text-gray-300 mx-auto mb-4" />
            <h3 className="text-xl font-semibold text-gray-700 mb-2">Search for Extensions</h3>
            <p className="text-gray-500">Start typing in the search box above to find extensions from the Microsoft Marketplace.</p>
          </div>
        )}

        {/* No Results State */}
        {!selectedResult && searchResults.length === 0 && searchQuery && !loading && (
          <div className="bg-white rounded-lg shadow-lg p-8 text-center">
            <AlertTriangle className="w-12 h-12 text-yellow-400 mx-auto mb-4" />
            <h3 className="text-lg font-semibold text-gray-700 mb-2">No Extensions Found</h3>
            <p className="text-gray-500">No extensions match "{searchQuery}". Try a different search term.</p>
          </div>
        )}

        {/* Clear Selection Button */}
        {selectedResult && (
          <div className="flex justify-center mt-4">
            <button
              onClick={onClearSelection}
              className="px-4 py-2 text-gray-600 hover:text-gray-800 hover:bg-gray-100 rounded-lg transition"
            >
              ← Back to search results
            </button>
          </div>
        )}

        {/* Install + Sync Dialog */}
        {showInstallSyncDialog && selectedResult && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div className="bg-white rounded-lg shadow-xl p-6 max-w-md w-full mx-4">
              <h3 className="text-xl font-bold mb-4 flex items-center">
                <ArrowRightLeft className="w-6 h-6 mr-2 text-purple-600" />
                Install + Sync
              </h3>
              
              <div className="mb-4 p-3 bg-gray-50 rounded-lg">
                <p className="text-sm text-gray-600">Extension:</p>
                <p className="font-semibold">{selectedResult.displayName || selectedResult.id}</p>
                <p className="text-xs text-gray-500">{selectedResult.id}</p>
              </div>

              <div className="mb-4">
                <p className="text-sm font-medium text-gray-700 mb-2">Install to:</p>
                <div className="p-3 bg-blue-50 rounded-lg">
                  <p className="font-medium text-blue-800">
                    {installTargetEditor === 'vscode' ? 'VS Code' : 
                     installTargetEditor === 'vscode-insiders' ? 'VS Code Insiders' : 'VSCodium'}
                  </p>
                  <p className="text-xs text-blue-600">Source for sync</p>
                </div>
              </div>

              <div className="mb-6">
                <p className="text-sm font-medium text-gray-700 mb-2">Then sync to (select targets):</p>
                <div className="space-y-2 max-h-48 overflow-y-auto">
                  {editorProfiles
                    .filter((p: EditorProfile) => !p.isVSCodeFamily)
                    .map((profile: EditorProfile) => {
                      const status = editorStatuses.find((s: EditorStatus) => s.editor.id === profile.id)
                      const isSelected = installSyncTargets.includes(profile.id)
                      const isAvailable = status?.isAvailable
                      return (
                        <label
                          key={profile.id}
                          className={`flex items-center p-3 rounded-lg border cursor-pointer transition ${
                            isSelected ? 'border-purple-500 bg-purple-50' : 
                            isAvailable ? 'border-gray-200 hover:border-gray-300' : 'border-gray-100 bg-gray-50 opacity-50'
                          }`}
                        >
                          <input
                            type="checkbox"
                            checked={isSelected}
                            onChange={() => onToggleInstallSyncTarget(profile.id)}
                            disabled={!isAvailable}
                            className="w-4 h-4 mr-3 text-purple-600 rounded"
                          />
                          <div className="flex-1">
                            <span className="font-medium">{profile.name}</span>
                            {!isAvailable && (
                              <p className="text-xs text-red-500">{status?.disabledReason || 'Not available'}</p>
                            )}
                          </div>
                        </label>
                      )
                    })}
                </div>
                {installSyncTargets.length === 0 && (
                  <p className="text-sm text-gray-500 mt-2">Select at least one target editor</p>
                )}
              </div>

              <div className="flex justify-end gap-3">
                <button
                  onClick={onCloseInstallSyncDialog}
                  disabled={installSyncInProgress}
                  className="px-4 py-2 text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200 transition disabled:opacity-50"
                >
                  Cancel
                </button>
                <button
                  onClick={onInstallAndSync}
                  disabled={installSyncInProgress || installSyncTargets.length === 0}
                  className="flex items-center space-x-2 px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {installSyncInProgress ? (
                    <>
                      <RefreshCw className="w-4 h-4 animate-spin" />
                      <span>Installing & Syncing...</span>
                    </>
                  ) : (
                    <>
                      <ArrowRightLeft className="w-4 h-4" />
                      <span>Install + Sync</span>
                    </>
                  )}
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}

// Sync View Component
function SyncView({
  editorProfiles,
  editorStatuses,
  extensions,
  syncSourceEditor,
  setSyncSourceEditor,
  syncTargetEditors,
  syncSelectedExtensions,
  syncReport,
  syncConflicts,
  showConflictDialog,
  loading,
  targetEditorExtensions,
  syncFilterMode,
  setSyncFilterMode,
  syncSearchFilter,
  setSyncSearchFilter,
  onToggleExtension,
  onSelectAll,
  onClearAll,
  onSelectMissing,
  onSelectPresent,
  onToggleTarget,
  onStartSync,
  onConfirmOverwrite,
  onCancelConflict,
  onEditorChange,
}: any) {
  const vsCodeFamilyEditors = editorProfiles.filter((p: EditorProfile) => p.isVSCodeFamily)
  const cloneEditors = editorProfiles.filter((p: EditorProfile) => !p.isVSCodeFamily)
  
  // Check if all target editor extensions have been loaded
  const targetsLoaded = syncTargetEditors.length === 0 || 
    syncTargetEditors.every((editorId: string) => targetEditorExtensions[editorId] !== undefined)
  
  // Helper to check if extension is missing from all targets
  const isMissingFromAllTargets = (extId: string) => {
    if (syncTargetEditors.length === 0 || !targetsLoaded) return false
    const extIdLower = extId.toLowerCase()
    return !syncTargetEditors.some((editorId: string) => 
      targetEditorExtensions[editorId]?.includes(extIdLower)
    )
  }
  
  // Helper to check if extension exists in any target
  const existsInAnyTarget = (extId: string) => {
    if (syncTargetEditors.length === 0 || !targetsLoaded) return false
    const extIdLower = extId.toLowerCase()
    return syncTargetEditors.some((editorId: string) => 
      targetEditorExtensions[editorId]?.includes(extIdLower)
    )
  }
  
  // Filter extensions based on search and filter mode
  const filteredExtensions = extensions.filter((ext: Extension) => {
    // Search filter
    if (syncSearchFilter && !ext.id.toLowerCase().includes(syncSearchFilter.toLowerCase())) {
      return false
    }
    // Mode filter
    if (syncFilterMode === 'missing') {
      const isMissing = isMissingFromAllTargets(ext.id)
      if (!isMissing) return false
    }
    if (syncFilterMode === 'present') {
      const isPresent = existsInAnyTarget(ext.id)
      if (!isPresent) return false
    }
    return true
  })
  
  // Count missing and present for display
  const missingCount = extensions.filter((ext: Extension) => isMissingFromAllTargets(ext.id)).length
  const presentCount = extensions.filter((ext: Extension) => existsInAnyTarget(ext.id)).length

  return (
    <div className="h-full overflow-y-auto p-6 bg-gray-50">
      <div className="container mx-auto max-w-6xl">
        <div className="bg-white rounded-lg shadow-lg p-6 mb-6">
          <h2 className="text-2xl font-bold mb-4 flex items-center">
            <ArrowRightLeft className="w-6 h-6 mr-2 text-blue-600" />
            Sync Extensions Between Editors
          </h2>
          <p className="text-gray-600 mb-6">
            Copy extensions from VS Code to other editors like Windsurf, Cursor, or Kiro.
          </p>

          {/* Source Editor Selection */}
          <div className="mb-6">
            <h3 className="font-semibold mb-3 flex items-center">
              <Monitor className="w-5 h-5 mr-2 text-gray-600" />
              Source Editor (copy from)
            </h3>
            <div className="flex flex-wrap gap-2">
              {vsCodeFamilyEditors.map((profile: EditorProfile) => {
                const status = editorStatuses.find((s: EditorStatus) => s.editor.id === profile.id)
                const isAvailable = status?.isAvailable
                return (
                  <button
                    key={profile.id}
                    onClick={() => {
                      setSyncSourceEditor(profile.id)
                      onEditorChange(profile.id)
                    }}
                    disabled={!isAvailable}
                    className={`px-4 py-2 rounded-lg border-2 transition ${
                      syncSourceEditor === profile.id
                        ? 'border-blue-500 bg-blue-50 text-blue-700'
                        : isAvailable
                        ? 'border-gray-300 hover:border-gray-400'
                        : 'border-gray-200 bg-gray-100 text-gray-400 cursor-not-allowed'
                    }`}
                    title={!isAvailable ? status?.disabledReason : profile.extensionsDir}
                  >
                    {profile.name}
                    {status && (
                      <span className="ml-2 text-xs text-gray-500">
                        ({status.extensionCount} exts)
                      </span>
                    )}
                  </button>
                )
              })}
            </div>
          </div>

          {/* Target Editors Selection */}
          <div className="mb-6">
            <h3 className="font-semibold mb-3 flex items-center">
              <Terminal className="w-5 h-5 mr-2 text-gray-600" />
              Target Editors (copy to)
            </h3>
            <div className="flex flex-wrap gap-2">
              {cloneEditors.map((profile: EditorProfile) => {
                const status = editorStatuses.find((s: EditorStatus) => s.editor.id === profile.id)
                const isSelected = syncTargetEditors.includes(profile.id)
                return (
                  <button
                    key={profile.id}
                    onClick={() => onToggleTarget(profile.id)}
                    className={`px-4 py-2 rounded-lg border-2 transition ${
                      isSelected
                        ? 'border-green-500 bg-green-50 text-green-700'
                        : 'border-gray-300 hover:border-gray-400'
                    }`}
                    title={profile.extensionsDir}
                  >
                    <span className="mr-2">{isSelected ? '✓' : '○'}</span>
                    {profile.name}
                    {status?.isAvailable && (
                      <span className="ml-2 text-xs text-gray-500">
                        ({status.extensionCount} exts)
                      </span>
                    )}
                  </button>
                )
              })}
            </div>
            {syncTargetEditors.length === 0 && (
              <p className="text-sm text-gray-500 mt-2">Select at least one target editor</p>
            )}
          </div>
        </div>

          {/* Extensions Selection */}
        <div className="bg-white rounded-lg shadow-lg p-6 mb-6">
          <div className="flex items-center justify-between mb-4">
            <h3 className="font-semibold flex items-center">
              Extensions to Sync ({syncSelectedExtensions.length} selected of {filteredExtensions.length} shown)
            </h3>
          </div>
          
          {/* Search and Filter Controls */}
          <div className="flex flex-wrap items-center gap-3 mb-4">
            {/* Search Input */}
            <div className="relative flex-1 min-w-48">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-4 h-4" />
              <input
                type="text"
                value={syncSearchFilter}
                onChange={(e) => setSyncSearchFilter(e.target.value)}
                placeholder="Filter by extension ID..."
                className="w-full pl-9 pr-3 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500"
              />
            </div>
            
            {/* Filter Mode Toggles */}
            <div className="flex rounded-lg border border-gray-300 overflow-hidden">
              <button
                type="button"
                onClick={() => setSyncFilterMode('all')}
                className={`px-3 py-2 text-sm ${syncFilterMode === 'all' ? 'bg-blue-600 text-white' : 'bg-white text-gray-700 hover:bg-gray-50'}`}
              >
                All ({extensions.length})
              </button>
              <button
                type="button"
                onClick={() => setSyncFilterMode('missing')}
                className={`px-3 py-2 text-sm border-l ${syncFilterMode === 'missing' ? 'bg-blue-600 text-white' : 'bg-white text-gray-700 hover:bg-gray-50'} disabled:opacity-50`}
                disabled={syncTargetEditors.length === 0 || !targetsLoaded}
                title={syncTargetEditors.length === 0 ? 'Select target editors first' : !targetsLoaded ? 'Loading...' : 'Show extensions missing from all targets'}
              >
                Missing ({targetsLoaded ? missingCount : '...'})
              </button>
              <button
                type="button"
                onClick={() => setSyncFilterMode('present')}
                className={`px-3 py-2 text-sm border-l ${syncFilterMode === 'present' ? 'bg-blue-600 text-white' : 'bg-white text-gray-700 hover:bg-gray-50'} disabled:opacity-50`}
                disabled={syncTargetEditors.length === 0 || !targetsLoaded}
                title={syncTargetEditors.length === 0 ? 'Select target editors first' : !targetsLoaded ? 'Loading...' : 'Show extensions already in at least one target'}
              >
                Present ({targetsLoaded ? presentCount : '...'})
              </button>
            </div>
            {syncTargetEditors.length > 0 && !targetsLoaded && (
              <div className="flex items-center text-sm text-blue-600">
                <RefreshCw className="w-3 h-3 animate-spin mr-1" />
                Loading target extensions...
              </div>
            )}
          </div>
          
          {/* Bulk Action Buttons */}
          <div className="flex flex-wrap gap-2 mb-4">
            <button
              onClick={() => onSelectAll(filteredExtensions)}
              className="px-3 py-1.5 text-sm bg-blue-100 text-blue-700 rounded hover:bg-blue-200"
            >
              Select All Visible
            </button>
            <button
              onClick={() => onSelectMissing(filteredExtensions)}
              disabled={syncTargetEditors.length === 0 || !targetsLoaded}
              className="px-3 py-1.5 text-sm bg-green-100 text-green-700 rounded hover:bg-green-200 disabled:opacity-50 disabled:cursor-not-allowed"
              title={!targetsLoaded ? 'Loading...' : 'Select all extensions missing from targets'}
            >
              Select Missing
            </button>
            <button
              onClick={() => onSelectPresent(filteredExtensions)}
              disabled={syncTargetEditors.length === 0 || !targetsLoaded}
              className="px-3 py-1.5 text-sm bg-yellow-100 text-yellow-700 rounded hover:bg-yellow-200 disabled:opacity-50 disabled:cursor-not-allowed"
              title={!targetsLoaded ? 'Loading...' : 'Select all extensions already in targets'}
            >
              Select Present
            </button>
            <button
              onClick={onClearAll}
              className="px-3 py-1.5 text-sm bg-gray-100 text-gray-700 rounded hover:bg-gray-200"
            >
              Clear Selection
            </button>
          </div>

          <div key={`ext-list-${syncFilterMode}-${syncSearchFilter}-${filteredExtensions.length}`} className="max-h-64 overflow-y-auto border rounded-lg">
            {filteredExtensions.length === 0 ? (
              <div className="p-4 text-center text-gray-500">
                {syncFilterMode === 'missing' ? 'No extensions missing from targets' : 
                 syncFilterMode === 'present' ? 'No extensions found in targets' :
                 'No extensions found. Select a source editor with extensions.'}
              </div>
            ) : (
              <div className="divide-y">
                {filteredExtensions.map((ext: Extension, idx: number) => {
                  // Check which target editors already have this extension
                  const extIdLower = ext.id.toLowerCase()
                  const existsInEditors = syncTargetEditors.filter((editorId: string) => 
                    targetEditorExtensions[editorId]?.includes(extIdLower)
                  )
                  const isNewToAll = existsInEditors.length === 0 && syncTargetEditors.length > 0
                  const isSelected = syncSelectedExtensions.includes(ext.id)
                  
                  return (
                    <label
                      key={`${ext.id}-${ext.version}-${idx}`}
                      className={`flex items-center p-3 cursor-pointer transition-colors ${
                        isSelected 
                          ? 'bg-blue-100 border-l-4 border-blue-500' 
                          : existsInEditors.length > 0 
                            ? 'bg-yellow-50 hover:bg-yellow-100' 
                            : 'hover:bg-gray-50'
                      }`}
                    >
                      <input
                        type="checkbox"
                        checked={isSelected}
                        onChange={() => onToggleExtension(ext.id)}
                        className="w-4 h-4 mr-3 text-blue-600 rounded"
                      />
                      <div className="flex-1 flex items-center justify-between">
                        <div>
                          <span className="font-medium">{ext.id}</span>
                          <span className="ml-2 text-sm text-gray-500">v{ext.version}</span>
                        </div>
                        <div className="flex items-center gap-1">
                          {isNewToAll && (
                            <span 
                              className="px-2 py-0.5 text-xs bg-blue-100 text-blue-700 rounded-full"
                              title="Not installed in any selected target editor - will be copied fresh"
                            >
                              Missing
                            </span>
                          )}
                          {existsInEditors.map((editorId: string) => {
                            const profile = editorProfiles.find((p: EditorProfile) => p.id === editorId)
                            return (
                              <span 
                                key={editorId}
                                className="px-2 py-0.5 text-xs bg-yellow-100 text-yellow-700 rounded-full"
                                title={`Already installed in ${profile?.name || editorId} - will be overwritten if synced`}
                              >
                                In {profile?.name || editorId}
                              </span>
                            )
                          })}
                        </div>
                      </div>
                    </label>
                  )
                })}
              </div>
            )}
          </div>
        </div>

        {/* Pre-Sync Preview Summary */}
        {syncSelectedExtensions.length > 0 && syncTargetEditors.length > 0 && (
          <div className="bg-white rounded-lg shadow-lg p-4 mb-6">
            <h4 className="font-semibold mb-3 text-gray-700">Sync Preview</h4>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
              {syncTargetEditors.map((editorId: string) => {
                const profile = editorProfiles.find((p: EditorProfile) => p.id === editorId)
                const targetExts = targetEditorExtensions[editorId] || []
                
                // Count new vs overwrite for selected extensions
                let newCount = 0
                let overwriteCount = 0
                syncSelectedExtensions.forEach((extId: string) => {
                  if (targetExts.includes(extId.toLowerCase())) {
                    overwriteCount++
                  } else {
                    newCount++
                  }
                })
                
                return (
                  <div key={editorId} className="border rounded-lg p-3 bg-gray-50">
                    <p className="font-medium text-sm mb-2">{profile?.name || editorId}</p>
                    <div className="flex gap-3 text-xs">
                      <span className="flex items-center text-blue-600">
                        <span className="w-2 h-2 bg-blue-500 rounded-full mr-1"></span>
                        {newCount} new
                      </span>
                      <span className="flex items-center text-yellow-600">
                        <span className="w-2 h-2 bg-yellow-500 rounded-full mr-1"></span>
                        {overwriteCount} overwrite
                      </span>
                    </div>
                  </div>
                )
              })}
            </div>
          </div>
        )}

        {/* Sync Button */}
        <div className="flex justify-center mb-6">
          <button
            onClick={onStartSync}
            disabled={loading || syncSelectedExtensions.length === 0 || syncTargetEditors.length === 0}
            className="flex items-center space-x-2 px-8 py-4 bg-green-600 text-white rounded-lg hover:bg-green-700 transition disabled:opacity-50 disabled:cursor-not-allowed text-lg font-semibold"
          >
            {loading ? (
              <>
                <RefreshCw className="w-5 h-5 animate-spin" />
                <span>Syncing...</span>
              </>
            ) : (
              <>
                <Play className="w-5 h-5" />
                <span>Sync {syncSelectedExtensions.length} Extensions</span>
              </>
            )}
          </button>
        </div>

        {/* Sync Report */}
        {syncReport && (
          <div className="bg-white rounded-lg shadow-lg p-6 mb-6">
            <h3 className="font-semibold mb-4">Sync Report</h3>
            <div className="grid grid-cols-3 gap-4 mb-4">
              <div className="bg-green-50 p-4 rounded-lg text-center">
                <p className="text-2xl font-bold text-green-600">{syncReport.totalCopied}</p>
                <p className="text-sm text-gray-600">Copied</p>
              </div>
              <div className="bg-yellow-50 p-4 rounded-lg text-center">
                <p className="text-2xl font-bold text-yellow-600">{syncReport.totalSkipped}</p>
                <p className="text-sm text-gray-600">Skipped</p>
              </div>
              <div className="bg-red-50 p-4 rounded-lg text-center">
                <p className="text-2xl font-bold text-red-600">{syncReport.totalErrors}</p>
                <p className="text-sm text-gray-600">Errors</p>
              </div>
            </div>
            
            {syncReport.results?.map((result: SyncResult, idx: number) => (
              <div key={idx} className="border rounded-lg p-4 mb-2">
                <div className="flex items-center justify-between mb-2">
                  <span className="font-medium">{result.targetEditor}</span>
                  <span className={result.success ? 'text-green-600' : 'text-red-600'}>
                    {result.success ? '✓ Success' : '✗ Failed'}
                  </span>
                </div>
                <p className="text-sm text-gray-600">
                  Copied: {result.copiedCount} | Skipped: {result.skippedCount} | Overwritten: {result.overwrittenCount}
                </p>
                {(result.errors?.length ?? 0) > 0 && (
                  <div className="mt-2 text-sm text-red-600">
                    {result.errors?.map((err: string, i: number) => (
                      <p key={i}>• {err}</p>
                    ))}
                  </div>
                )}
              </div>
            ))}
          </div>
        )}

        {/* Conflict Dialog */}
        {showConflictDialog && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div className="bg-white rounded-lg shadow-xl p-6 max-w-md w-full mx-4">
              <h3 className="text-lg font-bold mb-4 flex items-center">
                <AlertTriangle className="w-5 h-5 mr-2 text-yellow-600" />
                Conflicts Detected
              </h3>
              <p className="text-gray-600 mb-4">
                The following {syncConflicts.length} extension(s) already exist in the target editor(s):
              </p>
              <div className="max-h-40 overflow-y-auto bg-gray-50 rounded p-3 mb-4">
                {syncConflicts.map((conflict: string, idx: number) => (
                  <p key={idx} className="text-sm">• {conflict}</p>
                ))}
              </div>
              <div className="flex gap-3">
                <button
                  onClick={onCancelConflict}
                  className="flex-1 px-4 py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300"
                >
                  Skip Conflicts
                </button>
                <button
                  onClick={onConfirmOverwrite}
                  className="flex-1 px-4 py-2 bg-yellow-500 text-white rounded-lg hover:bg-yellow-600"
                >
                  Overwrite All
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}

// Settings View Component
function SettingsView({ 
  setVsynxCliStatus, 
  cliInstalling, 
  setCliInstalling 
}: any) {
  const [localCliStatus, setLocalCliStatus] = useState<any>(null)
  const [statusLoading, setStatusLoading] = useState(true)
  const [actionMessage, setActionMessage] = useState<string | null>(null)

  // Load CLI status on mount
  useEffect(() => {
    loadCliStatus()
  }, [])

  const loadCliStatus = async () => {
    setStatusLoading(true)
    try {
      const status = await GetCLIInstallStatus()
      setLocalCliStatus(status)
      if (setVsynxCliStatus) setVsynxCliStatus(status)
    } catch (err) {
      console.error('Failed to get CLI status:', err)
    }
    setStatusLoading(false)
  }

  const handleInstallCLI = async () => {
    if (setCliInstalling) setCliInstalling(true)
    setActionMessage(null)
    try {
      const result = await InstallCLI()
      setActionMessage(result.message)
      if (result.success) {
        loadCliStatus()
      }
    } catch (err: any) {
      setActionMessage(`Error: ${err.message || err}`)
    }
    if (setCliInstalling) setCliInstalling(false)
  }

  const handleUninstallCLI = async () => {
    if (setCliInstalling) setCliInstalling(true)
    setActionMessage(null)
    try {
      const result = await UninstallCLI()
      setActionMessage(result.message)
      if (result.success) {
        loadCliStatus()
      }
    } catch (err: any) {
      setActionMessage(`Error: ${err.message || err}`)
    }
    if (setCliInstalling) setCliInstalling(false)
  }

  return (
    <div className="h-full overflow-y-auto p-6 bg-gray-50">
      <div className="container mx-auto max-w-4xl">
        {/* About Section */}
        <div className="bg-white rounded-lg shadow-lg p-6 mb-6">
          <div className="flex items-center space-x-4 mb-6">
            <div className="w-16 h-16 bg-gradient-to-br from-blue-500 to-blue-700 rounded-xl flex items-center justify-center">
              <Shield className="w-10 h-10 text-white" />
            </div>
            <div>
              <h2 className="text-2xl font-bold text-gray-800">Vsynx Manager</h2>
              <p className="text-gray-600">Secure VS Code Extension Manager</p>
              <p className="text-sm text-gray-500">Version 1.0.0</p>
            </div>
          </div>
          <p className="text-gray-600 mb-4">
            Vsynx helps you securely manage VS Code extensions across multiple editors.
            Validate extensions against the Microsoft Marketplace, sync extensions between
            editors, and protect your development environment.
          </p>
        </div>

        {/* CLI Installation Section */}
        <div className="bg-white rounded-lg shadow-lg p-6 mb-6">
          <h3 className="text-lg font-bold mb-4 flex items-center">
            <Terminal className="w-5 h-5 mr-2 text-blue-600" />
            Command Line Interface (CLI)
          </h3>
          <p className="text-gray-600 mb-4">
            The vsynx CLI allows you to manage extensions from your terminal.
            Use it for scripting, automation, and quick access to all features.
          </p>

          {statusLoading ? (
            <div className="flex items-center text-gray-500">
              <RefreshCw className="w-4 h-4 animate-spin mr-2" />
              Checking CLI status...
            </div>
          ) : (
            <div className="space-y-4">
              {/* Status */}
              <div className={`p-4 rounded-lg ${localCliStatus?.installed ? 'bg-green-50 border border-green-200' : 'bg-yellow-50 border border-yellow-200'}`}>
                <div className="flex items-center justify-between">
                  <div className="flex items-center space-x-2">
                    {localCliStatus?.installed ? (
                      <>
                        <CheckCircle className="w-5 h-5 text-green-600" />
                        <span className="font-medium text-green-800">CLI Installed</span>
                      </>
                    ) : (
                      <>
                        <AlertTriangle className="w-5 h-5 text-yellow-600" />
                        <span className="font-medium text-yellow-800">CLI Not Installed</span>
                      </>
                    )}
                  </div>
                  {localCliStatus?.version && (
                    <span className="text-sm text-gray-500">{localCliStatus.version}</span>
                  )}
                </div>
                {localCliStatus?.path && (
                  <p className="text-sm text-gray-600 mt-2">
                    Path: <code className="bg-gray-100 px-2 py-0.5 rounded">{localCliStatus.path}</code>
                  </p>
                )}
                {!localCliStatus?.installed && localCliStatus?.instructions && (
                  <p className="text-sm text-gray-600 mt-2">{localCliStatus.instructions}</p>
                )}
              </div>

              {/* Actions */}
              <div className="flex gap-3">
                {!localCliStatus?.installed ? (
                  <button
                    onClick={handleInstallCLI}
                    disabled={cliInstalling || !localCliStatus?.canInstall}
                    className="flex items-center space-x-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    {cliInstalling ? (
                      <>
                        <RefreshCw className="w-4 h-4 animate-spin" />
                        <span>Installing...</span>
                      </>
                    ) : (
                      <>
                        <Download className="w-4 h-4" />
                        <span>Install CLI</span>
                      </>
                    )}
                  </button>
                ) : (
                  <button
                    onClick={handleUninstallCLI}
                    disabled={cliInstalling}
                    className="flex items-center space-x-2 px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 disabled:opacity-50"
                  >
                    {cliInstalling ? (
                      <>
                        <RefreshCw className="w-4 h-4 animate-spin" />
                        <span>Removing...</span>
                      </>
                    ) : (
                      <>
                        <XCircle className="w-4 h-4" />
                        <span>Uninstall CLI</span>
                      </>
                    )}
                  </button>
                )}
                <button
                  onClick={loadCliStatus}
                  disabled={statusLoading}
                  className="flex items-center space-x-2 px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 disabled:opacity-50"
                >
                  <RefreshCw className={`w-4 h-4 ${statusLoading ? 'animate-spin' : ''}`} />
                  <span>Refresh</span>
                </button>
              </div>

              {/* Action message */}
              {actionMessage && (
                <div className="p-3 bg-gray-100 rounded-lg">
                  <pre className="text-sm whitespace-pre-wrap">{actionMessage}</pre>
                </div>
              )}

              {/* CLI Usage Examples */}
              <div className="mt-4 p-4 bg-gray-50 rounded-lg">
                <h4 className="font-medium mb-2">CLI Usage Examples:</h4>
                <div className="space-y-2 text-sm font-mono">
                  <p><code className="bg-gray-200 px-2 py-0.5 rounded">vsynx validate ms-python.python</code></p>
                  <p><code className="bg-gray-200 px-2 py-0.5 rounded">vsynx audit --path ~/.vscode/extensions</code></p>
                  <p><code className="bg-gray-200 px-2 py-0.5 rounded">vsynx sync run --from vscode --to windsurf,cursor --all</code></p>
                  <p><code className="bg-gray-200 px-2 py-0.5 rounded">vsynx marketplace search python</code></p>
                  <p><code className="bg-gray-200 px-2 py-0.5 rounded">vsynx install ms-python.python --editor vscode</code></p>
                </div>
              </div>
            </div>
          )}
        </div>

        {/* Links Section */}
        <div className="bg-white rounded-lg shadow-lg p-6">
          <h3 className="text-lg font-bold mb-4">Resources</h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <a 
              href="https://github.com/vsynx/vsynx" 
              target="_blank" 
              rel="noopener noreferrer"
              className="flex items-center space-x-3 p-3 bg-gray-50 rounded-lg hover:bg-gray-100 transition"
            >
              <span className="text-2xl">📦</span>
              <div>
                <p className="font-medium">GitHub Repository</p>
                <p className="text-sm text-gray-500">Source code and issues</p>
              </div>
            </a>
            <a 
              href="https://vsynx.dev/docs" 
              target="_blank" 
              rel="noopener noreferrer"
              className="flex items-center space-x-3 p-3 bg-gray-50 rounded-lg hover:bg-gray-100 transition"
            >
              <span className="text-2xl">📚</span>
              <div>
                <p className="font-medium">Documentation</p>
                <p className="text-sm text-gray-500">Guides and API reference</p>
              </div>
            </a>
          </div>
        </div>
      </div>
    </div>
  )
}

export default App
