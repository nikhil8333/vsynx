package validation

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/yourusername/secureopenvsx/internal/marketplace"
	"github.com/yourusername/secureopenvsx/internal/models"
	"github.com/yourusername/secureopenvsx/internal/openvsx"
)

// Validator handles extension validation and trust classification
type Validator struct {
	marketplaceClient *marketplace.Client
	openvsxClient     *openvsx.Client
}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	return &Validator{
		marketplaceClient: marketplace.NewClient(),
		openvsxClient:     openvsx.NewClient(),
	}
}

// ValidateExtension validates an extension by comparing marketplace and OpenVSX metadata
func (v *Validator) ValidateExtension(extensionID string) (*models.ValidationResult, error) {
	log.Printf("[Validator] Starting validation for extension: %s", extensionID)
	result := &models.ValidationResult{
		ExtensionID:    extensionID,
		TrustLevel:     models.TrustLevelUnknown,
		ValidationTime: time.Now(),
	}

	// Fetch from Microsoft Marketplace
	marketplaceData, marketplaceErr := v.marketplaceClient.FetchMetadata(extensionID)
	if marketplaceErr != nil {
		result.Error = fmt.Sprintf("Marketplace error: %v", marketplaceErr)
	} else {
		result.MarketplaceData = marketplaceData
	}

	// Fetch from OpenVSX
	openvsxData, openvsxErr := v.openvsxClient.FetchMetadata(extensionID)
	if openvsxErr != nil {
		if result.Error != "" {
			result.Error += fmt.Sprintf("; OpenVSX error: %v", openvsxErr)
		} else {
			result.Error = fmt.Sprintf("OpenVSX error: %v", openvsxErr)
		}
	} else {
		result.OpenVSXData = openvsxData
	}

	// If both failed, return early
	if marketplaceErr != nil && openvsxErr != nil {
		result.Recommendation = "Cannot validate: both sources unavailable"
		return result, nil
	}

	// If only one source is available, mark as suspicious
	if marketplaceErr != nil {
		result.TrustLevel = models.TrustLevelSuspicious
		result.Differences = append(result.Differences, "Extension not found in Microsoft Marketplace")
		result.Recommendation = "Extension only exists in OpenVSX - verify authenticity manually"
		return result, nil
	}

	if openvsxErr != nil {
		result.TrustLevel = models.TrustLevelLegitimate
		result.Differences = append(result.Differences, "Extension not found in OpenVSX")
		result.Recommendation = "Extension verified from Microsoft Marketplace (OpenVSX unavailable)"
		return result, nil
	}

	// Compare metadata and classify trust level
	v.compareMetadata(result, marketplaceData, openvsxData)

	log.Printf("[Validator] Validation complete for %s: %s", extensionID, result.TrustLevel)
	return result, nil
}

// compareMetadata compares marketplace and OpenVSX metadata and determines trust level
func (v *Validator) compareMetadata(result *models.ValidationResult, marketplace, openvsx *models.ExtensionMetadata) {
	differences := []string{}

	// Add publisher verification info
	if marketplace.IsVerifiedPublisher {
		differences = append(differences, fmt.Sprintf("✓ Verified Publisher: %s", marketplace.Publisher))
		if marketplace.PublisherDomain != "" {
			differences = append(differences, fmt.Sprintf("  Publisher Domain: %s", marketplace.PublisherDomain))
		}
	}

	// Compare SHA256 hashes if available
	if marketplace.SHA256Hash != "" && openvsx.SHA256Hash != "" {
		if marketplace.SHA256Hash == openvsx.SHA256Hash {
			result.SHAMatch = true
			differences = append(differences, "✓ SHA256 hashes match")
		} else {
			result.SHAMatch = false
			result.SHAMismatchDetails = fmt.Sprintf("Marketplace: %s, OpenVSX: %s",
				marketplace.SHA256Hash[:16]+"...", openvsx.SHA256Hash[:16]+"...")
			differences = append(differences, "⚠ SHA256 hash mismatch - binaries are different!")
		}
	} else if marketplace.SHA256Hash != "" || openvsx.SHA256Hash != "" {
		result.SHAMatch = false
		differences = append(differences, "SHA256 hash not available from one source")
	}

	// Compare publisher
	if !strings.EqualFold(marketplace.Publisher, openvsx.Publisher) {
		differences = append(differences, fmt.Sprintf("Publisher mismatch: %s (marketplace) vs %s (OpenVSX)",
			marketplace.Publisher, openvsx.Publisher))
	}

	// Compare name
	if !strings.EqualFold(marketplace.Name, openvsx.Name) {
		differences = append(differences, fmt.Sprintf("Extension name mismatch: %s (marketplace) vs %s (OpenVSX)",
			marketplace.Name, openvsx.Name))
	}

	// Compare version
	if marketplace.Version != openvsx.Version {
		differences = append(differences, fmt.Sprintf("Version mismatch: %s (marketplace) vs %s (OpenVSX)",
			marketplace.Version, openvsx.Version))
	}

	// Compare repository URL
	if marketplace.RepositoryURL != "" && openvsx.RepositoryURL != "" {
		if !strings.EqualFold(normalizeURL(marketplace.RepositoryURL), normalizeURL(openvsx.RepositoryURL)) {
			differences = append(differences, fmt.Sprintf("Repository URL mismatch: %s (marketplace) vs %s (OpenVSX)",
				marketplace.RepositoryURL, openvsx.RepositoryURL))
		}
	}

	result.Differences = differences

	// Determine trust level based on differences and SHA match
	if result.SHAMismatchDetails != "" {
		// SHA mismatch is critical - binaries are different
		result.TrustLevel = models.TrustLevelMalicious
		result.Recommendation = "DANGER: SHA256 mismatch detected - binaries are DIFFERENT. Potential supply chain attack!"
	} else if len(differences) == 0 || (len(differences) <= 2 && marketplace.IsVerifiedPublisher) {
		result.TrustLevel = models.TrustLevelLegitimate
		if marketplace.IsVerifiedPublisher {
			result.Recommendation = "Extension is verified from trusted publisher - metadata matches across sources"
		} else {
			result.Recommendation = "Extension is verified - metadata matches across sources"
		}
	} else if containsCriticalDifference(differences) {
		result.TrustLevel = models.TrustLevelMalicious
		result.Recommendation = "DANGER: Critical metadata mismatches detected - do not use this extension"
	} else {
		result.TrustLevel = models.TrustLevelSuspicious
		result.Recommendation = "Warning: Minor metadata differences detected - verify before use"
	}
}

// containsCriticalDifference checks if differences contain critical mismatches
func containsCriticalDifference(differences []string) bool {
	for _, diff := range differences {
		lower := strings.ToLower(diff)
		// Publisher or name mismatch is critical
		if strings.Contains(lower, "publisher mismatch") || strings.Contains(lower, "name mismatch") {
			return true
		}
		// Repository URL mismatch is critical
		if strings.Contains(lower, "repository url mismatch") {
			return true
		}
	}
	return false
}

// normalizeURL normalizes a URL for comparison
func normalizeURL(url string) string {
	url = strings.ToLower(url)
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimSuffix(url, "/")
	url = strings.TrimSuffix(url, ".git")
	return url
}

// ComputeSHA256 computes the SHA256 hash of data
func ComputeSHA256(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// DownloadOfficialExtension downloads the official extension from Microsoft Marketplace
func (v *Validator) DownloadOfficialExtension(extensionID string) ([]byte, string, error) {
	metadata, err := v.marketplaceClient.FetchMetadata(extensionID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to fetch metadata: %w", err)
	}

	if metadata.DownloadURL == "" {
		return nil, "", fmt.Errorf("no download URL available for extension")
	}

	data, err := v.marketplaceClient.DownloadExtension(metadata.DownloadURL)
	if err != nil {
		return nil, "", fmt.Errorf("failed to download extension: %w", err)
	}

	hash := ComputeSHA256(data)
	return data, hash, nil
}
