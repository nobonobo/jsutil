package jsutil

import "syscall/js"

var (
	global = js.Global()
	array  = global.Get("Array")
	object = global.Get("Object")
)
