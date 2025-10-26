package basic

import (
	"fmt"
	"log"

	"github.com/fbarikzehi/gomap/mapper"
)

// Source represents the source data structure
type Source struct {
	Name    string
	Age     int
	Email   string
	Address Address
}

// Address represents an address
type Address struct {
	Street  string
	City    string
	Country string
}

// Destination represents the destination data structure
type Destination struct {
	Name    string
	Age     int
	Email   string
	Address Address
}

func main() {
	fmt.Println("=== Basic Mapping Examples ===")

	// Example 1: Simple mapping
	simpleMapping()

	// Example 2: Nested struct mapping
	nestedMapping()

	// Example 3: Slice mapping
	sliceMapping()

	// Example 4: Map mapping
	mapMapping()
}

func simpleMapping() {
	fmt.Println("\n1. Simple Mapping:")

	src := Source{
		Name:  "John Doe",
		Age:   30,
		Email: "john@example.com",
	}

	var dst Destination
	if err := mapper.Copy(&dst, src); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Source: %+v\n", src)
	fmt.Printf("Destination: %+v\n", dst)
}

func nestedMapping() {
	fmt.Println("\n2. Nested Struct Mapping:")

	src := Source{
		Name:  "Jane Smith",
		Age:   28,
		Email: "jane@example.com",
		Address: Address{
			Street:  "123 Main St",
			City:    "New York",
			Country: "USA",
		},
	}

	var dst Destination
	if err := mapper.Copy(&dst, src); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Source: %+v\n", src)
	fmt.Printf("Destination: %+v\n", dst)
}

func sliceMapping() {
	fmt.Println("\n3. Slice Mapping:")

	type Item struct {
		ID   int
		Name string
	}

	type Container struct {
		Items []Item
	}

	src := Container{
		Items: []Item{
			{ID: 1, Name: "Item 1"},
			{ID: 2, Name: "Item 2"},
			{ID: 3, Name: "Item 3"},
		},
	}

	var dst Container
	if err := mapper.Copy(&dst, src); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Source items: %d\n", len(src.Items))
	fmt.Printf("Destination items: %d\n", len(dst.Items))
	for i, item := range dst.Items {
		fmt.Printf("  [%d] %+v\n", i, item)
	}
}

func mapMapping() {
	fmt.Println("\n4. Map Mapping:")

	type Config struct {
		Settings map[string]string
	}

	src := Config{
		Settings: map[string]string{
			"theme":    "dark",
			"language": "en",
			"timezone": "UTC",
		},
	}

	var dst Config
	if err := mapper.Copy(&dst, src); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Source settings: %+v\n", src.Settings)
	fmt.Printf("Destination settings: %+v\n", dst.Settings)
}
