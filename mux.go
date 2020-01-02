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
type Args struct {
	data map[string]string
	md   Metadata
}

// NewArgs provies an empty Args
func NewArgs() Args {
	return Args{
		data: make(map[string]string),
		md:   make(Metadata),
	}
}

// Get returns the value of a named argument
func (a *Args) Get(key string) string {
	return a.data[key]
}

// Set stores the value of a named argument
func (a *Args) Set(key, value string) {
	a.data[key] = value
}

// Get returns the value of a metadata value
func (a *Args) GetMeta(key string) string {
	return a.md[key]
}

// Set stores the value of a metadata key
func (a *Args) SetMeta(key, value string) {
	a.md[key] = value
}

// SetMetadata assigns the provided Metadata to Args, overwriting any existing values
func (a *Args) SetMetadata(md Metadata) {
	a.md = md
}

// Metadata is a collection of arguments set by the caller of ExecuteWithMetadata
type Metadata map[string]string

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

// Execute is the same as ExecuteWithMetadata, except it only accepts an input
// and passes along an empty Metadata
func (m *Mux) Execute(input string) error {
	md := make(Metadata)
	return m.ExecuteWithMetadata(md, input)
}

// ExecuteWithMetadata accepts Metadata and an input command, searches the trie
// for a handler, runs it, or returns ErrNotFound if a handler was not located
func (m *Mux) ExecuteWithMetadata(md Metadata, input string) error {
	path := m.tokenize(input)
	h, a, err := m.trie.Get(path)
	if err != nil {
		return err
	}

	if h == nil {
		return ErrNotFound
	}

	a.SetMetadata(md)
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
