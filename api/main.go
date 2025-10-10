package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"dsmpartsfinder-api/routes"
	"dsmpartsfinder-api/scrapers"
	"dsmpartsfinder-api/siteclients"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Open database connection
	sqlClient, err := NewSQLClient("./sqlite.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer sqlClient.Close()

	// // Run database migrations
	// migrationRunner := NewMigrationRunner(sqlClient.db)
	// if err := migrationRunner.Run("./migrations"); err != nil {
	// 	log.Fatalf("Failed to run migrations: %v", err)
	// }

	// Initialize PartsService
	partsService := NewPartsService(sqlClient)

	sites, err := sqlClient.GetAllSites()
	if err != nil {
		log.Fatalf("Failed to get sites from database: %v", err)
	}

	// Register site clients dynamically based on DB entries
	for _, site := range sites {
		switch site.Name {
		case "SchadeAutos":
			client := siteclients.NewSchadeAutosClient(site.ID)
			partsService.RegisterSiteClient(site.ID, client)
			log.Printf("Registered SchadeAutos client (site ID: %d)", site.ID)
		case "Kleinanzeigen":
			client := scrapers.NewKleinanzeigenClient(site.ID)
			partsService.RegisterSiteClient(site.ID, client)
			log.Printf("Registered Kleinanzeigen client (site ID: %d)", site.ID)
		default:
			log.Printf("No client implementation for site '%s' (site ID: %d), skipping registration", site.Name, site.ID)
		}
	}

	// Initialize and start scheduler for automatic fetching
	scheduler := NewScheduler(partsService)
	if err := scheduler.Start(); err != nil {
		log.Fatalf("Failed to start scheduler: %v", err)
	}
	defer scheduler.Stop()

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

	// Register API endpoints from routes.go
	routes.RegisterAPIRoutes(r, sqlClient, partsService)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
