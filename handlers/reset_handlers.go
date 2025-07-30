package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/peteshima/cardgame-api/models"
	"github.com/peteshima/cardgame-api/validators"
)

// ResetDeck restores a game's deck to its original full state with all cards.
// It maintains the same deck type and count while shuffling the restored deck.
func (h *HandlerDependencies) ResetDeck(c *gin.Context) {
	gameID := validators.SanitizeString(c.Param("gameId"), 50)
	if !validators.ValidateUUID(gameID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid game ID format",
		})
		return
	}
	
	game, exists := h.GameService.ResetGameDeck(gameID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Game not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"game_id":        game.ID,
		"deck_name":      game.Deck.Name,
		"deck_type":      game.Deck.DeckType.String(),
		"message":        "Deck reset",
		"remaining_cards": game.Deck.RemainingCards(),
	})
}

func (h *HandlerDependencies) ResetDeckWithDecks(c *gin.Context) {
	gameID := validators.SanitizeString(c.Param("gameId"), 50)
	decksStr := validators.SanitizeString(c.Param("decks"), 10)
	
	if !validators.ValidateUUID(gameID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid game ID format",
		})
		return
	}
	
	numDecks, valid := validators.ValidateNumber(decksStr)
	if !valid || numDecks <= 0 || numDecks > 100 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid decks parameter (must be 1-100)",
		})
		return
	}

	game, exists := h.GameService.ResetGameDeckWithDecks(gameID, numDecks)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Game not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"game_id":        game.ID,
		"deck_name":      game.Deck.Name,
		"deck_type":      game.Deck.DeckType.String(),
		"message":        "Deck reset to " + decksStr + " decks",
		"remaining_cards": game.Deck.RemainingCards(),
		"num_decks":      numDecks,
	})
}

func (h *HandlerDependencies) ResetDeckWithType(c *gin.Context) {
	gameID := validators.SanitizeString(c.Param("gameId"), 50)
	decksStr := validators.SanitizeString(c.Param("decks"), 10)
	typeStr := validators.SanitizeString(c.Param("type"), 20)
	
	if !validators.ValidateUUID(gameID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid game ID format",
		})
		return
	}
	
	if !validators.ValidateDeckType(typeStr) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid deck type parameter",
		})
		return
	}
	
	numDecks, valid := validators.ValidateNumber(decksStr)
	if !valid || numDecks <= 0 || numDecks > 100 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid decks parameter (must be 1-100)",
		})
		return
	}
	
	deckType := models.ParseDeckType(typeStr)
	game, exists := h.GameService.ResetGameDeckWithType(gameID, numDecks, deckType)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Game not found",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"game_id":        game.ID,
		"deck_name":      game.Deck.Name,
		"deck_type":      game.Deck.DeckType.String(),
		"message":        "Deck reset to " + decksStr + " " + deckType.String() + " decks",
		"remaining_cards": game.Deck.RemainingCards(),
		"num_decks":      numDecks,
	})
}