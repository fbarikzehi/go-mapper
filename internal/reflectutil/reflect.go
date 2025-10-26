package reflectutil

import "reflect"

// IsNillable determines if a reflect.Kind can be nil.
// Used for safe nil checking during mapping operations.
func IsNillable(k reflect.Kind) bool {
	switch k {
	case reflect.Chan, reflect.Func, reflect.Interface,
		reflect.Map, reflect.Pointer, reflect.Slice:
		return true
	}
	return false
}

// IsPointerLike checks if a kind is a pointer-like type that can create
// circular references. This includes pointers, maps, slices, etc.
func IsPointerLike(k reflect.Kind) bool {
	switch k {
	case reflect.Pointer, reflect.Map, reflect.Slice,
		reflect.Chan, reflect.Func:
		return true
	}
	return false
}

// IsBasicType determines if a kind represents a basic Go type
// that can be directly assigned or needs simple type conversion.
func IsBasicType(k reflect.Kind) bool {
	switch k {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128,
		reflect.String:
		return true
	}
	return false
}

// EqualFold performs case-insensitive string comparison.
// Used for case-insensitive field name matching.
func EqualFold(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if ToLower(a[i]) != ToLower(b[i]) {
			return false
		}
	}
	return true
}

// ToLower converts a byte to lowercase
func ToLower(c byte) byte {
	if 'A' <= c && c <= 'Z' {
		return c + ('a' - 'A')
	}
	return c
}

// Min returns the minimum of two integers
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Max returns the maximum of two integers
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// GetFieldTag extracts tag value from struct field
func GetFieldTag(field reflect.StructField, tagName string) (string, bool) {
	tag := field.Tag.Get(tagName)
	if tag == "" || tag == "-" {
		return "", false
	}
	return tag, true
}

// IsZeroValue checks if a value is zero
func IsZeroValue(v reflect.Value) bool {
	if !v.IsValid() {
		return true
	}

	switch v.Kind() {
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Complex64, reflect.Complex128:
		return v.Complex() == 0
	case reflect.String:
		return v.String() == ""
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map,
		reflect.Ptr, reflect.Slice:
		return v.IsNil()
	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			if !IsZeroValue(v.Index(i)) {
				return false
			}
		}
		return true
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if !IsZeroValue(v.Field(i)) {
				return false
			}
		}
		return true
	}

	return false
}
