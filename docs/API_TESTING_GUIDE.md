# API Testing Guide

## Base URL
```
http://localhost:8080/api/v1
```

---

# Products API

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

---

# Orders API

## Endpoints Overview

### 1. Create Order
- **Method**: POST
- **Endpoint**: `/api/v1/orders`
- **Description**: Create a new order (DELIVERY or ON_SITE)
- **Returns**: 201 Created with order code for tracking

### 2. Track Order (Public)
- **Method**: GET
- **Endpoint**: `/api/v1/orders/track/:code`
- **Description**: Public tracking endpoint with limited information
- **Auth**: None required

### 3. Partial Update Order
- **Method**: PATCH
- **Endpoint**: `/api/v1/orders`
- **Description**: Update status, notes, payment (NO products allowed)

### 4. Modify Order
- **Method**: PUT
- **Endpoint**: `/api/v1/orders`
- **Description**: Full modification including products (auto-sets status to VERIFIED)

### 5. List Orders
- **Method**: GET
- **Endpoint**: `/api/v1/orders`
- **Description**: Get all orders with advanced filters
- **Query Parameters**:
  - `limit`: Number of results (default: 50, max: 100)
  - `offset`: Pagination offset (default: 0)
  - `date_from`: Filter from date (RFC3339 format)
  - `date_to`: Filter to date (RFC3339 format)
  - `status`: Filter by status (CREATED, VERIFIED, IN_PROGRESS, OUT_FOR_DELIVERY, DELIVERED, CANCELLED)
  - `sale_type`: Filter by sale type (DELIVERY, ON_SITE)
  - `product_id`: Filter by product ID
  - `product_name`: Filter by product name (partial match)
  - `min_total`: Minimum total amount (in cents)
  - `max_total`: Maximum total amount (in cents)

### 6. Get Order Metrics
- **Method**: GET
- **Endpoint**: `/api/v1/orders/metrics`
- **Description**: Get analytics and aggregated metrics
- **Query Parameters**: Same filters as list orders

### 7. Get Order by Code (Admin)
- **Method**: GET
- **Endpoint**: `/api/v1/orders/:code`
- **Description**: Get full order details (admin/internal use)

---

## Order Status Lifecycle

```
CREATED → VERIFIED → IN_PROGRESS → OUT_FOR_DELIVERY → DELIVERED
    ↓         ↓            ↓                ↓
CANCELLED  CANCELLED   CANCELLED       CANCELLED
```

**Valid Transitions:**
- CREATED → VERIFIED, IN_PROGRESS, CANCELLED
- VERIFIED → IN_PROGRESS, CANCELLED
- IN_PROGRESS → OUT_FOR_DELIVERY, DELIVERED (ON_SITE), CANCELLED
- OUT_FOR_DELIVERY → DELIVERED, CANCELLED
- DELIVERED → (terminal state)
- CANCELLED → (terminal state)

---

## Test Commands

### Test 1: Create DELIVERY Order (Full Example)

```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
  "sale_type": "DELIVERY",
  "products": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "limonada",
      "description": "HELADO: VAINILLA - FRESA",
      "observation": "sin azucar",
      "price": 10000,
      "quantity": 2
    },
    {
      "id": "660e8400-e29b-41d4-a716-446655440001",
      "name": "Ensalada de fruta 1",
      "description": "HELADO: VAINILLA - FRESA",
      "observation": "SIN UVAS",
      "price": 19900,
      "quantity": 1
    }
  ],
  "note": "Por favor entregar antes de las 3pm",
  "customer": {
    "identification": "3827994902",
    "id_type": "CC",
    "name": "CARLOS ARTURO MARTINEZ",
    "phone": "+573002003399"
  },
  "shipping_address": "Calle 123 #45-67, Apartamento 301, Bogotá",
  "payment_receipt_url": "https://example.com/receipts/payment123.jpg",
  "payment_account_id": "770e8400-e29b-41d4-a716-446655440002"
}'
```

**Expected Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": "880e8400-e29b-41d4-a716-446655440003",
    "code": "ORD-1735776000123456789-ab12cd34",
    "status": "CREATED",
    "sale_type": "DELIVERY",
    "total": 39900,
    "created_at": "2026-01-01T18:00:00Z",
    "updated_at": "2026-01-01T18:00:00Z"
  },
  "message": "Order created successfully"
}
```

**⚠️ SAVE THE ORDER CODE** - You'll need it for subsequent tests!

### Test 2: Create ON_SITE Order (Restaurant Table)

```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
  "sale_type": "ON_SITE",
  "products": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "Café Americano",
      "price": 5000,
      "quantity": 2
    },
    {
      "id": "660e8400-e29b-41d4-a716-446655440001",
      "name": "Croissant",
      "price": 8000,
      "quantity": 1
    }
  ],
  "table_number": 5,
  "note": "Mesa cerca de la ventana"
}'
```

**Expected Response:**
```json
{
  "success": true,
  "data": {
    "id": "...",
    "code": "ORD-1735776120000000000-xy45zw89",
    "status": "CREATED",
    "sale_type": "ON_SITE",
    "total": 18000,
    "created_at": "...",
    "updated_at": "..."
  },
  "message": "Order created successfully"
}
```

### Test 3: Track Order (Public - No Auth)

```bash
# Replace with your actual order code from Test 1 or 2
curl -X GET "http://localhost:8080/api/v1/orders/track/ORD-1735776000123456789-ab12cd34"
```

**Expected Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "code": "ORD-1735776000123456789-ab12cd34",
    "status": "CREATED",
    "customer_name": "CARLOS ARTURO MARTINEZ",
    "updated_at": "2026-01-01T18:00:00Z"
  }
}
```

**Note:** This endpoint only returns:
- Order code
- Current status
- Customer name (NOT identification or phone)
- Last update timestamp

### Test 4: Partial Update - Change Status to VERIFIED

```bash
curl -X PATCH http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
  "code": "ORD-1735776000123456789-ab12cd34",
  "status": "VERIFIED"
}'
```

**Expected Response (200 OK):**
Returns full order with status updated to "VERIFIED"

### Test 5: Partial Update - Add Note and Payment Info

```bash
curl -X PATCH http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
  "code": "ORD-1735776000123456789-ab12cd34",
  "note": "Cliente confirmó la dirección por teléfono",
  "payment_receipt_url": "https://example.com/receipts/updated-payment.jpg",
  "payment_account_id": "990e8400-e29b-41d4-a716-446655440099"
}'
```

### Test 6: Partial Update - Status Transition to IN_PROGRESS

```bash
curl -X PATCH http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
  "code": "ORD-1735776000123456789-ab12cd34",
  "status": "IN_PROGRESS"
}'
```

### Test 7: Try to Update Products via PATCH (Should FAIL with 400)

```bash
# This should return an error
curl -X PATCH http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
  "code": "ORD-1735776000123456789-ab12cd34",
  "products": [
    {
      "id": "test-product",
      "name": "Test",
      "price": 1000,
      "quantity": 1
    }
  ]
}'
```

**Expected Error (400 Bad Request):**
```json
{
  "success": false,
  "error": "products cannot be updated via PATCH, use PUT instead",
  "message": "Products cannot be updated via PATCH, use PUT instead"
}
```

### Test 8: Modify Order - Change Products (PUT)

```bash
# This WILL update products and automatically set status to VERIFIED
curl -X PUT http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
  "code": "ORD-1735776000123456789-ab12cd34",
  "products": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "limonada",
      "price": 10000,
      "quantity": 3
    },
    {
      "id": "770e8400-e29b-41d4-a716-446655440002",
      "name": "Jugo de naranja",
      "price": 12000,
      "quantity": 1
    }
  ],
  "shipping_address": "Nueva dirección: Carrera 7 #12-34, Apartamento 501",
  "note": "Cliente cambió pedido - agregó una limonada más"
}'
```

**Expected Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": "...",
    "code": "ORD-1735776000123456789-ab12cd34",
    "status": "VERIFIED",
    "sale_type": "DELIVERY",
    "products": [...],
    "total": 42000,
    "note": "Cliente cambió pedido - agregó una limonada más",
    "customer": {...},
    "shipping_address": "Nueva dirección: Carrera 7 #12-34, Apartamento 501",
    "created_at": "...",
    "updated_at": "..."
  },
  "message": "Order modified successfully"
}
```

**Note:** Status is automatically set to VERIFIED and total is recalculated

### Test 9: Update Status to OUT_FOR_DELIVERY

```bash
curl -X PATCH http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
  "code": "ORD-1735776000123456789-ab12cd34",
  "status": "OUT_FOR_DELIVERY"
}'
```

### Test 10: Try to Modify Order when OUT_FOR_DELIVERY (Should FAIL with 409)

```bash
# This should return an error
curl -X PUT http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
  "code": "ORD-1735776000123456789-ab12cd34",
  "products": [
    {
      "id": "test",
      "name": "test",
      "price": 1000,
      "quantity": 1
    }
  ]
}'
```

**Expected Error (409 Conflict):**
```json
{
  "success": false,
  "error": "order cannot be modified in current status",
  "message": "Failed to modify order"
}
```

### Test 11: Update Status to DELIVERED

```bash
curl -X PATCH http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
  "code": "ORD-1735776000123456789-ab12cd34",
  "status": "DELIVERED"
}'
```

### Test 12: Try Invalid Status Transition (Should FAIL with 409)

```bash
# Try to transition from DELIVERED to IN_PROGRESS (invalid)
curl -X PATCH http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
  "code": "ORD-1735776000123456789-ab12cd34",
  "status": "IN_PROGRESS"
}'
```

**Expected Error (409 Conflict):**
```json
{
  "success": false,
  "error": "invalid status transition",
  "message": "Failed to update order"
}
```

### Test 13: Cancel Order

```bash
# Create a new order first, then cancel it
curl -X PATCH http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
  "code": "ORD-1735776120000000000-xy45zw89",
  "status": "CANCELLED",
  "note": "Cliente canceló el pedido por teléfono"
}'
```

### Test 14: Get Order by Code (Admin - Full Details)

```bash
curl -X GET "http://localhost:8080/api/v1/orders/ORD-1735776000123456789-ab12cd34"
```

**Expected Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": "880e8400-e29b-41d4-a716-446655440003",
    "code": "ORD-1735776000123456789-ab12cd34",
    "status": "DELIVERED",
    "sale_type": "DELIVERY",
    "products": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "name": "limonada",
        "description": "HELADO: VAINILLA - FRESA",
        "observation": "sin azucar",
        "price": 10000,
        "quantity": 3
      },
      {
        "id": "770e8400-e29b-41d4-a716-446655440002",
        "name": "Jugo de naranja",
        "price": 12000,
        "quantity": 1
      }
    ],
    "total": 42000,
    "note": "Cliente cambió pedido - agregó una limonada más",
    "customer": {
      "identification": "3827994902",
      "id_type": "CC",
      "name": "CARLOS ARTURO MARTINEZ",
      "phone": "+573002003399"
    },
    "shipping_address": "Nueva dirección: Carrera 7 #12-34, Apartamento 501",
    "payment_receipt_url": "https://example.com/receipts/updated-payment.jpg",
    "payment_account_id": "990e8400-e29b-41d4-a716-446655440099",
    "created_at": "2026-01-01T18:00:00Z",
    "updated_at": "2026-01-01T19:30:00Z"
  }
}
```

### Test 15: List All Orders

```bash
curl -X GET "http://localhost:8080/api/v1/orders"
```

**Expected Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "id": "...",
      "code": "...",
      "status": "...",
      "sale_type": "...",
      "products": [...],
      "total": 42000,
      ...
    }
  ],
  "meta": {
    "current_page": 1,
    "total_pages": 1,
    "total_items": 2,
    "page_size": 50
  }
}
```

### Test 16: List Orders with Pagination

```bash
# Get first 10 orders
curl -X GET "http://localhost:8080/api/v1/orders?limit=10&offset=0"

# Get next 10 orders
curl -X GET "http://localhost:8080/api/v1/orders?limit=10&offset=10"
```

### Test 17: Filter Orders by Status

```bash
# Get all DELIVERED orders
curl -X GET "http://localhost:8080/api/v1/orders?status=DELIVERED"

# Get all IN_PROGRESS orders
curl -X GET "http://localhost:8080/api/v1/orders?status=IN_PROGRESS"

# Get all CANCELLED orders
curl -X GET "http://localhost:8080/api/v1/orders?status=CANCELLED"
```

### Test 18: Filter Orders by Sale Type

```bash
# Get all DELIVERY orders
curl -X GET "http://localhost:8080/api/v1/orders?sale_type=DELIVERY"

# Get all ON_SITE orders
curl -X GET "http://localhost:8080/api/v1/orders?sale_type=ON_SITE"
```

### Test 19: Filter Orders by Date Range

```bash
# Get orders from January 1, 2026
curl -X GET "http://localhost:8080/api/v1/orders?date_from=2026-01-01T00:00:00Z&date_to=2026-01-01T23:59:59Z"

# Get orders from last 7 days
curl -X GET "http://localhost:8080/api/v1/orders?date_from=2025-12-25T00:00:00Z&date_to=2026-01-01T23:59:59Z"
```

### Test 20: Filter Orders by Product

```bash
# Filter by product ID
curl -X GET "http://localhost:8080/api/v1/orders?product_id=550e8400-e29b-41d4-a716-446655440000"

# Filter by product name (partial match)
curl -X GET "http://localhost:8080/api/v1/orders?product_name=limonada"
```

### Test 21: Filter Orders by Total Range

```bash
# Get orders with total between 20000 and 50000 cents
curl -X GET "http://localhost:8080/api/v1/orders?min_total=20000&max_total=50000"

# Get orders over 100000 cents
curl -X GET "http://localhost:8080/api/v1/orders?min_total=100000"
```

### Test 22: Combined Filters

```bash
# Get DELIVERY orders with DELIVERED status from January, paginated
curl -X GET "http://localhost:8080/api/v1/orders?sale_type=DELIVERY&status=DELIVERED&date_from=2026-01-01T00:00:00Z&date_to=2026-01-31T23:59:59Z&limit=20"

# Get orders containing a specific product with minimum total
curl -X GET "http://localhost:8080/api/v1/orders?product_name=limonada&min_total=30000&status=DELIVERED"
```

### Test 23: Get Order Metrics (All Time)

```bash
curl -X GET "http://localhost:8080/api/v1/orders/metrics"
```

**Expected Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "metrics": {
      "total_sales": 1200000,
      "avg_ticket": 35000,
      "orders_by_status": {
        "CREATED": 5,
        "VERIFIED": 3,
        "IN_PROGRESS": 8,
        "OUT_FOR_DELIVERY": 2,
        "DELIVERED": 25,
        "CANCELLED": 2
      }
    },
    "top_products": [
      {
        "product_id": "550e8400-e29b-41d4-a716-446655440000",
        "name": "limonada",
        "total_quantity": 45,
        "total_revenue": 450000
      },
      {
        "product_id": "660e8400-e29b-41d4-a716-446655440001",
        "name": "Ensalada de fruta 1",
        "total_quantity": 30,
        "total_revenue": 597000
      }
    ]
  }
}
```

### Test 24: Get Metrics with Date Filter

```bash
# Get metrics for January 2026
curl -X GET "http://localhost:8080/api/v1/orders/metrics?date_from=2026-01-01T00:00:00Z&date_to=2026-01-31T23:59:59Z"
```

### Test 25: Get Metrics for Specific Status

```bash
# Get metrics for DELIVERED orders only
curl -X GET "http://localhost:8080/api/v1/orders/metrics?status=DELIVERED"
```

### Test 26: Get Metrics by Sale Type

```bash
# Get metrics for DELIVERY orders
curl -X GET "http://localhost:8080/api/v1/orders/metrics?sale_type=DELIVERY"

# Get metrics for ON_SITE orders
curl -X GET "http://localhost:8080/api/v1/orders/metrics?sale_type=ON_SITE"
```

---

## Validation Error Examples

### Invalid Enum Value (400 Bad Request)

```bash
# Invalid status value
curl -X PATCH http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
  "code": "ORD-1735776000123456789-ab12cd34",
  "status": "VERIFIEDed"
}'
```

**Expected Error:**
```json
{
  "success": false,
  "error": "Validation failed for field 'status': 'status' must be one of: CREATED, VERIFIED, IN_PROGRESS, OUT_FOR_DELIVERY, DELIVERED, CANCELLED",
  "message": "Validation failed",
  "details": [
    {
      "field": "status",
      "message": "'status' must be one of: CREATED, VERIFIED, IN_PROGRESS, OUT_FOR_DELIVERY, DELIVERED, CANCELLED"
    }
  ]
}
```

### Multiple Validation Errors (400 Bad Request)

```bash
# Invalid sale type and empty products
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
  "sale_type": "INVALID_TYPE",
  "products": []
}'
```

**Expected Error:**
```json
{
  "success": false,
  "error": "Validation failed for 2 field(s)",
  "message": "Validation failed",
  "details": [
    {
      "field": "sale_type",
      "message": "'sale_type' must be one of: DELIVERY, ON_SITE"
    },
    {
      "field": "products",
      "message": "'products' must contain at least 1 items"
    }
  ]
}
```

### Missing Required Fields (400 Bad Request)

```bash
# Missing products field
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
  "sale_type": "DELIVERY"
}'
```

**Expected Error:**
```json
{
  "success": false,
  "error": "Validation failed for field 'products': 'products' is required",
  "message": "Validation failed",
  "details": [
    {
      "field": "products",
      "message": "'products' is required"
    }
  ]
}
```

### Business Rule Violations (422 Unprocessable Entity)

Note: Business rule errors (422) have a different format than validation errors (400).

### Missing Required Fields (422 Unprocessable Entity)

```bash
# Missing customer for DELIVERY
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
  "sale_type": "DELIVERY",
  "products": [
    {"id": "test", "name": "test", "price": 1000, "quantity": 1}
  ],
  "shipping_address": "Test Address"
}'
```

**Expected Error:**
```json
{
  "success": false,
  "error": "validation error: customer information is required for delivery orders",
  "message": "Failed to create order"
}
```

### Missing Shipping Address for DELIVERY (422)

```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
  "sale_type": "DELIVERY",
  "products": [
    {"id": "test", "name": "test", "price": 1000, "quantity": 1}
  ],
  "customer": {
    "identification": "123456",
    "id_type": "CC",
    "name": "John Doe",
    "phone": "+573001234567"
  }
}'
```

**Expected Error:**
```json
{
  "success": false,
  "error": "validation error: shipping address is required for delivery orders",
  "message": "Failed to create order"
}
```

### Missing Table Number for ON_SITE (422)

```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
  "sale_type": "ON_SITE",
  "products": [
    {"id": "test", "name": "test", "price": 1000, "quantity": 1}
  ]
}'
```

**Expected Error:**
```json
{
  "success": false,
  "error": "validation error: table number is required for on-site orders",
  "message": "Failed to create order"
}
```

### No Products (422)

```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
  "sale_type": "ON_SITE",
  "products": [],
  "table_number": 5
}'
```

**Expected Error:**
```json
{
  "success": false,
  "error": "validation error: at least one product is required",
  "message": "Failed to create order"
}
```

### Order Not Found (404)

```bash
curl -X GET "http://localhost:8080/api/v1/orders/track/INVALID-CODE"
```

**Expected Error:**
```json
{
  "success": false,
  "error": "order not found",
  "message": "Order not found"
}
```

---

## Testing Workflow

### Complete Order Lifecycle Test

1. **Create Order** (Test 1)
   - Save the returned order code

2. **Track Order** (Test 3)
   - Verify public tracking works
   - Confirm limited data exposure

3. **Update to VERIFIED** (Test 4)
   - Confirm status transition works

4. **Modify Products** (Test 8)
   - Verify products can be changed
   - Confirm status auto-changes to VERIFIED
   - Verify total is recalculated

5. **Progress Order** (Tests 6, 9)
   - IN_PROGRESS → OUT_FOR_DELIVERY → DELIVERED

6. **Get Full Details** (Test 14)
   - Verify all fields are correct

7. **View in Metrics** (Test 23)
   - Confirm order appears in analytics

### Business Rules Validation

1. **PATCH vs PUT** (Tests 7, 8)
   - Confirm PATCH rejects products
   - Confirm PUT accepts products

2. **Status Transitions** (Test 12)
   - Try invalid transitions
   - Verify error responses

3. **Modification Constraints** (Test 10)
   - Try to modify OUT_FOR_DELIVERY order
   - Verify rejection

4. **Sale Type Validation** (Validation Examples)
   - Test DELIVERY requirements
   - Test ON_SITE requirements

---

## Common Use Cases

### Customer: Track My Order
```bash
# Customer receives order code via email/SMS
curl -X GET "http://localhost:8080/api/v1/orders/track/ORD-1735776000123456789-ab12cd34"
```

### Admin: View Today's Orders
```bash
curl -X GET "http://localhost:8080/api/v1/orders?date_from=2026-01-01T00:00:00Z&date_to=2026-01-01T23:59:59Z"
```

### Kitchen: Get IN_PROGRESS Orders
```bash
curl -X GET "http://localhost:8080/api/v1/orders?status=IN_PROGRESS"
```

### Delivery: Get OUT_FOR_DELIVERY Orders
```bash
curl -X GET "http://localhost:8080/api/v1/orders?status=OUT_FOR_DELIVERY&sale_type=DELIVERY"
```

### Manager: Daily Sales Report
```bash
curl -X GET "http://localhost:8080/api/v1/orders/metrics?date_from=2026-01-01T00:00:00Z&date_to=2026-01-01T23:59:59Z&status=DELIVERED"
```

### Support: Find Orders by Customer Phone
```bash
# Note: This requires adding the order code lookup or search functionality
curl -X GET "http://localhost:8080/api/v1/orders/ORD-CUSTOMER-CODE"
```

### Analytics: Best Selling Products
```bash
curl -X GET "http://localhost:8080/api/v1/orders/metrics?status=DELIVERED"
```

### Admin: Cancel Customer Request
```bash
curl -X PATCH http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
  "code": "ORD-XXXXX",
  "status": "CANCELLED",
  "note": "Customer requested cancellation via phone call"
}'
```

---

## Notes

1. **Price Format**: All prices are in cents/smallest currency unit (e.g., 10000 = $100.00 or 10000 COP)
2. **Order Codes**: Auto-generated format: `ORD-{nanosecond-timestamp}-{uuid8}`
3. **Total Calculation**: Backend automatically calculates total from products (sum of price × quantity)
4. **Status Transitions**: Only valid transitions allowed - invalid ones return 409 Conflict
5. **Modification Rules**:
   - PATCH: Status, notes, payment only (NO products)
   - PUT: All fields including products (auto-sets VERIFIED status)
   - Cannot modify: OUT_FOR_DELIVERY, DELIVERED, or CANCELLED orders
6. **Public Tracking**: Limited to code, status, customer_name, updated_at only
7. **Pagination**: Default limit is 50, maximum is 100
8. **Date Filters**: Use RFC3339 format (e.g., `2026-01-01T00:00:00Z`)
9. **Metrics**: Aggregated from all orders matching filters, includes top 10 products

---

## Quick Reference: HTTP Status Codes

| Code | Meaning | When |
|------|---------|------|
| 200 | OK | Successful GET, PATCH, PUT |
| 201 | Created | Successful POST (order created) |
| 400 | Bad Request | Products in PATCH, invalid JSON |
| 404 | Not Found | Order code not found |
| 409 | Conflict | Invalid status transition, cannot modify |
| 422 | Unprocessable Entity | Business rule violation, validation error |
| 500 | Internal Server Error | Database error, server error |

