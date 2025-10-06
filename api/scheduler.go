package main

import (
	"context"
	"log"
	"time"

	"dsmpartsfinder-api/siteclients"

	"github.com/robfig/cron/v3"
)

// Scheduler handles automatic scheduled tasks
type Scheduler struct {
	cron         *cron.Cron
	partsService *PartsService
}

// NewScheduler creates a new scheduler instance
func NewScheduler(partsService *PartsService) *Scheduler {
	// Create cron with seconds precision
	c := cron.New(cron.WithSeconds())

	return &Scheduler{
		cron:         c,
		partsService: partsService,
	}
}

// Start begins the scheduled tasks
func (s *Scheduler) Start() error {
	log.Println("[Scheduler] Setting up scheduled tasks...")

	// Schedule fetch at midnight (00:00:00)
	_, err := s.cron.AddFunc("0 0 0 * * *", func() {
		log.Println("[Scheduler] Running scheduled midnight fetch...")
		s.fetchAllParts()
	})
	if err != nil {
		return err
	}
	log.Println("[Scheduler] Scheduled: Fetch all parts at midnight (00:00)")

	// Schedule fetch at midday (12:00:00)
	_, err = s.cron.AddFunc("0 0 12 * * *", func() {
		log.Println("[Scheduler] Running scheduled midday fetch...")
		s.fetchAllParts()
	})
	if err != nil {
		return err
	}
	log.Println("[Scheduler] Scheduled: Fetch all parts at midday (12:00)")

	// Start the cron scheduler
	s.cron.Start()
	log.Println("[Scheduler] Scheduler started successfully")

	return nil
}

// Stop stops the scheduler
func (s *Scheduler) Stop() {
	log.Println("[Scheduler] Stopping scheduler...")
	s.cron.Stop()
	log.Println("[Scheduler] Scheduler stopped")
}

// fetchAllParts fetches parts from all registered sites
func (s *Scheduler) fetchAllParts() {
	startTime := time.Now()
	log.Println("[Scheduler] ========================================")
	log.Println("[Scheduler] Starting automatic parts fetch...")
	log.Println("[Scheduler] ========================================")

	// Get all registered site IDs
	siteIDs := s.partsService.GetRegisteredSiteIDs()
	if len(siteIDs) == 0 {
		log.Println("[Scheduler] WARNING: No site clients registered")
		return
	}

	log.Printf("[Scheduler] Fetching from %d site(s): %v", len(siteIDs), siteIDs)

	// Default search parameters - fetch everything
	params := siteclients.SearchParams{
		VehicleType: "P",
		Make:        "Mitsubishi",
		BaseModel:   "Eclipse",
		Model:       "",
		YearFrom:    1989,
		YearTo:      2012,
		Offset:      0,
		Limit:       10000, // High limit to get everything
	}

	// Track statistics
	totalParts := 0
	totalNew := 0
	totalErrors := 0

	// Fetch from each site
	ctx := context.Background()
	for _, siteID := range siteIDs {
		siteStartTime := time.Now()
		log.Printf("[Scheduler] Fetching from site ID: %d", siteID)

		parts, err := s.partsService.FetchAndStoreParts(ctx, siteID, params)
		if err != nil {
			log.Printf("[Scheduler] ERROR: Failed to fetch from site %d: %v", siteID, err)
			totalErrors++
			continue
		}

		siteDuration := time.Since(siteStartTime)
		totalParts += len(parts)
		totalNew += len(parts)

		log.Printf("[Scheduler] Site %d: Fetched %d new parts in %v", siteID, len(parts), siteDuration)
	}

	// Log summary
	duration := time.Since(startTime)
	log.Println("[Scheduler] ========================================")
	log.Printf("[Scheduler] Fetch completed in %v", duration)
	log.Printf("[Scheduler] Total new parts: %d", totalNew)
	log.Printf("[Scheduler] Sites processed: %d/%d", len(siteIDs)-totalErrors, len(siteIDs))
	if totalErrors > 0 {
		log.Printf("[Scheduler] Errors encountered: %d", totalErrors)
	}
	log.Println("[Scheduler] ========================================")
}

// GetNextRuns returns the next scheduled run times
func (s *Scheduler) GetNextRuns() []time.Time {
	entries := s.cron.Entries()
	nextRuns := make([]time.Time, 0, len(entries))
	for _, entry := range entries {
		nextRuns = append(nextRuns, entry.Next)
	}
	return nextRuns
}
