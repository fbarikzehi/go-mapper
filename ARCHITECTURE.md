# Architecture Documentation

## Overview

gomap is designed with performance, safety, and extensibility in mind. This document describes the internal architecture and design decisions.

## Component Overview

```
┌─────────────────────────────────────────────────────────┐
│                    Public API Layer                      │
│  ┌─────────────┐  ┌──────────────┐  ┌────────────────┐ │
│  │   Mapper    │  │   Options    │  │  Convenience   │ │
│  │   Struct    │  │  Functions   │  │   Functions    │ │
│  └─────────────┘  └──────────────┘  └────────────────┘ │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│                  Configuration Layer                     │
│  ┌──────────────────────────────────────────────────┐  │
│  │  Config struct with validation and cloning       │  │
│  └──────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│                   Execution Layer                        │
│  ┌─────────────┐  ┌──────────────┐  ┌────────────────┐ │
│  │  Context    │  │   Mapper     │  │    Object      │ │
│  │  Manager    │  │   Logic      │  │    Pool        │ │
│  └─────────────┘  └──────────────┘  └────────────────┘ │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│                Type-Specific Handlers                    │
│  ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐│
│  │Struct│ │Slice │ │ Map  │ │Pointer│ │Basic│ │Inter-││
│  │      │ │Array │ │      │ │      │ │Types│ │face  ││
│  └──────┘ └──────┘ └──────┘ └──────┘ └──────┘ └──────┘│
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│                  Utility Layer                           │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌────────────┐│
│  │Reflection│ │  Error   │ │ Circular │ │   String   ││
│  │ Helpers  │ │ Handler  │ │ Detection│ │  Matching  ││
│  └──────────┘ └──────────┘ └──────────┘ └────────────┘│
└─────────────────────────────────────────────────────────┘
```

## Key Components

### 1. Mapper

The main entry point that orchestrates the mapping process.

**Responsibilities:**

- Manages configuration
- Coordinates object pool
- Validates inputs
- Delegates to context

**Thread Safety:**

- Immutable configuration after creation
- Object pool with sync.Pool
- Each operation gets isolated context

### 2. Context

Per-operation state manager that tracks mapping progress.

**Responsibilities:**

- Tracks visited pointers (circular detection)
- Manages depth counter
- Collects errors
- Provides thread-safe access to shared state

**Lifecycle:**

- Created from pool at operation start
- Reset before reuse
- Returned to pool after operation

### 3. Configuration

Immutable configuration object with validation.

**Design Decisions:**

- Validated at creation time
- Cloned when needed
- Default values for safety
- Functional options pattern for API

### 4. Type Handlers

Specialized handlers for each reflect.Kind.

**Handler Types:**

- `mapStruct`: Handles struct-to-struct mapping
- `mapSlice`: Handles slices and arrays
- `mapMap`: Handles map types
- `mapPointer`: Handles pointer dereferencing
- `mapInterface`: Handles interface types
- `mapBasic`: Handles primitive types

**Common Pattern:**

```go
func (ctx *context) mapXXX(dst, src reflect.Value) error {
    // 1. Validate inputs
    // 2. Handle special cases
    // 3. Perform mapping
    // 4. Recursive calls for nested types
    return nil
}
```

## Design Patterns

### 1. Functional Options

Provides flexible, extensible configuration:

```go
type Option func(*Config)

func WithMaxDepth(depth int) Option {
    return func(c *Config) {
        c.MaxDepth = depth
    }
}
```

**Benefits:**

- Backward compatible
- Self-documenting
- Optional parameters
- Chainable

### 2. Object Pool

Reuses context objects to reduce allocations:

```go
pool: &sync.Pool{
    New: func() interface{} {
        return &context{
            visited: make(map[uintptr]reflect.Value),
            errors:  make([]error, 0),
        }
    },
}
```

**Benefits:**

- Reduced GC pressure
- Better performance
- Memory efficiency

### 3. Strategy Pattern

Custom converters allow pluggable behavior:

```go
type ConverterFunc func(src reflect.Value) (reflect.Value, error)

CustomConverters: map[reflect.Type]ConverterFunc
```

**Benefits:**

- Extensible
- Type-safe
- User-defined logic

### 4. Visitor Pattern

Recursive traversal of nested structures:

```go
func (ctx *context) mapValue(dst, src reflect.Value) error {
    // Visit node and recurse
    switch src.Kind() {
    case reflect.Struct:
        return ctx.mapStruct(dst, src)
    // ...
    }
}
```

## Performance Optimizations

### 1. Object Pooling

- Reuse context objects
- Reduce allocations
- ~30-40% performance improvement

### 2. Lazy Initialization

- Maps and slices created only when needed
- Nil checks before allocation

### 3. Early Returns

- Quick validation checks
- Skip unnecessary work
- Fail fast on errors

### 4. Inline Functions

- Small utility functions inlined
- Reduced function call overhead

### 5. Batch Operations

- Process collections efficiently
- Minimize lock contention

## Error Handling Strategy

### Three-Tier Approach

**1. Sentinel Errors**

```go
var ErrNilPointer = errors.New("mapper: nil pointer")
```

- Common, expected errors
- Easy to check with errors.Is()

**2. Detailed Errors**

```go
type MapError struct {
    Err      error
    SrcField string
    DstField string
}
```

- Context-rich error information
- Wraps underlying errors
- Useful for debugging

**3. Error Collection**

```go
ctx.errors []error
```

- Continue mapping after errors
- Report all errors at once
- Optional via error handler

## Thread Safety

### Immutable After Creation

- Config is read-only after creation
- Mapper is thread-safe

### Per-Operation Isolation

- Each operation gets own context
- No shared mutable state

### Synchronization Primitives

- sync.RWMutex for visited map
- sync.Pool for context reuse
- Atomic operations where applicable

## Testing Strategy

### Unit Tests

- Test each handler independently
- Edge cases and error conditions
- Table-driven tests

### Integration Tests

- End-to-end scenarios
- Real-world use cases
- Complex nested structures

### Benchmark Tests

- Performance regression detection
- Memory allocation tracking
- Comparison benchmarks

### Property-Based Tests

- Random input generation
- Invariant checking
- Edge case discovery

## Future Enhancements

### Planned Features

1. Field-level transformations
2. Validation hooks
3. Schema-based mapping
4. Code generation option
5. Streaming support for large datasets

### Performance Targets

- Sub-microsecond simple mappings
- Linear scaling with data size
- Minimal memory overhead

## Contributing

When adding features:

1. Maintain backward compatibility
2. Add comprehensive tests
3. Update documentation
4. Consider performance impact
5. Follow existing patterns
