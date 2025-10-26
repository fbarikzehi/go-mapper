package advanced

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/fbarikzehi/gomap/mapper"
)

func main() {
	fmt.Println("=== Advanced Mapping Examples ===")

	// Example 1: Tag-based mapping
	tagBasedMapping()

	// Example 2: Custom converters
	customConverters()

	// Example 3: Case-insensitive mapping
	caseInsensitiveMapping()

	// Example 4: Field name transformation
	fieldNameTransformation()

	// Example 5: Error handling
	errorHandling()

	// Example 6: Reusable mapper
	reusableMapper()
}

func tagBasedMapping() {
	fmt.Println("\n1. Tag-Based Mapping:")

	type SourceDTO struct {
		FullName    string `mapper:"name"`
		YearsOld    int    `mapper:"age"`
		ContactMail string `mapper:"email"`
	}

	type DestinationModel struct {
		Name  string
		Age   int
		Email string
	}

	src := SourceDTO{
		FullName:    "Alice Johnson",
		YearsOld:    35,
		ContactMail: "alice@example.com",
	}

	var dst DestinationModel
	if err := mapper.Copy(&dst, src, mapper.WithTagName("mapper")); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Source: %+v\n", src)
	fmt.Printf("Destination: %+v\n", dst)
}

func customConverters() {
	fmt.Println("\n2. Custom Type Converters:")

	type Event struct {
		Name      string
		Timestamp time.Time
		Priority  int
	}

	type EventDTO struct {
		Name      string
		Timestamp string
		Priority  string
	}

	src := Event{
		Name:      "System Update",
		Timestamp: time.Now(),
		Priority:  1,
	}

	timeConverter := func(v reflect.Value) (reflect.Value, error) {
		if t, ok := v.Interface().(time.Time); ok {
			return reflect.ValueOf(t.Format("2006-01-02 15:04:05")), nil
		}
		return v, nil
	}

	intConverter := func(v reflect.Value) (reflect.Value, error) {
		if v.Kind() == reflect.Int {
			return reflect.ValueOf(fmt.Sprintf("Level %d", v.Int())), nil
		}
		return v, nil
	}

	var dst EventDTO
	if err := mapper.Copy(&dst, src,
		mapper.WithCustomConverter(reflect.TypeOf(time.Time{}), timeConverter),
		mapper.WithCustomConverter(reflect.TypeOf(0), intConverter),
	); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Source: %+v\n", src)
	fmt.Printf("Destination: %+v\n", dst)
}

func caseInsensitiveMapping() {
	fmt.Println("\n3. Case-Insensitive Mapping:")

	type APIResponse struct {
		USERID    int
		USERNAME  string
		USEREMAIL string
	}

	type User struct {
		UserId    int
		UserName  string
		UserEmail string
	}

	src := APIResponse{
		USERID:    12345,
		USERNAME:  "bob_smith",
		USEREMAIL: "bob@example.com",
	}

	var dst User
	if err := mapper.Copy(&dst, src, mapper.WithCaseSensitive(false)); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Source: %+v\n", src)
	fmt.Printf("Destination: %+v\n", dst)
}

func fieldNameTransformation() {
	fmt.Println("\n4. Field Name Transformation:")

	type SnakeCaseModel struct {
		FirstName    string
		LastName     string
		EmailAddress string
	}

	type CamelCaseModel struct {
		FirstName    string
		LastName     string
		EmailAddress string
	}

	src := SnakeCaseModel{
		FirstName:    "Charlie",
		LastName:     "Brown",
		EmailAddress: "charlie@example.com",
	}

	fieldMapper := func(name string) string {
		// In real scenarios, convert snake_case to camelCase
		return name
	}

	var dst CamelCaseModel
	if err := mapper.Copy(&dst, src, mapper.WithFieldNameMapper(fieldMapper)); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Source: %+v\n", src)
	fmt.Printf("Destination: %+v\n", dst)
}

func errorHandling() {
	fmt.Println("\n5. Custom Error Handling:")

	type Source struct {
		Name  string
		Value int
	}

	type Destination struct {
		Name  string
		Value int
	}

	src := Source{
		Name:  "Test",
		Value: 42,
	}

	errorHandler := func(err error, srcField, dstField string) error {
		fmt.Printf("  Error mapping %s to %s: %v\n", srcField, dstField, err)
		return nil // Continue mapping
	}

	var dst Destination
	if err := mapper.Copy(&dst, src, mapper.WithErrorHandler(errorHandler)); err != nil {
		fmt.Printf("Mapping completed with errors: %v\n", err)
	} else {
		fmt.Printf("Mapping successful: %+v\n", dst)
	}
}

func reusableMapper() {
	fmt.Println("\n6. Reusable Mapper Configuration:")

	type Record struct {
		ID   int
		Name string
		Tags []string
	}

	m := mapper.NewMapper(
		mapper.WithMaxDepth(10),
		mapper.WithIgnoreUnexported(true),
		mapper.WithDeepCopy(true),
	)

	sources := []Record{
		{ID: 1, Name: "Record 1", Tags: []string{"tag1", "tag2"}},
		{ID: 2, Name: "Record 2", Tags: []string{"tag3", "tag4"}},
		{ID: 3, Name: "Record 3", Tags: []string{"tag5", "tag6"}},
	}

	destinations := make([]Record, len(sources))

	for i, src := range sources {
		if err := m.Map(&destinations[i], src); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf("Mapped %d records successfully\n", len(destinations))
	for i, dst := range destinations {
		fmt.Printf("  [%d] %+v\n", i, dst)
	}
}
