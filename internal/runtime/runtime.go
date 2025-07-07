package runtime

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	lua "github.com/Shopify/go-lua"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vrld/ansicht/internal/model"
	"github.com/vrld/ansicht/internal/service"
)

type Runtime struct {
	luaState *lua.State
	Messages *service.Messages
}

// Call key binding defined in config. If the binding exists, it is a function that
// will receive the given messages as arguments, e.g., `key['d'] = function(m1, m2, ...) [...] end`
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

	// string -> emit event
	if c.luaState.IsString(-1) {
		cmd, _ := c.parseTeaCommand(1)
		return cmd
	}

	if !c.luaState.IsFunction(-1) {
		return nil
	}

	c.luaState.Call(0, lua.MultipleReturns) // TODO: MultipleReturns; get number as luaState.Top() - top
	count_return_values := c.luaState.Top() - top - 1
	cmd, _ := c.parseTeaCommand(count_return_values)

	return cmd
}

func (c *Runtime) parseTeaCommand(count_return_values int) (tea.Cmd, error) {
	// TODO: these should be userdata, not string
	value, ok := c.luaState.ToString(-count_return_values) // first return value is the event name
	if !ok {
		return nil, fmt.Errorf("cannot parse value as string")
	}

	switch value {
	case "quit":
		return tea.Quit, nil

	case "query.new":
		return QueryNew, nil

	case "query.next":
		return QueryNext, nil

	case "query.prev":
		return QueryPrev, nil

	case "selection.toggle":
		return SelectionToggle, nil

	case "selection.invert":
		return SelectionInvert, nil

	case "refresh":
		// TODO: parse arguments, e.g.:
		// ids := tableToSlice(c.luaState, -count_return_values + 1, c.luaState.ToString)
		// return RefreshResults{ids: ids}, nil
		return RefreshResults, nil
	}

	return nil, fmt.Errorf("unknown signal: %s", value)
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

const defaultConfig = `
key.q = cmd.quit
key["ctrl+c"] = key.q
key["ctrl+d"] = key.q

key["/"] = cmd.query.new
key.left = cmd.query.prev
key.right = cmd.query.next

key[" "] = cmd.selection.toggle
key.i = cmd.selection.invert

key.r = cmd.refresh

local function messages_of_interest()
	local selected = messages.selected()
	local moi = { selected }
	for _, message in pairs(messages.marked()) do
		if selected ~= message then
			moi[#moi + 1] = message
		end
	end
	return moi
end

local function bind_tag_messages_of_interest(tags)
	return function()
		notmuch.tag(messages_of_interest(), table.unpack(tags))
		return cmd.refresh
	end
end

key.d = bind_tag_messages_of_interest{ "+deleted", "-unread", "-inbox" }
key.a = bind_tag_messages_of_interest{ "+archive", "-inbox" }
key.u = bind_tag_messages_of_interest{ "+unread" }

key.enter = function()
	-- TODO: async 'spawn(command, arg, arg, arg)'
	local message = messages.selected()
	local command = table.concat{
	  "/home/matthias/Projekte/übersicht.mail/einsicht/build/bin/einsicht",
		" ",
		message.filename,
		">/dev/null",
		" ",
		"2>&1"
	}
	if pcall(os.execute, command) then
		notmuch.tag(message, "-unread")
	end
	return cmd.refresh
end
`

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

	lua.NewLibrary(L, []lua.RegistryFunction{
		{Name: "all", Function: runtime.luaMessagesAll},
		{Name: "marked", Function: runtime.luaMessagesMarked},
		{Name: "selected", Function: runtime.luaMessagesSelected},
	})
	L.SetGlobal("messages")

	// TODO: make these userdata of the actual messages
	lua.DoString(L, `cmd = {
		query = {
			new = "query.new",
			next = "query.next",
			prev = "query.prev",
		},
		quit = "quit",
		refresh = "refresh",
		selection = {
			toggle = "selection.toggle",
			invert = "selection.invert",
		},
	}`)

	if err := lua.DoString(L, luaCode); err != nil {
		return nil, fmt.Errorf("error executing Lua config: %w", err)
	}

	return runtime, nil
}
