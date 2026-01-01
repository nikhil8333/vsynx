package openvsx

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient()
	if client == nil {
		t.Fatal("NewClient returned nil")
	}
	if client.httpClient == nil {
		t.Error("httpClient is nil")
	}
	if client.baseURL == "" {
		t.Error("baseURL is empty")
	}
	if client.baseURL != OpenVSXAPIURL {
		t.Errorf("baseURL = %s, want %s", client.baseURL, OpenVSXAPIURL)
	}
}

func TestParseExtensionID(t *testing.T) {
	tests := []struct {
		name              string
		extensionID       string
		expectedPublisher string
		expectedName      string
		expectError       bool
	}{
		{
			name:              "Valid extension ID",
			extensionID:       "ms-python.python",
			expectedPublisher: "ms-python",
			expectedName:      "python",
			expectError:       false,
		},
		{
			name:              "Another valid ID",
			extensionID:       "microsoft.vscode",
			expectedPublisher: "microsoft",
			expectedName:      "vscode",
			expectError:       false,
		},
		{
			name:              "Complex publisher name",
			extensionID:       "company-name.extension-name",
			expectedPublisher: "company-name",
			expectedName:      "extension-name",
			expectError:       false,
		},
		{
			name:        "Invalid - no dot",
			extensionID: "invalid",
			expectError: true,
		},
		{
			name:        "Invalid - empty",
			extensionID: "",
			expectError: true,
		},
		{
			name:              "Multiple dots - takes first",
			extensionID:       "pub.lisher.name",
			expectedPublisher: "pub",
			expectedName:      "lisher.name",
			expectError:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			publisher, name, err := parseExtensionID(tt.extensionID)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if publisher != tt.expectedPublisher {
				t.Errorf("publisher = %s, want %s", publisher, tt.expectedPublisher)
			}

			if name != tt.expectedName {
				t.Errorf("name = %s, want %s", name, tt.expectedName)
			}
		})
	}
}

func TestFetchMetadataErrorHandling(t *testing.T) {
	client := NewClient()

	// Test with invalid extension ID
	_, err := client.FetchMetadata("invalid")
	if err == nil {
		t.Error("Expected error for invalid extension ID")
	}

	// Test with non-existent extension
	_, err = client.FetchMetadata("nonexistent.extension12345")
	if err == nil {
		t.Error("Expected error for non-existent extension")
	}
}

func TestDownloadExtensionEmptyURL(t *testing.T) {
	client := NewClient()

	_, err := client.DownloadExtension("")
	if err == nil {
		t.Error("Expected error for empty download URL")
	}

	expectedMsg := "download URL is empty"
	if err.Error() != expectedMsg {
		t.Errorf("Error message = %s, want %s", err.Error(), expectedMsg)
	}
}
