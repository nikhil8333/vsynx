import '@testing-library/jest-dom'
import { vi } from 'vitest'

// Mock the Wails runtime
vi.mock('../wailsjs/runtime/runtime', () => ({
  EventsOn: vi.fn(),
  EventsOff: vi.fn(),
  EventsEmit: vi.fn(),
}))

// Mock the Wails Go bindings
vi.mock('../wailsjs/go/main/App', () => ({
  GetInstalledExtensions: vi.fn().mockResolvedValue([]),
  GetDefaultExtensionsPath: vi.fn().mockResolvedValue('/mock/path'),
  ValidateExtension: vi.fn().mockResolvedValue({
    extensionId: 'test.extension',
    trustLevel: 'Legitimate',
    recommendation: 'Extension is verified',
  }),
  AuditAllExtensions: vi.fn().mockResolvedValue({
    totalExtensions: 0,
    legitimateCount: 0,
    suspiciousCount: 0,
    maliciousCount: 0,
    results: [],
  }),
  SearchMarketplace: vi.fn().mockResolvedValue([]),
  GetEditorProfiles: vi.fn().mockResolvedValue([
    { id: 'vscode', name: 'VS Code', extensionsDir: '/mock/.vscode/extensions' },
    { id: 'windsurf', name: 'Windsurf', extensionsDir: '/mock/.windsurf/extensions' },
  ]),
  GetAllEditorStatuses: vi.fn().mockResolvedValue([
    { editor: { id: 'vscode', name: 'VS Code' }, isAvailable: true, extensionCount: 10 },
    { editor: { id: 'windsurf', name: 'Windsurf' }, isAvailable: true, extensionCount: 5 },
  ]),
  GetCLIStatus: vi.fn().mockResolvedValue({
    anyAvailable: true,
    vscodeAvailable: true,
  }),
  GetEditorExtensions: vi.fn().mockResolvedValue([]),
  SyncExtensions: vi.fn().mockResolvedValue({
    sourceEditor: 'vscode',
    results: [],
  }),
  InstallExtensionViaCLI: vi.fn().mockResolvedValue(null),
}))

// Mock window.matchMedia
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: vi.fn().mockImplementation((query: string) => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: vi.fn(),
    removeListener: vi.fn(),
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn(),
  })),
})
