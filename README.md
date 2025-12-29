# go-z3

A high-performance, idiomatic Go wrapper for the Z3 Theorem Prover.

Unlike other wrappers, go-z3 provides a clean, Go-first API while utilizing CGO to communicate directly with Z3's native C API for maximum performance. It features automatic memory management via Go finalizers, so you don't have to worry about manual reference counting.

> [!WARNING]
> This library uses CGO internally. You must have the Z3 development headers installed on your system to build and run this library.

## Features

- Idiomatic API: No more Lisp-style SMT-LIB strings. Use Go methods to build logic.

- Automatic Memory Management: Uses runtime.SetFinalizer to handle Z3_inc_ref and Z3_dec_ref automatically.

- Context Safety: Prevents mixing expressions between different Z3 contexts.

- Model Extraction: Easily evaluate symbolic variables into Go types.

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

## Roadmap

- [x] Basic Sorts (Int, Bool, Uninterpreted)
- [x] Core Logic (And, Or, Not, Eq, Gt, Lt)
- [x] Solver & Model Extraction
- [ ] FuncDecl: Support for uninterpreted functions (Struct fields)
- [ ] Quantifiers: Support for forall and exists
- [ ] Arrays & Bitvectors: Full SMT-LIB type support

## Contributing

Contributions are welcome!
