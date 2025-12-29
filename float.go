package z3

// In Z3, floating point follows the IEEE 754 standard.
// https://en.wikipedia.org/wiki/IEEE_754

/*
#include <z3.h>
#include <stdlib.h>
*/
import "C"

// Float32Sort returns the IEEE 754 single precision sort
func (ctx *Context) Float32Sort() *Sort {
	return &Sort{c: ctx, s: C.Z3_mk_fpa_sort_single(ctx.c)}
}

// Float64Sort returns the IEEE 754 double precision sort
func (ctx *Context) Float64Sort() *Sort {
	return &Sort{c: ctx, s: C.Z3_mk_fpa_sort_double(ctx.c)}
}

// RoundingModeSort returns the sort for rounding modes
func (ctx *Context) RoundingModeSort() *Sort {
	return &Sort{c: ctx, s: C.Z3_mk_fpa_rounding_mode_sort(ctx.c)}
}
