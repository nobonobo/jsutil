package jsutil

import "syscall/js"

type wrappedError js.Value

func (w wrappedError) Error() string {
	return js.Value(w).Call("toString").String()
}

func (w wrappedError) JSValue() js.Value {
	return js.Value(w)
}

// Await equivalent for js await statement.
func Await(promise js.Value) (res js.Value, err error) {
	ch := make(chan bool)
	then := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		res = args[0]
		close(ch)
		return js.Undefined()
	})
	defer then.Release()
	catch := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		err = wrappedError(args[0])
		close(ch)
		return js.Undefined()
	})
	defer catch.Release()
	promise.Call("then", then).Call("catch", catch)
	<-ch
	return
}
