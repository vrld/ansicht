package runtime

import (
	"fmt"
	"os"
	"path/filepath"

	lua "github.com/Shopify/go-lua"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vrld/ansicht/internal/model"
)

type RefreshResultsMsg struct{} // TODO: add args: which messages? => reqires more parsing
func RefreshResults() tea.Msg   { return RefreshResultsMsg{} }

// navigation
type QueryNewMsg struct{}
type QueryNextMsg struct{}
type QueryPrevMsg struct{}

func QueryNew() tea.Msg  { return QueryNewMsg{} }
func QueryNext() tea.Msg { return QueryNextMsg{} }
func QueryPrev() tea.Msg { return QueryPrevMsg{} }

// selection
type SelectionToggleMsg struct{}
type SelectionInvertMsg struct{}

func SelectionToggle() tea.Msg { return SelectionToggleMsg{} }
func SelectionInvert() tea.Msg { return SelectionInvertMsg{} }

type Config struct {
	luaState *lua.State
}

// Call key binding defined in config. If the binding exists, it is a function that
// will receive the given messages as arguments, e.g., `key['d'] = function(m1, m2, ...) [...] end`
func (c *Config) OnKey(keycode string, messages []*model.Message) tea.Cmd {
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

func (c *Config) pushMessage(message *model.Message) {
	c.luaState.CreateTable(0, 3)

	c.luaState.PushString(string(message.ID))
	c.luaState.SetField(-2, "id")

	c.luaState.PushString(string(message.ThreadID))
	c.luaState.SetField(-2, "thread_id")

	c.luaState.PushString(string(message.Filename))
	c.luaState.SetField(-2, "filename")
}

func (c *Config) parseTeaCommand(count_return_values int) (tea.Cmd, error) {
	value, ok := c.luaState.ToString(-count_return_values)  // first return value is the event name
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
	local command = std.concat_sep(
		" ",
	  "/home/matthias/Projekte/übersicht.mail/einsicht/build/bin/einsicht",
		msg.filename,
		">/dev/null",
		"2>&1"
	)
	if pcall(os.execute, command) then
		tag("-unread", msg)
	end
	return cmd.refresh
end
`

func LoadConfig() (*Config, error) {
	// Try XDG_CONFIG_HOME first
	if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
		content, err := os.ReadFile(filepath.Join(xdgConfigHome, "ansicht", "init.lua"))
		if err == nil {
			return ConfigFromString(string(content))
		}
	}

	if home, err := os.UserHomeDir(); err == nil {
		// maybe XDG_CONFIG_HOME was not set?
		content, err := os.ReadFile(filepath.Join(home, ".config", "ansicht", "init.lua"))
		if err == nil {
			return ConfigFromString(string(content))
		}
	}

	// no user config
	return ConfigFromString(defaultConfig)
}

func ConfigFromString(luaCode string) (*Config, error) {
	L := lua.NewState()
	lua.OpenLibraries(L)

	L.NewTable()
	L.SetGlobal("key")

	// TODO: consolidate runtime: what should be a Go closure?
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

	lua.DoString(L, `function tag(tags, ...)
		local shell_command = std.concat_sep(" ", "notmuch tag", tags, std.ids(...))
		os.execute(shell_command)
	end

	std = {}

	local function map_helper(f, i, ...)
		local N = select('#', ...)
		if i >= N then
			return f(select(N, ...))
		end
		return f(select(i, ...)), map_helper(f, i + 1, ...)
	end

	std.ids = function(...)
		local id = function(m) return "id:" .. m.id end
		return map_helper(id, 1, ...)
	end

	std.thread_ids = function(...)
		local thread = function(m) return "thread:" .. m.thread_id end
		return map_helper(thread, 1, ...)
	end

	std.filenames = function(...)
		local filename = function(m) return m.filename end
		return map_helper(filename, 1, ...)
	end

	local function concat_helper(accum, sep, s, ...)
		if select('#', ...) == 0 then
			accum[#accum + 1] = s
			return table.concat(accum)
		end
		accum[#accum + 1] = s
		if sep ~= nil then
			accum[#accum + 1] = sep
		end
		return concat_helper(accum, sep, ...)
	end
	std.concat = function(...) return concat_helper({}, nil, ...) end
	std.concat_sep = function(sep, ...) return concat_helper({}, sep, ...) end
	`)

	if err := lua.DoString(L, luaCode); err != nil {
		return nil, fmt.Errorf("error executing Lua config: %w", err)
	}

	return &Config{luaState: L}, nil
}
