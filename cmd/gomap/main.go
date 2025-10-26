package main

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/fbarikzehi/gomap/mapper"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Command-line flags
	showVersion := flag.Bool("version", false, "Show gomap version")
	flag.Parse()

	if *showVersion {
		fmt.Printf("gomap version: %s (commit: %s, built at: %s)\n", version, commit, date)
		return
	}

	fmt.Println("🚀 gomap - Go struct field mapper")

	// Example usage of the mapper
	type Source struct {
		Name  string
		Email string
		Age   int
	}

	type Destination struct {
		Name  string
		Email string
		Age   int
	}

	src := Source{
		Name:  "Alice",
		Email: "alice@example.com",
		Age:   28,
	}

	var dst Destination

	err := mapper.Copy(&dst, src,
		mapper.WithDeepCopy(true),
		mapper.WithTimeLayout(time.RFC3339),
		mapper.WithCustomConverter(reflect.TypeOf(time.Time{}),
			func(v reflect.Value) (reflect.Value, error) {
				if t, ok := v.Interface().(time.Time); ok {
					return reflect.ValueOf(t.Format(time.RFC3339)), nil
				}
				return v, nil
			}),
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Mapping error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n✅ Mapped struct:\n%+v\n", dst)
}
