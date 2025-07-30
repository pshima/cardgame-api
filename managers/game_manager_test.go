package managers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/peteshima/cardgame-api/models"
)

func TestNewGameManager(t *testing.T) {
	gm := NewGameManager()
	assert.NotNil(t, gm)
	assert.NotNil(t, gm.games)
	assert.Equal(t, 0, gm.GameCount())
}

func TestGameManagerCreateGame(t *testing.T) {
	gm := NewGameManager()
	
	game := gm.CreateGame(1)
	assert.NotNil(t, game)
	assert.NotEmpty(t, game.ID)
	assert.Equal(t, models.Blackjack, game.GameType)
	assert.Equal(t, models.Standard, game.Deck.DeckType)
	assert.Equal(t, 6, game.MaxPlayers)
	assert.Equal(t, 52, game.Deck.RemainingCards())
	assert.Equal(t, 1, gm.GameCount())
}

func TestGameManagerCreateGameWithType(t *testing.T) {
	gm := NewGameManager()
	
	game := gm.CreateGameWithType(2, models.Spanish21, models.Cribbage, 2)
	assert.NotNil(t, game)
	assert.Equal(t, models.Cribbage, game.GameType)
	assert.Equal(t, models.Spanish21, game.Deck.DeckType)
	assert.Equal(t, 2, game.MaxPlayers)
	assert.Equal(t, 96, game.Deck.RemainingCards()) // 2 Spanish21 decks
}

func TestGameManagerGetGame(t *testing.T) {
	gm := NewGameManager()
	
	// Create a game
	game := gm.CreateGame(1)
	assert.NotNil(t, game)
	
	// Retrieve the game
	retrieved, exists := gm.GetGame(game.ID)
	assert.True(t, exists)
	assert.NotNil(t, retrieved)
	assert.Equal(t, game.ID, retrieved.ID)
	
	// Try to get non-existent game
	missing, exists := gm.GetGame("non-existent-id")
	assert.False(t, exists)
	assert.Nil(t, missing)
}

func TestGameManagerDeleteGame(t *testing.T) {
	gm := NewGameManager()
	
	// Create multiple games
	game1 := gm.CreateGame(1)
	game2 := gm.CreateGame(2)
	assert.Equal(t, 2, gm.GameCount())
	
	// Delete first game
	deleted := gm.DeleteGame(game1.ID)
	assert.True(t, deleted)
	assert.Equal(t, 1, gm.GameCount())
	
	// Verify game is gone
	_, exists := gm.GetGame(game1.ID)
	assert.False(t, exists)
	
	// Verify second game still exists
	retrieved, exists := gm.GetGame(game2.ID)
	assert.True(t, exists)
	assert.Equal(t, game2.ID, retrieved.ID)
	
	// Try to delete non-existent game
	deleted = gm.DeleteGame("non-existent-id")
	assert.False(t, deleted)
	assert.Equal(t, 1, gm.GameCount())
}

func TestGameManagerListGames(t *testing.T) {
	gm := NewGameManager()
	
	// No games initially
	games := gm.ListGames()
	assert.Equal(t, 0, len(games))
	
	// Create some games
	game1 := gm.CreateGame(1)
	game2 := gm.CreateGame(2)
	
	// List games
	games = gm.ListGames()
	assert.Equal(t, 2, len(games))
	
	// Verify all games are included
	gameIDs := make(map[string]bool)
	for _, gameID := range games {
		gameIDs[gameID] = true
	}
	assert.True(t, gameIDs[game1.ID])
	assert.True(t, gameIDs[game2.ID])
}

func TestGameManagerGameCount(t *testing.T) {
	gm := NewGameManager()
	
	assert.Equal(t, 0, gm.GameCount())
	
	// Create games
	gm.CreateGame(1)
	assert.Equal(t, 1, gm.GameCount())
	
	game2 := gm.CreateGame(2)
	assert.Equal(t, 2, gm.GameCount())
	
	// Delete a game
	gm.DeleteGame(game2.ID)
	assert.Equal(t, 1, gm.GameCount())
}