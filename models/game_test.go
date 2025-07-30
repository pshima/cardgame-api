package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewGame(t *testing.T) {
	game := NewGame(1)
	
	assert.NotEmpty(t, game.ID)
	assert.Equal(t, Blackjack, game.GameType)
	assert.Equal(t, GameWaiting, game.Status)
	assert.NotNil(t, game.Deck)
	assert.Equal(t, 52, game.Deck.RemainingCards())
	assert.Equal(t, 0, len(game.Players))
	assert.NotNil(t, game.Dealer)
	assert.Equal(t, "dealer", game.Dealer.ID)
	assert.Equal(t, "Dealer", game.Dealer.Name)
	assert.Equal(t, 6, game.MaxPlayers)
	assert.Equal(t, 0, game.CurrentPlayer)
	assert.NotNil(t, game.DiscardPiles)
	assert.Contains(t, game.DiscardPiles, "main")
	assert.False(t, game.Created.IsZero())
	assert.False(t, game.LastUsed.IsZero())
}

func TestNewCustomGame(t *testing.T) {
	game := NewCustomGame(2, Spanish21)
	
	assert.Equal(t, Blackjack, game.GameType)
	assert.Equal(t, Spanish21, game.Deck.DeckType)
	assert.Equal(t, 96, game.Deck.RemainingCards()) // 2 Spanish21 decks = 96 cards
}

func TestNewGameWithType(t *testing.T) {
	game := NewGameWithType(1, Standard, Cribbage, 2)
	
	assert.Equal(t, Cribbage, game.GameType)
	assert.Equal(t, Standard, game.Deck.DeckType)
	assert.Equal(t, 2, game.MaxPlayers)
	assert.Equal(t, 52, game.Deck.RemainingCards())
}

func TestGameUpdateLastUsed(t *testing.T) {
	game := NewGame(1)
	originalTime := game.LastUsed
	
	// Wait a small amount to ensure time difference
	time.Sleep(1 * time.Millisecond)
	
	game.UpdateLastUsed()
	assert.True(t, game.LastUsed.After(originalTime))
}

func TestGameAddPlayer(t *testing.T) {
	game := NewGame(1)
	
	// Add first player
	player1 := game.AddPlayer("Alice")
	assert.NotNil(t, player1)
	assert.Equal(t, "Alice", player1.Name)
	assert.NotEmpty(t, player1.ID)
	assert.Equal(t, 0, len(player1.Hand))
	assert.Equal(t, 1, len(game.Players))
	
	// Add second player
	player2 := game.AddPlayer("Bob")
	assert.NotNil(t, player2)
	assert.Equal(t, "Bob", player2.Name)
	assert.NotEqual(t, player1.ID, player2.ID)
	assert.Equal(t, 2, len(game.Players))
	
	// Try to add too many players (default max is 6)
	for i := 3; i <= 6; i++ {
		player := game.AddPlayer("Player" + string(rune(i+'0')))
		assert.NotNil(t, player)
	}
	assert.Equal(t, 6, len(game.Players))
	
	// Adding 7th player should fail
	player7 := game.AddPlayer("Player7")
	assert.Nil(t, player7)
	assert.Equal(t, 6, len(game.Players))
}

func TestGameGetPlayer(t *testing.T) {
	game := NewGame(1)
	
	// Test getting dealer
	dealer := game.GetPlayer("dealer")
	assert.NotNil(t, dealer)
	assert.Equal(t, "dealer", dealer.ID)
	assert.Equal(t, "Dealer", dealer.Name)
	
	// Add a player and get them
	player := game.AddPlayer("Alice")
	assert.NotNil(t, player)
	
	retrieved := game.GetPlayer(player.ID)
	assert.NotNil(t, retrieved)
	assert.Equal(t, player.ID, retrieved.ID)
	assert.Equal(t, "Alice", retrieved.Name)
	
	// Test getting non-existent player
	missing := game.GetPlayer("non-existent-id")
	assert.Nil(t, missing)
}

func TestGameRemovePlayer(t *testing.T) {
	game := NewGame(1)
	
	// Add players
	player1 := game.AddPlayer("Alice")
	player2 := game.AddPlayer("Bob")
	player3 := game.AddPlayer("Charlie")
	assert.Equal(t, 3, len(game.Players))
	
	// Remove middle player
	removed := game.RemovePlayer(player2.ID)
	assert.True(t, removed)
	assert.Equal(t, 2, len(game.Players))
	
	// Verify correct players remain
	assert.Equal(t, player1.ID, game.Players[0].ID)
	assert.Equal(t, player3.ID, game.Players[1].ID)
	
	// Try to remove non-existent player
	removed = game.RemovePlayer("non-existent-id")
	assert.False(t, removed)
	assert.Equal(t, 2, len(game.Players))
	
	// Remove first player
	removed = game.RemovePlayer(player1.ID)
	assert.True(t, removed)
	assert.Equal(t, 1, len(game.Players))
	assert.Equal(t, player3.ID, game.Players[0].ID)
}

func TestGameDealToPlayer(t *testing.T) {
	game := NewGame(1)
	
	// Add a player
	player := game.AddPlayer("Alice")
	assert.NotNil(t, player)
	
	// Deal card face up
	card := game.DealToPlayer(player.ID, true)
	assert.NotNil(t, card)
	assert.True(t, card.FaceUp)
	assert.Equal(t, 1, len(player.Hand))
	assert.Equal(t, card, player.Hand[0])
	assert.Equal(t, 51, game.Deck.RemainingCards())
	
	// Deal card face down
	card2 := game.DealToPlayer(player.ID, false)
	assert.NotNil(t, card2)
	assert.False(t, card2.FaceUp)
	assert.Equal(t, 2, len(player.Hand))
	assert.Equal(t, 50, game.Deck.RemainingCards())
	
	// Deal to dealer
	dealerCard := game.DealToPlayer("dealer", true)
	assert.NotNil(t, dealerCard)
	assert.Equal(t, 1, len(game.Dealer.Hand))
	assert.Equal(t, 49, game.Deck.RemainingCards())
	
	// Try to deal to non-existent player
	invalidCard := game.DealToPlayer("non-existent", true)
	assert.Nil(t, invalidCard)
	assert.Equal(t, 49, game.Deck.RemainingCards())
}

func TestGameDealToPlayerEmptyDeck(t *testing.T) {
	game := NewGame(1)
	player := game.AddPlayer("Alice")
	
	// Deal all cards
	for i := 0; i < 52; i++ {
		card := game.DealToPlayer(player.ID, true)
		assert.NotNil(t, card)
	}
	
	// Deck should be empty
	assert.True(t, game.Deck.IsEmpty())
	
	// Try to deal from empty deck
	card := game.DealToPlayer(player.ID, true)
	assert.Nil(t, card)
}

func TestGameAddDiscardPile(t *testing.T) {
	game := NewGame(1)
	
	// Game should start with main discard pile
	assert.Contains(t, game.DiscardPiles, "main")
	
	// Add new discard pile
	pile := game.AddDiscardPile("custom", "Custom Pile")
	assert.NotNil(t, pile)
	assert.Equal(t, "custom", pile.ID)
	assert.Equal(t, "Custom Pile", pile.Name)
	assert.Equal(t, 0, len(pile.Cards))
	assert.Contains(t, game.DiscardPiles, "custom")
	
	// Try to add pile with duplicate ID
	duplicate := game.AddDiscardPile("custom", "Duplicate")
	assert.Nil(t, duplicate)
	assert.Equal(t, "Custom Pile", game.DiscardPiles["custom"].Name) // Original should remain
}

func TestGameGetDiscardPile(t *testing.T) {
	game := NewGame(1)
	
	// Get existing main pile
	mainPile := game.GetDiscardPile("main")
	assert.NotNil(t, mainPile)
	assert.Equal(t, "main", mainPile.ID)
	assert.Equal(t, "Main Discard Pile", mainPile.Name)
	
	// Add custom pile and get it
	game.AddDiscardPile("custom", "Custom Pile")
	customPile := game.GetDiscardPile("custom")
	assert.NotNil(t, customPile)
	assert.Equal(t, "custom", customPile.ID)
	
	// Try to get non-existent pile
	missing := game.GetDiscardPile("non-existent")
	assert.Nil(t, missing)
}