package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
)

// DeckType represents the different types of card decks supported by the API.
// Each deck type has different card compositions and rules for gameplay.
type DeckType int

const (
	Standard DeckType = iota  // Standard 52-card deck with all ranks 1-13
	Spanish21                 // Spanish 21 deck with 48 cards (no 10s)
)

// GameType represents the different card games supported by the API.
// Each game type has specific rules and gameplay mechanics implemented.
type GameType int

const (
	Blackjack GameType = iota  // Blackjack with dealer and automatic gameplay
	Poker                      // Poker game support
	War                        // War card game
	GoFish                     // Go Fish game
	Cribbage                   // Cribbage with full scoring and pegging
)

// String returns the string representation of a DeckType for API responses.
// This provides human-readable deck type names in JSON responses.
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

// String returns the string representation of a GameType for API responses.
// This provides human-readable game type names in JSON responses.
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

// Description returns a detailed explanation of what each deck type contains.
// This helps users understand the differences between deck types in API documentation.
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

// CardsPerDeck returns the number of cards in each deck type.
// This is used for deck initialization and game logic calculations.
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

// GetAllDeckTypes returns all supported deck types for API enumeration.
// This enables the /deck-types endpoint to list all available options.
func GetAllDeckTypes() []DeckType {
	return []DeckType{Standard, Spanish21}
}

// safeAdjectives contains family-friendly words used for generating random deck names.
// These adjectives are combined with nouns to create unique, memorable deck identifiers.
var safeAdjectives = []string{
	"Amazing", "Bright", "Clever", "Daring", "Eager", "Friendly", "Gentle", "Happy", "Jolly", "Kind",
	"Lucky", "Magic", "Noble", "Peaceful", "Quick", "Royal", "Smart", "Trusty", "Unique", "Wise",
	"Brave", "Calm", "Cool", "Fair", "Fast", "Good", "Great", "Nice", "Pure", "Safe",
	"Smooth", "Strong", "Sweet", "Warm", "Wild", "Young", "Agile", "Bold", "Bright", "Clean",
	"Clear", "Creative", "Curious", "Dancing", "Dreamy", "Electric", "Fantastic", "Glowing", "Golden", "Graceful",
	"Heroic", "Inspiring", "Joyful", "Laughing", "Mighty", "Mystic", "Perfect", "Playful", "Powerful", "Radiant",
	"Shining", "Singing", "Sparkling", "Sunny", "Swift", "Talented", "Vibrant", "Winning", "Zealous", "Cheerful",
}

// safeNouns contains family-friendly words used for generating random deck names.
// These nouns are combined with adjectives to create unique, memorable deck identifiers.
var safeNouns = []string{
	"Dragon", "Phoenix", "Eagle", "Tiger", "Lion", "Wolf", "Bear", "Fox", "Hawk", "Owl",
	"Star", "Moon", "Sun", "Cloud", "River", "Ocean", "Mountain", "Forest", "Garden", "Castle",
	"Knight", "Wizard", "Hero", "Champion", "Explorer", "Adventurer", "Captain", "Guardian", "Warrior", "Scout",
	"Arrow", "Sword", "Shield", "Crown", "Gem", "Crystal", "Diamond", "Gold", "Silver", "Treasure",
	"Thunder", "Lightning", "Rainbow", "Storm", "Wind", "Fire", "Ice", "Earth", "Sky", "Dawn",
	"Mystery", "Quest", "Journey", "Dream", "Hope", "Joy", "Peace", "Glory", "Honor", "Victory",
	"Magic", "Wonder", "Spirit", "Power", "Force", "Energy", "Light", "Flame", "Spark", "Glow",
}

// Suit represents the four traditional playing card suits.
// Used for card identification and game logic in various card games.
type Suit int

const (
	Hearts Suit = iota  // ♥ Red suit
	Diamonds            // ♦ Red suit
	Clubs               // ♣ Black suit  
	Spades              // ♠ Black suit
)

// String returns the name of the suit for display purposes.
// This provides human-readable suit names for API responses and card images.
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

// Rank represents the numerical value of a playing card (1-13).
// Ace=1, numbered cards=face value, Jack=11, Queen=12, King=13.
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

// String returns the display name for card ranks.
// Face cards return names (Ace, Jack, Queen, King) while number cards return their numeric value.
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

// Card represents a single playing card with rank, suit, and visibility.
// FaceUp determines whether the card is visible to players or hidden (like dealer's hole card).
type Card struct {
	Rank   Rank `json:"rank"`
	Suit   Suit `json:"suit"`
	FaceUp bool `json:"face_up"`
}

// CardWithImages extends Card with URLs to generated card images in multiple sizes.
// This provides card visuals for web interfaces while maintaining all card data.
type CardWithImages struct {
	Rank   Rank              `json:"rank"`
	Suit   Suit              `json:"suit"`
	FaceUp bool              `json:"face_up"`
	Images map[string]string `json:"images,omitempty"`
}

// String returns a human-readable representation of the card.
// Used for logging and debugging to identify specific cards.
func (c Card) String() string {
	return fmt.Sprintf("%s of %s", c.Rank, c.Suit)
}

// Value returns the base numeric value of the card (same as rank).
// This is used for basic card comparisons and non-game-specific operations.
func (c Card) Value() int {
	return int(c.Rank)
}

// BlackjackValue returns the point value of the card in Blackjack.
// Face cards are worth 10, Aces are 11 (soft value), and others are face value.
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

// CribbageValue returns the point value of the card in Cribbage.
// Face cards are worth 10, Aces are 1, and others are face value for counting.
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

// ToCardWithImages converts a Card to CardWithImages with URLs for generated card images.
// It creates URLs for three image sizes (icon, small, large) or shows card back if face down.
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

// ToCardWithImagesPtr safely converts a Card pointer to CardWithImages with image URLs.
// It handles nil pointers gracefully by returning an empty CardWithImages struct.
func ToCardWithImagesPtr(c *Card, baseURL string) CardWithImages {
	if c == nil {
		return CardWithImages{}
	}
	return c.ToCardWithImages(baseURL)
}

// generateDeckName creates a random, family-friendly name for new decks.
// It combines a random adjective with a random noun to ensure memorable, unique names.
func generateDeckName() string {
	rand.Seed(time.Now().UnixNano())
	adjective := safeAdjectives[rand.Intn(len(safeAdjectives))]
	noun := safeNouns[rand.Intn(len(safeNouns))]
	return adjective + " " + noun
}

// Deck represents a collection of playing cards with metadata.
// It maintains the card order for dealing and tracks the deck type for game rules.
type Deck struct {
	Cards    []Card   `json:"cards"`
	Name     string   `json:"name"`
	DeckType DeckType `json:"deck_type"`
}

// NewDeck creates a single standard 52-card deck.
// This is a convenience function that calls NewMultiDeck with 1 deck.
func NewDeck() *Deck {
	return NewMultiDeck(1)
}

// NewMultiDeck creates a deck with multiple standard decks shuffled together.
// This is commonly used in Blackjack where multiple decks reduce card counting effectiveness.
func NewMultiDeck(numDecks int) *Deck {
	return NewCustomDeck(numDecks, Standard)
}

// NewCustomDeck creates a deck with specified count and type (Standard or Spanish21).
// It handles deck type differences like Spanish21 missing 10s and ensures proper card generation.
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

// Reset restores the deck to a full single deck of the current type.
// All cards are restored and the deck is shuffled, maintaining the current deck type.
func (d *Deck) Reset() {
	d.ResetWithDecks(1)
}

// ResetWithDecks restores the deck with a specified number of decks.
// It maintains the current deck type while changing the number of deck copies.
func (d *Deck) ResetWithDecks(numDecks int) {
	d.ResetWithDecksAndType(numDecks, d.DeckType)
}

// ResetWithDecksAndType completely reconfigures the deck with new count and type.
// This allows changing both the number of decks and the deck type (Standard/Spanish21).
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

// Shuffle randomizes the order of all cards in the deck using Fisher-Yates algorithm.
// This ensures fair card distribution and prevents predictable card sequences.
func (d *Deck) Shuffle() {
	rand.Seed(time.Now().UnixNano())
	for i := len(d.Cards) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	}
}

// Deal removes and returns the top card from the deck.
// Returns nil if the deck is empty, allowing callers to handle empty deck scenarios.
func (d *Deck) Deal() *Card {
	if len(d.Cards) == 0 {
		return nil
	}
	card := d.Cards[0]
	d.Cards = d.Cards[1:]
	return &card
}

// RemainingCards returns the number of cards left in the deck.
// This is used for game logic and API responses to show deck status.
func (d *Deck) RemainingCards() int {
	return len(d.Cards)
}

// IsEmpty checks if the deck has no cards remaining.
// This is used to prevent dealing from empty decks and trigger deck resets.
func (d *Deck) IsEmpty() bool {
	return len(d.Cards) == 0
}

// Player represents a game participant with a unique ID, name, and hand of cards.
// Players can be human users or the dealer, identified by UUID or "dealer" respectively.
type Player struct {
	ID   string  `json:"id"`
	Name string  `json:"name"`
	Hand []*Card `json:"hand"`
}

// AddCard adds a new card to the player's hand.
// This is used when dealing cards or when players hit in blackjack.
func (p *Player) AddCard(card *Card) {
	p.Hand = append(p.Hand, card)
}

// RemoveCard removes and returns a card at the specified index from the player's hand.
// Returns nil if the index is invalid, used for discarding cards to piles.
func (p *Player) RemoveCard(cardIndex int) *Card {
	if cardIndex < 0 || cardIndex >= len(p.Hand) {
		return nil
	}
	card := p.Hand[cardIndex]
	p.Hand = append(p.Hand[:cardIndex], p.Hand[cardIndex+1:]...)
	return card
}

// HandSize returns the number of cards currently in the player's hand.
// This is used for game logic and API responses to show hand status.
func (p *Player) HandSize() int {
	return len(p.Hand)
}

// ClearHand removes all cards from the player's hand and returns them.
// This is used at the end of games to reset players for the next round.
func (p *Player) ClearHand() []*Card {
	cards := p.Hand
	p.Hand = []*Card{}
	return cards
}

// BlackjackHandValue calculates the optimal point value and detects blackjack.
// Returns total points and whether the hand is a blackjack (21 with exactly 2 cards).
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

// IsBusted checks if the player's hand value exceeds 21 in blackjack.
// A busted player automatically loses regardless of the dealer's hand.
func (p *Player) IsBusted() bool {
	value, _ := p.BlackjackHandValue()
	return value > 21
}

// HasBlackjack checks if the player has a natural blackjack (21 with exactly 2 cards).
// Blackjack beats a regular 21 and typically pays 3:2 in casino rules.
func (p *Player) HasBlackjack() bool {
	_, blackjack := p.BlackjackHandValue()
	return blackjack
}

// ScoreCribbageHand calculates the cribbage score for the player's hand plus starter card.
// This implements standard cribbage scoring: fifteens, pairs, runs, flush, and nobs.
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

// CribbagePlayValue returns the card's value during cribbage play phase.
// This is used for pegging and reaching the count of 31 during play.
func (c Card) CribbagePlayValue() int {
	return c.CribbageValue()
}

// ScoreCribbageCards calculates the total cribbage score for a collection of cards.
// This implements all cribbage scoring rules: fifteens, pairs, runs, flush, and nobs.
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

// TakeTopCard removes and returns the most recently discarded card from the pile.
// Returns nil if the pile is empty, used for drawing from discard piles.
func (dp *DiscardPile) TakeTopCard() *Card {
	if len(dp.Cards) == 0 {
		return nil
	}
	card := dp.Cards[len(dp.Cards)-1]
	dp.Cards = dp.Cards[:len(dp.Cards)-1]
	return card
}

// Size returns the number of cards currently in the discard pile.
// This is used for game logic and API responses to show pile status.
func (dp *DiscardPile) Size() int {
	return len(dp.Cards)
}

// Clear removes all cards from the discard pile and returns them.
// This is used for reshuffling cards back into the deck or resetting games.
func (dp *DiscardPile) Clear() []*Card {
	cards := dp.Cards
	dp.Cards = []*Card{}
	return cards
}

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

// CribbagePhase represents the current phase of a cribbage game.
// Cribbage has distinct phases with different rules and valid actions.
type CribbagePhase int

const (
	CribbageDeal CribbagePhase = iota
	CribbageDiscard
	CribbagePlay
	CribbageShow
	CribbageFinished
)

// String returns the string representation of the cribbage phase.
// This is used in JSON responses to show the current phase of cribbage gameplay.
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

// CribbageState holds all game state specific to cribbage gameplay.
// This includes phase tracking, scoring, and the crib collection.
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
		ID:   uuid.New().String(),
		Name: name,
		Hand: []*Card{},
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

// StartBlackjackGame begins a blackjack game by dealing initial cards to all players.
// Each player and dealer receives 2 cards, with dealer's second card face down.
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

// PlayDealer executes the dealer's turn according to standard blackjack rules.
// Dealer hits on 16 and below, stands on 17 and above, then finishes the game.
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

// GetGameResult calculates the final outcome for each player in a finished blackjack game.
// Returns a map of player IDs to results: "blackjack", "win", "push", "bust", or "lose".
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

// StartCribbageGame initializes a new cribbage game with 2 players.
// Deals 6 cards to each player and sets up the cribbage game state.
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

// CribbageDiscard handles players discarding 2 cards to the crib during the discard phase.
// Once both players discard, it cuts the starter card and moves to the play phase.
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

// GameManager provides thread-safe management of multiple concurrent card games.
// It uses read-write mutexes to allow concurrent read access while ensuring write safety.
type GameManager struct {
	games map[string]*Game
	mutex sync.RWMutex
}

// NewGameManager creates a new game manager with an empty game collection.
// This is used as a singleton to manage all active games in the application.
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

// CustomCard represents a user-defined card with flexible attributes and optional game compatibility.
// Cards can have string or numeric ranks, custom attributes, and tombstone deletion for data integrity.
type CustomCard struct {
	Index          int                    `json:"index"`
	Name           string                 `json:"name"`
	Rank           interface{}            `json:"rank,omitempty"`
	Suit           string                 `json:"suit,omitempty"`
	GameCompatible bool                   `json:"game_compatible"`
	Attributes     map[string]string      `json:"attributes"`
	Deleted        bool                   `json:"deleted"`
}

// UpdateGameCompatibility determines if the card can be used in traditional card games.
// Cards with numeric ranks and suits are compatible, while string ranks are not.
func (cc *CustomCard) UpdateGameCompatibility() {
	if cc.Rank == nil || cc.Suit == "" {
		cc.GameCompatible = false
		return
	}
	
	switch cc.Rank.(type) {
	case int, int32, int64, float32, float64:
		cc.GameCompatible = true
	case string:
		cc.GameCompatible = false
	default:
		cc.GameCompatible = false
	}
}

// GetNumericRank extracts the numeric value from the card's rank if possible.
// Returns the rank as an integer and whether the conversion was successful.
func (cc *CustomCard) GetNumericRank() (int, bool) {
	if cc.Rank == nil {
		return 0, false
	}
	
	switch v := cc.Rank.(type) {
	case int:
		return v, true
	case int32:
		return int(v), true
	case int64:
		return int(v), true
	case float32:
		return int(v), true
	case float64:
		return int(v), true
	default:
		return 0, false
	}
}

// CustomDeck represents a collection of user-defined custom cards with metadata.
// It tracks card indices for consistent referencing and usage timestamps for cleanup.
type CustomDeck struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Cards       []*CustomCard `json:"cards"`
	NextIndex   int           `json:"next_index"`
	Created     time.Time     `json:"created"`
	LastUsed    time.Time     `json:"last_used"`
}

// NewCustomDeckTemplate creates a new empty custom deck with the given name.
// Initializes all required fields including UUID and timestamps for proper tracking.
func NewCustomDeckTemplate(name string) *CustomDeck {
	return &CustomDeck{
		ID:        uuid.New().String(),
		Name:      name,
		Cards:     []*CustomCard{},
		NextIndex: 0,
		Created:   time.Now(),
		LastUsed:  time.Now(),
	}
}

func (cd *CustomDeck) UpdateLastUsed() {
	cd.LastUsed = time.Now()
}

func (cd *CustomDeck) AddCard(name string, rank interface{}, suit string, attributes map[string]string) *CustomCard {
	if attributes == nil {
		attributes = make(map[string]string)
	}
	
	card := &CustomCard{
		Index:      cd.NextIndex,
		Name:       name,
		Rank:       rank,
		Suit:       suit,
		Attributes: attributes,
		Deleted:    false,
	}
	
	card.UpdateGameCompatibility()
	cd.Cards = append(cd.Cards, card)
	cd.NextIndex++
	cd.UpdateLastUsed()
	
	return card
}

func (cd *CustomDeck) GetCard(index int) *CustomCard {
	for _, card := range cd.Cards {
		if card.Index == index {
			return card
		}
	}
	return nil
}

func (cd *CustomDeck) DeleteCard(index int) bool {
	card := cd.GetCard(index)
	if card == nil {
		return false
	}
	
	card.Deleted = true
	cd.UpdateLastUsed()
	return true
}

func (cd *CustomDeck) ListCards(includeDeleted bool) []*CustomCard {
	if includeDeleted {
		return cd.Cards
	}
	
	activeCards := []*CustomCard{}
	for _, card := range cd.Cards {
		if !card.Deleted {
			activeCards = append(activeCards, card)
		}
	}
	return activeCards
}

func (cd *CustomDeck) GetGameCompatibleCards() []*CustomCard {
	compatibleCards := []*CustomCard{}
	for _, card := range cd.Cards {
		if !card.Deleted && card.GameCompatible {
			compatibleCards = append(compatibleCards, card)
		}
	}
	return compatibleCards
}

func (cd *CustomDeck) CardCount() int {
	count := 0
	for _, card := range cd.Cards {
		if !card.Deleted {
			count++
		}
	}
	return count
}

type CustomDeckManager struct {
	decks map[string]*CustomDeck
	mutex sync.RWMutex
}

func NewCustomDeckManager() *CustomDeckManager {
	return &CustomDeckManager{
		decks: make(map[string]*CustomDeck),
	}
}

func (cdm *CustomDeckManager) CreateDeck(name string) *CustomDeck {
	cdm.mutex.Lock()
	defer cdm.mutex.Unlock()
	
	deck := NewCustomDeckTemplate(name)
	cdm.decks[deck.ID] = deck
	return deck
}

func (cdm *CustomDeckManager) GetDeck(deckID string) (*CustomDeck, bool) {
	cdm.mutex.RLock()
	defer cdm.mutex.RUnlock()
	
	deck, exists := cdm.decks[deckID]
	if exists {
		deck.UpdateLastUsed()
	}
	return deck, exists
}

func (cdm *CustomDeckManager) DeleteDeck(deckID string) bool {
	cdm.mutex.Lock()
	defer cdm.mutex.Unlock()
	
	_, exists := cdm.decks[deckID]
	if exists {
		delete(cdm.decks, deckID)
	}
	return exists
}

func (cdm *CustomDeckManager) ListDecks() []*CustomDeck {
	cdm.mutex.RLock()
	defer cdm.mutex.RUnlock()
	
	decks := make([]*CustomDeck, 0, len(cdm.decks))
	for _, deck := range cdm.decks {
		decks = append(decks, deck)
	}
	return decks
}