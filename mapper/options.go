// Package mapper provides functional configuration options
// for customizing the behavior of the Mapper.
//
// These options can be passed to NewMapper or Copy to control
// various aspects of mapping behavior, such as depth limits,
// tag usage, field matching rules, and custom type conversions.
//
// Example:
//
//	mapper.Copy(&dst, src,
//	    mapper.WithMaxDepth(5),
//	    mapper.WithIgnoreUnexported(true),
//	    mapper.WithJSONTag(true),
//	)
package mapper

import "reflect"

// Option represents a functional option for configuring a Mapper instance.
//
// Each Option modifies a Config struct, which determines how the Mapper
// performs struct-to-struct mapping operations.
type Option func(*Config)

// WithMaxDepth sets the maximum allowed depth for nested structure traversal.
// When the maximum depth is reached, mapping stops and returns ErrMaxDepthExceeded.
//
// Example:
//
//	mapper.Copy(&dst, src, mapper.WithMaxDepth(5))
func WithMaxDepth(depth int) Option {
	return func(c *Config) {
		c.MaxDepth = depth
	}
}

// WithTagName sets a custom struct tag name to use for field mapping.
//
// Example:
//
//	type Source struct {
//	    Name string `mapper:"full_name"`
//	}
//	mapper.Copy(&dst, src, mapper.WithTagName("mapper"))
func WithTagName(tag string) Option {
	return func(c *Config) {
		c.TagName = tag
	}
}

// WithIgnoreUnexported configures whether unexported fields (lowercase)
// should be ignored during mapping.
//
// Example:
//
//	mapper.Copy(&dst, src, mapper.WithIgnoreUnexported(true))
func WithIgnoreUnexported(ignore bool) Option {
	return func(c *Config) {
		c.IgnoreUnexported = ignore
	}
}

// WithDeepCopy enables or disables deep copying of complex types such as
// slices, maps, and nested structs.
//
// Example:
//
//	mapper.Copy(&dst, src, mapper.WithDeepCopy(true))
func WithDeepCopy(deep bool) Option {
	return func(c *Config) {
		c.DeepCopy = deep
	}
}

// WithZeroFields configures whether destination fields should be zeroed
// when the corresponding source field is a zero value.
//
// Example:
//
//	mapper.Copy(&dst, src, mapper.WithZeroFields(true))
func WithZeroFields(zero bool) Option {
	return func(c *Config) {
		c.ZeroFields = zero
	}
}

// WithIgnoreNilFields configures whether nil pointer fields in the source
// should be skipped during mapping.
//
// Example:
//
//	mapper.Copy(&dst, src, mapper.WithIgnoreNilFields(true))
func WithIgnoreNilFields(ignore bool) Option {
	return func(c *Config) {
		c.IgnoreNilFields = ignore
	}
}

// WithCaseSensitive controls whether field name matching is case-sensitive.
// If set to false, fields are matched case-insensitively.
//
// Example:
//
//	mapper.Copy(&dst, src, mapper.WithCaseSensitive(false))
func WithCaseSensitive(sensitive bool) Option {
	return func(c *Config) {
		c.CaseSensitive = sensitive
	}
}

// WithJSONTag enables support for JSON struct tags ("json") when matching
// source and destination fields.
//
// Example:
//
//	type Source struct {
//	    Name string `json:"full_name"`
//	}
//	mapper.Copy(&dst, src, mapper.WithJSONTag(true))
func WithJSONTag(use bool) Option {
	return func(c *Config) {
		c.UseJSONTag = use
	}
}

// WithCustomConverter registers a custom conversion function for a given type.
// The converter is used when mapping a value of that specific type.
//
// Example:
//
//	timeConverter := func(v reflect.Value) (reflect.Value, error) {
//	    if t, ok := v.Interface().(time.Time); ok {
//	        return reflect.ValueOf(t.Format(time.RFC3339)), nil
//	    }
//	    return v, nil
//	}
//
//	mapper.Copy(&dst, src,
//	    mapper.WithCustomConverter(reflect.TypeOf(time.Time{}), timeConverter))
func WithCustomConverter(typ reflect.Type, converter ConverterFunc) Option {
	return func(c *Config) {
		if c.CustomConverters == nil {
			c.CustomConverters = make(map[reflect.Type]ConverterFunc)
		}
		c.CustomConverters[typ] = converter
	}
}

// WithFieldNameMapper sets a custom function for transforming field names
// before matching. This is useful for converting between different naming
// conventions such as snake_case, camelCase, etc.
//
// Example:
//
//	mapper.Copy(&dst, src,
//	    mapper.WithFieldNameMapper(func(name string) string {
//	        return strings.ToLower(name)
//	    }))
func WithFieldNameMapper(mapper FieldNameMapperFunc) Option {
	return func(c *Config) {
		c.FieldNameMapper = mapper
	}
}

// WithErrorHandler registers a custom error handler that is invoked whenever
// a field mapping operation encounters an error. Returning nil continues
// the mapping process; returning an error stops it.
//
// Example:
//
//	mapper.Copy(&dst, src,
//	    mapper.WithErrorHandler(func(err error, srcField, dstField string) error {
//	        log.Printf("Mapping error: %v (src=%s, dst=%s)", err, srcField, dstField)
//	        return nil // Continue mapping
//	    }))
func WithErrorHandler(handler ErrorHandlerFunc) Option {
	return func(c *Config) {
		c.ErrorHandler = handler
	}
}

// WithSkipCircularCheck disables circular reference detection.
//
// ⚠️ Use with caution: only disable this if you are certain that
// your source data contains no circular references.
//
// Example:
//
//	mapper.Copy(&dst, src, mapper.WithSkipCircularCheck(true))
func WithSkipCircularCheck(skip bool) Option {
	return func(c *Config) {
		c.SkipCircularCheck = skip
	}
}

// WithTimeLayout specifies a custom time format for serializing or parsing
// time.Time values during mapping.
//
// Example:
//
//	mapper.Copy(&dst, src, mapper.WithTimeLayout("2006-01-02"))
func WithTimeLayout(layout string) Option {
	return func(c *Config) {
		c.TimeLayout = layout
	}
}

// WithMaxSliceCapacity defines an upper limit for slice allocation during mapping.
// This prevents excessive memory usage when mapping large slices.
//
// Example:
//
//	mapper.Copy(&dst, src, mapper.WithMaxSliceCapacity(10000))
func WithMaxSliceCapacity(capacity int) Option {
	return func(c *Config) {
		c.MaxSliceCapacity = capacity
	}
}

// WithAllowPrivateFields enables mapping of unexported (private) struct fields.
// ⚠️ This should be used cautiously, as it breaks Go's encapsulation guarantees.
//
// Example:
//
//	mapper.Copy(&dst, src, mapper.WithAllowPrivateFields(true))
func WithAllowPrivateFields(allow bool) Option {
	return func(c *Config) {
		c.AllowPrivateFields = allow
	}
}
