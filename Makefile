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
# Build the Go API
build-api: build-frontend
	@echo "Building Go API..."
	mkdir -p builds && \
	cd api && \
	go mod tidy && \
	go build -o ../builds/dsmpartsfinder . && \
	rm -rf ../builds/frontend && cp -r ../frontend/dist ../builds/frontend

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf frontend/dist
	rm -rf builds
	cd api && rm -rf logs
