# DSM Parts Finder

A full-stack web application for finding and managing parts for Diamond Star Motors (DSM) vehicles including Mitsubishi Eclipse, Eagle Talon, and Plymouth Laser.

## Architecture

- **Backend**: Go with Gin framework providing REST API
- **Frontend**: Vue 3 with Tailwind CSS for modern, responsive UI
- **Development**: Vite for fast frontend builds, Go modules for backend dependencies

## Quick Start

### Prerequisites

- Go 1.21 or higher
- Node.js 16+ and Yarn
- Git

### Backend Setup

1. Navigate to the API directory:
   ```bash
   cd api
   ```

2. Install Go dependencies:
   ```bash
   go mod tidy
   ```

3. Start the API server:
   ```bash
   go run main.go
   ```

The API will be available at `http://localhost:8080`

### Frontend Setup

1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```

2. Install Node dependencies:
   ```bash
   yarn install
   ```

3. Start the development server:
   ```bash
   yarn dev
   ```

The frontend will be available at `http://localhost:5173`

## Project Structure

```
dsmpartsfinder/
├── api/                    # Go backend API
│   ├── main.go            # Main API server file
│   ├── go.mod             # Go module dependencies
│   └── README.md          # Backend documentation
├── frontend/              # Vue 3 frontend
│   ├── src/               # Source code
│   │   ├── views/         # Page components
│   │   ├── components/    # Reusable components
│   │   ├── style.css      # Tailwind CSS styles
│   │   ├── App.vue        # Root component
│   │   └── main.js        # Application entry
│   ├── package.json       # Node dependencies
│   ├── vite.config.js     # Vite configuration
│   ├── tailwind.config.js # Tailwind CSS config
│   └── README.md          # Frontend documentation
└── README.md              # This file
```

## Features

### Current Features

- **API Health Check**: Monitor backend status
- **Parts Catalog**: Browse available DSM parts
- **Part Management**: Add new parts to inventory
- **Responsive Design**: Works on desktop and mobile
- **Modern UI**: Clean interface with Tailwind CSS

### Demo Data

The API includes demo parts data:
- Engine Oil Filter ($24.99)
- Brake Pads ($89.99)
- Air Filter ($45.99)
- Spark Plugs ($32.99)

## API Endpoints

### Health
- `GET /health` - API health status

### Parts
- `GET /api/v1/parts` - Get all parts
- `GET /api/v1/parts/:id` - Get specific part
- `POST /api/v1/parts` - Create new part

## Development

### Backend Development

The Go API uses:
- Gin web framework
- CORS middleware for frontend integration
- JSON request/response handling
- RESTful API design

To modify the API:
1. Edit `api/main.go`
2. Restart the server: `go run main.go`

### Frontend Development

The Vue 3 frontend uses:
- Vue Router for navigation
- Axios for API calls
- Tailwind CSS for styling
- Vite for development server

To modify the frontend:
1. Edit files in `frontend/src/`
2. Changes auto-reload in development

### Styling

The application uses Tailwind CSS with:
- Custom color schemes
- Responsive breakpoints
- Component utility classes
- Modern animations

## Building for Production

### Backend
```bash
cd api
go build -o dsmpartsfinder-api
```

### Frontend
```bash
cd frontend
yarn build
```

## Supported Vehicles

This application is designed for Diamond Star Motors vehicles:

- **Mitsubishi Eclipse** (1989-1999)
  - 1G: 1989-1994
  - 2G: 1995-1999

- **Eagle Talon** (1989-1998)
  - TSi, TSi AWD, ESi variants

- **Plymouth Laser** (1989-1994)
  - Base, RS, RS Turbo variants

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

### Code Style

- **Go**: Follow standard Go conventions
- **Vue**: Use Composition API and ES6+
- **CSS**: Use Tailwind utility classes
- **Commits**: Use descriptive commit messages

## Troubleshooting

### Common Issues

1. **API not responding**: Ensure Go server is running on port 8080
2. **Frontend build errors**: Check Node.js version (16+ required) and Yarn installation
3. **CORS errors**: Verify API CORS configuration matches frontend URL
4. **Module errors**: Run `go mod tidy` and `yarn install`

### Development Tips

- Use browser dev tools for frontend debugging
- Check Go server logs for API issues
- Ensure both frontend and backend are running for full functionality
- API proxy is configured in Vite for seamless development

## License

This project is for educational and demonstration purposes.

## Future Enhancements

- User authentication and authorization
- Advanced search and filtering
- Part images and detailed specifications
- Shopping cart and ordering system
- Inventory management
- Part compatibility checking
- Price comparison features
- User reviews and ratings