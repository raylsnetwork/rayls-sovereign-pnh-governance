package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/api/core"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/dto"
)

var _ core.TokenMetadataService = (*tokenMetadataService)(nil)

// tokenMetadataService fetches NFT metadata from external URLs
type tokenMetadataService struct {
	client *http.Client
}

// NewTokenMetadataService creates a new token metadata service
func NewTokenMetadataService() core.TokenMetadataService {
	return &tokenMetadataService{
		client: &http.Client{},
	}
}

// GetMetadata fetches token metadata from external URL
func (s *tokenMetadataService) GetMetadata(
	ctx context.Context,
	baseURL, ercId string,
) (*dto.TokenMetadataInfoDto, error) {
	if baseURL == "" {
		return &dto.TokenMetadataInfoDto{}, nil
	}

	url := buildMetadataURL(baseURL, ercId)
	if url == "" {
		return &dto.TokenMetadataInfoDto{}, nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return parseMetadataResponse(body)
}

// buildMetadataURL constructs the full metadata URL from base URL and token ID
func buildMetadataURL(baseURL, tokenId string) string {
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}

	url := baseURL + tokenId

	if strings.Contains(baseURL, "ipfs") {
		url = transformIPFSURL(url)
	}

	return url
}

// parseMetadataResponse attempts to parse metadata from response body
// Supports two common NFT metadata formats
func parseMetadataResponse(body []byte) (*dto.TokenMetadataInfoDto, error) {
	result := &dto.TokenMetadataInfoDto{}

	var tokenInfo dto.TokenMetadataStandartDto
	if err := json.Unmarshal(body, &tokenInfo); err != nil {
		return nil, err
	}

	if tokenInfo.Properties.Image.Description != "" {
		result.ImageUrl = tokenInfo.Properties.Image.Description
		if strings.Contains(result.ImageUrl, "ipfs") {
			result.ImageUrl = transformIPFSURL(result.ImageUrl)
		}
		result.Description = tokenInfo.Properties.Description.Description
		result.Name = tokenInfo.Properties.Name.Description
		return result, nil
	}

	// Try second format (simple with direct fields)
	var tokenInfo2 dto.TokenMetadataStandart2Dto
	if err := json.Unmarshal(body, &tokenInfo2); err != nil {
		return nil, err
	}

	if tokenInfo2.Image != "" {
		result.ImageUrl = tokenInfo2.Image
		if strings.Contains(result.ImageUrl, "ipfs") {
			result.ImageUrl = transformIPFSURL(result.ImageUrl)
		}
		result.Name = tokenInfo2.Name

		if tokenInfo2.Attributes != nil {
			jsonBytes, err := json.Marshal(tokenInfo2.Attributes)
			if err != nil {
				// Log but don't fail - return partial result
				fmt.Println("Error getting Attributes from metadata to JSON:", err)
				return result, nil
			}
			result.Description = string(jsonBytes)
		}
	}

	return result, nil
}

// transformIPFSURL converts "ipfs://<hash>" to "https://ipfs.io/ipfs/<hash>"
func transformIPFSURL(originalURL string) string {
	const baseIPFSURL = "https://ipfs.io/ipfs/"

	parts := strings.Split(originalURL, "://")
	if len(parts) != 2 {
		return ""
	}

	return baseIPFSURL + parts[1]
}
