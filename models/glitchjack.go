package models

import (
	"math/rand"
	"time"
)

// GlitchjackGame represents a Glitchjack game with the same rules as Blackjack
// but with a randomly generated deck composition from standard cards.
type GlitchjackGame struct {
	Game         *Game             `json:"game"`
	Results      map[string]string `json:"results,omitempty"`
	DealerHidden bool              `json:"dealer_hidden"`
}

// NewGlitchjackDeck creates a new deck for Glitchjack with 52 randomly selected cards
// from the standard deck. Cards can repeat (e.g., multiple Ace of Hearts).
func NewGlitchjackDeck() *Deck {
	rand.Seed(time.Now().UnixNano())
	
	deck := &Deck{
		Cards:    make([]Card, 0, 52),
		Name:     GenerateDeckName(),
		DeckType: Standard, // Using Standard type but with random composition
	}
	
	// Generate 52 random cards from standard deck possibilities
	for i := 0; i < 52; i++ {
		// Random suit (0-3)
		suit := Suit(rand.Intn(4))
		// Random rank (1-13)
		rank := Rank(rand.Intn(13) + 1)
		
		deck.Cards = append(deck.Cards, Card{
			Rank:   rank,
			Suit:   suit,
			FaceUp: false,
		})
	}
	
	deck.Shuffle()
	return deck
}

// CalculateGlitchjackHand calculates the value of a hand in Glitchjack.
// Uses the same rules as Blackjack for hand calculation.
func CalculateGlitchjackHand(cards []*Card) int {
	total := 0
	aces := 0
	
	for _, card := range cards {
		value := card.BlackjackValue()
		if card.Rank == Ace {
			aces++
		}
		total += value
	}
	
	// Adjust for aces (convert from 11 to 1 if needed)
	for aces > 0 && total > 21 {
		total -= 10
		aces--
	}
	
	return total
}

// IsGlitchjackBlackjack checks if a hand is a blackjack (21 with 2 cards).
func IsGlitchjackBlackjack(cards []*Card) bool {
	return len(cards) == 2 && CalculateGlitchjackHand(cards) == 21
}

// IsBlackjack is an alias for IsGlitchjackBlackjack for consistency with existing code.
func IsBlackjack(cards []*Card) bool {
	return IsGlitchjackBlackjack(cards)
}

// GlitchjackResult represents the possible outcomes of a Glitchjack game.
type GlitchjackResult string

const (
	GlitchjackResultBlackjack GlitchjackResult = "blackjack"
	GlitchjackResultWin       GlitchjackResult = "win"
	GlitchjackResultPush      GlitchjackResult = "push"
	GlitchjackResultBust      GlitchjackResult = "bust"
	GlitchjackResultLose      GlitchjackResult = "lose"
)

// String returns the string representation of a Glitchjack result.
func (gr GlitchjackResult) String() string {
	return string(gr)
}

