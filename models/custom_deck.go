package models

import (
	"time"

	"github.com/google/uuid"
)

// CustomCard represents a user-defined card with flexible attributes and optional game compatibility.
// Cards can have string or numeric ranks, custom attributes, and tombstone deletion for data integrity.
type CustomCard struct {
	Index          int                    `json:"index"`
	Name           string                 `json:"name"`
	Rank           interface{}            `json:"rank,omitempty"`
	Suit           string                 `json:"suit,omitempty"`
	GameCompatible bool                   `json:"game_compatible"`
	Attributes     map[string]string      `json:"attributes"`
	Deleted        bool                   `json:"deleted"`
}

// UpdateGameCompatibility determines if the card can be used in traditional card games.
// Cards with numeric ranks and suits are compatible, while string ranks are not.
func (cc *CustomCard) UpdateGameCompatibility() {
	if cc.Rank == nil || cc.Suit == "" {
		cc.GameCompatible = false
		return
	}
	
	switch cc.Rank.(type) {
	case int, int32, int64, float32, float64:
		cc.GameCompatible = true
	case string:
		cc.GameCompatible = false
	default:
		cc.GameCompatible = false
	}
}

// GetNumericRank extracts the numeric value from the card's rank if possible.
// Returns the rank as an integer and whether the conversion was successful.
func (cc *CustomCard) GetNumericRank() (int, bool) {
	if cc.Rank == nil {
		return 0, false
	}
	
	switch v := cc.Rank.(type) {
	case int:
		return v, true
	case int32:
		return int(v), true
	case int64:
		return int(v), true
	case float32:
		return int(v), true
	case float64:
		return int(v), true
	default:
		return 0, false
	}
}

// CustomDeck represents a collection of user-defined custom cards with metadata.
// It tracks card indices for consistent referencing and usage timestamps for cleanup.
type CustomDeck struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Cards       []*CustomCard `json:"cards"`
	NextIndex   int           `json:"next_index"`
	Created     time.Time     `json:"created"`
	LastUsed    time.Time     `json:"last_used"`
}

// NewCustomDeckTemplate creates a new empty custom deck with the given name.
// Initializes all required fields including UUID and timestamps for proper tracking.
func NewCustomDeckTemplate(name string) *CustomDeck {
	return &CustomDeck{
		ID:        uuid.New().String(),
		Name:      name,
		Cards:     []*CustomCard{},
		NextIndex: 0,
		Created:   time.Now(),
		LastUsed:  time.Now(),
	}
}

func (cd *CustomDeck) UpdateLastUsed() {
	cd.LastUsed = time.Now()
}

func (cd *CustomDeck) AddCard(name string, rank interface{}, suit string, attributes map[string]string) *CustomCard {
	if attributes == nil {
		attributes = make(map[string]string)
	}
	
	card := &CustomCard{
		Index:      cd.NextIndex,
		Name:       name,
		Rank:       rank,
		Suit:       suit,
		Attributes: attributes,
		Deleted:    false,
	}
	
	card.UpdateGameCompatibility()
	cd.Cards = append(cd.Cards, card)
	cd.NextIndex++
	cd.UpdateLastUsed()
	
	return card
}

func (cd *CustomDeck) GetCard(index int) *CustomCard {
	for _, card := range cd.Cards {
		if card.Index == index {
			return card
		}
	}
	return nil
}

func (cd *CustomDeck) DeleteCard(index int) bool {
	card := cd.GetCard(index)
	if card == nil {
		return false
	}
	
	card.Deleted = true
	cd.UpdateLastUsed()
	return true
}

func (cd *CustomDeck) ListCards(includeDeleted bool) []*CustomCard {
	if includeDeleted {
		return cd.Cards
	}
	
	activeCards := []*CustomCard{}
	for _, card := range cd.Cards {
		if !card.Deleted {
			activeCards = append(activeCards, card)
		}
	}
	return activeCards
}

func (cd *CustomDeck) GetGameCompatibleCards() []*CustomCard {
	compatibleCards := []*CustomCard{}
	for _, card := range cd.Cards {
		if !card.Deleted && card.GameCompatible {
			compatibleCards = append(compatibleCards, card)
		}
	}
	return compatibleCards
}

func (cd *CustomDeck) CardCount() int {
	count := 0
	for _, card := range cd.Cards {
		if !card.Deleted {
			count++
		}
	}
	return count
}