package scrapers

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"dsmpartsfinder-api/siteclients"

	"github.com/PuerkitoBio/goquery"
)

// KleinanzeigenClient implements scraping for kleinanzeigen.de
type KleinanzeigenClient struct {
	baseURL    string
	httpClient *http.Client
	siteID     int
}

// NewKleinanzeigenClient creates a new Kleinanzeigen scraper client
func NewKleinanzeigenClient(siteID int) *KleinanzeigenClient {
	return &KleinanzeigenClient{
		baseURL: "https://www.kleinanzeigen.de",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		siteID: siteID,
	}
}

// GetName returns the name of the site client
func (c *KleinanzeigenClient) GetName() string {
	return "Kleinanzeigen"
}

// GetSiteID returns the database ID of the site
func (c *KleinanzeigenClient) GetSiteID() int {
	return c.siteID
}

// FetchParts fetches parts from Kleinanzeigen based on search parameters
func (c *KleinanzeigenClient) FetchParts(ctx context.Context, params siteclients.SearchParams) ([]siteclients.Part, error) {
	log.Printf("[KleinanzeigenClient] Starting fetch with params: %+v", params)

	// Build search URL
	searchURL, err := c.buildSearchURL(params)
	if err != nil {
		return nil, fmt.Errorf("failed to build search URL: %w", err)
	}

	log.Printf("[KleinanzeigenClient] Search URL: %s", searchURL)

	// Fetch the page
	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers to mimic a browser
	// Note: Don't set Accept-Encoding - Go's http.Client will handle gzip decompression automatically
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:143.0) Gecko/20100101 Firefox/143.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	log.Printf("[KleinanzeigenClient] Got response with status: %d", resp.StatusCode)

	// Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Save HTML to file for debugging
	htmlFile := "/tmp/kleinanzeigen_response.html"
	if err := os.WriteFile(htmlFile, bodyBytes, 0644); err != nil {
		log.Printf("[KleinanzeigenClient] Warning: failed to save HTML to file: %v", err)
	} else {
		log.Printf("[KleinanzeigenClient] Saved HTML response to: %s (size: %d bytes)", htmlFile, len(bodyBytes))
	}

	// Check for common class names in the HTML
	htmlContent := string(bodyBytes)
	log.Printf("[KleinanzeigenClient] Checking for various selectors:")
	log.Printf("  - Contains 'aditem': %v", strings.Contains(htmlContent, "aditem"))
	log.Printf("  - Contains 'ad-listitem': %v", strings.Contains(htmlContent, "ad-listitem"))
	log.Printf("  - Contains 'article': %v", strings.Contains(htmlContent, "article"))
	log.Printf("  - Contains 'data-adid': %v", strings.Contains(htmlContent, "data-adid"))

	// Parse HTML
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(bodyBytes)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	log.Printf("[KleinanzeigenClient] HTML parsed successfully")

	// Initialize parts slice
	parts := make([]siteclients.Part, 0)

	// Use article.aditem selector - this finds the actual ad listings
	selector := "article.aditem"
	articleCount := doc.Find(selector).Length()
	log.Printf("[KleinanzeigenClient] Found %d elements with selector: %s", articleCount, selector)

	if articleCount == 0 {
		log.Printf("[KleinanzeigenClient] ERROR: No articles found")
		return parts, nil
	}

	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		log.Printf("[KleinanzeigenClient] Processing article %d of %d", i+1, articleCount)
		part, err := c.extractPart(ctx, s)
		if err != nil {
			log.Printf("[KleinanzeigenClient] ERROR: failed to extract part %d: %v", i, err)
			return
		}
		log.Printf("[KleinanzeigenClient] Successfully extracted part %d: ID=%s, Name=%s", i+1, part.ID, part.Name)
		parts = append(parts, part)
	})

	log.Printf("[KleinanzeigenClient] Successfully extracted %d out of %d parts", len(parts), articleCount)

	return parts, nil
}

// buildSearchURL constructs the search URL with parameters
func (c *KleinanzeigenClient) buildSearchURL(params siteclients.SearchParams) (string, error) {
	// Build the search keywords
	keywords := "Mitsubishi Eclipse"
	if params.Model != "" {
		keywords = fmt.Sprintf("Mitsubishi Eclipse %s", params.Model)
	}

	// Build query parameters
	queryParams := url.Values{}
	queryParams.Set("categoryId", "223") // Auto parts category
	queryParams.Set("keywords", keywords)
	queryParams.Set("locationStr", "Deutschland")
	queryParams.Set("radius", "0")
	queryParams.Set("sortingField", "")
	queryParams.Set("adType", "")
	queryParams.Set("posterType", "")
	queryParams.Set("maxPrice", "")
	queryParams.Set("minPrice", "")
	queryParams.Set("buyNowEnabled", "false")
	queryParams.Set("shippingCarrier", "")
	queryParams.Set("shipping", "")

	// Add pagination if needed (Kleinanzeigen uses page numbers)
	if params.Offset > 0 {
		page := (params.Offset / 25) + 1 // 25 items per page
		queryParams.Set("pageNum", fmt.Sprintf("%d", page))
	}

	searchURL := fmt.Sprintf("%s/s-suchanfrage.html?%s", c.baseURL, queryParams.Encode())
	return searchURL, nil
}

// extractPart extracts part information from an article element
func (c *KleinanzeigenClient) extractPart(ctx context.Context, s *goquery.Selection) (siteclients.Part, error) {
	part := siteclients.Part{
		SiteID: c.siteID,
	}

	// Extract ad ID (part ID)
	adID, exists := s.Attr("data-adid")
	if !exists || adID == "" {
		log.Printf("[KleinanzeigenClient] ERROR: missing data-adid attribute")
		return part, fmt.Errorf("missing data-adid")
	}
	log.Printf("[KleinanzeigenClient] Found ad ID: %s", adID)
	part.ID = adID

	// Extract relative URL
	relativeURL, exists := s.Attr("data-href")
	if !exists || relativeURL == "" {
		log.Printf("[KleinanzeigenClient] ERROR: missing data-href attribute")
		return part, fmt.Errorf("missing data-href")
	}
	part.URL = c.baseURL + relativeURL
	log.Printf("[KleinanzeigenClient] URL: %s", part.URL)

	// Extract title
	title := s.Find("h2 a.ellipsis").Text()
	title = strings.TrimSpace(title)
	if title == "" {
		log.Printf("[KleinanzeigenClient] ERROR: missing title")
		return part, fmt.Errorf("missing title")
	}
	part.Name = title
	log.Printf("[KleinanzeigenClient] Title: %s", title)

	// Extract description
	description := s.Find("p.aditem-main--middle--description").Text()
	description = strings.TrimSpace(description)
	part.Description = description

	// Extract price (optional)
	price := s.Find("p.aditem-main--middle--price-shipping--price").Text()
	price = strings.TrimSpace(price)
	if price != "" {
		// Add price to description
		part.Description = fmt.Sprintf("%s | Price: %s", part.Description, price)
	}

	// Extract location (optional)
	location := s.Find(".aditem-main--top--left").Text()
	location = strings.TrimSpace(location)
	if location != "" {
		// Clean up location (remove icon text)
		location = strings.ReplaceAll(location, "\n", " ")
		location = strings.TrimSpace(location)
		part.TypeName = location
	} else {
		part.TypeName = "Unknown Location"
	}

	// Extract image URL
	imgSrc, exists := s.Find(".imagebox img").Attr("src")
	if exists && imgSrc != "" {
		// Fetch and convert image to base64
		imageBase64, err := c.fetchImageAsBase64(ctx, imgSrc)
		if err != nil {
			log.Printf("[KleinanzeigenClient] Warning: failed to fetch image for part %s: %v", adID, err)
		} else {
			part.ImageBase64 = imageBase64
		}
	}

	return part, nil
}

// fetchImageAsBase64 fetches an image and returns it as base64
func (c *KleinanzeigenClient) fetchImageAsBase64(ctx context.Context, imageURL string) (string, error) {
	// Handle protocol-relative URLs
	if strings.HasPrefix(imageURL, "//") {
		imageURL = "https:" + imageURL
	}

	req, err := http.NewRequestWithContext(ctx, "GET", imageURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create image request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:143.0) Gecko/20100101 Firefox/143.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code for image: %d", resp.StatusCode)
	}

	// Read image data
	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read image data: %w", err)
	}

	// Convert to base64
	base64String := base64.StdEncoding.EncodeToString(imageData)

	return base64String, nil
}
