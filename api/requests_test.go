package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddPlayerRequest(t *testing.T) {
	// Test valid request
	request := AddPlayerRequest{
		Name: "Alice",
	}
	
	assert.Equal(t, "Alice", request.Name)
	assert.NotEmpty(t, request.Name)
}

func TestAddPlayerRequestEmpty(t *testing.T) {
	// Test empty request
	request := AddPlayerRequest{}
	
	assert.Equal(t, "", request.Name)
	assert.Empty(t, request.Name)
}

func TestAddPlayerRequestSpecialCharacters(t *testing.T) {
	// Test request with special characters
	request := AddPlayerRequest{
		Name: "Player-123_!@#",
	}
	
	assert.Equal(t, "Player-123_!@#", request.Name)
	assert.NotEmpty(t, request.Name)
}

func TestAddPlayerRequestLongName(t *testing.T) {
	// Test request with long name
	longName := "ThisIsAVeryLongPlayerNameThatExceedsNormalLimits"
	request := AddPlayerRequest{
		Name: longName,
	}
	
	assert.Equal(t, longName, request.Name)
	assert.NotEmpty(t, request.Name)
	assert.True(t, len(request.Name) > 30)
}

func TestAddPlayerRequestUnicodeCharacters(t *testing.T) {
	// Test request with unicode characters
	request := AddPlayerRequest{
		Name: "玩家123",
	}
	
	assert.Equal(t, "玩家123", request.Name)
	assert.NotEmpty(t, request.Name)
}

func TestAddPlayerRequestWhitespace(t *testing.T) {
	// Test request with whitespace
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "leading spaces",
			input:    "   Alice",
			expected: "   Alice",
		},
		{
			name:     "trailing spaces",
			input:    "Alice   ",
			expected: "Alice   ",
		},
		{
			name:     "internal spaces",
			input:    "Alice Bob",
			expected: "Alice Bob",
		},
		{
			name:     "only spaces",
			input:    "   ",
			expected: "   ",
		},
		{
			name:     "tabs and newlines",
			input:    "Alice\t\n",
			expected: "Alice\t\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := AddPlayerRequest{
				Name: test.input,
			}
			assert.Equal(t, test.expected, request.Name)
		})
	}
}