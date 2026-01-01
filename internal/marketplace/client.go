package marketplace

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/yourusername/secureopenvsx/internal/models"
)

const (
	MarketplaceAPIURL = "https://marketplace.visualstudio.com/_apis/public/gallery/extensionquery"
	APIVersion        = "7.0-preview.1"
	UserAgent         = "Vsynx/1.0"
)

// Client handles communication with the Microsoft Marketplace API
type Client struct {
	httpClient *http.Client
}

// NewClient creates a new marketplace API client
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// marketplaceQuery represents the request structure for the marketplace API
type marketplaceQuery struct {
	Filters []filter `json:"filters"`
	Flags   int      `json:"flags"`
}

type filter struct {
	Criteria []criterion `json:"criteria"`
}

type criterion struct {
	FilterType int    `json:"filterType"`
	Value      string `json:"value"`
}

// marketplaceResponse represents the response from the marketplace API
type marketplaceResponse struct {
	Results []struct {
		Extensions []struct {
			Publisher struct {
				PublisherName    string `json:"publisherName"`
				PublisherID      string `json:"publisherId"`
				Domain           string `json:"domain"`
				IsDomainVerified bool   `json:"isDomainVerified"`
				Flags            string `json:"flags"`
			} `json:"publisher"`
			ExtensionName    string `json:"extensionName"`
			DisplayName      string `json:"displayName"`
			ShortDescription string `json:"shortDescription"`
			Flags            string `json:"flags"`
			Versions         []struct {
				Version     string `json:"version"`
				LastUpdated string `json:"lastUpdated"`
				Flags       string `json:"flags"`
				Files       []struct {
					AssetType string `json:"assetType"`
					Source    string `json:"source"`
				} `json:"files"`
				Properties []struct {
					Key   string `json:"key"`
					Value string `json:"value"`
				} `json:"properties"`
			} `json:"versions"`
		} `json:"extensions"`
	} `json:"results"`
}

// SearchExtensions searches for extensions using keywords or partial names
func (c *Client) SearchExtensions(searchTerm string) ([]*models.ExtensionMetadata, error) {
	log.Printf("[Marketplace] Searching for extensions matching: %s", searchTerm)

	// Use FilterType 10 for text search (supports keywords/wildcards)
	query := marketplaceQuery{
		Filters: []filter{
			{
				Criteria: []criterion{
					{
						FilterType: 10, // SearchText (keyword search)
						Value:      searchTerm,
					},
				},
			},
		},
		Flags: 0x192, // Include versions, files, and other details
	}

	jsonData, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query: %w", err)
	}

	req, err := http.NewRequest("POST", MarketplaceAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", fmt.Sprintf("application/json; api-version=%s", APIVersion))
	req.Header.Set("User-Agent", UserAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Printf("[Marketplace] Search request failed for %s: %v", searchTerm, err)
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()
	log.Printf("[Marketplace] Search response status for %s: %d", searchTerm, resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("marketplace API returned status %d: %s", resp.StatusCode, string(body))
	}

	var apiResp marketplaceResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(apiResp.Results) == 0 || len(apiResp.Results[0].Extensions) == 0 {
		log.Printf("[Marketplace] No extensions found matching: %s", searchTerm)
		return []*models.ExtensionMetadata{}, nil
	}

	// Convert all results to ExtensionMetadata
	var results []*models.ExtensionMetadata
	for _, ext := range apiResp.Results[0].Extensions {
		if len(ext.Versions) == 0 {
			continue
		}

		latestVersion := ext.Versions[0]
		var downloadURL, repoURL string

		for _, file := range latestVersion.Files {
			if file.AssetType == "Microsoft.VisualStudio.Services.VSIXPackage" {
				downloadURL = file.Source
			}
		}

		for _, prop := range latestVersion.Properties {
			if prop.Key == "Microsoft.VisualStudio.Services.Links.Source" ||
				prop.Key == "Microsoft.VisualStudio.Services.Links.Repository" {
				repoURL = prop.Value
			}
		}

		lastUpdated, _ := time.Parse(time.RFC3339, latestVersion.LastUpdated)
		isVerified := ext.Publisher.IsDomainVerified || ext.Publisher.Flags == "verified"

		metadata := &models.ExtensionMetadata{
			ID:                  fmt.Sprintf("%s.%s", ext.Publisher.PublisherName, ext.ExtensionName),
			Publisher:           ext.Publisher.PublisherName,
			PublisherDomain:     ext.Publisher.Domain,
			IsVerifiedPublisher: isVerified,
			Name:                ext.ExtensionName,
			Version:             latestVersion.Version,
			DisplayName:         ext.DisplayName,
			Description:         ext.ShortDescription,
			RepositoryURL:       repoURL,
			DownloadURL:         downloadURL,
			LastUpdated:         lastUpdated,
			Source:              "marketplace",
		}
		results = append(results, metadata)
	}

	log.Printf("[Marketplace] Found %d extensions matching: %s", len(results), searchTerm)
	return results, nil
}

// FetchMetadata fetches extension metadata from the Microsoft Marketplace
func (c *Client) FetchMetadata(extensionID string) (*models.ExtensionMetadata, error) {
	log.Printf("[Marketplace] Fetching metadata for extension: %s", extensionID)
	query := marketplaceQuery{
		Filters: []filter{
			{
				Criteria: []criterion{
					{
						FilterType: 7, // Extension name (exact match)
						Value:      extensionID,
					},
				},
			},
		},
		Flags: 0x192, // Include versions, files, and other details
	}

	jsonData, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query: %w", err)
	}

	req, err := http.NewRequest("POST", MarketplaceAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", fmt.Sprintf("application/json; api-version=%s", APIVersion))
	req.Header.Set("User-Agent", UserAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Printf("[Marketplace] Request failed for %s: %v", extensionID, err)
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()
	log.Printf("[Marketplace] Response status for %s: %d", extensionID, resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("marketplace API returned status %d: %s", resp.StatusCode, string(body))
	}

	var apiResp marketplaceResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(apiResp.Results) == 0 || len(apiResp.Results[0].Extensions) == 0 {
		log.Printf("[Marketplace] Extension not found: %s", extensionID)
		return nil, fmt.Errorf("extension not found in marketplace: %s", extensionID)
	}

	ext := apiResp.Results[0].Extensions[0]
	if len(ext.Versions) == 0 {
		return nil, fmt.Errorf("no versions found for extension: %s", extensionID)
	}

	latestVersion := ext.Versions[0]

	// Extract download URL and repository URL
	var downloadURL, repoURL string
	for _, file := range latestVersion.Files {
		if file.AssetType == "Microsoft.VisualStudio.Services.VSIXPackage" {
			downloadURL = file.Source
		}
	}

	for _, prop := range latestVersion.Properties {
		if prop.Key == "Microsoft.VisualStudio.Services.Links.Source" ||
			prop.Key == "Microsoft.VisualStudio.Services.Links.Repository" {
			repoURL = prop.Value
		}
	}

	lastUpdated, _ := time.Parse(time.RFC3339, latestVersion.LastUpdated)

	// Check if publisher is verified (domain verified flag)
	isVerified := ext.Publisher.IsDomainVerified || ext.Publisher.Flags == "verified"

	metadata := &models.ExtensionMetadata{
		ID:                  fmt.Sprintf("%s.%s", ext.Publisher.PublisherName, ext.ExtensionName),
		Publisher:           ext.Publisher.PublisherName,
		PublisherDomain:     ext.Publisher.Domain,
		IsVerifiedPublisher: isVerified,
		Name:                ext.ExtensionName,
		Version:             latestVersion.Version,
		DisplayName:         ext.DisplayName,
		Description:         ext.ShortDescription,
		RepositoryURL:       repoURL,
		DownloadURL:         downloadURL,
		LastUpdated:         lastUpdated,
		Source:              "marketplace",
	}

	if isVerified {
		log.Printf("[Marketplace] Successfully fetched metadata for %s (version %s) - VERIFIED PUBLISHER", extensionID, metadata.Version)
	} else {
		log.Printf("[Marketplace] Successfully fetched metadata for %s (version %s)", extensionID, metadata.Version)
	}
	return metadata, nil
}

// DownloadExtension downloads the VSIX package from the marketplace
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
