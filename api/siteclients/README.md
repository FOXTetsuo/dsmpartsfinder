# Site Clients Package

This package provides a SOLID architecture for scraping car parts data from various supplier websites.

## Architecture Overview

The package follows the **Interface Segregation Principle** and **Open/Closed Principle** from SOLID:

- **SiteClient Interface**: Defines the contract that all site-specific implementations must follow
- **Concrete Implementations**: Each website (e.g., SchadeAutos) has its own client implementation
- **Extensibility**: New site clients can be added without modifying existing code

## Core Components

### 1. SiteClient Interface

```go
type SiteClient interface {
    GetName() string
    FetchParts(ctx context.Context, params SearchParams) ([]Part, error)
    GetSiteID() int
}
```

All site clients must implement these three methods:
- `GetName()`: Returns a human-readable name for the site
- `FetchParts()`: Fetches parts from the site based on search parameters
- `GetSiteID()`: Returns the database ID of the site this client represents

### 2. Part Model

```go
type Part struct {
    ID          string `json:"id"`
    Description string `json:"description"`
    TypeName    string `json:"type_name"`
    Name        string `json:"name"`
    ImageBase64 string `json:"image_base64"`
    URL         string `json:"url"`
    SiteID      int    `json:"site_id"`
}
```

This is the standardized format that all site clients must return parts in.

### 3. SearchParams

```go
type SearchParams struct {
    VehicleType string
    Make        string
    BaseModel   string
    Model       string
    YearFrom    int
    YearTo      int
    Offset      int
    Limit       int
}
```

## Existing Implementations

### SchadeAutos Client

The SchadeAutos client scrapes data from `https://www.schadeautos.nl`.

**Features:**
- POST request to `/parts/eng/search.json` API endpoint
- Parses JSON response containing part data
- Fetches product images and converts them to base64
- Handles relative and absolute image URLs
- Includes proper HTTP headers to mimic browser behavior

**Example Usage:**

```go
// Create the client
client := siteclients.NewSchadeAutosClient(1) // 1 is the site ID

// Define search parameters for Mitsubishi Eclipse D3
params := siteclients.SearchParams{
    VehicleType: "P",           // Passenger car
    Make:        "A0001E2D",    // Mitsubishi
    BaseModel:   "A0001FHK",    // Eclipse
    Model:       "A0001FHL",    // Eclipse D3
    YearFrom:    1960,
    YearTo:      2025,
    Offset:      0,
    Limit:       30,
}

// Fetch parts
ctx := context.Background()
parts, err := client.FetchParts(ctx, params)
if err != nil {
    log.Fatal(err)
}

for _, part := range parts {
    fmt.Printf("Part: %s - %s\n", part.Name, part.Description)
}
```

## Creating a New Site Client

To add support for a new website, follow these steps:

### Step 1: Create a New Client File

Create a new file in the `siteclients` package (e.g., `newSiteClient.go`):

```go
package siteclients

import (
    "context"
    "fmt"
    "net/http"
    "time"
)

type NewSiteClient struct {
    baseURL    string
    httpClient *http.Client
    siteID     int
}

func NewNewSiteClient(siteID int) *NewSiteClient {
    return &NewSiteClient{
        baseURL: "https://www.newsite.com",
        httpClient: &http.Client{
            Timeout: 30 * time.Second,
        },
        siteID: siteID,
    }
}

func (c *NewSiteClient) GetName() string {
    return "NewSite"
}

func (c *NewSiteClient) GetSiteID() int {
    return c.siteID
}

func (c *NewSiteClient) FetchParts(ctx context.Context, params SearchParams) ([]Part, error) {
    // TODO: Implement site-specific scraping logic
    // 1. Build HTTP request with appropriate parameters
    // 2. Execute request
    // 3. Parse response (JSON, HTML, XML, etc.)
    // 4. Convert to []Part format
    // 5. Fetch and convert images to base64
    return nil, fmt.Errorf("not implemented")
}
```

### Step 2: Register the Client in main.go

```go
// In main.go
newSiteClient := siteclients.NewNewSiteClient(2) // Use appropriate site ID
partsService.RegisterSiteClient(2, newSiteClient)
```

### Step 3: Test the Implementation

Create a test to verify your implementation works correctly:

```go
func TestNewSiteClient(t *testing.T) {
    client := NewNewSiteClient(2)

    params := SearchParams{
        // Set appropriate parameters
    }

    parts, err := client.FetchParts(context.Background(), params)
    if err != nil {
        t.Fatalf("FetchParts failed: %v", err)
    }

    if len(parts) == 0 {
        t.Error("Expected at least one part")
    }
}
```

## Best Practices

1. **Error Handling**: Always return descriptive errors
2. **Context Support**: Respect context cancellation for graceful shutdowns
3. **Timeouts**: Set appropriate HTTP client timeouts
4. **Rate Limiting**: Consider implementing rate limiting to avoid overwhelming target sites
5. **User Agent**: Use a proper User-Agent header
6. **Image Handling**: Always handle missing/broken images gracefully
7. **Logging**: Log important events and errors for debugging
8. **Testing**: Write tests for your client implementation

## API Integration

The site clients are integrated with the REST API through the `PartsService`:

### Endpoints

**Fetch and Store Parts:**
```bash
POST /api/parts/fetch
Content-Type: application/json

{
  "site_id": 1,
  "vehicle_type": "P",
  "make": "A0001E2D",
  "base_model": "A0001FHK",
  "model": "A0001FHL",
  "year_from": 1960,
  "year_to": 2025,
  "offset": 0,
  "limit": 30
}
```

**Get Parts for a Site:**
```bash
GET /api/sites/:id/parts?limit=50&offset=0
```

**Get All Parts:**
```bash
GET /api/parts?limit=50&offset=0
```

## Data Flow

1. **API Request** → Handler receives fetch request
2. **Service Layer** → PartsService selects appropriate SiteClient
3. **Client** → SiteClient makes HTTP request to target website
4. **Parsing** → Client parses response and converts to Part structs
5. **Image Fetching** → Client downloads and converts images to base64
6. **Storage** → PartsService stores parts in database
7. **Response** → API returns stored parts to client

## Future Enhancements

- [ ] Add caching layer to avoid redundant requests
- [ ] Implement background job queue for large scraping operations
- [ ] Add retry logic with exponential backoff
- [ ] Support for proxy rotation
- [ ] Implement rate limiting per site
- [ ] Add metrics and monitoring
- [ ] Support for incremental updates (only fetch new parts)
- [ ] Add support for pagination in API responses

## Troubleshooting

### Common Issues

**1. HTTP Request Fails**
- Check if the target website is accessible
- Verify firewall/proxy settings
- Ensure proper headers are set

**2. Parse Errors**
- Verify the response format matches your parser
- Check if the website updated their API/HTML structure
- Add logging to
