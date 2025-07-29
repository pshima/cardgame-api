package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupPlayerRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	gameManager = NewGameManager()
	r := gin.Default()

	r.GET("/deck-types", listDeckTypes)
	r.GET("/game/new", createNewGame)
	r.GET("/game/new/:decks", createNewGameWithDecks)
	r.GET("/game/new/:decks/:type", createNewGameWithType)
	r.GET("/game/new/:decks/:type/:players", createNewGameWithPlayers)
	r.GET("/game/:gameId/shuffle", shuffleDeck)
	r.GET("/game/:gameId", getGameInfo)
	r.GET("/game/:gameId/state", getGameState)
	r.POST("/game/:gameId/players", addPlayer)
	r.DELETE("/game/:gameId/players/:playerId", removePlayer)
	r.GET("/game/:gameId/deal", dealCard)
	r.GET("/game/:gameId/deal/:count", dealCards)
	r.GET("/game/:gameId/deal/player/:playerId", dealToPlayer)
	r.GET("/game/:gameId/deal/player/:playerId/:faceUp", dealToPlayerFaceUp)
	r.POST("/game/:gameId/discard/:pileId", discardToCard)
	r.GET("/game/:gameId/reset", resetDeck)
	r.GET("/game/:gameId/reset/:decks", resetDeckWithDecks)
	r.GET("/game/:gameId/reset/:decks/:type", resetDeckWithType)
	r.DELETE("/game/:gameId", deleteGame)
	r.GET("/games", listGames)

	return r
}

func TestCardFaceUpDown(t *testing.T) {
	card := Card{Rank: Ace, Suit: Hearts, FaceUp: true}
	assert.Equal(t, Ace, card.Rank)
	assert.Equal(t, Hearts, card.Suit)
	assert.True(t, card.FaceUp)
	
	card2 := Card{Rank: King, Suit: Spades, FaceUp: false}
	assert.False(t, card2.FaceUp)
}

func TestPlayerHandManagement(t *testing.T) {
	player := &Player{
		ID:   "test-player",
		Name: "Test Player",
		Hand: []*Card{},
	}
	
	// Test empty hand
	assert.Equal(t, 0, player.HandSize())
	
	// Add cards
	card1 := &Card{Rank: Ace, Suit: Hearts, FaceUp: true}
	card2 := &Card{Rank: King, Suit: Spades, FaceUp: false}
	
	player.AddCard(card1)
	player.AddCard(card2)
	assert.Equal(t, 2, player.HandSize())
	
	// Remove card by index
	removedCard := player.RemoveCard(0)
	assert.Equal(t, card1, removedCard)
	assert.Equal(t, 1, player.HandSize())
	
	// Invalid index
	nilCard := player.RemoveCard(10)
	assert.Nil(t, nilCard)
	
	// Clear hand
	clearedCards := player.ClearHand()
	assert.Equal(t, 1, len(clearedCards))
	assert.Equal(t, card2, clearedCards[0])
	assert.Equal(t, 0, player.HandSize())
}

func TestDiscardPile(t *testing.T) {
	pile := &DiscardPile{
		ID:    "test-pile",
		Name:  "Test Pile",
		Cards: []*Card{},
	}
	
	// Test empty pile
	assert.Equal(t, 0, pile.Size())
	assert.Nil(t, pile.TopCard())
	assert.Nil(t, pile.TakeTopCard())
	
	// Add cards
	card1 := &Card{Rank: Ace, Suit: Hearts, FaceUp: true}
	card2 := &Card{Rank: King, Suit: Spades, FaceUp: false}
	
	pile.AddCard(card1)
	pile.AddCard(card2)
	assert.Equal(t, 2, pile.Size())
	
	// Check top card
	topCard := pile.TopCard()
	assert.Equal(t, card2, topCard)
	assert.Equal(t, 2, pile.Size()) // Should not remove
	
	// Take top card
	takenCard := pile.TakeTopCard()
	assert.Equal(t, card2, takenCard)
	assert.Equal(t, 1, pile.Size())
	
	// Add multiple cards
	newCards := []*Card{
		{Rank: Queen, Suit: Diamonds, FaceUp: true},
		{Rank: Jack, Suit: Clubs, FaceUp: false},
	}
	pile.AddCards(newCards)
	assert.Equal(t, 3, pile.Size())
	
	// Clear pile
	clearedCards := pile.Clear()
	assert.Equal(t, 3, len(clearedCards))
	assert.Equal(t, 0, pile.Size())
}

func TestGameTypeEnum(t *testing.T) {
	tests := []struct {
		gameType GameType
		expected string
	}{
		{Blackjack, "Blackjack"},
		{Poker, "Poker"},
		{War, "War"},
		{GoFish, "GoFish"},
		{GameType(99), "Blackjack"}, // Default case
	}

	for _, test := range tests {
		result := test.gameType.String()
		assert.Equal(t, test.expected, result)
	}
}

func TestNewGameWithType(t *testing.T) {
	game := NewGameWithType(2, Spanish21, Poker, 4)
	
	assert.NotEmpty(t, game.ID)
	assert.Equal(t, Poker, game.GameType)
	assert.Equal(t, Spanish21, game.Deck.DeckType)
	assert.Equal(t, 96, game.Deck.RemainingCards()) // 2 Spanish21 decks
	assert.Equal(t, 4, game.MaxPlayers)
	assert.Equal(t, 0, len(game.Players))
	assert.NotNil(t, game.Dealer)
	assert.Equal(t, "dealer", game.Dealer.ID)
	assert.Equal(t, "Dealer", game.Dealer.Name)
	assert.Equal(t, 1, len(game.DiscardPiles))
	assert.NotNil(t, game.DiscardPiles["main"])
}

func TestGamePlayerManagement(t *testing.T) {
	game := NewGameWithType(1, Standard, Blackjack, 3)
	
	// Add players
	player1 := game.AddPlayer("Alice")
	assert.NotNil(t, player1)
	assert.Equal(t, "Alice", player1.Name)
	assert.Equal(t, 1, len(game.Players))
	
	player2 := game.AddPlayer("Bob")
	assert.NotNil(t, player2)
	assert.Equal(t, 2, len(game.Players))
	
	player3 := game.AddPlayer("Charlie")
	assert.NotNil(t, player3)
	assert.Equal(t, 3, len(game.Players))
	
	// Try to add beyond max
	player4 := game.AddPlayer("David")
	assert.Nil(t, player4)
	assert.Equal(t, 3, len(game.Players))
	
	// Get player
	foundPlayer := game.GetPlayer(player1.ID)
	assert.Equal(t, player1, foundPlayer)
	
	// Get dealer
	dealer := game.GetPlayer("dealer")
	assert.Equal(t, game.Dealer, dealer)
	
	// Get non-existent player
	notFound := game.GetPlayer("nonexistent")
	assert.Nil(t, notFound)
	
	// Remove player
	removed := game.RemovePlayer(player2.ID)
	assert.True(t, removed)
	assert.Equal(t, 2, len(game.Players))
	
	// Try to remove non-existent player
	notRemoved := game.RemovePlayer("nonexistent")
	assert.False(t, notRemoved)
}

func TestGameDealToPlayer(t *testing.T) {
	game := NewGameWithType(1, Standard, Blackjack, 2)
	player := game.AddPlayer("Test Player")
	
	// Deal face down
	card1 := game.DealToPlayer(player.ID, false)
	assert.NotNil(t, card1)
	assert.False(t, card1.FaceUp)
	assert.Equal(t, 1, player.HandSize())
	assert.Equal(t, 51, game.Deck.RemainingCards())
	
	// Deal face up
	card2 := game.DealToPlayer(player.ID, true)
	assert.NotNil(t, card2)
	assert.True(t, card2.FaceUp)
	assert.Equal(t, 2, player.HandSize())
	
	// Deal to dealer
	dealerCard := game.DealToPlayer("dealer", false)
	assert.NotNil(t, dealerCard)
	assert.Equal(t, 1, game.Dealer.HandSize())
	
	// Deal to non-existent player
	noCard := game.DealToPlayer("nonexistent", false)
	assert.Nil(t, noCard)
}

func TestGameDiscardPileManagement(t *testing.T) {
	game := NewGameWithType(1, Standard, Blackjack, 2)
	
	// Should have default main pile
	mainPile := game.GetDiscardPile("main")
	assert.NotNil(t, mainPile)
	assert.Equal(t, "main", mainPile.ID)
	assert.Equal(t, "Main Discard Pile", mainPile.Name)
	
	// Add new pile
	newPile := game.AddDiscardPile("side", "Side Pile")
	assert.NotNil(t, newPile)
	assert.Equal(t, "side", newPile.ID)
	assert.Equal(t, "Side Pile", newPile.Name)
	
	// Try to add duplicate pile
	duplicate := game.AddDiscardPile("side", "Duplicate")
	assert.Nil(t, duplicate)
	
	// Get pile
	foundPile := game.GetDiscardPile("side")
	assert.Equal(t, newPile, foundPile)
	
	// Get non-existent pile
	notFound := game.GetDiscardPile("nonexistent")
	assert.Nil(t, notFound)
}

func TestCreateNewGameWithPlayersEndpoint(t *testing.T) {
	router := setupPlayerRouter()

	tests := []struct {
		decks       string
		deckType    string
		players     string
		status      int
		expectedMax int
	}{
		{"1", "standard", "2", 200, 2},
		{"2", "spanish21", "4", 200, 4},
		{"1", "standard", "6", 200, 6},
		{"1", "standard", "0", 400, 0},   // Invalid players
		{"1", "standard", "11", 400, 0},  // Too many players
		{"0", "standard", "4", 400, 0},   // Invalid decks
		{"abc", "standard", "4", 400, 0}, // Invalid decks
	}

	for _, test := range tests {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/game/new/"+test.decks+"/"+test.deckType+"/"+test.players, nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, test.status, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		if test.status == 200 {
			assert.Contains(t, response, "game_id")
			assert.Contains(t, response, "game_type")
			assert.Equal(t, "Blackjack", response["game_type"])
			assert.Equal(t, float64(test.expectedMax), response["max_players"])
			assert.Equal(t, float64(0), response["current_players"])
		}
	}
}

func TestGetGameStateEndpoint(t *testing.T) {
	router := setupPlayerRouter()

	// Create a game
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/game/new/1/standard/4", nil)
	router.ServeHTTP(w, req)

	var gameResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &gameResponse)
	gameID := gameResponse["game_id"].(string)

	// Get game state
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/game/"+gameID+"/state", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, gameID, response["game_id"])
	assert.Equal(t, "Blackjack", response["game_type"])
	assert.Equal(t, "Standard", response["deck_type"])
	assert.Equal(t, float64(52), response["remaining_cards"])
	assert.Equal(t, float64(4), response["max_players"])
	assert.Equal(t, float64(0), response["current_players"])
	assert.Contains(t, response, "players")
	assert.Contains(t, response, "dealer")
	assert.Contains(t, response, "discard_piles")

	// Check discard piles structure
	discardPiles := response["discard_piles"].([]interface{})
	assert.Equal(t, 1, len(discardPiles))
	mainPile := discardPiles[0].(map[string]interface{})
	assert.Equal(t, "main", mainPile["id"])
	assert.Equal(t, "Main Discard Pile", mainPile["name"])
	assert.Equal(t, float64(0), mainPile["size"])
}

func TestAddPlayerEndpoint(t *testing.T) {
	router := setupPlayerRouter()

	// Create a game
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/game/new/1/standard/2", nil)
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

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, gameID, response["game_id"])
	assert.Contains(t, response, "player")
	player := response["player"].(map[string]interface{})
	assert.Equal(t, "Alice", player["name"])
	assert.NotEmpty(t, player["id"])

	// Add another player
	playerData2 := map[string]string{"name": "Bob"}
	jsonData2, _ := json.Marshal(playerData2)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/game/"+gameID+"/players", bytes.NewBuffer(jsonData2))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	// Try to add third player (should fail - max is 2)
	playerData3 := map[string]string{"name": "Charlie"}
	jsonData3, _ := json.Marshal(playerData3)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/game/"+gameID+"/players", bytes.NewBuffer(jsonData3))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)

	var errorResponse map[string]string
	json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.Equal(t, "Game is full", errorResponse["error"])
}

func TestRemovePlayerEndpoint(t *testing.T) {
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

	// Remove player
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/game/"+gameID+"/players/"+playerID, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, gameID, response["game_id"])
	assert.Equal(t, playerID, response["player_id"])

	// Try to remove non-existent player
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/game/"+gameID+"/players/nonexistent", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
}