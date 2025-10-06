package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/glebarez/go-sqlite"
)

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
}

type DemoItem struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       string `json:"price"`
}

type Site struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type CreateSiteRequest struct {
	Name string `json:"name" binding:"required"`
	URL  string `json:"url" binding:"required"`
}

type UpdateSiteRequest struct {
	Name string `json:"name" binding:"required"`
	URL  string `json:"url" binding:"required"`
}

func openSQLite() (*sql.DB, error) {
	// Connect to the SQLite database
	db, err := sql.Open("sqlite", "./sqlite.db")
	if err != nil {
		return nil, err
	}

	// Test the connection
	var sqliteVersion string
	err = db.QueryRow("select sqlite_version()").Scan(&sqliteVersion)
	if err != nil {
		db.Close()
		return nil, err
	}

	fmt.Println("Connected to SQLite database successfully. Version:", sqliteVersion)
	return db, nil
}

func main() {
	r := gin.Default()

	// Open database connection
	db, err := openSQLite()
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		return
	}
	defer db.Close()

	// Global error recovery middleware
	r.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		log.Printf("PANIC: %v", recovered)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"details": fmt.Sprintf("%v", recovered),
		})
	}))

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Demo API endpoints
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
			rows, err := db.Query("SELECT id, site_url, site_name FROM sites")
			if err != nil {
				log.Printf("ERROR: Failed to query sites: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Failed to query sites",
					"details": err.Error(),
				})
				return
			}
			defer rows.Close()

			var sites []Site
			for rows.Next() {
				var site Site
				err := rows.Scan(&site.ID, &site.URL, &site.Name)
				if err != nil {
					log.Printf("ERROR: Failed to scan site data: %v", err)
					c.JSON(http.StatusInternalServerError, gin.H{
						"error":   "Failed to scan site data",
						"details": err.Error(),
					})
					return
				}
				sites = append(sites, site)
			}

			if err = rows.Err(); err != nil {
				log.Printf("ERROR: Error iterating sites: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Error iterating sites",
					"details": err.Error(),
				})
				return
			}

			log.Printf("Successfully retrieved %d sites", len(sites))
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

			var site Site
			err = db.QueryRow("SELECT id, site_url, site_name FROM sites WHERE id = ?", id).Scan(&site.ID, &site.URL, &site.Name)
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Site not found",
				})
				return
			} else if err != nil {
				log.Printf("ERROR: Failed to query site: %v", err)
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

		// POST /api/sites - Create a new site
		api.POST("/sites", func(c *gin.Context) {
			var req CreateSiteRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "Invalid request body",
					"details": err.Error(),
				})
				return
			}

			result, err := db.Exec("INSERT INTO sites (site_url, site_name) VALUES (?, ?)", req.URL, req.Name)
			if err != nil {
				log.Printf("ERROR: Failed to create site: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Failed to create site",
					"details": err.Error(),
				})
				return
			}

			id, err := result.LastInsertId()
			if err != nil {
				log.Printf("ERROR: Failed to get last insert ID: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Failed to get created site ID",
					"details": err.Error(),
				})
				return
			}

			site := Site{
				ID:   int(id),
				Name: req.Name,
				URL:  req.URL,
			}

			log.Printf("Successfully created site with ID %d", id)
			c.JSON(http.StatusCreated, gin.H{
				"data":    site,
				"message": "Site created successfully",
			})
		})

		// PUT /api/sites/:id - Update a site
		api.PUT("/sites/:id", func(c *gin.Context) {
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Invalid site ID",
				})
				return
			}

			var req UpdateSiteRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "Invalid request body",
					"details": err.Error(),
				})
				return
			}

			result, err := db.Exec("UPDATE sites SET site_url = ?, site_name = ? WHERE id = ?", req.URL, req.Name, id)
			if err != nil {
				log.Printf("ERROR: Failed to update site: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Failed to update site",
					"details": err.Error(),
				})
				return
			}

			rowsAffected, err := result.RowsAffected()
			if err != nil {
				log.Printf("ERROR: Failed to get rows affected: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Failed to verify update",
					"details": err.Error(),
				})
				return
			}

			if rowsAffected == 0 {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Site not found",
				})
				return
			}

			site := Site{
				ID:   id,
				Name: req.Name,
				URL:  req.URL,
			}

			log.Printf("Successfully updated site with ID %d", id)
			c.JSON(http.StatusOK, gin.H{
				"data":    site,
				"message": "Site updated successfully",
			})
		})

		// DELETE /api/sites/:id - Delete a site
		api.DELETE("/sites/:id", func(c *gin.Context) {
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Invalid site ID",
				})
				return
			}

			result, err := db.Exec("DELETE FROM sites WHERE id = ?", id)
			if err != nil {
				log.Printf("ERROR: Failed to delete site: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Failed to delete site",
					"details": err.Error(),
				})
				return
			}

			rowsAffected, err := result.RowsAffected()
			if err != nil {
				log.Printf("ERROR: Failed to get rows affected: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Failed to verify deletion",
					"details": err.Error(),
				})
				return
			}

			if rowsAffected == 0 {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Site not found",
				})
				return
			}

			log.Printf("Successfully deleted site with ID %d", id)
			c.JSON(http.StatusOK, gin.H{
				"message": "Site deleted successfully",
			})
		})

	}

	// Start server on port 8080
	r.Run(":8080")
}
