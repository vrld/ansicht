package runtime

import (
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
type StatusSetMsg struct {
	Message string
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

func luaPushStatusSet(L *lua.State) int {
	if L.Top() < 1 || !L.IsString(1) {
		lua.Errorf(L, "missing string argument")
		panic("unreachable")
	}
	message, _ := L.ToString(1)
	L.PushUserData(StatusSetMsg{Message: message})
	return 1
}
