package routes

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	. "dsmpartsfinder-api/models"
	"dsmpartsfinder-api/siteclients"

	"github.com/gin-gonic/gin"
)

type SQLClient interface {
	GetAllSites() ([]Site, error)
	GetSiteByID(id int) (*Site, error)
	CreateSite(name, url string) (*Site, error)
	UpdateSite(id int, name, url string) (*Site, error)
	DeleteSite(id int) error

	GetAllParts(limit, offset int) ([]Part, error)
	GetPartByID(id int) (*Part, error)
	GetPartsBySiteID(siteID, limit, offset int) ([]Part, error)
	DeletePartsBySiteID(siteID int) error
}

type PartsService interface {
	FetchAndStoreParts(ctx context.Context, siteID int, params siteclients.SearchParams) ([]Part, error)
	GetRegisteredSiteIDs() []int
	GetAllParts(limit, offset int) ([]Part, error)
	GetFilteredParts(limit, offset int, typeFilter string, siteID int, newerThan time.Time) ([]Part, error)
	GetPartByID(id int) (*Part, error)
	GetPartsBySiteID(siteID, limit, offset int) ([]Part, error)
	DeletePartsBySiteID(siteID int) error
	GetTotalPartsCount() (int, error)
	GetFilteredPartsCount(typeFilter string, siteID int, newerThan time.Time) (int, error)
}

func RegisterAPIRoutes(r *gin.Engine, sqlClient SQLClient, partsService PartsService) {
	api := r.Group("/api")
	{
		// Health check endpoint
		api.GET("/health", func(c *gin.Context) {
			response := HealthResponse{
				Status:    "healthy",
				Timestamp: time.Now(),
				Message:   "DSM Parts Finder API is running",
			}
			c.JSON(http.StatusOK, response)
		})

		// GET /api/sites - Get all sites
		api.GET("/sites", func(c *gin.Context) {
			sites, err := sqlClient.GetAllSites()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Failed to query sites",
					"details": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"data":    sites,
				"message": "Sites retrieved successfully",
				"total":   len(sites),
			})
		})

		// GET /api/sites/:id - Get a single site by ID
		api.GET("/sites/:id", func(c *gin.Context) {
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Invalid site ID",
				})
				return
			}

			site, err := sqlClient.GetSiteByID(id)
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Site not found",
				})
				return
			} else if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Failed to query site",
					"details": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"data":    site,
				"message": "Site retrieved successfully",
			})
		})

		// POST /api/parts/fetch - Fetch parts from all sites
		api.POST("/parts/fetch", func(c *gin.Context) {
			log.Println("[POST /api/parts/fetch] Endpoint called")

			var req FetchPartsRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				log.Printf("[POST /api/parts/fetch] ERROR: Invalid request body: %v", err)
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "Invalid request body",
					"details": err.Error(),
				})
				return
			}

			log.Printf("[POST /api/parts/fetch] Request: SiteID=%d, Limit=%d", req.SiteID, req.Limit)

			// Convert to search params
			params := siteclients.SearchParams{}

			// Fetch and store parts
			parts, err := partsService.FetchAndStoreParts(c.Request.Context(), req.SiteID, params)
			if err != nil {
				log.Printf("[POST /api/parts/fetch] ERROR: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Failed to fetch and store parts",
					"details": err.Error(),
				})
				return
			}

			log.Printf("[POST /api/parts/fetch] Successfully fetched and stored %d parts", len(parts))

			c.JSON(http.StatusOK, gin.H{
				"data":    parts,
				"message": "Parts fetched and stored successfully",
				"total":   len(parts),
			})
		})

		// POST /api/parts/fetch-all - Fetch parts from all registered sites
		api.POST("/parts/fetch-all", func(c *gin.Context) {
			log.Println("[POST /api/parts/fetch-all] Endpoint called")

			var req struct {
				VehicleType string `json:"vehicle_type"`
				Make        string `json:"make"`
				BaseModel   string `json:"base_model"`
				Model       string `json:"model"`
				YearFrom    int    `json:"year_from"`
				YearTo      int    `json:"year_to"`
				Offset      int    `json:"offset"`
				Limit       int    `json:"limit"`
			}

			if err := c.ShouldBindJSON(&req); err != nil {
				// If no body provided, use defaults
				log.Printf("[POST /api/parts/fetch-all] No valid JSON body, using defaults. Error: %v", err)
				req.YearFrom = 1960
				req.YearTo = 2025
				req.Limit = 30
			}

			// Set defaults if not provided
			if req.YearFrom == 0 {
				req.YearFrom = 1960
			}
			if req.YearTo == 0 {
				req.YearTo = 2025
			}
			if req.Limit == 0 {
				req.Limit = 30
			}

			log.Printf("[POST /api/parts/fetch-all] Request params: YearFrom=%d, YearTo=%d, Limit=%d, Make=%s, Model=%s",
				req.YearFrom, req.YearTo, req.Limit, req.Make, req.Model)

			// Convert to search params
			params := siteclients.SearchParams{
				VehicleType: req.VehicleType,
				Make:        req.Make,
				BaseModel:   req.BaseModel,
				Model:       req.Model,
				YearFrom:    req.YearFrom,
				YearTo:      req.YearTo,
				Offset:      req.Offset,
				Limit:       req.Limit,
			}

			// Get all registered site IDs
			siteIDs := partsService.GetRegisteredSiteIDs()
			log.Printf("[POST /api/parts/fetch-all] Found %d registered site(s): %v", len(siteIDs), siteIDs)

			if len(siteIDs) == 0 {
				log.Println("[POST /api/parts/fetch-all] ERROR: No site clients registered")
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "No site clients registered",
				})
				return
			}

			// Fetch and store parts from all sites
			log.Println("[POST /api/parts/fetch-all] Starting to fetch from all sites...")
			allParts := make([]Part, 0)
			errors := make(map[int]string)

			for _, siteID := range siteIDs {
				log.Printf("[POST /api/parts/fetch-all] Fetching from site ID: %d", siteID)
				parts, err := partsService.FetchAndStoreParts(c.Request.Context(), siteID, params)
				if err != nil {
					errors[siteID] = err.Error()
					log.Printf("[POST /api/parts/fetch-all] ERROR fetching parts from site %d: %v", siteID, err)
					continue
				}
				log.Printf("[POST /api/parts/fetch-all] Got %d parts from site %d", len(parts), siteID)
				allParts = append(allParts, parts...)
			}

			log.Printf("[POST /api/parts/fetch-all] Total parts collected: %d", len(allParts))

			response := gin.H{
				"data":    allParts,
				"total":   len(allParts),
				"sites":   len(siteIDs),
				"message": "Parts fetched from all sites",
			}

			if len(errors) > 0 {
				log.Printf("[POST /api/parts/fetch-all] Encountered errors for %d sites: %v", len(errors), errors)
				response["errors"] = errors
			}

			log.Printf("[POST /api/parts/fetch-all] Returning response with %d parts", len(allParts))
			c.JSON(http.StatusOK, response)
		})

		// GET /api/parts - Get all parts with pagination
		api.GET("/parts", func(c *gin.Context) {
			limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
			offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
			typeFilter := c.Query("type")
			siteID, _ := strconv.Atoi(c.DefaultQuery("site_id", "0"))

			// If any filter is specified, use filtered endpoint
			if typeFilter != "" || siteID != 0 || c.Query("newer_than_hours") != "" {
				var newerThan time.Time
				if c.Query("newer_than_hours") != "" {
					hours, _ := strconv.Atoi(c.DefaultQuery("newer_than_hours", "72"))
					newerThan = time.Now().Add(-time.Duration(hours) * time.Hour)
				}

				log.Printf("[GET /api/parts] Called with filters: limit=%d, offset=%d, type=%s, site_id=%d, newer_than=%v",
					limit, offset, typeFilter, siteID, newerThan)

				parts, err := partsService.GetFilteredParts(limit, offset, typeFilter, siteID, newerThan)
				if err != nil {
					log.Printf("[GET /api/parts] ERROR: %v", err)
					c.JSON(http.StatusInternalServerError, gin.H{
						"error":   "Failed to query parts",
						"details": err.Error(),
					})
					return
				}

				total, err := partsService.GetFilteredPartsCount(typeFilter, siteID, newerThan)
				if err != nil {
					log.Printf("[GET /api/parts] ERROR getting filtered count: %v", err)
					c.JSON(http.StatusInternalServerError, gin.H{
						"error":   "Failed to get total parts count",
						"details": err.Error(),
					})
					return
				}

				c.JSON(http.StatusOK, gin.H{
					"data":    parts,
					"message": "Filtered parts retrieved successfully",
					"total":   total,
					"limit":   limit,
					"offset":  offset,
				})
				return
			}

			// Default unfiltered behavior
			parts, err := partsService.GetAllParts(limit, offset)
			if err != nil {
				log.Printf("[GET /api/parts] ERROR: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Failed to query parts",
					"details": err.Error(),
				})
				return
			}

			total, err := partsService.GetTotalPartsCount()
			if err != nil {
				log.Printf("[GET /api/parts] ERROR getting total count: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Failed to get total parts count",
					"details": err.Error(),
				})
				return
			}

			log.Printf("[GET /api/parts] Returning %d parts (limit=%d, offset=%d, total=%d)", len(parts), limit, offset, total)
			if len(parts) > 0 {
				log.Printf("[GET /api/parts] First part: ID=%d, PartID=%s, Name=%s", parts[0].ID, parts[0].PartID, parts[0].Name)
			}

			c.JSON(http.StatusOK, gin.H{
				"data":    parts,
				"message": "Parts retrieved successfully",
				"total":   total,
				"limit":   limit,
				"offset":  offset,
			})
		})

		// GET /api/parts/:id - Get a single part by ID
		api.GET("/parts/:id", func(c *gin.Context) {
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Invalid part ID",
				})
				return
			}

			part, err := partsService.GetPartByID(id)
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Part not found",
				})
				return
			} else if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Failed to query part",
					"details": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"data":    part,
				"message": "Part retrieved successfully",
			})
		})

		// GET /api/sites/:id/parts - Get all parts for a specific site
		api.GET("/sites/:id/parts", func(c *gin.Context) {
			siteID, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Invalid site ID",
				})
				return
			}

			limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
			offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

			parts, err := partsService.GetPartsBySiteID(siteID, limit, offset)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Failed to query parts for site",
					"details": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"data":    parts,
				"message": "Parts retrieved successfully",
				"total":   len(parts),
				"site_id": siteID,
				"limit":   limit,
				"offset":  offset,
			})
		})

		// DELETE /api/sites/:id/parts - Delete all parts for a specific site
		api.DELETE("/sites/:id/parts", func(c *gin.Context) {
			siteID, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Invalid site ID",
				})
				return
			}

			err = partsService.DeletePartsBySiteID(siteID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Failed to delete parts for site",
					"details": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "Parts deleted successfully",
				"site_id": siteID,
			})
		})
	}
}
