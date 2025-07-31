package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/peteshima/cardgame-api/config"
	"github.com/peteshima/cardgame-api/models"
	"github.com/peteshima/cardgame-api/validators"
)

// CreateNewGlitchjackGame creates a new Glitchjack game with default settings.
// Glitchjack uses the same rules as Blackjack but with a randomly generated deck composition.
func (h *HandlerDependencies) CreateNewGlitchjackGame(c *gin.Context) {
	h.Logger.Debug("Creating new Glitchjack game",
		zap.String("client_ip", c.ClientIP()),
	)

	game := h.GlitchjackService.CreateGlitchjackGame()
	
	// Update metrics
	h.updateGamesCreatedMetric(c, models.Standard, 1)

	h.Logger.Info("Glitchjack game created successfully",
		zap.String("game_id", game.ID),
		zap.String("game_type", game.GameType.String()),
		zap.Int("remaining_cards", game.Deck.RemainingCards()),
		zap.String("client_ip", c.ClientIP()),
	)

	c.JSON(http.StatusOK, gin.H{
		"game_id":         game.ID,
		"game_type":       game.GameType.String(),
		"deck_name":       game.Deck.Name,
		"message":         "New Glitchjack game created with random deck composition",
		"remaining_cards": game.Deck.RemainingCards(),
		"max_players":     game.MaxPlayers,
		"created":         game.Created,
	})
}

// CreateNewGlitchjackGameWithDecks creates a new Glitchjack game with specified number of decks.
// Each deck contains 52 randomly selected cards from standard playing cards.
func (h *HandlerDependencies) CreateNewGlitchjackGameWithDecks(c *gin.Context) {
	decksStr := validators.SanitizeString(c.Param("decks"), 10)
	numDecks, valid := validators.ValidateNumber(decksStr)
	if !valid || numDecks <= 0 || numDecks > 100 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid decks parameter (must be 1-100)",
		})
		return
	}

	game := h.GlitchjackService.CreateGlitchjackGameWithOptions(numDecks, 6)
	
	c.JSON(http.StatusOK, gin.H{
		"game_id":         game.ID,
		"game_type":       game.GameType.String(),
		"deck_name":       game.Deck.Name,
		"message":         "New Glitchjack game created with " + decksStr + " random decks",
		"remaining_cards": game.Deck.RemainingCards(),
		"num_decks":       numDecks,
		"max_players":     game.MaxPlayers,
		"created":         game.Created,
	})
}

// CreateNewGlitchjackGameWithPlayers creates a new Glitchjack game with specified decks and max players.
func (h *HandlerDependencies) CreateNewGlitchjackGameWithPlayers(c *gin.Context) {
	decksStr := validators.SanitizeString(c.Param("decks"), 10)
	playersStr := validators.SanitizeString(c.Param("players"), 10)
	
	numDecks, valid := validators.ValidateNumber(decksStr)
	if !valid || numDecks <= 0 || numDecks > 100 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid decks parameter (must be 1-100)",
		})
		return
	}
	
	maxPlayers, valid := validators.ValidateNumber(playersStr)
	if !valid || maxPlayers <= 0 || maxPlayers > 10 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid players parameter (must be 1-10)",
		})
		return
	}
	
	game := h.GlitchjackService.CreateGlitchjackGameWithOptions(numDecks, maxPlayers)
	
	c.JSON(http.StatusOK, gin.H{
		"game_id":         game.ID,
		"game_type":       game.GameType.String(),
		"deck_name":       game.Deck.Name,
		"message":         "New Glitchjack game created",
		"remaining_cards": game.Deck.RemainingCards(),
		"num_decks":       numDecks,
		"max_players":     maxPlayers,
		"created":         game.Created,
	})
}

// StartGlitchjackGame begins a Glitchjack game by dealing initial cards to all players.
func (h *HandlerDependencies) StartGlitchjackGame(c *gin.Context) {
	gameID := validators.SanitizeString(c.Param("gameId"), 50)
	if !validators.ValidateUUID(gameID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid game ID format",
		})
		return
	}
	
	game, success, message := h.GlitchjackService.StartGlitchjackGame(gameID)
	if !success {
		status := http.StatusBadRequest
		if message == "Game not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{
			"error": message,
		})
		return
	}
	
	baseURL := config.GetBaseURL(c)
	playersWithCards := convertPlayersWithImages(game.Players, baseURL)
	dealerWithCards := convertDealerInfo(game.Dealer, baseURL)
	
	c.JSON(http.StatusOK, gin.H{
		"game_id":         game.ID,
		"game_type":       game.GameType.String(),
		"status":          game.Status.String(),
		"current_player":  game.CurrentPlayer,
		"players":         playersWithCards,
		"dealer":          dealerWithCards,
		"message":         message,
		"remaining_cards": game.Deck.RemainingCards(),
	})
}

// GlitchjackHit handles a player taking another card in Glitchjack.
func (h *HandlerDependencies) GlitchjackHit(c *gin.Context) {
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
	
	game, player, success, errorMsg := h.GlitchjackService.PlayerHit(gameID, playerID)
	if !success {
		status := http.StatusBadRequest
		if errorMsg == "Game not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{
			"error": errorMsg,
		})
		return
	}
	
	baseURL := config.GetBaseURL(c)
	playerInfo := convertPlayerWithImages(player, baseURL)
	
	handValue := models.CalculateGlitchjackHand(player.Hand)
	response := gin.H{
		"game_id":        game.ID,
		"player":         playerInfo,
		"hand_value":     handValue,
		"busted":         player.Busted,
		"current_player": game.CurrentPlayer,
		"game_status":    game.Status.String(),
	}
	
	if player.Busted {
		response["message"] = "Player busted!"
	} else {
		response["message"] = "Card dealt"
	}
	
	// If game is finished, include dealer info
	if game.Status == models.GameFinished {
		response["dealer"] = convertDealerInfo(game.Dealer, baseURL)
	}
	
	c.JSON(http.StatusOK, response)
}

// GlitchjackStand handles a player choosing to stand in Glitchjack.
func (h *HandlerDependencies) GlitchjackStand(c *gin.Context) {
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
	
	game, player, success, errorMsg := h.GlitchjackService.PlayerStand(gameID, playerID)
	if !success {
		status := http.StatusBadRequest
		if errorMsg == "Game not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{
			"error": errorMsg,
		})
		return
	}
	
	baseURL := config.GetBaseURL(c)
	response := gin.H{
		"game_id":        game.ID,
		"player_id":      player.ID,
		"message":        "Player stands",
		"current_player": game.CurrentPlayer,
		"game_status":    game.Status.String(),
	}
	
	// If game is finished, include dealer info
	if game.Status == models.GameFinished {
		response["dealer"] = convertDealerInfo(game.Dealer, baseURL)
		response["message"] = "All players finished, dealer's turn complete"
	}
	
	c.JSON(http.StatusOK, response)
}

// GetGlitchjackResults returns the final results of a completed Glitchjack game.
func (h *HandlerDependencies) GetGlitchjackResults(c *gin.Context) {
	gameID := validators.SanitizeString(c.Param("gameId"), 50)
	if !validators.ValidateUUID(gameID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid game ID format",
		})
		return
	}
	
	game, results, success := h.GlitchjackService.GetGlitchjackResults(gameID)
	if !success {
		status := http.StatusBadRequest
		if game == nil {
			status = http.StatusNotFound
			c.JSON(status, gin.H{
				"error": "Game not found",
			})
		} else if game.GameType != models.Glitchjack {
			c.JSON(status, gin.H{
				"error": "Not a Glitchjack game",
			})
		} else {
			c.JSON(status, gin.H{
				"error": "Game not finished",
			})
		}
		return
	}
	
	baseURL := config.GetBaseURL(c)
	playersWithResults := make([]gin.H, 0, len(game.Players))
	
	for _, player := range game.Players {
		playerInfo := convertPlayerWithImages(player, baseURL)
		playerData := gin.H{
			"player":     playerInfo,
			"hand_value": models.CalculateGlitchjackHand(player.Hand),
			"result":     results[player.ID].String(),
		}
		playersWithResults = append(playersWithResults, playerData)
	}
	
	dealerInfo := convertDealerInfo(game.Dealer, baseURL)
	
	c.JSON(http.StatusOK, gin.H{
		"game_id":      game.ID,
		"game_type":    game.GameType.String(),
		"status":       game.Status.String(),
		"players":      playersWithResults,
		"dealer":       dealerInfo,
		"dealer_value": models.CalculateGlitchjackHand(game.Dealer.Hand),
	})
}