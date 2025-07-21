package runtime

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"

	lua "github.com/Shopify/go-lua"
	tea "github.com/charmbracelet/bubbletea"
)

type SpawnResultMsg struct {
	Command    []string
	ReturnCode int
	Stdout     string
	Stderr     string
	HandleID   int
	Timeout    bool
}

func (r *Runtime) HandleSpawnResult(msg SpawnResultMsg) tea.Cmd {
	handleID := msg.HandleID

	r.luaState.PushString(completeHandleKey(handleID))
	r.luaState.Table(lua.RegistryIndex)

	if r.luaState.TypeOf(-1) == lua.TypeFunction {
		// Push spawn result as arguments to callback
		r.luaState.CreateTable(0, 4)

		r.luaState.PushString("command")
		lPushStringTable(r.luaState, msg.Command)
		r.luaState.SetTable(-3)

		if !msg.Timeout {
			lSetFieldInteger(r.luaState, -1, "return_code", msg.ReturnCode)
		}
		lSetFieldBool(r.luaState, -1, "timeout", msg.Timeout)
		lSetFieldString(r.luaState, -1, "stdout", msg.Stdout)
		lSetFieldString(r.luaState, -1, "stderr", msg.Stderr)

		r.luaState.Call(1, 1)
		cmd, _ := r.getTeaCommand(-1)

		// Clean up both callbacks using derived keys
		lSetFieldNil(r.luaState, lua.RegistryIndex, completeHandleKey(handleID))

		return cmd
	}

	return nil
}

func completeHandleKey(handleId int) string {
	return fmt.Sprintf("ansicht.spawn_complete_callback_handle_%d", handleId)
}

// spawn{"command", "arg1", "arg2", ..., timeout=60, on_complete=function, on_timeout=function}
var spawnHandleId int

func (r *Runtime) luaSpawn(L *lua.State) int {
	if L.Top() < 1 || !L.IsTable(1) {
		lua.Errorf(L, "spawn expects a table argument")
		panic("unreachable")
	}

	// Extract command and arguments from array part
	var command []string
	count := L.RawLength(1)
	if count == 0 {
		lua.Errorf(L, "spawn requires at least one command argument")
		panic("unreachable")
	}

	for i := 1; i <= count; i++ {
		L.RawGetInt(1, i)
		if arg, ok := L.ToString(-1); ok {
			command = append(command, arg)
		}
		L.Pop(1)
	}

	if len(command) != count {
		lua.Errorf(L, "all command arguments must be strings")
		panic("unreachable")
	}

	// Extract timeoutMilliseconds from table
	var timeoutMilliseconds int
	L.Field(1, "timeout")
	if L.IsNumber(-1) {
		if timeoutSecs, ok := L.ToNumber(-1); ok && timeoutSecs > 0 {
			timeoutMilliseconds = int(timeoutSecs * 1000)
		}
	}
	L.Pop(1)

	// Register callbacks in the registry
	spawnHandleId++
	L.PushString(completeHandleKey(spawnHandleId))
	lFieldFunctionOrNil(L, 1, "next")
	L.SetTable(lua.RegistryIndex)

	go r.spawnCommand(command, time.Duration(timeoutMilliseconds)*time.Millisecond, spawnHandleId)

	return 0
}

func (r *Runtime) spawnCommand(command []string, timeout time.Duration, handleID int) {
	if len(command) == 0 {
		// TODO: send error event
		return
	}

	var ctx context.Context
	var cancel context.CancelFunc

	if timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	defer cancel()

	cmd := exec.CommandContext(ctx, command[0], command[1:]...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if ctx.Err() == context.DeadlineExceeded {
		r.SendMessage(SpawnResultMsg{
			Command:  command,
			Stdout:   stdout.String(),
			Stderr:   stderr.String(),
			HandleID: handleID,
			Timeout:  true,
		})
		return
	}

	returnCode := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			returnCode = exitError.ExitCode()
		} else {
			returnCode = -1
		}
	}

	r.SendMessage(SpawnResultMsg{
		Command:    command,
		ReturnCode: returnCode,
		Stdout:     stdout.String(),
		Stderr:     stderr.String(),
		HandleID:   handleID,
		Timeout:    false,
	})
}
