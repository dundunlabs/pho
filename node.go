package tra

import (
	"errors"
	"fmt"
	"strings"
)

type node struct {
	children map[any]*node
	routes   map[string]route
}

func (n *node) findNode(path string) *node {
	before, after, found := strings.Cut(path, "/")
	node := n.children[before]

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

	if n.children[before] == nil {
		n.children[before] = &node{}
	}

	child := n.children[before]

	if found {
		return child.addNode(after)
	} else {
		return child
	}
}

func (n *node) addRoute(r route) {
	path, method := r.path, r.method
	node := n.addNode(path)

	if node.routes == nil {
		node.routes = map[string]route{}
	}

	if _, ok := node.routes[method]; ok {
		panic(errors.New(fmt.Sprintf("cannot add %q route", path)))
	}

	node.routes[method] = r
}
