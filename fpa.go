package z3

/*
#include <z3.h>
*/
import "C"

// Rounding modes as defined in Z3
func (ctx *Context) RNA() *Expr { return ctx.wrap(C.Z3_mk_fpa_round_nearest_ties_to_away(ctx.c)) }
func (ctx *Context) RNE() *Expr { return ctx.wrap(C.Z3_mk_fpa_round_nearest_ties_to_even(ctx.c)) }
func (ctx *Context) RTP() *Expr { return ctx.wrap(C.Z3_mk_fpa_round_toward_positive(ctx.c)) }
func (ctx *Context) RTN() *Expr { return ctx.wrap(C.Z3_mk_fpa_round_toward_negative(ctx.c)) }
func (ctx *Context) RTZ() *Expr { return ctx.wrap(C.Z3_mk_fpa_round_toward_zero(ctx.c)) }
