package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/peteshima/cardgame-api/validators"
)

func TestValidateUUID(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
		desc     string
	}{
		{"123e4567-e89b-12d3-a456-426614174000", true, "valid UUID"},
		{"550e8400-e29b-41d4-a716-446655440000", true, "valid UUID variant"},
		{"not-a-uuid", false, "invalid format"},
		{"123e4567-e89b-12d3-a456", false, "incomplete UUID"},
		{"123e4567-e89b-12d3-a456-42661417400g", false, "invalid character"},
		{"", false, "empty string"},
		{"123e4567-e89b-12d3-a456-426614174000-extra", false, "extra characters"},
		{"<script>alert('xss')</script>", false, "XSS attempt"},
		{"../../../etc/passwd", false, "path traversal attempt"},
	}

	for _, test := range tests {
		result := validators.ValidateUUID(test.input)
		assert.Equal(t, test.expected, result, test.desc)
	}
}

func TestValidatePlayerID(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
		desc     string
	}{
		{"dealer", true, "special dealer ID"},
		{"123e4567-e89b-12d3-a456-426614174000", true, "valid UUID"},
		{"DEALER", false, "uppercase dealer (not allowed)"},
		{"dealer123", false, "dealer with extra chars"},
		{"player1", false, "invalid format"},
		{"", false, "empty string"},
	}

	for _, test := range tests {
		result := validators.ValidatePlayerID(test.input)
		assert.Equal(t, test.expected, result, test.desc)
	}
}

func TestValidateNumber(t *testing.T) {
	tests := []struct {
		input         string
		expectedNum   int
		expectedValid bool
		desc          string
	}{
		{"42", 42, true, "valid positive number"},
		{"0", 0, true, "zero"},
		{"100", 100, true, "valid larger number"},
		{"-5", 0, false, "negative number"},
		{"abc", 0, false, "non-numeric"},
		{"12.5", 0, false, "decimal"},
		{"1e10", 0, false, "scientific notation"},
		{"", 0, false, "empty string"},
		{"999999999999999999999", 0, false, "overflow"},
		{"42'; DROP TABLE games;--", 0, false, "SQL injection attempt"},
	}

	for _, test := range tests {
		num, valid := validators.ValidateNumber(test.input)
		assert.Equal(t, test.expectedValid, valid, test.desc)
		if valid {
			assert.Equal(t, test.expectedNum, num, test.desc)
		}
	}
}

func TestValidateDeckType(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
		desc     string
	}{
		{"standard", true, "valid deck type"},
		{"spanish21", true, "valid deck type"},
		{"Standard", true, "uppercase"},
		{"spanish-21", true, "with hyphen"},
		{"spanish_21", true, "with underscore"},
		{"custom123", true, "alphanumeric"},
		{"very-long-deck-type-name-that-exceeds-limit", false, "too long"},
		{"deck/type", false, "contains slash"},
		{"deck\\type", false, "contains backslash"},
		{"deck type", false, "contains space"},
		{"<script>", false, "XSS attempt"},
		{"", false, "empty string"},
	}

	for _, test := range tests {
		result := validators.ValidateDeckType(test.input)
		assert.Equal(t, test.expected, result, test.desc)
	}
}

func TestValidatePileID(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
		desc     string
	}{
		{"main", true, "simple alphanumeric"},
		{"discard_pile_1", true, "with underscore and number"},
		{"player-cards", true, "with hyphen"},
		{"ABC123", true, "uppercase and numbers"},
		{"x", true, "single character"},
		{"", false, "empty string"},
		{"pile/123", false, "contains slash"},
		{"pile<script>", false, "contains angle bracket"},
		{"this-is-a-very-long-pile-id-that-exceeds-the-fifty-character-limit", false, "too long"},
		{"pile id", false, "contains space"},
		{"pile@home", false, "contains @"},
	}

	for _, test := range tests {
		result := validators.ValidatePileID(test.input)
		assert.Equal(t, test.expected, result, test.desc)
	}
}

func TestValidateBoolean(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
		desc     string
	}{
		{"true", true, "lowercase true"},
		{"false", true, "lowercase false"},
		{"1", true, "numeric 1"},
		{"0", true, "numeric 0"},
		{"True", true, "uppercase True"},
		{"FALSE", true, "uppercase FALSE"},
		{"yes", false, "yes not allowed"},
		{"no", false, "no not allowed"},
		{"on", false, "on not allowed"},
		{"off", false, "off not allowed"},
		{"", false, "empty string"},
		{"2", false, "invalid number"},
		{"truee", false, "misspelled"},
	}

	for _, test := range tests {
		result := validators.ValidateBoolean(test.input)
		assert.Equal(t, test.expected, result, test.desc)
	}
}

func TestSanitizeString(t *testing.T) {
	tests := []struct {
		input     string
		maxLength int
		expected  string
		desc      string
	}{
		{"hello world", 20, "hello world", "normal string"},
		{"hello\x00world", 20, "helloworld", "null byte removed"},
		{"hello\nworld", 20, "helloworld", "newline removed"},
		{"hello\tworld", 20, "helloworld", "tab removed"},
		{"hello\x1bworld", 20, "helloworld", "escape character removed"},
		{"very long string that exceeds limit", 10, "very long ", "truncated to max length"},
		{"test\x00\x01\x02\x03", 20, "test", "multiple control chars removed"},
		{"", 10, "", "empty string"},
		{"normal", 0, "", "zero max length"},
	}

	for _, test := range tests {
		result := validators.SanitizeString(test.input, test.maxLength)
		assert.Equal(t, test.expected, result, test.desc)
	}
}