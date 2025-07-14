package runtime

import (
	"fmt"

	"github.com/Shopify/go-lua"
	tea "github.com/charmbracelet/bubbletea"
)

type RefreshResultsMsg struct{}
type QueryNewMsg struct {
	Query string
}
type QueryNextMsg struct{}
type QueryPrevMsg struct{}
type MarksToggleMsg struct{}
type MarksInvertMsg struct{}
type MarksClearMsg struct{}

type InputMsg struct {
	Placeholder string
	Handle      string
}

func luaPushRefresh(L *lua.State) int {
	L.PushUserData(RefreshResultsMsg{})
	return 1
}

func luaPushQuit(L *lua.State) int {
	L.PushUserData(tea.QuitMsg{})
	return 1
}

func luaPushQueryNew(L *lua.State) int {
	if L.Top() < 1 || !L.IsString(1) {
		lua.Errorf(L, "missing string argument")
		panic("unreachable")
	}
	query, _ := L.ToString(1)
	L.PushUserData(QueryNewMsg{Query: query})
	return 1
}

func luaPushQueryNext(L *lua.State) int {
	L.PushUserData(QueryNextMsg{})
	return 1
}

func luaPushQueryPrev(L *lua.State) int {
	L.PushUserData(QueryPrevMsg{})
	return 1
}

func luaPushMarksToggle(L *lua.State) int {
	L.PushUserData(MarksToggleMsg{})
	return 1
}

func luaPushMarksInvert(L *lua.State) int {
	L.PushUserData(MarksInvertMsg{})
	return 1
}

func luaPushMarksClear(L *lua.State) int {
	L.PushUserData(MarksClearMsg{})
	return 1
}

// event.input{ placeholder = 'string', with_input = function(input) return event end}
func luaPushInput(L *lua.State, handleId int) int {
	if L.Top() < 1 || !L.IsTable(1) {
		lua.Errorf(L, "missing table argument")
		panic("unreachable")
	}

	L.Field(1, "placeholder")
	if !(L.IsString(-1) || L.IsNil(-1)) {
		lua.Errorf(L, "placeholder must be a string or nil")
		panic("unreachable")
	}
	placeholder, _ := L.ToString(-1)

	handle := fmt.Sprintf("ansicht.input_callback_handle_%d", handleId)

	L.PushString(handle)
	L.Field(1, "with_input")
	if !(L.IsFunction(-1) || L.IsNil(-1)) {
		lua.Errorf(L, "with_input must be a function or nil")
		panic("unreachable")
	}
	L.SetTable(lua.RegistryIndex)

	L.PushUserData(InputMsg{
		Placeholder: placeholder,
		Handle:      handle,
	})
	return 1
}
