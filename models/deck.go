package models

import (
	"math/rand"
	"time"
)

// Deck represents a collection of playing cards with metadata.
// It maintains the card order for dealing and tracks the deck type for game rules.
type Deck struct {
	Cards    []Card   `json:"cards"`
	Name     string   `json:"name"`
	DeckType DeckType `json:"deck_type"`
}

// NewDeck creates a single standard 52-card deck.
// This is a convenience function that calls NewMultiDeck with 1 deck.
func NewDeck() *Deck {
	return NewMultiDeck(1)
}

// NewMultiDeck creates a deck with multiple standard decks shuffled together.
// This is commonly used in Blackjack where multiple decks reduce card counting effectiveness.
func NewMultiDeck(numDecks int) *Deck {
	return NewCustomDeck(numDecks, Standard)
}

// NewCustomDeck creates a deck with specified count and type (Standard or Spanish21).
// It handles deck type differences like Spanish21 missing 10s and ensures proper card generation.
func NewCustomDeck(numDecks int, deckType DeckType) *Deck {
	if numDecks <= 0 {
		numDecks = 1
	}
	
	cardsPerDeck := 52
	if deckType == Spanish21 {
		cardsPerDeck = 48 // No 10s (4 cards removed per deck)
	}
	
	deck := &Deck{
		Cards:    make([]Card, 0, cardsPerDeck*numDecks),
		Name:     GenerateDeckName(),
		DeckType: deckType,
	}
	deck.ResetWithDecks(numDecks)
	return deck
}

// Reset restores the deck to a full single deck of the current type.
// All cards are restored and the deck is shuffled, maintaining the current deck type.
func (d *Deck) Reset() {
	d.ResetWithDecks(1)
}

// ResetWithDecks restores the deck with a specified number of decks.
// It maintains the current deck type while changing the number of deck copies.
func (d *Deck) ResetWithDecks(numDecks int) {
	d.ResetWithDecksAndType(numDecks, d.DeckType)
}

// ResetWithDecksAndType completely reconfigures the deck with new count and type.
// This allows changing both the number of decks and the deck type (Standard/Spanish21).
func (d *Deck) ResetWithDecksAndType(numDecks int, deckType DeckType) {
	if numDecks <= 0 {
		numDecks = 1
	}
	
	d.DeckType = deckType
	cardsPerDeck := 52
	if deckType == Spanish21 {
		cardsPerDeck = 48
	}
	
	d.Cards = make([]Card, 0, cardsPerDeck*numDecks)
	
	for i := 0; i < numDecks; i++ {
		for suit := Hearts; suit <= Spades; suit++ {
			for rank := Ace; rank <= King; rank++ {
				// Skip 10s for Spanish 21
				if deckType == Spanish21 && rank == Ten {
					continue
				}
				d.Cards = append(d.Cards, Card{Rank: rank, Suit: suit, FaceUp: false})
			}
		}
	}
}

// Shuffle randomizes the order of all cards in the deck using Fisher-Yates algorithm.
// This ensures fair card distribution and prevents predictable card sequences.
func (d *Deck) Shuffle() {
	rand.Seed(time.Now().UnixNano())
	for i := len(d.Cards) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	}
}

// Deal removes and returns the top card from the deck.
// Returns nil if the deck is empty, allowing callers to handle empty deck scenarios.
func (d *Deck) Deal() *Card {
	if len(d.Cards) == 0 {
		return nil
	}
	card := d.Cards[0]
	d.Cards = d.Cards[1:]
	return &card
}

// RemainingCards returns the number of cards left in the deck.
// This is used for game logic and API responses to show deck status.
func (d *Deck) RemainingCards() int {
	return len(d.Cards)
}

// IsEmpty checks if the deck has no cards remaining.
// This is used to prevent dealing from empty decks and trigger deck resets.
func (d *Deck) IsEmpty() bool {
	return len(d.Cards) == 0
}

// DiscardPile represents a collection of discarded cards for various game mechanics.
type DiscardPile struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Cards []*Card `json:"cards"`
}

func (dp *DiscardPile) AddCard(card *Card) {
	dp.Cards = append(dp.Cards, card)
}

func (dp *DiscardPile) AddCards(cards []*Card) {
	dp.Cards = append(dp.Cards, cards...)
}

func (dp *DiscardPile) TopCard() *Card {
	if len(dp.Cards) == 0 {
		return nil
	}
	return dp.Cards[len(dp.Cards)-1]
}

// TakeTopCard removes and returns the most recently discarded card from the pile.
// Returns nil if the pile is empty, used for drawing from discard piles.
func (dp *DiscardPile) TakeTopCard() *Card {
	if len(dp.Cards) == 0 {
		return nil
	}
	card := dp.Cards[len(dp.Cards)-1]
	dp.Cards = dp.Cards[:len(dp.Cards)-1]
	return card
}

// Size returns the number of cards currently in the discard pile.
// This is used for game logic and API responses to show pile status.
func (dp *DiscardPile) Size() int {
	return len(dp.Cards)
}

// Clear removes all cards from the discard pile and returns them.
// This is used for reshuffling cards back into the deck or resetting games.
func (dp *DiscardPile) Clear() []*Card {
	cards := dp.Cards
	dp.Cards = []*Card{}
	return cards
}