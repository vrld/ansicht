package runtime

import (
	"github.com/Shopify/go-lua"
	tea "github.com/charmbracelet/bubbletea"
)

type RefreshResultsMsg struct{} // TODO: add args: which messages? => reqires more parsing
type QueryNewMsg struct{}
type QueryNextMsg struct{}
type QueryPrevMsg struct{}
type MarksToggleMsg struct{}
type MarksInvertMsg struct{}
type MarksClearMsg struct{}

func luaPushRefresh(L *lua.State) int {
	L.PushUserData(RefreshResultsMsg{})
	return 1
}

func luaPushQuit(L *lua.State) int {
	L.PushUserData(tea.QuitMsg{})
	return 1
}

func luaPushQueryNew(L *lua.State) int {
	L.PushUserData(QueryNewMsg{})
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
