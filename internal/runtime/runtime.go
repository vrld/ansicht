package runtime

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	_ "embed"

	lua "github.com/Shopify/go-lua"
	"github.com/vrld/ansicht/internal/service"
)

type Runtime struct {
	luaState        *lua.State
	countOpenInputs int
	Controller      ControllerAdapter
}

//go:embed default_config.lua
var defaultConfig string

func LoadRuntime() (*Runtime, error) {
	// Try XDG_CONFIG_HOME first
	if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
		content, err := os.ReadFile(filepath.Join(xdgConfigHome, "ansicht", "init.lua"))
		if err == nil {
			return runtimeFromString(string(content))
		}
	}

	if home, err := os.UserHomeDir(); err == nil {
		// maybe XDG_CONFIG_HOME was not set?
		content, err := os.ReadFile(filepath.Join(home, ".config", "ansicht", "init.lua"))
		if err == nil {
			return runtimeFromString(string(content))
		}
	}

	// no user config
	return runtimeFromString(defaultConfig)
}

func runtimeFromString(luaCode string) (*Runtime, error) {
	L := lua.NewState()
	lua.OpenLibraries(L)

	runtime := &Runtime{luaState: L, Controller: &NullAdapter{}}

	// Create key table
	L.NewTable()
	L.SetGlobal("key")

	lua.NewLibrary(L, []lua.RegistryFunction{
		{Name: "quit", Function: runtime.luaQuit},
		{Name: "refresh", Function: runtime.luaRefresh},
		{Name: "spawn", Function: runtime.luaSpawn},
		{Name: "tag", Function: luaNotmuchTag},
		{Name: "input", Function: runtime.luaInput},
		{Name: "set_theme", Function: runtime.luaSetTheme},
	})

	// status
	lua.NewLibrary(L, []lua.RegistryFunction{
		{Name: "set", Function: runtime.luaStatusSet},
		{Name: "get", Function: runtime.luaStatusGet},
	})
	L.SetField(-2, "status")

	// messages access
	lua.NewLibrary(L, []lua.RegistryFunction{
		{Name: "all", Function: runtime.luaMessagesAll},
		{Name: "marked", Function: runtime.luaMessagesMarked},
		{Name: "selected", Function: runtime.luaMessagesSelected},
	})
	L.SetField(-2, "messages")

	// query subgroup
	lua.NewLibrary(L, []lua.RegistryFunction{
		{Name: "new", Function: runtime.luaQueryNew},
		{Name: "next", Function: runtime.luaQuerySelectNext},
		{Name: "prev", Function: runtime.luaQuerySelectPrev},
	})
	L.SetField(-2, "query")

	// marks subgroup
	lua.NewLibrary(L, []lua.RegistryFunction{
		{Name: "toggle", Function: runtime.luaMarksToggle},
		{Name: "invert", Function: runtime.luaMarksInvert},
		{Name: "clear", Function: runtime.luaMarksClear},
	})
	L.SetField(-2, "marks")

	// log.<level>(message)  =>  real-log(LEVEL, message)
	lua.NewLibrary(L, []lua.RegistryFunction{
		{Name: "__index", Function: runtime.luaLogMetatableIndex},
	})
	L.PushValue(-1)
	L.SetMetaTable(-2)
	L.SetField(-2, "log")

	L.SetGlobal("ansicht")

	// load config
	if err := lua.DoString(L, luaCode); err != nil {
		return nil, fmt.Errorf("error executing Lua config: %w", err)
	}

	return runtime, nil
}

func (r *Runtime) OnStartup() {
	top := r.luaState.Top()
	defer r.luaState.SetTop(top)

	r.luaState.Global("Startup")
	if r.luaState.IsFunction(-1) {
		r.luaState.Call(0, 0)
	}
}

// Call key binding defined in config. If the binding exists, it must be a function
// that expects no arguments:
//
//	key.q = ansicht.quit
//	key.r = ansicht.refresh
//	key['d'] = function() ansicht.tag(ansicht.messages.selected(), "+deleted") end
func (r *Runtime) OnKey(keycode string) bool {
	// leave stack in clean state on early exit
	top := r.luaState.Top()
	defer r.luaState.SetTop(top)

	r.luaState.Global("key")
	if !r.luaState.IsTable(-1) {
		lua.Errorf(r.luaState, "Table `key` not found. Check your config.")
		panic("unreachable")
	}
	r.luaState.Field(-1, keycode)

	if r.luaState.IsNil(-1) {
		return false
	}

	if r.luaState.IsFunction(-1) {
		r.luaState.Call(0, 0)
		return true
	}

	lua.Errorf(r.luaState, "key['%s'] must be nil or a function.", keycode)
	panic("unreachable")
}

// Tag one or more messages
// notmuch.tag(message, tag1, tag2, ..., tag3)
// notmuch.tag({messages}, tag1, tag2, ..., tag3)
// equivalent to: `notmuch tag tag1 tag2 tag3 id:... id:... ...`
func luaNotmuchTag(L *lua.State) int {
	argc := L.Top()
	if argc < 1 {
		lua.Errorf(L, "invalid arguments")
		panic("unreachable")
	}

	var messageIds []string

	if isMessage(L, 1) {
		if id, ok := getMessageField(L, 1, "id"); ok {
			messageIds = append(messageIds, "id:"+id)
		}
	} else if L.IsTable(1) {
		count := L.RawLength(1)
		for i := 1; i <= count; i++ {
			L.RawGetInt(1, i)
			if id, ok := getMessageField(L, -1, "id"); ok {
				messageIds = append(messageIds, "id:"+id)
			}
			L.Pop(1)
		}
	} else {
		lua.Errorf(L, "Neither a message nor a table of messages")
		panic("unreachable")
	}

	if len(messageIds) > 0 {
		args := []string{"tag"}
		for i := 2; i <= L.Top(); i++ {
			if tag, ok := L.ToString(i); ok {
				args = append(args, tag)
			}
		}

		args = append(args, messageIds...)

		cmd := exec.Command("notmuch", args...)
		cmd.Run()
	}

	return 0
}

// returns the current status message
func (r *Runtime) luaStatusGet(L *lua.State) int {
	L.PushString(service.Status().Get())
	return 1
}

// returns closure that calls log()
// meant to be used as __index function
// effectively:
// ansicht.log.__index = function(_, level)
//
//	return function(message) log(level, message) end
//
// end
func (r *Runtime) luaLogMetatableIndex(L *lua.State) int {
	if key, ok := L.ToString(2); ok {
		key = strings.ToUpper(key)
		L.PushString(key)
		L.PushGoClosure(r.luaLogMessage, 1)
		return 1
	}

	service.Logger().Error("luaLogMetatableIndex: not a string")
	lua.Errorf(L, "what are you doing?")
	panic("unreachable")
}

// clearly better than defining separate functions
func (r *Runtime) luaLogMessage(L *lua.State) int {
	level, _ := L.ToString(lua.UpValueIndex(1))
	message, _ := lua.ToStringMeta(L, 1)
	service.Logger().Log(service.LogLevel(level), message)
	return 0
}
