package z3

import (
	"runtime"
	"testing"
)

func TestBasicLogic(t *testing.T) {
	cfg := NewConfig()
	defer cfg.Close()
	ctx := NewContext(cfg)

	solver := ctx.NewSolver()

	intSort := ctx.IntSort()
	x := ctx.Const("x", intSort)
	ten := ctx.Int(10, intSort)
	twenty := ctx.Int(20, intSort)

	// Rule: x + 10 > 20
	sum := ctx.Add(x, ten)
	constraint := ctx.GT(sum, twenty)

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

func TestBitvectorOverflow(t *testing.T) {
	ctx := NewContext(NewConfig())
	solver := ctx.NewSolver()

	// Use an 8-bit integer (0-255)
	bv8 := ctx.BVSort(8)
	x := ctx.Const("x", bv8)
	max := ctx.BVVal(255, 8)
	one := ctx.BVVal(1, 8)
	zero := ctx.BVVal(0, 8)

	// Rule: x == 255 AND (x + 1) == 0
	solver.Assert(ctx.Eq(x, max))
	solver.Assert(ctx.Eq(ctx.BVAdd(x, one), zero))

	if !solver.Check() {
		t.Fatal("Bitvector overflow logic failed: 255 + 1 should be 0 in 8-bit")
	}

	t.Log("Success: Verified 8-bit overflow edge case")
}

func TestArrayLogic(t *testing.T) {
	ctx := NewContext(NewConfig())
	solver := ctx.NewSolver()

	// Array mapping Int -> Int
	arrSort := ctx.ArraySort(ctx.IntSort(), ctx.IntSort())
	a := ctx.Const("a", arrSort)

	i1 := ctx.Int(1, ctx.IntSort())
	i2 := ctx.Int(2, ctx.IntSort())
	val100 := ctx.Int(100, ctx.IntSort())

	// b = a with [index 1] set to 100
	b := ctx.Store(a, i1, val100)

	// Assert: b[1] == 100 AND b[2] == a[2]
	solver.Assert(ctx.Eq(ctx.Select(b, i1), val100))
	solver.Assert(ctx.Eq(ctx.Select(b, i2), ctx.Select(a, i2)))

	if !solver.Check() {
		t.Fatal("Array Select/Store logic failed")
	}
}

func TestQuantifierContradiction(t *testing.T) {
	ctx := NewContext(NewConfig())
	solver := ctx.NewSolver()

	intSort := ctx.IntSort()
	u := ctx.Const("u", intSort)

	// Rule: "For all u, u > 10"
	forallRule := ctx.Forall([]*Expr{u}, ctx.GT(u, ctx.Int(10, intSort)))
	solver.Assert(forallRule)

	// Contradiction: "There exists some constant 'k' where k == 5"
	k := ctx.Const("k", intSort)
	solver.Assert(ctx.Eq(k, ctx.Int(5, intSort)))

	if solver.Check() {
		t.Fatal("Logic Error: Solver said SAT, but 'k=5' contradicts 'all u > 10'")
	}

	t.Log("Success: Quantifier correctly blocked the contradiction")
}

func TestMemoryStress(t *testing.T) {
	ctx := NewContext(NewConfig())

	// Create 100,000 expressions and throw them away
	for i := 0; i < 100000; i++ {
		_ = ctx.Int(i, ctx.IntSort())
		if i%10000 == 0 {
			runtime.GC() // Force GC to trigger finalizers
		}
	}

	t.Log("Success: Processed 100k expressions without crashing")
}
