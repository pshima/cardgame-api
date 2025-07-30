# Card Game API

A comprehensive Go/Gin API for card games with full blackjack support, player hand tracking, multi-deck functionality, and custom deck creation.

## Features

- **Multiple Card Games**: Blackjack, Poker, War, Go Fish, Cribbage
- **Multiple Deck Types**: Standard 52-card, Spanish 21 (48-card, no 10s)
- **Custom Decks**: Create completely free-form custom decks with custom cards, suits, ranks, and attributes
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

### Option 1: Using Go
```bash
go run .
```
Server runs on `http://localhost:8080`

### Option 2: Using Docker
```bash
# Build and run with Docker
docker build -t cardgame-api .
docker run -p 8080:8080 cardgame-api

# Or use Docker Compose
docker-compose up
```

### Option 3: Using Pre-built Image
```bash
docker run -p 8080:8080 \
  -e LOG_LEVEL=INFO \
  -e GIN_MODE=release \
  cardgame-api:latest
```

### Test Basic Endpoint
```bash
curl http://localhost:8080/hello
```

### View API Documentation
- **Interactive API Documentation**: http://localhost:8080/api-docs
- **OpenAPI Specification**: http://localhost:8080/openapi.yaml



## Custom Deck Creation and Management

The API supports creating completely free-form custom decks with custom cards that can have any rank, suit, and attributes you want.

### Step 1: Create a Custom Deck
```bash
curl -X POST http://localhost:8080/custom-decks \
  -H "Content-Type: application/json" \
  -d '{"name": "My Fantasy Deck"}'
```

**Response:**
```json
{
  "id": "dffd1c1d-c0f0-4c49-994c-6c7e47ffad6f",
  "name": "My Fantasy Deck",
  "message": "Custom deck created successfully",
  "created": "2025-07-29T15:40:24.662345-07:00"
}
```

### Step 2: Add Custom Cards

#### Game-Compatible Card (can be used in traditional games)
```bash
curl -X POST http://localhost:8080/custom-decks/dffd1c1d-c0f0-4c49-994c-6c7e47ffad6f/cards \
  -H "Content-Type: application/json" \
  -d '{
    "name": "9 of Hyenas",
    "rank": 9,
    "suit": "hyenas",
    "attributes": {
      "strength": "+1",
      "luck": "-1",
      "cursed": "true"
    }
  }'
```

**Response:**
```json
{
  "index": 0,
  "name": "9 of Hyenas",
  "rank": 9,
  "suit": "hyenas",
  "game_compatible": true,
  "attributes": {
    "cursed": "true",
    "luck": "-1",
    "strength": "+1"
  },
  "message": "Card added successfully"
}
```

#### Non-Game-Compatible Card (free-form)
```bash
curl -X POST http://localhost:8080/custom-decks/dffd1c1d-c0f0-4c49-994c-6c7e47ffad6f/cards \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Magic Potion",
    "attributes": {
      "effect": "heal 10 hp",
      "rarity": "common"
    }
  }'
```

**Response:**
```json
{
  "index": 1,
  "name": "Magic Potion",
  "rank": null,
  "suit": "",
  "game_compatible": false,
  "attributes": {
    "effect": "heal 10 hp",
    "rarity": "common"
  },
  "message": "Card added successfully"
}
```

#### Custom Suit Card
```bash
curl -X POST http://localhost:8080/custom-decks/dffd1c1d-c0f0-4c49-994c-6c7e47ffad6f/cards \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Ace of 4-Leaf Clovers",
    "rank": 1,
    "suit": "4-leaf-clovers",
    "attributes": {
      "luck": "+5",
      "symbol": "🍀"
    }
  }'
```

### Step 3: List Custom Decks
```bash
curl http://localhost:8080/custom-decks
```

**Response:**
```json
{
  "count": 1,
  "decks": [
    {
      "id": "dffd1c1d-c0f0-4c49-994c-6c7e47ffad6f",
      "name": "My Fantasy Deck",
      "card_count": 3,
      "created": "2025-07-29T15:40:24.662345-07:00",
      "last_used": "2025-07-29T15:40:45.906827-07:00"
    }
  ]
}
```

### Step 4: Get Deck with All Cards
```bash
curl http://localhost:8080/custom-decks/dffd1c1d-c0f0-4c49-994c-6c7e47ffad6f/cards
```

### Step 5: Get Specific Card
```bash
curl http://localhost:8080/custom-decks/dffd1c1d-c0f0-4c49-994c-6c7e47ffad6f/cards/0
```

### Step 6: Delete Card (Tombstone)
```bash
curl -X DELETE http://localhost:8080/custom-decks/dffd1c1d-c0f0-4c49-994c-6c7e47ffad6f/cards/1
```

The card is marked as deleted but remains queryable:
```bash
curl http://localhost:8080/custom-decks/dffd1c1d-c0f0-4c49-994c-6c7e47ffad6f/cards/1
```

**Response shows `"deleted": true`:**
```json
{
  "deck_id": "dffd1c1d-c0f0-4c49-994c-6c7e47ffad6f",
  "deck_name": "My Fantasy Deck",
  "index": 1,
  "name": "Magic Potion",
  "rank": null,
  "suit": "",
  "game_compatible": false,
  "attributes": {
    "effect": "heal 10 hp",
    "rarity": "common"
  },
  "deleted": true
}
```

### Custom Deck Features

- **Free-form Cards**: Create any card with any name, rank, suit, and attributes
- **Game Compatibility**: Cards with numeric ranks and suits are automatically marked as `game_compatible: true`
- **Custom Suits**: Use traditional suits (hearts, diamonds, clubs, spades) or create your own (hyenas, 4-leaf-clovers, etc.)
- **Flexible Ranks**: Use numbers, strings, or leave blank
- **Rich Attributes**: Add up to 100 custom key-value attributes per card
- **Tombstone Deletion**: Deleted cards remain queryable but are marked as `deleted: true`
- **Auto-indexing**: Cards get sequential indices (0, 1, 2...) for easy reference
- **Security**: All inputs are validated and sanitized

### Custom Deck Limits

- **Deck Name**: 1-128 characters
- **Cards per Deck**: Maximum 2,000 cards
- **Card Name**: Maximum 100 characters
- **Suit Name**: Maximum 50 characters
- **Attributes**: Maximum 100 attributes per card
- **Attribute Keys**: Maximum 50 characters each
- **Attribute Values**: Maximum 200 characters each

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

### Cribbage Game Flow
- `GET /game/new/cribbage` - Create new cribbage game (2 players, 1 deck)
- `POST /game/:gameId/cribbage/start` - Start cribbage game (deals 6 cards each)
- `POST /game/:gameId/cribbage/discard/:playerId` - Discard 2 cards to crib `{"card_indices": [0, 1]}`
- `POST /game/:gameId/cribbage/play/:playerId` - Play a card during play phase `{"card_index": 0}`
- `POST /game/:gameId/cribbage/go/:playerId` - Say "go" when can't play without exceeding 31
- `GET /game/:gameId/cribbage/show` - Score hands and crib (moves to next deal or ends game)
- `GET /game/:gameId/cribbage/state` - Get complete cribbage game state with scores

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

### Custom Deck Management
- `POST /custom-decks` - Create custom deck `{"name": "Deck Name"}`
- `GET /custom-decks` - List all custom decks with summaries
- `GET /custom-decks/:deckId` - Get custom deck details with all cards
- `DELETE /custom-decks/:deckId` - Delete custom deck permanently
- `POST /custom-decks/:deckId/cards` - Add card to deck `{"name": "Card Name", "rank": 9, "suit": "custom", "attributes": {...}}`
- `GET /custom-decks/:deckId/cards` - List cards in deck (`?include_deleted=true` for deleted cards)
- `GET /custom-decks/:deckId/cards/:cardIndex` - Get specific card by index
- `DELETE /custom-decks/:deckId/cards/:cardIndex` - Delete card (tombstone - remains queryable)

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
- **Cribbage**: Complete cribbage implementation with pegging, hand scoring, and crib
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

## Cribbage Rules Implemented

### Game Overview
- **Players**: Exactly 2 players
- **Cards**: Each player starts with 6 cards, discards 2 to the crib
- **Goal**: First player to reach 121 points wins
- **Scoring**: Points earned during play (pegging) and hand evaluation (show)

### Game Phases
1. **Deal**: Each player receives 6 cards
2. **Discard**: Each player puts 2 cards in the crib (dealer's bonus hand)
3. **Play**: Players alternate playing cards, keeping running total ≤ 31
4. **Show**: Score hands and crib, check for winner

### Cribbage Scoring

#### During Play (Pegging)
- **Fifteen**: Playing a card that makes the total exactly 15 (2 points)
- **Thirty-One**: Playing a card that makes the total exactly 31 (2 points)
- **Pairs**: Playing a card of the same rank as the previous card (2 points)
- **Runs**: Playing cards that form a sequence with recent cards (1 point per card)
- **Go**: Last card played when opponent can't play (1 point)

#### During Show (Hand Scoring)
- **Fifteens**: Any combination of cards totaling 15 (2 points each)
- **Pairs**: Two or more cards of the same rank (2 points per pair)
- **Runs**: Three or more consecutive ranks (1 point per card)
- **Flush**: All hand cards same suit (4 points, 5 if starter matches)
- **Nobs**: Jack in hand matching starter suit (1 point)

### Card Values
- **Ace**: 1 point
- **2-10**: Face value
- **Face Cards (J, Q, K)**: 10 points each

### Crib Rules
- **Ownership**: Crib belongs to the dealer
- **Scoring**: Scored after both players' hands
- **Composition**: 4 cards (2 from each player) plus starter card
- **Dealer Alternation**: Dealer alternates each hand

### Special Rules
- **His Heels**: If starter card is a Jack, dealer gets 2 points immediately
- **Pegging Order**: Non-dealer plays first in each round
- **Go Rules**: When a player can't play without exceeding 31, they say "go"
- **Last Card**: Player who plays the last card of a round gets 1 point

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

## Complete Cribbage Game Flow

Here's a step-by-step example of running a full cribbage game:

### Step 1: Create a Cribbage Game
```bash
# Create a new cribbage game (automatically configured for 2 players with 1 standard deck)
curl "http://localhost:8080/game/new/cribbage"
```

**Response:**
```json
{
  "game_id": "456e7890-e89b-12d3-a456-426614174111",
  "game_type": "Cribbage",
  "deck_name": "Swift Eagle",
  "deck_type": "Standard",
  "max_players": 2,
  "current_players": 0,
  "message": "New Cribbage game created",
  "remaining_cards": 52,
  "created": "2025-07-29T10:30:00Z"
}
```

### Step 2: Add Two Players
```bash
# Add first player
curl -X POST "http://localhost:8080/game/456e7890-e89b-12d3-a456-426614174111/players" \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice"}'

# Add second player  
curl -X POST "http://localhost:8080/game/456e7890-e89b-12d3-a456-426614174111/players" \
  -H "Content-Type: application/json" \
  -d '{"name": "Bob"}'
```

### Step 3: Start the Cribbage Game (Deal Phase)
```bash
# This deals 6 cards to each player and moves to discard phase
curl -X POST "http://localhost:8080/game/456e7890-e89b-12d3-a456-426614174111/cribbage/start"
```

**Response:**
```json
{
  "game_id": "456e7890-e89b-12d3-a456-426614174111",
  "game_type": "Cribbage",
  "status": "in_progress",
  "phase": "discard",
  "dealer": 0,
  "current_player": 1,
  "message": "Cribbage game started"
}
```

### Step 4: Check Game State
```bash
curl "http://localhost:8080/game/456e7890-e89b-12d3-a456-426614174111/cribbage/state"
```

**Response shows each player with 6 cards:**
```json
{
  "game_id": "456e7890-e89b-12d3-a456-426614174111",
  "game_type": "Cribbage",
  "status": "in_progress",
  "phase": "discard",
  "dealer": 0,
  "current_player": 1,
  "players": [
    {
      "id": "player-uuid-alice",
      "name": "Alice",
      "hand": [
        {"rank": 1, "suit": 0, "face_up": true, "images": {...}},
        {"rank": 5, "suit": 1, "face_up": true, "images": {...}},
        {"rank": 10, "suit": 2, "face_up": true, "images": {...}},
        {"rank": 11, "suit": 3, "face_up": true, "images": {...}},
        {"rank": 4, "suit": 0, "face_up": true, "images": {...}},
        {"rank": 9, "suit": 1, "face_up": true, "images": {...}}
      ],
      "hand_size": 6,
      "score": 0
    },
    {
      "id": "player-uuid-bob", 
      "name": "Bob",
      "hand": [...],
      "hand_size": 6,
      "score": 0
    }
  ],
  "crib": [],
  "crib_size": 0,
  "play_total": 0,
  "player_scores": [0, 0],
  "game_score": 121
}
```

### Step 5: Discard Phase - Players Put 2 Cards in Crib
```bash
# Bob (non-dealer) discards first - put cards at indices 4 and 5 into crib
curl -X POST "http://localhost:8080/game/456e7890-e89b-12d3-a456-426614174111/cribbage/discard/player-uuid-bob" \
  -H "Content-Type: application/json" \
  -d '{"card_indices": [4, 5]}'

# Alice (dealer) discards next
curl -X POST "http://localhost:8080/game/456e7890-e89b-12d3-a456-426614174111/cribbage/discard/player-uuid-alice" \
  -H "Content-Type: application/json" \
  -d '{"card_indices": [3, 5]}'
```

**Response after both players discard:**
```json
{
  "game_id": "456e7890-e89b-12d3-a456-426614174111",
  "player_id": "player-uuid-alice", 
  "player_name": "Alice",
  "phase": "play",
  "starter": {
    "rank": 6,
    "suit": 2,
    "face_up": true,
    "images": {...}
  },
  "message": "Cards discarded, starter cut, play phase begun"
}
```

### Step 6: Play Phase - Alternate Playing Cards (Pegging)
```bash
# Bob plays first card (non-dealer leads)
curl -X POST "http://localhost:8080/game/456e7890-e89b-12d3-a456-426614174111/cribbage/play/player-uuid-bob" \
  -H "Content-Type: application/json" \
  -d '{"card_index": 0}'

# Alice plays a card
curl -X POST "http://localhost:8080/game/456e7890-e89b-12d3-a456-426614174111/cribbage/play/player-uuid-alice" \
  -H "Content-Type: application/json" \
  -d '{"card_index": 1}'
```

**Response:**
```json
{
  "game_id": "456e7890-e89b-12d3-a456-426614174111",
  "player_id": "player-uuid-alice",
  "player_name": "Alice",
  "play_total": 15,
  "play_count": 2,
  "player_score": 2,
  "phase": "play",
  "current_player": 1,
  "message": "Card played"
}
```

### Step 7: Continue Play Until Cards Exhausted
Players continue alternating plays. If a player can't play without exceeding 31, they say "go":

```bash
# Player says "go" when they can't play
curl -X POST "http://localhost:8080/game/456e7890-e89b-12d3-a456-426614174111/cribbage/go/player-uuid-bob"
```

### Step 8: Show Phase - Score Hands and Crib
```bash
# After all cards are played, score the hands
curl "http://localhost:8080/game/456e7890-e89b-12d3-a456-426614174111/cribbage/show"
```

**Response:**
```json
{
  "game_id": "456e7890-e89b-12d3-a456-426614174111",
  "scores": {
    "player-uuid-bob": 8,
    "player-uuid-alice": 12,
    "crib": 6
  },
  "player_scores": [18, 8],
  "phase": "deal",
  "status": "in_progress"
}
```

### Step 9: Continue Until Game Ends at 121 Points
The game continues with new deals until one player reaches 121 points:

```json
{
  "game_id": "456e7890-e89b-12d3-a456-426614174111",
  "scores": {...},
  "player_scores": [121, 95],
  "phase": "finished",
  "status": "finished",
  "winner": "Alice",
  "winner_id": "player-uuid-alice"
}
```

## Observability and Debugging

The Card Game API includes comprehensive observability features to monitor application performance and troubleshoot issues in production.

### Logging

The application uses structured JSON logging with configurable log levels:

#### Environment Variables
- `LOG_LEVEL`: Set logging level (DEBUG, INFO, WARN, ERROR) - defaults to INFO
- `TRUSTED_PROXIES`: Comma-separated list of trusted proxy IPs for correct client IP extraction

#### Log Levels and Usage
- **DEBUG**: Detailed operations (card dealing, game state changes, validation steps)
- **INFO**: Normal operations (requests, game creation, successful operations)
- **WARN**: Validation failures, invalid requests, recoverable errors
- **ERROR**: System errors, failed operations, server failures

#### Example Log Output
```json
{
  "level": "info",
  "time": "2025-07-29T18:18:51-07:00",
  "caller": "cardgame-api/main.go:242",
  "message": "Request completed",
  "method": "GET",
  "path": "/game/abc123/deal",
  "query": "",
  "user_agent": "curl/8.7.1",
  "client_ip": "192.168.1.100",
  "game_id": "abc123",
  "status_code": 200,
  "latency": 0.000101375,
  "latency_human": "101.375µs"
}
```

### Metrics

OpenTelemetry metrics are collected for monitoring application performance:

#### Available Metrics
- `http_requests_total`: Counter of total HTTP requests by method, path, status code
- `http_request_duration_seconds`: Histogram of request latency in seconds
- `http_requests_in_flight`: Current number of HTTP requests being processed
- `active_games`: Current number of active games
- `active_custom_decks`: Current number of custom decks
- `cards_dealt_total`: Counter of total cards dealt
- `games_created_total`: Counter of total games created
- `api_errors_total`: Counter of API errors (5xx status codes)

#### Metrics Endpoints
- **Prometheus Format**: `GET /metrics` - Standard Prometheus metrics format
- **JSON Stats**: `GET /stats` - Human-readable JSON metrics summary

#### Example Stats Response
```json
{
  "service": {
    "name": "cardgame-api",
    "version": "1.0.0",
    "uptime": "2h34m15s"
  },
  "games": {
    "active_count": 15,
    "total_created": 247
  },
  "metrics": {
    "http_requests_total": 1523,
    "cards_dealt_total": 8934,
    "api_errors_total": 12
  },
  "system": {
    "timestamp": "2025-07-29T20:52:15Z",
    "log_level": "info"
  }
}
```

### Debugging Common Issues

#### Game Not Found Errors
1. Check logs for game creation: `game_id` field in request logs
2. Verify game exists: `curl http://localhost:8080/games`
3. Check if game was deleted: Look for DELETE operations in logs

#### Card Dealing Issues
- **No cards remaining**: Check `remaining_cards` in game state
- **Invalid game state**: Verify game status is not "finished"
- **Player not found**: Confirm player was added to game

#### Performance Issues
1. Monitor `/metrics` endpoint for high latency
2. Check `http_requests_in_flight` for request queuing
3. Review `http_request_duration_seconds` histogram

#### Log Analysis
```bash
# Filter by game ID
cat app.log | jq 'select(.game_id == "your-game-id")'

# Monitor error rates
cat app.log | jq 'select(.level == "error")'

# Track API performance
cat app.log | jq 'select(.latency != null) | {path, latency_human, status_code}'

# Watch live logs
tail -f app.log | jq 'select(.level == "error" or .status_code >= 400)'
```

#### Environment Configuration for Production
```bash
# Set log level to reduce verbosity
export LOG_LEVEL=INFO

# Configure trusted proxies for load balancers
export TRUSTED_PROXIES="10.0.1.100,192.168.1.0/24"

# Set Gin to release mode
export GIN_MODE=release
```

#### Monitoring Recommendations
1. Set up alerts on `api_errors_total` metric spikes
2. Monitor `http_request_duration_seconds` 95th percentile
3. Track `active_games` growth over time
4. Alert on sustained high `http_requests_in_flight`

### Common Log Patterns for Debugging

#### Successful Game Flow
```
INFO Request started: POST /game/new
DEBUG Creating new game: decks=1, type=standard
INFO Game created successfully: game_id=abc123
INFO Request completed: status_code=200
```

#### Card Dealing Sequence
```
INFO Request started: GET /game/abc123/deal
DEBUG Dealing card: game_id=abc123, remaining_cards=51
INFO Card dealt successfully: card="Ace of Hearts", remaining_cards=50
INFO Request completed: status_code=200
```

#### Error Scenarios
```
WARN Invalid game ID provided: game_id=invalid-id
WARN Game not found: game_id=nonexistent-game
ERROR Request failed: status_code=404
```

## Container Deployment

The Card Game API is fully containerized with Docker support for easy deployment and scaling.

### Docker Images

The application uses a multi-stage build process to create minimal, secure container images:
- **Build stage**: Uses Go 1.24.4 Alpine for compilation
- **Runtime stage**: Uses Alpine Linux with non-root user for security
- **Final image size**: ~15MB (excluding static assets)

### Building the Container

```bash
# Build the Docker image
docker build -t cardgame-api:latest .

# Build with custom tags
docker build -t cardgame-api:v1.0.0 -t cardgame-api:latest .
```

### Running with Docker

```bash
# Basic run
docker run -p 8080:8080 cardgame-api:latest

# Run with custom configuration
docker run -d \
  --name cardgame-api \
  -p 8080:8080 \
  -e LOG_LEVEL=DEBUG \
  -e PORT=8080 \
  -e GIN_MODE=release \
  -e TRUSTED_PROXIES="10.0.0.0/8" \
  --restart unless-stopped \
  cardgame-api:latest

# View logs
docker logs -f cardgame-api

# Check health
docker inspect cardgame-api --format='{{.State.Health.Status}}'
```

### Docker Compose

The included `docker-compose.yml` provides a complete development environment:

```bash
# Start the application
docker-compose up -d

# View logs
docker-compose logs -f

# Stop the application
docker-compose down

# Start with monitoring stack (Prometheus + Grafana)
docker-compose --profile monitoring up -d
```

### Environment Variables

Configure the container using these environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `LOG_LEVEL` | Logging level (DEBUG, INFO, WARN, ERROR) | `INFO` |
| `GIN_MODE` | Gin framework mode (debug, release) | `release` |
| `TRUSTED_PROXIES` | Comma-separated trusted proxy IPs | `""` |

### Health Checks

The container includes built-in health checks:
- **Endpoint**: `/hello`
- **Interval**: 30 seconds
- **Timeout**: 3 seconds
- **Retries**: 3

### Container Best Practices

1. **Security**:
   - Runs as non-root user (UID 1000)
   - Minimal Alpine base image
   - No unnecessary packages

2. **Resource Limits** (via docker-compose):
   - CPU: 1 core limit, 0.5 core reservation
   - Memory: 512MB limit, 256MB reservation

3. **Logging**:
   - JSON structured logs to stdout/stderr
   - Log rotation configured in docker-compose
   - Compatible with log aggregation systems

4. **Monitoring**:
   - Prometheus metrics at `/metrics`
   - Optional Prometheus/Grafana stack
   - Built-in health endpoint

### Kubernetes Deployment

Example Kubernetes deployment:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: cardgame-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: cardgame-api
  template:
    metadata:
      labels:
        app: cardgame-api
    spec:
      containers:
      - name: cardgame-api
        image: cardgame-api:latest
        ports:
        - containerPort: 8080
        env:
        - name: LOG_LEVEL
          value: "INFO"
        - name: GIN_MODE
          value: "release"
        livenessProbe:
          httpGet:
            path: /hello
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /hello
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          limits:
            cpu: "1"
            memory: "512Mi"
          requests:
            cpu: "500m"
            memory: "256Mi"
```

### Production Considerations

1. **Reverse Proxy**: Use nginx/traefik for SSL termination
2. **Secrets**: Use Docker secrets or Kubernetes secrets for sensitive data
3. **Persistence**: Games are in-memory; add Redis for persistence if needed
4. **Scaling**: Stateless design allows horizontal scaling
5. **Monitoring**: Connect to your existing Prometheus/Grafana stack

## Container Security

The Card Game API implements comprehensive security measures for containerized deployments.

### Security Features

#### 1. Secure Docker Image
- **Multi-stage build**: Minimal attack surface
- **Non-root user**: Runs as UID 65532
- **No shell**: Reduced exploitation risk
- **Read-only filesystem**: Prevents runtime modifications
- **Security labels**: Enables automated scanning

#### 2. Security Scanning
Run comprehensive security scans:
```bash
# Run all security scans
make security-scan

# Build with security validation
make docker-build-secure

# Run with enhanced security
make compose-security
```

#### 3. Secrets Management
Generate and manage secrets securely:
```bash
# Generate secure secrets
./scripts/generate-secrets.sh

# Run with secrets
docker-compose -f docker-compose.secrets.yml up
```

#### 4. Security Configurations

**docker-compose.security.yml** includes:
- Read-only root filesystem
- Dropped Linux capabilities
- No new privileges flag
- Seccomp and AppArmor profiles
- Resource limits (CPU, memory, PIDs)
- Network isolation
- Localhost-only port binding

### Security Best Practices

#### Image Security
1. **Regular Updates**: Rebuild images weekly for security patches
2. **Vulnerability Scanning**: Use Trivy/Clair before deployment
3. **SBOM Generation**: Track all dependencies
4. **Minimal Base**: Alpine Linux with only essential packages

#### Runtime Security
1. **Network Policies**: Restrict egress/ingress traffic
2. **Resource Limits**: Prevent DoS attacks
3. **Health Checks**: Ensure service availability
4. **Audit Logging**: Track all API access

#### Secrets Management
1. **Never hardcode**: Use environment variables or files
2. **Rotate regularly**: Implement key rotation
3. **Encrypt at rest**: Use Docker secrets or Kubernetes secrets
4. **Minimal exposure**: Mount secrets as read-only files

### Security Checklist

Before deploying to production:

- [ ] Run `make security-scan` and fix all HIGH/CRITICAL issues
- [ ] Generate new secrets with `./scripts/generate-secrets.sh`
- [ ] Review Dockerfile with `hadolint Dockerfile`
- [ ] Enable read-only root filesystem
- [ ] Configure resource limits
- [ ] Set up network policies
- [ ] Enable audit logging
- [ ] Configure HTTPS/TLS
- [ ] Set up WAF rules
- [ ] Implement rate limiting

### Kubernetes Security

Example secure Kubernetes deployment:

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: cardgame-api
  namespace: cardgame
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: cardgame-api
  namespace: cardgame
spec:
  replicas: 3
  selector:
    matchLabels:
      app: cardgame-api
  template:
    metadata:
      labels:
        app: cardgame-api
    spec:
      serviceAccountName: cardgame-api
      securityContext:
        runAsNonRoot: true
        runAsUser: 65532
        fsGroup: 65532
        seccompProfile:
          type: RuntimeDefault
      containers:
      - name: cardgame-api
        image: cardgame-api:latest
        imagePullPolicy: Always
        ports:
        - containerPort: 8080
          protocol: TCP
        securityContext:
          allowPrivilegeEscalation: false
          readOnlyRootFilesystem: true
          runAsNonRoot: true
          runAsUser: 65532
          capabilities:
            drop:
            - ALL
        env:
        - name: LOG_LEVEL
          value: "INFO"
        - name: GIN_MODE
          value: "release"
        - name: API_KEY
          valueFrom:
            secretKeyRef:
              name: cardgame-secrets
              key: api-key
        volumeMounts:
        - name: tmp
          mountPath: /tmp
        livenessProbe:
          httpGet:
            path: /hello
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /hello
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          limits:
            cpu: "1"
            memory: "512Mi"
          requests:
            cpu: "500m"
            memory: "256Mi"
      volumes:
      - name: tmp
        emptyDir: {}
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: cardgame-api-netpol
  namespace: cardgame
spec:
  podSelector:
    matchLabels:
      app: cardgame-api
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - namespaceSelector:
        matchLabels:
          name: ingress-nginx
    ports:
    - protocol: TCP
      port: 8080
  egress:
  - to:
    - namespaceSelector: {}
    ports:
    - protocol: TCP
      port: 53  # DNS
    - protocol: UDP
      port: 53  # DNS
```

### Compliance

The security implementation helps meet common compliance requirements:
- **PCI DSS**: Secure configuration, access controls, monitoring
- **HIPAA**: Encryption, audit logs, access management
- **SOC 2**: Security controls, monitoring, incident response
- **GDPR**: Data protection, audit trails, secure processing

## License

This project is licensed under the terms specified in the LICENSE file.