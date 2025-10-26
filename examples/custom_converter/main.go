package customconverter

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/fbarikzehi/gomap/mapper"
)

func main() {
	fmt.Println("=== Custom Converter Examples ===\n")

	// Example 1: Time formatting
	timeFormatting()

	// Example 2: Price conversion
	priceConversion()

	// Example 3: String transformations
	stringTransformations()

	// Example 4: Complex type conversion
	complexTypeConversion()
}

func timeFormatting() {
	fmt.Println("1. Time Formatting:")

	type Event struct {
		Name      string
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	type EventResponse struct {
		Name      string
		CreatedAt string
		UpdatedAt string
	}

	src := Event{
		Name:      "Conference 2024",
		CreatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		UpdatedAt: time.Now(),
	}

	// Custom time formatter
	timeConverter := func(v reflect.Value) (reflect.Value, error) {
		if t, ok := v.Interface().(time.Time); ok {
			formatted := t.Format("2006-01-02 15:04:05")
			return reflect.ValueOf(formatted), nil
		}
		return v, nil
	}

	var dst EventResponse
	err := mapper.Copy(&dst, src,
		mapper.WithCustomConverter(reflect.TypeOf(time.Time{}), timeConverter),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Source: %+v\n", src)
	fmt.Printf("Response: %+v\n", dst)
}

func priceConversion() {
	fmt.Println("2. Price Conversion (cents to dollars):")

	type Product struct {
		Name  string
		Price int // cents
	}

	type ProductDTO struct {
		Name  string
		Price float64 // dollars
	}

	src := Product{
		Name:  "Widget",
		Price: 1299, // $12.99
	}

	// Convert cents to dollars
	priceConverter := func(v reflect.Value) (reflect.Value, error) {
		if v.Kind() == reflect.Int {
			cents := v.Int()
			dollars := float64(cents) / 100.0
			return reflect.ValueOf(dollars), nil
		}
		return v, nil
	}

	var dst ProductDTO
	err := mapper.Copy(&dst, src,
		mapper.WithCustomConverter(reflect.TypeOf(0), priceConverter),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Source: %+v\n", src)
	fmt.Printf("DTO: %+v\n", dst)
}

func stringTransformations() {
	fmt.Println("3. String Transformations:")

	type User struct {
		Username string
		Email    string
		Bio      string
	}

	type UserDTO struct {
		Username string
		Email    string
		Bio      string
	}

	src := User{
		Username: "john_doe",
		Email:    "JOHN@EXAMPLE.COM",
		Bio:      "   Software Engineer   ",
	}

	// String normalizer
	stringConverter := func(v reflect.Value) (reflect.Value, error) {
		if v.Kind() == reflect.String {
			s := v.String()
			// Lowercase and trim
			normalized := strings.ToLower(strings.TrimSpace(s))
			return reflect.ValueOf(normalized), nil
		}
		return v, nil
	}

	var dst UserDTO
	err := mapper.Copy(&dst, src,
		mapper.WithCustomConverter(reflect.TypeOf(""), stringConverter),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Source: %+v\n", src)
	fmt.Printf("DTO (normalized): %+v\n", dst)
}

func complexTypeConversion() {
	fmt.Println("4. Complex Type Conversion:")

	type Order struct {
		OrderID    int
		CustomerID string
		Status     int // 0=pending, 1=processing, 2=completed
		Total      int // cents
	}

	type OrderDTO struct {
		OrderID    string
		CustomerID string
		Status     string
		Total      string
	}

	src := Order{
		OrderID:    12345,
		CustomerID: "CUST-001",
		Status:     1,
		Total:      5999,
	}

	// Int to formatted string
	intConverter := func(v reflect.Value) (reflect.Value, error) {
		if v.Kind() == reflect.Int {
			num := v.Int()

			// Check field context (simplified - in real usage you'd need more context)
			if num < 3 { // Assume it's status
				statuses := []string{"pending", "processing", "completed"}
				if num < int64(len(statuses)) {
					return reflect.ValueOf(statuses[num]), nil
				}
			}

			// Default: convert to string
			return reflect.ValueOf(strconv.FormatInt(num, 10)), nil
		}
		return v, nil
	}

	var dst OrderDTO
	err := mapper.Copy(&dst, src,
		mapper.WithCustomConverter(reflect.TypeOf(0), intConverter),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Source: %+v\n", src)
	fmt.Printf("DTO: %+v\n", dst)
}
