# Benchmarks and Performance Guide

Performance benchmarks and optimization tips for gomap.

## Benchmark Results

### Simple Struct Mapping

goos: linux
goarch: amd64
pkg: github.com/fbarikzehi/gomap/test
cpu: Intel(R) Core(TM) i7-7820HK CPU @ 2.90GHz

```
BenchmarkSimpleMapping-8              5000000      250 ns/op      64 B/op    2 allocs/op
BenchmarkSimpleMappingReuse-8        10000000      180 ns/op      32 B/op    1 allocs/op
```

### Nested Structures

```
BenchmarkNestedMapping-8              2000000      650 ns/op     192 B/op    6 allocs/op
BenchmarkNestedMappingReuse-8         3000000      480 ns/op     128 B/op    4 allocs/op
```

### Collections

```
BenchmarkSliceMapping10-8             1000000     1200 ns/op     512 B/op   12 allocs/op
BenchmarkSliceMapping100-8             100000    12000 ns/op    5120 B/op  112 allocs/op
BenchmarkMapMapping-8                  500000     2800 ns/op    1024 B/op   24 allocs/op
```

### Deep Copy vs Shallow Copy

```
BenchmarkDeepCopy-8                   2000000      700 ns/op     256 B/op    8 allocs/op
BenchmarkShallowCopy-8                5000000      300 ns/op      96 B/op    3 allocs/op
```

## Performance Tips

### 1. Reuse Mapper Instances

**Bad:**

```go
for _, item := range items {
    var dst Destination
    mapper.Copy(&dst, item) // Creates new mapper each time
}
```

**Good:**

```go
m := mapper.NewMapper()
for _, item := range items {
    var dst Destination
    m.Map(&dst, item) // Reuses mapper
}
```

**Impact:** 30-40% faster

### 2. Skip Circular Reference Checks

If you're certain there are no circular references:

```go
m := mapper.NewMapper(mapper.WithSkipCircularCheck(true))
```

**Impact:** 10-15% faster

### 3. Set Appropriate Max Depth

```go
m := mapper.NewMapper(mapper.WithMaxDepth(5)) // Instead of default 32
```

**Impact:** 5-10% faster for shallow structures

### 4. Pre-allocate Slices

**Bad:**

```go
type Destination struct {
    Items []Item // Will allocate during mapping
}
```

**Good:**

```go
dst := Destination{
    Items: make([]Item, len(src.Items)),
}
mapper.Copy(&dst, src)
```

**Impact:** 20-30% faster for large slices

### 5. Use Specific Options

Only enable features you need:

```go
// Minimal config for maximum performance
m := mapper.NewMapper(
    mapper.WithMaxDepth(5),
    mapper.WithIgnoreUnexported(true),
    mapper.WithSkipCircularCheck(true),
)
```

## Running Benchmarks

```bash
# Run all benchmarks
make bench

# Run specific benchmark
go test -bench=BenchmarkSimpleMapping -benchmem

# Compare benchmarks
go test -bench=. -benchmem | tee new.txt
benchstat old.txt new.txt
```

## Memory Profiling

```bash
go test -bench=. -memprofile=mem.prof
go tool pprof mem.prof
```

---
