package services

import (
	"github.com/peteshima/cardgame-api/managers"
	"github.com/peteshima/cardgame-api/models"
)

// GlitchjackService provides business logic operations for Glitchjack game management.
// Glitchjack follows the same rules as Blackjack but with randomly generated deck compositions.
type GlitchjackService struct {
	gameManager *managers.GameManager
}

// NewGlitchjackService creates a new Glitchjack service instance
func NewGlitchjackService(gameManager *managers.GameManager) *GlitchjackService {
	return &GlitchjackService{
		gameManager: gameManager,
	}
}

// CreateGlitchjackGame creates a new Glitchjack game with default settings
func (gs *GlitchjackService) CreateGlitchjackGame() *models.Game {
	return gs.CreateGlitchjackGameWithOptions(1, 6)
}

// CreateGlitchjackGameWithOptions creates a new Glitchjack game with specified options
func (gs *GlitchjackService) CreateGlitchjackGameWithOptions(numDecks int, maxPlayers int) *models.Game {
	// Create a new game with Glitchjack type using GameManager
	game := gs.gameManager.CreateGameWithType(1, models.Standard, models.Glitchjack, maxPlayers)
	
	// Replace the standard deck with Glitchjack deck(s)
	game.Deck = models.NewGlitchjackDeck()
	
	// For multiple decks, append additional random cards
	if numDecks > 1 {
		for i := 1; i < numDecks; i++ {
			additionalDeck := models.NewGlitchjackDeck()
			game.Deck.Cards = append(game.Deck.Cards, additionalDeck.Cards...)
		}
		// Shuffle the combined deck
		game.Deck.Shuffle()
	}
	
	return game
}

// StartGlitchjackGame initializes a new Glitchjack game by dealing initial cards
func (gs *GlitchjackService) StartGlitchjackGame(gameID string) (*models.Game, bool, string) {
	game, exists := gs.gameManager.GetGame(gameID)
	if !exists {
		return nil, false, "Game not found"
	}
	
	if game.GameType != models.Glitchjack {
		return game, false, "Not a Glitchjack game"
	}
	
	if game.Status != models.GameWaiting {
		return game, false, "Game already started"
	}
	
	if len(game.Players) == 0 {
		return game, false, "No players in game"
	}
	
	// Clear any existing cards and reset player state
	for _, player := range game.Players {
		player.Hand = []*models.Card{}
		player.Standing = false
		player.Busted = false
	}
	game.Dealer.Hand = []*models.Card{}
	game.Dealer.Standing = false
	game.Dealer.Busted = false
	
	// Deal initial cards (2 to each player and dealer)
	// First card to each player
	for _, player := range game.Players {
		card := game.DealToPlayer(player.ID, true)
		if card == nil {
			return game, false, "Not enough cards in deck"
		}
	}
	
	// First card to dealer (face up)
	dealerCard1 := game.Deck.Deal()
	if dealerCard1 == nil {
		return game, false, "Not enough cards for dealer"
	}
	dealerCard1.FaceUp = true
	game.Dealer.Hand = append(game.Dealer.Hand, dealerCard1)
	
	// Second card to each player
	for _, player := range game.Players {
		card := game.DealToPlayer(player.ID, true)
		if card == nil {
			return game, false, "Not enough cards in deck"
		}
	}
	
	// Second card to dealer (face down - hole card)
	dealerCard2 := game.Deck.Deal()
	if dealerCard2 == nil {
		return game, false, "Not enough cards for dealer"
	}
	dealerCard2.FaceUp = false
	game.Dealer.Hand = append(game.Dealer.Hand, dealerCard2)
	
	// Set game status
	game.Status = models.GameInProgress
	if len(game.Players) > 0 {
		game.CurrentPlayer = 0
	}
	
	return game, true, "Glitchjack game started"
}

// PlayerHit handles a player taking another card in Glitchjack
func (gs *GlitchjackService) PlayerHit(gameID string, playerID string) (*models.Game, *models.Player, bool, string) {
	game, exists := gs.gameManager.GetGame(gameID)
	if !exists {
		return nil, nil, false, "Game not found"
	}
	
	if game.GameType != models.Glitchjack {
		return game, nil, false, "Not a Glitchjack game"
	}
	
	if game.Status != models.GameInProgress {
		return game, nil, false, "Game not in progress"
	}
	
	player := game.GetPlayer(playerID)
	if player == nil {
		return game, nil, false, "Player not found"
	}
	
	if player.Standing || player.Busted {
		return game, player, false, "Player already finished"
	}
	
	// Check if it's the player's turn
	currentPlayerIndex := game.CurrentPlayer
	if currentPlayerIndex >= len(game.Players) || game.Players[currentPlayerIndex].ID != playerID {
		return game, player, false, "Not player's turn"
	}
	
	// Deal a card to the player
	card := game.DealToPlayer(playerID, true)
	if card == nil {
		return game, player, false, "No cards left in deck"
	}
	
	// Check if player busted
	handValue := models.CalculateGlitchjackHand(player.Hand)
	if handValue > 21 {
		player.Busted = true
		gs.advanceToNextPlayer(game)
	}
	
	return game, player, true, ""
}

// PlayerStand handles a player choosing to stand in Glitchjack
func (gs *GlitchjackService) PlayerStand(gameID string, playerID string) (*models.Game, *models.Player, bool, string) {
	game, exists := gs.gameManager.GetGame(gameID)
	if !exists {
		return nil, nil, false, "Game not found"
	}
	
	if game.GameType != models.Glitchjack {
		return game, nil, false, "Not a Glitchjack game"
	}
	
	if game.Status != models.GameInProgress {
		return game, nil, false, "Game not in progress"
	}
	
	player := game.GetPlayer(playerID)
	if player == nil {
		return game, nil, false, "Player not found"
	}
	
	if player.Standing || player.Busted {
		return game, player, false, "Player already finished"
	}
	
	// Check if it's the player's turn
	currentPlayerIndex := game.CurrentPlayer
	if currentPlayerIndex >= len(game.Players) || game.Players[currentPlayerIndex].ID != playerID {
		return game, player, false, "Not player's turn"
	}
	
	// Mark player as standing
	player.Standing = true
	gs.advanceToNextPlayer(game)
	
	return game, player, true, ""
}

// GetGlitchjackResults calculates and returns the results of a finished Glitchjack game
func (gs *GlitchjackService) GetGlitchjackResults(gameID string) (*models.Game, map[string]models.GlitchjackResult, bool) {
	game, exists := gs.gameManager.GetGame(gameID)
	if !exists {
		return nil, nil, false
	}
	
	if game.GameType != models.Glitchjack {
		return game, nil, false
	}
	
	if game.Status != models.GameFinished {
		return game, nil, false
	}
	
	// Calculate dealer's hand value
	dealerValue := models.CalculateGlitchjackHand(game.Dealer.Hand)
	dealerBlackjack := models.IsBlackjack(game.Dealer.Hand)
	
	results := make(map[string]models.GlitchjackResult)
	
	// Calculate results for each player
	for _, player := range game.Players {
		playerValue := models.CalculateGlitchjackHand(player.Hand)
		playerBlackjack := models.IsBlackjack(player.Hand)
		
		var result models.GlitchjackResult
		
		if player.Busted {
			result = models.GlitchjackResultBust
		} else if playerBlackjack && !dealerBlackjack {
			result = models.GlitchjackResultBlackjack
		} else if playerBlackjack && dealerBlackjack {
			result = models.GlitchjackResultPush
		} else if dealerValue > 21 {
			result = models.GlitchjackResultWin
		} else if playerValue > dealerValue {
			result = models.GlitchjackResultWin
		} else if playerValue == dealerValue {
			result = models.GlitchjackResultPush
		} else {
			result = models.GlitchjackResultLose
		}
		
		results[player.ID] = result
	}
	
	return game, results, true
}

// advanceToNextPlayer moves to the next player or dealer's turn
func (gs *GlitchjackService) advanceToNextPlayer(game *models.Game) {
	currentIndex := game.CurrentPlayer
	
	// Find next player who hasn't finished
	for i := currentIndex + 1; i < len(game.Players); i++ {
		if !game.Players[i].Standing && !game.Players[i].Busted {
			game.CurrentPlayer = i
			return
		}
	}
	
	// All players done, dealer's turn
	gs.playDealerHand(game)
}

// playDealerHand executes the dealer's turn according to Glitchjack rules
func (gs *GlitchjackService) playDealerHand(game *models.Game) {
	// Reveal dealer's hole card
	if len(game.Dealer.Hand) > 1 {
		game.Dealer.Hand[1].FaceUp = true
	}
	
	// Dealer hits on 16 or less, stands on 17 or more
	for {
		dealerValue := models.CalculateGlitchjackHand(game.Dealer.Hand)
		if dealerValue >= 17 {
			break
		}
		
		card := game.Deck.Deal()
		if card == nil {
			break
		}
		card.FaceUp = true
		game.Dealer.Hand = append(game.Dealer.Hand, card)
	}
	
	// Game is finished
	game.Status = models.GameFinished
	game.CurrentPlayer = -1
}