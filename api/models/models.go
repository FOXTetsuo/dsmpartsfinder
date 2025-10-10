package models

import "time"

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
}

// Site represents a parts supplier website
type Site struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

// CreateSiteRequest represents the request body for creating a site
type CreateSiteRequest struct {
	Name string `json:"name" binding:"required"`
	URL  string `json:"url" binding:"required"`
}

// UpdateSiteRequest represents the request body for updating a site
type UpdateSiteRequest struct {
	Name string `json:"name" binding:"required"`
	URL  string `json:"url" binding:"required"`
}

// Part represents a car part scraped from a site
type Part struct {
	ID          int       `json:"id"`
	PartID      string    `json:"part_id"`
	Description string    `json:"description"`
	TypeName    string    `json:"type_name"`
	Name        string    `json:"name"`
	ImageBase64 string    `json:"image_base64"`
	URL         string    `json:"url"`
	SiteID      int       `json:"site_id"`
	Price       string    `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	LastSeen    time.Time `json:"last_seen"`
}

// FetchPartsRequest represents the request body for fetching parts from a site
type FetchPartsRequest struct {
	SiteID      int    `json:"site_id" binding:"required"`
	VehicleType string `json:"vehicle_type" binding:"required"`
	Make        string `json:"make" binding:"required"`
	BaseModel   string `json:"base_model" binding:"required"`
	Model       string `json:"model" binding:"required"`
	YearFrom    int    `json:"year_from"`
	YearTo      int    `json:"year_to"`
	Offset      int    `json:"offset"`
	Limit       int    `json:"limit"`
}
