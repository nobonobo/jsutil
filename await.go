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
	promise.Call("then",
		Callback1(func(arg js.Value) interface{} {
			res = arg
			close(ch)
			return nil
		}),
	).Call("catch",
		Callback1(func(arg js.Value) interface{} {
			err = wrappedError(arg)
			close(ch)
			return nil
		}),
	)
	<-ch
	return
}
