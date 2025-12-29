package z3

import (
	"testing"
)

func TestBasicLogic(t *testing.T) {
	cfg := NewConfig()
	defer cfg.Close()
	ctx := NewContext(cfg)
	// No need to manually close ctx, Finalizer handles Z3_del_context

	// 2. Create the Solver
	solver := ctx.NewSolver()

	// 3. Define the World: Integer x
	intSort := ctx.IntSort()
	x := ctx.Const("x", intSort)
	ten := ctx.Int(10, intSort)
	twenty := ctx.Int(20, intSort)

	// 4. Build Expression: x + 10 > 20
	// Note: You'll need Add() and GT() in your expr.go
	sum := ctx.Add(x, ten)
	constraint := ctx.GT(sum, twenty)

	// 5. Assert and Check
	solver.Assert(constraint)

	if !solver.Check() {
		t.Fatal("Expected SAT (it is possible for x + 10 > 20), but got UNSAT")
	}

	if solver.Check() {
		model := solver.GetModel()
		val := model.Eval(x)
		t.Logf("Success! Z3 found a solution: x = %s", val)
	}

	t.Log("Success: Z3 found a solution where x + 10 > 20")
}
