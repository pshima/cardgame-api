package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/peteshima/cardgame-api/managers"
	"github.com/peteshima/cardgame-api/models"
)

func TestNewGameService(t *testing.T) {
	gm := managers.NewGameManager()
	gs := NewGameService(gm)
	
	assert.NotNil(t, gs)
	assert.Equal(t, gm, gs.gameManager)
}

func TestGameServiceCreateGame(t *testing.T) {
	gm := managers.NewGameManager()
	gs := NewGameService(gm)
	
	game := gs.CreateGame(1)
	assert.NotNil(t, game)
	assert.Equal(t, models.Blackjack, game.GameType)
	assert.Equal(t, models.Standard, game.Deck.DeckType)
	assert.Equal(t, 6, game.MaxPlayers)
}

func TestGameServiceGetGame(t *testing.T) {
	gm := managers.NewGameManager()
	gs := NewGameService(gm)
	
	// Create a game first
	game := gs.CreateGame(1)
	assert.NotNil(t, game)
	
	// Retrieve it
	retrieved, exists := gs.GetGame(game.ID)
	assert.True(t, exists)
	assert.NotNil(t, retrieved)
	assert.Equal(t, game.ID, retrieved.ID)
	
	// Try non-existent game
	missing, exists := gs.GetGame("non-existent")
	assert.False(t, exists)
	assert.Nil(t, missing)
}

func TestGameServiceShuffleGameDeck(t *testing.T) {
	gm := managers.NewGameManager()
	gs := NewGameService(gm)
	
	game := gs.CreateGame(1)
	assert.NotNil(t, game)
	
	// Shuffle should return the game
	shuffled, exists := gs.ShuffleGameDeck(game.ID)
	assert.True(t, exists)
	assert.NotNil(t, shuffled)
	assert.Equal(t, game.ID, shuffled.ID)
	
	// Try non-existent game
	missing, exists := gs.ShuffleGameDeck("non-existent")
	assert.False(t, exists)
	assert.Nil(t, missing)
}

func TestGameServiceResetGameDeck(t *testing.T) {
	gm := managers.NewGameManager()
	gs := NewGameService(gm)
	
	game := gs.CreateGame(1)
	assert.NotNil(t, game)
	
	// Deal some cards first
	card := game.Deck.Deal()
	assert.NotNil(t, card)
	assert.Equal(t, 51, game.Deck.RemainingCards())
	
	// Reset should restore full deck
	reset, exists := gs.ResetGameDeck(game.ID)
	assert.True(t, exists)
	assert.NotNil(t, reset)
	assert.Equal(t, 52, reset.Deck.RemainingCards())
	
	// Try non-existent game
	missing, exists := gs.ResetGameDeck("non-existent")
	assert.False(t, exists)
	assert.Nil(t, missing)
}

func TestGameServiceListGames(t *testing.T) {
	gm := managers.NewGameManager()
	gs := NewGameService(gm)
	
	// Initially empty
	games := gs.ListGames()
	assert.Equal(t, 0, len(games))
	
	// Create some games
	game1 := gs.CreateGame(1)
	game2 := gs.CreateGameWithType(2, models.Spanish21)
	
	games = gs.ListGames()
	assert.Equal(t, 2, len(games))
	
	// Check games are included
	found := make(map[string]bool)
	for _, gameID := range games {
		found[gameID] = true
	}
	assert.True(t, found[game1.ID])
	assert.True(t, found[game2.ID])
}

func TestGameServiceCreateGameWithDecks(t *testing.T) {
	gm := managers.NewGameManager()
	gs := NewGameService(gm)
	
	game := gs.CreateGameWithDecks(2)
	assert.NotNil(t, game)
	assert.Equal(t, models.Standard, game.Deck.DeckType)
	assert.Equal(t, 104, game.Deck.RemainingCards()) // 2 decks
}

func TestGameServiceCreateGameWithAllOptions(t *testing.T) {
	gm := managers.NewGameManager()
	gs := NewGameService(gm)
	
	game := gs.CreateGameWithAllOptions(2, models.Spanish21, models.Cribbage, 4)
	assert.NotNil(t, game)
	assert.Equal(t, models.Spanish21, game.Deck.DeckType)
	assert.Equal(t, models.Cribbage, game.GameType)
	assert.Equal(t, 4, game.MaxPlayers)
	assert.Equal(t, 96, game.Deck.RemainingCards()) // 2 Spanish21 decks
}

func TestGameServiceDeleteGame(t *testing.T) {
	gm := managers.NewGameManager()
	gs := NewGameService(gm)
	
	// Create a game
	game := gs.CreateGame(1)
	assert.NotNil(t, game)
	
	// Delete it
	deleted := gs.DeleteGame(game.ID)
	assert.True(t, deleted)
	
	// Verify it's gone
	_, exists := gs.GetGame(game.ID)
	assert.False(t, exists)
	
	// Try to delete non-existent game
	deleted = gs.DeleteGame("non-existent")
	assert.False(t, deleted)
}

func TestGameServiceGetGameCount(t *testing.T) {
	gm := managers.NewGameManager()
	gs := NewGameService(gm)
	
	assert.Equal(t, 0, gs.GetGameCount())
	
	gs.CreateGame(1)
	assert.Equal(t, 1, gs.GetGameCount())
	
	game2 := gs.CreateGame(2)
	assert.Equal(t, 2, gs.GetGameCount())
	
	gs.DeleteGame(game2.ID)
	assert.Equal(t, 1, gs.GetGameCount())
}

func TestGameServiceResetGameDeckWithDecks(t *testing.T) {
	gm := managers.NewGameManager()
	gs := NewGameService(gm)
	
	game := gs.CreateGame(1)
	assert.NotNil(t, game)
	assert.Equal(t, 52, game.Deck.RemainingCards())
	
	// Reset with more decks
	reset, exists := gs.ResetGameDeckWithDecks(game.ID, 3)
	assert.True(t, exists)
	assert.NotNil(t, reset)
	assert.Equal(t, 156, reset.Deck.RemainingCards()) // 3 decks
	
	// Try non-existent game
	missing, exists := gs.ResetGameDeckWithDecks("non-existent", 2)
	assert.False(t, exists)
	assert.Nil(t, missing)
}

func TestGameServiceResetGameDeckWithType(t *testing.T) {
	gm := managers.NewGameManager()
	gs := NewGameService(gm)
	
	game := gs.CreateGame(1)
	assert.NotNil(t, game)
	assert.Equal(t, models.Standard, game.Deck.DeckType)
	
	// Reset with different deck type
	reset, exists := gs.ResetGameDeckWithType(game.ID, 2, models.Spanish21)
	assert.True(t, exists)
	assert.NotNil(t, reset)
	assert.Equal(t, models.Spanish21, reset.Deck.DeckType)
	assert.Equal(t, 96, reset.Deck.RemainingCards()) // 2 Spanish21 decks
	
	// Try non-existent game
	missing, exists := gs.ResetGameDeckWithType("non-existent", 1, models.Standard)
	assert.False(t, exists)
	assert.Nil(t, missing)
}

func TestGameServiceAddPlayerToGame(t *testing.T) {
	gm := managers.NewGameManager()
	gs := NewGameService(gm)
	
	game := gs.CreateGame(1)
	assert.NotNil(t, game)
	
	// Add a player
	resultGame, player, success := gs.AddPlayerToGame(game.ID, "Alice")
	assert.True(t, success)
	assert.NotNil(t, resultGame)
	assert.NotNil(t, player)
	assert.Equal(t, "Alice", player.Name)
	assert.Equal(t, 1, len(resultGame.Players))
	
	// Try non-existent game
	missingGame, missingPlayer, success := gs.AddPlayerToGame("non-existent", "Bob")
	assert.False(t, success)
	assert.Nil(t, missingGame)
	assert.Nil(t, missingPlayer)
}

func TestGameServiceRemovePlayerFromGame(t *testing.T) {
	gm := managers.NewGameManager()
	gs := NewGameService(gm)
	
	game := gs.CreateGame(1)
	assert.NotNil(t, game)
	
	// Add a player first
	resultGame, player, success := gs.AddPlayerToGame(game.ID, "Alice")
	assert.True(t, success)
	assert.Equal(t, 1, len(resultGame.Players))
	
	// Remove the player
	resultGame, removed := gs.RemovePlayerFromGame(game.ID, player.ID)
	assert.True(t, removed)
	assert.NotNil(t, resultGame)
	assert.Equal(t, 0, len(resultGame.Players))
	
	// Try non-existent game
	missingGame, removed := gs.RemovePlayerFromGame("non-existent", "some-player")
	assert.False(t, removed)
	assert.Nil(t, missingGame)
}

func TestGameServiceDealCard(t *testing.T) {
	gm := managers.NewGameManager()
	gs := NewGameService(gm)
	
	game := gs.CreateGame(1)
	assert.NotNil(t, game)
	initialCards := game.Deck.RemainingCards()
	
	// Deal a card
	resultGame, card, success := gs.DealCard(game.ID)
	assert.True(t, success)
	assert.NotNil(t, resultGame)
	assert.NotNil(t, card)
	assert.True(t, card.FaceUp) // Should default to face up
	assert.Equal(t, initialCards-1, resultGame.Deck.RemainingCards())
	
	// Try non-existent game
	missingGame, missingCard, success := gs.DealCard("non-existent")
	assert.False(t, success)
	assert.Nil(t, missingGame)
	assert.Nil(t, missingCard)
}

func TestGameServiceDealCards(t *testing.T) {
	gm := managers.NewGameManager()
	gs := NewGameService(gm)
	
	game := gs.CreateGame(1)
	assert.NotNil(t, game)
	initialCards := game.Deck.RemainingCards()
	
	// Deal multiple cards
	resultGame, cards, success := gs.DealCards(game.ID, 5)
	assert.True(t, success)
	assert.NotNil(t, resultGame)
	assert.Equal(t, 5, len(cards))
	assert.Equal(t, initialCards-5, resultGame.Deck.RemainingCards())
	
	// All cards should be face up
	for _, card := range cards {
		assert.True(t, card.FaceUp)
	}
	
	// Try to deal more cards than available
	remaining := resultGame.Deck.RemainingCards()
	resultGame, cards, success = gs.DealCards(game.ID, remaining+1)
	assert.False(t, success)
	assert.NotNil(t, resultGame)
	assert.Nil(t, cards)
	
	// Try non-existent game
	missingGame, missingCards, success := gs.DealCards("non-existent", 3)
	assert.False(t, success)
	assert.Nil(t, missingGame)
	assert.Nil(t, missingCards)
}

func TestGameServiceDealToPlayer(t *testing.T) {
	gm := managers.NewGameManager()
	gs := NewGameService(gm)
	
	game := gs.CreateGame(1)
	assert.NotNil(t, game)
	
	// Add a player first
	resultGame, player, success := gs.AddPlayerToGame(game.ID, "Alice")
	assert.True(t, success)
	initialCards := resultGame.Deck.RemainingCards()
	
	// Deal to player face down
	resultGame, resultPlayer, card, success := gs.DealToPlayer(game.ID, player.ID, false)
	assert.True(t, success)
	assert.NotNil(t, resultGame)
	assert.NotNil(t, resultPlayer)
	assert.NotNil(t, card)
	assert.False(t, card.FaceUp)
	assert.Equal(t, 1, len(resultPlayer.Hand))
	assert.Equal(t, initialCards-1, resultGame.Deck.RemainingCards())
	
	// Deal to player face up
	resultGame, resultPlayer, card, success = gs.DealToPlayer(game.ID, player.ID, true)
	assert.True(t, success)
	assert.True(t, card.FaceUp)
	assert.Equal(t, 2, len(resultPlayer.Hand))
	
	// Try non-existent game
	missingGame, missingPlayer, missingCard, success := gs.DealToPlayer("non-existent", player.ID, true)
	assert.False(t, success)
	assert.Nil(t, missingGame)
	assert.Nil(t, missingPlayer)
	assert.Nil(t, missingCard)
	
	// Try non-existent player
	resultGame, missingPlayer, missingCard, success = gs.DealToPlayer(game.ID, "non-existent-player", true)
	assert.False(t, success)
	assert.NotNil(t, resultGame)
	assert.Nil(t, missingPlayer)
	assert.Nil(t, missingCard)
}