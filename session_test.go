package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupSessionRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	gameManager = NewGameManager()
	r := gin.Default()

	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, World!",
		})
	})

	r.GET("/deck-types", listDeckTypes)
	r.GET("/game/new", createNewGame)
	r.GET("/game/new/:decks", createNewGameWithDecks)
	r.GET("/game/new/:decks/:type", createNewGameWithType)
	r.GET("/game/:gameId/shuffle", shuffleDeck)
	r.GET("/game/:gameId", getGameInfo)
	r.GET("/game/:gameId/deal", dealCard)
	r.GET("/game/:gameId/deal/:count", dealCards)
	r.GET("/game/:gameId/reset", resetDeck)
	r.GET("/game/:gameId/reset/:decks", resetDeckWithDecks)
	r.GET("/game/:gameId/reset/:decks/:type", resetDeckWithType)
	r.DELETE("/game/:gameId", deleteGame)
	r.GET("/games", listGames)

	return r
}

func TestNewMultiDeck(t *testing.T) {
	tests := []struct {
		numDecks     int
		expectedCards int
	}{
		{1, 52},
		{2, 104},
		{6, 312},
		{8, 416},
		{0, 52}, // Should default to 1
		{-1, 52}, // Should default to 1
	}

	for _, test := range tests {
		deck := NewMultiDeck(test.numDecks)
		assert.Equal(t, test.expectedCards, len(deck.Cards))
		assert.Equal(t, test.expectedCards, deck.RemainingCards())
	}
}

func TestResetWithDecks(t *testing.T) {
	deck := NewDeck()
	deck.Deal() // Remove one card
	assert.Equal(t, 51, deck.RemainingCards())

	deck.ResetWithDecks(3)
	assert.Equal(t, 156, deck.RemainingCards())

	deck.ResetWithDecks(0) // Should default to 1
	assert.Equal(t, 52, deck.RemainingCards())
}

func TestGameManager(t *testing.T) {
	gm := NewGameManager()
	assert.Equal(t, 0, gm.GameCount())

	game1 := gm.CreateGame(1)
	assert.NotEmpty(t, game1.ID)
	assert.Equal(t, 52, game1.Deck.RemainingCards())
	assert.Equal(t, 1, gm.GameCount())

	game2 := gm.CreateGame(2)
	assert.NotEmpty(t, game2.ID)
	assert.NotEqual(t, game1.ID, game2.ID)
	assert.Equal(t, 104, game2.Deck.RemainingCards())
	assert.Equal(t, 2, gm.GameCount())

	retrievedGame, exists := gm.GetGame(game1.ID)
	assert.True(t, exists)
	assert.Equal(t, game1.ID, retrievedGame.ID)

	_, exists = gm.GetGame("nonexistent")
	assert.False(t, exists)

	deleted := gm.DeleteGame(game1.ID)
	assert.True(t, deleted)
	assert.Equal(t, 1, gm.GameCount())

	deleted = gm.DeleteGame("nonexistent")
	assert.False(t, deleted)

	gameIDs := gm.ListGames()
	assert.Equal(t, 1, len(gameIDs))
	assert.Equal(t, game2.ID, gameIDs[0])
}

func TestGameCleanup(t *testing.T) {
	gm := NewGameManager()
	
	game1 := gm.CreateGame(1)
	game2 := gm.CreateGame(1)
	
	// Simulate old game by setting last used time to 2 hours ago
	game1.LastUsed = time.Now().Add(-2 * time.Hour)
	
	deleted := gm.CleanupOldGames(1 * time.Hour)
	assert.Equal(t, 1, deleted)
	assert.Equal(t, 1, gm.GameCount())
	
	_, exists := gm.GetGame(game1.ID)
	assert.False(t, exists)
	
	_, exists = gm.GetGame(game2.ID)
	assert.True(t, exists)
}

func TestCreateNewGameEndpoint(t *testing.T) {
	router := setupSessionRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/game/new", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "game_id")
	assert.Equal(t, "New Standard game created with 1 deck", response["message"])
	assert.Equal(t, float64(52), response["remaining_cards"])
}

func TestCreateNewGameWithDecksEndpoint(t *testing.T) {
	router := setupSessionRouter()

	tests := []struct {
		decks    string
		expected int
		status   int
	}{
		{"2", 104, 200},
		{"6", 312, 200},
		{"8", 416, 200},
		{"0", 0, 400},
		{"-1", 0, 400},
		{"abc", 0, 400},
	}

	for _, test := range tests {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/game/new/"+test.decks, nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, test.status, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		
		if test.status == 200 {
			assert.Contains(t, response, "game_id")
			assert.Equal(t, float64(test.expected), response["remaining_cards"])
		} else {
			assert.Equal(t, "Invalid decks parameter", response["error"])
		}
	}
}

func TestGameNotFoundErrors(t *testing.T) {
	router := setupSessionRouter()
	nonexistentID := "nonexistent-game-id"

	endpoints := []string{
		"/game/" + nonexistentID + "/shuffle",
		"/game/" + nonexistentID,
		"/game/" + nonexistentID + "/deal",
		"/game/" + nonexistentID + "/deal/5",
		"/game/" + nonexistentID + "/reset",
		"/game/" + nonexistentID + "/reset/2",
	}

	for _, endpoint := range endpoints {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", endpoint, nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 404, w.Code)
		
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Game not found", response["error"])
	}
}

func TestCompleteSessionWorkflow(t *testing.T) {
	router := setupSessionRouter()

	// Create a new game with 2 decks
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/game/new/2", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	
	var gameResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &gameResponse)
	assert.NoError(t, err)
	gameID := gameResponse["game_id"].(string)
	assert.Equal(t, float64(104), gameResponse["remaining_cards"])

	// Shuffle the deck
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/game/"+gameID+"/shuffle", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	// Deal 5 cards
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/game/"+gameID+"/deal/5", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	
	var dealResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &dealResponse)
	assert.NoError(t, err)
	assert.Equal(t, gameID, dealResponse["game_id"])
	assert.Equal(t, float64(99), dealResponse["remaining_cards"])
	assert.Equal(t, float64(5), dealResponse["cards_dealt"])

	// Get game info
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/game/"+gameID, nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	
	var infoResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &infoResponse)
	assert.NoError(t, err)
	assert.Equal(t, gameID, infoResponse["game_id"])
	assert.Equal(t, float64(99), infoResponse["remaining_cards"])

	// Reset with 3 decks
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/game/"+gameID+"/reset/3", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	
	var resetResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &resetResponse)
	assert.NoError(t, err)
	assert.Equal(t, gameID, resetResponse["game_id"])
	assert.Equal(t, float64(156), resetResponse["remaining_cards"])

	// List games
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/games", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	
	var listResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &listResponse)
	assert.NoError(t, err)
	assert.Equal(t, float64(1), listResponse["game_count"])
	games := listResponse["games"].([]interface{})
	assert.Equal(t, 1, len(games))
	assert.Equal(t, gameID, games[0])

	// Delete game
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/game/"+gameID, nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	
	var deleteResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &deleteResponse)
	assert.NoError(t, err)
	assert.Equal(t, gameID, deleteResponse["game_id"])

	// Verify game is deleted
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/game/"+gameID, nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 404, w.Code)
}

func TestConcurrentGameAccess(t *testing.T) {
	router := setupSessionRouter()

	// Create multiple games
	var gameIDs []string
	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/game/new", nil)
		router.ServeHTTP(w, req)
		
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		gameIDs = append(gameIDs, response["game_id"].(string))
	}

	// Verify all games are independent
	for i, gameID := range gameIDs {
		// Deal different numbers of cards from each game
		dealCount := i + 1
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/game/"+gameID+"/deal/"+string(rune(dealCount+'0')), nil)
		router.ServeHTTP(w, req)
		
		var dealResponse map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &dealResponse)
		assert.Equal(t, float64(52-dealCount), dealResponse["remaining_cards"])
	}

	// Verify games are still independent
	for i, gameID := range gameIDs {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/game/"+gameID, nil)
		router.ServeHTTP(w, req)
		
		var infoResponse map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &infoResponse)
		expectedRemaining := 52 - (i + 1)
		assert.Equal(t, float64(expectedRemaining), infoResponse["remaining_cards"])
	}
}