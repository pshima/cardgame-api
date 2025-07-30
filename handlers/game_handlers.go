package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/peteshima/cardgame-api/config"
	"github.com/peteshima/cardgame-api/models"
	"github.com/peteshima/cardgame-api/validators"
	"github.com/peteshima/cardgame-api/api"
)

// CreateNewGame creates a new card game with default settings (1 standard deck, 6 max players).
// It logs the creation, updates metrics, and returns the game details including the unique game ID.
func (h *HandlerDependencies) CreateNewGame(c *gin.Context) {
	h.Logger.Debug("Creating new game",
		zap.String("client_ip", c.ClientIP()),
		zap.Int("decks", 1),
		zap.String("type", "standard"),
	)

	game := h.GameService.CreateGame(1)
	
	// Update metrics
	h.updateGamesCreatedMetric(c, models.Standard, 1)

	h.Logger.Info("Game created successfully",
		zap.String("game_id", game.ID),
		zap.String("deck_type", game.Deck.DeckType.String()),
		zap.Int("remaining_cards", game.Deck.RemainingCards()),
		zap.String("client_ip", c.ClientIP()),
	)

	c.JSON(http.StatusOK, gin.H{
		"game_id":        game.ID,
		"deck_name":      game.Deck.Name,
		"deck_type":      game.Deck.DeckType.String(),
		"message":        "New Standard game created with 1 deck",
		"remaining_cards": game.Deck.RemainingCards(),
		"created":        game.Created,
	})
}

// CreateNewGameWithDecks creates a new game with a specified number of decks.
// It validates the deck count (1-100) and creates a standard deck game with the requested number of decks.
func (h *HandlerDependencies) CreateNewGameWithDecks(c *gin.Context) {
	decksStr := validators.SanitizeString(c.Param("decks"), 10)
	numDecks, valid := validators.ValidateNumber(decksStr)
	if !valid || numDecks <= 0 || numDecks > 100 { // Reasonable upper limit
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid decks parameter (must be 1-100)",
		})
		return
	}

	game := h.GameService.CreateGameWithDecks(numDecks)
	c.JSON(http.StatusOK, gin.H{
		"game_id":        game.ID,
		"deck_name":      game.Deck.Name,
		"deck_type":      game.Deck.DeckType.String(),
		"message":        "New Standard game created with " + decksStr + " decks",
		"remaining_cards": game.Deck.RemainingCards(),
		"num_decks":      numDecks,
		"created":        game.Created,
	})
}

// CreateNewGameWithType creates a new game with specified deck count and type (standard or spanish21).
// It validates both parameters and creates a customized game based on the provided configuration.
func (h *HandlerDependencies) CreateNewGameWithType(c *gin.Context) {
	decksStr := validators.SanitizeString(c.Param("decks"), 10)
	typeStr := validators.SanitizeString(c.Param("type"), 20)
	
	numDecks, valid := validators.ValidateNumber(decksStr)
	if !valid || numDecks <= 0 || numDecks > 100 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid decks parameter (must be 1-100)",
		})
		return
	}
	
	if !validators.ValidateDeckType(typeStr) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid deck type parameter",
		})
		return
	}
	
	deckType := models.ParseDeckType(typeStr)
	game := h.GameService.CreateGameWithType(numDecks, deckType)
	
	c.JSON(http.StatusOK, gin.H{
		"game_id":        game.ID,
		"deck_name":      game.Deck.Name,
		"deck_type":      game.Deck.DeckType.String(),
		"message":        "New " + deckType.String() + " game created with " + decksStr + " decks",
		"remaining_cards": game.Deck.RemainingCards(),
		"num_decks":      numDecks,
		"created":        game.Created,
	})
}

func (h *HandlerDependencies) CreateNewGameWithPlayers(c *gin.Context) {
	decksStr := validators.SanitizeString(c.Param("decks"), 10)
	typeStr := validators.SanitizeString(c.Param("type"), 20)
	playersStr := validators.SanitizeString(c.Param("players"), 10)
	
	numDecks, valid := validators.ValidateNumber(decksStr)
	if !valid || numDecks <= 0 || numDecks > 100 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid decks parameter (must be 1-100)",
		})
		return
	}
	
	if !validators.ValidateDeckType(typeStr) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid deck type parameter",
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
	
	deckType := models.ParseDeckType(typeStr)
	game := h.GameService.CreateGameWithAllOptions(numDecks, deckType, models.Blackjack, maxPlayers)
	
	c.JSON(http.StatusOK, gin.H{
		"game_id":        game.ID,
		"game_type":      game.GameType.String(),
		"deck_name":      game.Deck.Name,
		"deck_type":      game.Deck.DeckType.String(),
		"max_players":    game.MaxPlayers,
		"current_players": len(game.Players),
		"message":        "New " + deckType.String() + " " + game.GameType.String() + " game created",
		"remaining_cards": game.Deck.RemainingCards(),
		"created":        game.Created,
	})
}

// ShuffleDeck randomizes the order of cards in an existing game's deck.
// It validates the game ID and performs an in-place shuffle of all remaining cards.
func (h *HandlerDependencies) ShuffleDeck(c *gin.Context) {
	gameID := validators.SanitizeString(c.Param("gameId"), 50)
	if !validators.ValidateUUID(gameID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid game ID format",
		})
		return
	}
	
	game, exists := h.GameService.ShuffleGameDeck(gameID)
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
		"message":        "Deck shuffled",
		"remaining_cards": game.Deck.RemainingCards(),
	})
}

// GetGameInfo retrieves basic information about a game including deck details and card count.
// It validates the game ID and returns deck name, type, remaining cards, and timestamps.
func (h *HandlerDependencies) GetGameInfo(c *gin.Context) {
	gameID := validators.SanitizeString(c.Param("gameId"), 50)
	if !validators.ValidateUUID(gameID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid game ID format",
		})
		return
	}
	
	game, exists := h.GameService.GetGame(gameID)
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
		"remaining_cards": game.Deck.RemainingCards(),
		"is_empty":       game.Deck.IsEmpty(),
		"created":        game.Created,
		"last_used":      game.LastUsed,
		"cards":          game.Deck.Cards,
	})
}

// GetGameState retrieves complete game state with blackjack values and card images
func (h *HandlerDependencies) GetGameState(c *gin.Context) {
	gameID := validators.SanitizeString(c.Param("gameId"), 50)
	if !validators.ValidateUUID(gameID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid game ID format",
		})
		return
	}
	
	game, exists := h.GameService.GetGame(gameID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Game not found",
		})
		return
	}

	baseURL := config.GetBaseURL(c)
	discardInfo := convertDiscardPiles(game.DiscardPiles)
	playersWithValues := convertPlayersWithImages(game.Players, baseURL)
	dealerInfo := convertDealerInfo(game.Dealer, baseURL)

	c.JSON(http.StatusOK, gin.H{
		"game_id":         game.ID,
		"game_type":       game.GameType.String(),
		"status":          game.Status.String(),
		"current_player":  game.CurrentPlayer,
		"deck_name":       game.Deck.Name,
		"deck_type":       game.Deck.DeckType.String(),
		"remaining_cards": game.Deck.RemainingCards(),
		"max_players":     game.MaxPlayers,
		"current_players": len(game.Players),
		"players":         playersWithValues,
		"dealer":          dealerInfo,
		"discard_piles":   discardInfo,
		"created":         game.Created,
		"last_used":       game.LastUsed,
	})
}

func (h *HandlerDependencies) AddPlayer(c *gin.Context) {
	gameID := validators.SanitizeString(c.Param("gameId"), 50)
	if !validators.ValidateUUID(gameID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid game ID format",
		})
		return
	}
	
	var request api.AddPlayerRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}
	
	// Sanitize and validate player name
	request.Name = validators.SanitizeString(request.Name, 50)
	if len(strings.TrimSpace(request.Name)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Player name cannot be empty",
		})
		return
	}

	game, player, success := h.GameService.AddPlayerToGame(gameID, request.Name)
	if !success {
		if game == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Game not found",
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Game is full",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"game_id":     game.ID,
		"player":      player,
		"message":     "Player added successfully",
	})
}

func (h *HandlerDependencies) RemovePlayer(c *gin.Context) {
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
	
	game, removed := h.GameService.RemovePlayerFromGame(gameID, playerID)
	if game == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Game not found",
		})
		return
	}
	
	if !removed {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Player not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"game_id":   game.ID,
		"player_id": playerID,
		"message":   "Player removed successfully",
	})
}

func (h *HandlerDependencies) ListGames(c *gin.Context) {
	gameIDs := h.GameService.ListGames()
	c.JSON(http.StatusOK, gin.H{
		"games":      gameIDs,
		"game_count": h.GameService.GetGameCount(),
	})
}

func (h *HandlerDependencies) DeleteGame(c *gin.Context) {
	gameID := validators.SanitizeString(c.Param("gameId"), 50)
	if !validators.ValidateUUID(gameID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid game ID format",
		})
		return
	}
	
	deleted := h.GameService.DeleteGame(gameID)
	if !deleted {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Game not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Game deleted successfully",
		"game_id": gameID,
	})
}