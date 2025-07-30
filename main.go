package main

import (
	"net/http"
	"runtime"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/peteshima/cardgame-api/config"
	"github.com/peteshima/cardgame-api/handlers"
	"github.com/peteshima/cardgame-api/managers"
	"github.com/peteshima/cardgame-api/middleware"
)

// Build information set via ldflags during compilation
var (
	Version    = "dev"     // Version string (set via -ldflags "-X main.Version=...")
	BuildTime  = "unknown" // Build timestamp (set via -ldflags "-X main.BuildTime=...")
	CommitHash = "unknown" // Git commit hash (set via -ldflags "-X main.CommitHash=...")
)

// main initializes the Card Game API server with logging, metrics, and all HTTP routes.
// It sets up observability, security middleware, and starts the web server on a configurable port.
func main() {
	// Record start time
	startTime := time.Now()
	
	// Initialize logger first
	logger := config.InitLogger()
	defer logger.Sync()

	logger.Info("Starting Card Game API",
		zap.String("version", Version),
		zap.String("build_time", BuildTime),
		zap.String("commit_hash", CommitHash),
		zap.String("go_version", runtime.Version()),
		zap.String("go_os", runtime.GOOS),
		zap.String("go_arch", runtime.GOARCH),
	)

	// Initialize metrics
	_, metricsRegistry := config.InitMetrics(logger)

	// Initialize managers
	gameManager := managers.NewGameManager()
	customDeckManager := managers.NewCustomDeckManager()

	logger.Info("Managers initialized successfully")

	// Create handler dependencies
	deps := handlers.NewHandlerDependencies(
		logger, 
		metricsRegistry, 
		gameManager, 
		customDeckManager, 
		startTime,
	)

	// Create Gin router without default middleware
	r := gin.New()

	// Add custom middleware
	r.Use(middleware.LogMiddleware(logger, metricsRegistry))
	r.Use(gin.Recovery())
	
	// Configure CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3001"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	
	// Configure trusted proxies for security
	trustedProxies := config.GetTrustedProxies()
	
	if err := r.SetTrustedProxies(trustedProxies); err != nil {
		logger.Fatal("Failed to set trusted proxies", zap.Error(err))
	}

	logger.Info("Trusted proxies configured", zap.Strings("proxies", trustedProxies))
	
	// Serve static files for card images
	r.Static("/static", "./static")
	
	// Metrics endpoints
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/stats", deps.GetStats)

	// Serve API documentation
	r.StaticFile("/openapi.yaml", "./openapi.yaml")
	r.StaticFile("/api-docs", "./api-docs.html")

	logger.Info("Static file routes configured")

	// Basic routes
	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, World!",
		})
	})

	// Version and build information endpoint
	r.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"version":     Version,
			"build_time":  BuildTime,
			"commit_hash": CommitHash,
			"go_version":  runtime.Version(),
			"go_os":       runtime.GOOS,
			"go_arch":     runtime.GOARCH,
		})
	})

	// Game management routes
	r.GET("/deck-types", deps.ListDeckTypes)
	r.GET("/game/new", deps.CreateNewGame)
	r.GET("/game/new/:decks", deps.CreateNewGameWithDecks)
	r.GET("/game/new/:decks/:type", deps.CreateNewGameWithType)
	r.GET("/game/new/:decks/:type/:players", deps.CreateNewGameWithPlayers)
	r.GET("/game/:gameId/shuffle", deps.ShuffleDeck)
	r.GET("/game/:gameId", deps.GetGameInfo)
	r.GET("/game/:gameId/state", deps.GetGameState)
	r.POST("/game/:gameId/players", deps.AddPlayer)
	r.DELETE("/game/:gameId/players/:playerId", deps.RemovePlayer)
	r.GET("/games", deps.ListGames)
	r.DELETE("/game/:gameId", deps.DeleteGame)

	// Card dealing routes
	r.GET("/game/:gameId/deal", deps.DealCard)
	r.GET("/game/:gameId/deal/:count", deps.DealCards)
	r.GET("/game/:gameId/deal/player/:playerId", deps.DealToPlayer)
	r.GET("/game/:gameId/deal/player/:playerId/:faceUp", deps.DealToPlayerFaceUp)
	r.POST("/game/:gameId/discard/:pileId", deps.DiscardToCard)

	// Deck reset routes
	r.GET("/game/:gameId/reset", deps.ResetDeck)
	r.GET("/game/:gameId/reset/:decks", deps.ResetDeckWithDecks)
	r.GET("/game/:gameId/reset/:decks/:type", deps.ResetDeckWithType)

	// Blackjack routes
	r.POST("/game/:gameId/start", deps.StartBlackjackGame)
	r.POST("/game/:gameId/hit/:playerId", deps.PlayerHit)
	r.POST("/game/:gameId/stand/:playerId", deps.PlayerStand)
	r.GET("/game/:gameId/results", deps.GetGameResults)
	
	// Cribbage routes
	r.GET("/game/new/cribbage", deps.CreateNewCribbageGame)
	r.POST("/game/:gameId/cribbage/start", deps.StartCribbageGame)
	r.POST("/game/:gameId/cribbage/discard/:playerId", deps.CribbageDiscard)
	
	// Custom deck routes
	r.POST("/custom-decks", deps.CreateCustomDeck)
	r.GET("/custom-decks", deps.ListCustomDecks)
	r.GET("/custom-decks/:deckId", deps.GetCustomDeck)

	// Get port from environment variable, default to 8080
	port := config.GetPort()

	logger.Info("All routes configured, starting server",
		zap.String("port", port),
		zap.String("metrics_endpoint", "/metrics"),
		zap.String("stats_endpoint", "/stats"),
	)

	if err := r.Run(":" + port); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}