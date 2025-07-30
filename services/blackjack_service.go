package services

import (
	"github.com/peteshima/cardgame-api/managers"
	"github.com/peteshima/cardgame-api/models"
)

// BlackjackService provides business logic operations for blackjack games
type BlackjackService struct {
	gameManager *managers.GameManager
}

// NewBlackjackService creates a new blackjack service instance
func NewBlackjackService(gameManager *managers.GameManager) *BlackjackService {
	return &BlackjackService{
		gameManager: gameManager,
	}
}

// StartBlackjackGame starts a new blackjack game
func (bs *BlackjackService) StartBlackjackGame(gameID string) (*models.Game, error) {
	game, exists := bs.gameManager.GetGame(gameID)
	if !exists {
		return nil, nil
	}
	
	err := game.StartBlackjackGame()
	return game, err
}

// PlayerHit handles a player hitting in blackjack
func (bs *BlackjackService) PlayerHit(gameID string, playerID string) (*models.Game, *models.Player, error) {
	game, exists := bs.gameManager.GetGame(gameID)
	if !exists {
		return nil, nil, nil
	}
	
	err := game.PlayerHit(playerID)
	if err != nil {
		return game, nil, err
	}
	
	player := game.GetPlayer(playerID)
	return game, player, nil
}

// PlayerStand handles a player standing in blackjack
func (bs *BlackjackService) PlayerStand(gameID string, playerID string) (*models.Game, *models.Player, error) {
	game, exists := bs.gameManager.GetGame(gameID)
	if !exists {
		return nil, nil, nil
	}
	
	err := game.PlayerStand(playerID)
	if err != nil {
		return game, nil, err
	}
	
	player := game.GetPlayer(playerID)
	return game, player, nil
}

// GetGameResults returns the final results of a blackjack game
func (bs *BlackjackService) GetGameResults(gameID string) (*models.Game, map[string]string, bool) {
	game, exists := bs.gameManager.GetGame(gameID)
	if !exists {
		return nil, nil, false
	}
	
	results := game.GetGameResult()
	return game, results, true
}