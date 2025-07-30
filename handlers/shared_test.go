package handlers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/peteshima/cardgame-api/managers"
	"github.com/peteshima/cardgame-api/middleware"
)

func TestNewHandlerDependencies(t *testing.T) {
	// Create test dependencies
	logger := zap.NewNop()
	metricsRegistry := &middleware.MetricsRegistry{}
	gameManager := managers.NewGameManager()
	customDeckManager := managers.NewCustomDeckManager()
	startTime := time.Now()

	// Create handler dependencies
	deps := NewHandlerDependencies(
		logger,
		metricsRegistry,
		gameManager,
		customDeckManager,
		startTime,
	)

	// Verify all fields are set correctly
	assert.NotNil(t, deps)
	assert.Equal(t, logger, deps.Logger)
	assert.Equal(t, metricsRegistry, deps.MetricsRegistry)
	assert.Equal(t, gameManager, deps.GameManager)
	assert.Equal(t, customDeckManager, deps.CustomDeckManager)
	assert.Equal(t, startTime, deps.StartTime)

	// Verify services are created
	assert.NotNil(t, deps.GameService)
	assert.NotNil(t, deps.BlackjackService)
	assert.NotNil(t, deps.CribbageService)
	assert.NotNil(t, deps.CustomDeckService)
}