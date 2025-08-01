package runtime

import "github.com/Shopify/go-lua"

type ControllerAdapter interface {
	Quit()
	Refresh()
	Status(message string)
	Notify(message string, level string, timeout float64)
	Input(prompt, placeholder string)
	SpawnResult(result SpawnResult)
	SetTheme(theme any)

	QueryNew(query string)
	QuerySelectNext()
	QuerySelectPrev()

	MarksToggle()
	MarksInvert()
	MarksClear()
}

type NullAdapter struct{}

func (a *NullAdapter) Quit()                          {}
func (a *NullAdapter) Refresh()                       {}
func (a *NullAdapter) Status(string)                  {}
func (a *NullAdapter) Notify(string, string, float64) {}
func (a *NullAdapter) Input(string, string)           {}
func (a *NullAdapter) SpawnResult(SpawnResult)        {}
func (a *NullAdapter) SetTheme(any)                   {}

func (a *NullAdapter) QueryNew(string)  {}
func (a *NullAdapter) QuerySelectNext() {}
func (a *NullAdapter) QuerySelectPrev() {}

func (a *NullAdapter) MarksToggle() {}
func (a *NullAdapter) MarksInvert() {}
func (a *NullAdapter) MarksClear()  {}

func (r *Runtime) luaQuit(L *lua.State) int {
	r.Controller.Quit()
	return 0
}

func (r *Runtime) luaRefresh(L *lua.State) int {
	r.Controller.Refresh()
	return 0
}

func (r *Runtime) luaStatusSet(L *lua.State) int {
	s, _ := lua.ToStringMeta(r.luaState, 1)
	r.Controller.Status(s)
	return 0
}

func (r *Runtime) luaQueryNew(L *lua.State) int {
	if s, ok := r.luaState.ToString(1); ok {
		r.Controller.QueryNew(s)
	}
	return 0
}

func (r *Runtime) luaQuerySelectNext(L *lua.State) int {
	r.Controller.QuerySelectNext()
	return 0
}

func (r *Runtime) luaQuerySelectPrev(L *lua.State) int {
	r.Controller.QuerySelectPrev()
	return 0
}

func (r *Runtime) luaMarksToggle(L *lua.State) int {
	r.Controller.MarksToggle()
	return 0
}

func (r *Runtime) luaMarksInvert(L *lua.State) int {
	r.Controller.MarksInvert()
	return 0
}

func (r *Runtime) luaMarksClear(L *lua.State) int {
	r.Controller.MarksClear()
	return 0
}

func (r *Runtime) luaNotify(L *lua.State) int {
	if !L.IsTable(1) {
		lua.Errorf(L, "ansicht.notify expects a table argument")
		panic("unreachable")
	}

	message, ok := lFieldString(L, 1, "message")
	if !ok {
		return 0
	}

	level := lFieldStringOrDefault(L, 1, "level", "info")
	timeout := lFieldNumberOrDefault(L, 1, "timeout", 0)
	r.Controller.Notify(message, level, timeout)

	return 0
}
