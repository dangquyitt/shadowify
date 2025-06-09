package ftsearch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildTsqueryExpression(t *testing.T) {
	tests := []struct {
		input    string
		options  []Options
		expected string
	}{
		{
			input:    "hello world",
			options:  nil,
			expected: "hello & world",
		},
		{
			input:    "hello        world",
			options:  []Options{},
			expected: "hello & world",
		},
		{
			input:    "hello world",
			options:  []Options{WithPrefixMatching()},
			expected: "hello:* & world:*",
		},
		{
			input:    "hello world",
			options:  []Options{WithOperator(OR)},
			expected: "hello | world",
		},
		{
			input:    "hello world",
			options:  []Options{WithPrefixMatching(), WithOperator(OR)},
			expected: "hello:* | world:*",
		},
	}

	for _, test := range tests {
		result := BuildTsqueryExpression(test.input, test.options...)
		assert.Equal(t, test.expected, result)
	}
}
