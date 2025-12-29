package z3

/*
#cgo LDFLAGS: -lz3
#include <z3.h>

// We define the error handler here ONCE to avoid "multiple definition" errors
extern void errorHandler(Z3_context c, Z3_error_code e);
*/
import "C"

// Config and other globals can go here
type Config struct {
	c C.Z3_config
}

func NewConfig() *Config {
	return &Config{c: C.Z3_mk_config()}
}

func (cfg *Config) Close() {
	C.Z3_del_config(cfg.c)
}
