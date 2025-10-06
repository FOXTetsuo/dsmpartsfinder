package scrapers

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
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
// Automatically fetches all pages until no more results are found
func (c *KleinanzeigenClient) FetchParts(ctx context.Context, params siteclients.SearchParams) ([]siteclients.Part, error) {
	log.Printf("[KleinanzeigenClient] Starting fetch with params: %+v", params)

	allParts := make([]siteclients.Part, 0)
	page := 1
	maxPages := 100    // Safety limit to prevent infinite loops
	itemsPerPage := 25 // Kleinanzeigen shows 25 items per page

	for page <= maxPages {
		log.Printf("[KleinanzeigenClient] Fetching page %d...", page)

		// Build search URL with page number
		searchURL, err := c.buildSearchURLWithPage(params, page)
		if err != nil {
			return nil, fmt.Errorf("failed to build search URL: %w", err)
		}

		log.Printf("[KleinanzeigenClient] Page %d URL: %s", page, searchURL)

		// Fetch the page
		pageParts, err := c.fetchSinglePage(ctx, searchURL)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch page %d: %w", page, err)
		}

		log.Printf("[KleinanzeigenClient] Page %d: got %d parts", page, len(pageParts))

		// If no parts found, we've reached the end
		if len(pageParts) == 0 {
			log.Printf("[KleinanzeigenClient] No more parts found on page %d, stopping", page)
			break
		}

		allParts = append(allParts, pageParts...)

		// If we got fewer parts than a full page, this is the last page
		if len(pageParts) < itemsPerPage {
			log.Printf("[KleinanzeigenClient] Got less than full page (%d < %d), this is the last page", len(pageParts), itemsPerPage)
			break
		}

		// Check if limit is set and we've reached it
		if params.Limit > 0 && len(allParts) >= params.Limit {
			log.Printf("[KleinanzeigenClient] Reached limit of %d parts, stopping", params.Limit)
			allParts = allParts[:params.Limit]
			break
		}

		page++
	}

	log.Printf("[KleinanzeigenClient] Finished fetching. Total parts: %d from %d page(s)", len(allParts), page)
	return allParts, nil
}

// fetchSinglePage fetches and parses a single page
func (c *KleinanzeigenClient) fetchSinglePage(ctx context.Context, searchURL string) ([]siteclients.Part, error) {
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

	// Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse HTML
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(bodyBytes)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	// Initialize parts slice
	parts := make([]siteclients.Part, 0)

	// Use article.aditem selector - this finds the actual ad listings
	selector := "article.aditem"
	articleCount := doc.Find(selector).Length()

	if articleCount == 0 {
		return parts, nil
	}

	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		part, err := c.extractPart(ctx, s)
		if err != nil {
			log.Printf("[KleinanzeigenClient] Warning: failed to extract part %d: %v", i, err)
			return
		}
		parts = append(parts, part)
	})

	log.Printf("[KleinanzeigenClient] Extracted %d parts from page", len(parts))
	return parts, nil
}

// buildSearchURLWithPage constructs the search URL with parameters and page number
func (c *KleinanzeigenClient) buildSearchURLWithPage(params siteclients.SearchParams, page int) (string, error) {
	// Build the search keywords
	keywords := "Mitsubishi Eclipse D30"
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

	// Add page number (Kleinanzeigen uses pageNum parameter)
	if page > 1 {
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
		return part, fmt.Errorf("missing data-adid")
	}
	part.ID = adID

	// Extract relative URL
	relativeURL, exists := s.Attr("data-href")
	if !exists || relativeURL == "" {
		return part, fmt.Errorf("missing data-href")
	}
	part.URL = c.baseURL + relativeURL

	// Extract title
	title := s.Find("h2 a.ellipsis").Text()
	title = strings.TrimSpace(title)
	if title == "" {
		return part, fmt.Errorf("missing title")
	}
	part.Name = title

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
