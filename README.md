# Go Map

[![Go Reference](https://pkg.go.dev/badge/github.com/fbarikzehi/gomap.svg)](https://pkg.go.dev/github.com/fbarikzehi/gomap)
[![Go Report Card](https://goreportcard.com/badge/github.com/fbarikzehi/gomap)](https://goreportcard.com/report/github.com/fbarikzehi/gomap)
[![CI](https://github.com/fbarikzehi/gomap/workflows/CI/badge.svg)](https://github.com/fbarikzehi/gomap/actions)
[![codecov](https://codecov.io/gh/fbarikzehi/gomap/branch/main/graph/badge.svg)](https://codecov.io/gh/fbarikzehi/gomap)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A high-performance, feature-rich Go library for mapping data between structs with deep copy support, customizable field mapping, and comprehensive type handling.

## Features

- 🚀 **High Performance** - Optimized with object pooling and minimal allocations
- 🔒 **Thread-Safe** - Safe for concurrent use
- 🔄 **Deep Copy** - Full recursive copying with circular reference detection
- 🎯 **Type Conversion** - Intelligent type conversion and custom converters
- 🏷️ **Tag Support** - Custom struct tags and JSON tag support
- 🔍 **Flexible Matching** - Case-sensitive and case-insensitive field matching
- 🛡️ **Error Handling** - Comprehensive error reporting with custom handlers
- 📦 **Zero Dependencies** - Pure Go implementation

## Installation

```bash
go get github.com/fbarikzehi/gomap
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/fbarikzehi/gomap"
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
    if err := mapper.Copy(&dst, src); err != nil {
        panic(err)
    }

    fmt.Printf("%+v\n", dst)
    // Output: {Name:Foo Bar Age:30 Email:foo@example.com}
}
```

## Documentation

- [Examples](docs/examples.md) - Comprehensive usage examples
- [Benchmarks](docs/benchmarks.md) - Performance benchmarks and optimization tips
- [API Reference](https://pkg.go.dev/github.com/fbarikzehi/gomap) - Complete API documentation

## Usage

### Basic Mapping

```go
var dst Destination
err := mapper.Copy(&dst, src)
```

### Reusable Mapper

```go
m := mapper.NewMapper(
    mapper.WithMaxDepth(10),
    mapper.WithIgnoreUnexported(true),
)

var dst Destination
err := m.Map(&dst, src)
```

### Tag-Based Mapping

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

### Custom Converters

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

### Case-Insensitive Mapping

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

## Configuration Options

| Option                        | Description                         | Default  |
| ----------------------------- | ----------------------------------- | -------- |
| `WithMaxDepth(int)`           | Maximum depth for nested structures | 32       |
| `WithTagName(string)`         | Custom struct tag name              | "mapper" |
| `WithIgnoreUnexported(bool)`  | Skip unexported fields              | true     |
| `WithDeepCopy(bool)`          | Enable deep copying                 | true     |
| `WithZeroFields(bool)`        | Zero destination on source zero     | false    |
| `WithIgnoreNilFields(bool)`   | Skip nil pointer fields             | false    |
| `WithCaseSensitive(bool)`     | Case-sensitive field matching       | true     |
| `WithJSONTag(bool)`           | Use JSON tags for mapping           | false    |
| `WithSkipCircularCheck(bool)` | Skip circular reference check       | false    |

## Performance

```
BenchmarkSimpleMapping-8       5000000    250 ns/op    64 B/op    2 allocs/op
BenchmarkNestedMapping-8       2000000    650 ns/op   192 B/op    6 allocs/op
BenchmarkSliceMapping-8         500000   3200 ns/op  2048 B/op   52 allocs/op
```

See [benchmarks](docs/benchmarks.md) for detailed performance analysis.

## Contributing

Contributions are welcome! Please read our [Contributing Guide](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by popular mapping libraries in other languages
- Built with performance and developer experience in mind
- Community-driven development

## Support

- 📫 [Open an issue](https://github.com/fbarikzehi/gomap/issues)
- 💬 [Start a discussion](https://github.com/fbarikzehi/gomap/discussions)
- ⭐ Star the project if you find it useful!
