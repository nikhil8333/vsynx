package models

import "time"

// TrustLevel represents the trust classification of an extension
type TrustLevel string

const (
	TrustLevelLegitimate TrustLevel = "Legitimate"
	TrustLevelSuspicious TrustLevel = "Suspicious"
	TrustLevelMalicious  TrustLevel = "Malicious"
	TrustLevelUnknown    TrustLevel = "Unknown"
)

// ExtensionMetadata represents metadata for a VS Code extension
type ExtensionMetadata struct {
	ID                 string            `json:"id"`
	Publisher          string            `json:"publisher"`
	PublisherDomain    string            `json:"publisherDomain,omitempty"`
	IsVerifiedPublisher bool             `json:"isVerifiedPublisher"`
	Name               string            `json:"name"`
	Version            string            `json:"version"`
	DisplayName        string            `json:"displayName"`
	Description        string            `json:"description"`
	RepositoryURL      string            `json:"repositoryUrl"`
	HomepageURL        string            `json:"homepageUrl"`
	SHA256Hash         string            `json:"sha256Hash,omitempty"`
	FileSize           int64             `json:"fileSize,omitempty"`
	LastUpdated        time.Time         `json:"lastUpdated"`
	DownloadURL        string            `json:"downloadUrl"`
	Source             string            `json:"source"` // "marketplace" or "openvsx"
	AdditionalData     map[string]string `json:"additionalData,omitempty"`
}

// ValidationResult represents the result of validating an extension
type ValidationResult struct {
	ExtensionID        string               `json:"extensionId"`
	TrustLevel         TrustLevel           `json:"trustLevel"`
	MarketplaceData    *ExtensionMetadata   `json:"marketplaceData,omitempty"`
	OpenVSXData        *ExtensionMetadata   `json:"openvsxData,omitempty"`
	InstalledData      *ExtensionMetadata   `json:"installedData,omitempty"`
	Differences        []string             `json:"differences,omitempty"`
	SHAMatch           bool                 `json:"shaMatch"`
	SHAMismatchDetails string               `json:"shaMismatchDetails,omitempty"`
	Recommendation     string               `json:"recommendation"`
	ValidationTime     time.Time            `json:"validationTime"`
	Error              string               `json:"error,omitempty"`
}

// InstalledExtension represents an extension installed in the editor
type InstalledExtension struct {
	ID           string    `json:"id"`
	Path         string    `json:"path"`
	Publisher    string    `json:"publisher"`
	Name         string    `json:"name"`
	Version      string    `json:"version"`
	IsEnabled    bool      `json:"isEnabled"`
	LastModified time.Time `json:"lastModified"`
}

// AuditReport represents a full audit of all installed extensions
type AuditReport struct {
	TotalExtensions      int                 `json:"totalExtensions"`
	LegitimateCount      int                 `json:"legitimateCount"`
	SuspiciousCount      int                 `json:"suspiciousCount"`
	MaliciousCount       int                 `json:"maliciousCount"`
	UnknownCount         int                 `json:"unknownCount"`
	Results              []ValidationResult  `json:"results"`
	AuditTime            time.Time           `json:"auditTime"`
}
