// pkg/goact/vdom/vnode.go
package vdom

import (
	"syscall/js"
)

// VNode represents a node in the Virtual DOM tree.
// It can be either an HTML element, a text node, or a fragment.
type VNode struct {
	Type     string
	Tag      string
	Props    map[string]string
	Events   map[string]func(js.Value, []js.Value) any
	JSEvents map[string]js.Func
	Children []*VNode
	Text     string
	DOM      js.Value

	// ===== hooks =====
	Component func() *VNode
	States []any
	hooksIndex int
	// =================
}

