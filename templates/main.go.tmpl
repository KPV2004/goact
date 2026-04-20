// test/main.go
package main

import (
	"fmt"
	"syscall/js"

	// 💡 แก้ไขการ Import ให้ถูกต้องตามชื่อ module ของคุณใน go.mod
	"github.com/KPV2004/goact/pkg/goact" 
	. "github.com/KPV2004/goact/pkg/goact/hooks"
	. "github.com/KPV2004/goact/pkg/goact/vdom"
)

// Todo is a functional component representing the entire Todo Application.
func Todo() *VNode {
	// 💡 ใช้ Tag Helpers และ Props Sugar ที่ได้รับการอัปเกรดแล้ว
	// ในไฟล์ h.go ของคุณต้องรองรับ 'Div', 'vdom.Button', etc.
	
	text, setText := UseState[string]("")
	todos, setTodos := UseState[[]string]([]string{})

	// handleInput captures text from the input field.
	handleInput := func(this js.Value, args []js.Value) any {
		// Event convention: "onInput" -> 'input' event in JS
		setText(args[0].Get("target").Get("value").String())
		return nil
	}

	// handleAdd adds the current text to the todo list.
	handleAdd := func(this js.Value, args []js.Value) any {
		// Event convention: "onClick" -> 'click' event in JS
		if text == "" {
			return nil
		}
		// setTodos causes a re-render.
		setTodos(append(todos, text))
		setText("") // Clear input field.
		return nil
	}

	// Build the todo list items as VNodes.
	todoItems := []*VNode{}
	for i, item := range todos {
		todoItems = append(todoItems, Div(
			map[string]any{"style": "padding: 10px; border-bottom: 1px solid #eee;"},
			Text(fmt.Sprintf("%d. %s", i+1, item)),
		))
	}

	// Return the final application structure using upgraded h.go helpers.
	// We wrap everything in a Fragment to return multiple root-level elements.
	return Div(map[string]any{
			"style": "padding: 50px; font-family: sans-serif; max-width: 400px; margin: 0 auto; text-align: center;"},
			H1(nil, Text("Goact Todo List")),

			// Input and Button controls.
			Div(map[string]any{"style": "margin-bottom: 20px;"},
				Input(map[string]any{
					"type":  "text",
					"value": text,
					"onInput": handleInput, // Automated event binding
				}),
				Button(map[string]any{
					"style": "padding: 10px; margin-left: 10px;",
					"onClick": handleAdd, // Automated event binding
				}, Text("Add Todo")),
			),

			// Render the list of todo items.
			Div(map[string]any{"style": "text-align: left;"}, todoItems...),
	)
}

func main() {
	// Mount the application to the 'root' element on the page.
	goact.Mount("root", Todo)
}
