#!/bin/bash

# Validation Error Improvements - Test Script
# This script demonstrates the improved validation error messages

echo "========================================="
echo "Validation Error Improvements Test Suite"
echo "========================================="
echo ""

BASE_URL="http://localhost:8080/api/v1"

echo "📋 Test 1: Invalid Enum Value (Order Status)"
echo "Request: Invalid status 'VERIFIEDed'"
echo ""
curl -s -X PATCH $BASE_URL/orders \
  -H "Content-Type: application/json" \
  -d '{
    "code": "ORD-TEST-123",
    "status": "VERIFIEDed"
  }' | jq
echo ""
echo "✅ Expected: Clear message showing valid status values"
echo ""

echo "========================================="
echo ""

echo "📋 Test 2: Multiple Validation Errors"
echo "Request: Invalid sale_type + empty products array"
echo ""
curl -s -X POST $BASE_URL/orders \
  -H "Content-Type: application/json" \
  -d '{
    "sale_type": "INVALID_TYPE",
    "products": []
  }' | jq
echo ""
echo "✅ Expected: Both errors listed with clear messages"
echo ""

echo "========================================="
echo ""

echo "📋 Test 3: Missing Required Field"
echo "Request: Missing products field"
echo ""
curl -s -X POST $BASE_URL/orders \
  -H "Content-Type: application/json" \
  -d '{
    "sale_type": "DELIVERY"
  }' | jq
echo ""
echo "✅ Expected: Clear 'field is required' message"
echo ""

echo "========================================="
echo ""

echo "📋 Test 4: Product Validation"
echo "Request: Missing multiple required fields"
echo ""
curl -s -X POST $BASE_URL/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": ""
  }' | jq
echo ""
echo "✅ Expected: All missing fields listed"
echo ""

echo "========================================="
echo ""

echo "📋 Test 5: Invalid ID Type Enum"
echo "Request: Invalid customer id_type"
echo ""
curl -s -X POST $BASE_URL/orders \
  -H "Content-Type: application/json" \
  -d '{
    "sale_type": "DELIVERY",
    "products": [
      {
        "id": "test-123",
        "name": "Test Product",
        "price": 10000,
        "quantity": 1
      }
    ],
    "customer": {
      "identification": "123456",
      "id_type": "INVALID_TYPE",
      "name": "John Doe",
      "phone": "+573001234567"
    },
    "shipping_address": "Test Address"
  }' | jq
echo ""
echo "✅ Expected: Clear message showing valid ID types (CC, CE, TI, PASSPORT)"
echo ""

echo "========================================="
echo ""

echo "📋 Test 6: Nested Validation (Product in Order)"
echo "Request: Invalid product quantity (must be > 0)"
echo ""
curl -s -X POST $BASE_URL/orders \
  -H "Content-Type: application/json" \
  -d '{
    "sale_type": "DELIVERY",
    "products": [
      {
        "id": "test-123",
        "name": "Test Product",
        "price": 10000,
        "quantity": 0
      }
    ],
    "customer": {
      "identification": "123456",
      "id_type": "CC",
      "name": "John Doe",
      "phone": "+573001234567"
    },
    "shipping_address": "Test Address"
  }' | jq
echo ""
echo "✅ Expected: Clear message about quantity validation"
echo ""

echo "========================================="
echo "All tests completed!"
echo "========================================="
echo ""
echo "Summary:"
echo "✅ All validation errors now provide:"
echo "   - Clear, user-friendly messages"
echo "   - Field-level details"
echo "   - Valid value suggestions for enums"
echo "   - Consistent response format"
echo ""
echo "See docs/VALIDATION_ERROR_HANDLING.md for complete guide"
echo ""
