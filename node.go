package tra

import (
	"errors"
	"fmt"
	"strings"
)

type nodeKey string

const (
	nodeKeyDynamic nodeKey = "dynamic"
)

type node struct {
	children map[any]*node
	routes   map[string]route
}

func (n *node) findChild(path string) *node {
	for _, key := range []any{path, nodeKeyDynamic} {
		child := n.children[key]
		if child != nil {
			return child
		}
	}

	return nil
}

func (n *node) findNode(path string) *node {
	before, after, found := strings.Cut(path, "/")
	node := n.findChild(before)

	if node == nil {
		return nil
	}

	if found {
		return node.findNode(after)
	} else {
		return node
	}
}

func (n *node) addNode(path string) *node {
	before, after, found := strings.Cut(path, "/")

	if n.children == nil {
		n.children = map[any]*node{}
	}

	var key any
	if strings.HasPrefix(before, ":") {
		key = nodeKeyDynamic
	} else {
		key = before
	}

	if n.children[key] == nil {
		n.children[key] = &node{}
	}

	child := n.children[key]

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
