package tra

import (
	"errors"
	"fmt"
	"strings"
)

type nodeKey string

const (
	nodeKeyDynamic  nodeKey = "dynamic"
	nodeKeyWildcard nodeKey = "wildcard"
)

type node struct {
	children map[any]*node
	routes   map[string]route
	dynamic  bool
	wilcard  bool
}

func (n *node) findChild(key any) *node {
	child := n.children[key]
	if child != nil {
		return child
	}
	if n.wilcard {
		return n
	}
	return nil
}

func (n *node) findNode(path string) *node {
	before, after, found := strings.Cut(path, "/")
	var node *node

	for _, key := range []any{before, nodeKeyDynamic, nodeKeyWildcard} {
		child := n.findChild(key)

		if child == nil {
			continue
		}

		if found {
			node = child.findNode(after)
		} else {
			node = child
		}
		break
	}

	return node
}

func (n *node) addChild(path string) *node {
	dynamic := strings.HasPrefix(path, ":")
	wilcard := strings.HasPrefix(path, "*")

	if n.children == nil {
		n.children = map[any]*node{}
	}

	var key any
	if dynamic {
		key = nodeKeyDynamic
	} else if wilcard {
		key = nodeKeyWildcard
	} else {
		key = path
	}

	if n.children[key] == nil {
		n.children[key] = &node{
			dynamic: dynamic,
			wilcard: wilcard,
		}
	}

	return n.children[key]
}

func (n *node) addNode(path string) *node {
	before, after, found := strings.Cut(path, "/")
	child := n.addChild(before)

	if found {
		return child.addNode(after)
	} else {
		return child
	}
}

func (n *node) addRoute(r route) {
	path, method := r.path, r.method
	node := n.addNode(strings.TrimPrefix(path, "/"))

	if node.routes == nil {
		node.routes = map[string]route{}
	}

	if _, ok := node.routes[method]; ok {
		panic(errors.New(fmt.Sprintf("cannot add %q route", path)))
	}

	node.routes[method] = r
}
