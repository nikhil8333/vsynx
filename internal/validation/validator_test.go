package validation

import (
	"strings"
	"testing"
	"time"

	"github.com/yourusername/secureopenvsx/internal/models"
)

func TestNewValidator(t *testing.T) {
	v := NewValidator()
	if v == nil {
		t.Fatal("NewValidator returned nil")
	}
	if v.marketplaceClient == nil {
		t.Error("marketplaceClient is nil")
	}
	if v.openvsxClient == nil {
		t.Error("openvsxClient is nil")
	}
}

func TestCompareMetadata(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name              string
		marketplace       *models.ExtensionMetadata
		openvsx           *models.ExtensionMetadata
		expectedTrust     models.TrustLevel
		expectedDiffCount int
	}{
		{
			name: "Identical metadata",
			marketplace: &models.ExtensionMetadata{
				ID:            "test.extension",
				Publisher:     "test",
				Name:          "extension",
				Version:       "1.0.0",
				RepositoryURL: "https://github.com/test/extension",
			},
			openvsx: &models.ExtensionMetadata{
				ID:            "test.extension",
				Publisher:     "test",
				Name:          "extension",
				Version:       "1.0.0",
				RepositoryURL: "https://github.com/test/extension",
			},
			expectedTrust:     models.TrustLevelLegitimate,
			expectedDiffCount: 0,
		},
		{
			name: "Different publishers - critical",
			marketplace: &models.ExtensionMetadata{
				ID:        "test.extension",
				Publisher: "test",
				Name:      "extension",
				Version:   "1.0.0",
			},
			openvsx: &models.ExtensionMetadata{
				ID:        "test.extension",
				Publisher: "malicious",
				Name:      "extension",
				Version:   "1.0.0",
			},
			expectedTrust:     models.TrustLevelMalicious,
			expectedDiffCount: 1,
		},
		{
			name: "Different versions - suspicious",
			marketplace: &models.ExtensionMetadata{
				ID:        "test.extension",
				Publisher: "test",
				Name:      "extension",
				Version:   "1.0.0",
			},
			openvsx: &models.ExtensionMetadata{
				ID:        "test.extension",
				Publisher: "test",
				Name:      "extension",
				Version:   "1.0.1",
			},
			expectedTrust:     models.TrustLevelSuspicious,
			expectedDiffCount: 1,
		},
		{
			name: "Different repository URLs - critical",
			marketplace: &models.ExtensionMetadata{
				ID:            "test.extension",
				Publisher:     "test",
				Name:          "extension",
				Version:       "1.0.0",
				RepositoryURL: "https://github.com/test/extension",
			},
			openvsx: &models.ExtensionMetadata{
				ID:            "test.extension",
				Publisher:     "test",
				Name:          "extension",
				Version:       "1.0.0",
				RepositoryURL: "https://github.com/malicious/extension",
			},
			expectedTrust:     models.TrustLevelMalicious,
			expectedDiffCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &models.ValidationResult{
				ExtensionID:    tt.marketplace.ID,
				ValidationTime: time.Now(),
			}

			validator.compareMetadata(result, tt.marketplace, tt.openvsx)

			if result.TrustLevel != tt.expectedTrust {
				t.Errorf("TrustLevel = %s, want %s", result.TrustLevel, tt.expectedTrust)
			}

			if len(result.Differences) != tt.expectedDiffCount {
				t.Errorf("Difference count = %d, want %d. Differences: %v",
					len(result.Differences), tt.expectedDiffCount, result.Differences)
			}

			if result.Recommendation == "" {
				t.Error("Recommendation is empty")
			}
		})
	}
}

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"https://github.com/test/repo", "github.com/test/repo"},
		{"http://github.com/test/repo/", "github.com/test/repo"},
		{"https://github.com/test/repo.git", "github.com/test/repo"},
		{"HTTPS://GITHUB.COM/TEST/REPO", "github.com/test/repo"},
		{"github.com/test/repo", "github.com/test/repo"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := normalizeURL(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeURL(%s) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestContainsCriticalDifference(t *testing.T) {
	tests := []struct {
		name        string
		differences []string
		expected    bool
	}{
		{
			name:        "No differences",
			differences: []string{},
			expected:    false,
		},
		{
			name:        "Publisher mismatch - critical",
			differences: []string{"Publisher mismatch: test vs malicious"},
			expected:    true,
		},
		{
			name:        "Name mismatch - critical",
			differences: []string{"Extension name mismatch: test vs fake"},
			expected:    true,
		},
		{
			name:        "Repository URL mismatch - critical",
			differences: []string{"Repository URL mismatch: url1 vs url2"},
			expected:    true,
		},
		{
			name:        "Version mismatch - not critical",
			differences: []string{"Version mismatch: 1.0.0 vs 1.0.1"},
			expected:    false,
		},
		{
			name:        "Mixed differences with critical",
			differences: []string{"Version mismatch: 1.0.0 vs 1.0.1", "Publisher mismatch: test vs fake"},
			expected:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsCriticalDifference(tt.differences)
			if result != tt.expected {
				t.Errorf("containsCriticalDifference() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestComputeSHA256(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected string
	}{
		{
			name:     "Empty data",
			data:     []byte{},
			expected: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			name:     "Simple string",
			data:     []byte("hello"),
			expected: "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
		},
		{
			name:     "Another string",
			data:     []byte("test data"),
			expected: "916f0027a575074ce72a331777c3478d6513f786a591bd892da1a577bf2335f9",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ComputeSHA256(tt.data)
			if result != tt.expected {
				t.Errorf("ComputeSHA256() = %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestValidateExtensionErrorHandling(t *testing.T) {
	validator := NewValidator()

	// Test with invalid extension ID
	result, err := validator.ValidateExtension("invalid-id-no-dot")
	if err != nil {
		t.Errorf("ValidateExtension should not return error for invalid ID, got: %v", err)
	}

	// Result should indicate errors
	if result.Error == "" {
		t.Error("Expected error message in result.Error")
	}

	if result.TrustLevel != models.TrustLevelUnknown {
		t.Errorf("TrustLevel should be Unknown for failed validation, got: %s", result.TrustLevel)
	}
}

func TestValidateExtensionStructure(t *testing.T) {
	validator := NewValidator()

	// This will fail to connect to real APIs, but tests the structure
	result, err := validator.ValidateExtension("test.extension")
	
	if err != nil {
		t.Errorf("ValidateExtension should not return error, got: %v", err)
	}

	if result == nil {
		t.Fatal("ValidateExtension returned nil result")
	}

	if result.ExtensionID != "test.extension" {
		t.Errorf("ExtensionID = %s, want test.extension", result.ExtensionID)
	}

	if result.ValidationTime.IsZero() {
		t.Error("ValidationTime should be set")
	}

	if result.Recommendation == "" {
		t.Error("Recommendation should not be empty")
	}
}

func TestCompareMetadataCaseInsensitive(t *testing.T) {
	validator := NewValidator()

	marketplace := &models.ExtensionMetadata{
		ID:        "Test.Extension",
		Publisher: "TEST",
		Name:      "EXTENSION",
		Version:   "1.0.0",
	}

	openvsx := &models.ExtensionMetadata{
		ID:        "test.extension",
		Publisher: "test",
		Name:      "extension",
		Version:   "1.0.0",
	}

	result := &models.ValidationResult{
		ExtensionID:    "test.extension",
		ValidationTime: time.Now(),
	}

	validator.compareMetadata(result, marketplace, openvsx)

	// Publisher and name comparison should be case-insensitive
	if result.TrustLevel != models.TrustLevelLegitimate {
		t.Errorf("TrustLevel = %s, want Legitimate (case-insensitive comparison failed)", result.TrustLevel)
	}

	if len(result.Differences) > 0 {
		t.Errorf("Expected no differences for case-insensitive match, got: %v", result.Differences)
	}
}

func TestCompareMetadataURLNormalization(t *testing.T) {
	validator := NewValidator()

	marketplace := &models.ExtensionMetadata{
		ID:            "test.extension",
		Publisher:     "test",
		Name:          "extension",
		Version:       "1.0.0",
		RepositoryURL: "https://github.com/test/repo.git",
	}

	openvsx := &models.ExtensionMetadata{
		ID:            "test.extension",
		Publisher:     "test",
		Name:          "extension",
		Version:       "1.0.0",
		RepositoryURL: "http://github.com/test/repo/",
	}

	result := &models.ValidationResult{
		ExtensionID:    "test.extension",
		ValidationTime: time.Now(),
	}

	validator.compareMetadata(result, marketplace, openvsx)

	// URLs should be normalized and match
	hasRepoMismatch := false
	for _, diff := range result.Differences {
		if strings.Contains(strings.ToLower(diff), "repository") {
			hasRepoMismatch = true
			break
		}
	}

	if hasRepoMismatch {
		t.Errorf("URLs should match after normalization, but got difference: %v", result.Differences)
	}
}
