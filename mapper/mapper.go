// Package mapper provides a high-performance, thread-safe utility for
// mapping one Go struct to another using reflection. It supports deep
// copying, custom type converters, circular reference detection, and
// field name mapping via struct tags.
//
// The Mapper is optimized for performance using object pooling and
// minimal allocations, and is safe for concurrent use.
//
// Key features:
//   - Thread-safe through sync.Pool
//   - Deep copy with circular reference detection
//   - Custom converters for specific types
//   - Field mapping via struct tags
//   - Case-sensitive or case-insensitive field matching
//   - Configurable mapping depth and nil handling
package mapper

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/fbarikzehi/gomap/internal/reflectutil"
)

// Mapper provides the main entry point for struct-to-struct mapping.
// It holds configuration options and manages a pool of reusable
// mapping contexts to minimize allocations.
type Mapper struct {
	config *Config    // Configuration for this mapper instance
	pool   *sync.Pool // Pool of reusable mapping contexts
}

// NewMapper creates and returns a new Mapper instance configured with
// the provided options. It validates configuration and initializes
// internal object pools for efficient reuse.
//
// Example:
//
//	mapper := NewMapper(
//	    WithMaxDepth(10),
//	    WithIgnoreUnexported(true),
//	    WithCustomConverter(timeType, timeConverter),
//	)
func NewMapper(opts ...Option) *Mapper {
	cfg := &Config{
		MaxDepth:          DefaultMaxDepth,
		IgnoreUnexported:  true,
		DeepCopy:          true,
		CaseSensitive:     true,
		CustomConverters:  make(map[reflect.Type]ConverterFunc),
		SkipCircularCheck: false,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return &Mapper{
		config: cfg,
		pool: &sync.Pool{
			New: func() interface{} {
				return &context{
					visited: make(map[uintptr]reflect.Value),
					errors:  make([]error, 0),
				}
			},
		},
	}
}

// Map performs the actual mapping from src to dst. The destination must
// be a pointer to a struct.
//
// Map performs a deep copy of all supported types (structs, slices, maps, etc.)
// and applies custom converters or tag-based field mapping as configured.
//
// Example:
//
//	var dst UserDTO
//	err := mapper.Map(&dst, srcUser)
//
// Returns an error if:
//   - dst or src is nil (ErrNilPointer)
//   - dst is not a pointer (ErrInvalidDestination)
//   - The mapping exceeds the maximum configured depth (ErrMaxDepthExceeded)
func (m *Mapper) Map(dst, src interface{}) error {
	if dst == nil || src == nil {
		return ErrNilPointer
	}

	dstVal := reflect.ValueOf(dst)
	if dstVal.Kind() != reflect.Ptr {
		return ErrInvalidDestination
	}

	srcVal := reflect.ValueOf(src)

	ctx := m.pool.Get().(*context)
	defer m.pool.Put(ctx)

	// Reset context before reuse
	for k := range ctx.visited {
		delete(ctx.visited, k)
	}
	ctx.errors = ctx.errors[:0]
	ctx.depth = 0
	ctx.config = m.config

	err := ctx.mapValue(dstVal.Elem(), srcVal)
	if err != nil {
		return err
	}

	if len(ctx.errors) > 0 {
		return fmt.Errorf("mapping completed with %d errors: %v", len(ctx.errors), ctx.errors[0])
	}

	return nil
}

// Copy is a convenience helper for performing a one-time struct mapping
// without explicitly creating a Mapper instance.
//
// Example:
//
//	var dst MyStruct
//	err := Copy(&dst, src, WithMaxDepth(5))
func Copy(dst, src interface{}, opts ...Option) error {
	m := NewMapper(opts...)
	return m.Map(dst, src)
}

// mapValue recursively maps a value from src to dst.
// It handles type routing, depth control, circular detection,
// and applies custom converters where applicable.
//
// Supported kinds:
//   - Structs
//   - Pointers
//   - Slices and arrays
//   - Maps
//   - Interfaces
//   - Basic types (numbers, strings, bools)
func (ctx *context) mapValue(dst, src reflect.Value) error {
	if !src.IsValid() {
		return nil
	}

	// Depth control
	if ctx.config.MaxDepth != NoDepthLimit && ctx.depth > ctx.config.MaxDepth {
		return ErrMaxDepthExceeded
	}

	// Handle nil source
	if reflectutil.IsNillable(src.Kind()) && src.IsNil() {
		if ctx.config.IgnoreNilFields {
			return nil
		}
		if dst.CanSet() && reflectutil.IsNillable(dst.Kind()) {
			dst.Set(reflect.Zero(dst.Type()))
		}
		return nil
	}

	// Circular reference detection
	if !ctx.config.SkipCircularCheck && reflectutil.IsPointerLike(src.Kind()) {
		if err := ctx.checkCircular(src); err != nil {
			return err
		}
	}

	// Custom converters
	if converter, ok := ctx.config.CustomConverters[src.Type()]; ok {
		converted, err := converter(src)
		if err != nil {
			return err
		}
		if dst.CanSet() && converted.Type().AssignableTo(dst.Type()) {
			dst.Set(converted)
		}
		return nil
	}

	ctx.depth++
	defer func() { ctx.depth-- }()

	switch src.Kind() {
	case reflect.Pointer:
		return ctx.mapPointer(dst, src)
	case reflect.Struct:
		return ctx.mapStruct(dst, src)
	case reflect.Map:
		return ctx.mapMap(dst, src)
	case reflect.Slice, reflect.Array:
		return ctx.mapSlice(dst, src)
	case reflect.Interface:
		return ctx.mapInterface(dst, src)
	default:
		return ctx.mapBasic(dst, src)
	}
}

// mapPointer handles mapping of pointer types by dereferencing and
// allocating destination pointers when necessary.
func (ctx *context) mapPointer(dst, src reflect.Value) error {
	if src.IsNil() {
		if ctx.config.IgnoreNilFields {
			return nil
		}
		if dst.CanSet() {
			dst.Set(reflect.Zero(dst.Type()))
		}
		return nil
	}

	srcElem := src.Elem()

	if dst.Kind() == reflect.Ptr {
		if dst.IsNil() && dst.CanSet() {
			dst.Set(reflect.New(dst.Type().Elem()))
		}
		return ctx.mapValue(dst.Elem(), srcElem)
	}

	return ctx.mapValue(dst, srcElem)
}

// mapStruct maps fields between two struct types. It respects configuration
// for field tags, case sensitivity, and unexported field handling.
func (ctx *context) mapStruct(dst, src reflect.Value) error {
	if dst.Kind() == reflect.Ptr {
		if dst.IsNil() && dst.CanSet() {
			dst.Set(reflect.New(dst.Type().Elem()))
		}
		return ctx.mapStruct(dst.Elem(), src)
	}

	if dst.Kind() != reflect.Struct {
		return nil
	}

	// Special case for time.Time
	if src.Type() == reflect.TypeOf(time.Time{}) {
		if dst.Type() == src.Type() && dst.CanSet() {
			dst.Set(src)
		}
		return nil
	}

	srcType := src.Type()
	dstType := dst.Type()

	for i := 0; i < src.NumField(); i++ {
		srcField := srcType.Field(i)

		// Skip unexported fields if configured
		if ctx.config.IgnoreUnexported && srcField.PkgPath != "" && !srcField.Anonymous {
			continue
		}

		// Tag filtering
		if ctx.config.TagName != "" {
			tag := srcField.Tag.Get(ctx.config.TagName)
			if tag == "" || tag == "-" {
				continue
			}
		}

		srcValue := src.Field(i)
		dstFieldName := ctx.getDestFieldName(srcField)
		dstField, found := ctx.findDstField(dstType, dstFieldName)
		if !found {
			continue
		}

		dstValue := dst.FieldByIndex(dstField.Index)
		if !dstValue.CanSet() {
			continue
		}

		// Zero field if configured
		if ctx.config.ZeroFields && srcValue.IsZero() {
			dstValue.Set(reflect.Zero(dstValue.Type()))
			continue
		}

		// Recursive field mapping
		if err := ctx.mapValue(dstValue, srcValue); err != nil {
			if ctx.config.ErrorHandler != nil {
				err = ctx.config.ErrorHandler(err, srcField.Name, dstField.Name)
			}
			if err != nil {
				ctx.addError(err)
			}
		}
	}

	return nil
}

// mapMap performs mapping between two maps, recursively mapping both keys
// and values. It creates a new destination map if needed.
func (ctx *context) mapMap(dst, src reflect.Value) error {
	if src.Kind() != reflect.Map || dst.Kind() != reflect.Map {
		return nil
	}

	if dst.IsNil() && dst.CanSet() {
		dst.Set(reflect.MakeMap(dst.Type()))
	}

	iter := src.MapRange()
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()

		newKey := reflect.New(dst.Type().Key()).Elem()
		newVal := reflect.New(dst.Type().Elem()).Elem()

		if err := ctx.mapValue(newKey, key); err != nil {
			ctx.addError(err)
			continue
		}
		if err := ctx.mapValue(newVal, value); err != nil {
			ctx.addError(err)
			continue
		}

		dst.SetMapIndex(newKey, newVal)
	}

	return nil
}

// mapSlice maps elements between slices and arrays. It allocates a
// new destination slice if necessary and maps elements recursively.
func (ctx *context) mapSlice(dst, src reflect.Value) error {
	if dst.Kind() != reflect.Slice && dst.Kind() != reflect.Array {
		return nil
	}

	srcLen := src.Len()

	if dst.Kind() == reflect.Slice {
		if dst.IsNil() || dst.Len() < srcLen {
			if dst.CanSet() {
				dst.Set(reflect.MakeSlice(dst.Type(), srcLen, srcLen))
			}
		}
	}

	length := min(dst.Len(), srcLen)
	for i := 0; i < length; i++ {
		if err := ctx.mapValue(dst.Index(i), src.Index(i)); err != nil {
			ctx.addError(fmt.Errorf("slice index %d: %w", i, err))
		}
	}

	return nil
}

// mapInterface handles mapping between interface values, extracting
// and mapping the underlying concrete types.
func (ctx *context) mapInterface(dst, src reflect.Value) error {
	if src.Kind() != reflect.Interface {
		return nil
	}

	if dst.Kind() != reflect.Interface {
		return ctx.mapValue(dst, src.Elem())
	}

	srcElem := src.Elem()
	if !srcElem.IsValid() {
		return nil
	}

	newDst := reflect.New(srcElem.Type()).Elem()
	if err := ctx.mapValue(newDst, srcElem); err != nil {
		return err
	}

	if dst.CanSet() {
		dst.Set(newDst)
	}

	return nil
}

// mapBasic handles assignment and conversion between basic types
// (e.g., numbers, strings, booleans). It performs direct assignment
// if the types match, otherwise attempts conversion if allowed.
func (ctx *context) mapBasic(dst, src reflect.Value) error {
	if !dst.CanSet() {
		return nil
	}

	if src.Type() == dst.Type() {
		dst.Set(src)
		return nil
	}

	if src.Type().ConvertibleTo(dst.Type()) {
		dst.Set(src.Convert(dst.Type()))
		return nil
	}

	return nil
}

// getDestFieldName determines the destination field name using
// struct tags, configuration options, or a custom field name mapper.
func (ctx *context) getDestFieldName(srcField reflect.StructField) string {
	if ctx.config.TagName != "" {
		if tag := srcField.Tag.Get(ctx.config.TagName); tag != "" && tag != "-" {
			return tag
		}
	}

	if ctx.config.UseJSONTag {
		if tag := srcField.Tag.Get("json"); tag != "" && tag != "-" {
			return tag
		}
	}

	if ctx.config.FieldNameMapper != nil {
		return ctx.config.FieldNameMapper(srcField.Name)
	}

	return srcField.Name
}

// findDstField locates the destination field in the target struct
// using case-sensitive or case-insensitive matching according to configuration.
func (ctx *context) findDstField(dstType reflect.Type, fieldName string) (reflect.StructField, bool) {
	if field, found := dstType.FieldByName(fieldName); found {
		return field, true
	}

	if !ctx.config.CaseSensitive {
		for i := 0; i < dstType.NumField(); i++ {
			field := dstType.Field(i)
			if reflectutil.EqualFold(field.Name, fieldName) {
				return field, true
			}
		}
	}

	return reflect.StructField{}, false
}
