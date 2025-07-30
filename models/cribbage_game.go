package models

import (
	"fmt"
)

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
	cribScore := scoreCribbageCards(cribWithStarter)
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