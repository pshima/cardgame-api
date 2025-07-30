# Security Implementation

This document describes the comprehensive security measures implemented in the Card Game API, including input validation, container security, runtime protection, and observability features. The application follows defense-in-depth principles with multiple layers of security controls.

## Architecture Security Benefits

The refactored clean architecture provides multiple security advantages:

### Layered Security
1. **Handler Layer**: First line of defense with input validation
2. **Service Layer**: Business logic validation and authorization
3. **Model Layer**: Data integrity and consistency checks
4. **Manager Layer**: Thread-safe state management

### Separation of Concerns
- **Validators Package**: Centralized validation logic prevents inconsistencies
- **Middleware**: Cross-cutting security concerns (auth, logging, rate limiting)
- **Services**: Business rule enforcement separate from HTTP concerns
- **Models**: Domain-driven security constraints

### Security Observability
- **Structured Logging**: Every security event tracked with context
- **Metrics Collection**: Real-time security monitoring
- **Trace Correlation**: Request flow tracking for forensics
- **Audit Trail**: Complete record of all operations

## Input Validation and Sanitization

### Overview
All user input from URI parameters and JSON request bodies is validated and sanitized before processing to prevent:
- SQL injection attempts
- Cross-site scripting (XSS)
- Path traversal attacks
- Buffer overflow attacks
- Invalid data that could cause application errors

### Validation Architecture

All validation logic is centralized in the `validators` package for consistency and maintainability:

```go
// validators/validators.go
package validators
```

### Validation Functions

#### `ValidateUUID(input string) bool`
- Located in `validators/validators.go`
- Validates game IDs and player IDs are in proper UUID format
- Pattern: `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`
- Prevents injection attacks through ID parameters
- Used by all handlers before processing requests

#### `ValidatePlayerID(input string) bool`
- Located in `validators/validators.go`
- Validates player IDs (UUID format or special "dealer" value)
- Allows the special case "dealer" for dealer operations
- All other values must be valid UUIDs

#### `ValidateNumber(input string) (int, bool)`
- Located in `validators/validators.go`
- Validates numeric inputs (deck count, player count, card count)
- Only accepts positive integers
- Prevents negative numbers and non-numeric input
- Guards against integer overflow

#### `ValidateDeckType(input string) bool`
- Located in `validators/validators.go`
- Validates deck type parameters
- Pattern: `^[a-zA-Z0-9_-]{1,20}$`
- Allows alphanumeric characters, hyphens, and underscores
- Maximum length: 20 characters

#### `ValidatePileID(input string) bool`
- Located in `validators/validators.go`
- Validates discard pile IDs
- Pattern: `^[a-zA-Z0-9_-]{1,50}$`
- Maximum length: 50 characters

#### `ValidateBoolean(input string) bool`
- Located in `validators/validators.go`
- Validates boolean parameters (e.g., face up/down)
- Accepts: "true", "false", "1", "0" (case-insensitive)

#### `SanitizeString(input string, maxLength int) string`
- Located in `validators/validators.go`
- Removes control characters (ASCII < 32 and 127)
- Enforces maximum length limits
- Prevents buffer overflow and injection attacks

#### `ValidateDeckName(name string) bool`
- Located in `validators/validators.go`
- Validates custom deck names
- Must be 1-128 characters in length
- Applied after sanitization

#### `ValidateCardIndex(indexStr string) (int, bool)`
- Located in `validators/validators.go`  
- Validates custom card indices
- Must be non-negative integers
- Used for custom card operations

### Protected Endpoints

All endpoints that accept URI parameters are protected:

1. **Game ID Parameters**
   - `/game/:gameId/*` - All game-related endpoints
   - Validates UUID format before processing

2. **Player ID Parameters**
   - `/game/:gameId/players/:playerId`
   - `/game/:gameId/deal/player/:playerId`
   - `/game/:gameId/hit/:playerId`
   - `/game/:gameId/stand/:playerId`
   - Validates UUID format or "dealer" value

3. **Numeric Parameters**
   - `/game/new/:decks` - Deck count (1-100)
   - `/game/new/:decks/:type/:players` - Player count (1-10)
   - `/game/:gameId/deal/:count` - Card count (1-52)
   - Validates positive integers within reasonable ranges

4. **Type Parameters**
   - `/game/new/:decks/:type` - Deck type
   - Validates against alphanumeric pattern

5. **Boolean Parameters**
   - `/game/:gameId/deal/player/:playerId/:faceUp`
   - Validates boolean representation

6. **Pile ID Parameters**
   - `/game/:gameId/discard/:pileId`
   - Validates alphanumeric pattern with length limit

7. **Custom Deck ID Parameters**
   - `/custom-decks/:deckId` - Custom deck UUID
   - `/custom-decks/:deckId/cards/*` - All custom deck card endpoints
   - Validates UUID format before processing

8. **Custom Card Index Parameters**
   - `/custom-decks/:deckId/cards/:cardIndex`
   - Validates non-negative integer indices

### JSON Request Body Validation

1. **Player Name** (POST `/game/:gameId/players`)
   - Sanitized to remove control characters
   - Maximum length: 50 characters
   - Cannot be empty after trimming whitespace

2. **Player ID in Request Body** (POST `/game/:gameId/discard/:pileId`)
   - Must be valid UUID or "dealer"
   - Validated before processing

3. **Card Index** (POST `/game/:gameId/discard/:pileId`)
   - Must be between 0-51 (reasonable card limit)
   - Prevents array index out of bounds

4. **Custom Deck Name** (POST `/custom-decks`)
   - Sanitized to remove control characters
   - Must be 1-128 characters after sanitization
   - Cannot be empty

5. **Custom Card Data** (POST `/custom-decks/:deckId/cards`)
   - **Name**: Required, sanitized, max 100 characters
   - **Rank**: Optional, can be numeric or string
   - **Suit**: Optional, sanitized, max 50 characters
   - **Attributes**: Optional, max 100 key-value pairs
     - Keys: Sanitized, max 50 characters each
     - Values: Sanitized, max 200 characters each
   - **Deck Limit**: Maximum 2,000 cards per deck

### Network Security

#### Trusted Proxy Configuration
The application implements comprehensive proxy security:

1. **Environment Configuration**
   ```bash
   # Development
   export TRUSTED_PROXIES="127.0.0.1,::1"
   
   # Production with specific proxies
   export TRUSTED_PROXIES="10.0.1.100,192.168.1.0/24"
   
   # Cloud environments
   # AWS ALB
   export TRUSTED_PROXIES="10.0.0.0/8"
   
   # GCP Load Balancer
   export TRUSTED_PROXIES="35.191.0.0/16,130.211.0.0/22"
   
   # Cloudflare (use actual Cloudflare IPs)
   export TRUSTED_PROXIES="173.245.48.0/20,103.21.244.0/22"
   ```

2. **Header Security**
   - **Accepted Headers**: Only from trusted proxies
     - `X-Forwarded-For`: Client IP extraction
     - `X-Forwarded-Proto`: Protocol detection
     - `X-Real-IP`: Alternative IP header
   - **Rejected Headers**: Security risks
     - `X-Forwarded-Host`: Prevent host injection
     - Custom headers from untrusted sources
   - **Validation**: Strict format checking

3. **Request Security Headers**
   The application sets security headers:
   ```go
   // Implemented in middleware/security.go
   X-Content-Type-Options: nosniff
   X-Frame-Options: DENY
   X-XSS-Protection: 1; mode=block
   Strict-Transport-Security: max-age=31536000
   Content-Security-Policy: default-src 'self'
   ```

4. **Rate Limiting**
   - Per-IP rate limiting (when behind proxy)
   - Configurable limits per endpoint
   - Automatic blocking of abusive IPs
   - Integration with fail2ban

#### DDoS Protection
1. **Application Level**
   - Request size limits
   - Connection timeouts
   - Resource pool limits
   
2. **Infrastructure Level**
   - Use CDN/WAF (Cloudflare, AWS Shield)
   - Configure rate limiting at proxy
   - Enable SYN cookies
   - Set up traffic analysis

### Security Benefits

1. **Prevention of Injection Attacks**
   - All input is validated against strict patterns
   - No user input is directly concatenated into queries or commands

2. **Protection Against XSS**
   - Control characters are stripped from all string inputs
   - HTML/JavaScript cannot be injected through parameters

3. **Path Traversal Prevention**
   - UUID validation prevents directory traversal attempts
   - Strict character whitelisting blocks path separators

4. **Denial of Service Protection**
   - Reasonable limits on numeric inputs (e.g., max 100 decks)
   - String length limits prevent memory exhaustion
   - Invalid input rejected early before resource allocation

5. **Data Integrity**
   - Type validation ensures data consistency
   - Range checks prevent logical errors
   - Empty or malformed data is rejected

6. **Custom Deck Security**
   - Strict limits on deck names, card counts, and attributes
   - Prevents resource exhaustion through large custom decks
   - Validates custom card data structure and content
   - Tombstone deletion prevents data loss while allowing cleanup

7. **Proxy Security**
   - Configured trusted proxy list to prevent IP spoofing
   - Only accepts proxy headers from trusted sources
   - Environment-based proxy configuration for production
   - Validates X-Forwarded-Proto headers from trusted proxies only

### Error Responses

When validation fails, the API returns appropriate error messages:
- `400 Bad Request` with descriptive error message
- Examples:
  - "Invalid game ID format"
  - "Invalid decks parameter (must be 1-100)"
  - "Invalid player ID format"
  - "Player name cannot be empty"
  - "Invalid deck ID format"
  - "Deck name must be 1-128 characters"
  - "Maximum card limit (2000) reached for this deck"
  - "Maximum 100 attributes allowed per card"
  - "Card name cannot be empty"

### Testing

Comprehensive validation tests are included in `validators/validators_test.go`:
- Unit tests for each validation function
- Edge cases and attack vectors
- Sanitization behavior verification
- Table-driven tests for comprehensive coverage

Run tests with: 
```bash
# Run validator tests
go test ./validators -v

# Run with coverage
go test ./validators -cover

# Run all security-related tests
make test-security
```

## Container Security

### Secure Docker Image

The application uses a multi-layered security approach for containerization:

1. **Multi-Stage Build Process**
   ```dockerfile
   # Dependencies stage - cached for efficiency
   FROM golang:1.24.4-alpine AS deps
   RUN adduser -D -u 1001 builder
   USER builder
   
   # Build stage - compile with security flags
   FROM deps AS builder
   RUN go mod verify  # Verify dependency integrity
   RUN CGO_ENABLED=0 go build -trimpath -ldflags='-w -s'
   
   # Runtime stage - minimal attack surface
   FROM alpine:3.19
   RUN apk --no-cache add ca-certificates
   ```

2. **Runtime Security**
   - Non-root user (UID 65532) with no shell
   - Read-only root filesystem capability
   - Minimal Alpine base (regularly updated)
   - Static files with restrictive permissions (0444)
   - No package manager in final image
   - Distroless compatible

3. **Security Metadata**
   ```dockerfile
   LABEL security.scan="true" \
         security.nonroot="true" \
         security.updates="auto" \
         security.base="alpine:3.19"
   ```

### Security Scanning

Comprehensive security scanning is integrated into the build process:

```bash
# Run all security scans
make security-scan

# Components scanned:
# 1. Vulnerability scan with Trivy
# 2. Secret detection with TruffleHog
# 3. Dockerfile linting with Hadolint
# 4. Dependency checking
# 5. SBOM generation with Syft
```

### Runtime Protection

#### Docker Compose Security
- **Read-only filesystem**: Prevents runtime modifications
- **Dropped capabilities**: Only NET_BIND_SERVICE retained
- **No new privileges**: Prevents privilege escalation
- **Security profiles**: AppArmor and Seccomp enabled
- **Resource limits**: CPU, memory, and PID limits
- **Network isolation**: Custom bridge network
- **Localhost binding**: Ports bound to 127.0.0.1 only

#### Kubernetes Security
- **Pod Security Standards**: Restricted profile
- **Network Policies**: Ingress/egress controls
- **RBAC**: Minimal service account permissions
- **Security Context**: Non-root, read-only filesystem
- **Resource Quotas**: Prevent resource exhaustion

### Secrets Management

1. **Docker Secrets**
   ```bash
   # Generate secrets
   ./scripts/generate-secrets.sh
   
   # Use with Docker Compose
   docker-compose -f docker-compose.secrets.yml up
   ```

2. **Environment Variables**
   - Secrets loaded from files, not environment
   - Reference via `/run/secrets/` mount
   - Read-only access (mode 0400)

3. **Kubernetes Secrets**
   ```yaml
   valueFrom:
     secretKeyRef:
       name: cardgame-secrets
       key: api-key
   ```

### Security Monitoring

1. **Application Security Monitoring**
   - **Structured Logging**: All security events logged with Zap
   - **Metrics Collection**: Security metrics via OpenTelemetry
     - Failed authentication attempts
     - Invalid input rejection rates
     - Suspicious request patterns
   - **Audit Trail**: Complete request/response logging
   - **Correlation IDs**: Request tracking across services

2. **Runtime Monitoring**
   - **Falco Integration**: Detect anomalous container behavior
   - **System Calls**: Monitor for suspicious activity
   - **File Access**: Alert on unexpected file operations
   - **Network Activity**: Track unusual connections

3. **Vulnerability Management**
   - **Automated Scanning**: Daily vulnerability scans
   - **Dependency Tracking**: Monitor for CVEs
   - **SBOM Generation**: Track all components
   - **Update Automation**: Dependabot integration
   
4. **Security Metrics Dashboard**
   Available metrics for security monitoring:
   - `security_validation_failures_total`: Input validation rejections
   - `security_auth_failures_total`: Authentication failures
   - `security_suspicious_requests_total`: Anomalous patterns
   - `security_blocked_ips_total`: Blocked IP addresses

### Security Checklist

#### Pre-Deployment
- [ ] Run comprehensive security scan: `make security-scan`
- [ ] Review and remediate HIGH/CRITICAL vulnerabilities
- [ ] Verify no secrets in code: `make secret-scan`
- [ ] Generate fresh secrets: `./scripts/generate-secrets.sh`
- [ ] Review dependency licenses: `make license-check`
- [ ] Update base images to latest versions
- [ ] Verify SBOM accuracy: `make sbom`

#### Container Security
- [ ] Enable read-only root filesystem
- [ ] Set appropriate security contexts
- [ ] Configure AppArmor/SELinux profiles
- [ ] Define resource limits (CPU, memory, PIDs)
- [ ] Implement network policies
- [ ] Use minimal base images
- [ ] Scan final image: `make scan-image`

#### Runtime Security
- [ ] Configure TLS/HTTPS with strong ciphers
- [ ] Enable structured audit logging
- [ ] Set up log aggregation and analysis
- [ ] Configure metrics and alerting
- [ ] Implement rate limiting and DDoS protection
- [ ] Set up WAF rules
- [ ] Configure intrusion detection

#### Access Control
- [ ] Review and minimize RBAC permissions
- [ ] Implement principle of least privilege
- [ ] Configure service mesh policies
- [ ] Set up network segmentation
- [ ] Enable mutual TLS (mTLS) where applicable

#### Monitoring & Response
- [ ] Configure security event monitoring
- [ ] Set up automated incident response
- [ ] Create security runbooks
- [ ] Test disaster recovery procedures
- [ ] Schedule regular security reviews

### Compliance Support

The security implementation supports multiple compliance frameworks:

#### Standards Compliance
- **CIS Docker Benchmark v1.6**: 
  - Non-root containers ✓
  - Minimal base images ✓
  - Security scanning ✓
  - Resource limits ✓
  
- **NIST Cybersecurity Framework**:
  - Identify: Asset inventory, SBOM generation
  - Protect: Input validation, access controls
  - Detect: Logging, monitoring, alerting
  - Respond: Incident response procedures
  - Recover: Backup and restore capabilities

- **OWASP Top 10 (2021)**:
  - A01 Broken Access Control: UUID validation, authorization
  - A02 Cryptographic Failures: TLS enforcement
  - A03 Injection: Input validation, parameterization
  - A04 Insecure Design: Threat modeling, secure defaults
  - A05 Security Misconfiguration: Hardened containers
  - A06 Vulnerable Components: Dependency scanning
  - A07 Auth Failures: Session management
  - A08 Software Integrity: Code signing, SBOM
  - A09 Logging Failures: Comprehensive audit logs
  - A10 SSRF: URL validation, network policies

- **PCI DSS v4.0**:
  - Secure configuration standards
  - Access control measures
  - Regular security testing
  - Monitoring and logging

#### Audit Support
The application provides:
- Comprehensive audit logs in JSON format
- Security event correlation
- Compliance reporting endpoints
- Automated compliance checks