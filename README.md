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
    "github.com/ezrantn/go-z3/z3"
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

## Contributing

Contributions are welcome!

## License

This tool is open-source and available under the [MIT License](https://github.com/ezrantn/go-z3/blob/main/LICENSE)
