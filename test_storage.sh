#!/bin/bash

echo "=== Testing Promptly Storage Backends ==="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Ensure data directory exists
mkdir -p data

echo -e "${YELLOW}1. Creating SQLite database...${NC}"
if [ -f "data/promptly.db" ]; then
    echo "⚠️  Database already exists at data/promptly.db"
    read -p "Do you want to overwrite it? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo "Backing up existing database..."
        mv data/promptly.db data/promptly.db.backup.$(date +%Y%m%d_%H%M%S)
        sqlite3 data/promptly.db < schema.sql
        echo "✓ Database created (backup saved)"
    else
        echo "✓ Using existing database"
    fi
else
    sqlite3 data/promptly.db < schema.sql
    echo "✓ Database created"
fi

echo -e "${YELLOW}2. Testing SQLite storage...${NC}"
go test ./internal/storage/sqlite -v

echo -e "${YELLOW}3. Testing JSON storage...${NC}"
go test ./internal/storage/jsonstore -v

echo -e "${YELLOW}4. Testing CLI help...${NC}"
./promptly serve --help

echo -e "${YELLOW}5. Available commands to test manually:${NC}"
echo "# Start with SQLite:"
echo "./promptly serve --storage=sqlite --db=data/promptly.db"
echo ""
echo "# Start with JSON (default):"
echo "./promptly serve --storage=json --data=data/prompts.json"
echo ""
echo "# Test API (run in another terminal):"
echo "curl http://localhost:8080/v1/personas"
echo "curl http://localhost:8080/health"

echo -e "${GREEN}=== Tests complete! ===${NC}"