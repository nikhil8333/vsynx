package models

// EditorType represents the type of code editor
type EditorType string

const (
	EditorVSCode         EditorType = "vscode"
	EditorVSCodeInsiders EditorType = "vscode-insiders"
	EditorVSCodium       EditorType = "vscodium"
	EditorWindsurf       EditorType = "windsurf"
	EditorCursor         EditorType = "cursor"
	EditorKiro           EditorType = "kiro"
	EditorCustom         EditorType = "custom"
)

// EditorProfile represents a code editor with its configuration
type EditorProfile struct {
	ID             EditorType `json:"id"`
	Name           string     `json:"name"`
	ExtensionsDir  string     `json:"extensionsDir"`
	IndexFile      string     `json:"indexFile"`
	CLICommand     string     `json:"cliCommand,omitempty"`
	IsVSCodeFamily bool       `json:"isVSCodeFamily"`
	IsCustom       bool       `json:"isCustom"`
}

// EditorStatus represents the status/availability of an editor
type EditorStatus struct {
	Editor          EditorProfile `json:"editor"`
	DirExists       bool          `json:"dirExists"`
	IndexFileExists bool          `json:"indexFileExists"`
	CLIAvailable    bool          `json:"cliAvailable"`
	CLIPath         string        `json:"cliPath,omitempty"`
	ExtensionCount  int           `json:"extensionCount"`
	DisabledReason  string        `json:"disabledReason,omitempty"`
	IsAvailable     bool          `json:"isAvailable"`
}

// SyncRequest represents a request to sync extensions between editors
type SyncRequest struct {
	SourceEditor       EditorType   `json:"sourceEditor"`
	TargetEditors      []EditorType `json:"targetEditors"`
	ExtensionIDs       []string     `json:"extensionIds"`
	OverwriteConflicts bool         `json:"overwriteConflicts"`
}

// SyncResult represents the result of syncing to a single target editor
type SyncResult struct {
	TargetEditor     EditorType `json:"targetEditor"`
	Success          bool       `json:"success"`
	CopiedCount      int        `json:"copiedCount"`
	SkippedCount     int        `json:"skippedCount"`
	OverwrittenCount int        `json:"overwrittenCount"`
	IndexUpdated     bool       `json:"indexUpdated"`
	Conflicts        []string   `json:"conflicts,omitempty"`
	Errors           []string   `json:"errors,omitempty"`
}

// SyncReport represents the full sync operation report
type SyncReport struct {
	SourceEditor EditorType   `json:"sourceEditor"`
	Results      []SyncResult `json:"results"`
	TotalCopied  int          `json:"totalCopied"`
	TotalSkipped int          `json:"totalSkipped"`
	TotalErrors  int          `json:"totalErrors"`
}

// CLIStatus represents the status of VS Code family CLI tools
type CLIStatus struct {
	VSCodeAvailable   bool   `json:"vscodeAvailable"`
	VSCodePath        string `json:"vscodePath,omitempty"`
	InsidersAvailable bool   `json:"insidersAvailable"`
	InsidersPath      string `json:"insidersPath,omitempty"`
	CodiumAvailable   bool   `json:"codiumAvailable"`
	CodiumPath        string `json:"codiumPath,omitempty"`
	AnyAvailable      bool   `json:"anyAvailable"`
	PreferredCLI      string `json:"preferredCli,omitempty"`
}

// InstallResult represents the result of installing an extension via CLI
type InstallResult struct {
	ExtensionID string `json:"extensionId"`
	Success     bool   `json:"success"`
	Message     string `json:"message,omitempty"`
	Error       string `json:"error,omitempty"`
}

// InstallReport represents the full install operation report
type InstallReport struct {
	TargetEditor EditorType      `json:"targetEditor"`
	CLIUsed      string          `json:"cliUsed"`
	Results      []InstallResult `json:"results"`
	TotalSuccess int             `json:"totalSuccess"`
	TotalFailed  int             `json:"totalFailed"`
}

// ExtensionIndexEntry represents an entry in extensions.json
type ExtensionIndexEntry struct {
	Identifier       ExtensionIdentifier `json:"identifier"`
	Version          string              `json:"version"`
	Location         ExtensionLocation   `json:"location,omitempty"`
	RelativeLocation string              `json:"relativeLocation"`
	Metadata         map[string]any      `json:"metadata,omitempty"`
}

// ExtensionIdentifier represents the identifier block in extensions.json
type ExtensionIdentifier struct {
	ID   string `json:"id"`
	UUID string `json:"uuid,omitempty"`
}

// ExtensionLocation represents the location block in extensions.json
type ExtensionLocation struct {
	Mid    int    `json:"$mid,omitempty"`
	Path   string `json:"path"`
	Scheme string `json:"scheme,omitempty"`
}
