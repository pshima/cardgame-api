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
	Cribbage
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
	case Cribbage:
		return "Cribbage"
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

func (c Card) CribbageValue() int {
	switch c.Rank {
	case Jack, Queen, King:
		return 10
	case Ace:
		return 1
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

// Cribbage scoring functions
func (p *Player) ScoreCribbageHand(starter *Card) int {
	if len(p.Hand) == 0 {
		return 0
	}
	
	// Create combined hand with starter card
	allCards := make([]*Card, len(p.Hand))
	copy(allCards, p.Hand)
	if starter != nil {
		allCards = append(allCards, starter)
	}
	
	return ScoreCribbageCards(allCards)
}

// CribbagePlayValue returns the value during the play phase (same as cribbage value)
func (c Card) CribbagePlayValue() int {
	return c.CribbageValue()
}

// Global cribbage scoring functions
func ScoreCribbageCards(cards []*Card) int {
	if len(cards) == 0 {
		return 0
	}
	
	score := 0
	
	// Score fifteens (2 points each)
	score += scoreFifteens(cards)
	
	// Score pairs (2 points each)
	score += scorePairs(cards)
	
	// Score runs (1 point per card)
	score += scoreRuns(cards)
	
	// Score flush (1 point per card if all same suit)
	score += scoreFlush(cards)
	
	// Score nobs (1 point if jack matches starter suit)
	score += scoreNobs(cards)
	
	return score
}

func scoreFifteens(cards []*Card) int {
	count := 0
	n := len(cards)
	
	// Check all possible combinations of cards
	for i := 1; i < (1 << n); i++ {
		sum := 0
		for j := 0; j < n; j++ {
			if i&(1<<j) != 0 {
				sum += cards[j].CribbageValue()
			}
		}
		if sum == 15 {
			count++
		}
	}
	
	return count * 2
}

func scorePairs(cards []*Card) int {
	rankCounts := make(map[Rank]int)
	for _, card := range cards {
		rankCounts[card.Rank]++
	}
	
	score := 0
	for _, count := range rankCounts {
		if count >= 2 {
			// n cards of same rank = n(n-1)/2 pairs
			score += count * (count - 1) / 2 * 2
		}
	}
	
	return score
}

func scoreRuns(cards []*Card) int {
	if len(cards) < 3 {
		return 0
	}
	
	rankCounts := make(map[Rank]int)
	for _, card := range cards {
		rankCounts[card.Rank]++
	}
	
	// Find consecutive ranks
	ranks := make([]Rank, 0, len(rankCounts))
	for rank := range rankCounts {
		ranks = append(ranks, rank)
	}
	
	// Sort ranks
	for i := 0; i < len(ranks)-1; i++ {
		for j := i + 1; j < len(ranks); j++ {
			if ranks[i] > ranks[j] {
				ranks[i], ranks[j] = ranks[j], ranks[i]
			}
		}
	}
	
	// Find longest consecutive sequence
	longestRun := 0
	currentRun := 1
	multiplier := rankCounts[ranks[0]]
	
	for i := 1; i < len(ranks); i++ {
		if int(ranks[i]) == int(ranks[i-1])+1 {
			currentRun++
			multiplier *= rankCounts[ranks[i]]
		} else {
			if currentRun >= 3 {
				longestRun = currentRun
				break
			}
			currentRun = 1
			multiplier = rankCounts[ranks[i]]
		}
	}
	
	if currentRun >= 3 {
		longestRun = currentRun
	}
	
	if longestRun >= 3 {
		return longestRun * multiplier
	}
	
	return 0
}

func scoreFlush(cards []*Card) int {
	if len(cards) < 4 {
		return 0
	}
	
	// Check if hand cards (excluding starter) are all same suit
	handSuit := cards[0].Suit
	handFlush := true
	
	// Assuming last card is starter if 5 cards, otherwise all are hand cards
	handSize := len(cards)
	if handSize == 5 {
		handSize = 4 // First 4 are hand, last is starter
	}
	
	for i := 1; i < handSize; i++ {
		if cards[i].Suit != handSuit {
			handFlush = false
			break
		}
	}
	
	if !handFlush {
		return 0
	}
	
	// If 5 cards and all same suit, score 5
	if len(cards) == 5 && cards[4].Suit == handSuit {
		return 5
	}
	
	// Otherwise score 4 for hand flush
	return 4
}

func scoreNobs(cards []*Card) int {
	if len(cards) != 5 {
		return 0
	}
	
	starter := cards[4] // Last card is starter
	
	// Check if any jack in hand matches starter suit
	for i := 0; i < 4; i++ {
		if cards[i].Rank == Jack && cards[i].Suit == starter.Suit {
			return 1
		}
	}
	
	return 0
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

type CribbagePhase int

const (
	CribbageDeal CribbagePhase = iota
	CribbageDiscard
	CribbagePlay
	CribbageShow
	CribbageFinished
)

func (cp CribbagePhase) String() string {
	switch cp {
	case CribbageDeal:
		return "deal"
	case CribbageDiscard:
		return "discard"
	case CribbagePlay:
		return "play"
	case CribbageShow:
		return "show"
	case CribbageFinished:
		return "finished"
	default:
		return "deal"
	}
}

type CribbageState struct {
	Phase         CribbagePhase `json:"phase"`
	Dealer        int           `json:"dealer"`
	Crib          []*Card       `json:"crib"`
	Starter       *Card         `json:"starter"`
	PlayedCards   []*Card       `json:"played_cards"`
	PlayTotal     int           `json:"play_total"`
	PlayCount     int           `json:"play_count"`
	PlayerScores  []int         `json:"player_scores"`
	GameScore     int           `json:"game_score"` // Target score (usually 121)
	CurrentGo     bool          `json:"current_go"`
	LastToPlay    int           `json:"last_to_play"`
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
	CribbageState *CribbageState         `json:"cribbage_state,omitempty"`
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

// Cribbage game methods
func (g *Game) StartCribbageGame() error {
	if len(g.Players) != 2 {
		return fmt.Errorf("cribbage requires exactly 2 players")
	}
	
	g.GameType = Cribbage
	g.Status = GameInProgress
	g.CribbageState = &CribbageState{
		Phase:        CribbageDeal,
		Dealer:       0,
		Crib:         []*Card{},
		PlayedCards:  []*Card{},
		PlayTotal:    0,
		PlayCount:    0,
		PlayerScores: make([]int, len(g.Players)),
		GameScore:    121,
		CurrentGo:    false,
		LastToPlay:   -1,
	}
	
	// Deal 6 cards to each player
	for i := 0; i < 6; i++ {
		for _, player := range g.Players {
			card := g.DealToPlayer(player.ID, true)
			if card == nil {
				return fmt.Errorf("not enough cards in deck")
			}
		}
	}
	
	g.CribbageState.Phase = CribbageDiscard
	g.CurrentPlayer = (g.CribbageState.Dealer + 1) % len(g.Players) // Non-dealer goes first
	
	return nil
}

func (g *Game) CribbageDiscard(playerID string, cardIndices []int) error {
	if g.CribbageState == nil || g.CribbageState.Phase != CribbageDiscard {
		return fmt.Errorf("not in discard phase")
	}
	
	player := g.GetPlayer(playerID)
	if player == nil {
		return fmt.Errorf("player not found")
	}
	
	if len(cardIndices) != 2 {
		return fmt.Errorf("must discard exactly 2 cards")
	}
	
	// Validate indices and discard cards to crib
	if len(player.Hand) != 6 {
		return fmt.Errorf("player must have 6 cards to discard")
	}
	
	// Sort indices in descending order to avoid index shifting
	for i := 0; i < len(cardIndices)-1; i++ {
		for j := i + 1; j < len(cardIndices); j++ {
			if cardIndices[i] < cardIndices[j] {
				cardIndices[i], cardIndices[j] = cardIndices[j], cardIndices[i]
			}
		}
	}
	
	// Remove cards from player's hand and add to crib
	for _, index := range cardIndices {
		if index < 0 || index >= len(player.Hand) {
			return fmt.Errorf("invalid card index: %d", index)
		}
		card := player.RemoveCard(index)
		if card != nil {
			g.CribbageState.Crib = append(g.CribbageState.Crib, card)
		}
	}
	
	// Check if both players have discarded
	if len(g.CribbageState.Crib) == 4 {
		// Cut starter card
		starter := g.Deck.Deal()
		if starter == nil {
			return fmt.Errorf("no cards remaining for starter")
		}
		starter.FaceUp = true
		g.CribbageState.Starter = starter
		
		// Check for "his heels" (Jack of same suit as starter = 2 points for dealer)
		if starter.Rank == Jack {
			dealerIndex := g.CribbageState.Dealer
			g.CribbageState.PlayerScores[dealerIndex] += 2
		}
		
		// Move to play phase
		g.CribbageState.Phase = CribbagePlay
		g.CurrentPlayer = (g.CribbageState.Dealer + 1) % len(g.Players) // Non-dealer plays first
	}
	
	return nil
}

func (g *Game) CribbagePlay(playerID string, cardIndex int) error {
	if g.CribbageState == nil || g.CribbageState.Phase != CribbagePlay {
		return fmt.Errorf("not in play phase")
	}
	
	player := g.GetPlayer(playerID)
	if player == nil {
		return fmt.Errorf("player not found")
	}
	
	playerIndex := -1
	for i, p := range g.Players {
		if p.ID == playerID {
			playerIndex = i
			break
		}
	}
	
	if playerIndex != g.CurrentPlayer {
		return fmt.Errorf("not your turn")
	}
	
	if cardIndex < 0 || cardIndex >= len(player.Hand) {
		return fmt.Errorf("invalid card index")
	}
	
	card := player.Hand[cardIndex]
	newTotal := g.CribbageState.PlayTotal + card.CribbagePlayValue()
	
	if newTotal > 31 {
		return fmt.Errorf("card would exceed 31")
	}
	
	// Play the card
	playedCard := player.RemoveCard(cardIndex)
	g.CribbageState.PlayedCards = append(g.CribbageState.PlayedCards, playedCard)
	g.CribbageState.PlayTotal = newTotal
	g.CribbageState.PlayCount++
	g.CribbageState.LastToPlay = playerIndex
	
	// Score pegging points
	points := g.scorePegging()
	g.CribbageState.PlayerScores[playerIndex] += points
	
	// Check for end of play round or game
	if newTotal == 31 || g.allHandsEmpty() {
		g.resetPlayRound()
	} else {
		g.CurrentPlayer = (g.CurrentPlayer + 1) % len(g.Players)
	}
	
	// Check if play phase is complete
	if g.allHandsEmpty() && g.CribbageState.PlayTotal == 0 {
		g.CribbageState.Phase = CribbageShow
		g.CurrentPlayer = (g.CribbageState.Dealer + 1) % len(g.Players) // Non-dealer shows first
	}
	
	return nil
}

func (g *Game) CribbageGo(playerID string) error {
	if g.CribbageState == nil || g.CribbageState.Phase != CribbagePlay {
		return fmt.Errorf("not in play phase")
	}
	
	playerIndex := -1
	for i, p := range g.Players {
		if p.ID == playerID {
			playerIndex = i
			break
		}
	}
	
	if playerIndex != g.CurrentPlayer {
		return fmt.Errorf("not your turn")
	}
	
	// Check if player can actually play (must say go if can't play)
	player := g.Players[playerIndex]
	canPlay := false
	for _, card := range player.Hand {
		if g.CribbageState.PlayTotal+card.CribbagePlayValue() <= 31 {
			canPlay = true
			break
		}
	}
	
	if canPlay {
		return fmt.Errorf("you must play a card if possible")
	}
	
	// Move to next player
	g.CurrentPlayer = (g.CurrentPlayer + 1) % len(g.Players)
	
	// If opponent also can't play, current player gets 1 point for "go"
	opponent := g.Players[g.CurrentPlayer]
	opponentCanPlay := false
	for _, card := range opponent.Hand {
		if g.CribbageState.PlayTotal+card.CribbagePlayValue() <= 31 {
			opponentCanPlay = true
			break
		}
	}
	
	if !opponentCanPlay {
		// Last to play gets 1 point for go
		g.CribbageState.PlayerScores[g.CribbageState.LastToPlay] += 1
		g.resetPlayRound()
	}
	
	return nil
}

func (g *Game) CribbageShow() map[string]interface{} {
	if g.CribbageState == nil || g.CribbageState.Phase != CribbageShow {
		return nil
	}
	
	scores := make(map[string]interface{})
	
	// Score non-dealer's hand first
	nonDealer := (g.CribbageState.Dealer + 1) % len(g.Players)
	nonDealerScore := g.Players[nonDealer].ScoreCribbageHand(g.CribbageState.Starter)
	g.CribbageState.PlayerScores[nonDealer] += nonDealerScore
	scores[g.Players[nonDealer].ID] = nonDealerScore
	
	// Score dealer's hand
	dealerScore := g.Players[g.CribbageState.Dealer].ScoreCribbageHand(g.CribbageState.Starter)
	g.CribbageState.PlayerScores[g.CribbageState.Dealer] += dealerScore
	scores[g.Players[g.CribbageState.Dealer].ID] = dealerScore
	
	// Score crib (dealer gets these points)
	cribWithStarter := make([]*Card, len(g.CribbageState.Crib))
	copy(cribWithStarter, g.CribbageState.Crib)
	cribWithStarter = append(cribWithStarter, g.CribbageState.Starter)
	cribScore := ScoreCribbageCards(cribWithStarter)
	g.CribbageState.PlayerScores[g.CribbageState.Dealer] += cribScore
	scores["crib"] = cribScore
	
	// Check for game winner
	for i, score := range g.CribbageState.PlayerScores {
		if score >= g.CribbageState.GameScore {
			g.Status = GameFinished
			g.CribbageState.Phase = CribbageFinished
			scores["winner"] = i
			return scores
		}
	}
	
	// Move to next hand
	g.CribbageState.Dealer = (g.CribbageState.Dealer + 1) % len(g.Players)
	g.CribbageState.Phase = CribbageDeal
	
	// Clear hands and prepare for next deal
	for _, player := range g.Players {
		player.ClearHand()
	}
	g.CribbageState.Crib = []*Card{}
	g.CribbageState.Starter = nil
	g.CribbageState.PlayedCards = []*Card{}
	g.CribbageState.PlayTotal = 0
	g.CribbageState.PlayCount = 0
	g.CribbageState.LastToPlay = -1
	
	return scores
}

func (g *Game) scorePegging() int {
	if len(g.CribbageState.PlayedCards) == 0 {
		return 0
	}
	
	points := 0
	playedCards := g.CribbageState.PlayedCards
	
	// Fifteen (2 points)
	if g.CribbageState.PlayTotal == 15 {
		points += 2
	}
	
	// Thirty-one (2 points)
	if g.CribbageState.PlayTotal == 31 {
		points += 2
	}
	
	// Pairs (2 points each)
	if len(playedCards) >= 2 {
		lastCard := playedCards[len(playedCards)-1]
		pairCount := 1
		
		for i := len(playedCards) - 2; i >= 0; i-- {
			if playedCards[i].Rank == lastCard.Rank {
				pairCount++
			} else {
				break
			}
		}
		
		if pairCount >= 2 {
			// n of a kind = n * (n-1) points
			points += pairCount * (pairCount - 1)
		}
	}
	
	// Runs (1 point per card)
	if len(playedCards) >= 3 {
		points += g.scorePlayRun()
	}
	
	return points
}

func (g *Game) scorePlayRun() int {
	playedCards := g.CribbageState.PlayedCards
	if len(playedCards) < 3 {
		return 0
	}
	
	// Check for run at end of played cards
	for runLength := len(playedCards); runLength >= 3; runLength-- {
		startIndex := len(playedCards) - runLength
		ranks := make([]Rank, runLength)
		
		for i := 0; i < runLength; i++ {
			ranks[i] = playedCards[startIndex+i].Rank
		}
		
		// Sort ranks to check for consecutive sequence
		for i := 0; i < len(ranks)-1; i++ {
			for j := i + 1; j < len(ranks); j++ {
				if ranks[i] > ranks[j] {
					ranks[i], ranks[j] = ranks[j], ranks[i]
				}
			}
		}
		
		// Check if consecutive
		isRun := true
		for i := 1; i < len(ranks); i++ {
			if int(ranks[i]) != int(ranks[i-1])+1 {
				isRun = false
				break
			}
		}
		
		if isRun {
			return runLength
		}
	}
	
	return 0
}

func (g *Game) allHandsEmpty() bool {
	for _, player := range g.Players {
		if len(player.Hand) > 0 {
			return false
		}
	}
	return true
}

func (g *Game) resetPlayRound() {
	// Last to play gets 1 point for last card (if not 31)
	if g.CribbageState.PlayTotal != 31 && g.CribbageState.LastToPlay >= 0 {
		g.CribbageState.PlayerScores[g.CribbageState.LastToPlay] += 1
	}
	
	g.CribbageState.PlayTotal = 0
	g.CribbageState.PlayedCards = []*Card{}
	g.CribbageState.LastToPlay = -1
	
	// Find next player who can play
	for i := 0; i < len(g.Players); i++ {
		if len(g.Players[i].Hand) > 0 {
			g.CurrentPlayer = i
			break
		}
	}
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