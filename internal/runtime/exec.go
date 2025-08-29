package runtime

import (
	"os/exec"
	"strings"

	lua "github.com/Shopify/go-lua"
)

const luaRegistryExecCompletHandle = "ansicht.exec_complete_handle"

func (r *Runtime) luaExec(L *lua.State) int {
	if L.Top() < 1 || !L.IsTable(1) {
		lua.Errorf(L, "exec expects a table argument")
		panic("unreachable")
	}

	// Extract command and arguments from array part
	count := L.RawLength(1)
	if count == 0 {
		lua.Errorf(L, "exec requires at least one command argument")
		panic("unreachable")
	}

	var command []string
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

	cmd := exec.Command(command[0], command[1:]...)

	if str, ok := lFieldString(L, 1, "stdin"); ok {
		cmd.Stdin = strings.NewReader(str)
	}

	L.PushString(luaRegistryExecCompletHandle)
	lFieldFunctionOrNil(L, 1, "next")
	L.SetTable(lua.RegistryIndex)

	r.Controller.Exec(cmd)

	return 0
}

func (r *Runtime) OnExecCommandResult(err error) {
	r.luaState.PushString(luaRegistryExecCompletHandle)
	r.luaState.Table(lua.RegistryIndex)

	if r.luaState.TypeOf(-1) == lua.TypeFunction {
		if err != nil {
			r.luaState.PushString(err.Error())
		} else {
			r.luaState.PushNil()
		}
		r.luaState.Call(1, 0)

		lSetFieldNil(r.luaState, lua.RegistryIndex, luaRegistryExecCompletHandle)
	}
}
