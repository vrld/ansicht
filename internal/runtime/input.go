package runtime

import (
	"fmt"

	"github.com/Shopify/go-lua"
	tea "github.com/charmbracelet/bubbletea"
)

type InputMsg struct {
	Placeholder string
	Prompt      string
	Handle      string
}


// event.input{ placeholder = 'string', with_input = function(input) return event end}
func luaPushInput(L *lua.State, handleId int) int {
	if L.Top() < 1 || !L.IsTable(1) {
		lua.Errorf(L, "missing table argument")
		panic("unreachable")
	}

	lFieldStringOrNil(L, 1, "placeholder")
	placeholder, _ := L.ToString(-1)

	lFieldStringOrNil(L, 1, "prompt")
	prompt, _ := L.ToString(-1)

	handle := fmt.Sprintf("ansicht.input_callback_handle_%d", handleId)

	L.PushString(handle)
	lFieldFunctionOrNil(L, 1, "with_input")
	L.SetTable(lua.RegistryIndex)

	L.PushUserData(InputMsg{
		Placeholder: placeholder,
		Prompt:      prompt,
		Handle:      handle,
	})
	return 1
}

func (r *Runtime) PushInputHandle(handle string) {
	r.inputCallbackStack = append(r.inputCallbackStack, InputCallbackHandle(handle))
}

func (r *Runtime) HandleInput(input string) tea.Cmd {
	nHandles := len(r.inputCallbackStack)
	if nHandles == 0 {
		return nil
	}

	callbackHandle := r.inputCallbackStack[nHandles-1]
	if nHandles == 1 {
		r.inputCallbackStack = nil
	} else {
		r.inputCallbackStack = r.inputCallbackStack[:nHandles-1]
	}

	r.luaState.PushString(string(callbackHandle))
	r.luaState.Table(lua.RegistryIndex)

	if r.luaState.TypeOf(-1) == lua.TypeFunction {
		r.luaState.PushString(input)
		r.luaState.Call(1, 1)
		cmd, _ := r.getTeaCommand(-1)

		r.luaState.PushString(string(callbackHandle))
		r.luaState.PushNil()
		r.luaState.SetTable(lua.RegistryIndex)

		return cmd
	}

	return nil
}
