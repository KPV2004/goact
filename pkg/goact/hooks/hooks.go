// pkg/goact/hooks/hooks.go
package hooks

// import (
// 	"fmt"
// )

// Internal state storage for the current rendering component.
var (
	// states stores all hook data sequentially.
	// It uses 'any' (Go 1.18+) to hold different types like int, string, slices.
	states []any
	
	// hookIndex keeps track of the current hook being processed during render.
	// CRITICAL: Hooks MUST be called in the exact same order every render.
	hookIndex int
	
	// updateFunc is a registered callback function that triggers a full
	// application re-render. This links state changes to the UI update cycle.
	updateFunc func()
)

// Init registers the main rendering function from the Goact core.
// This is called once during application Mount.
func Init(renderFn func()) {
	updateFunc = renderFn
}

// Reset moves the hookIndex pointer back to zero.
// This MUST be called by the Goact core immediately before every render cycle
// starts, ensuring UseState can read the correct data from the states array.
func Reset() {
	hookIndex = 0
}

// UseState is a Generic Hook that lets you add state to functional components.
//
// It returns two values:
// 1. The current state value (of type T).
// 2. A setter function (func(T)) to update the state and trigger a re-render.
//
// The value must be identifiable by its type T. If it's the first render,
// initialValue is stored. Subsequent renders return the stored value.
func UseState[T any](initialValue T) (T, func(T)) {
	// 1. Obtain the current index and increment for the next hook call.
	currentIndex := hookIndex
	hookIndex++

	// 2. Handle Initial Mount:
	// If the index exceeds current storage length, it's the first render.
	// We append the initial value to our flat states array.
	if currentIndex >= len(states) {
		states = append(states, initialValue)
	}

	// 3. Handle Re-render / Retrieve Data:
	// Cast the stored value (any) back to its specific type (T).
	// This will panic if T does not match the stored type (due to incorrect call order).
	state := states[currentIndex].(T)

	// 4. Create the State Setter (func(T)):
	// This is a closure that remembers which state object it is targeting
	// thanks to 'currentIndex'.
	setState := func(newValue T) {
		// Update the value in the flat storage array.
		states[currentIndex] = newValue
		
		// Log for debugging (optional).
		// fmt.Printf("Hooks: State updated at index %d, triggering re-render.\n", currentIndex)

		// Trigger the global update function to start a new VDOM Diff cycle.
		if updateFunc != nil {
			updateFunc()
		}
	}

	// 5. Return current value and setter.
	return state, setState
}
