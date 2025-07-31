# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview
This is a comprehensive Card Game API built with Go and the Gin web framework. The API provides full blackjack and cribbage gameplay functionality, custom deck creation, card images, player management, security features, complete OpenAPI documentation, and enterprise-grade observability. The codebase follows clean architecture principles with clear separation of concerns.

## Tech Stack
- **Language**: Go 1.24.4
- **Framework**: Gin (github.com/gin-gonic/gin)
- **Module**: github.com/peteshima/cardgame-api
- **UUID Generation**: Google UUID (github.com/google/uuid)
- **Logging**: Zap (go.uber.org/zap) - High-performance structured logging
- **Metrics**: OpenTelemetry (go.opentelemetry.io/otel) with Prometheus exporter
- **Image Generation**: Go's image package with custom rendering
- **Font Rendering**: golang.org/x/image/font with Go fonts
- **Testing**: Testify (github.com/stretchr/testify/assert)
- **Build Tools**: Make for automation, Docker for containerization

## Key Features
- **Complete Blackjack Implementation**: Full game flow with automatic dealer play
- **Glitchjack Support**: Blackjack variant with randomly generated deck composition
- **Cribbage Support**: Full 2-player cribbage implementation with pegging and scoring
- **Multiple Deck Types**: Standard (52 cards) and Spanish 21 (48 cards, no 10s)
- **Custom Deck Creation**: Free-form custom decks with custom cards, suits, ranks, and attributes
- **Card Images**: Auto-generated PNG images in 3 sizes (icon: 32x48, small: 64x90, large: 200x280) for standard cards
- **Game Compatibility**: Custom cards with numeric ranks can be used in traditional games
- **Tombstone Deletion**: Custom cards are marked as deleted but remain queryable
- **Security**: Comprehensive input validation and sanitization
- **API Documentation**: Complete OpenAPI 3.0 specification with interactive Swagger UI
- **Concurrent Games**: Thread-safe operations supporting multiple simultaneous games
- **Player Management**: Add/remove players with UUID-based identification
- **Observability**: Structured logging, OpenTelemetry metrics, health checks
- **Production Ready**: Docker support, cross-platform builds, monitoring endpoints

## Common Commands

### Development
- `make run` - Start the development server on port 8080 (preferred)
- `make dev` - Start with hot reload for development
- `make build` - Build the application binary for current platform
- `make build-all` - Build for all platforms (Linux, Windows, macOS)
- `make docker-build` - Build Docker image
- `make clean` - Clean build artifacts
- `go mod tidy` - Clean up dependencies

### Testing
- `make test` - Run all tests with coverage
- `make test-coverage` - Generate detailed coverage report
- `make coverage-html` - Generate HTML coverage report
- `make test-integration` - Run integration tests
- `go test ./handlers -v` - Test specific package
- `go test -run TestValidation` - Run specific test functions
- `./test_images.sh` - Test card image functionality
- `./test_docs.sh` - Test API documentation endpoints
- `./verify_cards.sh` - Verify all card images exist

### Card Generation
- `go run generate_cards.go` - Generate all card images (52 cards × 3 sizes + 3 backs = 159 images)

### Code Quality
- `make fmt` - Format all Go code
- `make lint` - Run golangci-lint checks
- `make vet` - Run go vet checks
- `make security-scan` - Run security vulnerability scans

### Dependencies
- `go get <package>` - Add a new dependency
- `go mod download` - Download dependencies

## Project Structure

### Application Architecture
```
├── main.go              - Application entry point and server setup
├── api/                 - Request/response DTOs
│   ├── requests.go      - Input validation structures
│   └── responses.go     - Output formatting structures
├── config/              - Configuration and setup
│   ├── config.go        - Environment configuration
│   ├── logger.go        - Zap logger setup
│   └── metrics.go       - OpenTelemetry metrics setup
├── handlers/            - HTTP request handlers
│   ├── blackjack_handlers.go    - Blackjack-specific endpoints
│   ├── cribbage_handlers.go     - Cribbage-specific endpoints
│   ├── custom_deck_handlers.go  - Custom deck endpoints
│   ├── deal_handlers.go         - Card dealing endpoints
│   ├── game_handlers.go         - Game management endpoints
│   ├── handler_dependencies.go  - Dependency injection
│   ├── reset_handlers.go        - Deck reset endpoints
│   └── system_handlers.go       - Health, metrics, docs
├── managers/            - State management
│   ├── game_manager.go          - Thread-safe game state
│   └── custom_deck_manager.go   - Custom deck storage
├── middleware/          - HTTP middleware
│   ├── logging.go       - Request/response logging
│   ├── metrics.go       - Metrics collection
│   └── recovery.go      - Panic recovery
├── models/              - Domain models
│   ├── card.go          - Card and deck entities
│   ├── player.go        - Player entity
│   ├── game.go          - Game entity
│   ├── deck.go          - Deck operations
│   ├── blackjack.go     - Blackjack-specific models
│   ├── cribbage.go      - Cribbage-specific models
│   └── custom_deck.go   - Custom deck models
├── services/            - Business logic
│   ├── game_service.go          - Core game operations
│   ├── blackjack_service.go     - Blackjack game logic
│   ├── cribbage_service.go      - Cribbage game logic
│   └── custom_deck_service.go   - Custom deck operations
└── validators/          - Input validation
    └── validators.go    - Validation functions
```

### Core Files
- `main.go` - Application entry point, server setup, and routing
- `go.mod` - Go module definition and dependencies
- `go.sum` - Dependency checksums
- `Makefile` - Build automation and common tasks

### Documentation
- `README.md` - Complete user documentation with API examples
- `openapi.yaml` - OpenAPI 3.0 specification (27 endpoints documented)
- `api-docs.html` - Interactive Swagger UI documentation page
- `SECURITY.md` - Security implementation documentation

### Testing
- `validation_test.go` - Security validation function tests
- `custom_deck_test.go` - Deck type and card generation tests
- `test_images.sh` - Card image functionality test script
- `test_docs.sh` - API documentation test script
- `verify_cards.sh` - Card image verification script
- `test_cards.html` - Visual card image test page

### Image Generation
- `generate_cards.go` - Card image generation script
- `static/cards/` - Generated card images directory structure
  - `icon/` - 32x48 pixel card images
  - `small/` - 64x90 pixel card images  
  - `large/` - 200x280 pixel card images

## API Endpoints (45 total)

### System & Monitoring
- `GET /hello` - Health check endpoint
- `GET /metrics` - Prometheus metrics endpoint
- `GET /stats` - Application statistics in JSON
- `GET /version` - Build version information
- `GET /api-docs` - Interactive API documentation
- `GET /openapi.yaml` - OpenAPI specification
- `GET /static/*` - Static file serving (card images)

### Game Management
- `GET /games` - List all active games
- `GET /game/new` - Create game (1 deck, standard, 6 max players)
- `GET /game/new/:decks` - Create game with specified decks
- `GET /game/new/:decks/:type` - Create game with deck type
- `GET /game/new/:decks/:type/:players` - Create game with all options
- `DELETE /game/:gameId` - Delete a game

### Game State
- `GET /game/:gameId` - Get basic game info
- `GET /game/:gameId/state` - Get complete game state with blackjack values
- `GET /game/:gameId/shuffle` - Shuffle deck
- `GET /game/:gameId/reset` - Reset deck to original state
- `GET /game/:gameId/reset/:decks` - Reset with different deck count
- `GET /game/:gameId/reset/:decks/:type` - Reset with different configuration

### Player Management
- `POST /game/:gameId/players` - Add player with JSON body: `{"name": "PlayerName"}`
- `DELETE /game/:gameId/players/:playerId` - Remove player

### Card Dealing
- `GET /game/:gameId/deal` - Deal single card (face up)
- `GET /game/:gameId/deal/:count` - Deal multiple cards
- `GET /game/:gameId/deal/player/:playerId` - Deal to player (face down)
- `GET /game/:gameId/deal/player/:playerId/:faceUp` - Deal with face control

### Discard Operations
- `POST /game/:gameId/discard/:pileId` - Discard card with JSON body

### Blackjack Gameplay
- `POST /game/:gameId/start` - Start blackjack game (deals initial cards)
- `POST /game/:gameId/hit/:playerId` - Player takes a card
- `POST /game/:gameId/stand/:playerId` - Player stands
- `GET /game/:gameId/results` - Get final game results

### Glitchjack Gameplay
- `GET /game/new/glitchjack` - Create new Glitchjack game (1 random deck, 6 max players)
- `GET /game/new/glitchjack/:decks` - Create Glitchjack game with multiple random decks
- `GET /game/new/glitchjack/:decks/:players` - Create Glitchjack game with specified decks and max players
- `POST /game/:gameId/glitchjack/start` - Start Glitchjack game (deals initial cards)
- `POST /game/:gameId/glitchjack/hit/:playerId` - Player takes a card
- `POST /game/:gameId/glitchjack/stand/:playerId` - Player stands
- `GET /game/:gameId/glitchjack/results` - Get final game results

### Cribbage Gameplay
- `GET /game/new/cribbage` - Create new cribbage game (2 players, 1 deck)
- `POST /game/:gameId/cribbage/start` - Start cribbage game (deals 6 cards each)
- `POST /game/:gameId/cribbage/discard/:playerId` - Discard 2 cards to crib
- `POST /game/:gameId/cribbage/play/:playerId` - Play a card during play phase
- `POST /game/:gameId/cribbage/go/:playerId` - Say "go" when can't play
- `GET /game/:gameId/cribbage/show` - Score hands and crib
- `GET /game/:gameId/cribbage/state` - Get complete cribbage game state

### Deck Types
- `GET /deck-types` - List available deck types with specifications

### Custom Deck Management
- `POST /custom-decks` - Create custom deck with JSON body: `{"name": "Deck Name"}`
- `GET /custom-decks` - List all custom decks with summary information
- `GET /custom-decks/:deckId` - Get custom deck details with all cards
- `DELETE /custom-decks/:deckId` - Delete custom deck permanently
- `POST /custom-decks/:deckId/cards` - Add card to deck with JSON body: `{"name": "Card Name", "rank": 9, "suit": "custom", "attributes": {...}}`
- `GET /custom-decks/:deckId/cards` - List cards in deck (query param `?include_deleted=true` for deleted cards)
- `GET /custom-decks/:deckId/cards/:cardIndex` - Get specific card by index
- `DELETE /custom-decks/:deckId/cards/:cardIndex` - Delete card (tombstone - remains queryable)

## Development Notes

### Server Configuration
- Server runs on port 8080 by default
- All API responses are in JSON format
- Uses Gin's default middleware for logging and recovery
- Serves static files from `./static` directory
- CORS headers not configured (add if needed for web apps)
- Trusted proxy configuration for production security

### Environment Variables
- `PORT` - Server port (default: 8080)
- `LOG_LEVEL` - Logging level: DEBUG, INFO, WARN, ERROR (default: INFO)
- `LOG_FORMAT` - Log format: json, console (default: json)
- `GIN_MODE` - Gin mode: debug, release (default: release)
- `TRUSTED_PROXIES` - Comma-separated list of trusted proxy IPs
- `METRICS_ENABLED` - Enable metrics collection (default: true)
- `OTEL_SERVICE_NAME` - OpenTelemetry service name (default: cardgame-api)
- `OTEL_RESOURCE_ATTRIBUTES` - Additional resource attributes

### Security Features
- All URI parameters validated with regex patterns
- Input sanitization removes control characters
- Length limits enforced on all string inputs
- UUID validation for game and player IDs
- Special handling for "dealer" player ID
- Numeric parameters have reasonable upper limits
- Custom deck validation: deck names (1-128 chars), cards (2000 max), attributes (100 max)

### Data Models
- **Ranks**: 1=Ace, 2-10=Number cards, 11=Jack, 12=Queen, 13=King
- **Suits**: 0=Hearts, 1=Diamonds, 2=Clubs, 3=Spades
- **Game Status**: waiting, in_progress, finished
- **Blackjack Results**: blackjack, win, push, bust, lose
- **Glitchjack Results**: Same as Blackjack (blackjack, win, push, bust, lose)
- **Custom Cards**: Free-form with name, optional rank/suit, attributes, game_compatible flag, tombstone deletion

### Card Images
- All cards include `images` object with URLs in responses
- Face-down cards show card back image
- Images auto-generated with proper symbol counts for each rank
- Traditional playing card layouts and colors

### Blackjack Logic
- Ace values: 11 or 1 (automatically optimized)
- Face cards: 10 points each
- Blackjack: 21 with exactly 2 cards
- Dealer: hits on 16, stands on 17
- Automatic dealer play when all players finish

### Glitchjack Logic
- Same hand values as Blackjack (Aces: 1 or 11, Face cards: 10)
- Same dealer rules (hits on 16, stands on 17)
- Same win conditions (blackjack, win, push, bust, lose)
- Random deck composition: 52 cards randomly selected from standard deck
- Duplicates allowed: Multiple copies of same card possible in single deck
- Card counting ineffective due to unpredictable deck composition

### Thread Safety
- GameManager uses read-write mutexes
- Safe for concurrent access across multiple games
- Each game has independent state

### Testing Strategy
- Unit tests for validation functions
- Integration tests for deck types
- Manual testing scripts for functionality
- Visual testing for card images

### Common Development Tasks

1. **Add new endpoint**:
   - Create handler function in appropriate `handlers/` file
   - Add routing in `main.go`
   - Add request/response DTOs in `api/`
   - Update OpenAPI specification
   - Add tests in `handlers/` test file

2. **Modify game logic**:
   - Update models in `models/` package
   - Modify service logic in `services/`
   - Update handler to use new service methods
   - Add unit tests for models and services
   - Add integration tests for complete flow

3. **Add validation**:
   - Add validation functions in `validators/validators.go`
   - Use validators in handler functions
   - Add tests in `validators/validators_test.go`

4. **Update documentation**:
   - Modify README.md for user-facing changes
   - Update CLAUDE.md for development guidance
   - Update openapi.yaml for API changes
   - Update inline code comments

5. **Add new game type**:
   - Create model in `models/` (e.g., `poker.go`)
   - Create service in `services/` (e.g., `poker_service.go`)
   - Create handlers in `handlers/` (e.g., `poker_handlers.go`)
   - Add routing in `main.go`
   - Update OpenAPI spec
   - Add comprehensive tests

6. **Add new metrics**:
   - Define metric in `config/metrics.go`
   - Instrument code in appropriate service/handler
   - Update metrics documentation

## Best Practices

### API Development
- Always update the OpenAPI documentation when changing endpoints
- Use appropriate HTTP status codes and error messages
- Validate all input parameters before processing
- Return consistent response formats

### Code Organization
- Follow clean architecture principles - handlers → services → models
- Keep business logic in services, not handlers
- Use dependency injection for testability
- Separate concerns into appropriate packages

### Code Guidelines
- Add code comments for all exported functions (1-3 lines)
- Use meaningful variable and function names
- Keep functions small and focused (< 50 lines preferred)
- Handle errors explicitly, don't ignore them
- Use early returns to reduce nesting

### Testing
- Write unit tests for all new functions
- Maintain code coverage above 80%
- Use table-driven tests for multiple scenarios
- Mock external dependencies in unit tests
- Write integration tests for complete workflows

### Security
- Validate and sanitize all user input
- Use parameterized queries (if adding database)
- Never log sensitive information
- Keep dependencies updated
- Run security scans before committing

### Performance
- Use read/write mutexes appropriately
- Avoid unnecessary allocations in hot paths
- Profile code for performance bottlenecks
- Use buffered channels where appropriate
- Consider caching for expensive operations

### Documentation
- Update README.md for user-facing features
- Update CLAUDE.md for development changes
- Keep OpenAPI spec in sync with code
- Document complex algorithms inline
- Update SECURITY.md for security changes