package main

import (
	"context"
	"fmt"
	"log"
	"time"

	. "dsmpartsfinder-api/models"
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
// It also updates last_seen for existing parts and optionally deletes stale parts
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

	// Check which parts already exist in the database
	log.Printf("[FetchAndStoreParts] Checking for existing parts in database")
	partIDs := make([]string, len(fetchedParts))
	for i, part := range fetchedParts {
		partIDs[i] = part.ID
	}

	existingParts, err := s.sqlClient.GetExistingPartIDs(partIDs, siteID)
	if err != nil {
		log.Printf("[FetchAndStoreParts] ERROR: Failed to check existing parts: %v", err)
		return nil, fmt.Errorf("failed to check existing parts: %w", err)
	}

	duplicateCount := len(existingParts)
	log.Printf("[FetchAndStoreParts] Found %d existing parts, %d new parts to insert", duplicateCount, len(fetchedParts)-duplicateCount)

	// Update last_seen for existing parts
	if len(existingParts) > 0 {
		existingPartIDs := make([]string, 0, len(existingParts))
		for partID := range existingParts {
			existingPartIDs = append(existingPartIDs, partID)
		}
		log.Printf("[FetchAndStoreParts] Updating last_seen for %d existing parts", len(existingPartIDs))
		if err := s.sqlClient.UpdateLastSeen(existingPartIDs, siteID); err != nil {
			log.Printf("[FetchAndStoreParts] WARNING: Failed to update last_seen: %v", err)
			// Don't fail the entire operation, just log the error
		}
	}

	// Delete stale parts (last seen more than 3 days ago)
	olderThan := time.Now().AddDate(0, 0, -3)
	deletedCount, err := s.sqlClient.DeleteStaleParts(siteID, olderThan)
	if err != nil {
		log.Printf("[FetchAndStoreParts] WARNING: Failed to delete stale parts: %v", err)
	} else {
		log.Printf("[FetchAndStoreParts] Deleted %d stale parts for site ID %d", deletedCount, siteID)
	}

	// Store only new parts in the database
	log.Printf("[FetchAndStoreParts] Starting to store new parts in database")
	storedParts := make([]Part, 0, len(fetchedParts)-duplicateCount)
	errorCount := 0
	insertedCount := 0

	for i, part := range fetchedParts {
		// Skip if part already exists
		if existingParts[part.ID] {
			if i < 3 { // Log first 3 skipped duplicates
				log.Printf("[FetchAndStoreParts] Skipping duplicate part: ID=%s, Name=%s", part.ID, part.Name)
			}
			continue
		}

		// Insert the new part
		storedPart, err := s.sqlClient.CreatePart(
			part.ID,
			part.Description,
			part.TypeName,
			part.Name,
			part.ImageBase64,
			part.URL,
			part.SiteID,
			part.Price,
			part.CreationDate,
		)
		if err != nil {
			// Log the error but continue with other parts
			if errorCount < 3 { // Log details for first 3 errors only
				log.Printf("[FetchAndStoreParts] Warning: failed to store part %s (index %d): %v", part.ID, i, err)
			}
			errorCount++
			continue
		}
		storedParts = append(storedParts, *storedPart)
		insertedCount++
		if insertedCount <= 3 { // Log first 3 successful stores
			log.Printf("[FetchAndStoreParts] Successfully stored part: ID=%s, DB_ID=%d, Name=%s", part.ID, storedPart.ID, storedPart.Name)
		}
	}

	log.Printf("[FetchAndStoreParts] Successfully stored %d new parts, skipped %d duplicates, %d errors out of %d fetched",
		len(storedParts), duplicateCount, errorCount, len(fetchedParts))

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
	parts, err := s.sqlClient.GetAllParts(limit, offset)
	if err != nil {
		log.Printf("[GetAllParts] ERROR: %v", err)
		return nil, err
	}
	return parts, nil
}

// GetFilteredParts retrieves filtered parts from the database
func (s *PartsService) GetFilteredParts(limit, offset int, typeFilter string, siteIDs []int, newerThan time.Time, search string, sortBy string, sortDesc bool) ([]Part, error) {
	log.Printf("[GetFilteredParts] Called with limit=%d, offset=%d, typeFilter=%s, siteIDs=%v, newerThan=%v, search=%s, sortBy=%s, sortDesc=%v",
		limit, offset, typeFilter, siteIDs, newerThan, search, sortBy, sortDesc)
	parts, err := s.sqlClient.GetFilteredParts(limit, offset, typeFilter, siteIDs, newerThan, search, sortBy, sortDesc)
	if err != nil {
		log.Printf("[GetFilteredParts] ERROR: %v", err)
		return nil, err
	}
	log.Printf("[GetFilteredParts] Retrieved %d parts from database", len(parts))
	return parts, nil
}

func (s *PartsService) GetTotalPartsCount() (int, error) {
	count, err := s.sqlClient.GetTotalPartsCount()
	if err != nil {
		log.Printf("[GetTotalPartsCount] ERROR: %v", err)
		return 0, err
	}
	log.Printf("[GetTotalPartsCount] Total count: %d", count)
	return count, nil
}

func (s *PartsService) GetFilteredPartsCount(typeFilter string, siteIDs []int, newerThan time.Time, search string) (int, error) {
	count, err := s.sqlClient.GetFilteredPartsCount(typeFilter, siteIDs, newerThan, search)
	if err != nil {
		log.Printf("[GetFilteredPartsCount] ERROR: %v", err)
		return 0, err
	}
	log.Printf("[GetFilteredPartsCount] Total count: %d (typeFilter=%s, siteIDs=%v, newerThan=%v, search=%s)",
		count, typeFilter, siteIDs, newerThan, search)
	return count, nil
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
