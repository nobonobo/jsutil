package jsutil

import "syscall/js"

// Releaser ...
type Releaser interface {
	Release()
}

type wrapper struct {
	fn func()
}

// Release ...
func (w *wrapper) Release() {
	w.fn()
}

// ReleaserFunc ...
func ReleaserFunc(fn func()) Releaser {
	return &wrapper{fn: fn}
}

// Bind event bind and return releaser.
func Bind(node js.Value, name string, callback func(res js.Value)) Releaser {
	fn := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		callback(args[0])
		return nil
	})
	node.Call("addEventListener", name, fn)
	return ReleaserFunc(func() {
		node.Call("removeEventListener", name, fn)
		fn.Release()
	})
}
