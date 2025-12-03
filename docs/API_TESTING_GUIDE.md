# Products API Testing Guide

## Base URL
```
http://localhost:8080/api/v1
```

## Endpoints Overview

### 1. Create Product
- **Method**: POST
- **Endpoint**: `/api/v1/products`
- **Description**: Create a new product with price variations and addons

### 2. Get Product by ID
- **Method**: GET
- **Endpoint**: `/api/v1/products/:id`
- **Description**: Get detailed information about a specific product

### 3. List Products by Company
- **Method**: GET
- **Endpoint**: `/api/v1/products/company/:company_id`
- **Description**: Get all products for a company with optional filters
- **Query Parameters**:
  - `limit`: Number of results (default: 50, max: 100)
  - `offset`: Pagination offset (default: 0)
  - `category`: Filter by category
  - `is_available`: Filter by availability (true/false)
  - `is_addon`: Filter addons only (true/false)

### 4. List Products by Sale Point
- **Method**: GET
- **Endpoint**: `/api/v1/products/sale-point/:sale_point_id`
- **Description**: Get all products for a sale point with optional filters
- **Query Parameters**: Same as company endpoint

### 5. Update Product
- **Method**: PUT
- **Endpoint**: `/api/v1/products/:id`
- **Description**: Update an existing product (partial updates supported)

### 6. Delete Product
- **Method**: DELETE
- **Endpoint**: `/api/v1/products/:id`
- **Description**: Delete a product

### 7. Get Categories by Company
- **Method**: GET
- **Endpoint**: `/api/v1/categories/company/:company_id`
- **Description**: Get all unique categories for a company

### 8. Get Categories by Sale Point
- **Method**: GET
- **Endpoint**: `/api/v1/categories/sale-point/:sale_point_id`
- **Description**: Get all unique categories for a sale point

---

## Test Commands

### Test 1: Create Ice Cream Product (Helado)

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
  "company_id": "550e8400-e29b-41d4-a716-446655440000",
  "sale_point_id": "650e8400-e29b-41d4-a716-446655440001",
  "name": "Helado",
  "description": "Bolas de helado melo",
  "category": "Helados",
  "photos": [
    "https://example.com/photos/helado-main.jpg",
    "https://example.com/photos/helado-secondary.jpg"
  ],
  "price_variations": [
    {
      "type": "1 bola",
      "price": 30000,
      "included_addons": {
        "max_selections": 2,
        "options": [
          {
            "id": "addon-001",
            "name": "chispitas de chocolate",
            "price": 0,
            "photos": [],
            "is_available": true
          }
        ]
      }
    },
    {
      "type": "2 bolas",
      "price": 32000,
      "included_addons": {
        "max_selections": 0,
        "options": []
      }
    }
  ],
  "available_addons": [
    {
      "id": "addon-002",
      "name": "chispitas de chocolate",
      "price": 5000,
      "photos": ["https://example.com/addons/chispitas.jpg"],
      "is_available": true
    },
    {
      "id": "addon-003",
      "name": "menta",
      "price": 5000,
      "photos": [],
      "is_available": true
    },
    {
      "id": "addon-004",
      "name": "hojuelas de cereal",
      "price": 5000,
      "photos": [],
      "is_available": true
    }
  ],
  "is_addon": false,
  "is_available": true,
  "is_unlimited_stock": true,
  "stock": null
}'
```

### Test 2: Create Natural Juice Product

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
  "company_id": "550e8400-e29b-41d4-a716-446655440000",
  "sale_point_id": "650e8400-e29b-41d4-a716-446655440001",
  "name": "Jugo natural en agua",
  "description": "Jugo natural en agua",
  "category": "Jugos",
  "photos": [],
  "price_variations": [
    {
      "type": "Pequeño",
      "price": 10000,
      "included_addons": {
        "max_selections": 1,
        "options": [
          {
            "id": "flavor-001",
            "name": "NARANJA",
            "price": 0,
            "photos": [],
            "is_available": true
          },
          {
            "id": "flavor-002",
            "name": "LIMON",
            "price": 0,
            "photos": [],
            "is_available": true
          }
        ]
      }
    },
    {
      "type": "Mediano",
      "price": 12000,
      "included_addons": {
        "max_selections": 1,
        "options": [
          {
            "id": "flavor-001",
            "name": "NARANJA",
            "price": 0,
            "photos": [],
            "is_available": true
          },
          {
            "id": "flavor-002",
            "name": "LIMON",
            "price": 0,
            "photos": [],
            "is_available": true
          }
        ]
      }
    }
  ],
  "available_addons": [],
  "is_addon": false,
  "is_available": true,
  "is_unlimited_stock": true,
  "stock": null
}'
```

### Test 3: Create Fruit Salad Product

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
  "company_id": "550e8400-e29b-41d4-a716-446655440000",
  "sale_point_id": "650e8400-e29b-41d4-a716-446655440001",
  "name": "Ensalada de fruta con helado",
  "description": "Ensalada de fruta con helado",
  "category": "Ensaladas de fruta",
  "photos": ["https://example.com/photos/ensalada-fruta.jpg"],
  "price_variations": [
    {
      "type": "16 Onzas",
      "price": 10000,
      "included_addons": {
        "max_selections": 2,
        "options": [
          {
            "id": "ice-001",
            "name": "CHOCOLATE",
            "price": 0,
            "photos": [],
            "is_available": true
          },
          {
            "id": "ice-002",
            "name": "VAINILLA",
            "price": 0,
            "photos": [],
            "is_available": true
          },
          {
            "id": "ice-003",
            "name": "RON CON PASAS",
            "price": 0,
            "photos": [],
            "is_available": true
          }
        ]
      }
    },
    {
      "type": "18 Onzas",
      "price": 12000,
      "included_addons": {
        "max_selections": 2,
        "options": [
          {
            "id": "ice-001",
            "name": "CHOCOLATE",
            "price": 0,
            "photos": [],
            "is_available": true
          },
          {
            "id": "ice-002",
            "name": "VAINILLA",
            "price": 0,
            "photos": [],
            "is_available": true
          },
          {
            "id": "ice-003",
            "name": "RON CON PASAS",
            "price": 0,
            "photos": [],
            "is_available": true
          }
        ]
      }
    },
    {
      "type": "20 Onzas",
      "price": 14000,
      "included_addons": {
        "max_selections": 3,
        "options": [
          {
            "id": "ice-001",
            "name": "CHOCOLATE",
            "price": 0,
            "photos": [],
            "is_available": true
          },
          {
            "id": "ice-002",
            "name": "VAINILLA",
            "price": 0,
            "photos": [],
            "is_available": true
          },
          {
            "id": "ice-003",
            "name": "RON CON PASAS",
            "price": 0,
            "photos": [],
            "is_available": true
          }
        ]
      }
    }
  ],
  "available_addons": [
    {
      "id": "addon-ice-001",
      "name": "CHOCOLATE",
      "price": 0,
      "photos": [],
      "is_available": true
    },
    {
      "id": "addon-ice-002",
      "name": "VAINILLA",
      "price": 0,
      "photos": [],
      "is_available": true
    },
    {
      "id": "addon-ice-003",
      "name": "RON CON PASAS",
      "price": 0,
      "photos": [],
      "is_available": true
    },
    {
      "id": "addon-extra-001",
      "name": "LECHERITA",
      "price": 3000,
      "photos": [],
      "is_available": true
    }
  ],
  "is_addon": false,
  "is_available": true,
  "is_unlimited_stock": true,
  "stock": null
}'
```

### Test 4: Get Product by ID

```bash
# First, save the product ID from the create response
# Then use it to get the product

curl -X GET http://localhost:8080/api/v1/products/{PRODUCT_ID}
```

### Test 5: List Products by Company

```bash
# Get all products for a company
curl -X GET "http://localhost:8080/api/v1/products/company/550e8400-e29b-41d4-a716-446655440000"

# Get products with pagination
curl -X GET "http://localhost:8080/api/v1/products/company/550e8400-e29b-41d4-a716-446655440000?limit=10&offset=0"

# Filter by category
curl -X GET "http://localhost:8080/api/v1/products/company/550e8400-e29b-41d4-a716-446655440000?category=Helados"

# Filter by availability
curl -X GET "http://localhost:8080/api/v1/products/company/550e8400-e29b-41d4-a716-446655440000?is_available=true"

# Combined filters
curl -X GET "http://localhost:8080/api/v1/products/company/550e8400-e29b-41d4-a716-446655440000?category=Jugos&is_available=true&limit=20"
```

### Test 6: List Products by Sale Point

```bash
# Get all products for a sale point
curl -X GET "http://localhost:8080/api/v1/products/sale-point/650e8400-e29b-41d4-a716-446655440001"

# With filters
curl -X GET "http://localhost:8080/api/v1/products/sale-point/650e8400-e29b-41d4-a716-446655440001?category=Ensaladas%20de%20fruta&is_available=true"
```

### Test 7: Get Categories

```bash
# Get categories by company
curl -X GET "http://localhost:8080/api/v1/categories/company/550e8400-e29b-41d4-a716-446655440000"

# Get categories by sale point
curl -X GET "http://localhost:8080/api/v1/categories/sale-point/650e8400-e29b-41d4-a716-446655440001"
```

### Test 8: Update Product

```bash
# Update product availability
curl -X PUT http://localhost:8080/api/v1/products/{PRODUCT_ID} \
  -H "Content-Type: application/json" \
  -d '{
  "is_available": false
}'

# Update product name and description
curl -X PUT http://localhost:8080/api/v1/products/{PRODUCT_ID} \
  -H "Content-Type: application/json" \
  -d '{
  "name": "Helado Premium",
  "description": "Bolas de helado premium artesanal"
}'

# Update price variations
curl -X PUT http://localhost:8080/api/v1/products/{PRODUCT_ID} \
  -H "Content-Type: application/json" \
  -d '{
  "price_variations": [
    {
      "type": "1 bola",
      "price": 35000,
      "included_addons": {
        "max_selections": 2,
        "options": []
      }
    },
    {
      "type": "2 bolas",
      "price": 37000,
      "included_addons": {
        "max_selections": 0,
        "options": []
      }
    }
  ]
}'
```

### Test 9: Delete Product

```bash
curl -X DELETE http://localhost:8080/api/v1/products/{PRODUCT_ID}
```

### Test 10: Health Check

```bash
curl -X GET http://localhost:8080/health
```

---

## Expected Response Formats

### Success Response
```json
{
  "success": true,
  "data": { ... },
  "message": "..."
}
```

### Error Response
```json
{
  "success": false,
  "error": "error details",
  "message": "user-friendly message"
}
```

### Product List Response (Simplified for lists)
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "name": "Product Name",
      "photos": ["url1", "url2"],
      "category": "Category Name",
      "min_price": 10000,
      "is_available": true
    }
  ]
}
```

### Product Detail Response (Full structure)
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "company_id": "uuid",
    "sale_point_id": "uuid",
    "name": "Product Name",
    "photos": ["url1"],
    "price_variations": [...],
    "category": "Category",
    "description": "Description",
    "is_addon": false,
    "is_available": true,
    "is_unlimited_stock": true,
    "stock": null,
    "available_addons": [...],
    "created_at": "2025-12-01T...",
    "updated_at": "2025-12-01T..."
  }
}
```

---

## Notes

1. **Price Format**: All prices are in cents/smallest currency unit (e.g., 30000 = $300.00 or 30000 COP)
2. **Stock Management**: 
   - If `is_unlimited_stock` is `true`, `stock` must be `null`
   - If `is_unlimited_stock` is `false`, `stock` must be a number >= 0
3. **UUIDs**: The system generates UUIDs automatically for products
4. **Addons**: Can have their own IDs for reference in orders
5. **Filtering**: All filters are optional and can be combined
6. **Pagination**: Default limit is 50, maximum is 100

---

## Testing Workflow

1. **Create test data**: Use Tests 1-3 to create sample products
2. **Save product IDs**: Note the IDs returned from create operations
3. **Test listing**: Use Tests 5-6 to list products with various filters
4. **Test categories**: Use Test 7 to verify categories are extracted correctly
5. **Test details**: Use Test 4 to get full product information
6. **Test updates**: Use Test 8 to modify products
7. **Test deletion**: Use Test 9 to remove products
8. **Verify changes**: Re-run list queries to confirm updates/deletions

---

## Common Use Cases

### Frontend Product List View
```bash
# Get products for a sale point with only available items
curl -X GET "http://localhost:8080/api/v1/products/sale-point/{SALE_POINT_ID}?is_available=true&limit=20"
```

### Frontend Product Detail View
```bash
# Get full product details including all variations and addons
curl -X GET "http://localhost:8080/api/v1/products/{PRODUCT_ID}"
```

### Frontend Category Filter
```bash
# Get all categories for filtering
curl -X GET "http://localhost:8080/api/v1/categories/sale-point/{SALE_POINT_ID}"

# Then filter by selected category
curl -X GET "http://localhost:8080/api/v1/products/sale-point/{SALE_POINT_ID}?category=Helados&is_available=true"
```

### Admin: Toggle Product Availability
```bash
curl -X PUT http://localhost:8080/api/v1/products/{PRODUCT_ID} \
  -H "Content-Type: application/json" \
  -d '{"is_available": false}'
```
