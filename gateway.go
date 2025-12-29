package z3

/*
#include <z3.h>
#include <stdio.h>

// This is a C function that can be called by Z3
// It can then call a Go function if we exported one
void errorHandler(Z3_context c, Z3_error_code e) {
    fprintf(stderr, "Z3 Error: %d\n", e);
}
*/
import "C"
