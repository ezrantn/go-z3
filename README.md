# go-z3

A high-performance, idiomatic Go wrapper for the Z3 Theorem Prover.

Unlike other wrappers, go-z3 provides a clean, Go-first API while utilizing CGO to communicate directly with Z3's native C API for maximum performance. It features automatic memory management via Go finalizers, so you don't have to worry about manual reference counting.

> [!WARNING]
> This library uses CGO internally. You must have the Z3 development headers installed on your system to build and run this library.

## Features

- Logical Operations: Core support for Propositional Logic (And, Or, Not, Xor, Implies).
- Bit-Vectors: Machine-precision arithmetic (8, 32, 64-bit) with support for bitwise operations and overflow modeling.
- Floating Point: Full IEEE 754 support (Single and Double precision) with configurable Rounding Modes and handling of NaN and ±∞.
- Functional Arrays: Model infinite mappings and memory states using functional Select and Store operations.
- Function Declarations: Define uninterpreted functions to model object properties, struct fields, and custom relations.
- Quantifiers: Support for First-Order Logic using Universal (∀) and Existential (∃) quantifiers for property verification.

## Installation

```bash
go get github.com/ezrantn/go-z3
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/ezrantn/go-z3"
)

func main() {
    // 1. Initialize Z3
    cfg := z3.NewConfig()
    defer cfg.Close()
    ctx := z3.NewContext(cfg)

    // 2. Create a Solver
    solver := ctx.NewSolver()

    // 3. Define Symbolic Variables
    intSort := ctx.IntSort()
    x := ctx.Const("x", intSort)
    ten := ctx.Int(10, intSort)

    // 4. Assert Constraints (x > 10)
    solver.Assert(ctx.GT(x, ten))

    // 5. Check and Evaluate
    if solver.Check() {
        model := solver.GetModel()
        fmt.Printf("Satisfiable! x = %s\n", model.Eval(x))
    } else {
        fmt.Println("Unsatisfiable.")
    }
}
```

## Benchmark

Environment:

- OS/Arch: linux/amd64
- CPU: 11th Gen Intel(R) Core(TM) i5-11400H @ 2.70GHz

| **Benchmark**                   | **Iterations** | **Time (ns/op)** | **Memory (B/op)** | **Allocs/op** |
| ------------------------------- | -------------- | ---------------- | ----------------- | ------------- |
| **BenchmarkExprCreation**       | 334,828        | 3,615            | 16                | 1             |
| **BenchmarkLogicModeling**      | 170,410        | 7,285            | 112               | 7             |
| **BenchmarkSolverCheck**        | 404            | 2,857,793        | 0                 | 0             |
| **BenchmarkFloatingPoint**      | 1,000,000      | 1,277            | 16                | 1             |
| **BenchmarkQuantifierCreation** | 627,232        | 1,754            | 24                | 2             |

The go-z3 library shows high efficiency across core operations, with expression creation and logic modeling completing in microseconds. Benchmarks indicate minimal CGO overhead, enabling over a million floating-point operations per second, while complex tasks like quantifiers and nested logic remain fast, ensuring the Go–Z3 bridge is not a bottleneck during AST construction.

During solving, go-z3 incurs 0 B/op and 0 allocs/op on the Go side, confirming that computation is fully offloaded to the native Z3 engine and does not stress the Go GC. Overall, go-z3 is faster and more memory-efficient than string-based approaches, making it well-suited for high-performance formal verification.

## Contributing

Contributions are welcome!

## License

This tool is open-source and available under the [MIT License](https://github.com/ezrantn/go-z3/blob/main/LICENSE)
