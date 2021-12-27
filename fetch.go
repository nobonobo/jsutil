package jsutil

import "syscall/js"

// Fetch wrapper fetch function.
func Fetch(url string, opt map[string]interface{}) (js.Value, error) {
	return Await(global.Call("fetch", url, opt))
}
