package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	gameManager = NewGameManager()
	r := gin.Default()

	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, World!",
		})
	})

	r.GET("/deck-types", listDeckTypes)
	r.GET("/game/new", createNewGame)
	r.GET("/game/new/:decks", createNewGameWithDecks)
	r.GET("/game/new/:decks/:type", createNewGameWithType)
	r.GET("/game/:gameId/shuffle", shuffleDeck)
	r.GET("/game/:gameId", getGameInfo)
	r.GET("/game/:gameId/deal", dealCard)
	r.GET("/game/:gameId/deal/:count", dealCards)
	r.GET("/game/:gameId/reset", resetDeck)
	r.GET("/game/:gameId/reset/:decks", resetDeckWithDecks)
	r.GET("/game/:gameId/reset/:decks/:type", resetDeckWithType)
	r.DELETE("/game/:gameId", deleteGame)
	r.GET("/games", listGames)

	return r
}

func TestHelloEndpoint(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/hello", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Hello, World!", response["message"])
}

func createTestGame(router *gin.Engine) string {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/game/new", nil)
	router.ServeHTTP(w, req)
	
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	return response["game_id"].(string)
}

func TestCreateNewDeck(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/game/new", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "New Standard game created with 1 deck", response["message"])
	assert.Equal(t, float64(52), response["remaining_cards"])
	assert.Contains(t, response, "game_id")
}

func TestShuffleDeck(t *testing.T) {
	router := setupRouter()
	gameID := createTestGame(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/game/"+gameID+"/shuffle", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Deck shuffled", response["message"])
	assert.Equal(t, float64(52), response["remaining_cards"])
	assert.Equal(t, gameID, response["game_id"])
}

func TestGetDeckInfo(t *testing.T) {
	router := setupRouter()
	gameID := createTestGame(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/game/"+gameID, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(52), response["remaining_cards"])
	assert.Equal(t, false, response["is_empty"])
	assert.Equal(t, gameID, response["game_id"])
	
	cards, ok := response["cards"].([]interface{})
	assert.True(t, ok)
	assert.Equal(t, 52, len(cards))
}

func TestDealCard(t *testing.T) {
	router := setupRouter()
	gameID := createTestGame(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/game/"+gameID+"/deal", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(51), response["remaining_cards"])
	assert.Equal(t, gameID, response["game_id"])
	
	card, ok := response["card"].(map[string]interface{})
	assert.True(t, ok)
	assert.NotNil(t, card["rank"])
	assert.NotNil(t, card["suit"])
}

func TestDealCardFromEmptyDeck(t *testing.T) {
	router := setupRouter()
	gameID := createTestGame(router)
	
	// Deal all 52 cards
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/game/"+gameID+"/deal/52", nil)
	router.ServeHTTP(w, req)

	// Now try to deal another card
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/game/"+gameID+"/deal", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "No cards remaining in deck", response["error"])
}

func TestDealMultipleCards(t *testing.T) {
	router := setupRouter()
	gameID := createTestGame(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/game/"+gameID+"/deal/5", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(47), response["remaining_cards"])
	assert.Equal(t, float64(5), response["cards_dealt"])
	assert.Equal(t, gameID, response["game_id"])
	
	cards, ok := response["cards"].([]interface{})
	assert.True(t, ok)
	assert.Equal(t, 5, len(cards))
}

func TestDealMultipleCardsInvalidCount(t *testing.T) {
	router := setupRouter()
	gameID := createTestGame(router)

	testCases := []string{"/game/"+gameID+"/deal/0", "/game/"+gameID+"/deal/-1", "/game/"+gameID+"/deal/abc"}
	
	for _, testCase := range testCases {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", testCase, nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
		
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid count parameter", response["error"])
	}
}

func TestDealMoreCardsThanAvailable(t *testing.T) {
	router := setupRouter()
	gameID := createTestGame(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/game/"+gameID+"/deal/100", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Not enough cards remaining in deck", response["error"])
}

func TestResetDeck(t *testing.T) {
	router := setupRouter()
	gameID := createTestGame(router)
	
	// Deal some cards first
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/game/"+gameID+"/deal/10", nil)
	router.ServeHTTP(w, req)

	// Reset the deck
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/game/"+gameID+"/reset", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Deck reset", response["message"])
	assert.Equal(t, float64(52), response["remaining_cards"])
	assert.Equal(t, gameID, response["game_id"])
}

func TestCompleteWorkflow(t *testing.T) {
	router := setupRouter()

	// Create new game
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/game/new", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	
	var gameResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &gameResponse)
	assert.NoError(t, err)
	gameID := gameResponse["game_id"].(string)

	// Shuffle deck
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/game/"+gameID+"/shuffle", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	// Deal 2 cards
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/game/"+gameID+"/deal/2", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	
	var dealResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &dealResponse)
	assert.NoError(t, err)
	assert.Equal(t, float64(50), dealResponse["remaining_cards"])
	assert.Equal(t, gameID, dealResponse["game_id"])

	// Get game info
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/game/"+gameID, nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	
	var deckResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &deckResponse)
	assert.NoError(t, err)
	assert.Equal(t, float64(50), deckResponse["remaining_cards"])
	assert.Equal(t, false, deckResponse["is_empty"])
	assert.Equal(t, gameID, deckResponse["game_id"])

	// Reset deck
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/game/"+gameID+"/reset", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	
	var resetResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &resetResponse)
	assert.NoError(t, err)
	assert.Equal(t, float64(52), resetResponse["remaining_cards"])
	assert.Equal(t, gameID, resetResponse["game_id"])
}