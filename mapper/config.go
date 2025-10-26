// Package mapper provides reflection-based object-to-object mapping utilities.
// This file defines the configuration system and related function types.
package mapper

import (
	"fmt"
	"maps"
	"reflect"
	"time"
)

// Configuration constants define mapper defaults and limits.
const (
	// DefaultMaxDepth is the default maximum depth for nested structure traversal.
	DefaultMaxDepth = 32

	// NoDepthLimit disables depth restriction (-1 means unlimited depth).
	NoDepthLimit = -1

	// DefaultTagName is the default struct tag key used for field mapping.
	DefaultTagName = "mapper"
)

// Config defines the complete configuration for a mapping operation.
//
// A Config can be customized using functional options or created manually.
// It controls depth limits, tag behavior, naming rules, converter functions,
// and advanced reflection behaviors.
//
// Example:
//
//	cfg := &mapper.Config{
//	    DeepCopy:         true,
//	    TagName:          "map",
//	    UseJSONTag:       true,
//	    IgnoreUnexported: true,
//	}
//
//	mapper.Copy(&dst, src, cfg)
type Config struct {
	// MaxDepth limits nested structure traversal depth.
	// Use NoDepthLimit (-1) for unlimited depth.
	MaxDepth int

	// TagName defines the struct tag key used to map field names.
	// If empty, tag-based mapping is disabled.
	TagName string

	// IgnoreUnexported skips unexported (private) fields during mapping.
	IgnoreUnexported bool

	// DeepCopy enables deep copying of struct fields and nested types.
	DeepCopy bool

	// ZeroFields sets destination fields to their zero value
	// when the corresponding source field is zero.
	ZeroFields bool

	// IgnoreNilFields skips mapping of nil pointer fields from the source.
	IgnoreNilFields bool

	// CaseSensitive enables case-sensitive field name matching.
	CaseSensitive bool

	// UseJSONTag allows JSON tag parsing (e.g., `json:"name"`) for field mapping.
	UseJSONTag bool

	// SkipCircularCheck disables circular reference detection.
	// Only disable this if you are certain your data has no circular references.
	SkipCircularCheck bool

	// CustomConverters defines per-type converter functions used
	// to transform values before assignment.
	CustomConverters map[reflect.Type]ConverterFunc

	// FieldNameMapper transforms field names between source and destination structs.
	FieldNameMapper FieldNameMapperFunc

	// ErrorHandler defines how errors encountered during mapping are handled.
	// Return nil to continue mapping despite the error.
	ErrorHandler ErrorHandlerFunc

	// TimeLayout specifies the layout string used for time.Time conversions.
	TimeLayout string

	// MaxSliceCapacity limits the maximum capacity allocated for slices.
	// Protects against excessive memory allocation.
	MaxSliceCapacity int

	// AllowPrivateFields enables copying of private/unexported fields via reflection.
	// ⚠️ Use with caution — this breaks encapsulation.
	AllowPrivateFields bool
}

// ConverterFunc defines a custom conversion function that transforms
// a reflected value into another reflected value (potentially of a different type).
type ConverterFunc func(src reflect.Value) (reflect.Value, error)

// FieldNameMapperFunc defines a function that transforms field names during mapping,
// allowing for case normalization, prefix/suffix handling, etc.
type FieldNameMapperFunc func(fieldName string) string

// ErrorHandlerFunc defines how mapping errors are processed.
//
// If the function returns nil, the mapper continues execution;
// otherwise, mapping is stopped and the returned error is propagated.
type ErrorHandlerFunc func(err error, srcField, dstField string) error

// defaultConfig returns a new Config populated with sensible defaults.
//
// It is used internally when no user configuration is provided.
func defaultConfig() *Config {
	return &Config{
		MaxDepth:           DefaultMaxDepth,
		TagName:            DefaultTagName,
		IgnoreUnexported:   true,
		DeepCopy:           true,
		CaseSensitive:      true,
		CustomConverters:   make(map[reflect.Type]ConverterFunc),
		TimeLayout:         time.RFC3339,
		MaxSliceCapacity:   1_048_576, // 1M elements
		AllowPrivateFields: false,
	}
}

// clone creates a deep copy of the configuration, including
// its CustomConverters map (by value, not by reference).
func (c *Config) clone() *Config {
	converters := make(map[reflect.Type]ConverterFunc, len(c.CustomConverters))
	maps.Copy(converters, c.CustomConverters)

	return &Config{
		MaxDepth:           c.MaxDepth,
		TagName:            c.TagName,
		IgnoreUnexported:   c.IgnoreUnexported,
		DeepCopy:           c.DeepCopy,
		ZeroFields:         c.ZeroFields,
		IgnoreNilFields:    c.IgnoreNilFields,
		CaseSensitive:      c.CaseSensitive,
		UseJSONTag:         c.UseJSONTag,
		SkipCircularCheck:  c.SkipCircularCheck,
		CustomConverters:   converters,
		FieldNameMapper:    c.FieldNameMapper,
		ErrorHandler:       c.ErrorHandler,
		TimeLayout:         c.TimeLayout,
		MaxSliceCapacity:   c.MaxSliceCapacity,
		AllowPrivateFields: c.AllowPrivateFields,
	}
}

// validate ensures that the configuration values are valid.
// Returns an error if any field contains an invalid setting.
func (c *Config) validate() error {
	if c.MaxDepth < NoDepthLimit {
		return fmt.Errorf("mapper: invalid MaxDepth %d (must be >= -1)", c.MaxDepth)
	}
	if c.MaxSliceCapacity < 0 {
		return fmt.Errorf("mapper: invalid MaxSliceCapacity %d (must be >= 0)", c.MaxSliceCapacity)
	}
	return nil
}
