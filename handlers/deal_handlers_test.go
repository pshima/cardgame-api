package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestDealCard(t *testing.T) {
	deps := setupTestHandler()
	
	// Create a game to test with
	game := deps.GameService.CreateGame(1)
	
	// Test valid game ID
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/game/"+game.ID+"/deal", nil)
	c.Params = gin.Params{gin.Param{Key: "gameId", Value: game.ID}}

	deps.DealCard(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "card")
	assert.Contains(t, w.Body.String(), "remaining_cards")

	// Test invalid game ID format
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/game/invalid-id/deal", nil)
	c.Params = gin.Params{gin.Param{Key: "gameId", Value: "invalid-id"}}

	deps.DealCard(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid game ID format")

	// Test non-existent game
	fakeUUID := "550e8400-e29b-41d4-a716-446655440000"
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/game/"+fakeUUID+"/deal", nil)
	c.Params = gin.Params{gin.Param{Key: "gameId", Value: fakeUUID}}

	deps.DealCard(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Game not found")
}

func TestDealCards(t *testing.T) {
	deps := setupTestHandler()
	
	// Create a game to test with
	game := deps.GameService.CreateGame(1)

	tests := []struct {
		name         string
		gameId       string
		count        string
		expectedCode int
		shouldContain string
	}{
		{
			name:         "valid deal 3 cards",
			gameId:       game.ID,
			count:        "3",
			expectedCode: http.StatusOK,
			shouldContain: "cards",
		},
		{
			name:         "invalid count - zero",
			gameId:       game.ID,
			count:        "0",
			expectedCode: http.StatusBadRequest,
			shouldContain: "Invalid count parameter",
		},
		{
			name:         "invalid count - too high",
			gameId:       game.ID,
			count:        "101",
			expectedCode: http.StatusBadRequest,
			shouldContain: "Invalid count parameter",
		},
		{
			name:         "invalid count - not a number",
			gameId:       game.ID,
			count:        "abc",
			expectedCode: http.StatusBadRequest,
			shouldContain: "Invalid count parameter",
		},
		{
			name:         "invalid game ID",
			gameId:       "invalid-id",
			count:        "5",
			expectedCode: http.StatusBadRequest,
			shouldContain: "Invalid game ID format",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest("GET", "/game/"+test.gameId+"/deal/"+test.count, nil)
			c.Params = gin.Params{
				gin.Param{Key: "gameId", Value: test.gameId},
				gin.Param{Key: "count", Value: test.count},
			}

			deps.DealCards(c)

			assert.Equal(t, test.expectedCode, w.Code)
			assert.Contains(t, w.Body.String(), test.shouldContain)
		})
	}
}

func TestDealToPlayer(t *testing.T) {
	deps := setupTestHandler()
	
	// Create a game and add a player
	game := deps.GameService.CreateGame(1)
	resultGame, player, success := deps.GameService.AddPlayerToGame(game.ID, "Alice")
	assert.True(t, success)

	tests := []struct {
		name         string
		gameId       string
		playerId     string
		expectedCode int
		shouldContain string
	}{
		{
			name:         "valid deal to player",
			gameId:       resultGame.ID,
			playerId:     player.ID,
			expectedCode: http.StatusOK,
			shouldContain: "card",
		},
		{
			name:         "deal to dealer",
			gameId:       resultGame.ID,
			playerId:     "dealer",
			expectedCode: http.StatusOK,
			shouldContain: "card",
		},
		{
			name:         "invalid game ID",
			gameId:       "invalid-id",
			playerId:     player.ID,
			expectedCode: http.StatusBadRequest,
			shouldContain: "Invalid game ID format",
		},
		{
			name:         "invalid player ID",
			gameId:       resultGame.ID,
			playerId:     "invalid-id",
			expectedCode: http.StatusBadRequest,
			shouldContain: "Invalid player ID format",
		},
		{
			name:         "non-existent game",
			gameId:       "550e8400-e29b-41d4-a716-446655440000",
			playerId:     player.ID,
			expectedCode: http.StatusNotFound,
			shouldContain: "Game not found",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest("GET", "/game/"+test.gameId+"/deal/player/"+test.playerId, nil)
			c.Params = gin.Params{
				gin.Param{Key: "gameId", Value: test.gameId},
				gin.Param{Key: "playerId", Value: test.playerId},
			}

			deps.DealToPlayer(c)

			assert.Equal(t, test.expectedCode, w.Code)
			assert.Contains(t, w.Body.String(), test.shouldContain)
		})
	}
}

func TestDealToPlayerWithFaceControl(t *testing.T) {
	deps := setupTestHandler()
	
	// Create a game and add a player
	game := deps.GameService.CreateGame(1)
	resultGame, player, success := deps.GameService.AddPlayerToGame(game.ID, "Alice")
	assert.True(t, success)

	tests := []struct {
		name         string
		gameId       string
		playerId     string
		faceUp       string
		expectedCode int
		shouldContain string
	}{
		{
			name:         "deal face up",
			gameId:       resultGame.ID,
			playerId:     player.ID,
			faceUp:       "true",
			expectedCode: http.StatusOK,
			shouldContain: "card",
		},
		{
			name:         "deal face down",
			gameId:       resultGame.ID,
			playerId:     player.ID,
			faceUp:       "false",
			expectedCode: http.StatusOK,
			shouldContain: "card",
		},
		{
			name:         "invalid face parameter",
			gameId:       resultGame.ID,
			playerId:     player.ID,
			faceUp:       "invalid",
			expectedCode: http.StatusBadRequest,
			shouldContain: "Invalid faceUp parameter",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest("GET", "/game/"+test.gameId+"/deal/player/"+test.playerId+"/"+test.faceUp, nil)
			c.Params = gin.Params{
				gin.Param{Key: "gameId", Value: test.gameId},
				gin.Param{Key: "playerId", Value: test.playerId},
				gin.Param{Key: "faceUp", Value: test.faceUp},
			}

			deps.DealToPlayerFaceUp(c)

			assert.Equal(t, test.expectedCode, w.Code)
			assert.Contains(t, w.Body.String(), test.shouldContain)
		})
	}
}