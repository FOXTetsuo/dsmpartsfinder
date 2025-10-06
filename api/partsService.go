package main

import (
	"context"
	"fmt"
	"log"

	"dsmpartsfinder-api/siteclients"
)

// PartsService manages the fetching and storage of parts from various site clients
type PartsService struct {
	sqlClient   *SQLClient
	siteClients map[int]siteclients.SiteClient
}

// NewPartsService creates a new PartsService
func NewPartsService(sqlClient *SQLClient) *PartsService {
	return &PartsService{
		sqlClient:   sqlClient,
		siteClients: make(map[int]siteclients.SiteClient),
	}
}

// RegisterSiteClient registers a site client for a specific site ID
func (s *PartsService) RegisterSiteClient(siteID int, client siteclients.SiteClient) {
	s.siteClients[siteID] = client
	log.Printf("Registered site client '%s' for site ID %d", client.GetName(), siteID)
}

// GetSiteClient retrieves a site client by site ID
func (s *PartsService) GetSiteClient(siteID int) (siteclients.SiteClient, error) {
	client, exists := s.siteClients[siteID]
	if !exists {
		return nil, fmt.Errorf("no site client registered for site ID %d", siteID)
	}
	return client, nil
}

// FetchAndStoreParts fetches parts from a site client and stores them in the database
func (s *PartsService) FetchAndStoreParts(ctx context.Context, siteID int, params siteclients.SearchParams) ([]Part, error) {
	log.Printf("[FetchAndStoreParts] Starting fetch for site ID: %d with params: %+v", siteID, params)

	// Get the appropriate site client
	client, err := s.GetSiteClient(siteID)
	if err != nil {
		log.Printf("[FetchAndStoreParts] ERROR: Failed to get site client for ID %d: %v", siteID, err)
		return nil, err
	}

	log.Printf("[FetchAndStoreParts] Fetching parts from %s (site ID: %d)", client.GetName(), siteID)

	// Fetch parts from the site
	fetchedParts, err := client.FetchParts(ctx, params)
	if err != nil {
		log.Printf("[FetchAndStoreParts] ERROR: Failed to fetch parts from %s: %v", client.GetName(), err)
		return nil, fmt.Errorf("failed to fetch parts from %s: %w", client.GetName(), err)
	}

	log.Printf("[FetchAndStoreParts] Fetched %d parts from %s", len(fetchedParts), client.GetName())

	if len(fetchedParts) > 0 {
		log.Printf("[FetchAndStoreParts] First part example: ID=%s, Name=%s, Type=%s",
			fetchedParts[0].ID, fetchedParts[0].Name, fetchedParts[0].TypeName)
	}

	// Store parts in the database
	log.Printf("[FetchAndStoreParts] Starting to store %d parts in database", len(fetchedParts))
	storedParts := make([]Part, 0, len(fetchedParts))
	duplicateCount := 0
	errorCount := 0

	for i, part := range fetchedParts {
		// Try to create the part, skip if it already exists (duplicate key)
		storedPart, err := s.sqlClient.CreatePart(
			part.ID,
			part.Description,
			part.TypeName,
			part.Name,
			part.ImageBase64,
			part.URL,
			part.SiteID,
		)
		if err != nil {
			// Log the error but continue with other parts
			if i < 3 { // Log details for first 3 errors only
				log.Printf("[FetchAndStoreParts] Warning: failed to store part %s (index %d): %v", part.ID, i, err)
			}
			errorCount++
			continue
		}
		storedParts = append(storedParts, *storedPart)
		if i < 3 { // Log first 3 successful stores
			log.Printf("[FetchAndStoreParts] Successfully stored part: ID=%s, DB_ID=%d, Name=%s", part.ID, storedPart.ID, storedPart.Name)
		}
	}

	duplicateCount = errorCount // Most errors are likely duplicates
	log.Printf("[FetchAndStoreParts] Successfully stored %d new parts, %d duplicates/errors out of %d fetched",
		len(storedParts), duplicateCount, len(fetchedParts))

	return storedParts, nil
}

// FetchPartsOnly fetches parts from a site client without storing them
func (s *PartsService) FetchPartsOnly(ctx context.Context, siteID int, params siteclients.SearchParams) ([]siteclients.Part, error) {
	client, err := s.GetSiteClient(siteID)
	if err != nil {
		return nil, err
	}

	log.Printf("Fetching parts from %s (site ID: %d) without storing", client.GetName(), siteID)

	parts, err := client.FetchParts(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch parts from %s: %w", client.GetName(), err)
	}

	log.Printf("Fetched %d parts from %s", len(parts), client.GetName())

	return parts, nil
}

// GetPartsBySiteID retrieves all parts for a specific site from the database
func (s *PartsService) GetPartsBySiteID(siteID int, limit, offset int) ([]Part, error) {
	log.Printf("[GetPartsBySiteID] Called with siteID=%d, limit=%d, offset=%d", siteID, limit, offset)
	parts, err := s.sqlClient.GetPartsBySiteID(siteID, limit, offset)
	if err != nil {
		log.Printf("[GetPartsBySiteID] ERROR: %v", err)
		return nil, err
	}
	log.Printf("[GetPartsBySiteID] Retrieved %d parts from database for site %d", len(parts), siteID)
	return parts, nil
}

// GetAllParts retrieves all parts from the database
func (s *PartsService) GetAllParts(limit, offset int) ([]Part, error) {
	log.Printf("[GetAllParts] Called with limit=%d, offset=%d", limit, offset)
	parts, err := s.sqlClient.GetAllParts(limit, offset)
	if err != nil {
		log.Printf("[GetAllParts] ERROR: %v", err)
		return nil, err
	}
	log.Printf("[GetAllParts] Retrieved %d parts from database", len(parts))
	return parts, nil
}

// GetPartByID retrieves a specific part by its ID
func (s *PartsService) GetPartByID(id int) (*Part, error) {
	return s.sqlClient.GetPartByID(id)
}

// DeletePartsBySiteID deletes all parts for a specific site
func (s *PartsService) DeletePartsBySiteID(siteID int) error {
	return s.sqlClient.DeletePartsBySiteID(siteID)
}

// GetRegisteredSiteIDs returns a list of all registered site IDs
func (s *PartsService) GetRegisteredSiteIDs() []int {
	siteIDs := make([]int, 0, len(s.siteClients))
	for id := range s.siteClients {
		siteIDs = append(siteIDs, id)
	}
	log.Printf("[GetRegisteredSiteIDs] Returning %d registered site IDs: %v", len(siteIDs), siteIDs)
	return siteIDs
}
