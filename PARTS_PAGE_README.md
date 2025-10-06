# Parts Page Documentation

## Overview

The Parts page is a new feature that allows you to view all parts stored in the database and fetch parts from all registered sites with a single click.

## Features

### 1. View All Parts
- Displays all parts from the database in a paginated table
- Shows part details including:
  - ID
  - Image (with thumbnail preview)
  - Part ID
  - Name
  - Description
  - Type
  - Site ID
  - Link to original source

### 2. Fetch Parts from All Sites
- Single button click to fetch parts from all registered site clients
- Automatically stores fetched parts in the database
- Shows success/error feedback
- Displays statistics about the fetch operation

### 3. Statistics
- Total number of parts in database
- Number of unique sites represented

## UI Components

The page uses **Naive UI** components for a modern, clean interface:
- `NCard` - For section containers
- `NButton` - For action buttons with loading states
- `NDataTable` - For displaying parts with pagination and sorting
- `NImage` - For part image previews
- `NTag` - For site ID badges
- `NStatistic` - For displaying counts
- `NAlert` - For success/error messages
- `NSpace` - For consistent spacing

## Backend Endpoints

### GET `/api/parts`
Retrieves all parts from the database with pagination.

**Query Parameters:**
- `limit` (default: 50) - Number of parts to return
- `offset` (default: 0) - Starting position

**Response:**
```json
{
  "data": [...],
  "message": "Parts retrieved successfully",
  "total": 100,
  "limit": 50,
  "offset": 0
}
```

### POST `/api/parts/fetch-all`
Fetches parts from all registered site clients and stores them in the database.

**Request Body (optional):**
```json
{
  "vehicle_type": "Car",
  "make": "Mitsubishi",
  "base_model": "Eclipse",
  "model": "GSX",
  "year_from": 1990,
  "year_to": 1999,
  "limit": 100,
  "offset": 0
}
```

**Response:**
```json
{
  "data": [...],
  "total": 150,
  "sites": 2,
  "message": "Parts fetched from all sites",
  "errors": {
    "1": "Error message if any site failed"
  }
}
```

## How to Use

1. **Navigate to the Parts page:**
   - Click "Parts" in the navigation menu
   - Or visit `http://localhost:5173/parts`

2. **View existing parts:**
   - Parts are loaded automatically when the page opens
   - Use pagination controls at the bottom of the table
   - Change page size using the dropdown (10, 20, 50, 100 items per page)
   - Sort by ID by clicking the column header

3. **Fetch parts from all sites:**
   - Click the "Fetch Parts from All Sites" button
   - Confirm the action in the dialog
   - Wait for the operation to complete (may take a while)
   - View success/error alerts at the top of the page
   - Parts list will automatically refresh after fetching

4. **View part details:**
   - Hover over truncated text to see full content in tooltip
   - Click "View Source" to open the original part listing in a new tab
   - Click on images to view them in full size

## Technical Details

### Frontend
- **Framework:** Vue 3 Composition API
- **UI Library:** Naive UI 2.43.1
- **HTTP Client:** Axios
- **Routing:** Vue Router 4

### Backend
- **Framework:** Gin (Go)
- **Service:** PartsService manages site clients and database operations
- **Concurrency:** Fetches from all sites sequentially, continues on errors

### Error Handling
- Individual site errors don't stop the entire fetch operation
- Errors are collected and displayed to the user
- Console logs provide detailed error information for debugging

## Future Enhancements

Potential improvements for the Parts page:
- [ ] Advanced filtering (by site, type, make, model)
- [ ] Search functionality
- [ ] Bulk delete operations
- [ ] Export to CSV/JSON
- [ ] Part comparison feature
- [ ] Favorites/bookmarks
- [ ] Real-time updates with WebSockets
- [ ] Parallel fetching from multiple sites
- [ ] Progress indicator for long-running fetch operations