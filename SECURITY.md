# Security Implementation

This document describes the comprehensive security measures implemented in the Card Game API, including input validation, container security, and runtime protection.

## Input Validation and Sanitization

### Overview
All user input from URI parameters and JSON request bodies is validated and sanitized before processing to prevent:
- SQL injection attempts
- Cross-site scripting (XSS)
- Path traversal attacks
- Buffer overflow attacks
- Invalid data that could cause application errors

### Validation Functions

#### `validateUUID(input string) bool`
- Validates game IDs and player IDs are in proper UUID format
- Pattern: `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`
- Prevents injection attacks through ID parameters

#### `validatePlayerID(input string) bool`
- Validates player IDs (UUID format or special "dealer" value)
- Allows the special case "dealer" for dealer operations
- All other values must be valid UUIDs

#### `validateNumber(input string) (int, bool)`
- Validates numeric inputs (deck count, player count, card count)
- Only accepts positive integers
- Prevents negative numbers and non-numeric input
- Guards against integer overflow

#### `validateDeckType(input string) bool`
- Validates deck type parameters
- Pattern: `^[a-zA-Z0-9_-]{1,20}$`
- Allows alphanumeric characters, hyphens, and underscores
- Maximum length: 20 characters

#### `validatePileID(input string) bool`
- Validates discard pile IDs
- Pattern: `^[a-zA-Z0-9_-]{1,50}$`
- Maximum length: 50 characters

#### `validateBoolean(input string) bool`
- Validates boolean parameters (e.g., face up/down)
- Accepts: "true", "false", "1", "0" (case-insensitive)

#### `sanitizeString(input string, maxLength int) string`
- Removes control characters (ASCII < 32 and 127)
- Enforces maximum length limits
- Prevents buffer overflow and injection attacks

#### `validateDeckName(name string) bool`
- Validates custom deck names
- Must be 1-128 characters in length
- Applied after sanitization

#### `validateCardIndex(indexStr string) (int, bool)`
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

### Proxy Security Configuration

#### Trusted Proxy Setup
The application implements secure proxy handling to prevent IP spoofing and header injection attacks:

1. **Environment Configuration**
   - Set `TRUSTED_PROXIES` environment variable with actual proxy IPs
   - Example: `TRUSTED_PROXIES="10.0.1.100,192.168.1.0/24,172.16.0.1"`
   - Defaults to localhost IPs for development

2. **Proxy Header Validation**
   - Only accepts `X-Forwarded-Proto` from trusted proxy IPs
   - Validates protocol values (only "http" or "https" accepted)
   - Ignores `X-Forwarded-Host` to prevent host header injection

3. **Production Deployment**
   ```bash
   # Example for Nginx reverse proxy
   export TRUSTED_PROXIES="10.0.1.100"
   export GIN_MODE="release"
   ```

4. **Common Proxy Configurations**
   - **Nginx/Apache**: Use the load balancer's internal IP
   - **Cloudflare**: Use Cloudflare's IP ranges
   - **AWS ALB**: Use the ALB's private IP range
   - **Google Cloud**: Use the load balancer IP range

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

Comprehensive validation tests are included in `validation_test.go`:
- Tests for each validation function
- Edge cases and attack vectors
- Sanitization behavior verification

Run tests with: `go test -v validation_test.go main.go card.go`

## Container Security

### Secure Docker Image

The application uses a hardened Docker image with multiple security layers:

1. **Multi-Stage Build**
   - Build stage runs as non-root user (UID 1001)
   - Dependencies verified with `go mod verify`
   - Minimal final image with only runtime requirements

2. **Runtime Security**
   - Runs as non-root user (UID 65532)
   - Read-only root filesystem
   - No shell access (`/sbin/nologin`)
   - Minimal Alpine base with security updates
   - Static files have read-only permissions (444)

3. **Security Labels**
   ```dockerfile
   LABEL security.scan="true" \
         security.nonroot="true" \
         security.updates="auto"
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

1. **Runtime Monitoring**
   - Falco integration for anomaly detection
   - Audit logging of all API access
   - Structured JSON logs for SIEM integration

2. **Vulnerability Management**
   - Weekly image rebuilds
   - Automated CVE scanning
   - Dependency update notifications

### Security Checklist

Before deploying:
- [ ] Run security scans: `make security-scan`
- [ ] Review scan results for HIGH/CRITICAL issues
- [ ] Generate fresh secrets
- [ ] Enable read-only filesystem
- [ ] Configure network policies
- [ ] Set resource limits
- [ ] Enable audit logging
- [ ] Configure TLS/HTTPS
- [ ] Implement rate limiting
- [ ] Review RBAC permissions

### Compliance Support

The security implementation supports:
- **CIS Docker Benchmark**: Hardened container configuration
- **NIST Guidelines**: Defense in depth approach
- **OWASP Top 10**: Protection against common vulnerabilities
- **PCI DSS**: Secure configuration and access controls