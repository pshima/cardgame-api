package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var gameManager *GameManager
var customDeckManager *CustomDeckManager

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

// validateDeckName checks if deck name is valid (1-128 characters)
func validateDeckName(name string) bool {
	return len(name) >= 1 && len(name) <= 128
}

// validateCardIndex checks if card index is valid
func validateCardIndex(indexStr string) (int, bool) {
	index, valid := validateNumber(indexStr)
	return index, valid
}

// getTrustedProxies returns the list of trusted proxy IPs from environment or defaults
func getTrustedProxies() []string {
	// Check if TRUSTED_PROXIES environment variable is set
	if envProxies := os.Getenv("TRUSTED_PROXIES"); envProxies != "" {
		// Split comma-separated proxy IPs
		proxies := strings.Split(envProxies, ",")
		for i, proxy := range proxies {
			proxies[i] = strings.TrimSpace(proxy)
		}
		return proxies
	}
	
	// Default trusted proxies for development
	return []string{
		"127.0.0.1", // localhost
		"::1",       // localhost IPv6
		// In production, set TRUSTED_PROXIES environment variable with your actual proxy IPs
		// Examples:
		// - Load balancer IP: "10.0.1.100"
		// - Private network range: "10.0.0.0/8"
		// - Cloudflare (if using): Use Cloudflare's IP ranges
		// - Google Cloud: Use Google's load balancer IP ranges
	}
}

// getBaseURL extracts the base URL from the request with security considerations
func getBaseURL(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	
	// Only trust proxy headers if the client IP is from a trusted proxy
	clientIP := c.ClientIP()
	
	// Check if request is from a trusted proxy by comparing with our trusted proxy list
	// Note: This is a simplified check. In production, you might want more sophisticated validation
	isTrustedProxy := clientIP == "127.0.0.1" || clientIP == "::1"
	
	if isTrustedProxy {
		// Only use forwarded headers from trusted proxies
		if proto := c.GetHeader("X-Forwarded-Proto"); proto == "https" || proto == "http" {
			scheme = proto
		}
	}
	
	// Always use the original Host header for security, don't trust X-Forwarded-Host
	// unless you have specific requirements and trusted proxy configuration
	host := c.Request.Host
	
	return fmt.Sprintf("%s://%s", scheme, host)
}

func main() {
	gameManager = NewGameManager()
	customDeckManager = NewCustomDeckManager()
	r := gin.Default()
	
	// Configure trusted proxies for security
	// Get trusted proxies from environment variable or use defaults
	trustedProxies := getTrustedProxies()
	
	if err := r.SetTrustedProxies(trustedProxies); err != nil {
		panic("Failed to set trusted proxies: " + err.Error())
	}
	
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
	
	// Cribbage endpoints
	r.GET("/game/new/cribbage", createNewCribbageGame)
	r.POST("/game/:gameId/cribbage/start", startCribbageGame)
	r.POST("/game/:gameId/cribbage/discard/:playerId", cribbageDiscard)
	r.POST("/game/:gameId/cribbage/play/:playerId", cribbagePlay)
	r.POST("/game/:gameId/cribbage/go/:playerId", cribbageGo)
	r.GET("/game/:gameId/cribbage/show", cribbageShow)
	r.GET("/game/:gameId/cribbage/state", getCribbageState)
	r.GET("/game/:gameId/reset", resetDeck)
	r.GET("/game/:gameId/reset/:decks", resetDeckWithDecks)
	r.GET("/game/:gameId/reset/:decks/:type", resetDeckWithType)
	r.DELETE("/game/:gameId", deleteGame)
	r.GET("/games", listGames)
	
	// Custom deck endpoints
	r.POST("/custom-decks", createCustomDeck)
	r.GET("/custom-decks", listCustomDecks)
	r.GET("/custom-decks/:deckId", getCustomDeck)
	r.DELETE("/custom-decks/:deckId", deleteCustomDeck)
	r.POST("/custom-decks/:deckId/cards", addCustomCard)
	r.GET("/custom-decks/:deckId/cards", listCustomCards)
	r.GET("/custom-decks/:deckId/cards/:cardIndex", getCustomCard)
	r.DELETE("/custom-decks/:deckId/cards/:cardIndex", deleteCustomCard)

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

// Cribbage API handlers
func createNewCribbageGame(c *gin.Context) {
	game := gameManager.CreateGameWithType(1, Standard, Cribbage, 2)
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

func startCribbageGame(c *gin.Context) {
	gameID := c.Param("gameId")
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

	err := game.StartCribbageGame()
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

func cribbageDiscard(c *gin.Context) {
	gameID := c.Param("gameId")
	playerID := c.Param("playerId")
	
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

	var request struct {
		CardIndices []int `json:"card_indices" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
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

	err := game.CribbageDiscard(playerID, request.CardIndices)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	player := game.GetPlayer(playerID)
	response := gin.H{
		"game_id":     game.ID,
		"player_id":   playerID,
		"player_name": player.Name,
		"phase":       game.CribbageState.Phase.String(),
		"message":     "Cards discarded to crib",
	}

	// If phase changed to play, include starter card
	if game.CribbageState.Phase == CribbagePlay && game.CribbageState.Starter != nil {
		baseURL := getBaseURL(c)
		response["starter"] = game.CribbageState.Starter.ToCardWithImages(baseURL)
		response["message"] = "Cards discarded, starter cut, play phase begun"
	}

	c.JSON(http.StatusOK, response)
}

func cribbagePlay(c *gin.Context) {
	gameID := c.Param("gameId")
	playerID := c.Param("playerId")
	
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

	var request struct {
		CardIndex int `json:"card_index" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
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

	err := game.CribbagePlay(playerID, request.CardIndex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	player := game.GetPlayer(playerID)
	playerIndex := -1
	for i, p := range game.Players {
		if p.ID == playerID {
			playerIndex = i
			break
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"game_id":      game.ID,
		"player_id":    playerID,
		"player_name":  player.Name,
		"play_total":   game.CribbageState.PlayTotal,
		"play_count":   game.CribbageState.PlayCount,
		"player_score": game.CribbageState.PlayerScores[playerIndex],
		"phase":        game.CribbageState.Phase.String(),
		"current_player": game.CurrentPlayer,
		"message":      "Card played",
	})
}

func cribbageGo(c *gin.Context) {
	gameID := c.Param("gameId")
	playerID := c.Param("playerId")
	
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

	err := game.CribbageGo(playerID)
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
		"play_total":     game.CribbageState.PlayTotal,
		"current_player": game.CurrentPlayer,
		"message":        player.Name + " says go",
	})
}

func cribbageShow(c *gin.Context) {
	gameID := c.Param("gameId")
	
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

	if game.CribbageState == nil || game.CribbageState.Phase != CribbageShow {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Not in show phase",
		})
		return
	}

	scores := game.CribbageShow()
	if scores == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to score hands",
		})
		return
	}

	response := gin.H{
		"game_id": game.ID,
		"scores":  scores,
		"player_scores": game.CribbageState.PlayerScores,
		"phase":   game.CribbageState.Phase.String(),
		"status":  game.Status.String(),
	}

	// Check if game is finished
	if game.Status == GameFinished {
		if winnerVal, hasWinner := scores["winner"]; hasWinner {
			if idx, ok := winnerVal.(int); ok && idx >= 0 && idx < len(game.Players) {
				response["winner"] = game.Players[idx].Name
				response["winner_id"] = game.Players[idx].ID
			}
		}
	}

	c.JSON(http.StatusOK, response)
}

func getCribbageState(c *gin.Context) {
	gameID := c.Param("gameId")
	
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

	if game.GameType != Cribbage || game.CribbageState == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Not a cribbage game",
		})
		return
	}

	baseURL := getBaseURL(c)
	
	// Convert player hands with images
	var playersWithImages []gin.H
	for i, player := range game.Players {
		var handWithImages []CardWithImages
		for _, card := range player.Hand {
			handWithImages = append(handWithImages, ToCardWithImagesPtr(card, baseURL))
		}
		
		playersWithImages = append(playersWithImages, gin.H{
			"id":         player.ID,
			"name":       player.Name,
			"hand":       handWithImages,
			"hand_size":  player.HandSize(),
			"score":      game.CribbageState.PlayerScores[i],
		})
	}

	// Convert crib cards with images
	var cribWithImages []CardWithImages
	for _, card := range game.CribbageState.Crib {
		cribWithImages = append(cribWithImages, ToCardWithImagesPtr(card, baseURL))
	}

	// Convert played cards with images
	var playedWithImages []CardWithImages
	for _, card := range game.CribbageState.PlayedCards {
		playedWithImages = append(playedWithImages, ToCardWithImagesPtr(card, baseURL))
	}

	response := gin.H{
		"game_id":        game.ID,
		"game_type":      game.GameType.String(),
		"status":         game.Status.String(),
		"phase":          game.CribbageState.Phase.String(),
		"dealer":         game.CribbageState.Dealer,
		"current_player": game.CurrentPlayer,
		"players":        playersWithImages,
		"crib":           cribWithImages,
		"crib_size":      len(game.CribbageState.Crib),
		"played_cards":   playedWithImages,
		"play_total":     game.CribbageState.PlayTotal,
		"play_count":     game.CribbageState.PlayCount,
		"player_scores":  game.CribbageState.PlayerScores,
		"game_score":     game.CribbageState.GameScore,
	}

	// Include starter if available
	if game.CribbageState.Starter != nil {
		response["starter"] = game.CribbageState.Starter.ToCardWithImages(baseURL)
	}

	c.JSON(http.StatusOK, response)
}

// Custom deck request/response structures
type CreateCustomDeckRequest struct {
	Name string `json:"name" binding:"required"`
}

type AddCustomCardRequest struct {
	Name       string            `json:"name" binding:"required"`
	Rank       interface{}       `json:"rank,omitempty"`
	Suit       string            `json:"suit,omitempty"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

// Custom deck handlers
func createCustomDeck(c *gin.Context) {
	var req CreateCustomDeckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON: " + err.Error(),
		})
		return
	}

	name := sanitizeString(req.Name, 128)
	if !validateDeckName(name) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Deck name must be 1-128 characters",
		})
		return
	}

	deck := customDeckManager.CreateDeck(name)
	c.JSON(http.StatusCreated, gin.H{
		"id":      deck.ID,
		"name":    deck.Name,
		"message": "Custom deck created successfully",
		"created": deck.Created,
	})
}

func listCustomDecks(c *gin.Context) {
	decks := customDeckManager.ListDecks()
	
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

func getCustomDeck(c *gin.Context) {
	deckID := sanitizeString(c.Param("deckId"), 50)
	if !validateUUID(deckID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid deck ID format",
		})
		return
	}

	deck, exists := customDeckManager.GetDeck(deckID)
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

func deleteCustomDeck(c *gin.Context) {
	deckID := sanitizeString(c.Param("deckId"), 50)
	if !validateUUID(deckID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid deck ID format",
		})
		return
	}

	deleted := customDeckManager.DeleteDeck(deckID)
	if !deleted {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Custom deck not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Custom deck deleted successfully",
	})
}

func addCustomCard(c *gin.Context) {
	deckID := sanitizeString(c.Param("deckId"), 50)
	if !validateUUID(deckID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid deck ID format",
		})
		return
	}

	deck, exists := customDeckManager.GetDeck(deckID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Custom deck not found",
		})
		return
	}

	if deck.CardCount() >= 2000 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Maximum card limit (2000) reached for this deck",
		})
		return
	}

	var req AddCustomCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON: " + err.Error(),
		})
		return
	}

	name := sanitizeString(req.Name, 100)
	if len(name) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Card name cannot be empty",
		})
		return
	}

	suit := sanitizeString(req.Suit, 50)
	
	if req.Attributes != nil && len(req.Attributes) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Maximum 100 attributes allowed per card",
		})
		return
	}

	sanitizedAttributes := make(map[string]string)
	if req.Attributes != nil {
		for k, v := range req.Attributes {
			cleanKey := sanitizeString(k, 50)
			cleanValue := sanitizeString(v, 200)
			if len(cleanKey) > 0 {
				sanitizedAttributes[cleanKey] = cleanValue
			}
		}
	}

	card := deck.AddCard(name, req.Rank, suit, sanitizedAttributes)
	
	c.JSON(http.StatusCreated, gin.H{
		"index":           card.Index,
		"name":            card.Name,
		"rank":            card.Rank,
		"suit":            card.Suit,
		"game_compatible": card.GameCompatible,
		"attributes":      card.Attributes,
		"message":         "Card added successfully",
	})
}

func listCustomCards(c *gin.Context) {
	deckID := sanitizeString(c.Param("deckId"), 50)
	if !validateUUID(deckID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid deck ID format",
		})
		return
	}

	deck, exists := customDeckManager.GetDeck(deckID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Custom deck not found",
		})
		return
	}

	includeDeleted := c.Query("include_deleted") == "true"
	cards := deck.ListCards(includeDeleted)

	c.JSON(http.StatusOK, gin.H{
		"deck_id":    deck.ID,
		"deck_name":  deck.Name,
		"cards":      cards,
		"card_count": len(cards),
	})
}

func getCustomCard(c *gin.Context) {
	deckID := sanitizeString(c.Param("deckId"), 50)
	if !validateUUID(deckID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid deck ID format",
		})
		return
	}

	cardIndexStr := sanitizeString(c.Param("cardIndex"), 10)
	cardIndex, valid := validateCardIndex(cardIndexStr)
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid card index",
		})
		return
	}

	deck, exists := customDeckManager.GetDeck(deckID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Custom deck not found",
		})
		return
	}

	card := deck.GetCard(cardIndex)
	if card == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Card not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"deck_id":         deck.ID,
		"deck_name":       deck.Name,
		"index":           card.Index,
		"name":            card.Name,
		"rank":            card.Rank,
		"suit":            card.Suit,
		"game_compatible": card.GameCompatible,
		"attributes":      card.Attributes,
		"deleted":         card.Deleted,
	})
}

func deleteCustomCard(c *gin.Context) {
	deckID := sanitizeString(c.Param("deckId"), 50)
	if !validateUUID(deckID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid deck ID format",
		})
		return
	}

	cardIndexStr := sanitizeString(c.Param("cardIndex"), 10)
	cardIndex, valid := validateCardIndex(cardIndexStr)
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid card index",
		})
		return
	}

	deck, exists := customDeckManager.GetDeck(deckID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Custom deck not found",
		})
		return
	}

	deleted := deck.DeleteCard(cardIndex)
	if !deleted {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Card not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Card marked as deleted (tombstone delete)",
	})
}