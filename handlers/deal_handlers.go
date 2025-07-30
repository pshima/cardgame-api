package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/peteshima/cardgame-api/config"
	"github.com/peteshima/cardgame-api/validators"
	"github.com/peteshima/cardgame-api/api"
)

// DealCard draws a single card from the top of a game's deck and returns it face-up.
// It validates the game exists, updates metrics, and logs the card dealing operation.
func (h *HandlerDependencies) DealCard(c *gin.Context) {
	gameID := validators.SanitizeString(c.Param("gameId"), 50)
	if !validators.ValidateUUID(gameID) {
		h.Logger.Warn("Invalid game ID provided",
			zap.String("game_id", gameID),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid game ID format",
		})
		return
	}
	
	game, card, success := h.GameService.DealCard(gameID)
	if !success {
		if game == nil {
			h.Logger.Warn("Game not found",
				zap.String("game_id", gameID),
				zap.String("client_ip", c.ClientIP()),
			)
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Game not found",
			})
		} else {
			h.Logger.Warn("No cards remaining in deck",
				zap.String("game_id", gameID),
				zap.String("client_ip", c.ClientIP()),
			)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "No cards remaining in deck",
			})
		}
		return
	}
	
	h.Logger.Debug("Dealing card",
		zap.String("game_id", gameID),
		zap.Int("remaining_cards", game.Deck.RemainingCards()),
		zap.String("client_ip", c.ClientIP()),
	)
	
	// Update metrics
	h.updateCardsDealtMetric(c, gameID, game, 1)
	
	baseURL := config.GetBaseURL(c)
	cardWithImages := card.ToCardWithImages(baseURL)
	
	h.Logger.Info("Card dealt successfully",
		zap.String("game_id", gameID),
		zap.String("card", fmt.Sprintf("%s of %s", card.Rank, card.Suit)),
		zap.Int("remaining_cards", game.Deck.RemainingCards()),
		zap.String("client_ip", c.ClientIP()),
	)
	
	c.JSON(http.StatusOK, gin.H{
		"game_id":        game.ID,
		"deck_name":      game.Deck.Name,
		"card":           cardWithImages,
		"remaining_cards": game.Deck.RemainingCards(),
	})
}

// DealCards draws multiple cards from a game's deck based on the count parameter.
// It validates the count (1-52), ensures enough cards remain, and returns all dealt cards.
func (h *HandlerDependencies) DealCards(c *gin.Context) {
	gameID := validators.SanitizeString(c.Param("gameId"), 50)
	countStr := validators.SanitizeString(c.Param("count"), 10)
	
	if !validators.ValidateUUID(gameID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid game ID format",
		})
		return
	}
	
	count, valid := validators.ValidateNumber(countStr)
	if !valid || count <= 0 || count > 52 { // Reasonable upper limit
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid count parameter (must be 1-52)",
		})
		return
	}

	game, cards, success := h.GameService.DealCards(gameID, count)
	if !success {
		if game == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Game not found",
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Not enough cards remaining in deck",
			})
		}
		return
	}

	baseURL := config.GetBaseURL(c)
	cardsWithImages := convertCardsWithImages(cards, baseURL)

	c.JSON(http.StatusOK, gin.H{
		"game_id":        game.ID,
		"deck_name":      game.Deck.Name,
		"cards":          cardsWithImages,
		"cards_dealt":    len(cardsWithImages),
		"remaining_cards": game.Deck.RemainingCards(),
	})
}

func (h *HandlerDependencies) DealToPlayer(c *gin.Context) {
	gameID := validators.SanitizeString(c.Param("gameId"), 50)
	playerID := validators.SanitizeString(c.Param("playerId"), 50)
	
	if !validators.ValidateUUID(gameID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid game ID format",
		})
		return
	}
	
	if !validators.ValidatePlayerID(playerID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid player ID format",
		})
		return
	}
	
	game, player, card, success := h.GameService.DealToPlayer(gameID, playerID, false) // Default to face down
	if !success {
		if game == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Game not found",
			})
		} else if player == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Player not found",
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "No cards remaining in deck",
			})
		}
		return
	}

	baseURL := config.GetBaseURL(c)
	cardWithImages := card.ToCardWithImages(baseURL)
	
	c.JSON(http.StatusOK, gin.H{
		"game_id":         game.ID,
		"player_id":       playerID,
		"player_name":     player.Name,
		"card":            cardWithImages,
		"hand_size":       player.HandSize(),
		"remaining_cards": game.Deck.RemainingCards(),
		"message":         "Card dealt to " + player.Name,
	})
}

func (h *HandlerDependencies) DealToPlayerFaceUp(c *gin.Context) {
	gameID := validators.SanitizeString(c.Param("gameId"), 50)
	playerID := validators.SanitizeString(c.Param("playerId"), 50)
	faceUpStr := validators.SanitizeString(c.Param("faceUp"), 10)
	
	if !validators.ValidateUUID(gameID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid game ID format",
		})
		return
	}
	
	if !validators.ValidatePlayerID(playerID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid player ID format",
		})
		return
	}
	
	if !validators.ValidateBoolean(faceUpStr) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid faceUp parameter (must be true/false/1/0)",
		})
		return
	}
	
	faceUp := faceUpStr == "true" || faceUpStr == "1"
	
	game, player, card, success := h.GameService.DealToPlayer(gameID, playerID, faceUp)
	if !success {
		if game == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Game not found",
			})
		} else if player == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Player not found",
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "No cards remaining in deck",
			})
		}
		return
	}

	baseURL := config.GetBaseURL(c)
	cardWithImages := card.ToCardWithImages(baseURL)
	
	c.JSON(http.StatusOK, gin.H{
		"game_id":         game.ID,
		"player_id":       playerID,
		"player_name":     player.Name,
		"card":            cardWithImages,
		"face_up":         faceUp,
		"hand_size":       player.HandSize(),
		"remaining_cards": game.Deck.RemainingCards(),
		"message":         "Card dealt to " + player.Name,
	})
}

func (h *HandlerDependencies) DiscardToCard(c *gin.Context) {
	gameID := validators.SanitizeString(c.Param("gameId"), 50)
	pileID := validators.SanitizeString(c.Param("pileId"), 50)
	
	if !validators.ValidateUUID(gameID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid game ID format",
		})
		return
	}
	
	if !validators.ValidatePileID(pileID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid pile ID format",
		})
		return
	}
	
	var request api.DiscardRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}
	
	// Sanitize and validate player ID from request body
	request.PlayerID = validators.SanitizeString(request.PlayerID, 50)
	if !validators.ValidatePlayerID(request.PlayerID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid player ID format in request body",
		})
		return
	}

	if request.CardIndex < 0 || request.CardIndex >= 52 { // Reasonable upper limit
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid card index (must be 0-51)",
		})
		return
	}

	game, player, pile, card, success := h.GameService.DiscardCard(gameID, pileID, request.PlayerID, request.CardIndex)
	if !success {
		if game == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Game not found",
			})
		} else if player == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Player not found",
			})
		} else if pile == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Discard pile not found",
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid card index",
			})
		}
		return
	}
	
	baseURL := config.GetBaseURL(c)
	cardWithImages := card.ToCardWithImages(baseURL)

	c.JSON(http.StatusOK, gin.H{
		"game_id":       game.ID,
		"player_id":     request.PlayerID,
		"player_name":   player.Name,
		"card":          cardWithImages,
		"pile_id":       pileID,
		"pile_name":     pile.Name,
		"pile_size":     pile.Size(),
		"hand_size":     player.HandSize(),
		"message":       "Card discarded to " + pile.Name,
	})
}