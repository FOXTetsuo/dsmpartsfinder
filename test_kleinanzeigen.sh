#!/bin/bash

# Test script for Kleinanzeigen scraper
# This script tests fetching parts from Kleinanzeigen.de

API_URL="http://localhost:8080"

echo "========================================="
echo "Testing Kleinanzeigen Scraper (Site ID: 2)"
echo "========================================="
echo ""

# Check if API is running
echo "1. Checking if API is running..."
if ! curl -s "${API_URL}/api/health" > /dev/null; then
    echo "❌ API is not running. Please start it with: cd api && go run ."
    exit 1
fi
echo "✅ API is running"
echo ""

# Check sites in database
echo "2. Checking registered sites..."
SITES=$(curl -s "${API_URL}/api/sites" | jq -r '.data[] | "\(.id): \(.name)"')
echo "$SITES"
echo ""

# Fetch parts from Kleinanzeigen
echo "3. Fetching parts from Kleinanzeigen (Site ID: 2)..."
echo "   This will search for: Mitsubishi Eclipse"
echo "   Limit: 25 parts"
echo ""

RESPONSE=$(curl -s -X POST "${API_URL}/api/parts/fetch" \
  -H "Content-Type: application/json" \
  -d '{
    "site_id": 2,
    "vehicle_type": "Car",
    "make": "Mitsubishi",
    "base_model": "Eclipse",
    "model": "D30",
    "year_from": 1995,
    "year_to": 2000,
    "limit": 25
  }')

# Check if successful
if echo "$RESPONSE" | jq -e '.error' > /dev/null 2>&1; then
    echo "❌ Error fetching parts:"
    echo "$RESPONSE" | jq '.'
    exit 1
fi

# Show results
TOTAL=$(echo "$RESPONSE" | jq -r '.total')
echo "✅ Successfully fetched $TOTAL parts from Kleinanzeigen"
echo ""

# Show first part as example
echo "4. Example part (first one):"
echo "$RESPONSE" | jq -r '.data[0] | "   ID: \(.part_id)\n   Name: \(.name)\n   Description: \(.description[0:100])...\n   URL: \(.url)\n   Location: \(.type_name)"'
echo ""

# Check database
echo "5. Checking database for Kleinanzeigen parts..."
cd api
KLEINANZEIGEN_COUNT=$(sqlite3 sqlite.db "SELECT COUNT(*) FROM parts WHERE site_id = 2;")
echo "   Total Kleinanzeigen parts in database: $KLEINANZEIGEN_COUNT"
cd ..
echo ""

echo "========================================="
echo "✅ Kleinanzeigen scraper test completed!"
echo "========================================="
echo ""
echo "View all parts in the Browse page: http://localhost:5173/browse"
echo "Filter by 'Kleinanzeigen' to see only these parts"
