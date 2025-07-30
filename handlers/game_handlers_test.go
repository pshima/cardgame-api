package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/peteshima/cardgame-api/managers"
	"github.com/peteshima/cardgame-api/middleware"
)

func setupTestHandler() *HandlerDependencies {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()
	var metricsRegistry *middleware.MetricsRegistry = nil // Use nil to avoid panic in tests
	gameManager := managers.NewGameManager()
	customDeckManager := managers.NewCustomDeckManager()
	startTime := time.Now()

	return NewHandlerDependencies(
		logger,
		metricsRegistry,
		gameManager,
		customDeckManager,
		startTime,
	)
}

func TestCreateNewGame(t *testing.T) {
	deps := setupTestHandler()
	
	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/game/new", nil)

	// Call handler
	deps.CreateNewGame(c)

	// Check response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "game_id")
	assert.Contains(t, w.Body.String(), "Standard")
	assert.Contains(t, w.Body.String(), "New Standard game created")
}

func TestCreateNewGameWithDecks(t *testing.T) {
	deps := setupTestHandler()

	tests := []struct {
		name         string
		decks        string
		expectedCode int
		shouldContain string
	}{
		{
			name:         "valid single deck",
			decks:        "1",
			expectedCode: http.StatusOK,
			shouldContain: "game_id",
		},
		{
			name:         "valid multiple decks",
			decks:        "3",
			expectedCode: http.StatusOK,
			shouldContain: "game_id",
		},
		{
			name:         "invalid deck count - zero",
			decks:        "0",
			expectedCode: http.StatusBadRequest,
			shouldContain: "Invalid decks parameter",
		},
		{
			name:         "invalid deck count - too high",
			decks:        "101",
			expectedCode: http.StatusBadRequest,
			shouldContain: "Invalid decks parameter",
		},
		{
			name:         "invalid deck count - not a number",
			decks:        "abc",
			expectedCode: http.StatusBadRequest,
			shouldContain: "Invalid decks parameter",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest("GET", "/game/new/"+test.decks, nil)
			c.Params = gin.Params{gin.Param{Key: "decks", Value: test.decks}}

			deps.CreateNewGameWithDecks(c)

			assert.Equal(t, test.expectedCode, w.Code)
			assert.Contains(t, w.Body.String(), test.shouldContain)
		})
	}
}

func TestCreateNewGameWithType(t *testing.T) {
	deps := setupTestHandler()

	tests := []struct {
		name         string
		decks        string
		deckType     string
		expectedCode int
		shouldContain string
	}{
		{
			name:         "valid standard deck",
			decks:        "1",
			deckType:     "standard",
			expectedCode: http.StatusOK,
			shouldContain: "Standard",
		},
		{
			name:         "valid spanish21 deck",
			decks:        "2",
			deckType:     "spanish21",
			expectedCode: http.StatusOK,
			shouldContain: "Spanish21",
		},
		{
			name:         "invalid deck type",
			decks:        "1",
			deckType:     "invalid",
			expectedCode: http.StatusBadRequest,
			shouldContain: "Invalid deck type",
		},
		{
			name:         "invalid deck count",
			decks:        "0",
			deckType:     "standard",
			expectedCode: http.StatusBadRequest,
			shouldContain: "Invalid decks parameter",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest("GET", "/game/new/"+test.decks+"/"+test.deckType, nil)
			c.Params = gin.Params{
				gin.Param{Key: "decks", Value: test.decks},
				gin.Param{Key: "type", Value: test.deckType},
			}

			deps.CreateNewGameWithType(c)

			assert.Equal(t, test.expectedCode, w.Code)
			assert.Contains(t, w.Body.String(), test.shouldContain)
		})
	}
}

func TestCreateNewGameWithPlayers(t *testing.T) {
	deps := setupTestHandler()

	tests := []struct {
		name         string
		decks        string
		deckType     string
		players      string
		expectedCode int
		shouldContain string
	}{
		{
			name:         "valid all options",
			decks:        "2",
			deckType:     "standard", 
			players:      "4",
			expectedCode: http.StatusOK,
			shouldContain: "Standard",
		},
		{
			name:         "invalid players count",
			decks:        "1",
			deckType:     "standard",
			players:      "0",
			expectedCode: http.StatusBadRequest,
			shouldContain: "Invalid players parameter",
		},
		{
			name:         "too many players",
			decks:        "1",
			deckType:     "standard",
			players:      "11",
			expectedCode: http.StatusBadRequest,
			shouldContain: "Invalid players parameter",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest("GET", "/game/new/"+test.decks+"/"+test.deckType+"/"+test.players, nil)
			c.Params = gin.Params{
				gin.Param{Key: "decks", Value: test.decks},
				gin.Param{Key: "type", Value: test.deckType},
				gin.Param{Key: "players", Value: test.players},
			}

			deps.CreateNewGameWithPlayers(c)

			assert.Equal(t, test.expectedCode, w.Code)
			assert.Contains(t, w.Body.String(), test.shouldContain)
		})
	}
}

func TestGetGameState(t *testing.T) {
	deps := setupTestHandler()
	
	// First create a game to test with
	game := deps.GameService.CreateGame(1)
	
	// Test valid game ID
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/game/"+game.ID+"/state", nil)
	c.Params = gin.Params{gin.Param{Key: "gameId", Value: game.ID}}

	deps.GetGameState(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "game_id")
	assert.Contains(t, w.Body.String(), game.ID)

	// Test invalid game ID
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/game/invalid-id/state", nil)
	c.Params = gin.Params{gin.Param{Key: "gameId", Value: "invalid-id"}}

	deps.GetGameState(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Game not found")
}

func TestListGames(t *testing.T) {
	deps := setupTestHandler()

	// Test with no games
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/games", nil)

	deps.ListGames(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "games")
	assert.Contains(t, w.Body.String(), "\"count\":0")

	// Create some games and test again
	deps.GameService.CreateGame(1)
	deps.GameService.CreateGame(2)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/games", nil)

	deps.ListGames(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "\"count\":2")
}

func TestDeleteGame(t *testing.T) {
	deps := setupTestHandler()
	
	// Create a game to delete
	game := deps.GameService.CreateGame(1)
	
	// Test deleting existing game
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("DELETE", "/game/"+game.ID, nil)
	c.Params = gin.Params{gin.Param{Key: "gameId", Value: game.ID}}

	deps.DeleteGame(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Game deleted successfully")

	// Test deleting non-existent game
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("DELETE", "/game/invalid-id", nil)
	c.Params = gin.Params{gin.Param{Key: "gameId", Value: "invalid-id"}}

	deps.DeleteGame(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Game not found")
}