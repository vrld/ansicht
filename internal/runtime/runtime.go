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

type Runtime struct {
	luaState *lua.State
	Messages *service.Messages
}

// Call key binding defined in config. If the binding exists, it is an event from the
// `event` table, or a  function with no arguments that returns an event or nil, e.g.,
//
//	key.q = event.quit()  -- this is a quit event instance
//	-- these are functions that return a refresh event:
//	key.r = event.refresh
//	key['d'] = function(m1, m2, ...) [...] return event.refresh() end
func (c *Runtime) OnKey(keycode string) tea.Cmd {
	// leave stack in clean state on early exit
	top := c.luaState.Top()
	defer c.luaState.SetTop(top)

	// get binding
	c.luaState.Global("key")
	if !c.luaState.IsTable(-1) {
		return nil
	}

	c.luaState.Field(-1, keycode)

	// Userdata -> emit event
	if cmd, err := c.getTeaCommand(-1); err == nil {
		return cmd
	}

	if !c.luaState.IsFunction(-1) {
		return nil
	}

	c.luaState.Call(0, 1)
	cmd, _ := c.getTeaCommand(-1)
	return cmd
}

func (c *Runtime) getTeaCommand(index int) (tea.Cmd, error) {
	if !c.luaState.IsUserData(index) {
		return nil, fmt.Errorf("not a userdata at index: %d", index)
	}

	value := c.luaState.ToUserData(index)
	switch value.(type) {
	case tea.QuitMsg:
		return tea.Quit, nil
	case RefreshResultsMsg, QueryNewMsg, QueryNextMsg, QueryPrevMsg, MarkToggleMsg, MarkInvertMsg:
		return (func() tea.Msg { return value }), nil
	default:
		return nil, fmt.Errorf("invalid message type: %v", value)
	}
}

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

func getMessageField(L *lua.State, index int, field string) (string, bool) {
	if !isMessage(L, index) {
		return "", false
	}

	L.Field(index, field)
	value, ok := L.ToString(-1)
	L.Pop(1)

	return value, ok
}

func (c *Runtime) luaMessagesAll(L *lua.State) int {
	pushMessagesTable(L, c.Messages.GetAll())
	return 1
}

func (c *Runtime) luaMessagesSelected(L *lua.State) int {
	pushMessage(L, c.Messages.GetSelected())
	return 1
}

func (c *Runtime) luaMessagesMarked(L *lua.State) int {
	pushMessagesTable(L, c.Messages.GetMarked())
	return 1
}

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

	runtime := &Runtime{luaState: L}

	// Create key table
	L.NewTable()
	L.SetGlobal("key")

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
		{Name: "toggle", Function: luaPushMarkToggle},
		{Name: "invert", Function: luaPushMarkInvert},
	})
	L.SetField(-2, "marks")

	L.SetGlobal("event")

	// load config
	if err := lua.DoString(L, luaCode); err != nil {
		return nil, fmt.Errorf("error executing Lua config: %w", err)
	}

	return runtime, nil
}
