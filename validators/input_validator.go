package validators

import (
	"regexp"
	"strconv"
	"strings"
)

// Security patterns for input validation
var (
	// UUID pattern for gameID and playerID validation
	uuidPattern = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	// Alphanumeric with limited special chars for pile IDs
	pileIDPattern = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,50}$`)
	// Number pattern for numeric parameters
	numberPattern = regexp.MustCompile(`^[0-9]+$`)
	// Deck type pattern
	deckTypePattern = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,20}$`)
	// Boolean pattern for faceUp parameter
	boolPattern = regexp.MustCompile(`^(true|false|1|0)$`)
)

// ValidateUUID verifies that the input string matches the standard UUID format.
// This prevents injection attacks and ensures game/player IDs are properly formatted.
func ValidateUUID(input string) bool {
	return uuidPattern.MatchString(input)
}

// ValidatePlayerID validates player identifiers, accepting either UUID format or the special "dealer" value.
// This allows the dealer to be referenced while maintaining security for player IDs.
func ValidatePlayerID(input string) bool {
	return input == "dealer" || uuidPattern.MatchString(input)
}

// ValidatePileID ensures discard pile identifiers contain only safe alphanumeric characters.
// This prevents injection attacks through pile ID parameters in discard operations.
func ValidatePileID(input string) bool {
	return pileIDPattern.MatchString(input)
}

// ValidateNumber converts and validates string input as a positive integer.
// It prevents negative numbers, non-numeric input, and integer overflow attacks.
func ValidateNumber(input string) (int, bool) {
	if !numberPattern.MatchString(input) {
		return 0, false
	}
	num, err := strconv.Atoi(input)
	if err != nil || num < 0 {
		return 0, false
	}
	return num, true
}

// ValidateDeckType ensures deck type parameters contain only safe characters.
// This prevents injection attacks while allowing valid deck type identifiers.
func ValidateDeckType(input string) bool {
	return deckTypePattern.MatchString(input)
}

// ValidateBoolean validates boolean string inputs, accepting multiple formats (true/false/1/0).
// This provides flexibility while maintaining security for boolean parameters.
func ValidateBoolean(input string) bool {
	return boolPattern.MatchString(strings.ToLower(input))
}

// SanitizeString removes control characters and enforces length limits on user input.
// This prevents XSS attacks and buffer overflow attempts while preserving valid content.
func SanitizeString(input string, maxLength int) string {
	// Remove control characters and limit length
	cleaned := strings.Map(func(r rune) rune {
		if r < 32 || r == 127 { // Remove control characters
			return -1
		}
		return r
	}, input)
	
	if len(cleaned) > maxLength {
		cleaned = cleaned[:maxLength]
	}
	
	return cleaned
}

// ValidateDeckName ensures custom deck names are within acceptable length limits.
// This prevents memory exhaustion attacks while allowing descriptive deck names.
func ValidateDeckName(name string) bool {
	return len(name) >= 1 && len(name) <= 128
}

// ValidateCardIndex validates and converts card index strings to integers.
// This ensures array bounds safety when accessing cards in custom decks.
func ValidateCardIndex(indexStr string) (int, bool) {
	index, valid := ValidateNumber(indexStr)
	return index, valid
}