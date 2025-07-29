package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var gameManager *GameManager

// Security patterns for input validation
var (
	// UUID pattern for gameID and playerID validation
	uuidPattern = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	// Alphanumeric with limited special chars for pile IDs
	pileIDPattern = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,50}$`)
	// Number pattern for numeric parameters
	numberPattern = regexp.MustCompile(`^[0-9]+$`)
	// Deck type pattern
	deckTypePattern = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,20}$`)
	// Boolean pattern for faceUp parameter
	boolPattern = regexp.MustCompile(`^(true|false|1|0)$`)
)

// validateUUID checks if the input is a valid UUID format
func validateUUID(input string) bool {
	return uuidPattern.MatchString(input)
}

// validatePlayerID checks if the input is a valid player ID (UUID or "dealer")
func validatePlayerID(input string) bool {
	return input == "dealer" || uuidPattern.MatchString(input)
}

// validatePileID checks if pile ID is safe
func validatePileID(input string) bool {
	return pileIDPattern.MatchString(input)
}

// validateNumber checks if input is a valid positive integer
func validateNumber(input string) (int, bool) {
	if !numberPattern.MatchString(input) {
		return 0, false
	}
	num, err := strconv.Atoi(input)
	if err != nil || num < 0 {
		return 0, false
	}
	return num, true
}

// validateDeckType checks if deck type input is safe
func validateDeckType(input string) bool {
	return deckTypePattern.MatchString(input)
}

// validateBoolean checks if input is a valid boolean representation
func validateBoolean(input string) bool {
	return boolPattern.MatchString(strings.ToLower(input))
}

// sanitizeString removes potentially dangerous characters and limits length
func sanitizeString(input string, maxLength int) string {
	// Remove control characters and limit length
	cleaned := strings.Map(func(r rune) rune {
		if r < 32 || r == 127 { // Remove control characters
			return -1
		}
		return r
	}, input)
	
	if len(cleaned) > maxLength {
		cleaned = cleaned[:maxLength]
	}
	
	return cleaned
}

// getBaseURL extracts the base URL from the request
func getBaseURL(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	
	// Try to get from X-Forwarded headers first (for proxies)
	if proto := c.GetHeader("X-Forwarded-Proto"); proto != "" {
		scheme = proto
	}
	
	host := c.Request.Host
	if forwardedHost := c.GetHeader("X-Forwarded-Host"); forwardedHost != "" {
		host = forwardedHost
	}
	
	return fmt.Sprintf("%s://%s", scheme, host)
}

func main() {
	gameManager = NewGameManager()
	r := gin.Default()
	
	// Serve static files for card images
	r.Static("/static", "./static")
	
	// Serve API documentation
	r.StaticFile("/openapi.yaml", "./openapi.yaml")
	r.StaticFile("/api-docs", "./api-docs.html")

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
	r.POST("/game/:gameId/start", startBlackjackGame)
	r.POST("/game/:gameId/hit/:playerId", playerHit)
	r.POST("/game/:gameId/stand/:playerId", playerStand)
	r.GET("/game/:gameId/results", getGameResults)
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
	decksStr := sanitizeString(c.Param("decks"), 10)
	numDecks, valid := validateNumber(decksStr)
	if !valid || numDecks <= 0 || numDecks > 100 { // Reasonable upper limit
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid decks parameter (must be 1-100)",
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
	decksStr := sanitizeString(c.Param("decks"), 10)
	typeStr := sanitizeString(c.Param("type"), 20)
	
	numDecks, valid := validateNumber(decksStr)
	if !valid || numDecks <= 0 || numDecks > 100 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid decks parameter (must be 1-100)",
		})
		return
	}
	
	if !validateDeckType(typeStr) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid deck type parameter",
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
	gameID := sanitizeString(c.Param("gameId"), 50)
	if !validateUUID(gameID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid game ID format",
		})
		return
	}
	
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
	gameID := sanitizeString(c.Param("gameId"), 50)
	if !validateUUID(gameID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid game ID format",
		})
		return
	}
	
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
	gameID := sanitizeString(c.Param("gameId"), 50)
	if !validateUUID(gameID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid game ID format",
		})
		return
	}
	
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
	
	// Default to face up for dealt cards
	card.FaceUp = true
	baseURL := getBaseURL(c)
	cardWithImages := card.ToCardWithImages(baseURL)
	
	c.JSON(http.StatusOK, gin.H{
		"game_id":        game.ID,
		"deck_name":      game.Deck.Name,
		"card":           cardWithImages,
		"remaining_cards": game.Deck.RemainingCards(),
	})
}

func dealCards(c *gin.Context) {
	gameID := sanitizeString(c.Param("gameId"), 50)
	countStr := sanitizeString(c.Param("count"), 10)
	
	if !validateUUID(gameID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid game ID format",
		})
		return
	}
	
	game, exists := gameManager.GetGame(gameID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Game not found",
		})
		return
	}

	count, valid := validateNumber(countStr)
	if !valid || count <= 0 || count > 52 { // Reasonable upper limit
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid count parameter (must be 1-52)",
		})
		return
	}

	if count > game.Deck.RemainingCards() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Not enough cards remaining in deck",
		})
		return
	}

	baseURL := getBaseURL(c)
	var cardsWithImages []CardWithImages
	for i := 0; i < count; i++ {
		card := game.Deck.Deal()
		if card == nil {
			break
		}
		// Default to face up for dealt cards
		card.FaceUp = true
		cardsWithImages = append(cardsWithImages, card.ToCardWithImages(baseURL))
	}

	c.JSON(http.StatusOK, gin.H{
		"game_id":        game.ID,
		"deck_name":      game.Deck.Name,
		"cards":          cardsWithImages,
		"cards_dealt":    len(cardsWithImages),
		"remaining_cards": game.Deck.RemainingCards(),
	})
}

func resetDeck(c *gin.Context) {
	gameID := sanitizeString(c.Param("gameId"), 50)
	if !validateUUID(gameID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid game ID format",
		})
		return
	}
	
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
	gameID := sanitizeString(c.Param("gameId"), 50)
	decksStr := sanitizeString(c.Param("decks"), 10)
	
	if !validateUUID(gameID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid game ID format",
		})
		return
	}
	
	game, exists := gameManager.GetGame(gameID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Game not found",
		})
		return
	}

	numDecks, valid := validateNumber(decksStr)
	if !valid || numDecks <= 0 || numDecks > 100 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid decks parameter (must be 1-100)",
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
	gameID := sanitizeString(c.Param("gameId"), 50)
	if !validateUUID(gameID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid game ID format",
		})
		return
	}
	
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
	gameID := sanitizeString(c.Param("gameId"), 50)
	decksStr := sanitizeString(c.Param("decks"), 10)
	typeStr := sanitizeString(c.Param("type"), 20)
	
	if !validateUUID(gameID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid game ID format",
		})
		return
	}
	
	if !validateDeckType(typeStr) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid deck type parameter",
		})
		return
	}
	
	game, exists := gameManager.GetGame(gameID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Game not found",
		})
		return
	}

	numDecks, valid := validateNumber(decksStr)
	if !valid || numDecks <= 0 || numDecks > 100 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid decks parameter (must be 1-100)",
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
	decksStr := sanitizeString(c.Param("decks"), 10)
	typeStr := sanitizeString(c.Param("type"), 20)
	playersStr := sanitizeString(c.Param("players"), 10)
	
	numDecks, valid := validateNumber(decksStr)
	if !valid || numDecks <= 0 || numDecks > 100 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid decks parameter (must be 1-100)",
		})
		return
	}
	
	if !validateDeckType(typeStr) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid deck type parameter",
		})
		return
	}
	
	maxPlayers, valid := validateNumber(playersStr)
	if !valid || maxPlayers <= 0 || maxPlayers > 10 {
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
	gameID := sanitizeString(c.Param("gameId"), 50)
	if !validateUUID(gameID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid game ID format",
		})
		return
	}
	
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

	// Add hand values for blackjack
	baseURL := getBaseURL(c)
	var playersWithValues []gin.H
	for _, player := range game.Players {
		handValue, hasBlackjack := player.BlackjackHandValue()
		
		// Convert cards to include images
		var handWithImages []CardWithImages
		for _, card := range player.Hand {
			handWithImages = append(handWithImages, ToCardWithImagesPtr(card, baseURL))
		}
		
		playersWithValues = append(playersWithValues, gin.H{
			"id":            player.ID,
			"name":          player.Name,
			"hand":          handWithImages,
			"hand_size":     player.HandSize(),
			"hand_value":    handValue,
			"has_blackjack": hasBlackjack,
			"is_busted":     player.IsBusted(),
		})
	}
	
	// Convert dealer cards with images
	var dealerHandWithImages []CardWithImages
	for _, card := range game.Dealer.Hand {
		dealerHandWithImages = append(dealerHandWithImages, ToCardWithImagesPtr(card, baseURL))
	}
	
	dealerValue, dealerBlackjack := game.Dealer.BlackjackHandValue()
	dealerInfo := gin.H{
		"id":            game.Dealer.ID,
		"name":          game.Dealer.Name,
		"hand":          dealerHandWithImages,
		"hand_size":     game.Dealer.HandSize(),
		"hand_value":    dealerValue,
		"has_blackjack": dealerBlackjack,
		"is_busted":     game.Dealer.IsBusted(),
	}

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

func addPlayer(c *gin.Context) {
	gameID := sanitizeString(c.Param("gameId"), 50)
	if !validateUUID(gameID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid game ID format",
		})
		return
	}
	
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
	
	// Sanitize and validate player name
	request.Name = sanitizeString(request.Name, 50)
	if len(strings.TrimSpace(request.Name)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Player name cannot be empty",
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
	gameID := sanitizeString(c.Param("gameId"), 50)
	playerID := sanitizeString(c.Param("playerId"), 50)
	
	if !validateUUID(gameID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid game ID format",
		})
		return
	}
	
	if !validatePlayerID(playerID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid player ID format",
		})
		return
	}
	
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
	gameID := sanitizeString(c.Param("gameId"), 50)
	playerID := sanitizeString(c.Param("playerId"), 50)
	
	if !validateUUID(gameID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid game ID format",
		})
		return
	}
	
	if !validatePlayerID(playerID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid player ID format",
		})
		return
	}
	
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

	baseURL := getBaseURL(c)
	cardWithImages := ToCardWithImagesPtr(card, baseURL)
	
	player := game.GetPlayer(playerID)
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

func dealToPlayerFaceUp(c *gin.Context) {
	gameID := sanitizeString(c.Param("gameId"), 50)
	playerID := sanitizeString(c.Param("playerId"), 50)
	faceUpStr := sanitizeString(c.Param("faceUp"), 10)
	
	if !validateUUID(gameID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid game ID format",
		})
		return
	}
	
	if !validatePlayerID(playerID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid player ID format",
		})
		return
	}
	
	if !validateBoolean(faceUpStr) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid faceUp parameter (must be true/false/1/0)",
		})
		return
	}
	
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

	baseURL := getBaseURL(c)
	cardWithImages := ToCardWithImagesPtr(card, baseURL)
	
	player := game.GetPlayer(playerID)
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

func discardToCard(c *gin.Context) {
	gameID := sanitizeString(c.Param("gameId"), 50)
	pileID := sanitizeString(c.Param("pileId"), 50)
	
	if !validateUUID(gameID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid game ID format",
		})
		return
	}
	
	if !validatePileID(pileID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid pile ID format",
		})
		return
	}
	
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
	
	// Sanitize and validate player ID from request body
	request.PlayerID = sanitizeString(request.PlayerID, 50)
	if !validatePlayerID(request.PlayerID) {
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
	
	baseURL := getBaseURL(c)
	cardWithImages := ToCardWithImagesPtr(card, baseURL)

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

func startBlackjackGame(c *gin.Context) {
	gameID := sanitizeString(c.Param("gameId"), 50)
	if !validateUUID(gameID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid game ID format",
		})
		return
	}
	
	game, exists := gameManager.GetGame(gameID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Game not found",
		})
		return
	}

	err := game.StartBlackjackGame()
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

func playerHit(c *gin.Context) {
	gameID := sanitizeString(c.Param("gameId"), 50)
	playerID := sanitizeString(c.Param("playerId"), 50)
	
	if !validateUUID(gameID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid game ID format",
		})
		return
	}
	
	if !validatePlayerID(playerID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid player ID format",
		})
		return
	}
	
	game, exists := gameManager.GetGame(gameID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Game not found",
		})
		return
	}

	err := game.PlayerHit(playerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	player := game.GetPlayer(playerID)
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

func playerStand(c *gin.Context) {
	gameID := sanitizeString(c.Param("gameId"), 50)
	playerID := sanitizeString(c.Param("playerId"), 50)
	
	if !validateUUID(gameID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid game ID format",
		})
		return
	}
	
	if !validatePlayerID(playerID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid player ID format",
		})
		return
	}
	
	game, exists := gameManager.GetGame(gameID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Game not found",
		})
		return
	}

	err := game.PlayerStand(playerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	player := game.GetPlayer(playerID)
	c.JSON(http.StatusOK, gin.H{
		"game_id":        game.ID,
		"player_id":      playerID,
		"player_name":    player.Name,
		"status":         game.Status.String(),
		"current_player": game.CurrentPlayer,
		"message":        player.Name + " stands",
	})
}

func getGameResults(c *gin.Context) {
	gameID := sanitizeString(c.Param("gameId"), 50)
	if !validateUUID(gameID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid game ID format",
		})
		return
	}
	
	game, exists := gameManager.GetGame(gameID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Game not found",
		})
		return
	}

	results := game.GetGameResult()
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