# DSM Parts Finder Frontend

A modern Vue 3 frontend application with Tailwind CSS for the DSM Parts Finder system.

## Features

- **Vue 3** with Composition API support
- **Tailwind CSS** for responsive, utility-first styling
- **Vue Router** for client-side routing
- **Axios** for API communication
- **Vite** for fast development and building
- Responsive design optimized for mobile and desktop
- Component-based architecture
- API proxy configuration for development

## Prerequisites

- Node.js 16+ and Yarn
- The backend API running on `http://localhost:8080`

## Setup

1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```

2. Install dependencies:
   ```bash
   yarn install
   ```

3. Start the development server:
   ```bash
   yarn dev
   ```

The application will be available at `http://localhost:5173`

## Available Scripts

- `yarn dev` - Start development server
- `yarn build` - Build for production
- `yarn preview` - Preview production build locally
- `yarn lint` - Run ESLint

## Project Structure

```
src/
├── components/          # Reusable Vue components
├── views/              # Page components
│   ├── Home.vue        # Landing page
│   └── Parts.vue       # Parts catalog page
├── router/             # Vue Router configuration
├── style.css           # Global Tailwind CSS styles
├── App.vue             # Root application component
└── main.js             # Application entry point
```

## API Integration

The frontend communicates with the Go backend API through:
- Vite proxy configuration (in `vite.config.js`)
- Axios HTTP client for API requests
- Base URL: `/api/v1/`

### API Endpoints Used

- `GET /health` - API health check
- `GET /api/v1/parts` - Fetch all parts
- `GET /api/v1/parts/:id` - Fetch single part
- `POST /api/v1/parts` - Create new part

## Styling

The application uses Tailwind CSS with:
- Custom color palette (primary/secondary themes)
- Responsive breakpoints
- Component utilities for common UI patterns
- Custom animations and transitions

### Key Tailwind Classes

- `.btn` - Button base styles
- `.card` - Card layout styles
- `.form-input` / `.form-textarea` - Form element styles
- `.alert` - Notification styles

## Pages

### Home (`/`)
- Hero section with branding
- Feature highlights
- Supported DSM vehicle models
- API status indicator
- Call-to-action sections

### Parts Catalog (`/parts`)
- Parts inventory display
- Add new part form
- Search and filtering
- Responsive grid layout
- Loading/error states

## Development

### Adding New Components

1. Create component in `src/components/`
2. Import and use in views or other components
3. Follow Vue 3 Composition API patterns

### Styling Guidelines

- Use Tailwind utility classes
- Follow mobile-first responsive design
- Maintain consistent spacing and typography
- Use the defined color palette

### API Integration

- Use Axios for HTTP requests
- Handle loading states
- Implement error handling
- Follow RESTful conventions

## Building for Production

```bash
yarn build
```

This creates optimized static files in the `dist/` directory.

## Browser Support

- Modern browsers (Chrome, Firefox, Safari, Edge)
- ES2015+ support required
- Responsive design for mobile devices

## Troubleshooting

### Development Server Issues

1. Ensure Node.js 16+ is installed
2. Clear Yarn cache: `yarn cache clean`
3. Delete `node_modules` and run `yarn install` again

### API Connection Issues

1. Verify the backend API is running on port 8080
2. Check the proxy configuration in `vite.config.js`
3. Ensure CORS is properly configured in the backend

### Build Issues

1. Check for TypeScript/ESLint errors
2. Verify all dependencies are installed
3. Ensure all imports are correct

## Contributing

1. Follow Vue 3 and ES6+ best practices
2. Use Tailwind CSS for styling
3. Maintain responsive design principles
4. Test components thoroughly
5. Keep code clean and well-documented