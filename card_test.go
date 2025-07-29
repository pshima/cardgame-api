package main

import (
	"testing"
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
		{Suit(99), "Unknown"},
	}

	for _, test := range tests {
		result := test.suit.String()
		if result != test.expected {
			t.Errorf("Suit(%d).String() = %s; expected %s", test.suit, result, test.expected)
		}
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
	}

	for _, test := range tests {
		result := test.rank.String()
		if result != test.expected {
			t.Errorf("Rank(%d).String() = %s; expected %s", test.rank, result, test.expected)
		}
	}
}

func TestCardString(t *testing.T) {
	card := Card{Rank: Ace, Suit: Spades}
	expected := "Ace of Spades"
	result := card.String()
	if result != expected {
		t.Errorf("Card.String() = %s; expected %s", result, expected)
	}
}

func TestCardValue(t *testing.T) {
	tests := []struct {
		card     Card
		expected int
	}{
		{Card{Rank: Ace, Suit: Hearts}, 1},
		{Card{Rank: Five, Suit: Clubs}, 5},
		{Card{Rank: Ten, Suit: Diamonds}, 10},
		{Card{Rank: Jack, Suit: Spades}, 11},
		{Card{Rank: Queen, Suit: Hearts}, 12},
		{Card{Rank: King, Suit: Clubs}, 13},
	}

	for _, test := range tests {
		result := test.card.Value()
		if result != test.expected {
			t.Errorf("Card{%s, %s}.Value() = %d; expected %d", test.card.Rank, test.card.Suit, result, test.expected)
		}
	}
}

func TestNewDeck(t *testing.T) {
	deck := NewDeck()
	
	if deck == nil {
		t.Fatal("NewDeck() returned nil")
	}
	
	if len(deck.Cards) != 52 {
		t.Errorf("NewDeck() created deck with %d cards; expected 52", len(deck.Cards))
	}
	
	if deck.RemainingCards() != 52 {
		t.Errorf("NewDeck() RemainingCards() = %d; expected 52", deck.RemainingCards())
	}
	
	if deck.IsEmpty() {
		t.Error("NewDeck() created empty deck; expected full deck")
	}
}

func TestDeckReset(t *testing.T) {
	deck := NewDeck()
	
	deck.Deal()
	deck.Deal()
	
	if deck.RemainingCards() != 50 {
		t.Errorf("After dealing 2 cards, deck has %d cards; expected 50", deck.RemainingCards())
	}
	
	deck.Reset()
	
	if deck.RemainingCards() != 52 {
		t.Errorf("After reset, deck has %d cards; expected 52", deck.RemainingCards())
	}
}

func TestDeckShuffle(t *testing.T) {
	deck1 := NewDeck()
	
	originalOrder := make([]Card, len(deck1.Cards))
	copy(originalOrder, deck1.Cards)
	
	deck1.Shuffle()
	
	if len(deck1.Cards) != 52 {
		t.Errorf("After shuffle, deck has %d cards; expected 52", len(deck1.Cards))
	}
	
	cardCounts := make(map[Card]int)
	for _, card := range deck1.Cards {
		cardCounts[card]++
	}
	
	if len(cardCounts) != 52 {
		t.Errorf("After shuffle, deck has %d unique cards; expected 52", len(cardCounts))
	}
	
	for card, count := range cardCounts {
		if count != 1 {
			t.Errorf("Card %s appears %d times; expected 1", card, count)
		}
	}
	
	same := true
	for i, card := range deck1.Cards {
		if card != originalOrder[i] {
			same = false
			break
		}
	}
	
	if same {
		t.Error("Shuffle did not change card order (very unlikely but possible)")
	}
}

func TestDeckDeal(t *testing.T) {
	deck := NewDeck()
	originalFirst := deck.Cards[0]
	
	card := deck.Deal()
	
	if card == nil {
		t.Fatal("Deal() returned nil from full deck")
	}
	
	if *card != originalFirst {
		t.Errorf("Deal() returned %s; expected %s", card, originalFirst)
	}
	
	if deck.RemainingCards() != 51 {
		t.Errorf("After dealing 1 card, deck has %d cards; expected 51", deck.RemainingCards())
	}
}

func TestDeckDealEmpty(t *testing.T) {
	deck := &Deck{Cards: []Card{}}
	
	card := deck.Deal()
	
	if card != nil {
		t.Errorf("Deal() from empty deck returned %s; expected nil", card)
	}
	
	if !deck.IsEmpty() {
		t.Error("Empty deck IsEmpty() returned false; expected true")
	}
}

func TestDeckIsEmpty(t *testing.T) {
	deck := NewDeck()
	
	if deck.IsEmpty() {
		t.Error("Full deck IsEmpty() returned true; expected false")
	}
	
	for i := 0; i < 52; i++ {
		deck.Deal()
	}
	
	if !deck.IsEmpty() {
		t.Error("Empty deck IsEmpty() returned false; expected true")
	}
}

func TestDeckRemainingCards(t *testing.T) {
	deck := NewDeck()
	
	for i := 52; i > 0; i-- {
		if deck.RemainingCards() != i {
			t.Errorf("RemainingCards() = %d; expected %d", deck.RemainingCards(), i)
		}
		deck.Deal()
	}
	
	if deck.RemainingCards() != 0 {
		t.Errorf("RemainingCards() after dealing all cards = %d; expected 0", deck.RemainingCards())
	}
}

func TestDeckContainsAllCards(t *testing.T) {
	deck := NewDeck()
	
	expectedCards := make(map[Card]bool)
	for suit := Hearts; suit <= Spades; suit++ {
		for rank := Ace; rank <= King; rank++ {
			expectedCards[Card{Rank: rank, Suit: suit}] = true
		}
	}
	
	if len(expectedCards) != 52 {
		t.Errorf("Expected 52 unique cards, got %d", len(expectedCards))
	}
	
	for _, card := range deck.Cards {
		if !expectedCards[card] {
			t.Errorf("Deck contains unexpected card: %s", card)
		}
		delete(expectedCards, card)
	}
	
	if len(expectedCards) > 0 {
		t.Errorf("Deck is missing %d cards", len(expectedCards))
		for card := range expectedCards {
			t.Errorf("Missing card: %s", card)
		}
	}
}