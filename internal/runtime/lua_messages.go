package runtime

import (
	lua "github.com/Shopify/go-lua"
	"github.com/vrld/ansicht/internal/model"
	"github.com/vrld/ansicht/internal/service"
)

// put all messages on the stack
func (r *Runtime) luaMessagesAll(L *lua.State) int {
	pushMessagesTable(L, service.Messages().GetAll())
	return 1
}

// put selected/highligted message on the stack
func (r *Runtime) luaMessagesSelected(L *lua.State) int {
	pushMessage(L, service.Messages().GetSelected())
	return 1
}

// put marked messages on the stack
func (r *Runtime) luaMessagesMarked(L *lua.State) int {
	pushMessagesTable(L, service.Messages().GetMarked())
	return 1
}

// pushes a single message on the stack:
// { __type = "ansicht.Message", id = "...", thread_id = "...", filename = "..." }
const LUA_TYPE_ID_MESSAGE = "ansicht.Message"

func pushMessage(L *lua.State, message *model.Message) int {
	L.CreateTable(0, 4)
	L.PushString(LUA_TYPE_ID_MESSAGE)
	L.SetField(-2, "__type")

	L.PushString(string(message.ID))
	L.SetField(-2, "id")

	L.PushString(string(message.ThreadID))
	L.SetField(-2, "thread_id")

	L.PushString(string(message.Filename))
	L.SetField(-2, "filename")

	return 1
}

func isMessage(L *lua.State, index int) bool {
	if !L.IsTable(index) {
		return false
	}

	name, _ := lFieldString(L, index, "__type")
	return name == LUA_TYPE_ID_MESSAGE
}

// pushes a table of messages on the stack
func pushMessagesTable(L *lua.State, messages []*model.Message) {
	L.CreateTable(len(messages), 0)
	for i, msg := range messages {
		pushMessage(L, msg)
		L.RawSetInt(-2, i+1)
	}
}

// returns message[field] where message is the message table at `index` on the stack
// converts objects to string according to Lua rules
func getMessageField(L *lua.State, index int, field string) (string, bool) {
	if !isMessage(L, index) {
		return "", false
	}

	return lFieldString(L, index, field)
}
