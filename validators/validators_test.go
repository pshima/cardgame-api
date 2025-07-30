package validators

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
		result := ValidateUUID(test.input)
		assert.Equal(t, test.expected, result, test.desc)
	}
}

func TestValidatePlayerID(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
		desc     string
	}{
		{"123e4567-e89b-12d3-a456-426614174000", true, "valid UUID"},
		{"dealer", true, "dealer keyword"},
		{"DEALER", false, "dealer must be lowercase"},
		{"Dealer", false, "dealer must be lowercase"},
		{"not-a-uuid", false, "invalid format"},
		{"", false, "empty string"},
		{"<script>alert('xss')</script>", false, "XSS attempt"},
		{"../../../etc/passwd", false, "path traversal attempt"},
	}

	for _, test := range tests {
		result := ValidatePlayerID(test.input)
		assert.Equal(t, test.expected, result, test.desc)
	}
}

func TestValidateNumber(t *testing.T) {
	tests := []struct {
		input        string
		expectedNum  int
		expectedValid bool
		desc         string
	}{
		{"1", 1, true, "valid number 1"},
		{"5", 5, true, "valid number 5"},
		{"100", 100, true, "valid number 100"},
		{"0", 0, true, "zero allowed"},
		{"-1", 0, false, "negative number"},
		{"101", 101, true, "number 101 allowed"},
		{"abc", 0, false, "non-numeric"},
		{"", 0, false, "empty string"},
		{"1.5", 0, false, "decimal number"},
		{"1e10", 0, false, "scientific notation"},
		{" 5 ", 0, false, "number with spaces"},
		{"5abc", 0, false, "number with letters"},
	}

	for _, test := range tests {
		num, valid := ValidateNumber(test.input)
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
		{"standard", true, "standard deck"},
		{"spanish21", true, "spanish21 deck"},
		{"Spanish21", true, "Spanish21 deck"},
		{"SPANISH21", true, "SPANISH21 deck"},
		{"spanish_21", true, "spanish_21 deck"},
		{"spanish-21", true, "spanish-21 deck"},
		{"normal", true, "normal deck"},
		{"regular", true, "regular deck"},
		{"invalid", true, "invalid deck type allowed by pattern"},
		{"", false, "empty string"},
		{"<script>", false, "XSS attempt"},
		{"../../../etc/passwd", false, "path traversal"},
	}

	for _, test := range tests {
		result := ValidateDeckType(test.input)
		assert.Equal(t, test.expected, result, test.desc)
	}
}

func TestValidatePileID(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
		desc     string
	}{
		{"main", true, "main pile"},
		{"custom", true, "custom pile"},
		{"pile123", true, "alphanumeric pile"},
		{"test-pile", true, "pile with dash"},
		{"test_pile", true, "pile with underscore"},
		{"", false, "empty string"},
		{"a", true, "single character allowed"},
		{"this-is-a-very-long-pile-name-that-exceeds-limits-123456789012345", false, "too long"},
		{"pile with spaces", false, "spaces not allowed"},
		{"pile<script>", false, "XSS attempt"},
		{"../pile", false, "path traversal"},
		{"pile\nname", false, "newline character"},
		{"pile\ttab", false, "tab character"},
	}

	for _, test := range tests {
		result := ValidatePileID(test.input)
		assert.Equal(t, test.expected, result, test.desc)
	}
}

func TestValidateBoolean(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
		desc     string
	}{
		{"true", true, "true value"},
		{"false", true, "false value"},
		{"TRUE", true, "uppercase true allowed"},
		{"FALSE", true, "uppercase false allowed"},
		{"True", true, "capitalized true allowed"},
		{"False", true, "capitalized false allowed"},
		{"1", true, "numeric 1 allowed"},
		{"0", true, "numeric 0 allowed"},
		{"yes", false, "yes value"},
		{"no", false, "no value"},
		{"", false, "empty string"},
		{"<script>", false, "XSS attempt"},
	}

	for _, test := range tests {
		result := ValidateBoolean(test.input)
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
		{"hello", 10, "hello", "normal string"},
		{"hello world", 5, "hello", "truncated string"},
		{"test\x00string", 20, "teststring", "null character removed"},
		{"test\nstring", 20, "teststring", "newline removed"},
		{"test\tstring", 20, "teststring", "tab removed"},
		{"test\rstring", 20, "teststring", "carriage return removed"},
		{"test\x01\x02\x03string", 20, "teststring", "control characters removed"},
		{"", 10, "", "empty string"},
		{"hello\x7fworld", 20, "helloworld", "DEL character removed"},
		{"normal text with spaces", 30, "normal text with spaces", "spaces preserved"},
		{"übung", 10, "übung", "unicode preserved"},
	}

	for _, test := range tests {
		result := SanitizeString(test.input, test.maxLength)
		assert.Equal(t, test.expected, result, test.desc)
	}
}