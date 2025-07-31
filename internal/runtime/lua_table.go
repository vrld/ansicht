package runtime

import "github.com/Shopify/go-lua"

func lFieldFunctionOrNil(L *lua.State, index int, key string) {
	L.Field(index, key)
	if !(L.IsFunction(-1) || L.IsNil(-1)) {
		lua.Errorf(L, "%s must be a string or nil", key)
		panic("unreachable")
	}
}

func lSetFieldNil(L *lua.State, index int, key string) {
	L.PushString(key)
	L.PushNil()
	if index == lua.RegistryIndex {
		L.SetTable(lua.RegistryIndex)
	} else {
		L.SetTable(index - 2)
	}
}

func lSetFieldString(L *lua.State, index int, key string, value string) {
	L.PushString(key)
	L.PushString(value)
	if index != lua.RegistryIndex {
		L.SetTable(lua.RegistryIndex)
	} else {
		L.SetTable(index - 2)
	}
}

func lSetFieldInteger(L *lua.State, index int, key string, value int) {
	L.PushString(key)
	L.PushInteger(value)
	if index != lua.RegistryIndex {
		L.SetTable(lua.RegistryIndex)
	} else {
		L.SetTable(index - 2)
	}
}

func lSetFieldBool(L *lua.State, index int, key string, value bool) {
	L.PushString(key)
	L.PushBoolean(value)
	if index != lua.RegistryIndex {
		L.SetTable(lua.RegistryIndex)
	} else {
		L.SetTable(index - 2)
	}
}

func lPushStringTable(L *lua.State, slice []string) {
	L.CreateTable(len(slice), 0)
	for i, arg := range slice {
		L.PushString(arg)
		L.RawSetInt(-2, i+1)
	}
}

func lFieldString(L *lua.State, index int, key string) (string, bool) {
	L.Field(index, key)
	defer L.Pop(1)
	if value, ok := L.ToString(-1); ok {
		return value, ok
	}
	return "", false
}

func lFieldStringOrDefault(L *lua.State, index int, key string, defaultValue string) string {
	if value, ok := lFieldString(L, index, key); ok {
		return value
	}
	return defaultValue
}

func lFieldNumber(L *lua.State, index int, key string) (float64, bool) {
	L.Field(index, key)
	defer L.Pop(1)
	if value, ok := L.ToNumber(-1); ok {
		return value, ok
	}
	return 0, false
}
