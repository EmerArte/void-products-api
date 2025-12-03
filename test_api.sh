#!/bin/bash

BASE_URL="http://localhost:8080"

echo "🧪 Testing Products API"
echo "======================"

# Health check
echo -e "\n1️⃣  Health Check:"
curl -s "$BASE_URL/health" | jq '.'

# Create product
echo -e "\n2️⃣  Creating a product:"
RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/products" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "MacBook Pro",
    "description": "High-performance laptop for developers",
    "price": 2499.99,
    "stock": 25,
    "category": "Electronics"
  }')
echo "$RESPONSE" | jq '.'
PRODUCT_ID=$(echo "$RESPONSE" | jq -r '.data.id')

# Get all products
echo -e "\n3️⃣  Getting all products:"
curl -s "$BASE_URL/api/v1/products?limit=5&offset=0" | jq '.'

# Get product by ID
if [ ! -z "$PRODUCT_ID" ]; then
  echo -e "\n4️⃣  Getting product by ID ($PRODUCT_ID):"
  curl -s "$BASE_URL/api/v1/products/$PRODUCT_ID" | jq '.'

  # Update product
  echo -e "\n5️⃣  Updating product:"
  curl -s -X PUT "$BASE_URL/api/v1/products/$PRODUCT_ID" \
    -H "Content-Type: application/json" \
    -d '{
      "price": 2299.99,
      "stock": 30
    }' | jq '.'

  # Delete product
  echo -e "\n6️⃣  Deleting product:"
  curl -s -X DELETE "$BASE_URL/api/v1/products/$PRODUCT_ID" | jq '.'
fi

echo -e "\n✅ API tests completed!"
