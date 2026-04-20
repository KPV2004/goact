// pkg/goact/vdom/h.go
package vdom

import (
	"strings"
	"syscall/js"
)

// H creates a Virtual DOM element.
//
// Props are passed as a map[string]any. Keys that start with "on" (e.g., "onClick")
// are automatically registered as event listeners. The value must be a
// func(js.Value, []js.Value) any. Other keys are treated as standard HTML attributes
// and their values are converted to strings.
func H(tag string, props map[string]any, children ...*VNode) *VNode {
	
	// We make separate maps to distinguish between standard HTML attributes
	// and interactive JavaScript events, optimizing the DOM creation phase.
	processedProps := make(map[string]string)
	processedEvents := make(map[string]func(js.Value, []js.Value) any)

	for key, value := range props {
		// Event detection: Goact convention is "on" prefix for events.
		if strings.HasPrefix(key, "on") && len(key) > 2 {
			// e.g., "onClick" -> "click"
			eventName := strings.ToLower(key[2:]) 
			if handler, ok := value.(func(js.Value, []js.Value) any); ok {
				processedEvents[eventName] = handler
			}
		} else {
			// Attribute handling: convert value to string.
			// TODO: Add support for boolean attributes (e.g., 'checked', 'disabled')
			// as well as map-based style definitions in the next iteration.
			if strVal, ok := value.(string); ok {
				processedProps[key] = strVal
			} else if intVal, ok := value.(int); ok {
				processedProps[key] = js.ValueOf(intVal).String()
			}
		}
	}

	return &VNode{
		Type:     "element",
		Tag:      tag,
		Props:    processedProps,
		Events:   processedEvents,
		Children: children,
	}
}


// ==========================================
// 2. Tag Helpers (สิ่งที่ทำให้เกิด Error undefined)
// ==========================================

// Fragment is a special component that acts as a ghost container, 
// allowing a component to return multiple root elements simultaneously.
func Fragment(children ...*VNode) *VNode {
	return &VNode{
		Type:     "fragment",
		Children: children,
	}
}

// Div is a helper function to create a <div> element.
func Div(props map[string]any, children ...*VNode) *VNode {
	return H("div", props, children...)
}

// H1 is a helper function to create an <h1> element.
func H1(props map[string]any, children ...*VNode) *VNode {
	return H("h1", props, children...)
}

// Button is a helper function to create a <button> element.
func Button(props map[string]any, children ...*VNode) *VNode {
	return H("button", props, children...)
}

// Input is a helper function to create an <input> element.
// Note: It is self-closing and does not accept children.
func Input(props map[string]any) *VNode {
	return H("input", props)
}

// Text is a helper function to create a text content VNode.
func Text(content string) *VNode {
	return &VNode{
		Type: "text",
		Text: content,
	}
}
