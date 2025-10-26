.PHONY: build build-frontend build-api clean run

# Default target
build: build-frontend build-api

# Build the frontend
build-frontend:
	@echo "Building frontend..."
	cd frontend && \
	if [ ! -d "node_modules" ]; then \
		echo "Installing dependencies..."; \
		yarn install; \
	fi && \
	yarn build


build-api: build-frontend
	@echo "Copying frontend dist to api folder..."
	rm -rf api/frontend/dist
	mkdir -p api/frontend
	cp -r frontend/dist api/frontend/
	@echo "Building Go API with embedded frontend..."
	mkdir -p builds && \
	cd api && \
	go mod tidy && \
	go build -o ../builds/dsmpartsfinder .

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf frontend/dist
	rm -rf builds
	cd api && rm -rf logs

# Run in development mode (without embedding)
run:
	@echo "Running in development mode..."
	cd api && go run .
