package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var gameManager *GameManager

func main() {
	gameManager = NewGameManager()
	r := gin.Default()

	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, World!",
		})
	})

	r.GET("/deck-types", listDeckTypes)
	r.GET("/game/new", createNewGame)
	r.GET("/game/new/:decks", createNewGameWithDecks)
	r.GET("/game/new/:decks/:type", createNewGameWithType)
	r.GET("/game/new/:decks/:type/:players", createNewGameWithPlayers)
	r.GET("/game/:gameId/shuffle", shuffleDeck)
	r.GET("/game/:gameId", getGameInfo)
	r.GET("/game/:gameId/state", getGameState)
	r.POST("/game/:gameId/players", addPlayer)
	r.DELETE("/game/:gameId/players/:playerId", removePlayer)
	r.GET("/game/:gameId/deal", dealCard)
	r.GET("/game/:gameId/deal/:count", dealCards)
	r.GET("/game/:gameId/deal/player/:playerId", dealToPlayer)
	r.GET("/game/:gameId/deal/player/:playerId/:faceUp", dealToPlayerFaceUp)
	r.POST("/game/:gameId/discard/:pileId", discardToCard)
	r.GET("/game/:gameId/reset", resetDeck)
	r.GET("/game/:gameId/reset/:decks", resetDeckWithDecks)
	r.GET("/game/:gameId/reset/:decks/:type", resetDeckWithType)
	r.DELETE("/game/:gameId", deleteGame)
	r.GET("/games", listGames)

	r.Run(":8080")
}

func parseDeckType(typeStr string) DeckType {
	switch strings.ToLower(typeStr) {
	case "spanish21", "spanish_21", "spanish-21":
		return Spanish21
	case "standard", "normal", "regular":
		return Standard
	default:
		return Standard
	}
}

func createNewGame(c *gin.Context) {
	game := gameManager.CreateGame(1)
	c.JSON(http.StatusOK, gin.H{
		"game_id":        game.ID,
		"deck_name":      game.Deck.Name,
		"deck_type":      game.Deck.DeckType.String(),
		"message":        "New Standard game created with 1 deck",
		"remaining_cards": game.Deck.RemainingCards(),
		"created":        game.Created,
	})
}

func createNewGameWithDecks(c *gin.Context) {
	decksStr := c.Param("decks")
	numDecks, err := strconv.Atoi(decksStr)
	if err != nil || numDecks <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid decks parameter",
		})
		return
	}

	game := gameManager.CreateGame(numDecks)
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

func createNewGameWithType(c *gin.Context) {
	decksStr := c.Param("decks")
	typeStr := c.Param("type")
	
	numDecks, err := strconv.Atoi(decksStr)
	if err != nil || numDecks <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid decks parameter",
		})
		return
	}
	
	deckType := parseDeckType(typeStr)
	game := gameManager.CreateCustomGame(numDecks, deckType)
	
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

func shuffleDeck(c *gin.Context) {
	gameID := c.Param("gameId")
	game, exists := gameManager.GetGame(gameID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Game not found",
		})
		return
	}

	game.Deck.Shuffle()
	c.JSON(http.StatusOK, gin.H{
		"game_id":        game.ID,
		"deck_name":      game.Deck.Name,
		"deck_type":      game.Deck.DeckType.String(),
		"message":        "Deck shuffled",
		"remaining_cards": game.Deck.RemainingCards(),
	})
}

func getGameInfo(c *gin.Context) {
	gameID := c.Param("gameId")
	game, exists := gameManager.GetGame(gameID)
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

func dealCard(c *gin.Context) {
	gameID := c.Param("gameId")
	game, exists := gameManager.GetGame(gameID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Game not found",
		})
		return
	}

	card := game.Deck.Deal()
	if card == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No cards remaining in deck",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"game_id":        game.ID,
		"deck_name":      game.Deck.Name,
		"card":           card,
		"remaining_cards": game.Deck.RemainingCards(),
	})
}

func dealCards(c *gin.Context) {
	gameID := c.Param("gameId")
	countStr := c.Param("count")
	
	game, exists := gameManager.GetGame(gameID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Game not found",
		})
		return
	}

	count, err := strconv.Atoi(countStr)
	if err != nil || count <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid count parameter",
		})
		return
	}

	if count > game.Deck.RemainingCards() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Not enough cards remaining in deck",
		})
		return
	}

	var cards []*Card
	for i := 0; i < count; i++ {
		card := game.Deck.Deal()
		if card == nil {
			break
		}
		cards = append(cards, card)
	}

	c.JSON(http.StatusOK, gin.H{
		"game_id":        game.ID,
		"deck_name":      game.Deck.Name,
		"cards":          cards,
		"cards_dealt":    len(cards),
		"remaining_cards": game.Deck.RemainingCards(),
	})
}

func resetDeck(c *gin.Context) {
	gameID := c.Param("gameId")
	game, exists := gameManager.GetGame(gameID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Game not found",
		})
		return
	}

	game.Deck.Reset()
	c.JSON(http.StatusOK, gin.H{
		"game_id":        game.ID,
		"deck_name":      game.Deck.Name,
		"deck_type":      game.Deck.DeckType.String(),
		"message":        "Deck reset",
		"remaining_cards": game.Deck.RemainingCards(),
	})
}

func resetDeckWithDecks(c *gin.Context) {
	gameID := c.Param("gameId")
	decksStr := c.Param("decks")
	
	game, exists := gameManager.GetGame(gameID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Game not found",
		})
		return
	}

	numDecks, err := strconv.Atoi(decksStr)
	if err != nil || numDecks <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid decks parameter",
		})
		return
	}

	game.Deck.ResetWithDecks(numDecks)
	c.JSON(http.StatusOK, gin.H{
		"game_id":        game.ID,
		"deck_name":      game.Deck.Name,
		"deck_type":      game.Deck.DeckType.String(),
		"message":        "Deck reset to " + decksStr + " decks",
		"remaining_cards": game.Deck.RemainingCards(),
		"num_decks":      numDecks,
	})
}

func deleteGame(c *gin.Context) {
	gameID := c.Param("gameId")
	deleted := gameManager.DeleteGame(gameID)
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

func resetDeckWithType(c *gin.Context) {
	gameID := c.Param("gameId")
	decksStr := c.Param("decks")
	typeStr := c.Param("type")
	
	game, exists := gameManager.GetGame(gameID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Game not found",
		})
		return
	}

	numDecks, err := strconv.Atoi(decksStr)
	if err != nil || numDecks <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid decks parameter",
		})
		return
	}
	
	deckType := parseDeckType(typeStr)
	game.Deck.ResetWithDecksAndType(numDecks, deckType)
	
	c.JSON(http.StatusOK, gin.H{
		"game_id":        game.ID,
		"deck_name":      game.Deck.Name,
		"deck_type":      game.Deck.DeckType.String(),
		"message":        "Deck reset to " + decksStr + " " + deckType.String() + " decks",
		"remaining_cards": game.Deck.RemainingCards(),
		"num_decks":      numDecks,
	})
}

func listDeckTypes(c *gin.Context) {
	deckTypes := GetAllDeckTypes()
	var deckInfo []gin.H
	
	for _, deckType := range deckTypes {
		deckInfo = append(deckInfo, gin.H{
			"id":             int(deckType),
			"type":           deckType.String(),
			"name":           deckType.String(),
			"description":    deckType.Description(),
			"cards_per_deck": deckType.CardsPerDeck(),
		})
	}
	
	c.JSON(http.StatusOK, gin.H{
		"deck_types": deckInfo,
		"count":      len(deckTypes),
	})
}

func listGames(c *gin.Context) {
	gameIDs := gameManager.ListGames()
	c.JSON(http.StatusOK, gin.H{
		"games":      gameIDs,
		"game_count": gameManager.GameCount(),
	})
}

func createNewGameWithPlayers(c *gin.Context) {
	decksStr := c.Param("decks")
	typeStr := c.Param("type")
	playersStr := c.Param("players")
	
	numDecks, err := strconv.Atoi(decksStr)
	if err != nil || numDecks <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid decks parameter",
		})
		return
	}
	
	maxPlayers, err := strconv.Atoi(playersStr)
	if err != nil || maxPlayers <= 0 || maxPlayers > 10 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid players parameter (must be 1-10)",
		})
		return
	}
	
	deckType := parseDeckType(typeStr)
	game := gameManager.CreateGameWithType(numDecks, deckType, Blackjack, maxPlayers)
	
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

func getGameState(c *gin.Context) {
	gameID := c.Param("gameId")
	game, exists := gameManager.GetGame(gameID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Game not found",
		})
		return
	}

	var discardInfo []gin.H
	for _, pile := range game.DiscardPiles {
		discardInfo = append(discardInfo, gin.H{
			"id":    pile.ID,
			"name":  pile.Name,
			"size":  pile.Size(),
			"cards": pile.Cards,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"game_id":         game.ID,
		"game_type":       game.GameType.String(),
		"deck_name":       game.Deck.Name,
		"deck_type":       game.Deck.DeckType.String(),
		"remaining_cards": game.Deck.RemainingCards(),
		"max_players":     game.MaxPlayers,
		"current_players": len(game.Players),
		"players":         game.Players,
		"dealer":          game.Dealer,
		"discard_piles":   discardInfo,
		"created":         game.Created,
		"last_used":       game.LastUsed,
	})
}

func addPlayer(c *gin.Context) {
	gameID := c.Param("gameId")
	game, exists := gameManager.GetGame(gameID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Game not found",
		})
		return
	}

	var request struct {
		Name string `json:"name" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	player := game.AddPlayer(request.Name)
	if player == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Game is full",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"game_id":     game.ID,
		"player":      player,
		"message":     "Player added successfully",
	})
}

func removePlayer(c *gin.Context) {
	gameID := c.Param("gameId")
	playerID := c.Param("playerId")
	
	game, exists := gameManager.GetGame(gameID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Game not found",
		})
		return
	}

	removed := game.RemovePlayer(playerID)
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

func dealToPlayer(c *gin.Context) {
	gameID := c.Param("gameId")
	playerID := c.Param("playerId")
	
	game, exists := gameManager.GetGame(gameID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Game not found",
		})
		return
	}

	card := game.DealToPlayer(playerID, false) // Default to face down
	if card == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to deal card (player not found or no cards remaining)",
		})
		return
	}

	player := game.GetPlayer(playerID)
	c.JSON(http.StatusOK, gin.H{
		"game_id":         game.ID,
		"player_id":       playerID,
		"player_name":     player.Name,
		"card":            card,
		"hand_size":       player.HandSize(),
		"remaining_cards": game.Deck.RemainingCards(),
		"message":         "Card dealt to " + player.Name,
	})
}

func dealToPlayerFaceUp(c *gin.Context) {
	gameID := c.Param("gameId")
	playerID := c.Param("playerId")
	faceUpStr := c.Param("faceUp")
	
	game, exists := gameManager.GetGame(gameID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Game not found",
		})
		return
	}

	faceUp := strings.ToLower(faceUpStr) == "true" || faceUpStr == "1"
	
	card := game.DealToPlayer(playerID, faceUp)
	if card == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to deal card (player not found or no cards remaining)",
		})
		return
	}

	player := game.GetPlayer(playerID)
	c.JSON(http.StatusOK, gin.H{
		"game_id":         game.ID,
		"player_id":       playerID,
		"player_name":     player.Name,
		"card":            card,
		"face_up":         faceUp,
		"hand_size":       player.HandSize(),
		"remaining_cards": game.Deck.RemainingCards(),
		"message":         "Card dealt to " + player.Name,
	})
}

func discardToCard(c *gin.Context) {
	gameID := c.Param("gameId")
	pileID := c.Param("pileId")
	
	game, exists := gameManager.GetGame(gameID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Game not found",
		})
		return
	}

	var request struct {
		PlayerID  string `json:"player_id" binding:"required"`
		CardIndex int    `json:"card_index"`
	}
	
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if request.CardIndex < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid card index",
		})
		return
	}

	player := game.GetPlayer(request.PlayerID)
	if player == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Player not found",
		})
		return
	}

	card := player.RemoveCard(request.CardIndex)
	if card == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid card index",
		})
		return
	}

	pile := game.GetDiscardPile(pileID)
	if pile == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Discard pile not found",
		})
		return
	}

	pile.AddCard(card)

	c.JSON(http.StatusOK, gin.H{
		"game_id":       game.ID,
		"player_id":     request.PlayerID,
		"player_name":   player.Name,
		"card":          card,
		"pile_id":       pileID,
		"pile_name":     pile.Name,
		"pile_size":     pile.Size(),
		"hand_size":     player.HandSize(),
		"message":       "Card discarded to " + pile.Name,
	})
}