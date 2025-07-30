package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/peteshima/cardgame-api/validators"
	"github.com/peteshima/cardgame-api/api"
)

func (h *HandlerDependencies) CreateCustomDeck(c *gin.Context) {
	var req api.CreateCustomDeckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON: " + err.Error(),
		})
		return
	}

	name := validators.SanitizeString(req.Name, 128)
	if !validators.ValidateDeckName(name) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Deck name must be 1-128 characters",
		})
		return
	}

	deck := h.CustomDeckService.CreateCustomDeck(name)
	c.JSON(http.StatusCreated, gin.H{
		"id":      deck.ID,
		"name":    deck.Name,
		"message": "Custom deck created successfully",
		"created": deck.Created,
	})
}

func (h *HandlerDependencies) ListCustomDecks(c *gin.Context) {
	decks := h.CustomDeckService.ListCustomDecks()
	
	deckSummaries := make([]gin.H, len(decks))
	for i, deck := range decks {
		deckSummaries[i] = gin.H{
			"id":         deck.ID,
			"name":       deck.Name,
			"card_count": deck.CardCount(),
			"created":    deck.Created,
			"last_used":  deck.LastUsed,
		}
	}
	
	c.JSON(http.StatusOK, gin.H{
		"decks": deckSummaries,
		"count": len(decks),
	})
}

func (h *HandlerDependencies) GetCustomDeck(c *gin.Context) {
	deckID := validators.SanitizeString(c.Param("deckId"), 50)
	if !validators.ValidateUUID(deckID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid deck ID format",
		})
		return
	}

	deck, exists := h.CustomDeckService.GetCustomDeck(deckID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Custom deck not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         deck.ID,
		"name":       deck.Name,
		"card_count": deck.CardCount(),
		"cards":      deck.ListCards(false),
		"created":    deck.Created,
		"last_used":  deck.LastUsed,
	})
}

// Additional custom deck handlers would follow the same pattern