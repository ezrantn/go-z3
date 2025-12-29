package z3

/*
#include <z3.h>
#include <stdlib.h>
*/
import "C"

// Forall creates a universal quantifier: "For all vars, body is true"
// Example: Forall([]*Expr{u}, ctx.GT(uAge, eighteen))
func (ctx *Context) Forall(vars []*Expr, body *Expr) *Expr {
	if len(vars) == 0 {
		return body
	}

	// 1. Convert Go Exprs to C Apps (variables in quantifiers must be constants/apps)
	cVars := make([]C.Z3_app, len(vars))
	for i, v := range vars {
		cVars[i] = C.Z3_to_app(ctx.c, v.ast)
	}

	// 2. Call Z3_mk_forall_const
	// Params: context, weight, num_bound, bound_vars, num_patterns, patterns, body
	res := C.Z3_mk_forall_const(
		ctx.c,
		0,                 // weight (default 0)
		C.uint(len(vars)), // number of bound variables
		&cVars[0],         // the variables themselves
		0, nil,            // patterns (used for advanced E-matching, leave 0/nil for now)
		body.ast, // the logical body
	)

	return ctx.wrap(res)
}

// Exists creates an existential quantifier: "There exists vars such that body is true"
func (ctx *Context) Exists(vars []*Expr, body *Expr) *Expr {
	if len(vars) == 0 {
		return body
	}

	cVars := make([]C.Z3_app, len(vars))
	for i, v := range vars {
		cVars[i] = C.Z3_to_app(ctx.c, v.ast)
	}

	res := C.Z3_mk_exists_const(
		ctx.c,
		0,
		C.uint(len(vars)),
		&cVars[0],
		0, nil,
		body.ast,
	)

	return ctx.wrap(res)
}
