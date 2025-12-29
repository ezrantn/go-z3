package z3

/*
#include <z3.h>
*/
import "C"
import "runtime"

type Optimize struct {
	ctx *Context
	o   C.Z3_optimize
}

func (ctx *Context) NewOptimize() *Optimize {
	o := C.Z3_mk_optimize(ctx.c)
	opt := &Optimize{ctx: ctx, o: o}
	C.Z3_optimize_inc_ref(ctx.c, o)

	runtime.SetFinalizer(opt, func(opt *Optimize) {
		C.Z3_optimize_dec_ref(opt.ctx.c, opt.o)
	})
	return opt
}

func (o *Optimize) Assert(e *Expr) {
	C.Z3_optimize_assert(o.ctx.c, o.o, e.ast)
}

// Maximize adds an objective to maximize the value of an expression
func (o *Optimize) Maximize(e *Expr) {
	C.Z3_optimize_maximize(o.ctx.c, o.o, e.ast)
}

// Minimize adds an objective to minimize the value of an expression
func (o *Optimize) Minimize(e *Expr) {
	C.Z3_optimize_minimize(o.ctx.c, o.o, e.ast)
}

func (o *Optimize) Check() bool {
	// 0 args for simple check
	return int(C.Z3_optimize_check(o.ctx.c, o.o, 0, nil)) == 1
}

func (o *Optimize) GetModel() *Model {
	m := C.Z3_optimize_get_model(o.ctx.c, o.o)
	C.Z3_model_inc_ref(o.ctx.c, m)
	return &Model{ctx: o.ctx, m: m}
}
