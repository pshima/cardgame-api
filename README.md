# Card Game API

A comprehensive Go/Gin API for card games with full blackjack support, player hand tracking, and multi-deck functionality.

## Features

- **Multiple Card Games**: Blackjack, Poker, War, Go Fish
- **Multiple Deck Types**: Standard 52-card, Spanish 21 (48-card, no 10s)
- **Player Management**: Add/remove players, track individual hands
- **Blackjack Logic**: Hand value calculation, automatic dealer play, winner determination
- **Face Up/Down Cards**: Full control over card visibility
- **Multi-Pile Discard System**: Support for multiple discard piles
- **Concurrent Games**: Thread-safe operations for multiple simultaneous games
- **Session Management**: UUID-based game sessions with cleanup
- **Card Images**: Auto-generated PNG images for all cards in icon (32x48), small (64x90), and large (200x280) formats
- **Image URLs**: All card responses include URLs for card images in three sizes
- **Security**: Comprehensive input validation and sanitization for all parameters
- **API Documentation**: Complete OpenAPI 3.0 specification with interactive Swagger UI

## Quick Start

### 1. Start the Server
```bash
go run .
```
Server runs on `http://localhost:8080`

### 2. Test Basic Endpoint
```bash
curl http://localhost:8080/hello
```

### 3. View API Documentation
- **Interactive API Documentation**: http://localhost:8080/api-docs
- **OpenAPI Specification**: http://localhost:8080/openapi.yaml

## Complete Blackjack Game Flow

Here's a step-by-step example of running a full blackjack game with 2 players:

### Step 1: Create a Blackjack Game
```bash
# Create game with 2 decks, standard cards, max 4 players
curl "http://localhost:8080/game/new/2/standard/4"
```

**Response:**
```json
{
  "game_id": "123e4567-e89b-12d3-a456-426614174000",
  "game_type": "Blackjack",
  "deck_name": "Golden Phoenix",
  "deck_type": "Standard",
  "max_players": 4,
  "current_players": 0,
  "message": "New Standard Blackjack game created",
  "remaining_cards": 104,
  "created": "2025-07-29T10:30:00Z"
}
```

### Step 2: Add Players
```bash
# Add first player
curl -X POST "http://localhost:8080/game/123e4567-e89b-12d3-a456-426614174000/players" \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice"}'

# Add second player  
curl -X POST "http://localhost:8080/game/123e4567-e89b-12d3-a456-426614174000/players" \
  -H "Content-Type: application/json" \
  -d '{"name": "Bob"}'
```

**Response:**
```json
{
  "game_id": "123e4567-e89b-12d3-a456-426614174000",
  "player": {
    "id": "player-uuid-alice",
    "name": "Alice",
    "hand": []
  },
  "message": "Player added successfully"
}
```

### Step 3: Shuffle the Deck
```bash
curl "http://localhost:8080/game/123e4567-e89b-12d3-a456-426614174000/shuffle"
```

### Step 4: Start the Blackjack Game
```bash
# This deals 2 cards to each player and dealer automatically
curl -X POST "http://localhost:8080/game/123e4567-e89b-12d3-a456-426614174000/start"
```

**Response:**
```json
{
  "game_id": "123e4567-e89b-12d3-a456-426614174000",
  "status": "in_progress",
  "message": "Blackjack game started",
  "current_player": 0
}
```

### Step 5: Check Game State
```bash
curl "http://localhost:8080/game/123e4567-e89b-12d3-a456-426614174000/state"
```

**Response:**
```json
{
  "game_id": "123e4567-e89b-12d3-a456-426614174000",
  "game_type": "Blackjack",
  "status": "in_progress",
  "current_player": 0,
  "remaining_cards": 98,
  "players": [
    {
      "id": "player-uuid-alice",
      "name": "Alice",
      "hand": [
        {
          "rank": 10, 
          "suit": 2, 
          "face_up": true,
          "images": {
            "icon": "http://localhost:8080/static/cards/icon/10_2.png",
            "small": "http://localhost:8080/static/cards/small/10_2.png",
            "large": "http://localhost:8080/static/cards/large/10_2.png"
          }
        },
        {
          "rank": 1, 
          "suit": 0, 
          "face_up": true,
          "images": {
            "icon": "http://localhost:8080/static/cards/icon/1_0.png",
            "small": "http://localhost:8080/static/cards/small/1_0.png",
            "large": "http://localhost:8080/static/cards/large/1_0.png"
          }
        }
      ],
      "hand_size": 2,
      "hand_value": 21,
      "has_blackjack": true,
      "is_busted": false
    },
    {
      "id": "player-uuid-bob", 
      "name": "Bob",
      "hand": [
        {"rank": 7, "suit": 1, "face_up": true},
        {"rank": 5, "suit": 3, "face_up": true}
      ],
      "hand_size": 2,
      "hand_value": 12,
      "has_blackjack": false,
      "is_busted": false
    }
  ],
  "dealer": {
    "id": "dealer",
    "name": "Dealer", 
    "hand": [
      {"rank": 13, "suit": 2, "face_up": false},
      {"rank": 6, "suit": 1, "face_up": true}
    ],
    "hand_size": 2,
    "hand_value": 16,
    "has_blackjack": false,
    "is_busted": false
  }
}
```

### Step 6: Player Actions

#### Alice has blackjack, so Bob plays first
```bash
# Bob hits (takes another card)
curl -X POST "http://localhost:8080/game/123e4567-e89b-12d3-a456-426614174000/hit/player-uuid-bob"
```

**Response:**
```json
{
  "game_id": "123e4567-e89b-12d3-a456-426614174000",
  "player_id": "player-uuid-bob",
  "player_name": "Bob",
  "hand_value": 21,
  "hand_size": 3,
  "has_blackjack": false,
  "is_busted": false,
  "message": "Card dealt to Bob"
}
```

#### Bob got 21! Now he stands
```bash
curl -X POST "http://localhost:8080/game/123e4567-e89b-12d3-a456-426614174000/stand/player-uuid-bob"
```

**Response:**
```json
{
  "game_id": "123e4567-e89b-12d3-a456-426614174000",
  "player_id": "player-uuid-bob",
  "player_name": "Bob", 
  "status": "finished",
  "current_player": 2,
  "message": "Bob stands"
}
```

### Step 7: Get Final Results
```bash
curl "http://localhost:8080/game/123e4567-e89b-12d3-a456-426614174000/results"
```

**Response:**
```json
{
  "game_id": "123e4567-e89b-12d3-a456-426614174000",
  "status": "finished",
  "dealer": {
    "hand_value": 16,
    "has_blackjack": false,
    "is_busted": false
  },
  "players": [
    {
      "player_id": "player-uuid-alice",
      "player_name": "Alice",
      "hand_value": 21,
      "has_blackjack": true,
      "is_busted": false,
      "result": "blackjack"
    },
    {
      "player_id": "player-uuid-bob",
      "player_name": "Bob", 
      "hand_value": 21,
      "has_blackjack": false,
      "is_busted": false,
      "result": "win"
    }
  ],
  "results": {
    "player-uuid-alice": "blackjack",
    "player-uuid-bob": "win"
  }
}
```

## API Endpoints

### Game Management
- `GET /deck-types` - List available deck types
- `GET /game/new` - Create new game (1 standard deck, 6 max players)
- `GET /game/new/:decks` - Create game with multiple decks
- `GET /game/new/:decks/:type` - Create game with specific deck type
- `GET /game/new/:decks/:type/:players` - Create game with max player limit
- `GET /games` - List all active games
- `DELETE /game/:gameId` - Delete a game

### Game State
- `GET /game/:gameId` - Get basic game info
- `GET /game/:gameId/state` - Get complete game state with hand values
- `GET /game/:gameId/shuffle` - Shuffle the deck

### Player Management  
- `POST /game/:gameId/players` - Add player `{"name": "PlayerName"}`
- `DELETE /game/:gameId/players/:playerId` - Remove player

### Blackjack Game Flow
- `POST /game/:gameId/start` - Start blackjack game (deals initial cards)
- `POST /game/:gameId/hit/:playerId` - Player takes a card
- `POST /game/:gameId/stand/:playerId` - Player stands (ends turn)
- `GET /game/:gameId/results` - Get game results and winners

### Manual Card Dealing (Advanced)
- `GET /game/:gameId/deal` - Deal one card from deck
- `GET /game/:gameId/deal/:count` - Deal multiple cards from deck  
- `GET /game/:gameId/deal/player/:playerId` - Deal card to specific player (face down)
- `GET /game/:gameId/deal/player/:playerId/:faceUp` - Deal card with face up/down control

### Discard & Reset
- `POST /game/:gameId/discard/:pileId` - Discard card `{"player_id": "id", "card_index": 0}`
- `GET /game/:gameId/reset` - Reset deck with same configuration
- `GET /game/:gameId/reset/:decks` - Reset with different deck count
- `GET /game/:gameId/reset/:decks/:type` - Reset with different deck type

## Deck Types

### Standard (52 cards)
- **ID**: 0
- **Cards**: Ace through King in all 4 suits
- **Total**: 52 cards per deck

### Spanish 21 (48 cards)  
- **ID**: 1
- **Cards**: Ace through King in all 4 suits, **excluding all 10s**
- **Total**: 48 cards per deck
- **Use Case**: Spanish Blackjack variant

## Game Types

- **Blackjack**: Full blackjack implementation with automatic dealer play
- **Poker**: Framework ready (cards, players, hands)
- **War**: Framework ready
- **Go Fish**: Framework ready

## Blackjack Rules Implemented

### Hand Values
- **Number cards (2-9)**: Face value
- **Face cards (J, Q, K)**: 10 points each
- **Aces**: 11 points (automatically converted to 1 if hand would bust)

### Game Flow
1. Players join game (up to configured maximum)
2. Game starts - deals 2 cards to each player and dealer
3. Players' cards are face up, dealer's first card is face down
4. Players take turns hitting or standing
5. When all players finish, dealer plays automatically
6. Dealer hits on 16 or less, stands on 17 or more
7. Winners determined by comparing final hand values

### Win Conditions
- **Blackjack**: 21 with exactly 2 cards (beats regular 21)
- **Win**: Higher value than dealer without busting
- **Push**: Same value as dealer (tie)
- **Bust**: Hand value over 21 (automatic loss)
- **Lose**: Lower value than dealer

## Advanced Features

### Concurrent Games
The API supports multiple simultaneous games with thread-safe operations:

```bash
# Create multiple games
curl "http://localhost:8080/game/new/1/standard/2"  # Game 1
curl "http://localhost:8080/game/new/6/spanish21/4" # Game 2  

# List all games
curl "http://localhost:8080/games"
```

### Face Up/Down Cards
Full control over card visibility:

```bash
# Deal face down card
curl "http://localhost:8080/game/GAME_ID/deal/player/PLAYER_ID/false"

# Deal face up card  
curl "http://localhost:8080/game/GAME_ID/deal/player/PLAYER_ID/true"
```

### Multi-Deck Support
Perfect for casino-style blackjack:

```bash
# Create 8-deck blackjack game (typical casino setup)
curl "http://localhost:8080/game/new/8/standard/7"
```

### Spanish 21 Support
```bash
# Create Spanish 21 game (no 10s, 6 decks typical)
curl "http://localhost:8080/game/new/6/spanish21/6"
```

## Error Handling

The API returns appropriate HTTP status codes and error messages:

### Common Error Responses
```json
{
  "error": "Game not found"
}
```

```json  
{
  "error": "Player not found"
}
```

```json
{
  "error": "Game is full"
}
```

```json
{
  "error": "No cards remaining in deck"
}
```

## Development

### Running Tests

#### Unit Tests
```bash
# Run all tests
go test ./...

# Run specific test
go test -run TestBlackjackHandValue

# Run with verbose output
go test -v ./...

# Run validation tests specifically
go test -v validation_test.go main.go card.go
```

#### Integration Tests
```bash
# Test card image functionality
./test_images.sh

# Test API documentation endpoints
./test_docs.sh

# Verify all card images exist
./verify_cards.sh

# Visual card testing (open in browser)
# Start server first, then open test_cards.html
```

#### Test Coverage
- **Unit Tests**: Validation functions, deck types, card generation
- **Integration Tests**: Card images, API documentation, blackjack gameplay
- **Visual Tests**: Card image rendering and layout verification
- **Security Tests**: Input validation and sanitization functions

### Building
```bash
# Build binary
go build

# Run
./cardgame-api
```

## Card Images

All cards are automatically rendered as PNG images in three sizes:
- **Icon**: 32x48 pixels - Perfect for UI indicators and small displays
- **Small**: 64x90 pixels - Good for hand displays and medium UI elements  
- **Large**: 200x280 pixels - Full-size card display

### Image URLs in Responses

Every card object in API responses includes an `images` object with URLs:

```json
{
  "rank": 1,
  "suit": 0,
  "face_up": true,
  "images": {
    "icon": "http://localhost:8080/static/cards/icon/1_0.png",
    "small": "http://localhost:8080/static/cards/small/1_0.png",
    "large": "http://localhost:8080/static/cards/large/1_0.png"
  }
}
```

Face-down cards return the card back image:
```json
{
  "rank": 13,
  "suit": 2,
  "face_up": false,
  "images": {
    "icon": "http://localhost:8080/static/cards/icon/back.png",
    "small": "http://localhost:8080/static/cards/small/back.png",
    "large": "http://localhost:8080/static/cards/large/back.png"
  }
}
```

### Generating Card Images

Card images are pre-generated using the included generator:
```bash
go run generate_cards.go
```

This creates all 52 cards plus card backs in all three sizes (157 total images).

## Security Features

The API includes comprehensive security measures to protect against malicious input:

### Input Validation
- **UUID Validation**: All game and player IDs must be valid UUID format
- **Parameter Sanitization**: Control characters removed from all string inputs
- **Length Limits**: Maximum length enforced on all text inputs (player names: 50 chars)
- **Numeric Bounds**: Reasonable limits on numeric parameters (decks: 1-100, players: 1-10)
- **Pattern Matching**: Deck types and pile IDs validated against safe character sets

### Security Protections
- **SQL Injection Prevention**: All input validated before processing
- **XSS Protection**: Control characters stripped from responses
- **Path Traversal Prevention**: UUID validation blocks directory traversal attempts
- **DoS Protection**: Reasonable limits prevent resource exhaustion
- **Data Integrity**: Type validation ensures consistent data structures

### Validation Details
See [SECURITY.md](SECURITY.md) for complete security implementation documentation including:
- Validation function specifications
- Input sanitization processes
- Error handling approaches
- Security testing methodology

## API Documentation

The API is fully documented with OpenAPI 3.0 specification:

- **Interactive Documentation**: Visit http://localhost:8080/api-docs while the server is running for a complete interactive API reference with the ability to test endpoints directly
- **OpenAPI Spec**: Available at http://localhost:8080/openapi.yaml
- **Complete Coverage**: All 27 endpoints documented with request/response schemas, examples, and validation rules

### Key Documentation Features:
- **Request/Response Examples**: Every endpoint includes realistic examples
- **Schema Validation**: Complete data models for all request and response objects
- **Security Documentation**: Input validation and sanitization details
- **Error Responses**: Documented error cases with appropriate HTTP status codes
- **Try It Out**: Interactive testing directly from the documentation

## Tech Stack

- **Language**: Go 1.24.4
- **Framework**: Gin HTTP web framework
- **UUID**: Google UUID for unique identifiers
- **Testing**: Testify for assertions
- **Concurrency**: Go sync package for thread safety
- **Image Generation**: Go's image package with custom rendering

## Architecture

### Core Components
- **Card**: Represents individual playing cards with face up/down state
- **Deck**: Manages collections of cards with shuffle/deal operations
- **Player**: Manages individual player hands and blackjack calculations  
- **Game**: Orchestrates game flow, player management, and rules
- **GameManager**: Thread-safe management of multiple concurrent games

### Thread Safety
All game operations are protected by read-write mutexes to ensure safe concurrent access across multiple games and players.

## License

This project is licensed under the terms specified in the LICENSE file.