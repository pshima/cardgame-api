package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
)

type DeckType int

const (
	Standard DeckType = iota
	Spanish21
)

type GameType int

const (
	Blackjack GameType = iota
	Poker
	War
	GoFish
)

func (dt DeckType) String() string {
	switch dt {
	case Standard:
		return "Standard"
	case Spanish21:
		return "Spanish21"
	default:
		return "Standard"
	}
}

func (gt GameType) String() string {
	switch gt {
	case Blackjack:
		return "Blackjack"
	case Poker:
		return "Poker"
	case War:
		return "War"
	case GoFish:
		return "GoFish"
	default:
		return "Blackjack"
	}
}

func (dt DeckType) Description() string {
	switch dt {
	case Standard:
		return "Traditional 52-card deck with all ranks from Ace to King in all four suits"
	case Spanish21:
		return "Spanish 21 deck with 48 cards - all 10s removed, perfect for Spanish Blackjack"
	default:
		return "Traditional 52-card deck with all ranks from Ace to King in all four suits"
	}
}

func (dt DeckType) CardsPerDeck() int {
	switch dt {
	case Standard:
		return 52
	case Spanish21:
		return 48
	default:
		return 52
	}
}

func GetAllDeckTypes() []DeckType {
	return []DeckType{Standard, Spanish21}
}

var safeAdjectives = []string{
	"Amazing", "Bright", "Clever", "Daring", "Eager", "Friendly", "Gentle", "Happy", "Jolly", "Kind",
	"Lucky", "Magic", "Noble", "Peaceful", "Quick", "Royal", "Smart", "Trusty", "Unique", "Wise",
	"Brave", "Calm", "Cool", "Fair", "Fast", "Good", "Great", "Nice", "Pure", "Safe",
	"Smooth", "Strong", "Sweet", "Warm", "Wild", "Young", "Agile", "Bold", "Bright", "Clean",
	"Clear", "Creative", "Curious", "Dancing", "Dreamy", "Electric", "Fantastic", "Glowing", "Golden", "Graceful",
	"Heroic", "Inspiring", "Joyful", "Laughing", "Mighty", "Mystic", "Perfect", "Playful", "Powerful", "Radiant",
	"Shining", "Singing", "Sparkling", "Sunny", "Swift", "Talented", "Vibrant", "Winning", "Zealous", "Cheerful",
}

var safeNouns = []string{
	"Dragon", "Phoenix", "Eagle", "Tiger", "Lion", "Wolf", "Bear", "Fox", "Hawk", "Owl",
	"Star", "Moon", "Sun", "Cloud", "River", "Ocean", "Mountain", "Forest", "Garden", "Castle",
	"Knight", "Wizard", "Hero", "Champion", "Explorer", "Adventurer", "Captain", "Guardian", "Warrior", "Scout",
	"Arrow", "Sword", "Shield", "Crown", "Gem", "Crystal", "Diamond", "Gold", "Silver", "Treasure",
	"Thunder", "Lightning", "Rainbow", "Storm", "Wind", "Fire", "Ice", "Earth", "Sky", "Dawn",
	"Mystery", "Quest", "Journey", "Dream", "Hope", "Joy", "Peace", "Glory", "Honor", "Victory",
	"Magic", "Wonder", "Spirit", "Power", "Force", "Energy", "Light", "Flame", "Spark", "Glow",
}

type Suit int

const (
	Hearts Suit = iota
	Diamonds
	Clubs
	Spades
)

func (s Suit) String() string {
	switch s {
	case Hearts:
		return "Hearts"
	case Diamonds:
		return "Diamonds"
	case Clubs:
		return "Clubs"
	case Spades:
		return "Spades"
	default:
		return "Unknown"
	}
}

type Rank int

const (
	Ace Rank = iota + 1
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

func (r Rank) String() string {
	switch r {
	case Ace:
		return "Ace"
	case Jack:
		return "Jack"
	case Queen:
		return "Queen"
	case King:
		return "King"
	default:
		return fmt.Sprintf("%d", int(r))
	}
}

type Card struct {
	Rank   Rank `json:"rank"`
	Suit   Suit `json:"suit"`
	FaceUp bool `json:"face_up"`
}

// CardWithImages extends Card with image URLs
type CardWithImages struct {
	Rank   Rank              `json:"rank"`
	Suit   Suit              `json:"suit"`
	FaceUp bool              `json:"face_up"`
	Images map[string]string `json:"images,omitempty"`
}

func (c Card) String() string {
	return fmt.Sprintf("%s of %s", c.Rank, c.Suit)
}

func (c Card) Value() int {
	return int(c.Rank)
}

func (c Card) BlackjackValue() int {
	switch c.Rank {
	case Jack, Queen, King:
		return 10
	case Ace:
		return 11 // Will be handled as 1 or 11 in hand calculation
	default:
		return int(c.Rank)
	}
}

// ToCardWithImages converts a Card to CardWithImages with image URLs
func (c Card) ToCardWithImages(baseURL string) CardWithImages {
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	
	var images map[string]string
	if c.FaceUp {
		// Generate filename: rank_suit (e.g., "1_0" for Ace of Hearts)
		filename := fmt.Sprintf("%d_%d", int(c.Rank), int(c.Suit))
		images = map[string]string{
			"icon":  fmt.Sprintf("%s/static/cards/icon/%s.png", baseURL, filename),
			"small": fmt.Sprintf("%s/static/cards/small/%s.png", baseURL, filename),
			"large": fmt.Sprintf("%s/static/cards/large/%s.png", baseURL, filename),
		}
	} else {
		// Card back images
		images = map[string]string{
			"icon":  fmt.Sprintf("%s/static/cards/icon/back.png", baseURL),
			"small": fmt.Sprintf("%s/static/cards/small/back.png", baseURL),
			"large": fmt.Sprintf("%s/static/cards/large/back.png", baseURL),
		}
	}
	
	return CardWithImages{
		Rank:   c.Rank,
		Suit:   c.Suit,
		FaceUp: c.FaceUp,
		Images: images,
	}
}

// ToCardWithImagesPtr converts a *Card to CardWithImages
func ToCardWithImagesPtr(c *Card, baseURL string) CardWithImages {
	if c == nil {
		return CardWithImages{}
	}
	return c.ToCardWithImages(baseURL)
}

func generateDeckName() string {
	rand.Seed(time.Now().UnixNano())
	adjective := safeAdjectives[rand.Intn(len(safeAdjectives))]
	noun := safeNouns[rand.Intn(len(safeNouns))]
	return adjective + " " + noun
}

type Deck struct {
	Cards    []Card   `json:"cards"`
	Name     string   `json:"name"`
	DeckType DeckType `json:"deck_type"`
}

func NewDeck() *Deck {
	return NewMultiDeck(1)
}

func NewMultiDeck(numDecks int) *Deck {
	return NewCustomDeck(numDecks, Standard)
}

func NewCustomDeck(numDecks int, deckType DeckType) *Deck {
	if numDecks <= 0 {
		numDecks = 1
	}
	
	cardsPerDeck := 52
	if deckType == Spanish21 {
		cardsPerDeck = 48 // No 10s (4 cards removed per deck)
	}
	
	deck := &Deck{
		Cards:    make([]Card, 0, cardsPerDeck*numDecks),
		Name:     generateDeckName(),
		DeckType: deckType,
	}
	deck.ResetWithDecks(numDecks)
	return deck
}

func (d *Deck) Reset() {
	d.ResetWithDecks(1)
}

func (d *Deck) ResetWithDecks(numDecks int) {
	d.ResetWithDecksAndType(numDecks, d.DeckType)
}

func (d *Deck) ResetWithDecksAndType(numDecks int, deckType DeckType) {
	if numDecks <= 0 {
		numDecks = 1
	}
	
	d.DeckType = deckType
	cardsPerDeck := 52
	if deckType == Spanish21 {
		cardsPerDeck = 48
	}
	
	d.Cards = make([]Card, 0, cardsPerDeck*numDecks)
	
	for i := 0; i < numDecks; i++ {
		for suit := Hearts; suit <= Spades; suit++ {
			for rank := Ace; rank <= King; rank++ {
				// Skip 10s for Spanish 21
				if deckType == Spanish21 && rank == Ten {
					continue
				}
				d.Cards = append(d.Cards, Card{Rank: rank, Suit: suit, FaceUp: false})
			}
		}
	}
}

func (d *Deck) Shuffle() {
	rand.Seed(time.Now().UnixNano())
	for i := len(d.Cards) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	}
}

func (d *Deck) Deal() *Card {
	if len(d.Cards) == 0 {
		return nil
	}
	card := d.Cards[0]
	d.Cards = d.Cards[1:]
	return &card
}

func (d *Deck) RemainingCards() int {
	return len(d.Cards)
}

func (d *Deck) IsEmpty() bool {
	return len(d.Cards) == 0
}

type Player struct {
	ID   string  `json:"id"`
	Name string  `json:"name"`
	Hand []*Card `json:"hand"`
}

func (p *Player) AddCard(card *Card) {
	p.Hand = append(p.Hand, card)
}

func (p *Player) RemoveCard(cardIndex int) *Card {
	if cardIndex < 0 || cardIndex >= len(p.Hand) {
		return nil
	}
	card := p.Hand[cardIndex]
	p.Hand = append(p.Hand[:cardIndex], p.Hand[cardIndex+1:]...)
	return card
}

func (p *Player) HandSize() int {
	return len(p.Hand)
}

func (p *Player) ClearHand() []*Card {
	cards := p.Hand
	p.Hand = []*Card{}
	return cards
}

func (p *Player) BlackjackHandValue() (int, bool) {
	total := 0
	aces := 0
	
	for _, card := range p.Hand {
		value := card.BlackjackValue()
		if card.Rank == Ace {
			aces++
		}
		total += value
	}
	
	// Adjust for aces (convert from 11 to 1 if needed)
	for aces > 0 && total > 21 {
		total -= 10
		aces--
	}
	
	// Check for blackjack (21 with exactly 2 cards)
	blackjack := total == 21 && len(p.Hand) == 2
	
	return total, blackjack
}

func (p *Player) IsBusted() bool {
	value, _ := p.BlackjackHandValue()
	return value > 21
}

func (p *Player) HasBlackjack() bool {
	_, blackjack := p.BlackjackHandValue()
	return blackjack
}

type DiscardPile struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Cards []*Card `json:"cards"`
}

func (dp *DiscardPile) AddCard(card *Card) {
	dp.Cards = append(dp.Cards, card)
}

func (dp *DiscardPile) AddCards(cards []*Card) {
	dp.Cards = append(dp.Cards, cards...)
}

func (dp *DiscardPile) TopCard() *Card {
	if len(dp.Cards) == 0 {
		return nil
	}
	return dp.Cards[len(dp.Cards)-1]
}

func (dp *DiscardPile) TakeTopCard() *Card {
	if len(dp.Cards) == 0 {
		return nil
	}
	card := dp.Cards[len(dp.Cards)-1]
	dp.Cards = dp.Cards[:len(dp.Cards)-1]
	return card
}

func (dp *DiscardPile) Size() int {
	return len(dp.Cards)
}

func (dp *DiscardPile) Clear() []*Card {
	cards := dp.Cards
	dp.Cards = []*Card{}
	return cards
}

type GameStatus int

const (
	GameWaiting GameStatus = iota
	GameInProgress
	GameFinished
)

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
	Created      time.Time               `json:"created"`
	LastUsed     time.Time               `json:"last_used"`
}

func NewGame(numDecks int) *Game {
	return NewCustomGame(numDecks, Standard)
}

func NewCustomGame(numDecks int, deckType DeckType) *Game {
	return NewGameWithType(numDecks, deckType, Blackjack, 6)
}

func NewGameWithType(numDecks int, deckType DeckType, gameType GameType, maxPlayers int) *Game {
	game := &Game{
		ID:            uuid.New().String(),
		GameType:      gameType,
		Status:        GameWaiting,
		Deck:          NewCustomDeck(numDecks, deckType),
		Players:       []*Player{},
		Dealer:        &Player{ID: "dealer", Name: "Dealer", Hand: []*Card{}},
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

func (g *Game) UpdateLastUsed() {
	g.LastUsed = time.Now()
}

func (g *Game) AddPlayer(name string) *Player {
	if len(g.Players) >= g.MaxPlayers {
		return nil
	}
	
	player := &Player{
		ID:   uuid.New().String(),
		Name: name,
		Hand: []*Card{},
	}
	g.Players = append(g.Players, player)
	return player
}

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

func (g *Game) RemovePlayer(playerID string) bool {
	for i, player := range g.Players {
		if player.ID == playerID {
			g.Players = append(g.Players[:i], g.Players[i+1:]...)
			return true
		}
	}
	return false
}

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

func (g *Game) GetDiscardPile(id string) *DiscardPile {
	return g.DiscardPiles[id]
}

func (g *Game) StartBlackjackGame() error {
	if len(g.Players) == 0 {
		return fmt.Errorf("no players in game")
	}
	
	g.Status = GameInProgress
	g.CurrentPlayer = 0
	
	// Deal initial two cards to each player and dealer
	for i := 0; i < 2; i++ {
		// Deal to players
		for _, player := range g.Players {
			card := g.DealToPlayer(player.ID, true) // Face up for players
			if card == nil {
				return fmt.Errorf("not enough cards in deck")
			}
		}
		
		// Deal to dealer (first card face down, second face up)
		faceUp := i == 1
		card := g.DealToPlayer("dealer", faceUp)
		if card == nil {
			return fmt.Errorf("not enough cards in deck")
		}
	}
	
	return nil
}

func (g *Game) PlayerHit(playerID string) error {
	player := g.GetPlayer(playerID)
	if player == nil {
		return fmt.Errorf("player not found")
	}
	
	if g.Status != GameInProgress {
		return fmt.Errorf("game is not in progress")
	}
	
	card := g.DealToPlayer(playerID, true)
	if card == nil {
		return fmt.Errorf("no cards remaining in deck")
	}
	
	return nil
}

func (g *Game) PlayerStand(playerID string) error {
	player := g.GetPlayer(playerID)
	if player == nil {
		return fmt.Errorf("player not found")
	}
	
	if g.Status != GameInProgress {
		return fmt.Errorf("game is not in progress")
	}
	
	// Move to next player or dealer
	g.CurrentPlayer++
	if g.CurrentPlayer >= len(g.Players) {
		// All players finished, play dealer
		return g.PlayDealer()
	}
	
	return nil
}

func (g *Game) PlayDealer() error {
	// Reveal dealer's hole card
	if len(g.Dealer.Hand) > 0 {
		g.Dealer.Hand[0].FaceUp = true
	}
	
	// Dealer hits on 16 or less, stands on 17 or more
	for {
		value, _ := g.Dealer.BlackjackHandValue()
		if value >= 17 {
			break
		}
		
		card := g.DealToPlayer("dealer", true)
		if card == nil {
			break
		}
	}
	
	g.Status = GameFinished
	return nil
}

func (g *Game) GetGameResult() map[string]string {
	if g.Status != GameFinished {
		return map[string]string{"status": "game not finished"}
	}
	
	results := make(map[string]string)
	dealerValue, dealerBlackjack := g.Dealer.BlackjackHandValue()
	dealerBusted := g.Dealer.IsBusted()
	
	for _, player := range g.Players {
		playerValue, playerBlackjack := player.BlackjackHandValue()
		playerBusted := player.IsBusted()
		
		if playerBusted {
			results[player.ID] = "bust"
		} else if playerBlackjack && !dealerBlackjack {
			results[player.ID] = "blackjack"
		} else if dealerBusted {
			results[player.ID] = "win"
		} else if playerValue > dealerValue {
			results[player.ID] = "win"
		} else if playerValue == dealerValue {
			results[player.ID] = "push"
		} else {
			results[player.ID] = "lose"
		}
	}
	
	return results
}

type GameManager struct {
	games map[string]*Game
	mutex sync.RWMutex
}

func NewGameManager() *GameManager {
	return &GameManager{
		games: make(map[string]*Game),
	}
}

func (gm *GameManager) CreateGame(numDecks int) *Game {
	return gm.CreateCustomGame(numDecks, Standard)
}

func (gm *GameManager) CreateCustomGame(numDecks int, deckType DeckType) *Game {
	return gm.CreateGameWithType(numDecks, deckType, Blackjack, 6)
}

func (gm *GameManager) CreateGameWithType(numDecks int, deckType DeckType, gameType GameType, maxPlayers int) *Game {
	gm.mutex.Lock()
	defer gm.mutex.Unlock()
	
	game := NewGameWithType(numDecks, deckType, gameType, maxPlayers)
	gm.games[game.ID] = game
	return game
}

func (gm *GameManager) GetGame(gameID string) (*Game, bool) {
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

func (gm *GameManager) ListGames() []string {
	gm.mutex.RLock()
	defer gm.mutex.RUnlock()
	
	gameIDs := make([]string, 0, len(gm.games))
	for id := range gm.games {
		gameIDs = append(gameIDs, id)
	}
	return gameIDs
}

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

func (gm *GameManager) GameCount() int {
	gm.mutex.RLock()
	defer gm.mutex.RUnlock()
	return len(gm.games)
}