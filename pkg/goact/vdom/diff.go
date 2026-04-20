// pkg/goact/vdom/diff.go
package vdom

import (
	"syscall/js"
)
func Diff(parentDOM js.Value, oldNode, newNode *VNode) {
	if oldNode == nil {
		parentDOM.Call("appendChild", CreateDOM(newNode))
		return
	}
	if newNode == nil {
		parentDOM.Call("removeChild", oldNode.DOM)
		return
	}
	if oldNode.Type != newNode.Type || oldNode.Tag != newNode.Tag {
		newDOM := CreateDOM(newNode)
		parentDOM.Call("replaceChild", newDOM, oldNode.DOM)
		return
	}

	newNode.DOM = oldNode.DOM

	if newNode.Type == "text" {
		if oldNode.Text != newNode.Text {
			newNode.DOM.Set("nodeValue", newNode.Text)
		}
		return
	}

	// ==========================================
	// พระเอกของเราอยู่ตรงนี้ครับ: ระบบถอด Event เก่า
	// ==========================================
	newNode.JSEvents = make(map[string]js.Func)
	
	// 1. สั่งลบ Event อันเก่า (ที่จำว่า count = 0) ออกจาก DOM
	for eventName, oldCb := range oldNode.JSEvents {
		newNode.DOM.Call("removeEventListener", eventName, oldCb)
		oldCb.Release() 
	}
	
	// 2. เอา Event อันใหม่ (ที่จำว่า count = 1, 2, 3...) เสียบเข้าไปแทน
	for eventName, newHandler := range newNode.Events {
		cb := js.FuncOf(newHandler)
		newNode.JSEvents[eventName] = cb
		newNode.DOM.Call("addEventListener", eventName, cb)
	}
	// ==========================================

	// อัปเดต Props
	for key, val := range newNode.Props {
		if oldNode.Props[key] != val {
			if key == "value" {
				newNode.DOM.Set("value", val)
			} else {
				newNode.DOM.Call("setAttribute", key, val)
			}
		}
	}

	// ลูปเช็ค Diff ลูกๆ
	maxLen := len(newNode.Children)
	if len(oldNode.Children) > maxLen {
		maxLen = len(oldNode.Children)
	}

	for i := 0; i < maxLen; i++ {
		var oldChild, newChild *VNode
		if i < len(oldNode.Children) {
			oldChild = oldNode.Children[i]
		}
		if i < len(newNode.Children) {
			newChild = newNode.Children[i]
		}
		Diff(newNode.DOM, oldChild, newChild)
	}
}
