package models

import "time"

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
