package services

import (
	"github.com/peteshima/cardgame-api/managers"
	"github.com/peteshima/cardgame-api/models"
)

// CustomDeckService provides business logic operations for custom deck management
type CustomDeckService struct {
	customDeckManager *managers.CustomDeckManager
}

// NewCustomDeckService creates a new custom deck service instance
func NewCustomDeckService(customDeckManager *managers.CustomDeckManager) *CustomDeckService {
	return &CustomDeckService{
		customDeckManager: customDeckManager,
	}
}

// CreateCustomDeck creates a new custom deck
func (cds *CustomDeckService) CreateCustomDeck(name string) *models.CustomDeck {
	return cds.customDeckManager.CreateDeck(name)
}

// ListCustomDecks returns all custom decks
func (cds *CustomDeckService) ListCustomDecks() []*models.CustomDeck {
	return cds.customDeckManager.ListDecks()
}

// GetCustomDeck retrieves a custom deck by ID
func (cds *CustomDeckService) GetCustomDeck(deckID string) (*models.CustomDeck, bool) {
	return cds.customDeckManager.GetDeck(deckID)
}

// DeleteCustomDeck removes a custom deck
func (cds *CustomDeckService) DeleteCustomDeck(deckID string) bool {
	return cds.customDeckManager.DeleteDeck(deckID)
}

// AddCustomCard adds a card to a custom deck
func (cds *CustomDeckService) AddCustomCard(deckID string, name string, rank interface{}, suit string, attributes map[string]string) (*models.CustomDeck, *models.CustomCard, bool) {
	deck, exists := cds.customDeckManager.GetDeck(deckID)
	if !exists {
		return nil, nil, false
	}
	
	if deck.CardCount() >= 2000 {
		return deck, nil, false // Max limit reached
	}
	
	card := deck.AddCard(name, rank, suit, attributes)
	return deck, card, true
}

// ListCustomCards returns cards in a custom deck
func (cds *CustomDeckService) ListCustomCards(deckID string, includeDeleted bool) (*models.CustomDeck, []*models.CustomCard, bool) {
	deck, exists := cds.customDeckManager.GetDeck(deckID)
	if !exists {
		return nil, nil, false
	}
	
	cards := deck.ListCards(includeDeleted)
	return deck, cards, true
}

// GetCustomCard retrieves a specific card from a custom deck
func (cds *CustomDeckService) GetCustomCard(deckID string, cardIndex int) (*models.CustomDeck, *models.CustomCard, bool) {
	deck, exists := cds.customDeckManager.GetDeck(deckID)
	if !exists {
		return nil, nil, false
	}
	
	card := deck.GetCard(cardIndex)
	if card == nil {
		return deck, nil, false
	}
	
	return deck, card, true
}

// DeleteCustomCard marks a card as deleted (tombstone deletion)
func (cds *CustomDeckService) DeleteCustomCard(deckID string, cardIndex int) (*models.CustomDeck, bool) {
	deck, exists := cds.customDeckManager.GetDeck(deckID)
	if !exists {
		return nil, false
	}
	
	deleted := deck.DeleteCard(cardIndex)
	return deck, deleted
}