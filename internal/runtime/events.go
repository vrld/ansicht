package runtime

import (
	"github.com/Shopify/go-lua"
	tea "github.com/charmbracelet/bubbletea"
)

type RefreshResultsMsg struct{} // TODO: add args: which messages? => reqires more parsing
type QueryNewMsg struct{}
type QueryNextMsg struct{}
type QueryPrevMsg struct{}
type MarkToggleMsg struct{}
type MarkInvertMsg struct{}

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

func luaPushMarkToggle(L *lua.State) int {
	L.PushUserData(MarkToggleMsg{})
	return 1
}

func luaPushMarkInvert(L *lua.State) int {
	L.PushUserData(MarkInvertMsg{})
	return 1
}
