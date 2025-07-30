package middleware

import (
	"go.opentelemetry.io/otel/metric"
)

// NewMetricsRegistry creates and initializes all application metrics with proper descriptions.
// It returns an error if any metric fails to initialize, ensuring all metrics are properly configured.
func NewMetricsRegistry(meter metric.Meter) (*MetricsRegistry, error) {
	httpRequestsTotal, err := meter.Int64Counter(
		"http_requests_total",
		metric.WithDescription("Total number of HTTP requests"),
	)
	if err != nil {
		return nil, err
	}

	httpRequestDuration, err := meter.Float64Histogram(
		"http_request_duration_seconds",
		metric.WithDescription("HTTP request duration in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	httpRequestsInFlight, err := meter.Int64UpDownCounter(
		"http_requests_in_flight",
		metric.WithDescription("Current number of HTTP requests being processed"),
	)
	if err != nil {
		return nil, err
	}

	activeGames, err := meter.Int64UpDownCounter(
		"active_games",
		metric.WithDescription("Current number of active games"),
	)
	if err != nil {
		return nil, err
	}

	activeCustomDecks, err := meter.Int64UpDownCounter(
		"active_custom_decks",
		metric.WithDescription("Current number of custom decks"),
	)
	if err != nil {
		return nil, err
	}

	cardsDealt, err := meter.Int64Counter(
		"cards_dealt_total",
		metric.WithDescription("Total number of cards dealt"),
	)
	if err != nil {
		return nil, err
	}

	gamesCreated, err := meter.Int64Counter(
		"games_created_total",
		metric.WithDescription("Total number of games created"),
	)
	if err != nil {
		return nil, err
	}

	apiErrors, err := meter.Int64Counter(
		"api_errors_total",
		metric.WithDescription("Total number of API errors"),
	)
	if err != nil {
		return nil, err
	}

	return &MetricsRegistry{
		HttpRequestsTotal:     httpRequestsTotal,
		HttpRequestDuration:   httpRequestDuration,
		HttpRequestsInFlight:  httpRequestsInFlight,
		ActiveGames:          activeGames,
		ActiveCustomDecks:    activeCustomDecks,
		CardsDealt:           cardsDealt,
		GamesCreated:         gamesCreated,
		ApiErrors:            apiErrors,
	}, nil
}