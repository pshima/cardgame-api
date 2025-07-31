package models

import (
	"time"

	"github.com/google/uuid"
)

// GameStatus represents the current state of a game session.
// Used to control game flow and determine valid actions.
type GameStatus int

const (
	GameWaiting GameStatus = iota
	GameInProgress
	GameFinished
)

// String returns the string representation of the game status.
// This is used in JSON responses and logging to show human-readable status.
func (gs GameStatus) String() string {
	switch gs {
	case GameWaiting:
		return "waiting"
	case GameInProgress:
		return "in_progress"
	case GameFinished:
		return "finished"
	default:
		return "waiting"
	}
}

// Game represents a complete card game session with players, deck, and game state.
// It supports multiple game types (Blackjack, Cribbage) and manages all game operations.
type Game struct {
	ID           string                  `json:"id"`
	GameType     GameType                `json:"game_type"`
	Status       GameStatus              `json:"status"`
	Deck         *Deck                   `json:"deck"`
	Players      []*Player               `json:"players"`
	Dealer       *Player                 `json:"dealer"`
	DiscardPiles map[string]*DiscardPile `json:"discard_piles"`
	MaxPlayers   int                     `json:"max_players"`
	CurrentPlayer int                    `json:"current_player"`
	CribbageState *CribbageState         `json:"cribbage_state,omitempty"`
	Created      time.Time               `json:"created"`
	LastUsed     time.Time               `json:"last_used"`
}

// NewGame creates a new blackjack game with the specified number of standard decks.
// This is a convenience function that defaults to standard 52-card decks.
func NewGame(numDecks int) *Game {
	return NewCustomGame(numDecks, Standard)
}

// NewCustomGame creates a new blackjack game with specified deck count and type.
// This allows using different deck types like Spanish21 while defaulting to blackjack.
func NewCustomGame(numDecks int, deckType DeckType) *Game {
	return NewGameWithType(numDecks, deckType, Blackjack, 6)
}

// NewGameWithType creates a fully customized game with all parameters specified.
// This is the most flexible constructor supporting different games, decks, and player limits.
func NewGameWithType(numDecks int, deckType DeckType, gameType GameType, maxPlayers int) *Game {
	game := &Game{
		ID:            uuid.New().String(),
		GameType:      gameType,
		Status:        GameWaiting,
		Deck:          NewCustomDeck(numDecks, deckType),
		Players:       []*Player{},
		Dealer:        &Player{ID: "dealer", Name: "Dealer", Hand: []*Card{}, Standing: false, Busted: false},
		DiscardPiles:  make(map[string]*DiscardPile),
		MaxPlayers:    maxPlayers,
		CurrentPlayer: 0,
		Created:       time.Now(),
		LastUsed:      time.Now(),
	}
	
	// Create a default discard pile
	game.DiscardPiles["main"] = &DiscardPile{
		ID:    "main",
		Name:  "Main Discard Pile",
		Cards: []*Card{},
	}
	
	return game
}

// UpdateLastUsed updates the game's last activity timestamp.
// This is used for game cleanup and determining inactive games for garbage collection.
func (g *Game) UpdateLastUsed() {
	g.LastUsed = time.Now()
}

// AddPlayer creates and adds a new player to the game.
// Returns nil if the game is at maximum capacity, otherwise returns the new player.
func (g *Game) AddPlayer(name string) *Player {
	if len(g.Players) >= g.MaxPlayers {
		return nil
	}
	
	player := &Player{
		ID:       uuid.New().String(),
		Name:     name,
		Hand:     []*Card{},
		Standing: false,
		Busted:   false,
	}
	g.Players = append(g.Players, player)
	return player
}

// GetPlayer retrieves a player by ID, including the special "dealer" player.
// Returns nil if the player is not found in the game.
func (g *Game) GetPlayer(playerID string) *Player {
	if playerID == "dealer" {
		return g.Dealer
	}
	
	for _, player := range g.Players {
		if player.ID == playerID {
			return player
		}
	}
	return nil
}

// RemovePlayer removes a player from the game by their ID.
// Returns true if the player was found and removed, false otherwise.
func (g *Game) RemovePlayer(playerID string) bool {
	for i, player := range g.Players {
		if player.ID == playerID {
			g.Players = append(g.Players[:i], g.Players[i+1:]...)
			return true
		}
	}
	return false
}

// DealToPlayer deals a single card from the deck to the specified player.
// The card's face-up status can be controlled, used for initial deals and hits.
func (g *Game) DealToPlayer(playerID string, faceUp bool) *Card {
	player := g.GetPlayer(playerID)
	if player == nil {
		return nil
	}
	
	card := g.Deck.Deal()
	if card == nil {
		return nil
	}
	
	card.FaceUp = faceUp
	player.AddCard(card)
	return card
}

// AddDiscardPile creates a new named discard pile for the game.
// Returns nil if a pile with the same ID already exists.
func (g *Game) AddDiscardPile(id, name string) *DiscardPile {
	if _, exists := g.DiscardPiles[id]; exists {
		return nil
	}
	
	pile := &DiscardPile{
		ID:    id,
		Name:  name,
		Cards: []*Card{},
	}
	g.DiscardPiles[id] = pile
	return pile
}

// GetDiscardPile retrieves a discard pile by its ID.
// Returns nil if the pile doesn't exist in the game.
func (g *Game) GetDiscardPile(id string) *DiscardPile {
	return g.DiscardPiles[id]
}