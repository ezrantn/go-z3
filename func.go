package z3

/*
#include <z3.h>
#include <stdlib.h>
*/
import "C"
import (
	"runtime"
	"unsafe"
)

type FuncDecl struct {
	c *Context
	d C.Z3_func_decl
}

// CreateFuncDecl defines a function: Name(Domain) -> Range
// For a struct field: FieldName(StructSort) -> FieldTypeSort
func (ctx *Context) CreateFuncDecl(name string, domain []*Sort, rangeSort *Sort) *FuncDecl {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	symbol := C.Z3_mk_string_symbol(ctx.c, cName)

	// Convert Go Sort slice to C Sort array
	cDomain := make([]C.Z3_sort, len(domain))
	for i, s := range domain {
		cDomain[i] = s.s
	}

	var domPtr *C.Z3_sort
	if len(cDomain) > 0 {
		domPtr = &cDomain[0]
	}

	d := C.Z3_mk_func_decl(ctx.c, symbol, C.uint(len(domain)), domPtr, rangeSort.s)

	fd := &FuncDecl{c: ctx, d: d}
	C.Z3_inc_ref(ctx.c, C.Z3_func_decl_to_ast(ctx.c, d))

	runtime.SetFinalizer(fd, func(fd *FuncDecl) {
		C.Z3_dec_ref(fd.c.c, C.Z3_func_decl_to_ast(fd.c.c, fd.d))
	})

	return fd
}
