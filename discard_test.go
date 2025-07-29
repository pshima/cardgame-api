package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDealToPlayerEndpoint(t *testing.T) {
	router := setupPlayerRouter()

	// Create a game and add a player
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/game/new/1/standard/4", nil)
	router.ServeHTTP(w, req)

	var gameResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &gameResponse)
	gameID := gameResponse["game_id"].(string)

	// Add player
	playerData := map[string]string{"name": "Alice"}
	jsonData, _ := json.Marshal(playerData)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/game/"+gameID+"/players", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var addResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &addResponse)
	player := addResponse["player"].(map[string]interface{})
	playerID := player["id"].(string)

	// Deal card to player (face down by default)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/game/"+gameID+"/deal/player/"+playerID, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var dealResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &dealResponse)
	assert.NoError(t, err)

	assert.Equal(t, gameID, dealResponse["game_id"])
	assert.Equal(t, playerID, dealResponse["player_id"])
	assert.Equal(t, "Alice", dealResponse["player_name"])
	assert.Contains(t, dealResponse, "card")
	assert.Equal(t, float64(1), dealResponse["hand_size"])
	assert.Equal(t, float64(51), dealResponse["remaining_cards"])

	// Deal to dealer
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/game/"+gameID+"/deal/player/dealer", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var dealerResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &dealerResponse)
	assert.Equal(t, "dealer", dealerResponse["player_id"])
	assert.Equal(t, "Dealer", dealerResponse["player_name"])

	// Try to deal to non-existent player
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/game/"+gameID+"/deal/player/nonexistent", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestDealToPlayerFaceUpEndpoint(t *testing.T) {
	router := setupPlayerRouter()

	// Create a game and add a player
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/game/new/1/standard/4", nil)
	router.ServeHTTP(w, req)

	var gameResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &gameResponse)
	gameID := gameResponse["game_id"].(string)

	// Add player
	playerData := map[string]string{"name": "Alice"}
	jsonData, _ := json.Marshal(playerData)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/game/"+gameID+"/players", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var addResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &addResponse)
	player := addResponse["player"].(map[string]interface{})
	playerID := player["id"].(string)

	// Deal face up card
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/game/"+gameID+"/deal/player/"+playerID+"/true", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var dealResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &dealResponse)
	assert.NoError(t, err)

	assert.Equal(t, true, dealResponse["face_up"])
	card := dealResponse["card"].(map[string]interface{})
	assert.Equal(t, true, card["face_up"])

	// Deal face down card
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/game/"+gameID+"/deal/player/"+playerID+"/false", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var dealResponse2 map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &dealResponse2)
	assert.Equal(t, false, dealResponse2["face_up"])
	card2 := dealResponse2["card"].(map[string]interface{})
	assert.Equal(t, false, card2["face_up"])

	// Test with "1" for true
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/game/"+gameID+"/deal/player/"+playerID+"/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var dealResponse3 map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &dealResponse3)
	assert.Equal(t, true, dealResponse3["face_up"])
}

func TestDiscardToCardEndpoint(t *testing.T) {
	router := setupPlayerRouter()

	// Create a game and add a player
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/game/new/1/standard/4", nil)
	router.ServeHTTP(w, req)

	var gameResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &gameResponse)
	gameID := gameResponse["game_id"].(string)

	// Add player
	playerData := map[string]string{"name": "Alice"}
	jsonData, _ := json.Marshal(playerData)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/game/"+gameID+"/players", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var addResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &addResponse)
	player := addResponse["player"].(map[string]interface{})
	playerID := player["id"].(string)

	// Deal some cards to player
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/game/"+gameID+"/deal/player/"+playerID+"/true", nil)
	router.ServeHTTP(w, req)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/game/"+gameID+"/deal/player/"+playerID+"/false", nil)
	router.ServeHTTP(w, req)

	// Now player should have 2 cards, discard the first one (index 0)
	discardData := map[string]interface{}{
		"player_id":  playerID,
		"card_index": 0,
	}
	jsonData, _ = json.Marshal(discardData)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/game/"+gameID+"/discard/main", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var discardResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &discardResponse)
	assert.NoError(t, err)

	assert.Equal(t, gameID, discardResponse["game_id"])
	assert.Equal(t, playerID, discardResponse["player_id"])
	assert.Equal(t, "Alice", discardResponse["player_name"])
	assert.Equal(t, "main", discardResponse["pile_id"])
	assert.Equal(t, "Main Discard Pile", discardResponse["pile_name"])
	assert.Equal(t, float64(1), discardResponse["pile_size"])
	assert.Equal(t, float64(1), discardResponse["hand_size"])
	assert.Contains(t, discardResponse, "card")

	// Test discarding to non-existent pile
	discardData2 := map[string]interface{}{
		"player_id":  playerID,
		"card_index": 0,
	}
	jsonData2, _ := json.Marshal(discardData2)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/game/"+gameID+"/discard/nonexistent", bytes.NewBuffer(jsonData2))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)

	// Test discarding from non-existent player
	discardData3 := map[string]interface{}{
		"player_id":  "nonexistent",
		"card_index": 0,
	}
	jsonData3, _ := json.Marshal(discardData3)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/game/"+gameID+"/discard/main", bytes.NewBuffer(jsonData3))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)

	// Test discarding with invalid card index
	discardData4 := map[string]interface{}{
		"player_id":  playerID,
		"card_index": 10,
	}
	jsonData4, _ := json.Marshal(discardData4)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/game/"+gameID+"/discard/main", bytes.NewBuffer(jsonData4))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestCompleteBlackjackGameWorkflow(t *testing.T) {
	router := setupPlayerRouter()

	// Create a blackjack game with 2 decks for 4 players
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/game/new/2/standard/4", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var gameResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &gameResponse)
	gameID := gameResponse["game_id"].(string)
	assert.Equal(t, "Blackjack", gameResponse["game_type"])
	assert.Equal(t, float64(104), gameResponse["remaining_cards"])

	// Add two players
	players := []string{"Alice", "Bob"}
	var playerIDs []string

	for _, name := range players {
		playerData := map[string]string{"name": name}
		jsonData, _ := json.Marshal(playerData)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/game/"+gameID+"/players", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)

		var addResponse map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &addResponse)
		player := addResponse["player"].(map[string]interface{})
		playerIDs = append(playerIDs, player["id"].(string))
	}

	// Shuffle the deck
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/game/"+gameID+"/shuffle", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	// Deal initial blackjack hands (2 cards each, first face down, second face up)
	for _, playerID := range playerIDs {
		// First card face down
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/game/"+gameID+"/deal/player/"+playerID+"/false", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)

		// Second card face up
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/game/"+gameID+"/deal/player/"+playerID+"/true", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
	}

	// Deal dealer cards (first face down, second face up)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/game/"+gameID+"/deal/player/dealer/false", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/game/"+gameID+"/deal/player/dealer/true", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	// Check game state
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/game/"+gameID+"/state", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var stateResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &stateResponse)

	assert.Equal(t, gameID, stateResponse["game_id"])
	assert.Equal(t, "Blackjack", stateResponse["game_type"])
	assert.Equal(t, float64(2), stateResponse["current_players"])
	assert.Equal(t, float64(98), stateResponse["remaining_cards"]) // 104 - 6 cards dealt

	// Check players have cards
	players_data := stateResponse["players"].([]interface{})
	assert.Equal(t, 2, len(players_data))

	for _, playerData := range players_data {
		player := playerData.(map[string]interface{})
		hand := player["hand"].([]interface{})
		assert.Equal(t, 2, len(hand))

		// Check face up/down status
		card1 := hand[0].(map[string]interface{})
		card2 := hand[1].(map[string]interface{})
		assert.Equal(t, false, card1["face_up"]) // First card face down
		assert.Equal(t, true, card2["face_up"])  // Second card face up
	}

	// Check dealer has cards
	dealer := stateResponse["dealer"].(map[string]interface{})
	dealerHand := dealer["hand"].([]interface{})
	assert.Equal(t, 2, len(dealerHand))

	// Simulate player discarding a card
	discardData := map[string]interface{}{
		"player_id":  playerIDs[0],
		"card_index": 0,
	}
	jsonData, _ := json.Marshal(discardData)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/game/"+gameID+"/discard/main", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	// Check final state
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/game/"+gameID+"/state", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var finalState map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &finalState)

	// Check discard pile has a card
	discardPiles := finalState["discard_piles"].([]interface{})
	mainPile := discardPiles[0].(map[string]interface{})
	assert.Equal(t, float64(1), mainPile["size"])

	// Check first player now has 1 card
	finalPlayers := finalState["players"].([]interface{})
	firstPlayer := finalPlayers[0].(map[string]interface{})
	firstPlayerHand := firstPlayer["hand"].([]interface{})
	assert.Equal(t, 1, len(firstPlayerHand))
}

func TestSpanish21GameWithPlayers(t *testing.T) {
	router := setupPlayerRouter()

	// Create a Spanish21 game
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/game/new/1/spanish21/3", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var gameResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &gameResponse)
	gameID := gameResponse["game_id"].(string)
	assert.Equal(t, "Spanish21", gameResponse["deck_type"])
	assert.Equal(t, float64(48), gameResponse["remaining_cards"])

	// Add player
	playerData := map[string]string{"name": "Carlos"}
	jsonData, _ := json.Marshal(playerData)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/game/"+gameID+"/players", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var addResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &addResponse)
	player := addResponse["player"].(map[string]interface{})
	playerID := player["id"].(string)

	// Deal several cards and verify no 10s are dealt
	dealtCards := []interface{}{}
	for i := 0; i < 10; i++ {
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/game/"+gameID+"/deal/player/"+playerID+"/true", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)

		var dealResponse map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &dealResponse)
		card := dealResponse["card"].(map[string]interface{})
		dealtCards = append(dealtCards, card)

		// Verify no 10s (rank 10 = Ten constant)
		rank := card["rank"].(float64)
		assert.NotEqual(t, float64(Ten), rank, "Spanish21 should not deal 10s")
	}

	// Check game state
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/game/"+gameID+"/state", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var stateResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &stateResponse)
	assert.Equal(t, float64(38), stateResponse["remaining_cards"]) // 48 - 10 dealt
	assert.Equal(t, "Spanish21", stateResponse["deck_type"])
}