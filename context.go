package z3

/*
#include <z3.h>
*/
import "C"
import "runtime"

type Context struct {
	c C.Z3_context
}

func NewContext(cfg *Config) *Context {
	ctx := &Context{
		c: C.Z3_mk_context(cfg.c),
	}

	// Clean up memory via Go's GC
	runtime.SetFinalizer(ctx, func(c *Context) {
		C.Z3_del_context(c.c)
	})
	return ctx
}
