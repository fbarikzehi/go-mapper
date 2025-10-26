package gomap_test

import (
	"testing"
	"time"

	"github.com/fbarikzehi/gomap/mapper"
)

type BenchPerson struct {
	ID        int64
	Name      string
	Email     string
	Age       int
	CreatedAt time.Time
	Address   *BenchAddress
	Tags      []string
	Meta      map[string]interface{}
}

type BenchAddress struct {
	Street  string
	City    string
	Country string
	Zip     string
}

func BenchmarkSimpleStruct(b *testing.B) {
	src := BenchPerson{
		ID:    1,
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   30,
	}
	var dst BenchPerson

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mapper.Copy(&dst, src)
	}
}

func BenchmarkReusableMapper(b *testing.B) {
	m := mapper.NewMapper()
	src := BenchPerson{
		ID:    1,
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   30,
	}
	var dst BenchPerson

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Map(&dst, src)
	}
}

func BenchmarkDeepStruct(b *testing.B) {
	src := BenchPerson{
		ID:    1,
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   30,
		Address: &BenchAddress{
			Street:  "123 Main St",
			City:    "New York",
			Country: "USA",
			Zip:     "10001",
		},
		Tags: []string{"tag1", "tag2", "tag3"},
		Meta: map[string]interface{}{
			"key1": "value1",
			"key2": 123,
			"key3": true,
		},
	}
	var dst BenchPerson

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mapper.Copy(&dst, src)
	}
}
