# DSM Parts Finder API

A simple Go REST API backend for the DSM Parts Finder application.

## Prerequisites

- Go 1.21 or higher
- Git

## Setup

1. Navigate to the api directory:
   ```bash
   cd api
   ```

2. Initialize and download dependencies:
   ```bash
   go mod tidy
   ```

3. Run the server:
   ```bash
   go run main.go
   ```

The API will start on `http://localhost:8080`

## API Endpoints

### Health Check
- `GET /health` - Returns API health status

### Parts API
- `GET /api/v1/parts` - Get all parts (demo data)
- `GET /api/v1/parts/:id` - Get a specific part by ID
- `POST /api/v1/parts` - Create a new part

## Example Usage

```bash
# Health check
curl http://localhost:8080/health

# Get all parts
curl http://localhost:8080/api/v1/parts

# Get specific part
curl http://localhost:8080/api/v1/parts/1

# Create new part
curl -X POST http://localhost:8080/api/v1/parts \
  -H "Content-Type: application/json" \
  -d '{"name":"Test Part","description":"A test part","price":"$19.99"}'
```

## CORS Configuration

The API is configured to accept requests from:
- `http://localhost:3000` (typical React dev server)
- `http://localhost:5173` (Vite dev server)

## Development

The server uses Gin framework with the following features:
- CORS middleware for frontend integration
- JSON request/response handling
- Basic error handling
- RESTful API structure

To modify the API, edit `main.go` and restart the server.