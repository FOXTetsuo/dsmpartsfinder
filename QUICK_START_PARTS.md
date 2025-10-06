# Quick Start Guide: Testing the Parts Page

## Prerequisites

Before you begin, ensure you have:
- Go 1.20+ installed
- Node.js 16+ and yarn installed
- The database initialized with at least one site registered

## Step 1: Start the Backend API

```bash
cd dsmpartsfinder/api
go run .
```

The API should start on `http://localhost:8080`

You should see:
```
Registered SchadeAutos client
Registered site client 'SchadeAutos' for site ID 1
[GIN-debug] Listening and serving HTTP on :8080
```

## Step 2: Start the Frontend

Open a new terminal:

```bash
cd dsmpartsfinder/frontend
yarn install  # if first time
yarn dev
```

The frontend should start on `http://localhost:5173`

## Step 3: Test the Parts Page

1. **Navigate to Parts Page:**
   - Open your browser to `http://localhost:5173`
   - Click "Parts" in the navigation menu
   - You should see the Parts Inventory page

2. **View Existing Parts:**
   - The page will automatically load any existing parts from the database
   - If the database is empty, you'll see an empty table

3. **Fetch Parts from All Sites:**
   - Click the blue "Fetch Parts from All Sites" button
   - Confirm the action when prompted
   - Wait for the operation to complete (may take 30-60 seconds)
   - You should see:
     - A loading spinner while fetching
     - A success alert with the number of parts fetched
     - The parts table automatically refreshes with new data

4. **Interact with the Table:**
   - **Pagination:** Use controls at bottom to navigate pages
   - **Page Size:** Change from 20 to 10, 50, or 100 items per page
   - **Sorting:** Click the ID column header to sort
   - **View Images:** Click on any part image to see it full size
   - **View Source:** Click "View Source" to open the original listing
   - **Tooltips:** Hover over truncated text to see full content

## Expected Results

### After Clicking "Fetch Parts from All Sites":

1. **Success Alert (Green):**
   ```
   Successfully fetched 100 parts from 1 site(s)
   ```

2. **Parts Table Populated:**
   - Shows parts with images, IDs, names, descriptions
   - Each row has a "Site 1" badge
   - "View Source" links are clickable

3. **Statistics Card:**
   - Total Parts: 100 (or however many were fetched)
   - Unique Sites: 1 (or more if you have multiple sites)

### What You Can Test:

✅ **Fetch Functionality:**
- Click "Fetch Parts from All Sites" multiple times
- Should see message about fetching and storing parts
- Duplicate parts should be skipped (check backend logs)

✅ **Table Features:**
- Pagination works correctly
- Changing page size updates display
- Sorting by ID works
- Image previews load and expand on click

✅ **Refresh:**
- Click the "Refresh" button
- Table should reload from database
- Should maintain current page/size settings

✅ **Navigation:**
- Switch between Home, Sites, and Parts pages
- Navigation highlights current page
- Mobile menu works on small screens

## Troubleshooting

### Backend not starting?
- Check if port 8080 is already in use
- Verify `sqlite.db` file exists in the api directory
- Check Go dependencies: `go mod tidy`

### Frontend not starting?
- Delete `node_modules` and run `yarn install` again
- Check if port 5173 is available
- Verify Node.js version: `node --version` (should be 16+)

### "No site clients registered" error?
- You need at least one site registered in the database
- The SchadeAutos client is registered for site ID 1
- Check that your database has a site with ID 1

### Parts not loading?
- Check browser console for errors (F12)
- Verify API is running: visit `http://localhost:8080/api/health`
- Check API logs for error messages

### CORS errors?
- Make sure frontend is running on `http://localhost:5173`
- Backend CORS config allows this origin by default

## API Endpoints Reference

Test these directly with curl or Postman:

```bash
# Health check
curl http://localhost:8080/api/health

# Get all parts
curl http://localhost:8080/api/parts?limit=10&offset=0

# Fetch from all sites
curl -X POST http://localhost:8080/api/parts/fetch-all \
  -H "Content-Type: application/json" \
  -d '{"limit": 50}'

# Get all sites
curl http://localhost:8080/api/sites
```

## Next Steps

Once everything is working:

1. **Add More Sites:** Go to the Sites page and add more parts sources
2. **Register Site Clients:** Update `main.go` to register clients for your new sites
3. **Customize Fetch Parameters:** Modify the fetch request to filter by make/model/year
4. **Explore the Code:** Check out `Parts.vue` for frontend and `partsService.go` for backend logic

## Need Help?

- Check `PARTS_PAGE_README.md` for detailed documentation
- Review backend logs in the API terminal
- Check browser console for frontend errors
- Inspect network tab to see API requests/responses