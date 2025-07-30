package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/peteshima/cardgame-api/config"
	"github.com/peteshima/cardgame-api/validators"
	"github.com/peteshima/cardgame-api/api"
)

func (h *HandlerDependencies) CreateNewCribbageGame(c *gin.Context) {
	game := h.CribbageService.CreateCribbageGame()
	c.JSON(http.StatusOK, gin.H{
		"game_id":        game.ID,
		"game_type":      game.GameType.String(),
		"deck_name":      game.Deck.Name,
		"deck_type":      game.Deck.DeckType.String(),
		"max_players":    game.MaxPlayers,
		"current_players": len(game.Players),
		"message":        "New Cribbage game created",
		"remaining_cards": game.Deck.RemainingCards(),
		"created":        game.Created,
	})
}

func (h *HandlerDependencies) StartCribbageGame(c *gin.Context) {
	gameID := c.Param("gameId")
	if !validators.ValidateUUID(gameID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid game ID format",
		})
		return
	}

	game, err := h.CribbageService.StartCribbageGame(gameID)
	if game == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Game not found",
		})
		return
	}
	
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"game_id":        game.ID,
		"game_type":      game.GameType.String(),
		"status":         game.Status.String(),
		"phase":          game.CribbageState.Phase.String(),
		"dealer":         game.CribbageState.Dealer,
		"current_player": game.CurrentPlayer,
		"message":        "Cribbage game started",
	})
}

func (h *HandlerDependencies) CribbageDiscard(c *gin.Context) {
	gameID := c.Param("gameId")
	playerID := c.Param("playerId")
	
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

	var request api.CribbageDiscardRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	game, player, err := h.CribbageService.CribbageDiscard(gameID, playerID, request.CardIndices)
	if game == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Game not found",
		})
		return
	}
	
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := gin.H{
		"game_id":     game.ID,
		"player_id":   playerID,
		"player_name": player.Name,
		"phase":       game.CribbageState.Phase.String(),
		"message":     "Cards discarded to crib",
	}

	// If phase changed to play, include starter card
	if game.CribbageState.Phase.String() == "play" && game.CribbageState.Starter != nil {
		baseURL := config.GetBaseURL(c)
		response["starter"] = game.CribbageState.Starter.ToCardWithImages(baseURL)
		response["message"] = "Cards discarded, starter cut, play phase begun"
	}

	c.JSON(http.StatusOK, response)
}

// Additional cribbage handlers would go here but I'll keep it minimal for now
// The pattern is the same: validate inputs, call service, return response