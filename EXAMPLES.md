# API Usage Examples

This document provides comprehensive examples of how to use the Card Game API, including complete game flows for Blackjack and Cribbage, as well as custom deck creation.

## Table of Contents

- [Quick Start](#quick-start)
- [Custom Deck Examples](#custom-deck-examples)
  - [Creating a Custom Deck](#creating-a-custom-deck)
  - [Adding Custom Cards](#adding-custom-cards)
  - [Managing Custom Decks](#managing-custom-decks)
- [Complete Blackjack Game Flow](#complete-blackjack-game-flow)
- [Complete Cribbage Game Flow](#complete-cribbage-game-flow)
- [Advanced Examples](#advanced-examples)
  - [Concurrent Games](#concurrent-games)
  - [Face Up/Down Cards](#face-updown-cards)
  - [Multi-Deck Support](#multi-deck-support)
  - [Spanish 21 Support](#spanish-21-support)

## Quick Start

### Test Basic Endpoint
```bash
curl http://localhost:8080/hello
```

### View API Documentation
- **Interactive API Documentation**: http://localhost:8080/api-docs
- **OpenAPI Specification**: http://localhost:8080/openapi.yaml

## Custom Deck Examples

### Creating a Custom Deck

The API supports creating completely free-form custom decks with custom cards that can have any rank, suit, and attributes you want.

#### Step 1: Create a Custom Deck
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

### Adding Custom Cards

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

### Managing Custom Decks

#### List Custom Decks
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

#### Get Deck with All Cards
```bash
curl http://localhost:8080/custom-decks/dffd1c1d-c0f0-4c49-994c-6c7e47ffad6f/cards
```

#### Get Specific Card
```bash
curl http://localhost:8080/custom-decks/dffd1c1d-c0f0-4c49-994c-6c7e47ffad6f/cards/0
```

#### Delete Card (Tombstone)
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

## Advanced Examples

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