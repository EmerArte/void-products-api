# CORS Configuration Guide

## Overview

The API now includes comprehensive CORS (Cross-Origin Resource Sharing) configuration with support for environment variables, making it production-ready and secure.

## Configuration

CORS settings are configured via environment variables in `.env` file or directly in your deployment environment.

### Environment Variables

```bash
# Allow all origins (development only)
CORS_ALLOWED_ORIGINS=*

# Allow specific origins (production recommended)
CORS_ALLOWED_ORIGINS=http://localhost:3000,https://myapp.com,https://admin.myapp.com

# Customize allowed HTTP methods
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS

# Customize allowed headers
CORS_ALLOWED_HEADERS=Content-Type,Authorization,X-Requested-With,X-API-Key
```

### Default Values

If not specified, the following defaults are used:

- **CORS_ALLOWED_ORIGINS**: `*` (allow all origins)
- **CORS_ALLOWED_METHODS**: `GET,POST,PUT,DELETE,OPTIONS`
- **CORS_ALLOWED_HEADERS**: `Content-Type,Authorization,X-Requested-With`

## Configuration Examples

### Development Environment

For local development, allow all origins:

```bash
# .env
CORS_ALLOWED_ORIGINS=*
```

### Production Environment - Single Frontend

For production with a single frontend application:

```bash
# .env or environment variables
CORS_ALLOWED_ORIGINS=https://myapp.com
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOWED_HEADERS=Content-Type,Authorization
```

### Production Environment - Multiple Frontends

For production with multiple frontend applications (e.g., web app and admin panel):

```bash
CORS_ALLOWED_ORIGINS=https://myapp.com,https://admin.myapp.com,https://mobile.myapp.com
```

### Staging Environment

For staging with both staging and local development:

```bash
CORS_ALLOWED_ORIGINS=https://staging.myapp.com,http://localhost:3000,http://localhost:3001
```

## CORS Headers

The middleware sets the following CORS headers:

| Header | Description | Value |
|--------|-------------|-------|
| `Access-Control-Allow-Origin` | Allowed origins | Based on configuration |
| `Access-Control-Allow-Credentials` | Allow credentials | `true` |
| `Access-Control-Allow-Methods` | Allowed HTTP methods | Configurable |
| `Access-Control-Allow-Headers` | Allowed request headers | Configurable |
| `Access-Control-Max-Age` | Preflight cache duration | `86400` (24 hours) |

## How It Works

### Wildcard Mode (`*`)

When `CORS_ALLOWED_ORIGINS=*`:
- All origins are allowed
- `Access-Control-Allow-Origin: *` is sent in response
- Suitable for public APIs or development

### Specific Origins Mode

When specific origins are configured:
1. The middleware checks the `Origin` header from the request
2. If the origin is in the allowed list, it's echoed back in `Access-Control-Allow-Origin`
3. If the origin is NOT in the allowed list, the request is blocked with `403 Forbidden`
4. If no `Origin` header is present, the request is allowed (same-origin or non-browser clients)

### Preflight Requests

The middleware automatically handles OPTIONS preflight requests:
- Returns `204 No Content`
- Includes all CORS headers
- Does not execute the route handler

## Testing CORS

### Test 1: Basic Request (Wildcard)

```bash
curl -I http://localhost:8080/health
```

Expected headers:
```
Access-Control-Allow-Origin: *
Access-Control-Allow-Credentials: true
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
Access-Control-Allow-Headers: Content-Type, Authorization, X-Requested-With
```

### Test 2: Preflight Request

```bash
curl -X OPTIONS -I http://localhost:8080/api/v1/products
```

Expected:
- Status: `204 No Content`
- All CORS headers present

### Test 3: Request with Allowed Origin

```bash
curl -H "Origin: http://localhost:3000" -I http://localhost:8080/health
```

Expected (if `http://localhost:3000` is in allowed list):
```
Access-Control-Allow-Origin: http://localhost:3000
```

### Test 4: Request with Non-Allowed Origin

```bash
curl -H "Origin: http://evil.com" -I http://localhost:8080/health
```

Expected (if `http://evil.com` is NOT in allowed list):
- Status: `403 Forbidden`
- No `Access-Control-Allow-Origin` header

## Docker Deployment

### docker-compose.yml

```yaml
version: '3.8'
services:
  api:
    image: products-api:latest
    environment:
      - SERVER_PORT=8080
      - CORS_ALLOWED_ORIGINS=https://myapp.com,https://admin.myapp.com
      - CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
      - CORS_ALLOWED_HEADERS=Content-Type,Authorization,X-API-Key
      - DATABASE_URI=mongodb://mongo:27017
    ports:
      - "8080:8080"
```

## Kubernetes Deployment

### ConfigMap

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: api-config
data:
  CORS_ALLOWED_ORIGINS: "https://myapp.com,https://admin.myapp.com"
  CORS_ALLOWED_METHODS: "GET,POST,PUT,DELETE,OPTIONS"
  CORS_ALLOWED_HEADERS: "Content-Type,Authorization,X-API-Key"
```

### Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: products-api
spec:
  replicas: 3
  template:
    spec:
      containers:
      - name: api
        image: products-api:latest
        envFrom:
        - configMapRef:
            name: api-config
```

## Security Best Practices

### ✅ DO

1. **Production**: Always specify exact allowed origins
   ```bash
   CORS_ALLOWED_ORIGINS=https://myapp.com
   ```

2. **Use HTTPS**: In production, only allow HTTPS origins
   ```bash
   CORS_ALLOWED_ORIGINS=https://myapp.com,https://admin.myapp.com
   ```

3. **Minimal Headers**: Only allow headers your frontend actually uses
   ```bash
   CORS_ALLOWED_HEADERS=Content-Type,Authorization
   ```

4. **Minimal Methods**: Only allow HTTP methods your API supports
   ```bash
   CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE
   ```

### ❌ DON'T

1. **Production Wildcard**: Never use `*` in production
   ```bash
   # ❌ Bad for production
   CORS_ALLOWED_ORIGINS=*
   ```

2. **Mixed Protocols**: Don't mix HTTP and HTTPS in production
   ```bash
   # ❌ Security risk
   CORS_ALLOWED_ORIGINS=http://myapp.com,https://myapp.com
   ```

3. **Overly Permissive**: Don't allow unnecessary headers or methods
   ```bash
   # ❌ Too permissive
   CORS_ALLOWED_HEADERS=*
   ```

## Troubleshooting

### Issue: CORS error in browser console

**Error**: `Access to fetch at 'http://localhost:8080/api/v1/products' from origin 'http://localhost:3000' has been blocked by CORS policy`

**Solution**: Add your frontend origin to `CORS_ALLOWED_ORIGINS`:
```bash
CORS_ALLOWED_ORIGINS=http://localhost:3000
```

### Issue: 403 Forbidden on API requests

**Cause**: Your origin is not in the allowed list

**Solution**: Check your frontend URL and add it to the allowed origins:
```bash
# Check the Origin header your browser is sending
# Add it to allowed origins
CORS_ALLOWED_ORIGINS=https://your-actual-frontend-url.com
```

### Issue: Preflight request fails

**Cause**: Required headers or methods not in allowed list

**Solution**: Add the headers/methods your frontend needs:
```bash
CORS_ALLOWED_HEADERS=Content-Type,Authorization,X-Custom-Header
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
```

### Issue: Credentials not working

**Cause**: Browser doesn't send cookies/credentials

**Solution**: 
1. Ensure `Access-Control-Allow-Credentials: true` is set (automatic in this implementation)
2. In your frontend, set `credentials: 'include'`:
   ```javascript
   fetch('http://localhost:8080/api/v1/products', {
     credentials: 'include'
   })
   ```

## Environment-Specific Examples

### Local Development (.env)
```bash
SERVER_PORT=8080
SERVER_MODE=debug
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:3001
```

### Staging (.env)
```bash
SERVER_PORT=8080
SERVER_MODE=release
CORS_ALLOWED_ORIGINS=https://staging.myapp.com
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOWED_HEADERS=Content-Type,Authorization
```

### Production (Environment Variables)
```bash
SERVER_PORT=8080
SERVER_MODE=release
CORS_ALLOWED_ORIGINS=https://myapp.com,https://www.myapp.com
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE
CORS_ALLOWED_HEADERS=Content-Type,Authorization
```

## Integration with Frontend

### JavaScript/TypeScript (Fetch API)

```javascript
// Development - wildcard allowed
const response = await fetch('http://localhost:8080/api/v1/products', {
  method: 'GET',
  headers: {
    'Content-Type': 'application/json',
  }
});

// Production - specific origin
const response = await fetch('https://api.myapp.com/api/v1/products', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': 'Bearer <token>'
  },
  credentials: 'include', // If using cookies
  body: JSON.stringify(data)
});
```

### React/Next.js

```typescript
// lib/api.ts
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export async function fetchProducts() {
  const response = await fetch(`${API_BASE_URL}/api/v1/products/company/${companyId}`, {
    headers: {
      'Content-Type': 'application/json',
    }
  });
  
  if (!response.ok) {
    throw new Error('Failed to fetch products');
  }
  
  return response.json();
}
```

## Summary

✅ CORS is fully configurable via environment variables  
✅ Supports wildcard (`*`) for development  
✅ Supports specific origins list for production  
✅ Automatic preflight handling  
✅ Secure by default with customizable headers and methods  
✅ Production-ready with proper origin validation  
