// pkg/goact/goact.go
package goact

import (
	"syscall/js"

	// 💡 แก้ไขการ Import ให้ถูกต้องตามโครงสร้างโปรเจกต์
	// แทนที่ "github.com/KPV2004/goact" ด้วย path จริงใน go.mod ของคุณ
	"github.com/KPV2004/goact/pkg/goact/hooks"
	"github.com/KPV2004/goact/pkg/goact/vdom"
)

// Internal state for managing the application lifecycle.
var (
	// rootDOM holds the reference to the actual HTML container element on the page.
	rootDOM js.Value
	// currentVNode keeps track of the currently rendered Virtual DOM tree for diffing.
	currentVNode *vdom.VNode
	// rootComponent stores the user's main application function.
	rootComponent func() *vdom.VNode
)

// internalRender executes the core rendering lifecycle.
//
// It resets the hook index, requests a new VDOM tree from the root component,
// and either performs an initial mount or an efficient update using the diff algorithm.
func internalRender() {
	// 💡 ปลดล็อก hooks.Reset()
	// CRITICAL: Reset the hook index pointer before calling the component.
	// This ensures that UseState returns the correct state based on call order.
	hooks.Reset()

	// Call the user's app function to get the intended VDOM structure.
	newVTree := rootComponent()

	if currentVNode == nil {
		// Initial Mount phase.
		// We clear the container's innerHTML just in case there's static loading text.
		rootDOM.Set("innerHTML", "")
		rootDOM.Call("appendChild", vdom.CreateDOM(newVTree))
	} else {
		// Update phase: Perform VDOM Diffing to update only changed parts.
		vdom.Diff(rootDOM, currentVNode, newVTree)
	}

	// Update the reference for the next render cycle.
	currentVNode = newVTree
}

// Mount attaches the Goact application to a specific HTML element on the page.
//
// elementID should match the ID of a container in your index.html (e.g., "root").
// app is the main functional component of your application.
func Mount(elementID string, app func() *vdom.VNode) {
	// Obtian document reference once.
	doc := js.Global().Get("document")
	
	// Get the container element.
	rootDOM = doc.Call("getElementById", elementID)
	
	// TODO: Add error checking here if rootDOM is null (element not found).

	// Store the application function.
	rootComponent = app

	// 💡 ปลดล็อก hooks.Init()
	// CRITICAL: Initialize the hooks package, linking it to Goact's internal
	// rendering function so that state updates trigger re-renders.
	hooks.Init(internalRender)

	// Trigger the initial render immediately.
	internalRender()

	// Block the main Go goroutine. WebAssembly apps must not exit, otherwise
	// they stop running. This allows JS callbacks to continue functioning.
	select {}
}
