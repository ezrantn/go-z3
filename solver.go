package z3

/*
#include <z3.h>
*/
import "C"
import "runtime"

type Solver struct {
	ctx *Context
	s   C.Z3_solver
}

func (ctx *Context) NewSolver() *Solver {
	s := &Solver{
		ctx: ctx,
		s:   C.Z3_mk_solver(ctx.c),
	}

	C.Z3_solver_inc_ref(ctx.c, s.s)

	runtime.SetFinalizer(s, func(s *Solver) {
		C.Z3_solver_dec_ref(s.ctx.c, s.s)
	})

	return s
}

func (s *Solver) Assert(e *Expr) {
	C.Z3_solver_assert(s.ctx.c, s.s, e.ast)
}

func (s *Solver) Check() bool {
	// Z3_L_TRUE is an enum value 1.
	// We cast to int to ensure the comparison works across all Go versions.
	return int(C.Z3_solver_check(s.ctx.c, s.s)) == 1
}
