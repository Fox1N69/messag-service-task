#!/bin/bash

BASE_URL="http://localhost:4000"

# Test Endpoint 1: Add Client (POST)
echo "Testing Endpoint: Add Client"
hey -m POST -H 'Content-Type: application/json' -D '{"client_name": "Test Client", "version": 1, "image": "test_image", "cpu": "2x Intel Xeon", "memory": "16GB", "priority": 0.75, "need_restart": false}' -c 10 -z 200ms -q 100 -n 1000 ${BASE_URL}/api/client/add
echo ""

# Test Endpoint 2: Update Client by ID (PATCH)
echo "Testing Endpoint: Update Client by ID"
# Replace :id with an actual client ID, here we assume it is 123 for demonstration
hey -m PATCH -H 'Content-Type: application/json' -d '{"client_name": "Updated Client Name"}' -c 10 -z 10s -q 100 ${BASE_URL}/api/client/123
echo ""

# Test Endpoint 3: Delete Client by ID (DELETE)
echo "Testing Endpoint: Delete Client by ID"
# Replace :id with an actual client ID, here we assume it is 123 for demonstration
hey -m DELETE -c 10 -z 10s -q 100 ${BASE_URL}/api/client/123
echo ""

# Test Endpoint 4: Update Algorithm Status by ID (PATCH)
echo "Testing Endpoint: Update Algorithm Status by ID"
# Replace :id with an actual algorithm ID, here we assume it is 456 for demonstration
hey -m PATCH -H 'Content-Type: application/json' -d '{"vwap": true}' -c 10 -z 10s -q 100 ${BASE_URL}/api/client/algorithm/456
echo ""

# Test Endpoint 5: Create Algorithm (POST)
echo "Testing Endpoint: Create Algorithm"
hey -m POST -H 'Content-Type: application/json' -d '{"client_id": 123, "vwap": false, "twap": true, "hft": true}' -c 10 -z 10s -q 100 ${BASE_URL}/api/client/algorithm/create
echo ""