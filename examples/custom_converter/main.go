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
	fmt.Println("=== Custom Converter Examples ===")

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
	fmt.Println("\n1. Time Formatting:")

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

	timeConverter := func(v reflect.Value) (reflect.Value, error) {
		if t, ok := v.Interface().(time.Time); ok {
			return reflect.ValueOf(t.Format("2006-01-02 15:04:05")), nil
		}
		return v, nil
	}

	var dst EventResponse
	if err := mapper.Copy(&dst, src,
		mapper.WithCustomConverter(reflect.TypeOf(time.Time{}), timeConverter),
	); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Source: %+v\n", src)
	fmt.Printf("Response: %+v\n", dst)
}

func priceConversion() {
	fmt.Println("\n2. Price Conversion (cents to dollars):")

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

	priceConverter := func(v reflect.Value) (reflect.Value, error) {
		if v.Kind() == reflect.Int {
			return reflect.ValueOf(float64(v.Int()) / 100.0), nil
		}
		return v, nil
	}

	var dst ProductDTO
	if err := mapper.Copy(&dst, src,
		mapper.WithCustomConverter(reflect.TypeOf(0), priceConverter),
	); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Source: %+v\n", src)
	fmt.Printf("DTO: %+v\n", dst)
}

func stringTransformations() {
	fmt.Println("\n3. String Transformations:")

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

	stringConverter := func(v reflect.Value) (reflect.Value, error) {
		if v.Kind() == reflect.String {
			return reflect.ValueOf(strings.ToLower(strings.TrimSpace(v.String()))), nil
		}
		return v, nil
	}

	var dst UserDTO
	if err := mapper.Copy(&dst, src,
		mapper.WithCustomConverter(reflect.TypeOf(""), stringConverter),
	); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Source: %+v\n", src)
	fmt.Printf("DTO (normalized): %+v\n", dst)
}

func complexTypeConversion() {
	fmt.Println("\n4. Complex Type Conversion:")

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

	intConverter := func(v reflect.Value) (reflect.Value, error) {
		if v.Kind() == reflect.Int {
			num := v.Int()
			if num >= 0 && num <= 2 {
				statuses := []string{"pending", "processing", "completed"}
				return reflect.ValueOf(statuses[num]), nil
			}
			return reflect.ValueOf(strconv.FormatInt(num, 10)), nil
		}
		return v, nil
	}

	var dst OrderDTO
	if err := mapper.Copy(&dst, src,
		mapper.WithCustomConverter(reflect.TypeOf(0), intConverter),
	); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Source: %+v\n", src)
	fmt.Printf("DTO: %+v\n", dst)
}
