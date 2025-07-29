package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeckTypeString(t *testing.T) {
	tests := []struct {
		deckType DeckType
		expected string
	}{
		{Standard, "Standard"},
		{Spanish21, "Spanish21"},
		{DeckType(99), "Standard"}, // Default case
	}

	for _, test := range tests {
		result := test.deckType.String()
		assert.Equal(t, test.expected, result)
	}
}

func TestDeckTypeDescription(t *testing.T) {
	tests := []struct {
		deckType DeckType
		expected string
	}{
		{Standard, "Traditional 52-card deck with all ranks from Ace to King in all four suits"},
		{Spanish21, "Spanish 21 deck with 48 cards - all 10s removed, perfect for Spanish Blackjack"},
		{DeckType(99), "Traditional 52-card deck with all ranks from Ace to King in all four suits"}, // Default case
	}

	for _, test := range tests {
		result := test.deckType.Description()
		assert.Equal(t, test.expected, result)
	}
}

func TestDeckTypeCardsPerDeck(t *testing.T) {
	tests := []struct {
		deckType DeckType
		expected int
	}{
		{Standard, 52},
		{Spanish21, 48},
		{DeckType(99), 52}, // Default case
	}

	for _, test := range tests {
		result := test.deckType.CardsPerDeck()
		assert.Equal(t, test.expected, result)
	}
}

func TestGetAllDeckTypes(t *testing.T) {
	deckTypes := GetAllDeckTypes()
	assert.Equal(t, 2, len(deckTypes))
	assert.Contains(t, deckTypes, Standard)
	assert.Contains(t, deckTypes, Spanish21)
}

func TestGenerateDeckName(t *testing.T) {
	// Test that deck names are generated
	name1 := generateDeckName()
	name2 := generateDeckName()
	
	assert.NotEmpty(t, name1)
	assert.NotEmpty(t, name2)
	assert.Contains(t, name1, " ") // Should contain a space between adjective and noun
	assert.Contains(t, name2, " ")
	
	// Names should be different (very high probability)
	// This test might occasionally fail due to randomness, but it's very unlikely
	assert.NotEqual(t, name1, name2)
	
	// Test that names contain valid words
	parts1 := strings.Split(name1, " ")
	parts2 := strings.Split(name2, " ")
	assert.Equal(t, 2, len(parts1))
	assert.Equal(t, 2, len(parts2))
}

func TestDeckNameWordSafety(t *testing.T) {
	// Test multiple generated names to ensure they're appropriate
	for i := 0; i < 50; i++ {
		name := generateDeckName()
		parts := strings.Split(name, " ")
		assert.Equal(t, 2, len(parts), "Name should have exactly 2 parts: "+name)
		
		adjective := parts[0]
		noun := parts[1]
		
		// Check that the words are from our safe lists
		adjFound := false
		for _, safeAdj := range safeAdjectives {
			if adjective == safeAdj {
				adjFound = true
				break
			}
		}
		assert.True(t, adjFound, "Adjective '"+adjective+"' not found in safe list")
		
		nounFound := false
		for _, safeNoun := range safeNouns {
			if noun == safeNoun {
				nounFound = true
				break
			}
		}
		assert.True(t, nounFound, "Noun '"+noun+"' not found in safe list")
	}
}

func TestNewCustomDeck(t *testing.T) {
	tests := []struct {
		numDecks     int
		deckType     DeckType
		expectedCards int
	}{
		{1, Standard, 52},
		{2, Standard, 104},
		{1, Spanish21, 48}, // No 10s
		{2, Spanish21, 96}, // No 10s, 2 decks
		{6, Spanish21, 288}, // 6 decks for Spanish 21
		{0, Standard, 52}, // Should default to 1
		{-1, Spanish21, 48}, // Should default to 1
	}

	for _, test := range tests {
		deck := NewCustomDeck(test.numDecks, test.deckType)
		assert.Equal(t, test.expectedCards, len(deck.Cards))
		assert.Equal(t, test.expectedCards, deck.RemainingCards())
		assert.Equal(t, test.deckType, deck.DeckType)
		assert.NotEmpty(t, deck.Name)
		assert.Contains(t, deck.Name, " ")
	}
}

func TestSpanish21DeckComposition(t *testing.T) {
	deck := NewCustomDeck(1, Spanish21)
	
	// Should have 48 cards (52 - 4 tens)
	assert.Equal(t, 48, len(deck.Cards))
	
	// Count cards by rank
	rankCounts := make(map[Rank]int)
	for _, card := range deck.Cards {
		rankCounts[card.Rank]++
	}
	
	// Should have 4 of each rank except Ten
	for rank := Ace; rank <= King; rank++ {
		if rank == Ten {
			assert.Equal(t, 0, rankCounts[rank], "Spanish21 deck should have no 10s")
		} else {
			assert.Equal(t, 4, rankCounts[rank], "Should have 4 of each rank except 10")
		}
	}
	
	// Verify all suits are present for non-Ten cards
	suitCounts := make(map[Suit]int)
	for _, card := range deck.Cards {
		suitCounts[card.Suit]++
	}
	
	// Should have 12 cards of each suit (13 - 1 ten)
	for suit := Hearts; suit <= Spades; suit++ {
		assert.Equal(t, 12, suitCounts[suit])
	}
}

func TestMultiDeckSpanish21(t *testing.T) {
	deck := NewCustomDeck(2, Spanish21)
	
	// Should have 96 cards (48 * 2)
	assert.Equal(t, 96, len(deck.Cards))
	
	// Count cards by rank
	rankCounts := make(map[Rank]int)
	for _, card := range deck.Cards {
		rankCounts[card.Rank]++
	}
	
	// Should have 8 of each rank except Ten (4 per deck * 2 decks)
	for rank := Ace; rank <= King; rank++ {
		if rank == Ten {
			assert.Equal(t, 0, rankCounts[rank])
		} else {
			assert.Equal(t, 8, rankCounts[rank])
		}
	}
}

func TestResetWithDecksAndType(t *testing.T) {
	// Start with a standard deck
	deck := NewCustomDeck(1, Standard)
	assert.Equal(t, 52, len(deck.Cards))
	assert.Equal(t, Standard, deck.DeckType)
	
	// Deal some cards
	deck.Deal()
	deck.Deal()
	assert.Equal(t, 50, deck.RemainingCards())
	
	// Reset to Spanish21
	deck.ResetWithDecksAndType(2, Spanish21)
	assert.Equal(t, 96, len(deck.Cards))
	assert.Equal(t, Spanish21, deck.DeckType)
	
	// Verify no 10s in the reset deck
	for _, card := range deck.Cards {
		assert.NotEqual(t, Ten, card.Rank, "Spanish21 deck should not contain 10s")
	}
}

func TestParseDeckType(t *testing.T) {
	tests := []struct {
		input    string
		expected DeckType
	}{
		{"spanish21", Spanish21},
		{"Spanish21", Spanish21},
		{"SPANISH21", Spanish21},
		{"spanish_21", Spanish21},
		{"spanish-21", Spanish21},
		{"standard", Standard},
		{"Standard", Standard},
		{"normal", Standard},
		{"regular", Standard},
		{"invalid", Standard}, // Default
		{"", Standard},        // Default
	}

	for _, test := range tests {
		result := parseDeckType(test.input)
		assert.Equal(t, test.expected, result, "Input: "+test.input)
	}
}

func TestCreateNewGameWithTypeEndpoint(t *testing.T) {
	router := setupSessionRouter()

	tests := []struct {
		decks        string
		deckType     string
		expectedCards int
		expectedType  string
		status       int
	}{
		{"1", "standard", 52, "Standard", 200},
		{"1", "spanish21", 48, "Spanish21", 200},
		{"2", "spanish21", 96, "Spanish21", 200},
		{"6", "spanish21", 288, "Spanish21", 200},
		{"8", "standard", 416, "Standard", 200},
		{"0", "spanish21", 0, "", 400}, // Invalid decks
		{"-1", "standard", 0, "", 400}, // Invalid decks
		{"abc", "spanish21", 0, "", 400}, // Invalid decks
	}

	for _, test := range tests {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/game/new/"+test.decks+"/"+test.deckType, nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, test.status, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		
		if test.status == 200 {
			assert.Contains(t, response, "game_id")
			assert.Contains(t, response, "deck_name")
			assert.Equal(t, test.expectedType, response["deck_type"])
			assert.Equal(t, float64(test.expectedCards), response["remaining_cards"])
			assert.Contains(t, response["message"], test.expectedType)
		} else {
			assert.Equal(t, "Invalid decks parameter", response["error"])
		}
	}
}

func TestResetDeckWithTypeEndpoint(t *testing.T) {
	router := setupSessionRouter()

	// Create a standard game first
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/game/new/1/standard", nil)
	router.ServeHTTP(w, req)
	
	var gameResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &gameResponse)
	gameID := gameResponse["game_id"].(string)

	// Reset to Spanish21
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/game/"+gameID+"/reset/2/spanish21", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	
	var resetResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resetResponse)
	assert.NoError(t, err)
	assert.Equal(t, gameID, resetResponse["game_id"])
	assert.Equal(t, "Spanish21", resetResponse["deck_type"])
	assert.Equal(t, float64(96), resetResponse["remaining_cards"]) // 2 Spanish21 decks
	assert.Contains(t, resetResponse["message"], "Spanish21")
}

func TestSpanish21GameWorkflow(t *testing.T) {
	router := setupSessionRouter()

	// Create Spanish21 game with 6 decks (typical for Spanish21)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/game/new/6/spanish21", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	
	var gameResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &gameResponse)
	assert.NoError(t, err)
	gameID := gameResponse["game_id"].(string)
	assert.Equal(t, "Spanish21", gameResponse["deck_type"])
	assert.Equal(t, float64(288), gameResponse["remaining_cards"]) // 6 * 48

	// Shuffle the deck
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/game/"+gameID+"/shuffle", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	
	var shuffleResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &shuffleResponse)
	assert.Equal(t, "Spanish21", shuffleResponse["deck_type"])

	// Deal some cards and verify no 10s
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/game/"+gameID+"/deal/20", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	
	var dealResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &dealResponse)
	assert.Equal(t, float64(268), dealResponse["remaining_cards"]) // 288 - 20
	
	cards := dealResponse["cards"].([]interface{})
	assert.Equal(t, 20, len(cards))
	
	// Verify no 10s were dealt
	for _, cardInterface := range cards {
		card := cardInterface.(map[string]interface{})
		rank := card["rank"].(float64)
		assert.NotEqual(t, float64(Ten), rank, "Spanish21 should not deal 10s")
	}

	// Get game info
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/game/"+gameID, nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	
	var infoResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &infoResponse)
	assert.Equal(t, "Spanish21", infoResponse["deck_type"])
	assert.Contains(t, infoResponse, "deck_name")
}

func TestConcurrentDifferentDeckTypes(t *testing.T) {
	router := setupSessionRouter()

	// Create multiple games with different deck types
	var gameIDs []string
	var expectedCards []int
	var expectedTypes []string

	games := []struct {
		decks string
		deckType string
		expectedCards int
		expectedType string
	}{
		{"1", "standard", 52, "Standard"},
		{"1", "spanish21", 48, "Spanish21"},
		{"2", "standard", 104, "Standard"},
		{"2", "spanish21", 96, "Spanish21"},
		{"6", "spanish21", 288, "Spanish21"},
	}

	for _, game := range games {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/game/new/"+game.decks+"/"+game.deckType, nil)
		router.ServeHTTP(w, req)
		
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		gameIDs = append(gameIDs, response["game_id"].(string))
		expectedCards = append(expectedCards, game.expectedCards)
		expectedTypes = append(expectedTypes, game.expectedType)
	}

	// Verify all games maintain their independent state
	for i, gameID := range gameIDs {
		// Deal different numbers of cards from each game
		dealCount := i + 1
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/game/"+gameID+"/deal/"+string(rune(dealCount+'0')), nil)
		router.ServeHTTP(w, req)
		
		var dealResponse map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &dealResponse)
		expectedRemaining := expectedCards[i] - dealCount
		assert.Equal(t, float64(expectedRemaining), dealResponse["remaining_cards"])
	}

	// Verify games still have correct types and independent state
	for i, gameID := range gameIDs {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/game/"+gameID, nil)
		router.ServeHTTP(w, req)
		
		var infoResponse map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &infoResponse)
		assert.Equal(t, expectedTypes[i], infoResponse["deck_type"])
		
		expectedRemaining := expectedCards[i] - (i + 1)
		assert.Equal(t, float64(expectedRemaining), infoResponse["remaining_cards"])
	}
}

func TestListDeckTypesEndpoint(t *testing.T) {
	router := setupSessionRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/deck-types", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	// Check response structure
	assert.Contains(t, response, "deck_types")
	assert.Contains(t, response, "count")
	assert.Equal(t, float64(2), response["count"])
	
	// Check deck types array
	deckTypes := response["deck_types"].([]interface{})
	assert.Equal(t, 2, len(deckTypes))
	
	// Check Standard deck type
	standardFound := false
	spanish21Found := false
	
	for _, deckTypeInterface := range deckTypes {
		deckType := deckTypeInterface.(map[string]interface{})
		
		assert.Contains(t, deckType, "id")
		assert.Contains(t, deckType, "type")
		assert.Contains(t, deckType, "name")
		assert.Contains(t, deckType, "description")
		assert.Contains(t, deckType, "cards_per_deck")
		
		if deckType["type"] == "Standard" {
			standardFound = true
			assert.Equal(t, float64(0), deckType["id"])
			assert.Equal(t, "Standard", deckType["name"])
			assert.Equal(t, "Traditional 52-card deck with all ranks from Ace to King in all four suits", deckType["description"])
			assert.Equal(t, float64(52), deckType["cards_per_deck"])
		}
		
		if deckType["type"] == "Spanish21" {
			spanish21Found = true
			assert.Equal(t, float64(1), deckType["id"])
			assert.Equal(t, "Spanish21", deckType["name"])
			assert.Equal(t, "Spanish 21 deck with 48 cards - all 10s removed, perfect for Spanish Blackjack", deckType["description"])
			assert.Equal(t, float64(48), deckType["cards_per_deck"])
		}
	}
	
	assert.True(t, standardFound, "Standard deck type should be present")
	assert.True(t, spanish21Found, "Spanish21 deck type should be present")
}