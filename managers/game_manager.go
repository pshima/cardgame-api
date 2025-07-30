package managers

import (
	"sync"
	"time"

	"github.com/peteshima/cardgame-api/models"
)

// GameManager provides thread-safe management of multiple concurrent card games.
// It uses read-write mutexes to allow concurrent read access while ensuring write safety.
type GameManager struct {
	games map[string]*models.Game
	mutex sync.RWMutex
}

// NewGameManager creates a new game manager with an empty game collection.
// This is used as a singleton to manage all active games in the application.
func NewGameManager() *GameManager {
	return &GameManager{
		games: make(map[string]*models.Game),
	}
}

func (gm *GameManager) CreateGame(numDecks int) *models.Game {
	return gm.CreateCustomGame(numDecks, models.Standard)
}

func (gm *GameManager) CreateCustomGame(numDecks int, deckType models.DeckType) *models.Game {
	return gm.CreateGameWithType(numDecks, deckType, models.Blackjack, 6)
}

func (gm *GameManager) CreateGameWithType(numDecks int, deckType models.DeckType, gameType models.GameType, maxPlayers int) *models.Game {
	gm.mutex.Lock()
	defer gm.mutex.Unlock()
	
	game := models.NewGameWithType(numDecks, deckType, gameType, maxPlayers)
	gm.games[game.ID] = game
	return game
}

func (gm *GameManager) GetGame(gameID string) (*models.Game, bool) {
	gm.mutex.RLock()
	defer gm.mutex.RUnlock()
	
	game, exists := gm.games[gameID]
	if exists {
		game.UpdateLastUsed()
	}
	return game, exists
}

func (gm *GameManager) DeleteGame(gameID string) bool {
	gm.mutex.Lock()
	defer gm.mutex.Unlock()
	
	_, exists := gm.games[gameID]
	if exists {
		delete(gm.games, gameID)
	}
	return exists
}

// ListGames returns a slice of all active game IDs.
// This method is thread-safe and provides a snapshot of current games.
func (gm *GameManager) ListGames() []string {
	gm.mutex.RLock()
	defer gm.mutex.RUnlock()
	
	gameIDs := make([]string, 0, len(gm.games))
	for id := range gm.games {
		gameIDs = append(gameIDs, id)
	}
	return gameIDs
}

// CleanupOldGames removes games that haven't been used within the specified duration.
// Returns the number of games deleted, used for memory management and cleanup.
func (gm *GameManager) CleanupOldGames(maxAge time.Duration) int {
	gm.mutex.Lock()
	defer gm.mutex.Unlock()
	
	cutoff := time.Now().Add(-maxAge)
	deleted := 0
	
	for id, game := range gm.games {
		if game.LastUsed.Before(cutoff) {
			delete(gm.games, id)
			deleted++
		}
	}
	
	return deleted
}

// GameCount returns the current number of active games.
// This method is thread-safe and used for monitoring and metrics.
func (gm *GameManager) GameCount() int {
	gm.mutex.RLock()
	defer gm.mutex.RUnlock()
	return len(gm.games)
}