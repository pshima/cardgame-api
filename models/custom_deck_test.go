package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCustomDeckTemplate(t *testing.T) {
	deck := NewCustomDeckTemplate("Test Deck")
	
	assert.NotNil(t, deck)
	assert.NotEmpty(t, deck.ID)
	assert.Equal(t, "Test Deck", deck.Name)
	assert.Equal(t, 0, len(deck.Cards))
	assert.Equal(t, 0, deck.NextIndex)
	assert.False(t, deck.Created.IsZero())
	assert.False(t, deck.LastUsed.IsZero())
}

func TestCustomDeckUpdateLastUsed(t *testing.T) {
	deck := NewCustomDeckTemplate("Test Deck")
	originalTime := deck.LastUsed
	
	// Wait a bit to ensure time difference
	deck.UpdateLastUsed()
	
	assert.True(t, deck.LastUsed.After(originalTime) || deck.LastUsed.Equal(originalTime))
}

func TestCustomCardUpdateGameCompatibility(t *testing.T) {
	tests := []struct {
		name     string
		rank     interface{}
		suit     string
		expected bool
	}{
		{"numeric rank with suit", 5, "hearts", true},
		{"int32 rank with suit", int32(10), "spades", true},
		{"int64 rank with suit", int64(1), "diamonds", true},
		{"float32 rank with suit", float32(7.0), "clubs", true},
		{"float64 rank with suit", float64(13.0), "hearts", true},
		{"string rank with suit", "Ace", "hearts", false},
		{"numeric rank no suit", 5, "", false},
		{"no rank with suit", nil, "hearts", false},
		{"no rank no suit", nil, "", false},
		{"invalid type", []int{1, 2}, "hearts", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			card := &CustomCard{
				Name: test.name,
				Rank: test.rank,
				Suit: test.suit,
			}
			
			card.UpdateGameCompatibility()
			assert.Equal(t, test.expected, card.GameCompatible, test.name)
		})
	}
}

func TestCustomCardGetNumericRank(t *testing.T) {
	tests := []struct {
		name         string
		rank         interface{}
		expectedRank int
		expectedOk   bool
	}{
		{"int rank", 5, 5, true},
		{"int32 rank", int32(10), 10, true},
		{"int64 rank", int64(1), 1, true},
		{"float32 rank", float32(7.0), 7, true},
		{"float64 rank", float64(13.0), 13, true},
		{"string rank", "Ace", 0, false},
		{"nil rank", nil, 0, false},
		{"slice rank", []int{1, 2}, 0, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			card := &CustomCard{
				Name: test.name,
				Rank: test.rank,
			}
			
			rank, ok := card.GetNumericRank()
			assert.Equal(t, test.expectedOk, ok, test.name)
			if ok {
				assert.Equal(t, test.expectedRank, rank, test.name)
			}
		})
	}
}

func TestCustomCardAttributes(t *testing.T) {
	card := &CustomCard{
		Name: "Test Card",
		Rank: 5,
		Suit: "hearts",
		Attributes: map[string]string{
			"power":       "high",
			"element":     "fire",
			"description": "A powerful fire card",
		},
		GameCompatible: true,
	}

	// Test attributes are properly stored
	assert.Equal(t, "high", card.Attributes["power"])
	assert.Equal(t, "fire", card.Attributes["element"])
	assert.Equal(t, "A powerful fire card", card.Attributes["description"])
	
	// Test non-existent attribute
	assert.Equal(t, "", card.Attributes["non-existent"])
	
	// Test attribute modification
	card.Attributes["power"] = "maximum"
	assert.Equal(t, "maximum", card.Attributes["power"])
}

func TestCustomCardIndex(t *testing.T) {
	card := &CustomCard{
		Index: 42,
		Name:  "Test Card",
	}
	
	assert.Equal(t, 42, card.Index)
	
	// Test index modification
	card.Index = 100
	assert.Equal(t, 100, card.Index)
}

func TestCustomCardDeleted(t *testing.T) {
	card := &CustomCard{
		Name:    "Test Card",
		Deleted: false,
	}
	
	assert.False(t, card.Deleted)
	
	// Test tombstone deletion
	card.Deleted = true
	assert.True(t, card.Deleted)
}

func TestCustomDeckCardManagement(t *testing.T) {
	deck := NewCustomDeckTemplate("Test Deck")
	
	// Add some cards manually to test deck structure
	card1 := &CustomCard{
		Index: 0,
		Name:  "First Card",
		Rank:  1,
		Suit:  "hearts",
	}
	card1.UpdateGameCompatibility()
	
	card2 := &CustomCard{
		Index: 1,
		Name:  "Second Card",
		Rank:  "King",
		Suit:  "spades",
	}
	card2.UpdateGameCompatibility()
	
	deck.Cards = append(deck.Cards, card1, card2)
	deck.NextIndex = 2
	
	// Verify cards are properly stored
	assert.Equal(t, 2, len(deck.Cards))
	assert.Equal(t, card1, deck.Cards[0])
	assert.Equal(t, card2, deck.Cards[1])
	
	// Verify game compatibility
	assert.True(t, card1.GameCompatible)   // Numeric rank
	assert.False(t, card2.GameCompatible)  // String rank
	
	// Test NextIndex
	assert.Equal(t, 2, deck.NextIndex)
}