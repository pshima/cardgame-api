package services

import (
	"github.com/peteshima/cardgame-api/managers"
	"github.com/peteshima/cardgame-api/models"
)

// GameService provides business logic operations for game management
type GameService struct {
	gameManager *managers.GameManager
}

// NewGameService creates a new game service instance
func NewGameService(gameManager *managers.GameManager) *GameService {
	return &GameService{
		gameManager: gameManager,
	}
}

// CreateGame creates a new standard game
func (gs *GameService) CreateGame(numDecks int) *models.Game {
	return gs.gameManager.CreateGame(numDecks)
}

// CreateGameWithDecks creates a game with specified deck count
func (gs *GameService) CreateGameWithDecks(numDecks int) *models.Game {
	return gs.gameManager.CreateCustomGame(numDecks, models.Standard)
}

// CreateGameWithType creates a game with specified deck count and type
func (gs *GameService) CreateGameWithType(numDecks int, deckType models.DeckType) *models.Game {
	return gs.gameManager.CreateCustomGame(numDecks, deckType)
}

// CreateGameWithAllOptions creates a game with all options specified
func (gs *GameService) CreateGameWithAllOptions(numDecks int, deckType models.DeckType, gameType models.GameType, maxPlayers int) *models.Game {
	return gs.gameManager.CreateGameWithType(numDecks, deckType, gameType, maxPlayers)
}

// GetGame retrieves a game by ID
func (gs *GameService) GetGame(gameID string) (*models.Game, bool) {
	return gs.gameManager.GetGame(gameID)
}

// DeleteGame removes a game
func (gs *GameService) DeleteGame(gameID string) bool {
	return gs.gameManager.DeleteGame(gameID)
}

// ListGames returns all active games
func (gs *GameService) ListGames() []string {
	return gs.gameManager.ListGames()
}

// GetGameCount returns the number of active games
func (gs *GameService) GetGameCount() int {
	return gs.gameManager.GameCount()
}

// ShuffleGameDeck shuffles the deck for a game
func (gs *GameService) ShuffleGameDeck(gameID string) (*models.Game, bool) {
	game, exists := gs.gameManager.GetGame(gameID)
	if !exists {
		return nil, false
	}
	
	game.Deck.Shuffle()
	return game, true
}

// ResetGameDeck resets a game's deck
func (gs *GameService) ResetGameDeck(gameID string) (*models.Game, bool) {
	game, exists := gs.gameManager.GetGame(gameID)
	if !exists {
		return nil, false
	}
	
	game.Deck.Reset()
	return game, true
}

// ResetGameDeckWithDecks resets a game's deck with specified number of decks
func (gs *GameService) ResetGameDeckWithDecks(gameID string, numDecks int) (*models.Game, bool) {
	game, exists := gs.gameManager.GetGame(gameID)
	if !exists {
		return nil, false
	}
	
	game.Deck.ResetWithDecks(numDecks)
	return game, true
}

// ResetGameDeckWithType resets a game's deck with specified decks and type
func (gs *GameService) ResetGameDeckWithType(gameID string, numDecks int, deckType models.DeckType) (*models.Game, bool) {
	game, exists := gs.gameManager.GetGame(gameID)
	if !exists {
		return nil, false
	}
	
	game.Deck.ResetWithDecksAndType(numDecks, deckType)
	return game, true
}

// AddPlayerToGame adds a player to a game
func (gs *GameService) AddPlayerToGame(gameID string, playerName string) (*models.Game, *models.Player, bool) {
	game, exists := gs.gameManager.GetGame(gameID)
	if !exists {
		return nil, nil, false
	}
	
	player := game.AddPlayer(playerName)
	if player == nil {
		return game, nil, false
	}
	
	return game, player, true
}

// RemovePlayerFromGame removes a player from a game
func (gs *GameService) RemovePlayerFromGame(gameID string, playerID string) (*models.Game, bool) {
	game, exists := gs.gameManager.GetGame(gameID)
	if !exists {
		return nil, false
	}
	
	removed := game.RemovePlayer(playerID)
	return game, removed
}

// DealCard deals a single card from a game
func (gs *GameService) DealCard(gameID string) (*models.Game, *models.Card, bool) {
	game, exists := gs.gameManager.GetGame(gameID)
	if !exists {
		return nil, nil, false
	}
	
	card := game.Deck.Deal()
	if card == nil {
		return game, nil, false
	}
	
	// Default to face up for dealt cards
	card.FaceUp = true
	return game, card, true
}

// DealCards deals multiple cards from a game
func (gs *GameService) DealCards(gameID string, count int) (*models.Game, []*models.Card, bool) {
	game, exists := gs.gameManager.GetGame(gameID)
	if !exists {
		return nil, nil, false
	}
	
	if count > game.Deck.RemainingCards() {
		return game, nil, false
	}
	
	var cards []*models.Card
	for i := 0; i < count; i++ {
		card := game.Deck.Deal()
		if card == nil {
			break
		}
		// Default to face up for dealt cards
		card.FaceUp = true
		cards = append(cards, card)
	}
	
	return game, cards, true
}

// DealToPlayer deals a card to a specific player
func (gs *GameService) DealToPlayer(gameID string, playerID string, faceUp bool) (*models.Game, *models.Player, *models.Card, bool) {
	game, exists := gs.gameManager.GetGame(gameID)
	if !exists {
		return nil, nil, nil, false
	}
	
	player := game.GetPlayer(playerID)
	if player == nil {
		return game, nil, nil, false
	}
	
	card := game.DealToPlayer(playerID, faceUp)
	if card == nil {
		return game, player, nil, false
	}
	
	return game, player, card, true
}

// DiscardCard handles card discarding to piles
func (gs *GameService) DiscardCard(gameID string, pileID string, playerID string, cardIndex int) (*models.Game, *models.Player, *models.DiscardPile, *models.Card, bool) {
	game, exists := gs.gameManager.GetGame(gameID)
	if !exists {
		return nil, nil, nil, nil, false
	}
	
	player := game.GetPlayer(playerID)
	if player == nil {
		return game, nil, nil, nil, false
	}
	
	pile := game.GetDiscardPile(pileID)
	if pile == nil {
		return game, player, nil, nil, false
	}
	
	card := player.RemoveCard(cardIndex)
	if card == nil {
		return game, player, pile, nil, false
	}
	
	pile.AddCard(card)
	return game, player, pile, card, true
}