package models

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
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

// ParseDeckType converts string deck type parameters to the corresponding DeckType enum.
// It supports multiple string formats (spanish21, spanish_21, spanish-21) for flexibility.
func ParseDeckType(typeStr string) DeckType {
	typeStr = strings.ToLower(typeStr)
	switch typeStr {
	case "spanish21", "spanish_21", "spanish-21":
		return Spanish21
	case "standard", "normal", "regular":
		return Standard
	default:
		return Standard
	}
}

// SafeAdjectives contains family-friendly words used for generating random deck names.
// These adjectives are combined with nouns to create unique, memorable deck identifiers.
var SafeAdjectives = []string{
	"Amazing", "Bright", "Clever", "Daring", "Eager", "Friendly", "Gentle", "Happy", "Jolly", "Kind",
	"Lucky", "Magic", "Noble", "Peaceful", "Quick", "Royal", "Smart", "Trusty", "Unique", "Wise",
	"Brave", "Calm", "Cool", "Fair", "Fast", "Good", "Great", "Nice", "Pure", "Safe",
	"Smooth", "Strong", "Sweet", "Warm", "Wild", "Young", "Agile", "Bold", "Bright", "Clean",
	"Clear", "Creative", "Curious", "Dancing", "Dreamy", "Electric", "Fantastic", "Glowing", "Golden", "Graceful",
	"Heroic", "Inspiring", "Joyful", "Laughing", "Mighty", "Mystic", "Perfect", "Playful", "Powerful", "Radiant",
	"Shining", "Singing", "Sparkling", "Sunny", "Swift", "Talented", "Vibrant", "Winning", "Zealous", "Cheerful",
}

// SafeNouns contains family-friendly words used for generating random deck names.
// These nouns are combined with adjectives to create unique, memorable deck identifiers.
var SafeNouns = []string{
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

// CribbagePlayValue returns the card's value during cribbage play phase.
// This is used for pegging and reaching the count of 31 during play.
func (c Card) CribbagePlayValue() int {
	return c.CribbageValue()
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

// GenerateDeckName creates a random, family-friendly name for new decks.
// It combines a random adjective with a random noun to ensure memorable, unique names.
func GenerateDeckName() string {
	rand.Seed(time.Now().UnixNano())
	adjective := SafeAdjectives[rand.Intn(len(SafeAdjectives))]
	noun := SafeNouns[rand.Intn(len(SafeNouns))]
	return adjective + " " + noun
}