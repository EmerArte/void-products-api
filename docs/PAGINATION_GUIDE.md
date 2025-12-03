# Pagination Guide

## Overview

All product list endpoints now return proper pagination metadata in the response, including:
- `current_page`: Current page number (1-indexed)
- `total_pages`: Total number of pages available
- `total_items`: Total number of items matching the query
- `page_size`: Number of items per page (limit)

## Endpoints with Pagination

### 1. Get Products by Company ID

**Endpoint:** `GET /api/v1/products/company/:company_id`

**Query Parameters:**
- `limit` (optional, default: 50, max: 100): Number of items per page
- `offset` (optional, default: 0): Number of items to skip
- `category` (optional): Filter by category
- `is_available` (optional): Filter by availability (true/false)
- `is_addon` (optional): Filter addons only (true/false)

**Example Request:**
```bash
curl -X GET "http://localhost:8080/api/v1/products/company/550e8400-e29b-41d4-a716-446655440000?limit=10&offset=0"
```

**Example Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": "e88faead-d88c-441e-a247-e400c9758987",
      "name": "Helado Premium",
      "photos": ["https://example.com/photos/helado-main.jpg"],
      "category": "Helados",
      "min_price": 30000,
      "is_available": true
    }
  ],
  "meta": {
    "current_page": 1,
    "total_pages": 1,
    "total_items": 2,
    "page_size": 10
  }
}
```

### 2. Get Products by Sale Point ID

**Endpoint:** `GET /api/v1/products/sale-point/:sale_point_id`

**Query Parameters:** Same as company endpoint

**Example Request:**
```bash
curl -X GET "http://localhost:8080/api/v1/products/sale-point/660e8400-e29b-41d4-a716-446655440001?limit=20&offset=0"
```

**Example Response:**
```json
{
  "success": true,
  "data": [...],
  "meta": {
    "current_page": 1,
    "total_pages": 5,
    "total_items": 95,
    "page_size": 20
  }
}
```

## Pagination Calculation

### How `current_page` is calculated:
```
current_page = (offset / limit) + 1
```

**Examples:**
- `offset=0, limit=10` → `current_page=1` (first page)
- `offset=10, limit=10` → `current_page=2` (second page)
- `offset=20, limit=10` → `current_page=3` (third page)

### How `total_pages` is calculated:
```
total_pages = ceil(total_items / page_size)
```

**Examples:**
- `total_items=25, page_size=10` → `total_pages=3`
- `total_items=100, page_size=50` → `total_pages=2`
- `total_items=5, page_size=10` → `total_pages=1`

## Usage Examples

### Example 1: Paginate through all products (10 per page)

**Page 1:**
```bash
curl -X GET "http://localhost:8080/api/v1/products/company/COMPANY_ID?limit=10&offset=0"
```

**Page 2:**
```bash
curl -X GET "http://localhost:8080/api/v1/products/company/COMPANY_ID?limit=10&offset=10"
```

**Page 3:**
```bash
curl -X GET "http://localhost:8080/api/v1/products/company/COMPANY_ID?limit=10&offset=20"
```

### Example 2: Paginate with filters

**Get available products in "Helados" category, page 1:**
```bash
curl -X GET "http://localhost:8080/api/v1/products/company/COMPANY_ID?category=Helados&is_available=true&limit=10&offset=0"
```

Response will include `total_items` that match the filter criteria.

### Example 3: Get all items on single page (default)

If you don't specify `limit` and `offset`, defaults are used:
```bash
curl -X GET "http://localhost:8080/api/v1/products/company/COMPANY_ID"
```

This returns up to 50 items (page 1) with full pagination metadata.

## Best Practices

1. **Always use pagination**: Don't fetch all items at once, especially for large datasets
2. **Respect the limits**: Maximum limit is 100 items per page
3. **Use filters**: Combine pagination with filters to reduce the dataset
4. **Check total_pages**: Use this to know if there are more pages to fetch
5. **Calculate next offset**: `next_offset = current_offset + limit`

## Frontend Integration Example

### JavaScript/TypeScript
```javascript
async function fetchProducts(companyId, page = 1, pageSize = 10, filters = {}) {
  const offset = (page - 1) * pageSize;
  const params = new URLSearchParams({
    limit: pageSize,
    offset: offset,
    ...filters
  });
  
  const response = await fetch(
    `http://localhost:8080/api/v1/products/company/${companyId}?${params}`
  );
  
  const data = await response.json();
  
  return {
    products: data.data,
    currentPage: data.meta.current_page,
    totalPages: data.meta.total_pages,
    totalItems: data.meta.total_items,
    hasNextPage: data.meta.current_page < data.meta.total_pages,
    hasPrevPage: data.meta.current_page > 1
  };
}

// Usage
const result = await fetchProducts('550e8400-e29b-41d4-a716-446655440000', 1, 10, {
  category: 'Helados',
  is_available: true
});

console.log(`Page ${result.currentPage} of ${result.totalPages}`);
console.log(`Total products: ${result.totalItems}`);
console.log('Products:', result.products);
```

## Notes

- Pagination metadata is **only** included in product list endpoints
- Single product endpoint (`GET /api/v1/products/:id`) returns the product directly without pagination
- Category endpoints (`GET /api/v1/categories/*`) do not include pagination (they return all unique categories)
- The `total_items` count respects all applied filters
- Empty results return `total_items=0`, `total_pages=1`, `current_page=1`
