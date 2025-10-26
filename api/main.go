package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"dsmpartsfinder-api/routes"
	"dsmpartsfinder-api/scrapers"
	"dsmpartsfinder-api/siteclients"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
)

//go:embed migrations
var migrationsFS embed.FS

//go:embed frontend/dist
var frontendFS embed.FS

func main() {
	// Get debug mode
	debug := os.Getenv("DEBUG")
	if debug == "" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	logsDir := "logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		log.Fatalf("Failed to create logs directory: %v", err)
	}

	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	// Get configurable database path
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "./sqlite.db"
	}

	// Get configurable port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting DSM Parts finder on port %s", port)
	log.Printf("Opening database connection on %s", dbPath)
	log.Printf("Application is running in %s mode", gin.Mode())

	// If in release mode, log to file
	if gin.Mode() == gin.ReleaseMode {
		currentDate := time.Now().Format("2006-01-02")
		logFilePath := filepath.Join(logsDir, currentDate+".log")
		logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}
		log.SetOutput(logFile)
		// Note: not closing the file, as it should remain open for the duration
	}

	r := gin.Default()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Open database connection
	sqlClient, err := NewSQLClient("./sqlite.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer sqlClient.Close()

	subFS, err := fs.Sub(migrationsFS, "migrations")
	if err != nil {
		log.Fatalf("Failed to create sub FS: %v", err)
	}
	provider, err := goose.NewProvider(goose.DialectSQLite3, sqlClient.db, subFS)
	if err != nil {
		log.Fatalf("Failed to create migration provider: %v", err)
	}
	_, err = provider.Up(ctx)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

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
		case "Kleinanzeigen":
			client := scrapers.NewKleinanzeigenClient(site.ID)
			partsService.RegisterSiteClient(site.ID, client)
		case "Ebay":
			clientID := os.Getenv("EBAY_CLIENT_ID")
			clientSecret := os.Getenv("EBAY_CLIENT_SECRET")
			client := siteclients.NewEbayClient(site.ID, clientID, clientSecret, false)
			partsService.RegisterSiteClient(site.ID, client)
		default:
			log.Printf("No client implementation for site '%s' (site ID: %d), skipping registration", site.Name, site.ID)
		}
	}

	// Initialize and start scheduler for automatic fetching
	scheduler := NewScheduler(partsService)
	go func() {
		if err := scheduler.Start(); err != nil {
			log.Printf("Scheduler error: %v", err)
		}
	}()
	defer scheduler.Stop()

	// Global error recovery middleware
	r.Use(gin.CustomRecovery(func(c *gin.Context, recovered any) {
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

	// Serve embedded frontend files
	frontendSubFS, err := fs.Sub(frontendFS, "frontend/dist")
	if err != nil {
		log.Fatalf("Failed to create frontend sub FS: %v", err)
	}

	// SPA fallback and static file serving from embedded FS
	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
			return
		}

		// Remove leading slash for fs.Sub compatibility
		path := strings.TrimPrefix(c.Request.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}

		// Try to read the file from embedded FS
		data, err := fs.ReadFile(frontendSubFS, path)
		if err != nil {
			// File not found, serve index.html for SPA routing
			data, err = fs.ReadFile(frontendSubFS, "index.html")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load frontend"})
				return
			}
			c.Data(http.StatusOK, "text/html; charset=utf-8", data)
			return
		}

		// Determine content type based on file extension
		contentType := getContentType(path)
		c.Data(http.StatusOK, contentType, data)
	})

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getContentType(path string) string {
	ext := filepath.Ext(path)
	switch ext {
	case ".html":
		return "text/html; charset=utf-8"
	case ".css":
		return "text/css; charset=utf-8"
	case ".js":
		return "application/javascript; charset=utf-8"
	case ".json":
		return "application/json; charset=utf-8"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".svg":
		return "image/svg+xml"
	case ".ico":
		return "image/x-icon"
	case ".woff":
		return "font/woff"
	case ".woff2":
		return "font/woff2"
	case ".ttf":
		return "font/ttf"
	case ".eot":
		return "application/vnd.ms-fontobject"
	default:
		return "application/octet-stream"
	}
}
