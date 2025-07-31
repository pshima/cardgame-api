package models

// Player represents a game participant with a unique ID, name, and hand of cards.
// Players can be human users or the dealer, identified by UUID or "dealer" respectively.
type Player struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Hand     []*Card `json:"hand"`
	Standing bool    `json:"standing,omitempty"`
	Busted   bool    `json:"busted,omitempty"`
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