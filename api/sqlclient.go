package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/glebarez/go-sqlite"
)

// SQLClient wraps database operations for the DSM Parts Finder
type SQLClient struct {
	db *sql.DB
}

// NewSQLClient creates and initializes a new SQLClient
func NewSQLClient(dbPath string) (*SQLClient, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	var sqliteVersion string
	err = db.QueryRow("SELECT sqlite_version()").Scan(&sqliteVersion)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to verify database connection: %w", err)
	}

	log.Printf("Connected to SQLite database successfully. Version: %s", sqliteVersion)

	return &SQLClient{db: db}, nil
}

// Close closes the database connection
func (c *SQLClient) Close() error {
	return c.db.Close()
}

// logError logs an error with a context message
func logError(context string, err error) {
	log.Printf("ERROR: %s: %v", context, err)
}

// logSuccess logs a success message
func logSuccess(message string) {
	log.Printf("SUCCESS: %s", message)
}

// GetAllSites retrieves all sites from the database
func (c *SQLClient) GetAllSites() ([]Site, error) {
	rows, err := c.db.Query("SELECT id, site_url, site_name FROM sites")
	if err != nil {
		logError("Failed to query sites", err)
		return nil, err
	}
	defer rows.Close()

	var sites []Site
	for rows.Next() {
		var site Site
		err := rows.Scan(&site.ID, &site.URL, &site.Name)
		if err != nil {
			logError("Failed to scan site data", err)
			return nil, err
		}
		sites = append(sites, site)
	}

	if err = rows.Err(); err != nil {
		logError("Error iterating sites", err)
		return nil, err
	}

	logSuccess(fmt.Sprintf("Retrieved %d sites", len(sites)))
	return sites, nil
}

// GetSiteByID retrieves a single site by its ID
func (c *SQLClient) GetSiteByID(id int) (*Site, error) {
	var site Site
	err := c.db.QueryRow("SELECT id, site_url, site_name FROM sites WHERE id = ?", id).
		Scan(&site.ID, &site.URL, &site.Name)

	if err == sql.ErrNoRows {
		return nil, sql.ErrNoRows
	} else if err != nil {
		logError(fmt.Sprintf("Failed to query site with ID %d", id), err)
		return nil, err
	}

	logSuccess(fmt.Sprintf("Retrieved site with ID %d", id))
	return &site, nil
}

// CreateSite creates a new site in the database
func (c *SQLClient) CreateSite(name, url string) (*Site, error) {
	result, err := c.db.Exec("INSERT INTO sites (site_url, site_name) VALUES (?, ?)", url, name)
	if err != nil {
		logError("Failed to create site", err)
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		logError("Failed to get last insert ID", err)
		return nil, err
	}

	site := &Site{
		ID:   int(id),
		Name: name,
		URL:  url,
	}

	logSuccess(fmt.Sprintf("Created site with ID %d", id))
	return site, nil
}

// UpdateSite updates an existing site in the database
func (c *SQLClient) UpdateSite(id int, name, url string) (*Site, error) {
	result, err := c.db.Exec("UPDATE sites SET site_url = ?, site_name = ? WHERE id = ?", url, name, id)
	if err != nil {
		logError(fmt.Sprintf("Failed to update site with ID %d", id), err)
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logError("Failed to get rows affected", err)
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, sql.ErrNoRows
	}

	site := &Site{
		ID:   id,
		Name: name,
		URL:  url,
	}

	logSuccess(fmt.Sprintf("Updated site with ID %d", id))
	return site, nil
}

// DeleteSite deletes a site from the database
func (c *SQLClient) DeleteSite(id int) error {
	result, err := c.db.Exec("DELETE FROM sites WHERE id = ?", id)
	if err != nil {
		logError(fmt.Sprintf("Failed to delete site with ID %d", id), err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logError("Failed to get rows affected", err)
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	logSuccess(fmt.Sprintf("Deleted site with ID %d", id))
	return nil
}

// CreatePart creates a new part in the database
func (c *SQLClient) CreatePart(partID, description, typeName, name, imageBase64, url string, siteID int) (*Part, error) {
	result, err := c.db.Exec(`
		INSERT INTO parts (part_id, description, type_name, name, image_base64, url, site_id, last_seen)
		VALUES (?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
	`, partID, description, typeName, name, imageBase64, url, siteID)
	if err != nil {
		logError("Failed to create part", err)
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		logError("Failed to get last insert ID for part", err)
		return nil, err
	}

	part := &Part{
		ID:          int(id),
		PartID:      partID,
		Description: description,
		TypeName:    typeName,
		Name:        name,
		ImageBase64: imageBase64,
		URL:         url,
		SiteID:      siteID,
	}

	logSuccess(fmt.Sprintf("Created part with ID %d", id))
	return part, nil
}

// GetPartByID retrieves a single part by its database ID
func (c *SQLClient) GetPartByID(id int) (*Part, error) {
	var part Part
	err := c.db.QueryRow(`
		SELECT id, part_id, description, type_name, name, image_base64, url, site_id, created_at, updated_at, last_seen
		FROM parts WHERE id = ?
	`, id).Scan(
		&part.ID, &part.PartID, &part.Description, &part.TypeName,
		&part.Name, &part.ImageBase64, &part.URL, &part.SiteID,
		&part.CreatedAt, &part.UpdatedAt, &part.LastSeen,
	)

	if err == sql.ErrNoRows {
		return nil, sql.ErrNoRows
	} else if err != nil {
		logError(fmt.Sprintf("Failed to query part with ID %d", id), err)
		return nil, err
	}

	logSuccess(fmt.Sprintf("Retrieved part with ID %d", id))
	return &part, nil
}

// GetPartsBySiteID retrieves all parts for a specific site
func (c *SQLClient) GetPartsBySiteID(siteID int, limit, offset int) ([]Part, error) {
	query := `
		SELECT id, part_id, description, type_name, name, image_base64, url, site_id, created_at, updated_at, last_seen
		FROM parts
		WHERE site_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := c.db.Query(query, siteID, limit, offset)
	if err != nil {
		logError(fmt.Sprintf("Failed to query parts for site ID %d", siteID), err)
		return nil, err
	}
	defer rows.Close()

	var parts []Part
	for rows.Next() {
		var part Part
		err := rows.Scan(
			&part.ID, &part.PartID, &part.Description, &part.TypeName,
			&part.Name, &part.ImageBase64, &part.URL, &part.SiteID,
			&part.CreatedAt, &part.UpdatedAt, &part.LastSeen,
		)
		if err != nil {
			logError("Failed to scan part data", err)
			return nil, err
		}
		parts = append(parts, part)
	}

	if err = rows.Err(); err != nil {
		logError("Error iterating parts", err)
		return nil, err
	}

	logSuccess(fmt.Sprintf("Retrieved %d parts for site ID %d", len(parts), siteID))
	return parts, nil
}

// GetAllParts retrieves all parts from the database
func (c *SQLClient) GetAllParts(limit, offset int) ([]Part, error) {
	query := `
		SELECT id, part_id, description, type_name, name, image_base64, url, site_id, created_at, updated_at, last_seen
		FROM parts
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := c.db.Query(query, limit, offset)
	if err != nil {
		logError("Failed to query parts", err)
		return nil, err
	}
	defer rows.Close()

	var parts []Part
	for rows.Next() {
		var part Part
		err := rows.Scan(
			&part.ID, &part.PartID, &part.Description, &part.TypeName,
			&part.Name, &part.ImageBase64, &part.URL, &part.SiteID,
			&part.CreatedAt, &part.UpdatedAt, &part.LastSeen,
		)
		if err != nil {
			logError("Failed to scan part data", err)
			return nil, err
		}
		parts = append(parts, part)
	}

	if err = rows.Err(); err != nil {
		logError("Error iterating parts", err)
		return nil, err
	}

	logSuccess(fmt.Sprintf("Retrieved %d parts", len(parts)))
	return parts, nil
}

// GetExistingPartIDs checks which part IDs already exist in the database for a given site
// Returns a map where keys are part IDs that exist
func (c *SQLClient) GetExistingPartIDs(partIDs []string, siteID int) (map[string]bool, error) {
	if len(partIDs) == 0 {
		return make(map[string]bool), nil
	}

	// Build a query with placeholders for the IN clause
	placeholders := make([]string, len(partIDs))
	args := make([]interface{}, len(partIDs)+1)
	args[0] = siteID

	for i, partID := range partIDs {
		placeholders[i] = "?"
		args[i+1] = partID
	}

	query := fmt.Sprintf(`
		SELECT part_id FROM parts
		WHERE site_id = ? AND part_id IN (%s)
	`, strings.Join(placeholders, ","))

	rows, err := c.db.Query(query, args...)
	if err != nil {
		logError("Failed to query existing part IDs", err)
		return nil, err
	}
	defer rows.Close()

	existingParts := make(map[string]bool)
	for rows.Next() {
		var partID string
		if err := rows.Scan(&partID); err != nil {
			logError("Failed to scan part ID", err)
			return nil, err
		}
		existingParts[partID] = true
	}

	if err = rows.Err(); err != nil {
		logError("Error iterating existing part IDs", err)
		return nil, err
	}

	log.Printf("Found %d existing parts out of %d checked for site ID %d", len(existingParts), len(partIDs), siteID)
	return existingParts, nil
}

// UpdateLastSeen updates the last_seen timestamp for multiple parts
func (c *SQLClient) UpdateLastSeen(partIDs []string, siteID int) error {
	if len(partIDs) == 0 {
		return nil
	}

	placeholders := make([]string, len(partIDs))
	args := make([]interface{}, len(partIDs)+1)
	args[0] = siteID

	for i, partID := range partIDs {
		placeholders[i] = "?"
		args[i+1] = partID
	}

	query := fmt.Sprintf(`
		UPDATE parts
		SET last_seen = CURRENT_TIMESTAMP
		WHERE site_id = ? AND part_id IN (%s)
	`, strings.Join(placeholders, ","))

	result, err := c.db.Exec(query, args...)
	if err != nil {
		logError("Failed to update last_seen timestamps", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logError("Failed to get rows affected for last_seen update", err)
		return err
	}

	log.Printf("Updated last_seen for %d parts in site ID %d", rowsAffected, siteID)
	return nil
}

// DeleteStaleParts deletes parts that haven't been seen since a specific time
func (c *SQLClient) DeleteStaleParts(siteID int, olderThan time.Time) (int64, error) {
	result, err := c.db.Exec(`
		DELETE FROM parts
		WHERE site_id = ? AND last_seen < ?
	`, siteID, olderThan)
	if err != nil {
		logError(fmt.Sprintf("Failed to delete stale parts for site ID %d", siteID), err)
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logError("Failed to get rows affected for stale parts deletion", err)
		return 0, err
	}

	log.Printf("Deleted %d stale parts for site ID %d (older than %s)", rowsAffected, siteID, olderThan.Format("2006-01-02 15:04:05"))
	return rowsAffected, nil
}

// UpdatePart updates an existing part in the database
func (c *SQLClient) UpdatePart(id int, partID, description, typeName, name, imageBase64, url string, siteID int) (*Part, error) {
	result, err := c.db.Exec(`
		UPDATE parts
		SET part_id = ?, description = ?, type_name = ?, name = ?, image_base64 = ?, url = ?, site_id = ?, updated_at = CURRENT_TIMESTAMP, last_seen = CURRENT_TIMESTAMP
		WHERE id = ?
	`, partID, description, typeName, name, imageBase64, url, siteID, id)
	if err != nil {
		logError(fmt.Sprintf("Failed to update part with ID %d", id), err)
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logError("Failed to get rows affected", err)
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, sql.ErrNoRows
	}

	// Fetch and return the updated part
	return c.GetPartByID(id)
}

// DeletePart deletes a part from the database
func (c *SQLClient) DeletePart(id int) error {
	result, err := c.db.Exec("DELETE FROM parts WHERE id = ?", id)
	if err != nil {
		logError(fmt.Sprintf("Failed to delete part with ID %d", id), err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logError("Failed to get rows affected", err)
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	logSuccess(fmt.Sprintf("Deleted part with ID %d", id))
	return nil
}

// DeletePartsBySiteID deletes all parts for a specific site
func (c *SQLClient) DeletePartsBySiteID(siteID int) error {
	result, err := c.db.Exec("DELETE FROM parts WHERE site_id = ?", siteID)
	if err != nil {
		logError(fmt.Sprintf("Failed to delete parts for site ID %d", siteID), err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logError("Failed to get rows affected", err)
		return err
	}

	logSuccess(fmt.Sprintf("Deleted %d parts for site ID %d", rowsAffected, siteID))
	return nil
}
