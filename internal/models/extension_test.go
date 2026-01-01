package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTrustLevelConstants(t *testing.T) {
	tests := []struct {
		name     string
		level    TrustLevel
		expected string
	}{
		{"Legitimate", TrustLevelLegitimate, "Legitimate"},
		{"Suspicious", TrustLevelSuspicious, "Suspicious"},
		{"Malicious", TrustLevelMalicious, "Malicious"},
		{"Unknown", TrustLevelUnknown, "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.level) != tt.expected {
				t.Errorf("TrustLevel %s = %s, want %s", tt.name, tt.level, tt.expected)
			}
		})
	}
}

func TestExtensionMetadataJSON(t *testing.T) {
	now := time.Now()
	metadata := ExtensionMetadata{
		ID:            "test.extension",
		Publisher:     "test",
		Name:          "extension",
		Version:       "1.0.0",
		DisplayName:   "Test Extension",
		Description:   "A test extension",
		RepositoryURL: "https://github.com/test/extension",
		HomepageURL:   "https://test.com",
		SHA256Hash:    "abc123",
		LastUpdated:   now,
		DownloadURL:   "https://example.com/download",
		Source:        "marketplace",
		AdditionalData: map[string]string{
			"key": "value",
		},
	}

	// Test JSON marshaling
	data, err := json.Marshal(metadata)
	if err != nil {
		t.Fatalf("Failed to marshal metadata: %v", err)
	}

	// Test JSON unmarshaling
	var decoded ExtensionMetadata
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal metadata: %v", err)
	}

	// Verify fields
	if decoded.ID != metadata.ID {
		t.Errorf("ID = %s, want %s", decoded.ID, metadata.ID)
	}
	if decoded.Publisher != metadata.Publisher {
		t.Errorf("Publisher = %s, want %s", decoded.Publisher, metadata.Publisher)
	}
	if decoded.Version != metadata.Version {
		t.Errorf("Version = %s, want %s", decoded.Version, metadata.Version)
	}
}

func TestValidationResultJSON(t *testing.T) {
	result := ValidationResult{
		ExtensionID: "test.extension",
		TrustLevel:  TrustLevelLegitimate,
		MarketplaceData: &ExtensionMetadata{
			ID:        "test.extension",
			Publisher: "test",
			Version:   "1.0.0",
		},
		OpenVSXData: &ExtensionMetadata{
			ID:        "test.extension",
			Publisher: "test",
			Version:   "1.0.0",
		},
		Differences:    []string{},
		Recommendation: "Extension is verified",
		ValidationTime: time.Now(),
	}

	// Test JSON marshaling
	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal result: %v", err)
	}

	// Test JSON unmarshaling
	var decoded ValidationResult
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal result: %v", err)
	}

	// Verify fields
	if decoded.ExtensionID != result.ExtensionID {
		t.Errorf("ExtensionID = %s, want %s", decoded.ExtensionID, result.ExtensionID)
	}
	if decoded.TrustLevel != result.TrustLevel {
		t.Errorf("TrustLevel = %s, want %s", decoded.TrustLevel, result.TrustLevel)
	}
}

func TestInstalledExtension(t *testing.T) {
	ext := InstalledExtension{
		ID:           "test.extension",
		Path:         "/path/to/extension",
		Publisher:    "test",
		Name:         "extension",
		Version:      "1.0.0",
		IsEnabled:    true,
		LastModified: time.Now(),
	}

	if ext.ID != "test.extension" {
		t.Errorf("ID = %s, want test.extension", ext.ID)
	}
	if !ext.IsEnabled {
		t.Error("IsEnabled should be true")
	}
}

func TestAuditReport(t *testing.T) {
	report := AuditReport{
		TotalExtensions: 10,
		LegitimateCount: 7,
		SuspiciousCount: 2,
		MaliciousCount:  1,
		UnknownCount:    0,
		Results: []ValidationResult{
			{
				ExtensionID: "test.extension",
				TrustLevel:  TrustLevelLegitimate,
			},
		},
		AuditTime: time.Now(),
	}

	// Verify counts add up
	total := report.LegitimateCount + report.SuspiciousCount + report.MaliciousCount + report.UnknownCount
	if total != report.TotalExtensions {
		t.Errorf("Count mismatch: %d + %d + %d + %d = %d, want %d",
			report.LegitimateCount, report.SuspiciousCount, report.MaliciousCount, report.UnknownCount,
			total, report.TotalExtensions)
	}

	// Test JSON serialization
	data, err := json.Marshal(report)
	if err != nil {
		t.Fatalf("Failed to marshal report: %v", err)
	}

	var decoded AuditReport
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal report: %v", err)
	}

	if decoded.TotalExtensions != report.TotalExtensions {
		t.Errorf("TotalExtensions = %d, want %d", decoded.TotalExtensions, report.TotalExtensions)
	}
}
