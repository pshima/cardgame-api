package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDeck(t *testing.T) {
	deck := NewDeck()
	assert.Equal(t, 52, len(deck.Cards))
	assert.Equal(t, Standard, deck.DeckType)
	assert.NotEmpty(t, deck.Name)
	assert.Contains(t, deck.Name, " ")
}

func TestNewCustomDeck(t *testing.T) {
	tests := []struct {
		numDecks     int
		deckType     DeckType
		expectedCards int
	}{
		{1, Standard, 52},
		{2, Standard, 104},
		{1, Spanish21, 48}, // No 10s
		{2, Spanish21, 96}, // No 10s, 2 decks
		{6, Spanish21, 288}, // 6 decks for Spanish 21
		{0, Standard, 52}, // Should default to 1
		{-1, Spanish21, 48}, // Should default to 1
	}

	for _, test := range tests {
		deck := NewCustomDeck(test.numDecks, test.deckType)
		assert.Equal(t, test.expectedCards, len(deck.Cards))
		assert.Equal(t, test.expectedCards, deck.RemainingCards())
		assert.Equal(t, test.deckType, deck.DeckType)
		assert.NotEmpty(t, deck.Name)
		assert.Contains(t, deck.Name, " ")
	}
}

func TestSpanish21DeckComposition(t *testing.T) {
	deck := NewCustomDeck(1, Spanish21)
	
	// Should have 48 cards (52 - 4 tens)
	assert.Equal(t, 48, len(deck.Cards))
	
	// Count cards by rank
	rankCounts := make(map[Rank]int)
	for _, card := range deck.Cards {
		rankCounts[card.Rank]++
	}
	
	// Should have 4 of each rank except Ten
	for rank := Ace; rank <= King; rank++ {
		if rank == Ten {
			assert.Equal(t, 0, rankCounts[rank], "Spanish21 deck should have no 10s")
		} else {
			assert.Equal(t, 4, rankCounts[rank], "Should have 4 of each rank except 10")
		}
	}
	
	// Verify all suits are present for non-Ten cards
	suitCounts := make(map[Suit]int)
	for _, card := range deck.Cards {
		suitCounts[card.Suit]++
	}
	
	// Should have 12 cards of each suit (13 - 1 ten)
	for suit := Hearts; suit <= Spades; suit++ {
		assert.Equal(t, 12, suitCounts[suit])
	}
}

func TestDeckReset(t *testing.T) {
	deck := NewDeck()
	originalLength := len(deck.Cards)
	
	// Deal some cards
	deck.Deal()
	deck.Deal()
	assert.Equal(t, originalLength-2, len(deck.Cards))
	
	// Reset deck
	deck.Reset()
	assert.Equal(t, originalLength, len(deck.Cards))
	assert.Equal(t, originalLength, deck.RemainingCards())
	assert.False(t, deck.IsEmpty())
}

func TestDeckShuffle(t *testing.T) {
	deck := NewDeck()
	originalOrder := make([]Card, len(deck.Cards))
	copy(originalOrder, deck.Cards)
	
	deck.Shuffle()
	
	// After shuffling, the order should be different (very high probability)
	// We'll check if at least one card is in a different position
	different := false
	for i, card := range deck.Cards {
		if card.Rank != originalOrder[i].Rank || card.Suit != originalOrder[i].Suit {
			different = true
			break
		}
	}
	assert.True(t, different, "Deck should be shuffled")
	
	// Should still have same number of cards
	assert.Equal(t, len(originalOrder), len(deck.Cards))
}

func TestDeckDeal(t *testing.T) {
	deck := NewDeck()
	originalCount := deck.RemainingCards()
	
	card := deck.Deal()
	assert.NotNil(t, card)
	assert.Equal(t, originalCount-1, deck.RemainingCards())
	
	// Test that dealing reduces the deck
	for i := 0; i < 10; i++ {
		card := deck.Deal()
		assert.NotNil(t, card)
		assert.Equal(t, originalCount-i-2, deck.RemainingCards())
	}
}

func TestDeckDealEmpty(t *testing.T) {
	deck := NewDeck()
	
	// Deal all cards
	for i := 0; i < 52; i++ {
		card := deck.Deal()
		assert.NotNil(t, card)
	}
	
	// Now deck should be empty
	assert.True(t, deck.IsEmpty())
	assert.Equal(t, 0, deck.RemainingCards())
	
	// Dealing from empty deck should return nil
	card := deck.Deal()
	assert.Nil(t, card)
}

func TestDeckIsEmpty(t *testing.T) {
	deck := NewDeck()
	assert.False(t, deck.IsEmpty())
	
	// Deal all cards
	for len(deck.Cards) > 0 {
		deck.Deal()
	}
	
	assert.True(t, deck.IsEmpty())
}

func TestDeckRemainingCards(t *testing.T) {
	deck := NewDeck()
	assert.Equal(t, 52, deck.RemainingCards())
	
	deck.Deal()
	assert.Equal(t, 51, deck.RemainingCards())
	
	deck.Deal()
	deck.Deal()
	assert.Equal(t, 49, deck.RemainingCards())
}

func TestResetWithDecksAndType(t *testing.T) {
	// Start with a standard deck
	deck := NewCustomDeck(1, Standard)
	assert.Equal(t, 52, len(deck.Cards))
	assert.Equal(t, Standard, deck.DeckType)
	
	// Deal some cards
	deck.Deal()
	deck.Deal()
	assert.Equal(t, 50, deck.RemainingCards())
	
	// Reset to Spanish21
	deck.ResetWithDecksAndType(2, Spanish21)
	assert.Equal(t, 96, len(deck.Cards))
	assert.Equal(t, Spanish21, deck.DeckType)
	
	// Verify no 10s in the reset deck
	for _, card := range deck.Cards {
		assert.NotEqual(t, Ten, card.Rank, "Spanish21 deck should not contain 10s")
	}
}

func TestMultiDeckStandard(t *testing.T) {
	deck := NewCustomDeck(2, Standard)
	
	// Should have 104 cards (52 * 2)
	assert.Equal(t, 104, len(deck.Cards))
	
	// Count cards by rank
	rankCounts := make(map[Rank]int)
	for _, card := range deck.Cards {
		rankCounts[card.Rank]++
	}
	
	// Should have 8 of each rank (4 per deck * 2 decks)
	for rank := Ace; rank <= King; rank++ {
		assert.Equal(t, 8, rankCounts[rank])
	}
}

func TestDiscardPile(t *testing.T) {
	pile := &DiscardPile{
		ID:    "test",
		Name:  "Test Pile",
		Cards: []*Card{},
	}
	
	assert.Equal(t, 0, pile.Size())
	
	card := &Card{Rank: Ace, Suit: Hearts, FaceUp: true}
	pile.AddCard(card)
	
	assert.Equal(t, 1, pile.Size())
	assert.Equal(t, card, pile.Cards[0])
}