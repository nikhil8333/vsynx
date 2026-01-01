package marketplace

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
}

func TestMarketplaceConstants(t *testing.T) {
	if MarketplaceAPIURL == "" {
		t.Error("MarketplaceAPIURL is empty")
	}

	expectedURL := "https://marketplace.visualstudio.com/_apis/public/gallery/extensionquery"
	if MarketplaceAPIURL != expectedURL {
		t.Errorf("MarketplaceAPIURL = %s, want %s", MarketplaceAPIURL, expectedURL)
	}

	if UserAgent == "" {
		t.Error("UserAgent is empty")
	}
}

func TestFetchMetadataErrorHandling(t *testing.T) {
	client := NewClient()

	// Test with non-existent extension
	_, err := client.FetchMetadata("nonexistent.extension12345xyz")
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

func TestDownloadExtensionInvalidURL(t *testing.T) {
	client := NewClient()

	_, err := client.DownloadExtension("not-a-valid-url")
	if err == nil {
		t.Error("Expected error for invalid URL")
	}
}
