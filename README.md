# Struct Mapper

A high-performance, feature-rich Go package for mapping data between structs with deep copy support, customizable field mapping, and comprehensive type handling.

## Features

- ✅ **Deep Copying**: Full deep copy support for nested structures
- ✅ **Type Safety**: Comprehensive type checking and conversion
- ✅ **Performance Optimized**: Object pooling and efficient memory usage
- ✅ **Circular Reference Detection**: Prevents infinite loops
- ✅ **Custom Converters**: Define custom type conversion logic
- ✅ **Tag Support**: Use struct tags for field mapping
- ✅ **Case-Insensitive Mapping**: Optional case-insensitive field matching
- ✅ **Pointer Handling**: Smart pointer dereferencing and allocation
- ✅ **Collection Support**: Maps, slices, and arrays
- ✅ **Error Handling**: Comprehensive error reporting
- ✅ **Concurrent Safe**: Thread-safe operations
- ✅ **Zero Allocation Options**: Minimize memory allocations

## Installation

```bash
go get https://github.com/fbarikzehi/go-mapper/mapper
```

## Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
    "https://github.com/fbarikzehi/go-mapper/mapper"
)

type Source struct {
    Name  string
    Age   int
    Email string
}

type Destination struct {
    Name  string
    Age   int
    Email string
}

func main() {
    src := Source{
        Name:  "Foo Bar",
        Age:   30,
        Email: "foo@example.com",
    }

    var dst Destination
    err := mapper.Copy(&dst, src)
    if err != nil {
        panic(err)
    }

    fmt.Printf("%+v\n", dst)
}
```

### Reusable Mapper

```go
// Create mapper once, reuse multiple times
m := mapper.NewMapper(
    mapper.WithMaxDepth(10),
    mapper.WithIgnoreUnexported(true),
)

var dst Destination
err := m.Map(&dst, src)
```

## Configuration Options

### WithMaxDepth

Limit the depth of nested structure traversal:

```go
mapper.Copy(&dst, src, mapper.WithMaxDepth(5))
```

### WithTagName

Use custom struct tags for field mapping:

```go
type Source struct {
    FullName string `mapper:"name"`
    Years    int    `mapper:"age"`
}

type Destination struct {
    Name string
    Age  int
}

mapper.Copy(&dst, src, mapper.WithTagName("mapper"))
```

### WithIgnoreUnexported

Control whether to copy unexported fields:

```go
mapper.Copy(&dst, src, mapper.WithIgnoreUnexported(true))
```

### WithDeepCopy

Enable or disable deep copying:

```go
mapper.Copy(&dst, src, mapper.WithDeepCopy(true))
```

### WithZeroFields

Handle zero-value fields:

```go
mapper.Copy(&dst, src, mapper.WithZeroFields(true))
```

### WithIgnoreNilFields

Skip nil pointer fields:

```go
mapper.Copy(&dst, src, mapper.WithIgnoreNilFields(true))
```

### WithCaseSensitive

Enable case-insensitive field matching:

```go
type Source struct {
    USERNAME string
    USERID   int
}

type Destination struct {
    Username string
    UserId   int
}

mapper.Copy(&dst, src, mapper.WithCaseSensitive(false))
```

### WithJSONTag

Use JSON tags for field mapping:

```go
type Source struct {
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
}

mapper.Copy(&dst, src, mapper.WithJSONTag(true))
```

### WithCustomConverter

Define custom type converters:

```go
import "time"

timeConverter := func(v reflect.Value) (reflect.Value, error) {
    if t, ok := v.Interface().(time.Time); ok {
        formatted := t.Format("2006-01-02")
        return reflect.ValueOf(formatted), nil
    }
    return v, nil
}

mapper.Copy(&dst, src,
    mapper.WithCustomConverter(reflect.TypeOf(time.Time{}), timeConverter),
)
```

### WithFieldNameMapper

Transform field names during mapping:

```go
// Convert camelCase to snake_case
fieldMapper := func(fieldName string) string {
    result := ""
    for i, r := range fieldName {
        if i > 0 && r >= 'A' && r <= 'Z' {
            result += "_"
        }
        result += string(r)
    }
    return result
}

mapper.Copy(&dst, src, mapper.WithFieldNameMapper(fieldMapper))
```

### WithErrorHandler

Custom error handling:

```go
errorHandler := func(err error, srcField, dstField string) error {
    log.Printf("Error mapping %s to %s: %v", srcField, dstField, err)
    return nil // Continue mapping
}

mapper.Copy(&dst, src, mapper.WithErrorHandler(errorHandler))
```

## Advanced Usage

### Nested Structures

```go
type Address struct {
    Street  string
    City    string
    ZipCode string
}

type Person struct {
    Name    string
    Address Address
}

src := Person{
    Name: "Alice",
    Address: Address{
        Street:  "123 Main St",
        City:    "New York",
        ZipCode: "10001",
    },
}

var dst Person
mapper.Copy(&dst, src)
```

### Slices and Arrays

```go
type Source struct {
    Items []Item
}

type Destination struct {
    Items []Item
}

src := Source{
    Items: []Item{
        {ID: 1, Name: "Item 1"},
        {ID: 2, Name: "Item 2"},
    },
}

var dst Destination
mapper.Copy(&dst, src)
```

### Maps

```go
type Source struct {
    Attributes map[string]string
}

type Destination struct {
    Attributes map[string]string
}

src := Source{
    Attributes: map[string]string{
        "color": "blue",
        "size":  "large",
    },
}

var dst Destination
mapper.Copy(&dst, src)
```

### Pointers

```go
type Source struct {
    Name  string
    Value *int
}

type Destination struct {
    Name  string
    Value *int
}

value := 42
src := Source{
    Name:  "Test",
    Value: &value,
}

var dst Destination
mapper.Copy(&dst, src)

// dst.Value is a deep copy, different pointer than src.Value
```

## Use Cases

### 1. API Layer Mapping

```go
// Database model
type UserModel struct {
    ID        int
    Username  string
    Email     string
    Password  string // Don't expose
    CreatedAt time.Time
}

// API response
type UserResponse struct {
    ID        int    `json:"id"`
    Username  string `json:"username"`
    Email     string `json:"email"`
    CreatedAt string `json:"created_at"`
}

// Map database model to API response
func ToUserResponse(model UserModel) UserResponse {
    var response UserResponse

    timeConverter := func(v reflect.Value) (reflect.Value, error) {
        if t, ok := v.Interface().(time.Time); ok {
            return reflect.ValueOf(t.Format(time.RFC3339)), nil
        }
        return v, nil
    }

    mapper.Copy(&response, model,
        mapper.WithCustomConverter(reflect.TypeOf(time.Time{}), timeConverter),
    )

    return response
}
```

### 2. DTO to Domain Entity

```go
// Data Transfer Object
type CreateProductDTO struct {
    Name        string
    Description string
    Price       int // cents
    CategoryID  int
}

// Domain Entity
type Product struct {
    ID          int
    Name        string
    Description string
    Price       float64 // dollars
    CategoryID  int
    CreatedAt   time.Time
}

func CreateProduct(dto CreateProductDTO) Product {
    var product Product

    priceConverter := func(v reflect.Value) (reflect.Value, error) {
        if v.Kind() == reflect.Int {
            cents := v.Int()
            dollars := float64(cents) / 100.0
            return reflect.ValueOf(dollars), nil
        }
        return v, nil
    }

    mapper.Copy(&product, dto,
        mapper.WithCustomConverter(reflect.TypeOf(0), priceConverter),
    )

    product.CreatedAt = time.Now()

    return product
}
```

### 3. Configuration Merging

```go
type Config struct {
    Host     string
    Port     int
    Timeout  time.Duration
    Features map[string]bool
}

func MergeConfig(base, override Config) Config {
    var merged Config
    mapper.Copy(&merged, base)

    // Merge override values
    mapper.Copy(&merged, override, mapper.WithIgnoreNilFields(true))

    return merged
}
```

### 4. Caching Layer

```go
type CacheEntry struct {
    Key       string
    Value     interface{}
    ExpiresAt time.Time
}

type StoredCacheEntry struct {
    Key       string
    Value     []byte
    ExpiresAt int64
}

func ToCacheEntry(stored StoredCacheEntry) CacheEntry {
    var entry CacheEntry
    mapper.Copy(&entry, stored)
    return entry
}
```

## Performance Tips

1. **Reuse Mapper Instance**: Create once, use many times
2. **Use Specific Options**: Only enable features you need
3. **Avoid Deep Nesting**: Set appropriate MaxDepth
4. **Skip Circular Check**: If you're sure there are no circular references
5. **Batch Operations**: Map multiple items in a single session

```go
// Good: Reuse mapper
m := mapper.NewMapper(
    mapper.WithMaxDepth(5),
    mapper.WithSkipCircularCheck(true),
)

for _, src := range sources {
    var dst Destination
    m.Map(&dst, src)
}

// Avoid: Creating new mapper each time
for _, src := range sources {
    var dst Destination
    mapper.Copy(&dst, src) // Creates new mapper internally
}
```

## Error Handling

The package provides specific error types:

```go
var (
    ErrNilPointer         // nil pointer provided
    ErrUnsupportedType    // unsupported type
    ErrInvalidDestination // destination must be a pointer
    ErrTypeMismatch       // type mismatch
    ErrMaxDepthExceeded   // max depth exceeded
    ErrCircularReference  // circular reference detected
)
```

Handle errors appropriately:

```go
err := mapper.Copy(&dst, src)
if err != nil {
    switch {
    case errors.Is(err, mapper.ErrCircularReference):
        // Handle circular reference
    case errors.Is(err, mapper.ErrMaxDepthExceeded):
        // Handle depth exceeded
    default:
        // Handle other errors
    }
}
```

## Best Practices

1. **Always pass destination as pointer**: `mapper.Copy(&dst, src)`
2. **Initialize complex types**: Pre-allocate maps and slices when possible
3. **Use tags for clarity**: Document field mappings with struct tags
4. **Handle errors**: Don't ignore mapping errors in production
5. **Test edge cases**: Test with nil values, empty collections, circular references
6. **Profile performance**: Use benchmarks for critical paths
7. **Document converters**: Comment custom converter logic
8. **Version compatibility**: Consider struct versioning for APIs

## Thread Safety

The mapper is thread-safe and can be used concurrently:

```go
m := mapper.NewMapper()

var wg sync.WaitGroup
for i := 0; i < 100; i++ {
    wg.Add(1)
    go func(src Source) {
        defer wg.Done()
        var dst Destination
        m.Map(&dst, src)
    }(sources[i])
}
wg.Wait()
```

## Benchmarks

```
BenchmarkSimpleMapping-8       5000000    250 ns/op    64 B/op    2 allocs/op
BenchmarkNestedMapping-8       2000000    650 ns/op   192 B/op    6 allocs/op
BenchmarkSliceMapping-8         500000   3200 ns/op  2048 B/op   52 allocs/op
```

## License

MIT License

## Contributing

Contributions are welcome! Please submit pull requests or issues on GitHub.

## Support

For questions and support, please open an issue on GitHub.
