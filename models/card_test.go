package models

import (
	"testing"
	"strings"

	"github.com/stretchr/testify/assert"
)

func TestSuitString(t *testing.T) {
	tests := []struct {
		suit     Suit
		expected string
	}{
		{Hearts, "Hearts"},
		{Diamonds, "Diamonds"},
		{Clubs, "Clubs"},
		{Spades, "Spades"},
		{Suit(99), "Unknown"}, // Default case
	}

	for _, test := range tests {
		result := test.suit.String()
		assert.Equal(t, test.expected, result)
	}
}

func TestRankString(t *testing.T) {
	tests := []struct {
		rank     Rank
		expected string
	}{
		{Ace, "Ace"},
		{Two, "2"},
		{Three, "3"},
		{Four, "4"},
		{Five, "5"},
		{Six, "6"},
		{Seven, "7"},
		{Eight, "8"},
		{Nine, "9"},
		{Ten, "10"},
		{Jack, "Jack"},
		{Queen, "Queen"},
		{King, "King"},
		{Rank(99), "99"}, // Default case
	}

	for _, test := range tests {
		result := test.rank.String()
		assert.Equal(t, test.expected, result)
	}
}

func TestCardString(t *testing.T) {
	card := Card{Rank: Ace, Suit: Hearts, FaceUp: true}
	expected := "Ace of Hearts"
	assert.Equal(t, expected, card.String())

	cardTwo := Card{Rank: King, Suit: Spades, FaceUp: false}
	expected = "King of Spades"
	assert.Equal(t, expected, cardTwo.String())
}

func TestCardValue(t *testing.T) {
	tests := []struct {
		card     Card
		expected int
	}{
		{Card{Rank: Ace, Suit: Hearts}, 1},
		{Card{Rank: Two, Suit: Hearts}, 2},
		{Card{Rank: Ten, Suit: Hearts}, 10},
		{Card{Rank: Jack, Suit: Hearts}, 11},
		{Card{Rank: Queen, Suit: Hearts}, 12},
		{Card{Rank: King, Suit: Hearts}, 13},
	}

	for _, test := range tests {
		value := test.card.Value()
		assert.Equal(t, test.expected, value)
	}
}

func TestCardBlackjackValue(t *testing.T) {
	tests := []struct {
		card     Card
		expected int
	}{
		{Card{Rank: Ace, Suit: Hearts}, 11},
		{Card{Rank: Two, Suit: Hearts}, 2},
		{Card{Rank: Ten, Suit: Hearts}, 10},
		{Card{Rank: Jack, Suit: Hearts}, 10},
		{Card{Rank: Queen, Suit: Hearts}, 10},
		{Card{Rank: King, Suit: Hearts}, 10},
	}

	for _, test := range tests {
		result := test.card.BlackjackValue()
		assert.Equal(t, test.expected, result)
	}
}

func TestCardCribbageValue(t *testing.T) {
	tests := []struct {
		card     Card
		expected int
	}{
		{Card{Rank: Ace, Suit: Hearts}, 1},
		{Card{Rank: Two, Suit: Hearts}, 2},
		{Card{Rank: Ten, Suit: Hearts}, 10},
		{Card{Rank: Jack, Suit: Hearts}, 10},
		{Card{Rank: Queen, Suit: Hearts}, 10},
		{Card{Rank: King, Suit: Hearts}, 10},
	}

	for _, test := range tests {
		result := test.card.CribbageValue()
		assert.Equal(t, test.expected, result)
	}
}

func TestCribbagePlayValue(t *testing.T) {
	tests := []struct {
		card     Card
		expected int
	}{
		{Card{Rank: Ace, Suit: Hearts}, 1},
		{Card{Rank: Two, Suit: Hearts}, 2},
		{Card{Rank: Ten, Suit: Hearts}, 10},
		{Card{Rank: Jack, Suit: Hearts}, 10},
		{Card{Rank: Queen, Suit: Hearts}, 10},
		{Card{Rank: King, Suit: Hearts}, 10},
	}

	for _, test := range tests {
		result := test.card.CribbagePlayValue()
		assert.Equal(t, test.expected, result)
	}
}

func TestToCardWithImages(t *testing.T) {
	baseURL := "http://localhost:8080"
	
	// Test face up card
	card := Card{Rank: Ace, Suit: Hearts, FaceUp: true}
	result := card.ToCardWithImages(baseURL)
	
	assert.Equal(t, Ace, result.Rank)
	assert.Equal(t, Hearts, result.Suit)
	assert.True(t, result.FaceUp)
	assert.NotNil(t, result.Images)
	assert.Contains(t, result.Images["icon"], "1_0.png")
	assert.Contains(t, result.Images["small"], "1_0.png")
	assert.Contains(t, result.Images["large"], "1_0.png")
	
	// Test face down card
	cardDown := Card{Rank: King, Suit: Spades, FaceUp: false}
	resultDown := cardDown.ToCardWithImages(baseURL)
	
	assert.Equal(t, King, resultDown.Rank)
	assert.Equal(t, Spades, resultDown.Suit)
	assert.False(t, resultDown.FaceUp)
	assert.Contains(t, resultDown.Images["icon"], "back.png")
	assert.Contains(t, resultDown.Images["small"], "back.png")
	assert.Contains(t, resultDown.Images["large"], "back.png")
}

func TestToCardWithImagesPtr(t *testing.T) {
	baseURL := "http://localhost:8080"
	
	// Test nil pointer
	result := ToCardWithImagesPtr(nil, baseURL)
	assert.Equal(t, CardWithImages{}, result)
	
	// Test valid pointer
	card := &Card{Rank: Ace, Suit: Hearts, FaceUp: true}
	result = ToCardWithImagesPtr(card, baseURL)
	assert.Equal(t, Ace, result.Rank)
	assert.Equal(t, Hearts, result.Suit)
	assert.True(t, result.FaceUp)
}

func TestDeckTypeString(t *testing.T) {
	tests := []struct {
		deckType DeckType
		expected string
	}{
		{Standard, "Standard"},
		{Spanish21, "Spanish21"},
		{DeckType(99), "Standard"}, // Default case
	}

	for _, test := range tests {
		result := test.deckType.String()
		assert.Equal(t, test.expected, result)
	}
}

func TestDeckTypeDescription(t *testing.T) {
	tests := []struct {
		deckType DeckType
		expected string
	}{
		{Standard, "Traditional 52-card deck with all ranks from Ace to King in all four suits"},
		{Spanish21, "Spanish 21 deck with 48 cards - all 10s removed, perfect for Spanish Blackjack"},
		{DeckType(99), "Traditional 52-card deck with all ranks from Ace to King in all four suits"}, // Default case
	}

	for _, test := range tests {
		result := test.deckType.Description()
		assert.Equal(t, test.expected, result)
	}
}

func TestDeckTypeCardsPerDeck(t *testing.T) {
	tests := []struct {
		deckType DeckType
		expected int
	}{
		{Standard, 52},
		{Spanish21, 48},
		{DeckType(99), 52}, // Default case
	}

	for _, test := range tests {
		result := test.deckType.CardsPerDeck()
		assert.Equal(t, test.expected, result)
	}
}

func TestGetAllDeckTypes(t *testing.T) {
	deckTypes := GetAllDeckTypes()
	assert.Equal(t, 2, len(deckTypes))
	assert.Contains(t, deckTypes, Standard)
	assert.Contains(t, deckTypes, Spanish21)
}

func TestParseDeckType(t *testing.T) {
	tests := []struct {
		input    string
		expected DeckType
	}{
		{"spanish21", Spanish21},
		{"Spanish21", Spanish21},
		{"SPANISH21", Spanish21},
		{"spanish_21", Spanish21},
		{"spanish-21", Spanish21},
		{"standard", Standard},
		{"Standard", Standard},
		{"normal", Standard},
		{"regular", Standard},
		{"invalid", Standard}, // Default
		{"", Standard},        // Default
	}

	for _, test := range tests {
		result := ParseDeckType(test.input)
		assert.Equal(t, test.expected, result, "Input: "+test.input)
	}
}

func TestGenerateDeckName(t *testing.T) {
	// Test that deck names are generated
	name1 := GenerateDeckName()
	name2 := GenerateDeckName()
	
	assert.NotEmpty(t, name1)
	assert.NotEmpty(t, name2)
	assert.Contains(t, name1, " ") // Should contain a space between adjective and noun
	assert.Contains(t, name2, " ")
	
	// Test that names contain valid words
	parts1 := strings.Split(name1, " ")
	parts2 := strings.Split(name2, " ")
	assert.Equal(t, 2, len(parts1))
	assert.Equal(t, 2, len(parts2))
}

func TestGameTypeString(t *testing.T) {
	tests := []struct {
		gameType GameType
		expected string
	}{
		{Blackjack, "Blackjack"},
		{Poker, "Poker"},
		{Cribbage, "Cribbage"},
		{GameType(99), "Blackjack"}, // Default case
	}

	for _, test := range tests {
		result := test.gameType.String()
		assert.Equal(t, test.expected, result)
	}
}

func TestGameStatusString(t *testing.T) {
	tests := []struct {
		status   GameStatus
		expected string
	}{
		{GameWaiting, "waiting"},
		{GameInProgress, "in_progress"},
		{GameFinished, "finished"},
		{GameStatus(99), "waiting"}, // Default case
	}

	for _, test := range tests {
		result := test.status.String()
		assert.Equal(t, test.expected, result)
	}
}