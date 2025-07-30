package managers

import (
	"sync"

	"github.com/peteshima/cardgame-api/models"
)

type CustomDeckManager struct {
	decks map[string]*models.CustomDeck
	mutex sync.RWMutex
}

func NewCustomDeckManager() *CustomDeckManager {
	return &CustomDeckManager{
		decks: make(map[string]*models.CustomDeck),
	}
}

func (cdm *CustomDeckManager) CreateDeck(name string) *models.CustomDeck {
	cdm.mutex.Lock()
	defer cdm.mutex.Unlock()
	
	deck := models.NewCustomDeckTemplate(name)
	cdm.decks[deck.ID] = deck
	return deck
}

func (cdm *CustomDeckManager) GetDeck(deckID string) (*models.CustomDeck, bool) {
	cdm.mutex.RLock()
	defer cdm.mutex.RUnlock()
	
	deck, exists := cdm.decks[deckID]
	if exists {
		deck.UpdateLastUsed()
	}
	return deck, exists
}

func (cdm *CustomDeckManager) DeleteDeck(deckID string) bool {
	cdm.mutex.Lock()
	defer cdm.mutex.Unlock()
	
	_, exists := cdm.decks[deckID]
	if exists {
		delete(cdm.decks, deckID)
	}
	return exists
}

func (cdm *CustomDeckManager) ListDecks() []*models.CustomDeck {
	cdm.mutex.RLock()
	defer cdm.mutex.RUnlock()
	
	decks := make([]*models.CustomDeck, 0, len(cdm.decks))
	for _, deck := range cdm.decks {
		decks = append(decks, deck)
	}
	return decks
}