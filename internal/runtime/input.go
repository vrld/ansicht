package runtime

import (
	"fmt"

	"github.com/Shopify/go-lua"
)

func inputCallbackHandleString(id int) string {
	return fmt.Sprintf("ansicht.input_callback_handle_%d", id)
}

// event.input{ placeholder = 'string', with_input = function(input) return event end}
func (r *Runtime) luaInput(L *lua.State) int {
	if L.Top() < 1 || !L.IsTable(1) {
		lua.Errorf(L, "missing table argument")
		panic("unreachable")
	}

	lFieldStringOrNil(L, 1, "placeholder")
	placeholder, _ := L.ToString(-1)

	lFieldStringOrNil(L, 1, "prompt")
	prompt, _ := L.ToString(-1)

	r.countOpenInputs++
	L.PushString(inputCallbackHandleString(r.countOpenInputs))
	lFieldFunctionOrNil(L, 1, "with_input")
	L.SetTable(lua.RegistryIndex)

	r.Controller.Input(prompt, placeholder)
	return 0
}

func (r *Runtime) HandleInput(input string) {
	if r.countOpenInputs <= 0 {
		return
	}

	handle := inputCallbackHandleString(r.countOpenInputs)
	r.luaState.PushString(handle)
	r.luaState.Table(lua.RegistryIndex)

	if r.luaState.TypeOf(-1) == lua.TypeFunction {
		r.luaState.PushString(input)
		r.luaState.Call(1, 0)

		r.luaState.PushString(handle)
		r.luaState.PushNil()
		r.luaState.SetTable(lua.RegistryIndex)
	}
}
