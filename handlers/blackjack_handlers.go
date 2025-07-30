package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/peteshima/cardgame-api/validators"
)

func (h *HandlerDependencies) StartBlackjackGame(c *gin.Context) {
	gameID := validators.SanitizeString(c.Param("gameId"), 50)
	if !validators.ValidateUUID(gameID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid game ID format",
		})
		return
	}
	
	game, err := h.BlackjackService.StartBlackjackGame(gameID)
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
		"game_id": game.ID,
		"status":  game.Status.String(),
		"message": "Blackjack game started",
		"current_player": game.CurrentPlayer,
	})
}

func (h *HandlerDependencies) PlayerHit(c *gin.Context) {
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
	
	game, player, err := h.BlackjackService.PlayerHit(gameID, playerID)
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

	handValue, hasBlackjack := player.BlackjackHandValue()
	isBusted := player.IsBusted()

	c.JSON(http.StatusOK, gin.H{
		"game_id":      game.ID,
		"player_id":    playerID,
		"player_name":  player.Name,
		"hand_value":   handValue,
		"hand_size":    player.HandSize(),
		"has_blackjack": hasBlackjack,
		"is_busted":    isBusted,
		"message":      "Card dealt to " + player.Name,
	})
}

func (h *HandlerDependencies) PlayerStand(c *gin.Context) {
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
	
	game, player, err := h.BlackjackService.PlayerStand(gameID, playerID)
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
		"player_id":      playerID,
		"player_name":    player.Name,
		"status":         game.Status.String(),
		"current_player": game.CurrentPlayer,
		"message":        player.Name + " stands",
	})
}

func (h *HandlerDependencies) GetGameResults(c *gin.Context) {
	gameID := validators.SanitizeString(c.Param("gameId"), 50)
	if !validators.ValidateUUID(gameID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid game ID format",
		})
		return
	}
	
	game, results, exists := h.BlackjackService.GetGameResults(gameID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Game not found",
		})
		return
	}

	dealerValue, dealerBlackjack := game.Dealer.BlackjackHandValue()

	// Add player details to results
	playerResults := make([]gin.H, 0)
	for _, player := range game.Players {
		playerValue, playerBlackjack := player.BlackjackHandValue()
		playerResults = append(playerResults, gin.H{
			"player_id":     player.ID,
			"player_name":   player.Name,
			"hand_value":    playerValue,
			"has_blackjack": playerBlackjack,
			"is_busted":     player.IsBusted(),
			"result":        results[player.ID],
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"game_id": game.ID,
		"status":  game.Status.String(),
		"dealer": gin.H{
			"hand_value":    dealerValue,
			"has_blackjack": dealerBlackjack,
			"is_busted":     game.Dealer.IsBusted(),
		},
		"players": playerResults,
		"results": results,
	})
}