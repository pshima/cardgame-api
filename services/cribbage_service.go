package services

import (
	"github.com/peteshima/cardgame-api/managers"
	"github.com/peteshima/cardgame-api/models"
)

// CribbageService provides business logic operations for cribbage games
type CribbageService struct {
	gameManager *managers.GameManager
}

// NewCribbageService creates a new cribbage service instance
func NewCribbageService(gameManager *managers.GameManager) *CribbageService {
	return &CribbageService{
		gameManager: gameManager,
	}
}

// CreateCribbageGame creates a new cribbage game
func (cs *CribbageService) CreateCribbageGame() *models.Game {
	return cs.gameManager.CreateGameWithType(1, models.Standard, models.Cribbage, 2)
}

// StartCribbageGame starts a cribbage game
func (cs *CribbageService) StartCribbageGame(gameID string) (*models.Game, error) {
	game, exists := cs.gameManager.GetGame(gameID)
	if !exists {
		return nil, nil
	}
	
	err := game.StartCribbageGame()
	return game, err
}

// CribbageDiscard handles the discard phase of cribbage
func (cs *CribbageService) CribbageDiscard(gameID string, playerID string, cardIndices []int) (*models.Game, *models.Player, error) {
	game, exists := cs.gameManager.GetGame(gameID)
	if !exists {
		return nil, nil, nil
	}
	
	err := game.CribbageDiscard(playerID, cardIndices)
	if err != nil {
		return game, nil, err
	}
	
	player := game.GetPlayer(playerID)
	return game, player, nil
}

// CribbagePlay handles the play phase of cribbage
func (cs *CribbageService) CribbagePlay(gameID string, playerID string, cardIndex int) (*models.Game, *models.Player, error) {
	game, exists := cs.gameManager.GetGame(gameID)
	if !exists {
		return nil, nil, nil
	}
	
	err := game.CribbagePlay(playerID, cardIndex)
	if err != nil {
		return game, nil, err
	}
	
	player := game.GetPlayer(playerID)
	return game, player, nil
}

// CribbageGo handles the "go" action in cribbage play
func (cs *CribbageService) CribbageGo(gameID string, playerID string) (*models.Game, *models.Player, error) {
	game, exists := cs.gameManager.GetGame(gameID)
	if !exists {
		return nil, nil, nil
	}
	
	err := game.CribbageGo(playerID)
	if err != nil {
		return game, nil, err
	}
	
	player := game.GetPlayer(playerID)
	return game, player, nil
}

// CribbageShow handles the show phase of cribbage
func (cs *CribbageService) CribbageShow(gameID string) (*models.Game, map[string]interface{}, bool) {
	game, exists := cs.gameManager.GetGame(gameID)
	if !exists {
		return nil, nil, false
	}
	
	scores := game.CribbageShow()
	if scores == nil {
		return game, nil, false
	}
	
	return game, scores, true
}