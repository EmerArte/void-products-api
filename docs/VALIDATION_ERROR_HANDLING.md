# Validation Error Handling - Implementation Guide

## Overview

This document describes the validation error handling improvements implemented across all API endpoints that accept request bodies.

## Problem

Previously, validation errors returned raw validator messages that were not user-friendly:

**Before:**
```json
{
  "success": false,
  "error": "Key: 'PartialUpdateOrderRequest.Status' Error:Field validation for 'Status' failed on the 'oneof' tag",
  "message": "Invalid request body"
}
```

This error message:
- ❌ Exposes internal struct field names
- ❌ Uses technical validation tags
- ❌ Doesn't explain what values are valid
- ❌ Hard for frontend developers to parse

## Solution

Implemented a comprehensive validation error handler that provides:
- ✅ User-friendly field names (snake_case)
- ✅ Clear, actionable error messages
- ✅ Field-level error details
- ✅ Valid values for enum fields
- ✅ Consistent error response format

**After:**
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

## Architecture

Following Clean Architecture principles, the implementation is organized as:

```
internal/
├── handler/
│   ├── validation.go       # Validation error formatter (NEW)
│   ├── order_handler.go    # Updated with validation handling
│   └── product_handler.go  # Updated with validation handling
└── response/
    └── api_response.go     # Updated with ValidationError response
```

### Files Created/Modified

#### 1. `internal/handler/validation.go` (NEW)

Centralized validation error handling with:

- **`FormatValidationErrors(error)`**: Converts validator errors to user-friendly messages
- **`formatFieldName(string)`**: Converts PascalCase to snake_case
- **`getValidationMessage(FieldError)`**: Maps validation tags to readable messages

Supported validation tags:
- `required` → "'field' is required"
- `oneof` → "'field' must be one of: value1, value2"
- `min`/`max` → Context-aware messages for strings, numbers, arrays
- `email`, `url`, `uuid` → Format validation messages
- `gte`, `lte`, `gt`, `lt` → Numeric comparison messages
- And 20+ more validation rules

#### 2. `internal/response/api_response.go` (UPDATED)

Added new response types:

```go
// ValidationErrorDetail represents a field-level validation error
type ValidationErrorDetail struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}

// ValidationErrorResponse with field details
type ValidationErrorResponse struct {
    Success bool                    `json:"success"`
    Error   string                  `json:"error"`
    Message string                  `json:"message,omitempty"`
    Details []ValidationErrorDetail `json:"details,omitempty"`
}

// ValidationError sends structured validation error response
func ValidationError(c *gin.Context, statusCode int, errorMsg string, message string, details []ValidationErrorDetail)
```

#### 3. Handlers Updated

Both `order_handler.go` and `product_handler.go` now use the validation error handler in all endpoints that accept request bodies:

**Orders:**
- ✅ `POST /api/v1/orders` (Create)
- ✅ `PATCH /api/v1/orders` (PartialUpdate)
- ✅ `PUT /api/v1/orders` (Modify)

**Products:**
- ✅ `POST /api/v1/products` (Create)
- ✅ `PUT /api/v1/products/:id` (Update)

## Response Format

### Single Field Error

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

### Multiple Field Errors

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

### Missing Required Field

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

## Example Test Cases

### Test 1: Invalid Enum Value

**Request:**
```bash
curl -X PATCH http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "code": "ORD-TEST-123",
    "status": "VERIFIEDed"
  }'
```

**Response:**
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

### Test 2: Multiple Validation Errors

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "sale_type": "INVALID_TYPE",
    "products": []
  }'
```

**Response:**
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

### Test 3: Missing Required Field

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "sale_type": "DELIVERY"
  }'
```

**Response:**
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

## Frontend Integration

### Parsing Validation Errors

Frontend developers can now easily parse and display validation errors:

**TypeScript Example:**
```typescript
interface ValidationErrorDetail {
  field: string;
  message: string;
}

interface ValidationErrorResponse {
  success: boolean;
  error: string;
  message: string;
  details?: ValidationErrorDetail[];
}

// Handle API error
async function createOrder(data: OrderData) {
  try {
    const response = await fetch('/api/v1/orders', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data)
    });
    
    if (!response.ok) {
      const error: ValidationErrorResponse = await response.json();
      
      if (error.details) {
        // Display field-level errors
        error.details.forEach(detail => {
          showFieldError(detail.field, detail.message);
        });
      } else {
        // Display general error
        showGeneralError(error.error);
      }
    }
  } catch (err) {
    console.error('Network error:', err);
  }
}
```

**React Hook Form Example:**
```typescript
const onSubmit = async (data: OrderFormData) => {
  try {
    await createOrder(data);
  } catch (error) {
    if (error.details) {
      error.details.forEach((detail: ValidationErrorDetail) => {
        setError(detail.field as any, {
          type: 'server',
          message: detail.message
        });
      });
    }
  }
};
```

## Validation Tag Reference

| Tag | Example Message |
|-----|-----------------|
| `required` | `'products' is required` |
| `oneof=val1 val2` | `'status' must be one of: val1, val2` |
| `min=5` (string) | `'name' must be at least 5 characters long` |
| `max=100` (string) | `'description' must be at most 100 characters long` |
| `min=1` (array) | `'products' must contain at least 1 items` |
| `gte=0` | `'price' must be greater than or equal to 0` |
| `gt=0` | `'quantity' must be greater than 0` |
| `lte=100` | `'limit' must be less than or equal to 100` |
| `email` | `'email' must be a valid email address` |
| `url` | `'payment_receipt_url' must be a valid URL` |
| `uuid` | `'product_id' must be a valid UUID` |

## Benefits

### For Users
- ✅ Clear, actionable error messages
- ✅ Know exactly what values are valid
- ✅ Better user experience

### For Frontend Developers
- ✅ Easy to parse programmatically
- ✅ Can map errors to form fields automatically
- ✅ Consistent error format across all endpoints
- ✅ Field-level validation feedback

### For Backend Developers
- ✅ Centralized validation logic
- ✅ Easy to add new validation messages
- ✅ Consistent across all handlers
- ✅ Follows Clean Architecture principles
- ✅ Maintainable and extensible

## Best Practices

### When to Use

Use validation error handling for:
- ✅ Request body validation errors
- ✅ Field format errors (email, URL, UUID)
- ✅ Enum validation (oneof)
- ✅ Range validation (min, max, gte, lte)
- ✅ Required field errors

### When NOT to Use

Don't use for:
- ❌ Business logic errors (use domain errors)
- ❌ Not found errors (404)
- ❌ Conflict errors (409)
- ❌ Server errors (500)

## Extending Validation Messages

To add support for new validation tags:

1. Open `internal/handler/validation.go`
2. Add a new case in the `getValidationMessage()` switch statement:

```go
case "your_tag":
    return fmt.Sprintf("'%s' your custom message", field)
```

## Testing

All validation error handling has been tested with:
- ✅ Single field errors
- ✅ Multiple field errors
- ✅ Required field errors
- ✅ Enum validation (oneof)
- ✅ Nested field validation
- ✅ Array validation

Run tests:
```bash
go test ./internal/handler/... -v
```

## Summary

✅ **Validation errors are now user-friendly and actionable**  
✅ **Consistent across all endpoints**  
✅ **Easy for frontend to parse and display**  
✅ **Follows Clean Architecture principles**  
✅ **Maintainable and extensible**

The validation error handling is production-ready and provides a professional API experience!
