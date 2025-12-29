package z3

/*
#include <z3.h>
#include <stdlib.h>
#include <stdint.h>
*/
import "C"
import (
	"runtime"
	"unsafe"
)

type Expr struct {
	ctx *Context
	ast C.Z3_ast
}

// wrap is our internal helper to handle Z3 reference counting
func (ctx *Context) wrap(ast C.Z3_ast) *Expr {
	e := &Expr{ctx: ctx, ast: ast}
	C.Z3_inc_ref(ctx.c, ast)
	runtime.SetFinalizer(e, func(e *Expr) {
		C.Z3_dec_ref(e.ctx.c, e.ast)
	})
	return e
}

// Const creates a symbolic variable
func (ctx *Context) Const(name string, sort *Sort) *Expr {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	symbol := C.Z3_mk_string_symbol(ctx.c, cname)
	return ctx.wrap(C.Z3_mk_const(ctx.c, symbol, sort.s))
}

// Int creates a numeral integer constant
func (ctx *Context) Int(val int, sort *Sort) *Expr {
	return ctx.wrap(C.Z3_mk_int(ctx.c, C.int(val), sort.s))
}

func (ctx *Context) Eq(l, r *Expr) *Expr {
	return ctx.wrap(C.Z3_mk_eq(ctx.c, l.ast, r.ast))
}

func (ctx *Context) GT(l, r *Expr) *Expr {
	return ctx.wrap(C.Z3_mk_gt(ctx.c, l.ast, r.ast))
}

func (ctx *Context) LT(l, r *Expr) *Expr {
	return ctx.wrap(C.Z3_mk_lt(ctx.c, l.ast, r.ast))
}

func (ctx *Context) And(args ...*Expr) *Expr {
	cArgs := make([]C.Z3_ast, len(args))
	for i, arg := range args {
		cArgs[i] = arg.ast
	}
	// Note: index 0 is safe because Z3 handles empty/single args if we pass count
	var ptr *C.Z3_ast
	if len(cArgs) > 0 {
		ptr = &cArgs[0]
	}
	return ctx.wrap(C.Z3_mk_and(ctx.c, C.uint(len(args)), ptr))
}

// Not negates the given boolean expression: !e
func (ctx *Context) Not(e *Expr) *Expr {
	return ctx.wrap(C.Z3_mk_not(ctx.c, e.ast))
}

// Or performs logical OR: args[0] || args[1] || ...
func (ctx *Context) Or(args ...*Expr) *Expr {
	cArgs := make([]C.Z3_ast, len(args))
	for i, arg := range args {
		cArgs[i] = arg.ast
	}

	var ptr *C.Z3_ast
	if len(cArgs) > 0 {
		ptr = &cArgs[0]
	}

	return ctx.wrap(C.Z3_mk_or(ctx.c, C.uint(len(args)), ptr))
}

// Add performs addition: l + r
func (ctx *Context) Add(args ...*Expr) *Expr {
	cArgs := make([]C.Z3_ast, len(args))
	for i, arg := range args {
		cArgs[i] = arg.ast
	}
	var ptr *C.Z3_ast
	if len(cArgs) > 0 {
		ptr = &cArgs[0]
	}
	return ctx.wrap(C.Z3_mk_add(ctx.c, C.uint(len(args)), ptr))
}

// Apply calls a function with the given arguments
func (ctx *Context) Apply(f *FuncDecl, args ...*Expr) *Expr {
	cArgs := make([]C.Z3_ast, len(args))
	for i, arg := range args {
		cArgs[i] = arg.ast
	}

	var ptr *C.Z3_ast
	if len(cArgs) > 0 {
		ptr = &cArgs[0]
	}

	return ctx.wrap(C.Z3_mk_app(ctx.c, f.d, C.uint(len(args)), ptr))
}

// BVVal creates a bit-vector numeral
// Cast to C.int64_t which matches Z3's internal 64-bit expectation
// If your compiler still complains about 'long', use C.long(val)
func (ctx *Context) BVVal(val int64, bits uint) *Expr {
	sort := ctx.BVSort(bits)
	return ctx.wrap(C.Z3_mk_int64(ctx.c, C.int64_t(val), sort.s))
}

// BVAdd performs bit-vector addition (wraps on overflow)
func (ctx *Context) BVAdd(l, r *Expr) *Expr {
	return ctx.wrap(C.Z3_mk_bvadd(ctx.c, l.ast, r.ast))
}

// BVUgt is Unsigned Greater Than for bit-vectors
func (ctx *Context) BVUgt(l, r *Expr) *Expr {
	return ctx.wrap(C.Z3_mk_bvuge(ctx.c, l.ast, r.ast))
}

// Select reads a value from an array: array[index]
func (ctx *Context) Select(array, index *Expr) *Expr {
	return ctx.wrap(C.Z3_mk_select(ctx.c, array.ast, index.ast))
}

// Store updates an array: returns a NEW array where array[index] = value
func (ctx *Context) Store(array, index, value *Expr) *Expr {
	return ctx.wrap(C.Z3_mk_store(ctx.c, array.ast, index.ast, value.ast))
}

// FloatVal creates a floating point constant from a float64
func (ctx *Context) FloatVal(val float64, sort *Sort) *Expr {
	return ctx.wrap(C.Z3_mk_fpa_numeral_double(ctx.c, C.double(val), sort.s))
}

// FPAAdd performs: l + r using rounding mode rm
func (ctx *Context) FPAAdd(rm, l, r *Expr) *Expr {
	return ctx.wrap(C.Z3_mk_fpa_add(ctx.c, rm.ast, l.ast, r.ast))
}

// FPADiv performs: l / r using rounding mode rm
func (ctx *Context) FPADiv(rm, l, r *Expr) *Expr {
	return ctx.wrap(C.Z3_mk_fpa_div(ctx.c, rm.ast, l.ast, r.ast))
}

// FPAEq performs floating point equality (handles NaN correctly)
func (ctx *Context) FPAEq(l, r *Expr) *Expr {
	return ctx.wrap(C.Z3_mk_fpa_eq(ctx.c, l.ast, r.ast))
}

// FPALt performs: l < r
func (ctx *Context) FPALt(l, r *Expr) *Expr {
	return ctx.wrap(C.Z3_mk_fpa_lt(ctx.c, l.ast, r.ast))
}

// FPANeg returns the additive inverse: -e
func (ctx *Context) FPANeg(e *Expr) *Expr {
	return ctx.wrap(C.Z3_mk_fpa_neg(ctx.c, e.ast))
}

// FPAIsNaN returns true if the expression is NaN
func (ctx *Context) FPAIsNaN(e *Expr) *Expr {
	return ctx.wrap(C.Z3_mk_fpa_is_nan(ctx.c, e.ast))
}

// FPAGt is Floating Point Greater Than: l > r
func (ctx *Context) FPAGt(l, r *Expr) *Expr {
	return ctx.wrap(C.Z3_mk_fpa_gt(ctx.c, l.ast, r.ast))
}

// FPAToIEEEBV converts a Float expression to its bit-level Bitvector representation
func (ctx *Context) FPAToIEEEBV(e *Expr) *Expr {
	return ctx.wrap(C.Z3_mk_fpa_to_ieee_bv(ctx.c, e.ast))
}
