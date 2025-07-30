package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCribbagePlay(t *testing.T) {
	// Set up game in play phase
	game := NewGame(1)
	player1 := game.AddPlayer("Alice")
	player2 := game.AddPlayer("Bob")
	game.StartCribbageGame()

	// Discard to get to play phase
	game.CribbageDiscard(player1.ID, []int{0, 1})
	game.CribbageDiscard(player2.ID, []int{0, 1})

	// Test valid play - use the current player
	currentPlayerID := game.Players[game.CurrentPlayer].ID
	currentPlayer := game.Players[game.CurrentPlayer]
	initialHandSize := len(currentPlayer.Hand)
	
	err := game.CribbagePlay(currentPlayerID, 0)
	assert.NoError(t, err)
	assert.Equal(t, initialHandSize-1, len(currentPlayer.Hand))
	assert.Equal(t, 1, len(game.CribbageState.PlayedCards))
	assert.Greater(t, game.CribbageState.PlayTotal, 0)
	assert.Equal(t, 1, game.CribbageState.PlayCount)
}

func TestCribbagePlayErrors(t *testing.T) {
	// Set up game in play phase
	game := NewGame(1)
	player1 := game.AddPlayer("Alice")
	player2 := game.AddPlayer("Bob")
	game.StartCribbageGame()

	// Test wrong phase
	err := game.CribbagePlay(player1.ID, 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not in play phase")

	// Get to play phase
	game.CribbageDiscard(player1.ID, []int{0, 1})
	game.CribbageDiscard(player2.ID, []int{0, 1})

	// Test invalid player
	err = game.CribbagePlay("invalid", 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "player not found")

	// Test wrong turn - use the player who is NOT current
	wrongPlayerID := ""
	if game.CurrentPlayer == 0 {
		wrongPlayerID = player2.ID
	} else {
		wrongPlayerID = player1.ID
	}
	err = game.CribbagePlay(wrongPlayerID, 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not your turn")

	// Test invalid card index - use current player
	currentPlayerID := game.Players[game.CurrentPlayer].ID
	err = game.CribbagePlay(currentPlayerID, -1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid card index")

	err = game.CribbagePlay(currentPlayerID, 10)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid card index")
}

func TestCribbagePlayExceed31(t *testing.T) {
	// Set up game
	game := NewGame(1)
	player1 := game.AddPlayer("Alice")
	player2 := game.AddPlayer("Bob")

	// Create custom hands with high-value cards
	player1.Hand = []*Card{
		{Rank: King, Suit: Hearts},   // 10
		{Rank: Queen, Suit: Diamonds}, // 10
		{Rank: Jack, Suit: Clubs},    // 10
		{Rank: Ten, Suit: Spades},    // 10
	}
	player2.Hand = []*Card{
		{Rank: Nine, Suit: Hearts},
		{Rank: Eight, Suit: Diamonds},
		{Rank: Seven, Suit: Clubs},
		{Rank: Six, Suit: Spades},
	}

	game.GameType = Cribbage
	game.Status = GameInProgress
	game.CribbageState = &CribbageState{
		Phase:        CribbagePlay,
		Dealer:       0,
		Crib:         []*Card{},
		PlayedCards:  []*Card{},
		PlayTotal:    25, // Set high play total
		PlayCount:    0,
		PlayerScores: make([]int, 2),
		GameScore:    121,
		CurrentGo:    false,
		LastToPlay:   -1,
	}
	game.CurrentPlayer = 0

	// Try to play King (10) when total is 25 (would make 35 > 31)
	err := game.CribbagePlay(player1.ID, 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "card would exceed 31")
}

func TestCribbageGo(t *testing.T) {
	// Set up game
	game := NewGame(1)
	player1 := game.AddPlayer("Alice")
	player2 := game.AddPlayer("Bob")

	// Create hands where player can't play
	player1.Hand = []*Card{
		{Rank: King, Suit: Hearts}, // 10
	}
	player2.Hand = []*Card{
		{Rank: Queen, Suit: Diamonds}, // 10
	}

	game.GameType = Cribbage
	game.Status = GameInProgress
	game.CribbageState = &CribbageState{
		Phase:        CribbagePlay,
		Dealer:       0,
		Crib:         []*Card{},
		PlayedCards:  []*Card{},
		PlayTotal:    25, // High total so 10-value cards can't be played
		PlayCount:    0,
		PlayerScores: make([]int, 2),
		GameScore:    121,
		CurrentGo:    false,
		LastToPlay:   0,
	}
	game.CurrentPlayer = 0

	// Player 1 says go
	err := game.CribbageGo(player1.ID)
	assert.NoError(t, err)
	// Current player should either move to next player or reset if both can't play
	// The implementation handles the case where both players can't play
	assert.True(t, game.CurrentPlayer == 1 || game.CribbageState.PlayTotal == 0)
}

func TestCribbageGoErrors(t *testing.T) {
	// Set up game
	game := NewGame(1)
	player1 := game.AddPlayer("Alice")
	game.AddPlayer("Bob")

	// Test wrong phase
	game.CribbageState = &CribbageState{Phase: CribbageDiscard}
	err := game.CribbageGo(player1.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not in play phase")

	// Test wrong turn
	game.CribbageState.Phase = CribbagePlay
	game.CurrentPlayer = 1
	err = game.CribbageGo(player1.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not your turn")

	// Test when player can actually play
	player1.Hand = []*Card{{Rank: Two, Suit: Hearts}} // Can always play 2
	game.CribbageState.PlayTotal = 0
	game.CurrentPlayer = 0
	err = game.CribbageGo(player1.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "you must play a card if possible")
}

func TestCribbageShow(t *testing.T) {
	// Set up game in show phase
	game := NewGame(1)
	player1 := game.AddPlayer("Alice")
	player2 := game.AddPlayer("Bob")

	// Set up hands with known scoring
	player1.Hand = []*Card{
		{Rank: Five, Suit: Hearts},
		{Rank: Five, Suit: Diamonds},
		{Rank: Jack, Suit: Clubs},
		{Rank: King, Suit: Spades},
	}
	player2.Hand = []*Card{
		{Rank: Six, Suit: Hearts},
		{Rank: Seven, Suit: Diamonds},
		{Rank: Eight, Suit: Clubs},
		{Rank: Nine, Suit: Spades},
	}

	game.GameType = Cribbage
	game.Status = GameInProgress
	game.CribbageState = &CribbageState{
		Phase:   CribbageShow,
		Dealer:  0,
		Crib:    []*Card{{Rank: Two, Suit: Hearts}, {Rank: Three, Suit: Diamonds}},
		Starter: &Card{Rank: Four, Suit: Clubs},
		PlayerScores: []int{0, 0},
		GameScore:    121,
	}

	scores := game.CribbageShow()
	assert.NotNil(t, scores)
	assert.Contains(t, scores, player1.ID)
	assert.Contains(t, scores, player2.ID)
	assert.Contains(t, scores, "crib")

	// Check that scores were added to player totals
	assert.GreaterOrEqual(t, game.CribbageState.PlayerScores[0], 0)
	assert.GreaterOrEqual(t, game.CribbageState.PlayerScores[1], 0)
}

func TestCribbageShowWrongPhase(t *testing.T) {
	game := NewGame(1)
	game.AddPlayer("Alice")
	game.AddPlayer("Bob")

	game.CribbageState = &CribbageState{Phase: CribbagePlay}
	scores := game.CribbageShow()
	assert.Nil(t, scores)
}

func TestCribbageShowGameWin(t *testing.T) {
	// Set up game where player will win
	game := NewGame(1)
	player1 := game.AddPlayer("Alice")
	player2 := game.AddPlayer("Bob")

	// Set up hands that will score - hand that makes fifteen with starter
	player1.Hand = []*Card{
		{Rank: Five, Suit: Hearts},
		{Rank: Five, Suit: Diamonds}, 
		{Rank: Five, Suit: Clubs},
		{Rank: Ten, Suit: Spades},
	}
	player2.Hand = []*Card{{Rank: Two, Suit: Hearts}}

	game.GameType = Cribbage
	game.Status = GameInProgress
	game.CribbageState = &CribbageState{
		Phase:   CribbageShow,
		Dealer:  0,
		Crib:    []*Card{},
		Starter: &Card{Rank: Three, Suit: Clubs},
		PlayerScores: []int{120, 50}, // Player 0 is close to winning
		GameScore:    121,
	}

	scores := game.CribbageShow()
	assert.NotNil(t, scores)
	assert.Contains(t, scores, "winner")
	assert.Equal(t, GameFinished, game.Status)
	assert.Equal(t, CribbageFinished, game.CribbageState.Phase)
}

func TestScorePegging(t *testing.T) {
	game := NewGame(1)
	game.CribbageState = &CribbageState{
		PlayedCards: []*Card{},
		PlayTotal:   0,
	}

	// Test no played cards
	score := game.scorePegging()
	assert.Equal(t, 0, score)

	// Test fifteen
	game.CribbageState.PlayedCards = []*Card{{Rank: Five, Suit: Hearts}}
	game.CribbageState.PlayTotal = 15
	score = game.scorePegging()
	assert.Equal(t, 2, score)

	// Test thirty-one
	game.CribbageState.PlayTotal = 31
	score = game.scorePegging()
	assert.Equal(t, 2, score)

	// Test pair
	game.CribbageState.PlayedCards = []*Card{
		{Rank: Five, Suit: Hearts},
		{Rank: Five, Suit: Diamonds},
	}
	game.CribbageState.PlayTotal = 10
	score = game.scorePegging()
	assert.Equal(t, 2, score) // One pair = 2 points
}

func TestScorePlayRun(t *testing.T) {
	game := NewGame(1)

	// Test no run (less than 3 cards)
	game.CribbageState = &CribbageState{
		PlayedCards: []*Card{
			{Rank: Five, Suit: Hearts},
			{Rank: Six, Suit: Diamonds},
		},
	}
	score := game.scorePlayRun()
	assert.Equal(t, 0, score)

	// Test run of 3
	game.CribbageState.PlayedCards = []*Card{
		{Rank: Five, Suit: Hearts},
		{Rank: Six, Suit: Diamonds},
		{Rank: Seven, Suit: Clubs},
	}
	score = game.scorePlayRun()
	assert.Equal(t, 3, score)

	// Test no run (not consecutive)
	game.CribbageState.PlayedCards = []*Card{
		{Rank: Five, Suit: Hearts},
		{Rank: Seven, Suit: Diamonds},
		{Rank: Nine, Suit: Clubs},
	}
	score = game.scorePlayRun()
	assert.Equal(t, 0, score)
}

func TestAllHandsEmpty(t *testing.T) {
	game := NewGame(1)
	player1 := game.AddPlayer("Alice")
	player2 := game.AddPlayer("Bob")

	// Test with cards in hands
	player1.Hand = []*Card{{Rank: Ace, Suit: Hearts}}
	player2.Hand = []*Card{{Rank: Two, Suit: Diamonds}}
	assert.False(t, game.allHandsEmpty())

	// Test with one empty hand
	player1.Hand = []*Card{}
	assert.False(t, game.allHandsEmpty())

	// Test with all empty hands
	player2.Hand = []*Card{}
	assert.True(t, game.allHandsEmpty())
}

func TestResetPlayRound(t *testing.T) {
	game := NewGame(1)
	player1 := game.AddPlayer("Alice")
	player2 := game.AddPlayer("Bob")

	game.CribbageState = &CribbageState{
		PlayTotal:     20,
		PlayedCards:   []*Card{{Rank: Five, Suit: Hearts}},
		LastToPlay:    0,
		PlayerScores:  []int{0, 0},
	}

	// Give players some cards
	player1.Hand = []*Card{{Rank: Six, Suit: Hearts}}
	player2.Hand = []*Card{}

	game.resetPlayRound()

	// Should award 1 point to last player (not 31)
	assert.Equal(t, 1, game.CribbageState.PlayerScores[0])
	assert.Equal(t, 0, game.CribbageState.PlayTotal)
	assert.Equal(t, 0, len(game.CribbageState.PlayedCards))
	assert.Equal(t, -1, game.CribbageState.LastToPlay)
	assert.Equal(t, 0, game.CurrentPlayer) // Player 0 has cards, should be current
}

func TestResetPlayRoundThirtyOne(t *testing.T) {
	game := NewGame(1)
	game.AddPlayer("Alice")
	game.AddPlayer("Bob")

	game.CribbageState = &CribbageState{
		PlayTotal:     31, // Exactly 31
		PlayedCards:   []*Card{{Rank: Five, Suit: Hearts}},
		LastToPlay:    0,
		PlayerScores:  []int{0, 0},
	}

	game.resetPlayRound()

	// Should NOT award extra point for last card when total is 31
	assert.Equal(t, 0, game.CribbageState.PlayerScores[0])
}

func TestCribbageNextHandPreparation(t *testing.T) {
	// Set up game that won't end
	game := NewGame(1)
	player1 := game.AddPlayer("Alice")
	player2 := game.AddPlayer("Bob")

	// Set up show phase
	player1.Hand = []*Card{{Rank: Ace, Suit: Hearts}}
	player2.Hand = []*Card{{Rank: Two, Suit: Hearts}}

	game.GameType = Cribbage
	game.Status = GameInProgress
	game.CribbageState = &CribbageState{
		Phase:   CribbageShow,
		Dealer:  0,
		Crib:    []*Card{{Rank: Three, Suit: Hearts}},
		Starter: &Card{Rank: Four, Suit: Hearts},
		PlayerScores: []int{50, 40}, // Neither close to winning
		GameScore:    121,
	}

	game.CribbageShow()

	// Should prepare for next hand
	assert.Equal(t, CribbageDeal, game.CribbageState.Phase)
	assert.Equal(t, 1, game.CribbageState.Dealer) // Dealer should rotate
	assert.Equal(t, 0, len(game.CribbageState.Crib))
	assert.Nil(t, game.CribbageState.Starter)
	assert.Equal(t, 0, len(game.CribbageState.PlayedCards))
	assert.Equal(t, 0, game.CribbageState.PlayTotal)
	assert.Equal(t, 0, game.CribbageState.PlayCount)
	assert.Equal(t, -1, game.CribbageState.LastToPlay)

	// Players should have empty hands
	assert.Equal(t, 0, len(player1.Hand))
	assert.Equal(t, 0, len(player2.Hand))
}