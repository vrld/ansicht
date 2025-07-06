package runtime

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	lua "github.com/Shopify/go-lua"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vrld/ansicht/internal/model"
)

type Runtime struct {
	luaState *lua.State
}

// Call key binding defined in config. If the binding exists, it is a function that
// will receive the given messages as arguments, e.g., `key['d'] = function(m1, m2, ...) [...] end`
func (c *Runtime) OnKey(keycode string, messages []*model.Message) tea.Cmd {
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

	// call the binding with all selected messages and return the next command
	for _, message := range messages {
		c.pushMessage(message)
	}

	c.luaState.Call(len(messages), lua.MultipleReturns) // TODO: MultipleReturns; get number as luaState.Top() - top
	count_return_values := c.luaState.Top() - top - 1
	cmd, _ := c.parseTeaCommand(count_return_values)

	return cmd
}

func (c *Runtime) pushMessage(message *model.Message) {
	c.luaState.PushUserData(message)

	lua.NewMetaTable(c.luaState, "Message")
	c.luaState.CreateTable(0, 3)
	c.luaState.PushString(string(message.ID))
	c.luaState.SetField(-2, "id")
	c.luaState.PushString(string(message.ThreadID))
	c.luaState.SetField(-2, "thread_id")
	c.luaState.PushString(string(message.Filename))
	c.luaState.SetField(-2, "filename")
	c.luaState.SetField(-2, "__index")
	c.luaState.SetMetaTable(-2)
}

// Helper function to extract field values from message userdata
func (c *Runtime) getMessageField(L *lua.State, index int, field string) (string, bool) {
	if !L.IsUserData(index) {
		return "", false
	}
	L.Field(index, field)
	value, ok := L.ToString(-1)
	L.Pop(1)
	return value, ok
}

// Go implementation of tag function
func (c *Runtime) setNotmuchTag(L *lua.State) int {
	argc := L.Top()
	if argc < 1 {
		return 0
	}

	tags, _ := L.ToString(1)
	var messageIds []string

	for i := 2; i <= argc; i++ {
		if id, ok := c.getMessageField(L, i, "id"); ok {
			messageIds = append(messageIds, "id:"+id)
		}
	}

	if len(messageIds) > 0 {
		args := []string{"tag"}
		args = append(args, strings.Split(tags, " ")...)
		args = append(args, messageIds...)
		cmd := exec.Command("notmuch", args...)
		cmd.Run()
	}

	return 0
}

func (c *Runtime) parseTeaCommand(count_return_values int) (tea.Cmd, error) {
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

const defaultConfig = `
key.q = cmd.quit
key["ctrl+c"] = key.q
key["ctrl+d"] = key.q

key["/"] = cmd.query.new
key.left = cmd.query.prev
key.right = cmd.query.next

key[" "] = cmd.selection.toggle
key["I"] = cmd.selection.invert

key["r"] = cmd.refresh

key["d"] = function(...) tag("+deleted -unread -inbox", ...) return cmd.refresh end
key["a"] = function(...) tag("+archive -inbox", ...) return cmd.refresh end

key["enter"] = function(msg)
	-- TODO: make this async
	local command = table.concat{
	  "/home/matthias/Projekte/übersicht.mail/einsicht/build/bin/einsicht",
		" ",
		msg.filename,
		">/dev/null",
		" ",
		"2>&1"
	}
	if pcall(os.execute, command) then
		tag("-unread", msg)
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
	L.Register("tag", runtime.setNotmuchTag)

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
