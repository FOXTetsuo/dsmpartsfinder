package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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

func main() {
	r := gin.Default()

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		response := HealthResponse{
			Status:    "healthy",
			Timestamp: time.Now(),
			Message:   "DSM Parts Finder API is running",
		}
		c.JSON(http.StatusOK, response)
	})

	// Demo API endpoints
	api := r.Group("/api/v1")
	{
		// GET /api/v1/parts - Demo parts list
		api.GET("/parts", func(c *gin.Context) {
			demoParts := []DemoItem{
				{ID: 1, Name: "Engine Oil Filter", Description: "High-quality oil filter for DSM engines", Price: "$24.99"},
				{ID: 2, Name: "Brake Pads", Description: "Performance brake pads front set", Price: "$89.99"},
				{ID: 3, Name: "Air Filter", Description: "K&N style air filter", Price: "$45.99"},
				{ID: 4, Name: "Spark Plugs", Description: "NGK iridium spark plugs set of 4", Price: "$32.99"},
			}
			c.JSON(http.StatusOK, gin.H{
				"data":    demoParts,
				"message": "Parts retrieved successfully",
				"total":   len(demoParts),
			})
		})

		// GET /api/v1/parts/:id - Get single part
		api.GET("/parts/:id", func(c *gin.Context) {
			id := c.Param("id")
			// Demo response for any ID
			demoPart := DemoItem{
				ID:          1,
				Name:        "Engine Oil Filter",
				Description: "High-quality oil filter for DSM engines",
				Price:       "$24.99",
			}
			c.JSON(http.StatusOK, gin.H{
				"data":    demoPart,
				"message": "Part retrieved successfully",
				"id":      id,
			})
		})

		// POST /api/v1/parts - Create new part (demo)
		api.POST("/parts", func(c *gin.Context) {
			var newPart DemoItem
			if err := c.ShouldBindJSON(&newPart); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "Invalid JSON",
					"message": err.Error(),
				})
				return
			}

			// Demo response
			newPart.ID = 999 // Mock generated ID
			c.JSON(http.StatusCreated, gin.H{
				"data":    newPart,
				"message": "Part created successfully",
			})
		})
	}

	// Start server on port 8080
	r.Run(":8080")
}
