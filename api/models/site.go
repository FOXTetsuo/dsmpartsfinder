package models

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
