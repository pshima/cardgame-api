# Security Implementation

This document describes the security measures implemented in the Card Game API to protect against malicious input.

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

### Error Responses

When validation fails, the API returns appropriate error messages:
- `400 Bad Request` with descriptive error message
- Examples:
  - "Invalid game ID format"
  - "Invalid decks parameter (must be 1-100)"
  - "Invalid player ID format"
  - "Player name cannot be empty"

### Testing

Comprehensive validation tests are included in `validation_test.go`:
- Tests for each validation function
- Edge cases and attack vectors
- Sanitization behavior verification

Run tests with: `go test -v validation_test.go main.go card.go`