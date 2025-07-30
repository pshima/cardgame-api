package models

import (
	"fmt"
)

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
	
	return scoreCribbageCards(allCards)
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

// scoreCribbageCards calculates the total cribbage score for a collection of cards.
// This implements all cribbage scoring rules: fifteens, pairs, runs, flush, and nobs.
func scoreCribbageCards(cards []*Card) int {
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