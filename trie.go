package cmdmux

import (
	"fmt"
	"strings"
)

type nodeClass uint8

const (
	rootNodeClass nodeClass = iota
	argNodeClass
	staticNodeClass
)

type node struct {
	path     string
	children []*node
	handle   Handler
	class    nodeClass
}

func isArg(token string) bool {
	if strings.HasPrefix(token, ":") {
		return true
	}

	return false
}

func (n *node) getChild(token string) *node {
	for i, c := range n.children {
		if c.path == token {
			return n.children[i]
		}
	}

	for i, c := range n.children {
		if c.class == argNodeClass {
			return n.children[i]
		}
	}

	return nil
}

func (n *node) Get(path []string) (Handler, Args, error) {
	trie := n

	argList := make(Args)

	for _, token := range path {
		c := trie.getChild(token)
		if c == nil {
			return nil, argList, fmt.Errorf("not found")
		}

		if c.class == argNodeClass {
			argList[c.path] = token
		}

		trie = c
	}

	return trie.handle, argList, nil
}

func (n *node) InsertChild(path []string, handle Handler) error {
	if len(path) < 1 {
		return ErrNotFound
	}

	trie := n
	for _, token := range path {
		child := trie.getChild(token)
		if child == nil {
			child = &node{
				path:  token,
				class: staticNodeClass,
			}
			if isArg(token) {
				child.class = argNodeClass
			}

			trie.children = append(trie.children, child)
		}
		trie = child
	}

	trie.handle = handle

	return nil
}
