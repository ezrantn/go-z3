package z3

import (
	"fmt"
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

func TestFloatAssociativity(t *testing.T) {
	ctx := NewContext(NewConfig())
	solver := ctx.NewSolver()
	f32 := ctx.Float32Sort()
	rm := ctx.RNE() // Round to nearest even

	// Constants
	a := ctx.Const("a", f32)
	b := ctx.Const("b", f32)
	c := ctx.Const("c", f32)

	// lhs = (a + b) + c
	lhs := ctx.FPAAdd(rm, ctx.FPAAdd(rm, a, b), c)
	// rhs = a + (b + c)
	rhs := ctx.FPAAdd(rm, a, ctx.FPAAdd(rm, b, c))

	// Assert that they are NOT equal
	solver.Assert(ctx.Not(ctx.FPAEq(lhs, rhs)))

	if solver.Check() {
		m := solver.GetModel()
		t.Logf("Found case where (a+b)+c != a+(b+c):")
		t.Logf("a: %s, b: %s, c: %s", m.Eval(a), m.Eval(b), m.Eval(c))
	} else {
		t.Fatal("Z3 says float addition is always associative (this is wrong!)")
	}
}

func TestDeMorgan(t *testing.T) {
	ctx := NewContext(NewConfig())
	solver := ctx.NewSolver()

	a := ctx.Const("a", ctx.BoolSort())
	b := ctx.Const("b", ctx.BoolSort())

	// lhs = !(a && b)
	lhs := ctx.Not(ctx.And(a, b))
	// rhs = !a || !b
	rhs := ctx.Or(ctx.Not(a), ctx.Not(b))

	// Theorem: lhs == rhs. To prove it, we assert !(lhs == rhs) and expect UNSAT.
	solver.Assert(ctx.Not(ctx.Eq(lhs, rhs)))

	if solver.Check() {
		t.Fatal("De Morgan's law failed! Solver found a counter-example where !(a&&b) != !a||!b")
	}
	t.Log("Success: Boolean logic core is sound.")
}

func TestFloatToBitvector(t *testing.T) {
	ctx := NewContext(NewConfig())
	solver := ctx.NewSolver()

	f32 := ctx.Float32Sort()
	// Create a float: 1.0
	val := ctx.FloatVal(1.0, f32)

	// Cast float bits to BV32
	// Note: You'll need ctx.FPAToIEEEBV()
	bv := ctx.FPAToIEEEBV(val)

	// In IEEE 754, 1.0 is represented as 0x3f800000
	expected := ctx.BVVal(0x3f800000, 32)
	solver.Assert(ctx.Eq(bv, expected))

	if !solver.Check() {
		t.Fatal("Float to Bitvector cast failed. 1.0 bits should be 0x3f800000")
	}
}

func TestArrayOfStructs(t *testing.T) {
	ctx := NewContext(NewConfig())
	solver := ctx.NewSolver()

	userSort := ctx.CreateSort("User")
	ageField := ctx.CreateFuncDecl("Age", []*Sort{userSort}, ctx.IntSort())
	userArray := ctx.ArraySort(ctx.IntSort(), userSort)

	users := ctx.Const("users", userArray)
	idx := ctx.Int(0, ctx.IntSort())

	// Access: users[0].Age
	userAtZero := ctx.Select(users, idx)
	ageAtZero := ctx.Apply(ageField, userAtZero)

	// Assert users[0].Age == 25
	solver.Assert(ctx.Eq(ageAtZero, ctx.Int(25, ctx.IntSort())))

	if !solver.Check() {
		t.Fatal("Failed to solve nested Array and FuncDecl logic")
	}

	m := solver.GetModel()
	t.Logf("Verified: users[0].Age = %s", m.Eval(ageAtZero))
}

func TestConcurrentContexts(t *testing.T) {
	for i := 0; i < 10; i++ {
		go func(id int) {
			cfg := NewConfig()
			ctx := NewContext(cfg)
			x := ctx.Const(fmt.Sprintf("x%d", id), ctx.IntSort())
			solver := ctx.NewSolver()
			solver.Assert(ctx.GT(x, ctx.Int(10, ctx.IntSort())))
			if !solver.Check() {
				t.Errorf("Concurrent solver %d failed", id)
			}
		}(i)
	}
}
