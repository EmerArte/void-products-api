# Validation Error Quick Reference

## What Changed?

### Before ❌
```json
{
  "success": false,
  "error": "Key: 'PartialUpdateOrderRequest.Status' Error:Field validation for 'Status' failed on the 'oneof' tag",
  "message": "Invalid request body"
}
```

### After ✅
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

## Affected Endpoints

| Endpoint | Method | Status |
|----------|--------|--------|
| `/api/v1/orders` | POST | ✅ Updated |
| `/api/v1/orders` | PATCH | ✅ Updated |
| `/api/v1/orders` | PUT | ✅ Updated |
| `/api/v1/products` | POST | ✅ Updated |
| `/api/v1/products/:id` | PUT | ✅ Updated |

## Common Validation Messages

### Order Fields

| Field | Invalid Value | Message |
|-------|---------------|---------|
| `status` | `"VERIFIEDed"` | `'status' must be one of: CREATED, VERIFIED, IN_PROGRESS, OUT_FOR_DELIVERY, DELIVERED, CANCELLED` |
| `sale_type` | `"INVALID"` | `'sale_type' must be one of: DELIVERY, ON_SITE` |
| `products` | `[]` (empty) | `'products' must contain at least 1 items` |
| `products` | missing | `'products' is required` |
| `customer.id_type` | `"INVALID"` | `'id_type' must be one of: CC, CE, TI, PASSPORT` |
| `customer.name` | missing (DELIVERY) | `'name' is required` |
| `customer.phone` | missing (DELIVERY) | `'phone' is required` |
| `shipping_address` | missing (DELIVERY) | `'shipping_address' is required` |
| `table_number` | missing (ON_SITE) | `'table_number' is required` |

### Product Fields

| Field | Validation | Message |
|-------|------------|---------|
| `name` | required | `'name' is required` |
| `category` | required | `'category' is required` |
| `company_id` | required | `'company_id' is required` |
| `sale_point_id` | required | `'sale_point_id' is required` |
| `price_variations` | required | `'price_variations' is required` |
| `price_variations` | min 1 | `'price_variations' must contain at least 1 items` |

## Response Structure

```json
{
  "success": false,              // Always false for errors
  "error": "string",             // Summary message
  "message": "Validation failed", // User-friendly message
  "details": [                   // Array of field-level errors
    {
      "field": "field_name",     // Snake_case field name
      "message": "Clear error"   // Specific, actionable message
    }
  ]
}
```

## Frontend Parsing

### JavaScript/TypeScript
```typescript
fetch('/api/v1/orders', {
  method: 'POST',
  body: JSON.stringify(orderData)
})
.then(res => res.json())
.then(data => {
  if (!data.success && data.details) {
    data.details.forEach(error => {
      console.log(`${error.field}: ${error.message}`);
    });
  }
});
```

### React Hook Form
```typescript
if (error.details) {
  error.details.forEach(detail => {
    setError(detail.field as any, {
      type: 'server',
      message: detail.message
    });
  });
}
```

## Testing

### Run Test Suite
```bash
./test_validation_errors.sh
```

### Manual Test
```bash
curl -X PATCH http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{"code": "TEST", "status": "INVALID"}' | jq
```

## Documentation

- **Complete Guide**: [docs/VALIDATION_ERROR_HANDLING.md](docs/VALIDATION_ERROR_HANDLING.md)
- **API Testing**: [docs/API_TESTING_GUIDE.md](docs/API_TESTING_GUIDE.md)
- **Implementation Summary**: [VALIDATION_IMPROVEMENTS_SUMMARY.md](VALIDATION_IMPROVEMENTS_SUMMARY.md)

## Key Benefits

✅ **Clear Messages** - Users know exactly what's wrong  
✅ **Actionable** - Shows valid values for enums  
✅ **Structured** - Easy to parse programmatically  
✅ **Consistent** - Same format across all endpoints  
✅ **Field-Level** - Pinpoints exact problem fields  

## HTTP Status Codes

| Code | Type | When |
|------|------|------|
| 400 | Validation Error | Format/type errors, missing required fields |
| 422 | Business Logic | Domain validation (e.g., customer required for DELIVERY) |
| 409 | Conflict | State conflicts (e.g., invalid status transition) |
| 404 | Not Found | Resource doesn't exist |

**Note:** Validation errors (400) now have `details` array. Business logic errors (422) use standard error format.
