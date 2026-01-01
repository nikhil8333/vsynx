package openvsx

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/yourusername/secureopenvsx/internal/models"
)

const (
	OpenVSXAPIURL = "https://open-vsx.org/api"
	UserAgent     = "SecureVSX/1.0"
)

// Client handles communication with the OpenVSX registry API
type Client struct {
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a new OpenVSX API client
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: OpenVSXAPIURL,
	}
}

// openVSXExtension represents the response from OpenVSX API
type openVSXExtension struct {
	Namespace   string `json:"namespace"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	Repository  string `json:"repository"`
	Homepage    string `json:"homepage"`
	Files       struct {
		Download string `json:"download"`
	} `json:"files"`
	Timestamp string `json:"timestamp"`
}

// FetchMetadata fetches extension metadata from the OpenVSX registry
func (c *Client) FetchMetadata(extensionID string) (*models.ExtensionMetadata, error) {
	log.Printf("[OpenVSX] Fetching metadata for extension: %s", extensionID)
	// Parse extensionID (format: publisher.name)
	publisher, name, err := parseExtensionID(extensionID)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/%s/%s", c.baseURL, publisher, name)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", UserAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Printf("[OpenVSX] Request failed for %s: %v", extensionID, err)
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()
	log.Printf("[OpenVSX] Response status for %s: %d", extensionID, resp.StatusCode)

	if resp.StatusCode == http.StatusNotFound {
		log.Printf("[OpenVSX] Extension not found: %s", extensionID)
		return nil, fmt.Errorf("extension not found in OpenVSX: %s", extensionID)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("OpenVSX API returned status %d: %s", resp.StatusCode, string(body))
	}

	var ext openVSXExtension
	if err := json.NewDecoder(resp.Body).Decode(&ext); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	lastUpdated, _ := time.Parse(time.RFC3339, ext.Timestamp)

	metadata := &models.ExtensionMetadata{
		ID:            fmt.Sprintf("%s.%s", ext.Namespace, ext.Name),
		Publisher:     ext.Namespace,
		Name:          ext.Name,
		Version:       ext.Version,
		DisplayName:   ext.DisplayName,
		Description:   ext.Description,
		RepositoryURL: ext.Repository,
		HomepageURL:   ext.Homepage,
		DownloadURL:   ext.Files.Download,
		LastUpdated:   lastUpdated,
		Source:        "openvsx",
	}

	log.Printf("[OpenVSX] Successfully fetched metadata for %s (version %s)", extensionID, metadata.Version)
	return metadata, nil
}

// DownloadExtension downloads the VSIX package from OpenVSX
func (c *Client) DownloadExtension(downloadURL string) ([]byte, error) {
	if downloadURL == "" {
		return nil, fmt.Errorf("download URL is empty")
	}

	req, err := http.NewRequest("GET", downloadURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create download request: %w", err)
	}

	req.Header.Set("User-Agent", UserAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download extension: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read download data: %w", err)
	}

	return data, nil
}

// parseExtensionID parses an extension ID into publisher and name
func parseExtensionID(extensionID string) (publisher, name string, err error) {
	for i := 0; i < len(extensionID); i++ {
		if extensionID[i] == '.' {
			return extensionID[:i], extensionID[i+1:], nil
		}
	}
	return "", "", fmt.Errorf("invalid extension ID format: %s (expected publisher.name)", extensionID)
}
