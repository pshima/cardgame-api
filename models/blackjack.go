package models

import "fmt"

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