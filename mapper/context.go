// Package mapper provides reflection-based object-to-object mapping utilities.
// This file defines the internal context type used to manage per-operation state.
package mapper

import (
	"reflect"
	"sync"

	"github.com/fbarikzehi/gomap/internal/reflectutil"
)

// context represents the internal state of a single mapping operation.
//
// It is created for each Copy() call and reused from a sync.Pool to minimize
// allocations. The context tracks recursion depth, visited references
// (for circular reference detection), and mapping errors.
//
// The context is concurrency-safe for use across recursive or concurrent
// mapping paths within a single operation, but it is not intended for
// sharing between independent Copy() calls.
type context struct {
	// visited tracks visited pointers to detect circular references
	visited map[uintptr]reflect.Value

	// depth represents the current recursion depth
	depth int

	// config holds the active mapping configuration
	config *Config

	// errors accumulates errors encountered during mapping
	errors []error

	// mu protects concurrent access to visited and errors
	mu sync.RWMutex
}

// newContext initializes and returns a new mapping context.
// The returned context is not shared and should be obtained
// from a pool or newly allocated per Copy() call.
func newContext(cfg *Config) *context {
	return &context{
		visited: make(map[uintptr]reflect.Value),
		config:  cfg,
		errors:  make([]error, 0),
	}
}

// reset clears the context to its initial state so it can be safely reused.
// It preserves allocated memory to minimize GC pressure.
func (ctx *context) reset(cfg *Config) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	for k := range ctx.visited {
		delete(ctx.visited, k)
	}

	ctx.errors = ctx.errors[:0]
	ctx.depth = 0
	ctx.config = cfg
}

// checkCircular detects circular references by tracking visited pointers.
// It returns ErrCircularReference if the given value has been seen before.
//
// Non-pointer values and invalid reflect.Values are ignored.
func (ctx *context) checkCircular(v reflect.Value) error {
	if !v.IsValid() || !reflectutil.IsPointerLike(v.Kind()) {
		return nil
	}

	ptr := v.Pointer()
	if ptr == 0 {
		return nil
	}

	ctx.mu.RLock()
	_, exists := ctx.visited[ptr]
	ctx.mu.RUnlock()

	if exists {
		return ErrCircularReference
	}

	ctx.mu.Lock()
	ctx.visited[ptr] = v
	ctx.mu.Unlock()
	return nil
}

// addError appends an error to the context's error list.
// Nil errors are ignored.
func (ctx *context) addError(err error) {
	if err == nil {
		return
	}
	ctx.mu.Lock()
	ctx.errors = append(ctx.errors, err)
	ctx.mu.Unlock()
}

// hasErrors reports whether any errors were recorded in the context.
func (ctx *context) hasErrors() bool {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	return len(ctx.errors) > 0
}

// firstError returns the first recorded error, or nil if none exist.
func (ctx *context) firstError() error {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()

	if len(ctx.errors) == 0 {
		return nil
	}
	return ctx.errors[0]
}

// allErrors returns a shallow copy of all recorded errors.
//
// The returned slice is safe to read and modify, as it is detached
// from the internal error list.
func (ctx *context) allErrors() []error {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()

	errors := make([]error, len(ctx.errors))
	copy(errors, ctx.errors)
	return errors
}
