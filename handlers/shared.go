package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"

	"github.com/peteshima/cardgame-api/managers"
	"github.com/peteshima/cardgame-api/middleware"
	"github.com/peteshima/cardgame-api/models"
	"github.com/peteshima/cardgame-api/services"
)

// HandlerDependencies contains all the dependencies needed by handlers
type HandlerDependencies struct {
	Logger              *zap.Logger
	MetricsRegistry     *middleware.MetricsRegistry
	GameService         *services.GameService
	BlackjackService    *services.BlackjackService
	CribbageService     *services.CribbageService
	GlitchjackService   *services.GlitchjackService
	CustomDeckService   *services.CustomDeckService
	GameManager         *managers.GameManager
	CustomDeckManager   *managers.CustomDeckManager
	StartTime           time.Time
}

// NewHandlerDependencies creates a new HandlerDependencies instance
func NewHandlerDependencies(
	logger *zap.Logger,
	metricsRegistry *middleware.MetricsRegistry,
	gameManager *managers.GameManager,
	customDeckManager *managers.CustomDeckManager,
	startTime time.Time,
) *HandlerDependencies {
	return &HandlerDependencies{
		Logger:              logger,
		MetricsRegistry:     metricsRegistry,
		GameService:         services.NewGameService(gameManager),
		BlackjackService:    services.NewBlackjackService(gameManager),
		CribbageService:     services.NewCribbageService(gameManager),
		GlitchjackService:   services.NewGlitchjackService(gameManager),
		CustomDeckService:   services.NewCustomDeckService(customDeckManager),
		GameManager:         gameManager,
		CustomDeckManager:   customDeckManager,
		StartTime:           startTime,
	}
}

// GetStats provides a JSON endpoint with application metrics and health information.
// This enables monitoring and debugging by exposing key performance and business metrics.
func (h *HandlerDependencies) GetStats(c *gin.Context) {
	ctx := context.Background()
	
	// Get current metrics state
	stats := gin.H{
		"service": gin.H{
			"name":    "cardgame-api",
			"version": "1.0.0",
			"uptime":  time.Since(h.StartTime).String(),
		},
		"games": gin.H{
			"active_count":    h.GameManager.GameCount(),
			"total_created":   0, // Will be tracked via metrics
		},
		"custom_decks": gin.H{
			"active_count": len(h.CustomDeckManager.ListDecks()),
		},
		"metrics": gin.H{
			"http_requests_total":        0, // These will show current counter values
			"http_request_duration_avg":  0.0,
			"http_requests_in_flight":    0,
			"cards_dealt_total":         0,
			"api_errors_total":          0,
		},
		"system": gin.H{
			"timestamp": time.Now().Format(time.RFC3339),
		},
	}

	// Update metrics gauges
	h.MetricsRegistry.ActiveGames.Add(ctx, int64(h.GameManager.GameCount()))
	h.MetricsRegistry.ActiveCustomDecks.Add(ctx, int64(len(h.CustomDeckManager.ListDecks())))

	h.Logger.Debug("Stats endpoint accessed",
		zap.String("client_ip", c.ClientIP()),
		zap.String("user_agent", c.Request.UserAgent()),
	)

	c.JSON(http.StatusOK, stats)
}


// updateCardsDealtMetric updates the cards dealt metric
func (h *HandlerDependencies) updateCardsDealtMetric(c *gin.Context, gameID string, game *models.Game, count int) {
	h.MetricsRegistry.CardsDealt.Add(c.Request.Context(), int64(count),
		metric.WithAttributes(
			attribute.String("game_id", gameID),
			attribute.String("deck_type", game.Deck.DeckType.String()),
		),
	)
}

// updateGamesCreatedMetric updates the games created metric
func (h *HandlerDependencies) updateGamesCreatedMetric(c *gin.Context, deckType models.DeckType, deckCount int) {
	h.MetricsRegistry.GamesCreated.Add(c.Request.Context(), 1,
		metric.WithAttributes(
			attribute.String("deck_type", deckType.String()),
			attribute.Int("deck_count", deckCount),
		),
	)
}

// listDeckTypes returns all available deck types
func (h *HandlerDependencies) ListDeckTypes(c *gin.Context) {
	deckTypes := models.GetAllDeckTypes()
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

// convertCardsWithImages converts cards to include image URLs
func convertCardsWithImages(cards []*models.Card, baseURL string) []models.CardWithImages {
	var cardsWithImages []models.CardWithImages
	for _, card := range cards {
		cardsWithImages = append(cardsWithImages, models.ToCardWithImagesPtr(card, baseURL))
	}
	return cardsWithImages
}

// convertPlayersWithImages converts players to include card images and blackjack values
func convertPlayersWithImages(players []*models.Player, baseURL string) []gin.H {
	var playersWithValues []gin.H
	for _, player := range players {
		handValue, hasBlackjack := player.BlackjackHandValue()
		
		// Convert cards to include images
		handWithImages := convertCardsWithImages(player.Hand, baseURL)
		
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
	return playersWithValues
}

// convertDealerInfo converts dealer information with images and blackjack values
func convertDealerInfo(dealer *models.Player, baseURL string) gin.H {
	dealerHandWithImages := convertCardsWithImages(dealer.Hand, baseURL)
	dealerValue, dealerBlackjack := dealer.BlackjackHandValue()
	
	return gin.H{
		"id":            dealer.ID,
		"name":          dealer.Name,
		"hand":          dealerHandWithImages,
		"hand_size":     dealer.HandSize(),
		"hand_value":    dealerValue,
		"has_blackjack": dealerBlackjack,
		"is_busted":     dealer.IsBusted(),
	}
}

// convertPlayerWithImages converts a single player to include card images
func convertPlayerWithImages(player *models.Player, baseURL string) gin.H {
	handWithImages := convertCardsWithImages(player.Hand, baseURL)
	handValue, hasBlackjack := player.BlackjackHandValue()
	
	return gin.H{
		"id":            player.ID,
		"name":          player.Name,
		"hand":          handWithImages,
		"hand_size":     player.HandSize(),
		"hand_value":    handValue,
		"has_blackjack": hasBlackjack,
		"is_busted":     player.IsBusted(),
		"standing":      player.Standing,
		"busted":        player.Busted,
	}
}

// convertDiscardPiles converts discard piles to gin.H format
func convertDiscardPiles(discardPiles map[string]*models.DiscardPile) []gin.H {
	var discardInfo []gin.H
	for _, pile := range discardPiles {
		discardInfo = append(discardInfo, gin.H{
			"id":    pile.ID,
			"name":  pile.Name,
			"size":  pile.Size(),
			"cards": pile.Cards,
		})
	}
	return discardInfo
}