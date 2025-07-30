package middleware

import (
	"testing"
	"context"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/metric/noop"
)

func TestNewMetricsRegistry(t *testing.T) {
	// Create a test meter
	meter := noop.NewMeterProvider().Meter("test")
	
	// Create metrics registry
	registry, err := NewMetricsRegistry(meter)
	assert.NoError(t, err)
	
	// Verify all metrics are created
	assert.NotNil(t, registry)
	assert.NotNil(t, registry.HttpRequestsTotal)
	assert.NotNil(t, registry.HttpRequestDuration)
	assert.NotNil(t, registry.HttpRequestsInFlight)
	assert.NotNil(t, registry.ActiveGames)
	assert.NotNil(t, registry.ActiveCustomDecks)
	assert.NotNil(t, registry.CardsDealt)
	assert.NotNil(t, registry.GamesCreated)
	assert.NotNil(t, registry.ApiErrors)
}

func TestMetricsRegistryCreation(t *testing.T) {
	// Create a test meter
	meter := noop.NewMeterProvider().Meter("test")
	
	// Create metrics registry
	registry, err := NewMetricsRegistry(meter)
	assert.NoError(t, err)
	ctx := context.Background()

	// Test that metrics can be used without panicking
	assert.NotPanics(t, func() {
		registry.HttpRequestsTotal.Add(ctx, 1)
		registry.HttpRequestDuration.Record(ctx, 0.5)
		registry.HttpRequestsInFlight.Add(ctx, 1)
		registry.ActiveGames.Add(ctx, 1)
		registry.CardsDealt.Add(ctx, 5)
		registry.GamesCreated.Add(ctx, 1)
		registry.ApiErrors.Add(ctx, 1)
	})
}

func TestMetricsRegistryWithNilMeter(t *testing.T) {
	// Test that creating registry with nil meter returns error
	registry, err := NewMetricsRegistry(nil)
	assert.Error(t, err)
	assert.Nil(t, registry)
}