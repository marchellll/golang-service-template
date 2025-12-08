# Security Vulnerability Audit Report

**Date:** 2025-01-27
**Project:** golang-service-template
**Auditor:** Automated Security Scan

## Executive Summary

This security audit identified **8 vulnerabilities** across different severity levels:
- **2 CRITICAL** issues requiring immediate attention
- **2 HIGH** severity issues
- **4 MEDIUM/LOW** severity issues

---

## CRITICAL Vulnerabilities

### 1. Hardcoded JWT Secret Key
**Location:** `internal/middleware/authentication.go:20`
**Severity:** CRITICAL
**Risk:** Complete authentication bypass possible if secret is exposed

```go
keyFunc := func(ctx context.Context) (interface{}, error) {
    return []byte("secret"), nil // TODO: replace with actual secret from config
}
```

**Impact:**
- Attackers can forge JWT tokens if they discover the secret
- All authenticated endpoints become vulnerable
- No way to rotate secrets without code changes

**Recommendation:**
- Move JWT secret to environment variable (e.g., `JWT_SECRET`)
- Use a strong, randomly generated secret (minimum 32 bytes)
- Implement secret rotation capability
- Never commit secrets to version control

**Fix:**
```go
keyFunc := func(ctx context.Context) (interface{}, error) {
    secret := getenv("JWT_SECRET")
    if secret == "" {
        logger.Fatal().Msg("JWT_SECRET environment variable is required")
    }
    return []byte(secret), nil
}
```

---

### 2. Missing Authorization Checks
**Location:** `internal/service/task.go` (Get, Update, Delete methods)
**Severity:** CRITICAL
**Risk:** Users can access, modify, or delete any task regardless of ownership

**Current Behavior:**
- `Get()` - No ownership check before returning task
- `Update()` - No ownership check before updating task
- `Delete()` - No ownership check before deleting task
- Only `FindByUserId()` filters by user, but other endpoints don't

**Impact:**
- Any authenticated user can read/modify/delete tasks belonging to other users
- Data breach and unauthorized data modification
- Violation of data privacy principles

**Recommendation:**
- Add ownership checks in all task operations
- Verify `task.CreatedBy` matches the authenticated user's ID
- Return 403 Forbidden if user doesn't own the resource

**Fix Example:**
```go
func (s *taskService) Get(ctx context.Context, id string) (*model.Task, error) {
    // ... existing code ...

    // Add authorization check
    userId := ctx.Value(middleware.ContextKeyUserId).(string)
    if entity.CreatedBy != userId {
        return nil, errz.NewPrettyError(http.StatusForbidden, "forbidden", "you don't have permission to access this task", nil)
    }

    return entity, nil
}
```

---

## HIGH Severity Vulnerabilities

### 3. CORS Misconfiguration
**Location:** `internal/app/routes.go:27` and `internal/app/socketio.go:40`
**Severity:** HIGH
**Risk:** Cross-Origin attacks, credential theft

**Issues:**
1. **Echo CORS:** `echo_middleware.CORS()` with no configuration allows all origins
2. **Socket.IO CORS:** `Origin: "*"` with `Credentials: true` is invalid and insecure

**Impact:**
- Any website can make authenticated requests to your API
- CSRF attacks become easier
- Credentials can be stolen via malicious websites

**Recommendation:**
- Configure CORS to allow only trusted origins
- Remove `Credentials: true` if using wildcard origin (or vice versa)
- Use environment variable for allowed origins

**Fix:**
```go
e.Use(echo_middleware.CORSWithConfig(echo_middleware.CORSConfig{
    AllowOrigins: []string{getenv("ALLOWED_ORIGINS")}, // e.g., "https://yourdomain.com"
    AllowCredentials: true,
    AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.PATCH, echo.DELETE},
}))
```

---

### 4. PostgreSQL SSL Disabled
**Location:** `internal/app/db.go:39`
**Severity:** HIGH
**Risk:** Database credentials and data transmitted in plaintext

```go
dsn := fmt.Sprintf(
    "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
    ...
)
```

**Impact:**
- Database credentials exposed in network traffic
- All database queries/responses visible to network sniffers
- Violates compliance requirements (PCI-DSS, HIPAA, etc.)

**Recommendation:**
- Enable SSL/TLS for database connections
- Use `sslmode=require` (or `verify-full` for production)
- Configure proper SSL certificates
- make a sensible default in .env.template for local dev

**Fix:**
```go
sslmode := getenv("DB_SSLMODE")
if sslmode == "" {
    sslmode = "require" // Default to secure
}
dsn := fmt.Sprintf(
    "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
    config.DbConfig.Host, config.DbConfig.Port,
    config.DbConfig.Username, config.DbConfig.Password,
    config.DbConfig.DBName, sslmode,
)
```

---

## MEDIUM Severity Vulnerabilities

### 5. Hardcoded Development Passwords
**Location:** `compose.yml:31-33, 51, 108`
**Severity:** MEDIUM
**Risk:** Weak default credentials in development environment

**Issues:**
- MySQL root password: `the_root_password`
- MySQL user password: `the_service_password`
- PostgreSQL password: `the_service_password`
- Grafana admin password: `admin123`

**Impact:**
- If compose.yml is committed to version control, credentials are exposed
- Developers might use same passwords in production
- Easy to guess/default passwords

**Recommendation:**
- Use environment variables for all passwords
- Generate random passwords for development
- Document that these are development-only credentials

---

### 6. No Rate Limiting
**Location:** `internal/app/routes.go`
**Severity:** MEDIUM
**Risk:** Brute force attacks, DoS, resource exhaustion

**Impact:**
- Attackers can brute force authentication endpoints
- API endpoints vulnerable to DoS attacks
- No protection against automated abuse

**Recommendation:**
- Implement rate limiting middleware
- Use Redis-based rate limiting for distributed systems
- Set different limits for authenticated vs unauthenticated endpoints
- Consider using `github.com/labstack/echo/v4/middleware` rate limiter or `golang.org/x/time/rate`

---

### 7. Hardcoded JWT Issuer and Audience
**Location:** `internal/middleware/authentication.go:27-28`
**Severity:** MEDIUM
**Risk:** Token validation may be too permissive

```go
validator.HS256,
"http://example.com/", // TODO: replace with actual issuer URL from config
[]string{"audience"},  // TODO: replace with actual audience
```

**Impact:**
- Tokens from other services might be accepted if they share the secret
- Difficult to distinguish between different services/environments

**Recommendation:**
- Move issuer and audience to configuration
- Use environment-specific values
- Validate issuer and audience strictly

---

### 8. No Input Size Limits
**Location:** `internal/handler/task.go`
**Severity:** LOW
**Risk:** DoS via large payloads

**Impact:**
- Attackers can send extremely large request bodies
- Memory exhaustion possible
- Database storage abuse

**Recommendation:**
- Set request body size limits in Echo
- Validate input length in handlers
- Consider using `echo_middleware.BodyLimit()`

**Fix:**
```go
e.Use(echo_middleware.BodyLimit("1M")) // Limit to 1MB
```

---

## Additional Security Recommendations

### 9. Security Headers Missing
**Recommendation:** Add security headers middleware:
- `X-Content-Type-Options: nosniff`
- `X-Frame-Options: DENY`
- `X-XSS-Protection: 1; mode=block`
- `Strict-Transport-Security` (for HTTPS)
- `Content-Security-Policy`

### 10. Error Information Disclosure
**Location:** Various error handlers
**Recommendation:** Ensure error messages don't leak sensitive information (database structure, file paths, etc.)

### 11. Logging Sensitive Data
**Location:** `internal/middleware/logger.go:83`
**Good:** Already filters sensitive fields, but verify all sensitive data is filtered

### 12. Dependency Vulnerabilities
**Recommendation:**
- Run `govulncheck` regularly (requires Go 1.24+)
- Use `go list -m -u` to check for outdated dependencies
- Consider using Dependabot or similar tools
- Review security advisories for all dependencies

### 13. Health Check Endpoints Exposure
**Location:** `internal/app/routes.go:43-49`
**Recommendation:** Consider restricting `/healthz` endpoints to internal networks only

---

## Summary of Required Actions

### Immediate (Critical):
1. ✅ Replace hardcoded JWT secret with environment variable
2. ✅ Add authorization checks to all task operations

### High Priority:
3. ✅ Configure CORS properly
4. ✅ Enable SSL for PostgreSQL connections

### Medium Priority:
5. ✅ Move hardcoded passwords to environment variables
6. ✅ Implement rate limiting
7. ✅ Configure JWT issuer/audience from config
8. ✅ Add request body size limits

### Low Priority:
9. ✅ Add security headers
10. ✅ Review and update dependencies regularly

---

## Testing Recommendations

- Add security tests for authorization checks
- Test CORS configuration with different origins
- Verify rate limiting works as expected
- Test with invalid/malformed JWT tokens
- Perform penetration testing on authentication flow

---

**Note:** This audit focused on code-level vulnerabilities. Additional security measures should include:
- Regular dependency updates
- Security monitoring and logging
- Regular penetration testing
- Security code reviews
- Incident response planning

