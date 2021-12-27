package jsutil

import (
	"syscall/js"
)

// JS2Bytes convert from TypedArray for JS to byte slice for Go.
func JS2Bytes(dv js.Value) []byte {
	b := make([]byte, dv.Get("byteLength").Int())
	js.CopyBytesToGo(b, global.Get("Uint8Array").New(dv.Get("buffer")))
	return b
}

// Bytes2JS convert from byte slice for Go to Uint8Array for JS.
func Bytes2JS(b []byte) js.Value {
	res := global.Get("Uint8Array").New(len(b))
	js.CopyBytesToJS(res, b)
	return res
}

// Callback0 make auto-release callback without params.
func Callback0(fn func() interface{}) js.Func {
	var cb js.Func
	cb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		defer cb.Release()
		return fn()
	})
	return cb
}

// Callback1 make auto-release callback with 1 param.
func Callback1(fn func(res js.Value) interface{}) js.Func {
	var cb js.Func
	cb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		defer cb.Release()
		return fn(args[0])
	})
	return cb
}

// CallbackN make auto-release callback with multiple params.
func CallbackN(fn func(res []js.Value) interface{}) js.Func {
	var cb js.Func
	cb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		defer cb.Release()
		return fn(args)
	})
	return cb
}

// RequestAnimationFrame function call for 30 or 60 fps.
// return value: tick chan
func RequestAnimationFrame(ch <-chan bool, callback func(dt float64)) {
	var cb js.Func
	lastID := -1
	lastTick := 0
	cb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		tick := args[0].Int()
		dt := float64(tick-lastTick) / 1000.0
		lastTick = tick
		go callback(dt)
		b, ok := <-ch
		if !b || !ok {
			global.Call("cancelAnimationFrame", lastID)
			cb.Release()
		}
		lastID = global.Call("requestAnimationFrame", cb).Int()
		return nil
	})
	go cb.Invoke(js.ValueOf(0.0))
}

// IsArray checking value is array type.
func IsArray(item js.Value) bool {
	return array.Call("isArray", item).Bool()
}

// JS2Go JS values convert to Go values.
func JS2Go(obj js.Value) interface{} {
	switch obj.Type() {
	default:
		return obj
	case js.TypeBoolean:
		return obj.Bool()
	case js.TypeNumber:
		return obj.Float()
	case js.TypeString:
		return obj.String()
	case js.TypeObject:
		if IsArray(obj) {
			res := []interface{}{}
			for i := 0; i < obj.Length(); i++ {
				res = append(res, JS2Go(obj.Index(i)))
			}
			return res
		}
		res := map[string]interface{}{}
		entries := object.Call("entries", obj)
		for i := 0; i < entries.Length(); i++ {
			entry := entries.Index(i)
			key, value := entry.Index(0).String(), entry.Index(1)
			res[key] = JS2Go(value)
		}
		return res
	}
}
