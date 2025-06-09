package ftsearch

import (
	"fmt"
	"strings"
)

type Operator string
type Options func(*ftsearch)

const (
	AND Operator = "&"
	OR  Operator = "|"
	NOT Operator = "!"
)

func BuildTsqueryExpression(input string, options ...Options) string {
	f := New(options...)
	return f.buildTsqueryExpression(input)
}

func WithPrefixMatching() Options {
	return func(f *ftsearch) {
		f.isPrefixMatching = true
	}
}

func WithOperator(operator Operator) Options {
	return func(f *ftsearch) {
		f.operator = operator
	}
}

type ftsearch struct {
	// Default is false
	isPrefixMatching bool
	// Default operator is AND
	operator Operator
}

func New(options ...Options) *ftsearch {
	f := &ftsearch{
		isPrefixMatching: false,
		operator:         AND,
	}
	for _, option := range options {
		option(f)
	}
	return f
}

func (f *ftsearch) buildTsqueryExpression(input string) string {
	lexemes := strings.Fields(input)

	if f.isPrefixMatching {
		for i, lexeme := range lexemes {
			lexemes[i] = lexeme + ":*"
		}
	}

	joinOperator := fmt.Sprintf(" %s ", f.operator)
	return strings.Join(lexemes, joinOperator)
}
