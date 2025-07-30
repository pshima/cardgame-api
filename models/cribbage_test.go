package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCribbagePhaseString(t *testing.T) {
	tests := []struct {
		phase    CribbagePhase
		expected string
	}{
		{CribbageDeal, "deal"},
		{CribbageDiscard, "discard"},
		{CribbagePlay, "play"},
		{CribbageShow, "show"},
		{CribbageFinished, "finished"},
		{CribbagePhase(99), "deal"}, // Invalid phase defaults to deal
	}

	for _, test := range tests {
		result := test.phase.String()
		assert.Equal(t, test.expected, result)
	}
}

func TestStartCribbageGame(t *testing.T) {
	// Test successful cribbage game start
	game := NewGame(1)
	player1 := game.AddPlayer("Alice")
	player2 := game.AddPlayer("Bob")

	err := game.StartCribbageGame()
	assert.NoError(t, err)

	// Verify game state
	assert.Equal(t, Cribbage, game.GameType)
	assert.Equal(t, GameInProgress, game.Status)
	assert.NotNil(t, game.CribbageState)
	assert.Equal(t, CribbageDiscard, game.CribbageState.Phase)
	assert.Equal(t, 0, game.CribbageState.Dealer)
	assert.Equal(t, 121, game.CribbageState.GameScore)
	assert.Equal(t, 1, game.CurrentPlayer) // Non-dealer goes first

	// Verify players have 6 cards each
	assert.Equal(t, 6, len(player1.Hand))
	assert.Equal(t, 6, len(player2.Hand))

	// Verify crib is empty initially
	assert.Equal(t, 0, len(game.CribbageState.Crib))

	// Verify player scores initialized
	assert.Equal(t, 2, len(game.CribbageState.PlayerScores))
	assert.Equal(t, 0, game.CribbageState.PlayerScores[0])
	assert.Equal(t, 0, game.CribbageState.PlayerScores[1])
}

func TestStartCribbageGameWrongPlayerCount(t *testing.T) {
	// Test with 1 player
	game := NewGame(1)
	game.AddPlayer("Alice")

	err := game.StartCribbageGame()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cribbage requires exactly 2 players")

	// Test with 3 players
	game = NewGame(1)
	game.AddPlayer("Alice")
	game.AddPlayer("Bob")
	game.AddPlayer("Charlie")

	err = game.StartCribbageGame()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cribbage requires exactly 2 players")
}

func TestCribbageDiscard(t *testing.T) {
	// Set up game
	game := NewGame(1)
	player1 := game.AddPlayer("Alice")
	player2 := game.AddPlayer("Bob")
	game.StartCribbageGame()

	// Test first player discard
	err := game.CribbageDiscard(player1.ID, []int{0, 1})
	assert.NoError(t, err)
	assert.Equal(t, 4, len(player1.Hand)) // Should have 4 cards left
	assert.Equal(t, 2, len(game.CribbageState.Crib)) // Crib should have 2 cards

	// Test second player discard
	err = game.CribbageDiscard(player2.ID, []int{0, 1})
	assert.NoError(t, err)
	assert.Equal(t, 4, len(player2.Hand)) // Should have 4 cards left
	assert.Equal(t, 4, len(game.CribbageState.Crib)) // Crib should have 4 cards

	// Should move to play phase and cut starter
	assert.Equal(t, CribbagePlay, game.CribbageState.Phase)
	assert.NotNil(t, game.CribbageState.Starter)
	assert.True(t, game.CribbageState.Starter.FaceUp)
	assert.Equal(t, 1, game.CurrentPlayer) // Non-dealer plays first
}

func TestCribbageDiscardErrors(t *testing.T) {
	// Set up game
	game := NewGame(1)
	player1 := game.AddPlayer("Alice")
	_ = game.AddPlayer("Bob")
	game.StartCribbageGame()

	// Test wrong number of cards
	err := game.CribbageDiscard(player1.ID, []int{0})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must discard exactly 2 cards")

	err = game.CribbageDiscard(player1.ID, []int{0, 1, 2})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must discard exactly 2 cards")

	// Test invalid card indices - need fresh game since hand size checked first
	freshGame := NewGame(1)
	freshPlayer := freshGame.AddPlayer("TestPlayer")
	_ = freshGame.AddPlayer("Player2")
	freshGame.StartCribbageGame()
	
	err = freshGame.CribbageDiscard(freshPlayer.ID, []int{-1, 0})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid card index")

	// Test another fresh game for second invalid index test
	freshGame2 := NewGame(1)
	freshPlayer2 := freshGame2.AddPlayer("TestPlayer")
	_ = freshGame2.AddPlayer("Player2")
	freshGame2.StartCribbageGame()
	
	err = freshGame2.CribbageDiscard(freshPlayer2.ID, []int{0, 10})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid card index")

	// Test invalid player
	err = game.CribbageDiscard("invalid", []int{0, 1})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "player not found")

	// Test wrong phase
	game.CribbageState.Phase = CribbagePlay
	err = game.CribbageDiscard(player1.ID, []int{0, 1})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not in discard phase")
}

func TestCribbageDiscardHisHeels(t *testing.T) {
	// Set up game
	game := NewGame(1)
	player1 := game.AddPlayer("Alice")
	player2 := game.AddPlayer("Bob")
	game.StartCribbageGame()

	// Mock the deck to ensure Jack is cut as starter
	game.Deck.Cards = []Card{
		{Rank: Jack, Suit: Hearts, FaceUp: false},
	}

	// Both players discard
	game.CribbageDiscard(player1.ID, []int{0, 1})
	game.CribbageDiscard(player2.ID, []int{0, 1})

	// Dealer should get 2 points for "his heels"
	assert.Equal(t, 2, game.CribbageState.PlayerScores[0]) // Dealer is player 0
	assert.Equal(t, 0, game.CribbageState.PlayerScores[1]) // Non-dealer gets 0
}

func TestScoreCribbageHand(t *testing.T) {
	player := &Player{
		ID:   "test",
		Name: "Test Player",
		Hand: []*Card{},
	}

	// Test empty hand
	score := player.ScoreCribbageHand(nil)
	assert.Equal(t, 0, score)

	// Test hand with starter
	player.Hand = []*Card{
		{Rank: Five, Suit: Hearts},
		{Rank: Ten, Suit: Diamonds},
		{Rank: King, Suit: Clubs},
		{Rank: Ace, Suit: Spades},
	}
	starter := &Card{Rank: Four, Suit: Hearts}

	score = player.ScoreCribbageHand(starter)
	assert.Greater(t, score, 0) // Should score something for fifteens

	// Test hand without starter
	score = player.ScoreCribbageHand(nil)
	assert.GreaterOrEqual(t, score, 0)
}

func TestScoreFifteens(t *testing.T) {
	// Test simple fifteen
	cards := []*Card{
		{Rank: Five, Suit: Hearts},
		{Rank: Ten, Suit: Diamonds},
	}
	score := scoreFifteens(cards)
	assert.Equal(t, 2, score) // One fifteen = 2 points

	// Test multiple fifteens (3 fives = 15)
	cards = []*Card{
		{Rank: Five, Suit: Hearts},
		{Rank: Five, Suit: Diamonds},
		{Rank: Five, Suit: Clubs},
	}
	score = scoreFifteens(cards)
	assert.Equal(t, 2, score) // One fifteen (5+5+5) = 2 points

	// Test no fifteens
	cards = []*Card{
		{Rank: Two, Suit: Hearts},
		{Rank: Three, Suit: Diamonds},
	}
	score = scoreFifteens(cards)
	assert.Equal(t, 0, score)
}

func TestScorePairs(t *testing.T) {
	// Test no pairs
	cards := []*Card{
		{Rank: Two, Suit: Hearts},
		{Rank: Three, Suit: Diamonds},
	}
	score := scorePairs(cards)
	assert.Equal(t, 0, score)

	// Test one pair
	cards = []*Card{
		{Rank: Five, Suit: Hearts},
		{Rank: Five, Suit: Diamonds},
		{Rank: Three, Suit: Clubs},
	}
	score = scorePairs(cards)
	assert.Equal(t, 2, score) // One pair = 2 points

	// Test three of a kind
	cards = []*Card{
		{Rank: Five, Suit: Hearts},
		{Rank: Five, Suit: Diamonds},
		{Rank: Five, Suit: Clubs},
	}
	score = scorePairs(cards)
	assert.Equal(t, 6, score) // Three of a kind = 6 points (3 pairs)
}

func TestScoreRuns(t *testing.T) {
	// Test no run (less than 3 cards)
	cards := []*Card{
		{Rank: Five, Suit: Hearts},
		{Rank: Six, Suit: Diamonds},
	}
	score := scoreRuns(cards)
	assert.Equal(t, 0, score)

	// Test no run (not consecutive)
	cards = []*Card{
		{Rank: Five, Suit: Hearts},
		{Rank: Seven, Suit: Diamonds},
		{Rank: Nine, Suit: Clubs},
	}
	score = scoreRuns(cards)
	assert.Equal(t, 0, score)

	// Test run of 3
	cards = []*Card{
		{Rank: Five, Suit: Hearts},
		{Rank: Six, Suit: Diamonds},
		{Rank: Seven, Suit: Clubs},
	}
	score = scoreRuns(cards)
	assert.Equal(t, 3, score) // Run of 3 = 3 points

	// Test run of 4
	cards = []*Card{
		{Rank: Five, Suit: Hearts},
		{Rank: Six, Suit: Diamonds},
		{Rank: Seven, Suit: Clubs},
		{Rank: Eight, Suit: Spades},
	}
	score = scoreRuns(cards)
	assert.Equal(t, 4, score) // Run of 4 = 4 points
}

func TestScoreFlush(t *testing.T) {
	// Test no flush (less than 4 cards)
	cards := []*Card{
		{Rank: Five, Suit: Hearts},
		{Rank: Six, Suit: Hearts},
		{Rank: Seven, Suit: Hearts},
	}
	score := scoreFlush(cards)
	assert.Equal(t, 0, score)

	// Test hand flush (4 cards same suit)
	cards = []*Card{
		{Rank: Five, Suit: Hearts},
		{Rank: Six, Suit: Hearts},
		{Rank: Seven, Suit: Hearts},
		{Rank: Eight, Suit: Hearts},
	}
	score = scoreFlush(cards)
	assert.Equal(t, 4, score)

	// Test full flush (5 cards same suit)
	cards = []*Card{
		{Rank: Five, Suit: Hearts},
		{Rank: Six, Suit: Hearts},
		{Rank: Seven, Suit: Hearts},
		{Rank: Eight, Suit: Hearts},
		{Rank: Nine, Suit: Hearts}, // Starter
	}
	score = scoreFlush(cards)
	assert.Equal(t, 5, score)

	// Test no flush (mixed suits)
	cards = []*Card{
		{Rank: Five, Suit: Hearts},
		{Rank: Six, Suit: Diamonds},
		{Rank: Seven, Suit: Hearts},
		{Rank: Eight, Suit: Hearts},
	}
	score = scoreFlush(cards)
	assert.Equal(t, 0, score)
}

func TestScoreNobs(t *testing.T) {
	// Test no nobs (not 5 cards)
	cards := []*Card{
		{Rank: Jack, Suit: Hearts},
		{Rank: Six, Suit: Hearts},
		{Rank: Seven, Suit: Hearts},
		{Rank: Eight, Suit: Hearts},
	}
	score := scoreNobs(cards)
	assert.Equal(t, 0, score)

	// Test nobs (jack matches starter suit)
	cards = []*Card{
		{Rank: Jack, Suit: Hearts},
		{Rank: Six, Suit: Diamonds},
		{Rank: Seven, Suit: Clubs},
		{Rank: Eight, Suit: Spades},
		{Rank: Nine, Suit: Hearts}, // Starter - Hearts
	}
	score = scoreNobs(cards)
	assert.Equal(t, 1, score)

	// Test no nobs (jack doesn't match starter suit)
	cards = []*Card{
		{Rank: Jack, Suit: Hearts},
		{Rank: Six, Suit: Diamonds},
		{Rank: Seven, Suit: Clubs},
		{Rank: Eight, Suit: Spades},
		{Rank: Nine, Suit: Diamonds}, // Starter - Diamonds
	}
	score = scoreNobs(cards)
	assert.Equal(t, 0, score)

	// Test no nobs (no jack in hand)
	cards = []*Card{
		{Rank: Five, Suit: Hearts},
		{Rank: Six, Suit: Diamonds},
		{Rank: Seven, Suit: Clubs},
		{Rank: Eight, Suit: Spades},
		{Rank: Nine, Suit: Hearts}, // Starter
	}
	score = scoreNobs(cards)
	assert.Equal(t, 0, score)
}

func TestScoreCribbageCards(t *testing.T) {
	// Test empty cards
	score := scoreCribbageCards([]*Card{})
	assert.Equal(t, 0, score)

	// Test cards with multiple scoring combinations
	cards := []*Card{
		{Rank: Five, Suit: Hearts},
		{Rank: Five, Suit: Diamonds},
		{Rank: Ten, Suit: Clubs},
		{Rank: Jack, Suit: Spades},
	}
	score = scoreCribbageCards(cards)
	assert.Greater(t, score, 0) // Should score for pairs and fifteens
}