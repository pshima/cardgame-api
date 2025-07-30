package managers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCustomDeckManager(t *testing.T) {
	cdm := NewCustomDeckManager()
	assert.NotNil(t, cdm)
	assert.NotNil(t, cdm.decks)
	assert.Equal(t, 0, len(cdm.ListDecks()))
}

func TestCustomDeckManagerCreateDeck(t *testing.T) {
	cdm := NewCustomDeckManager()
	
	deck := cdm.CreateDeck("Test Deck")
	assert.NotNil(t, deck)
	assert.NotEmpty(t, deck.ID)
	assert.Equal(t, "Test Deck", deck.Name)
	assert.Equal(t, 0, len(deck.Cards))
	assert.False(t, deck.Created.IsZero())
	assert.False(t, deck.LastUsed.IsZero())
	assert.Equal(t, 1, len(cdm.ListDecks()))
}

func TestCustomDeckManagerGetDeck(t *testing.T) {
	cdm := NewCustomDeckManager()
	
	// Create a deck
	deck := cdm.CreateDeck("Test Deck")
	assert.NotNil(t, deck)
	
	// Retrieve the deck
	retrieved, exists := cdm.GetDeck(deck.ID)
	assert.True(t, exists)
	assert.NotNil(t, retrieved)
	assert.Equal(t, deck.ID, retrieved.ID)
	assert.Equal(t, deck.Name, retrieved.Name)
	
	// Try to get non-existent deck
	missing, exists := cdm.GetDeck("non-existent-id")
	assert.False(t, exists)
	assert.Nil(t, missing)
}

func TestCustomDeckManagerDeleteDeck(t *testing.T) {
	cdm := NewCustomDeckManager()
	
	// Create multiple decks
	deck1 := cdm.CreateDeck("Deck 1")
	deck2 := cdm.CreateDeck("Deck 2")
	assert.Equal(t, 2, len(cdm.ListDecks()))
	
	// Delete first deck
	deleted := cdm.DeleteDeck(deck1.ID)
	assert.True(t, deleted)
	assert.Equal(t, 1, len(cdm.ListDecks()))
	
	// Verify deck is gone
	_, exists := cdm.GetDeck(deck1.ID)
	assert.False(t, exists)
	
	// Verify second deck still exists
	retrieved, exists := cdm.GetDeck(deck2.ID)
	assert.True(t, exists)
	assert.Equal(t, deck2.ID, retrieved.ID)
	
	// Try to delete non-existent deck
	deleted = cdm.DeleteDeck("non-existent-id")
	assert.False(t, deleted)
	assert.Equal(t, 1, len(cdm.ListDecks()))
}

func TestCustomDeckManagerListDecks(t *testing.T) {
	cdm := NewCustomDeckManager()
	
	// No decks initially
	decks := cdm.ListDecks()
	assert.Equal(t, 0, len(decks))
	
	// Create some decks
	deck1 := cdm.CreateDeck("Deck 1")
	deck2 := cdm.CreateDeck("Deck 2")
	deck3 := cdm.CreateDeck("Deck 3")
	
	// List decks
	decks = cdm.ListDecks()
	assert.Equal(t, 3, len(decks))
	
	// Verify all decks are included
	deckIDs := make(map[string]bool)
	for _, deck := range decks {
		deckIDs[deck.ID] = true
	}
	assert.True(t, deckIDs[deck1.ID])
	assert.True(t, deckIDs[deck2.ID])
	assert.True(t, deckIDs[deck3.ID])
}

// Note: AddCard functionality is handled by the service layer, not the manager

// Note: DeleteCard functionality is handled by the service layer, not the manager

// Note: GetCards functionality is handled by the service layer, not the manager

// Note: GetCard functionality is handled by the service layer, not the manager