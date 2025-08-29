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
// { __type = "ansicht.Message", id = "...", thread_id = "...", filename = "...", ... }
const LUA_TYPE_ID_MESSAGE = "ansicht.Message"

func pushMessage(L *lua.State, message *model.Message) int {
	L.CreateTable(0, 10)
	lSetFieldString(L, -1, "__type", LUA_TYPE_ID_MESSAGE)
	lSetFieldString(L, -1, "id", string(message.ID))
	lSetFieldString(L, -1, "thread_id", message.ThreadID)
	lSetFieldString(L, -1, "date", message.Date.String())
	lSetFieldString(L, -1, "filename", string(message.Filename))
	lSetFieldString(L, -1, "from", message.From)
	lSetFieldString(L, -1, "to", message.To)
	lSetFieldString(L, -1, "subject", message.Subject)

	L.CreateTable(len(message.Tags), 0)
	for i, tag := range message.Tags {
		L.PushString(tag)
		L.RawSetInt(-2, i + 1)
	}
	L.SetField(-2, "tags")

	L.CreateTable(0, 6)
	lSetFieldBool(L, -1, "draft", message.Flags.Draft)
	lSetFieldBool(L, -1, "flagged", message.Flags.Flagged)
	lSetFieldBool(L, -1, "passed", message.Flags.Passed)
	lSetFieldBool(L, -1, "replied", message.Flags.Replied)
	lSetFieldBool(L, -1, "seen", message.Flags.Seen)
	lSetFieldBool(L, -1, "trashed", message.Flags.Trashed)
	L.SetField(-2, "flags")

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
