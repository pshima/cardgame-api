package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartBlackjackGame(t *testing.T) {
	game := NewGame(1)
	
	// Add some players
	player1 := game.AddPlayer("Alice")
	player2 := game.AddPlayer("Bob")
	assert.NotNil(t, player1)
	assert.NotNil(t, player2)
	
	// Start blackjack game
	err := game.StartBlackjackGame()
	assert.NoError(t, err)
	
	// Each player should have 2 cards
	assert.Equal(t, 2, len(player1.Hand))
	assert.Equal(t, 2, len(player2.Hand))
	
	// Dealer should have 2 cards (first face down, second face up)
	assert.NotNil(t, game.Dealer)
	assert.Equal(t, 2, len(game.Dealer.Hand))
	assert.False(t, game.Dealer.Hand[0].FaceUp) // Dealer hole card
	assert.True(t, game.Dealer.Hand[1].FaceUp)  // Dealer up card
	
	// Game status should be in progress
	assert.Equal(t, GameInProgress, game.Status)
	
	// Deck should have fewer cards
	assert.Equal(t, 46, game.Deck.RemainingCards()) // 52 - 6 cards dealt
	
	// Test with no players
	emptyGame := NewGame(1)
	err = emptyGame.StartBlackjackGame()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no players")
}

func TestPlayerHit(t *testing.T) {
	game := NewGame(1)
	player := game.AddPlayer("Alice")
	assert.NotNil(t, player)
	
	// Start game to set status
	game.Status = GameInProgress
	
	// Give player initial cards
	player.AddCard(&Card{Rank: Five, Suit: Hearts, FaceUp: true})
	player.AddCard(&Card{Rank: Six, Suit: Spades, FaceUp: true})
	initialHandSize := len(player.Hand)
	initialDeckSize := game.Deck.RemainingCards()
	
	// Player hits
	err := game.PlayerHit(player.ID)
	assert.NoError(t, err)
	assert.Equal(t, initialHandSize+1, len(player.Hand))
	assert.Equal(t, initialDeckSize-1, game.Deck.RemainingCards())
	
	// Test invalid player
	err = game.PlayerHit("non-existent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "player not found")
	
	// Test game not in progress
	game.Status = GameFinished
	err = game.PlayerHit(player.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "game is not in progress")
}

func TestPlayerStand(t *testing.T) {
	game := NewGame(1)
	player := game.AddPlayer("Alice")
	assert.NotNil(t, player)
	
	// Start game to set status and current player
	game.Status = GameInProgress
	game.CurrentPlayer = 0
	
	// Give player some cards
	player.AddCard(&Card{Rank: King, Suit: Hearts, FaceUp: true})
	player.AddCard(&Card{Rank: Nine, Suit: Spades, FaceUp: true})
	
	// Player stands
	err := game.PlayerStand(player.ID)
	assert.NoError(t, err)
	
	// Current player should advance
	assert.Equal(t, 1, game.CurrentPlayer)
	
	// Test invalid player
	err = game.PlayerStand("non-existent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "player not found")
	
	// Test game not in progress
	game.Status = GameFinished
	err = game.PlayerStand(player.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "game is not in progress")
}

func TestPlayDealer(t *testing.T) {
	game := NewGame(1)
	
	// Set up dealer with initial cards
	game.Dealer = &Player{
		ID:   "dealer",
		Name: "Dealer",
		Hand: []*Card{
			{Rank: Ten, Suit: Hearts, FaceUp: false}, // Hole card
			{Rank: Six, Suit: Spades, FaceUp: true},   // Up card
		},
	}
	initialDeckSize := game.Deck.RemainingCards()
	
	// Play dealer (should hit on 16, stand on 17)
	err := game.PlayDealer()
	assert.NoError(t, err)
	
	// Dealer should have revealed hole card
	assert.True(t, game.Dealer.Hand[0].FaceUp)
	
	// Dealer should have hit (16 is less than 17)
	assert.True(t, len(game.Dealer.Hand) > 2)
	assert.True(t, game.Deck.RemainingCards() < initialDeckSize)
	
	// Final hand value should be 17 or higher (unless busted)
	value, _ := game.Dealer.BlackjackHandValue()
	assert.True(t, value >= 17 || value > 21) // Either stands on 17+ or busted
	
	// Game should be finished
	assert.Equal(t, GameFinished, game.Status)
}

func TestGetGameResult(t *testing.T) {
	game := NewGame(1)
	player := game.AddPlayer("Alice")
	
	// Set up dealer
	game.Dealer = &Player{
		ID:   "dealer",
		Name: "Dealer",
	}
	
	// Game must be finished to get results
	game.Status = GameFinished
	
	// Test blackjack win
	player.Hand = []*Card{
		{Rank: Ace, Suit: Hearts, FaceUp: true},
		{Rank: King, Suit: Spades, FaceUp: true},
	}
	game.Dealer.Hand = []*Card{
		{Rank: Ten, Suit: Hearts, FaceUp: true},
		{Rank: Nine, Suit: Spades, FaceUp: true},
	}
	
	results := game.GetGameResult()
	assert.Equal(t, "blackjack", results[player.ID])
	
	// Test regular win
	player.Hand = []*Card{
		{Rank: Ten, Suit: Hearts, FaceUp: true},
		{Rank: Nine, Suit: Spades, FaceUp: true},
	}
	game.Dealer.Hand = []*Card{
		{Rank: Ten, Suit: Hearts, FaceUp: true},
		{Rank: Eight, Suit: Spades, FaceUp: true},
	}
	
	results = game.GetGameResult()
	assert.Equal(t, "win", results[player.ID])
	
	// Test push (tie)
	game.Dealer.Hand = []*Card{
		{Rank: Ten, Suit: Hearts, FaceUp: true},
		{Rank: Nine, Suit: Clubs, FaceUp: true},
	}
	
	results = game.GetGameResult()
	assert.Equal(t, "push", results[player.ID])
	
	// Test player bust
	player.Hand = []*Card{
		{Rank: Ten, Suit: Hearts, FaceUp: true},
		{Rank: Nine, Suit: Spades, FaceUp: true},
		{Rank: Five, Suit: Clubs, FaceUp: true},
	}
	
	results = game.GetGameResult()
	assert.Equal(t, "bust", results[player.ID])
	
	// Test dealer bust, player doesn't bust
	player.Hand = []*Card{
		{Rank: Ten, Suit: Hearts, FaceUp: true},
		{Rank: Nine, Suit: Spades, FaceUp: true},
	}
	game.Dealer.Hand = []*Card{
		{Rank: Ten, Suit: Hearts, FaceUp: true},
		{Rank: Nine, Suit: Spades, FaceUp: true},
		{Rank: Five, Suit: Clubs, FaceUp: true},
	}
	
	results = game.GetGameResult()
	assert.Equal(t, "win", results[player.ID])
	
	// Test lose
	player.Hand = []*Card{
		{Rank: Ten, Suit: Hearts, FaceUp: true},
		{Rank: Eight, Suit: Spades, FaceUp: true},
	}
	game.Dealer.Hand = []*Card{
		{Rank: Ten, Suit: Hearts, FaceUp: true},
		{Rank: Nine, Suit: Spades, FaceUp: true},
	}
	
	results = game.GetGameResult()
	assert.Equal(t, "lose", results[player.ID])
	
	// Test game not finished
	game.Status = GameInProgress
	results = game.GetGameResult()
	assert.Equal(t, "game not finished", results["status"])
}