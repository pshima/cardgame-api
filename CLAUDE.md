# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview
This is a comprehensive Card Game API built with Go and the Gin web framework. The API provides full blackjack gameplay functionality, custom deck creation, card images, player management, security features, and complete OpenAPI documentation.

## Tech Stack
- **Language**: Go 1.24.4
- **Framework**: Gin (github.com/gin-gonic/gin)
- **Module**: github.com/peteshima/cardgame-api
- **UUID Generation**: Google UUID (github.com/google/uuid)
- **Image Generation**: Go's image package with custom rendering
- **Font Rendering**: golang.org/x/image/font with Go fonts
- **Testing**: Testify (github.com/stretchr/testify/assert)

## Key Features
- **Complete Blackjack Implementation**: Full game flow with automatic dealer play
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

## Common Commands

### Development
- `go run .` - Start the development server on port 8080 (preferred)
- `go run main.go` - Alternative way to start server
- `go build` - Build the application binary
- `go mod tidy` - Clean up dependencies

### Testing
- `go test ./...` - Run all tests
- `go test -v ./...` - Run tests with verbose output
- `go test -run TestValidation` - Run specific test functions
- `./test_images.sh` - Test card image functionality
- `./test_docs.sh` - Test API documentation endpoints
- `./verify_cards.sh` - Verify all card images exist

### Card Generation
- `go run generate_cards.go` - Generate all card images (52 cards × 3 sizes + 3 backs = 159 images)

### Dependencies
- `go get <package>` - Add a new dependency
- `go mod download` - Download dependencies

## Project Structure

### Core Files
- `main.go` - Main application with all HTTP handlers and routing
- `card.go` - Card, Player, Game, and GameManager structs with blackjack logic
- `go.mod` - Go module definition and dependencies
- `go.sum` - Dependency checksums

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

## API Endpoints (35 total)

### System
- `GET /hello` - Health check endpoint
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
- `TRUSTED_PROXIES` - Comma-separated list of trusted proxy IPs (e.g., "10.0.1.100,192.168.1.0/24")
- `GIN_MODE` - Set to "release" for production deployment

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
1. **Add new endpoint**: Update main.go routing, add handler function, update OpenAPI spec
2. **Modify game logic**: Update card.go structs and methods, add tests
3. **Add validation**: Update validation functions in main.go, add tests in validation_test.go
4. **Update documentation**: Modify README.md and openapi.yaml
5. **Add new card size**: Update generate_cards.go constants and functions
6. **Custom deck features**: Modify CustomCard/CustomDeck structs in card.go, update CustomDeckManager methods

## Best Practices
- Always update the openapi documentation whenever the api changes

## Code Guidelines
- If you add new functions or classes, make sure to add code comments, 1-3 lines describing what it does and why it is needed

## Security Considerations
- Make sure to take security considerations in for all updates

## Testing Guidelines
- When making any changes please make sure all tests pass and fix them if they do not.  Code coverage should be 80% or above