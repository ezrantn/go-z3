package benches

import (
	"testing"

	"github.com/ezrantn/go-z3"
)

// BenchmarkExprCreation measures how fast we can create 10,000 Int constants.
// This tests the CGO overhead and the finalizer registration.
func BenchmarkExprCreation(b *testing.B) {
	cfg := z3.NewConfig()
	ctx := z3.NewContext(cfg)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ctx.Int(i, ctx.IntSort())
	}
}

// BenchmarkLogicModeling measures the time to build a complex nested expression.
func BenchmarkLogicModeling(b *testing.B) {
	cfg := z3.NewConfig()
	ctx := z3.NewContext(cfg)
	intSort := ctx.IntSort()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x := ctx.Const("x", intSort)
		y := ctx.Int(100, intSort)
		_ = ctx.And(ctx.GT(x, y), ctx.LT(x, ctx.Int(200, intSort)))
	}
}

// BenchmarkSolverCheck measures how long it takes for Z3 to solve a problem.
func BenchmarkSolverCheck(b *testing.B) {
	cfg := z3.NewConfig()
	ctx := z3.NewContext(cfg)
	solver := ctx.NewSolver()

	x := ctx.Const("x", ctx.IntSort())
	solver.Assert(ctx.GT(x, ctx.Int(1000, ctx.IntSort())))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !solver.Check() {
			b.Fatal("Expected SAT")
		}
	}
}

// BenchmarkFloatingPoint measures the cost of IEEE 754 operations
func BenchmarkFloatingPoint(b *testing.B) {
	ctx := z3.NewContext(z3.NewConfig())
	f64 := ctx.Float64Sort()
	rm := ctx.RNE()
	val := ctx.FloatVal(1.234, f64)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ctx.FPAAdd(rm, val, val)
	}
}

// BenchmarkQuantifierCreation measures the cost of building First-Order Logic
func BenchmarkQuantifierCreation(b *testing.B) {
	ctx := z3.NewContext(z3.NewConfig())
	x := ctx.Const("x", ctx.IntSort())
	body := ctx.GT(x, ctx.Int(10, ctx.IntSort()))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ctx.Forall([]*z3.Expr{x}, body)
	}
}
