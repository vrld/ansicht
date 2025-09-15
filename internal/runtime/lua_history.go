package runtime

import (
	"os"
	"path/filepath"

	"github.com/Shopify/go-lua"
	"github.com/vrld/ansicht/internal/service"
)

func luaHistoryClear(L *lua.State) int {
	if L.Top() == 0 {
		service.InputHistory().ClearAll()
	} else {
		service.InputHistory().Clear(lua.CheckString(L, 1))
	}
	return 0
}

func luaHistoryCount(L *lua.State) int {
	L.PushInteger(service.InputHistory().Count(lua.CheckString(L, 1)))
	return 1
}

func luaHistoryAdd(L *lua.State) int {
	prompt := lua.CheckString(L, 1)
	service.InputHistory().Add(prompt, lua.CheckString(L, 2))
	return 0
}

func luaHistoryGet(L *lua.State) int {
	prompt := lua.CheckString(L, 1)
	if !L.IsNil(2) {
		L.PushString(service.InputHistory().Get(prompt, lua.CheckInteger(L, 2)-1))
		return 1
	}

	count := service.InputHistory().Count(prompt)
	for i := range count {
		L.PushString(service.InputHistory().Get(prompt, i))
	}
	return count
}

func luaHistoryLoad(L *lua.State) int {
	xdgCacheHome := findXdgDir("XDG_CACHE_HOME", ".cache")
	if xdgCacheHome == "" {
		panic("cannot determine XDG_CACHE_HOME")
	}

	historyPath := filepath.Join(xdgCacheHome, "ansicht", "history")
	err := service.InputHistory().Load(historyPath)
	if err != nil {
		L.PushString(err.Error())
		return 1
	}
	return 0
}

func luaHistorySave(L *lua.State) int {
	xdgCacheHome := findXdgDir("XDG_CACHE_HOME", ".cache")
	if xdgCacheHome == "" {
		panic("cannot determine XDG_CACHE_HOME")
	}

	os.Mkdir(filepath.Join(xdgCacheHome, "ansicht"), 0600)
	historyPath := filepath.Join(xdgCacheHome, "ansicht", "history")
	err := service.InputHistory().Save(historyPath)
	if err != nil {
		L.PushString(err.Error())
		return 1
	}
	return 0
}
