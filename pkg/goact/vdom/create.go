// pkg/goact/vdom/create.go
package vdom

import (
	"syscall/js"
)

// CreateDOM instantiates a real browser DOM tree from a Virtual DOM `VNode` tree.
//
// It recursively creates HTML elements, sets their attributes, generates text nodes,
// and attaches event listeners as defined in the provided `VNode`.
// It returns a standard JavaScript DOM object (as js.Value) that is ready to be
// mounted to the actual page.
//
// It also ensures that interactive events are attached and their callback
// references (`js.Func`) are remembered, enabling proper memory management
// when the DOM element is removed or updated later.
func CreateDOM(vnode *VNode) js.Value {
	// Obtian the global 'document' object once for performance.
	doc := js.Global().Get("document")

	
	// Text Nodes are the simplest form; we handle them first.
	// This also saves the 'js.Value' reference back into the VNode for diffing.
	if vnode.Type == "text" {
		vnode.DOM = doc.Call("createTextNode", vnode.Text)
		return vnode.DOM
	}

	// For standard elements, we create the HTMLElement and set its standard properties.
	el := doc.Call("createElement", vnode.Tag)

	// Set standard attributes (e.g., 'class', 'id', 'style' as a string).
	for key, val := range vnode.Props {
		el.Call("setAttribute", key, val)
	}

	
	// Create and attach Event Listeners.
	// We MUST initialize and store the 'js.Func' callback references in the VNode.
	// This is CRITICAL for preventing memory leaks in WebAssembly by allowing us
	// to call '.Release()' on these functions during full DOM tree removal.
	vnode.JSEvents = make(map[string]js.Func)
	for eventName, handler := range vnode.Events {
		cb := js.FuncOf(handler)
		vnode.JSEvents[eventName] = cb
		el.Call("addEventListener", eventName, cb)
	}

	// Recursively handle all children by calling CreateDOM on each.
	for _, child := range vnode.Children {
		el.Call("appendChild", CreateDOM(child))
	}

	// Cache the final DOM node reference back to the VNode.
	vnode.DOM = el
	return el
}
