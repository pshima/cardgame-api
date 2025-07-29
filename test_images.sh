#!/bin/bash

# Test script to verify card images are working

echo "Testing Card Game API Image Functionality"
echo "========================================"

# Start the server in background if not already running
if ! curl -s http://localhost:8080/hello > /dev/null 2>&1; then
    echo "Starting server..."
    go run . &
    SERVER_PID=$!
    sleep 2
else
    echo "Server already running"
fi

# Create a new game
echo -e "\n1. Creating new game..."
GAME_RESPONSE=$(curl -s http://localhost:8080/game/new)
GAME_ID=$(echo $GAME_RESPONSE | grep -o '"game_id":"[^"]*"' | cut -d'"' -f4)
echo "Game ID: $GAME_ID"

# Deal a single card
echo -e "\n2. Dealing a single card..."
CARD_RESPONSE=$(curl -s "http://localhost:8080/game/$GAME_ID/deal")
echo $CARD_RESPONSE | python3 -m json.tool | grep -A 5 '"images"'

# Deal multiple cards
echo -e "\n3. Dealing 5 cards..."
CARDS_RESPONSE=$(curl -s "http://localhost:8080/game/$GAME_ID/deal/5")
echo "Cards dealt. First card images:"
echo $CARDS_RESPONSE | python3 -m json.tool | head -30 | grep -A 5 '"images"' | head -6

# Add a player
echo -e "\n4. Adding a player..."
curl -s -X POST "http://localhost:8080/game/$GAME_ID/players" \
  -H "Content-Type: application/json" \
  -d '{"name": "Test Player"}' > /dev/null

PLAYER_RESPONSE=$(curl -s -X POST "http://localhost:8080/game/$GAME_ID/players" \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice"}')
PLAYER_ID=$(echo $PLAYER_RESPONSE | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
echo "Player ID: $PLAYER_ID"

# Deal to player face down
echo -e "\n5. Dealing face-down card to player..."
FACEDOWN_RESPONSE=$(curl -s "http://localhost:8080/game/$GAME_ID/deal/player/$PLAYER_ID/false")
echo "Face-down card images (should show back):"
echo $FACEDOWN_RESPONSE | python3 -m json.tool | grep -A 5 '"images"'

# Deal to player face up
echo -e "\n6. Dealing face-up card to player..."
FACEUP_RESPONSE=$(curl -s "http://localhost:8080/game/$GAME_ID/deal/player/$PLAYER_ID/true")
echo "Face-up card images:"
echo $FACEUP_RESPONSE | python3 -m json.tool | grep -A 5 '"images"'

# Get game state
echo -e "\n7. Getting game state with all cards..."
STATE_RESPONSE=$(curl -s "http://localhost:8080/game/$GAME_ID/state")
echo "Number of image URLs in response:"
echo $STATE_RESPONSE | grep -o '"icon"' | wc -l

# Test image accessibility
echo -e "\n8. Testing image accessibility..."
FIRST_IMAGE_URL=$(echo $FACEUP_RESPONSE | grep -o '"icon":"[^"]*"' | cut -d'"' -f4 | head -1)
echo "Testing URL: $FIRST_IMAGE_URL"
HTTP_STATUS=$(curl -s -o /dev/null -w "%{http_code}" "$FIRST_IMAGE_URL")
if [ "$HTTP_STATUS" == "200" ]; then
    echo "✓ Image is accessible (HTTP $HTTP_STATUS)"
else
    echo "✗ Image not accessible (HTTP $HTTP_STATUS)"
fi

# Cleanup
echo -e "\n9. Cleaning up..."
curl -s -X DELETE "http://localhost:8080/game/$GAME_ID" > /dev/null
echo "Game deleted"

# Kill server if we started it
if [ ! -z "$SERVER_PID" ]; then
    echo "Stopping server..."
    kill $SERVER_PID 2>/dev/null
fi

echo -e "\nTest complete!"