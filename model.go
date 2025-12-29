package z3

/*
#include <z3.h>
*/
import "C"

type Model struct {
	ctx *Context
	m   C.Z3_model
}

func (s *Solver) GetModel() *Model {
	m := C.Z3_solver_get_model(s.ctx.c, s.s)
	C.Z3_model_inc_ref(s.ctx.c, m)
	return &Model{ctx: s.ctx, m: m}
}

func (m *Model) Eval(e *Expr) string {
	var res C.Z3_ast
	// C.Z3_model_eval returns a Z3_bool (which Go sees as a bool).
	// We cast it to a Go bool to be safe and compare it to true.
	if bool(C.Z3_model_eval(m.ctx.c, m.m, e.ast, C.bool(true), &res)) != true {
		return "unknown"
	}

	// Convert the resulting AST to a string
	return C.GoString(C.Z3_ast_to_string(m.ctx.c, res))
}
