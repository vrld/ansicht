package runtime

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	_ "embed"

	lua "github.com/Shopify/go-lua"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vrld/ansicht/internal/model"
	"github.com/vrld/ansicht/internal/service"
)

type InputCallbackHandle string

type Runtime struct {
	luaState           *lua.State
	Messages           *service.Messages
	Status             *service.Status
	inputCallbackStack []InputCallbackHandle
	SendMessage        func(tea.Msg)
}

//go:embed default_config.lua
var defaultConfig string

//go:embed lua_runtime.lua
var luaRuntime string

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

	runtime := &Runtime{luaState: L, SendMessage: func(tea.Msg) {}}

	// Create key table
	L.NewTable()
	L.SetGlobal("key")

	L.PushGoFunction(runtime.luaSpawn)
	L.SetGlobal("spawn")

	// Register tag function globally
	lua.NewLibrary(L, []lua.RegistryFunction{
		{Name: "tag", Function: luaNotmuchTag},
	})
	L.SetGlobal("notmuch")

	// messages access
	lua.NewLibrary(L, []lua.RegistryFunction{
		{Name: "all", Function: runtime.luaMessagesAll},
		{Name: "marked", Function: runtime.luaMessagesMarked},
		{Name: "selected", Function: runtime.luaMessagesSelected},
	})
	L.SetGlobal("messages")

	// commands / events
	lua.NewLibrary(L, []lua.RegistryFunction{
		{Name: "refresh", Function: luaPushRefresh},
		{Name: "quit", Function: luaPushQuit},
		{Name: "input", Function: func(L *lua.State) int {
			return luaPushInput(L, len(runtime.inputCallbackStack))
		}},
		{Name: "status", Function: luaPushStatusSet},
	})

	// query subgroup
	lua.NewLibrary(L, []lua.RegistryFunction{
		{Name: "new", Function: luaPushQueryNew},
		{Name: "next", Function: luaPushQueryNext},
		{Name: "prev", Function: luaPushQueryPrev},
	})
	L.SetField(-2, "query")

	// marks subgroup
	lua.NewLibrary(L, []lua.RegistryFunction{
		{Name: "toggle", Function: luaPushMarksToggle},
		{Name: "invert", Function: luaPushMarksInvert},
		{Name: "clear", Function: luaPushMarksClear},
	})
	L.SetField(-2, "marks")

	L.SetGlobal("event")

	// status functions
	lua.NewLibrary(L, []lua.RegistryFunction{
		{Name: "get", Function: runtime.luaStatusGet},
	})
	L.SetGlobal("status")

	if err := lua.DoString(L, luaRuntime); err != nil {
		panic(err)
	}

	// load config
	if err := lua.DoString(L, luaCode); err != nil {
		return nil, fmt.Errorf("error executing Lua config: %w", err)
	}

	return runtime, nil
}

func (r *Runtime) OnStartup() tea.Cmd {
	top := r.luaState.Top()
	defer r.luaState.SetTop(top)

	r.luaState.Global("Startup")
	if !r.luaState.IsFunction(-1) {
		return nil
	}

	r.luaState.Call(0, 1)
	cmd, _ := r.getTeaCommand(-1)

	// TODO: load style and other config here

	return cmd
}

// Call key binding defined in config. If the binding exists, it is an event from the
// `event` table, or a  function with no arguments that returns an event or nil, e.g.,
//
//	key.q = event.quit()  -- this is a quit event instance
//	-- these are functions that return a refresh event:
//	key.r = event.refresh
//	key['d'] = function(m1, m2, ...) [...] return event.refresh() end
func (r *Runtime) OnKey(keycode string) tea.Cmd {
	// leave stack in clean state on early exit
	top := r.luaState.Top()
	defer r.luaState.SetTop(top)

	r.luaState.Global("key")
	if !r.luaState.IsTable(-1) {
		lua.Errorf(r.luaState, "Table `key` not found. Check your config.")
		panic("unreachable")
	}
	r.luaState.Field(-1, keycode)

	cmd, _ := r.getTeaCommand(-1)
	return cmd
}

func (r *Runtime) getTeaCommand(index int) (tea.Cmd, error) {
	// execute functions until we get to the userdata
	// this makes it possible to wrap events in multiple functions, see lua_runtime.lua
	for r.luaState.IsFunction(-1) {
		// TODO: exit when nesting is too deep?
		r.luaState.Call(0, 1)
	}

	// userdata => command
	if r.luaState.IsUserData(index) {
		value := r.luaState.ToUserData(index)

		// NOTE: invalid events will be ignored by the event loop, so we don't need to filter
		return (func() tea.Msg { return value }), nil
	}

	// table of events => batch of commandss
	if r.luaState.IsTable(index) {
		var cmds []tea.Cmd
		count := r.luaState.RawLength(index)
		for i := 1; i <= count; i++ {
			r.luaState.RawGetInt(index, i)
			if cmd, err := r.getTeaCommand(-1); err == nil && cmd != nil {
				cmds = append(cmds, cmd)
			}
			r.luaState.Pop(1)
		}

		if len(cmds) > 0 {
			return tea.Batch(cmds...), nil
		}
	}

	return nil, fmt.Errorf("not a userdata or table at index: %d", index)
}

// lua api
// notmuch.tag({messages}, tag1, tag2, ..., tag3)
// notmuch.tag(message, tag1, tag2, ..., tag3)
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

// returns message[field] where message is the message userdata at `index` on the stack
// converts objects to string according to Lua rules
func getMessageField(L *lua.State, index int, field string) (string, bool) {
	if !isMessage(L, index) {
		return "", false
	}

	L.Field(index, field)
	value, ok := L.ToString(-1)
	L.Pop(1)

	return value, ok
}

// put all messages on the stack
func (r *Runtime) luaMessagesAll(L *lua.State) int {
	pushMessagesTable(L, r.Messages.GetAll())
	return 1
}

// put selected/highligted message on the stack
func (r *Runtime) luaMessagesSelected(L *lua.State) int {
	pushMessage(L, r.Messages.GetSelected())
	return 1
}

// put marked messages on the stack
func (r *Runtime) luaMessagesMarked(L *lua.State) int {
	pushMessagesTable(L, r.Messages.GetMarked())
	return 1
}

// pushes a table of messages on the stack
func pushMessagesTable(L *lua.State, messages []*model.Message) {
	L.CreateTable(len(messages), 0)
	for i, msg := range messages {
		pushMessage(L, msg)
		L.RawSetInt(-2, i+1)
	}
}

func pushMessage(L *lua.State, message *model.Message) int {
	L.PushUserData(message)

	L.CreateTable(0, 5)
	L.PushString("ansicht.Message")
	L.SetField(-2, "__name")

	L.PushString(string(message.ID))
	L.SetField(-2, "id")

	L.PushString(string(message.ThreadID))
	L.SetField(-2, "thread_id")

	L.PushString(string(message.Filename))
	L.SetField(-2, "filename")

	// metatable.__index = metatable
	L.PushValue(-1)
	L.SetField(-2, "__index")

	L.SetMetaTable(-2)

	return 1
}

func isMessage(L *lua.State, index int) bool {
	if !L.IsUserData(index) {
		return false
	}

	// leave a tidy stack
	top := L.Top()
	defer L.SetTop(top)

	// check metatable.__name
	if !L.MetaTable(index) {
		return false
	}

	L.Field(-1, "__name")
	if name, ok := L.ToString(-1); ok {
		return name == "ansicht.Message"
	}

	return false
}

// returns the current status message
func (r *Runtime) luaStatusGet(L *lua.State) int {
	L.PushString(r.Status.Get())
	return 1
}
