package cmdmux

import (
	"errors"
	"strings"
	"text/scanner"
	"unicode"
)

var (
	// ErrNotFound is returned by the Mux when a Handler could not be found
	// along the search path
	ErrNotFound = errors.New("cmdmux: handler not found")
)

// Mux is a text based command multiplexer. It matches an input command
// against a list of registered command templates using a rudimentary
// prefix trie
type Mux struct {
	trie *node
	// Tokenizer specifies a custom function to call to tokenize input commands.
	// If Tokenizer is nil, DefaultTokenizer will be used
	Tokenizer TokenizerFunc
}

// Args is a collection of arguments extracted from an input command
type Args map[string]string

// The Handler type is used by functions registered to a command path
// in the prefix trie
type Handler func(Args) error

// TokenizerFunc functions accept a string, returning a string slice containing
// its individual ident tokens
type TokenizerFunc func(string) []string

// NewMux returns a fresh Mux in its default state, ready to be used
func NewMux() *Mux {
	return &Mux{
		trie: &node{
			class: rootNodeClass,
		},
	}
}

func (m *Mux) tokenize(input string) []string {
	if m.Tokenizer != nil {
		return m.Tokenizer(input)
	}

	return DefaultTokenizer(input)
}

// Execute accepts an input command, searches the trie for a handler, runs
// it, or returns ErrNotFound if a handler was not located
func (m *Mux) Execute(input string) error {
	path := m.tokenize(input)
	h, a, err := m.trie.Get(path)
	if err != nil {
		return err
	}

	if h == nil {
		return ErrNotFound
	}

	return h(a)
}

// Handle registers a Handler function for the provided command template
func (m *Mux) Handle(template string, h Handler) {
	path := m.tokenize(template)
	m.trie.InsertChild(path, h)
}

// DefaultTokenizer is the default strategy for tokenizing command templates and inputs
func DefaultTokenizer(input string) []string {
	var s scanner.Scanner
	var tokens []string

	s.Init(strings.NewReader(input))

	s.IsIdentRune = func(ch rune, i int) bool {
		return ch == ':' && i == 0 || unicode.IsLetter(ch) || unicode.IsDigit(ch) && i > 0 || ch == '-' && i > 0
	}

	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		tokens = append(tokens, s.TokenText())
	}

	return tokens
}
