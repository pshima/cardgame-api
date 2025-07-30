package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlayerAddCard(t *testing.T) {
	player := &Player{
		ID:   "test-player",
		Name: "Test Player",
		Hand: []*Card{},
	}
	
	card := &Card{Rank: Ace, Suit: Hearts, FaceUp: true}
	player.AddCard(card)
	
	assert.Equal(t, 1, len(player.Hand))
	assert.Equal(t, card, player.Hand[0])
}

func TestPlayerRemoveCard(t *testing.T) {
	player := &Player{
		ID:   "test-player",
		Name: "Test Player",
		Hand: []*Card{
			{Rank: Ace, Suit: Hearts, FaceUp: true},
			{Rank: King, Suit: Spades, FaceUp: true},
			{Rank: Queen, Suit: Diamonds, FaceUp: true},
		},
	}
	
	// Remove middle card
	removed := player.RemoveCard(1)
	assert.NotNil(t, removed)
	assert.Equal(t, King, removed.Rank)
	assert.Equal(t, Spades, removed.Suit)
	assert.Equal(t, 2, len(player.Hand))
	
	// Verify remaining cards
	assert.Equal(t, Ace, player.Hand[0].Rank)
	assert.Equal(t, Queen, player.Hand[1].Rank)
	
	// Test invalid index
	removed = player.RemoveCard(10)
	assert.Nil(t, removed)
	assert.Equal(t, 2, len(player.Hand))
	
	// Test negative index
	removed = player.RemoveCard(-1)
	assert.Nil(t, removed)
	assert.Equal(t, 2, len(player.Hand))
}

func TestPlayerHandSize(t *testing.T) {
	player := &Player{
		ID:   "test-player",
		Name: "Test Player",
		Hand: []*Card{},
	}
	
	assert.Equal(t, 0, player.HandSize())
	
	player.AddCard(&Card{Rank: Ace, Suit: Hearts})
	assert.Equal(t, 1, player.HandSize())
	
	player.AddCard(&Card{Rank: King, Suit: Spades})
	player.AddCard(&Card{Rank: Queen, Suit: Diamonds})
	assert.Equal(t, 3, player.HandSize())
}

func TestPlayerClearHand(t *testing.T) {
	player := &Player{
		ID:   "test-player",
		Name: "Test Player",
		Hand: []*Card{
			{Rank: Ace, Suit: Hearts},
			{Rank: King, Suit: Spades},
		},
	}
	
	assert.Equal(t, 2, player.HandSize())
	
	player.ClearHand()
	assert.Equal(t, 0, player.HandSize())
	assert.Equal(t, 0, len(player.Hand))
}

func TestPlayerBlackjackHandValue(t *testing.T) {
	player := &Player{
		ID:   "test-player",
		Name: "Test Player",
		Hand: []*Card{},
	}
	
	// Test empty hand
	value, hasBlackjack := player.BlackjackHandValue()
	assert.Equal(t, 0, value)
	assert.False(t, hasBlackjack)
	
	// Test blackjack (Ace + 10-value card)
	player.Hand = []*Card{
		{Rank: Ace, Suit: Hearts},
		{Rank: King, Suit: Spades},
	}
	value, hasBlackjack = player.BlackjackHandValue()
	assert.Equal(t, 21, value)
	assert.True(t, hasBlackjack)
	
	// Test 21 but not blackjack (more than 2 cards)
	player.Hand = []*Card{
		{Rank: Seven, Suit: Hearts},
		{Rank: Seven, Suit: Spades},
		{Rank: Seven, Suit: Diamonds},
	}
	value, hasBlackjack = player.BlackjackHandValue()
	assert.Equal(t, 21, value)
	assert.False(t, hasBlackjack)
	
	// Test soft ace (Ace counts as 1)
	player.Hand = []*Card{
		{Rank: Ace, Suit: Hearts},
		{Rank: Nine, Suit: Spades},
		{Rank: Five, Suit: Diamonds},
	}
	value, hasBlackjack = player.BlackjackHandValue()
	assert.Equal(t, 15, value) // Ace=1, 9, 5
	assert.False(t, hasBlackjack)
	
	// Test multiple aces
	player.Hand = []*Card{
		{Rank: Ace, Suit: Hearts},
		{Rank: Ace, Suit: Spades},
		{Rank: Nine, Suit: Diamonds},
	}
	value, hasBlackjack = player.BlackjackHandValue()
	assert.Equal(t, 21, value) // Ace=11, Ace=1, 9
	assert.False(t, hasBlackjack)
}

func TestPlayerIsBusted(t *testing.T) {
	player := &Player{
		ID:   "test-player",
		Name: "Test Player",
		Hand: []*Card{},
	}
	
	// Empty hand - not busted
	assert.False(t, player.IsBusted())
	
	// Under 21 - not busted
	player.Hand = []*Card{
		{Rank: Ten, Suit: Hearts},
		{Rank: Nine, Suit: Spades},
	}
	assert.False(t, player.IsBusted())
	
	// Exactly 21 - not busted
	player.Hand = []*Card{
		{Rank: Ace, Suit: Hearts},
		{Rank: King, Suit: Spades},
	}
	assert.False(t, player.IsBusted())
	
	// Over 21 - busted
	player.Hand = []*Card{
		{Rank: King, Suit: Hearts},
		{Rank: Queen, Suit: Spades},
		{Rank: Five, Suit: Diamonds},
	}
	assert.True(t, player.IsBusted())
}

func TestPlayerHasBlackjack(t *testing.T) {
	player := &Player{
		ID:   "test-player",
		Name: "Test Player",
		Hand: []*Card{},
	}
	
	// Empty hand
	assert.False(t, player.HasBlackjack())
	
	// Blackjack
	player.Hand = []*Card{
		{Rank: Ace, Suit: Hearts},
		{Rank: Jack, Suit: Spades},
	}
	assert.True(t, player.HasBlackjack())
	
	// 21 but not blackjack
	player.Hand = []*Card{
		{Rank: Seven, Suit: Hearts},
		{Rank: Seven, Suit: Spades},
		{Rank: Seven, Suit: Diamonds},
	}
	assert.False(t, player.HasBlackjack())
	
	// Less than 21
	player.Hand = []*Card{
		{Rank: Ten, Suit: Hearts},
		{Rank: Nine, Suit: Spades},
	}
	assert.False(t, player.HasBlackjack())
}

func TestPlayerScoreCribbageHand(t *testing.T) {
	player := &Player{
		ID:   "test-player",
		Name: "Test Player",
		Hand: []*Card{},
	}
	
	// Empty hand
	starter := &Card{Rank: Five, Suit: Hearts}
	score := player.ScoreCribbageHand(starter)
	assert.Equal(t, 0, score)
	
	// Test a simple fifteen (5 + 10)
	player.Hand = []*Card{
		{Rank: Five, Suit: Spades},
		{Rank: Jack, Suit: Hearts},
		{Rank: Two, Suit: Diamonds},
		{Rank: Three, Suit: Clubs},
	}
	starter = &Card{Rank: Four, Suit: Hearts}
	score = player.ScoreCribbageHand(starter)
	assert.Greater(t, score, 0) // Should score points for fifteen and run
	
	// Test with nil starter
	score = player.ScoreCribbageHand(nil)
	assert.Greater(t, score, 0) // Should still score hand alone
}