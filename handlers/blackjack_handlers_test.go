package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/peteshima/cardgame-api/api"
	"github.com/peteshima/cardgame-api/models"
)

func TestStartBlackjackGame(t *testing.T) {
	deps := setupTestHandler()
	
	// Create a game and add players
	game := deps.GameService.CreateGame(1)
	deps.GameService.AddPlayerToGame(game.ID, "Alice")
	deps.GameService.AddPlayerToGame(game.ID, "Bob")

	// Test starting blackjack game
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/game/"+game.ID+"/start", nil)
	c.Params = gin.Params{gin.Param{Key: "gameId", Value: game.ID}}

	deps.StartBlackjackGame(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Blackjack game started")

	// Test invalid game ID
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/game/invalid-id/start", nil)
	c.Params = gin.Params{gin.Param{Key: "gameId", Value: "invalid-id"}}

	deps.StartBlackjackGame(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid game ID format")

	// Test non-existent game
	fakeUUID := "550e8400-e29b-41d4-a716-446655440000"
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/game/"+fakeUUID+"/start", nil)
	c.Params = gin.Params{gin.Param{Key: "gameId", Value: fakeUUID}}

	deps.StartBlackjackGame(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Game not found")
}

func TestPlayerHit(t *testing.T) {
	deps := setupTestHandler()
	
	// Create a game, add player, and start blackjack
	game := deps.GameService.CreateGame(1)
	resultGame, player, success := deps.GameService.AddPlayerToGame(game.ID, "Alice")
	assert.True(t, success)

	// Start the blackjack game first
	actualGame, _ := deps.GameService.GetGame(resultGame.ID)
	actualGame.StartBlackjackGame()

	// Test player hit
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/game/"+resultGame.ID+"/hit/"+player.ID, nil)
	c.Params = gin.Params{
		gin.Param{Key: "gameId", Value: resultGame.ID},
		gin.Param{Key: "playerId", Value: player.ID},
	}

	deps.PlayerHit(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Player hit")

	// Test invalid game ID
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/game/invalid-id/hit/"+player.ID, nil)
	c.Params = gin.Params{
		gin.Param{Key: "gameId", Value: "invalid-id"},
		gin.Param{Key: "playerId", Value: player.ID},
	}

	deps.PlayerHit(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid game ID format")

	// Test invalid player ID
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/game/"+resultGame.ID+"/hit/invalid-id", nil)
	c.Params = gin.Params{
		gin.Param{Key: "gameId", Value: resultGame.ID},
		gin.Param{Key: "playerId", Value: "invalid-id"},
	}

	deps.PlayerHit(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid player ID format")
}

func TestPlayerStand(t *testing.T) {
	deps := setupTestHandler()
	
	// Create a game, add player, and start blackjack
	game := deps.GameService.CreateGame(1)
	resultGame, player, success := deps.GameService.AddPlayerToGame(game.ID, "Alice")
	assert.True(t, success)

	// Start the blackjack game first
	actualGame, _ := deps.GameService.GetGame(resultGame.ID)
	actualGame.StartBlackjackGame()

	// Test player stand
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/game/"+resultGame.ID+"/stand/"+player.ID, nil)
	c.Params = gin.Params{
		gin.Param{Key: "gameId", Value: resultGame.ID},
		gin.Param{Key: "playerId", Value: player.ID},
	}

	deps.PlayerStand(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Player stands")

	// Test invalid game ID
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/game/invalid-id/stand/"+player.ID, nil)
	c.Params = gin.Params{
		gin.Param{Key: "gameId", Value: "invalid-id"},
		gin.Param{Key: "playerId", Value: player.ID},
	}

	deps.PlayerStand(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid game ID format")
}

func TestGetGameResults(t *testing.T) {
	deps := setupTestHandler()
	
	// Create a game, add player
	game := deps.GameService.CreateGame(1)
	resultGame, _, success := deps.GameService.AddPlayerToGame(game.ID, "Alice")
	assert.True(t, success)

	// Get the actual game instance and finish it
	actualGame, _ := deps.GameService.GetGame(resultGame.ID)
	actualGame.StartBlackjackGame()
	actualGame.Status = models.GameFinished // Set to finished to get results

	// Test get results
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/game/"+resultGame.ID+"/results", nil)
	c.Params = gin.Params{gin.Param{Key: "gameId", Value: resultGame.ID}}

	deps.GetGameResults(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "results")

	// Test invalid game ID
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/game/invalid-id/results", nil)
	c.Params = gin.Params{gin.Param{Key: "gameId", Value: "invalid-id"}}

	deps.GetGameResults(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid game ID format")

	// Test non-existent game
	fakeUUID := "550e8400-e29b-41d4-a716-446655440000"
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/game/"+fakeUUID+"/results", nil)
	c.Params = gin.Params{gin.Param{Key: "gameId", Value: fakeUUID}}

	deps.GetGameResults(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Game not found")
}

func TestAddPlayer(t *testing.T) {
	deps := setupTestHandler()
	
	// Create a game
	game := deps.GameService.CreateGame(1)

	// Test adding a player
	addPlayerRequest := api.AddPlayerRequest{
		Name: "Alice",
	}
	jsonData, _ := json.Marshal(addPlayerRequest)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/game/"+game.ID+"/players", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{gin.Param{Key: "gameId", Value: game.ID}}

	deps.AddPlayer(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Alice")
	assert.Contains(t, w.Body.String(), "player_id")

	// Test invalid JSON
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/game/"+game.ID+"/players", bytes.NewBuffer([]byte("invalid json")))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{gin.Param{Key: "gameId", Value: game.ID}}

	deps.AddPlayer(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid JSON")

	// Test invalid game ID
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/game/invalid-id/players", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{gin.Param{Key: "gameId", Value: "invalid-id"}}

	deps.AddPlayer(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid game ID format")
}

func TestRemovePlayer(t *testing.T) {
	deps := setupTestHandler()
	
	// Create a game and add a player
	game := deps.GameService.CreateGame(1)
	resultGame, player, success := deps.GameService.AddPlayerToGame(game.ID, "Alice")
	assert.True(t, success)

	// Test removing the player
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("DELETE", "/game/"+resultGame.ID+"/players/"+player.ID, nil)
	c.Params = gin.Params{
		gin.Param{Key: "gameId", Value: resultGame.ID},
		gin.Param{Key: "playerId", Value: player.ID},
	}

	deps.RemovePlayer(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Player removed successfully")

	// Test invalid game ID
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("DELETE", "/game/invalid-id/players/"+player.ID, nil)
	c.Params = gin.Params{
		gin.Param{Key: "gameId", Value: "invalid-id"},
		gin.Param{Key: "playerId", Value: player.ID},
	}

	deps.RemovePlayer(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid game ID format")

	// Test invalid player ID
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("DELETE", "/game/"+resultGame.ID+"/players/invalid-id", nil)
	c.Params = gin.Params{
		gin.Param{Key: "gameId", Value: resultGame.ID},
		gin.Param{Key: "playerId", Value: "invalid-id"},
	}

	deps.RemovePlayer(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid player ID format")
}