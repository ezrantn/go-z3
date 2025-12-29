package z3

/*
#include <z3.h>
#include <stdlib.h>
*/
import "C"
import "unsafe"

type Sort struct {
	c *Context
	s C.Z3_sort
}

// BoolSort returns the built-in Boolean type
func (ctx *Context) BoolSort() *Sort {
	return &Sort{c: ctx, s: C.Z3_mk_bool_sort(ctx.c)}
}

// IntSort returns the built-in Integer type
func (ctx *Context) IntSort() *Sort {
	return &Sort{c: ctx, s: C.Z3_mk_int_sort(ctx.c)}
}

// CreateSort creates a custom "Uninterpreted" sort (like 'User' or 'Profile')
func (ctx *Context) CreateSort(name string) *Sort {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName)) // Requires "unsafe" and "stdlib.h" in preamble

	symbol := C.Z3_mk_string_symbol(ctx.c, cName)
	return &Sort{c: ctx, s: C.Z3_mk_uninterpreted_sort(ctx.c, symbol)}
}
