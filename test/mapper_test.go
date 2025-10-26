package gomap_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/fbarikzehi/gomap/mapper"
)

type TestPerson struct {
	Name     string
	Age      int
	Email    string
	Created  time.Time
	Address  *TestAddress
	Tags     []string
	Metadata map[string]interface{}
}

type TestAddress struct {
	Street  string
	City    string
	Country string
}

func TestBasicMapping(t *testing.T) {
	tests := []struct {
		name    string
		src     interface{}
		dst     interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name: "simple struct",
			src: TestPerson{
				Name:  "John",
				Age:   30,
				Email: "john@example.com",
			},
			dst:     &TestPerson{},
			wantErr: false,
		},
		{
			name: "nested struct with pointers",
			src: TestPerson{
				Name: "Alice",
				Address: &TestAddress{
					Street:  "123 Main St",
					City:    "NY",
					Country: "USA",
				},
			},
			dst:     &TestPerson{},
			wantErr: false,
		},
		// Add more test cases...
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mapper.Copy(tt.dst, tt.src)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.src, reflect.ValueOf(tt.dst).Elem().Interface())
		})
	}
}

func TestCustomConverters(t *testing.T) {
	type TimeStruct struct {
		Created time.Time
	}

	type StringStruct struct {
		Created string
	}

	timeConverter := func(v reflect.Value) (reflect.Value, error) {
		if t, ok := v.Interface().(time.Time); ok {
			return reflect.ValueOf(t.Format(time.RFC3339)), nil
		}
		return v, nil
	}

	src := TimeStruct{Created: time.Now()}
	var dst StringStruct

	err := mapper.Copy(&dst, src, mapper.WithCustomConverter(reflect.TypeOf(time.Time{}), timeConverter))
	require.NoError(t, err)
	assert.NotEmpty(t, dst.Created)
}
