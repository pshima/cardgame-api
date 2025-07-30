package api

// CreateCustomDeckRequest represents the request body for creating custom decks
type CreateCustomDeckRequest struct {
	Name string `json:"name" binding:"required"`
}

// AddCustomCardRequest represents the request body for adding cards to custom decks
type AddCustomCardRequest struct {
	Name       string            `json:"name" binding:"required"`
	Rank       interface{}       `json:"rank,omitempty"`
	Suit       string            `json:"suit,omitempty"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

// AddPlayerRequest represents the request body for adding players to games
type AddPlayerRequest struct {
	Name string `json:"name" binding:"required"`
}

// DiscardRequest represents the request body for discarding cards
type DiscardRequest struct {
	PlayerID  string `json:"player_id" binding:"required"`
	CardIndex int    `json:"card_index"`
}

// CribbageDiscardRequest represents the request body for cribbage discard phase
type CribbageDiscardRequest struct {
	CardIndices []int `json:"card_indices" binding:"required"`
}

// CribbagePlayRequest represents the request body for cribbage play phase
type CribbagePlayRequest struct {
	CardIndex int `json:"card_index" binding:"required"`
}