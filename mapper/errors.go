// Package mapper defines common error types and sentinel values
// used throughout the mapping process.
//
// These errors provide structured and context-rich diagnostics
// for debugging and error handling in mapping operations.
package mapper

import (
	"errors"
	"fmt"
)

// Common sentinel errors returned by Mapper operations.
//
// These can be compared using errors.Is() for fine-grained
// error handling. They are designed to describe specific mapping
// failure conditions.
var (
	// ErrNilPointer indicates that a nil pointer was provided
	// as either the source or destination in a mapping operation.
	ErrNilPointer = errors.New("mapper: nil pointer provided")

	// ErrUnsupportedType indicates that a value of an unsupported
	// or unhandled type was encountered during mapping.
	ErrUnsupportedType = errors.New("mapper: unsupported type")

	// ErrInvalidDestination indicates that the destination argument
	// must be a pointer type but was not.
	ErrInvalidDestination = errors.New("mapper: destination must be a pointer")

	// ErrTypeMismatch indicates a type incompatibility between
	// source and destination fields that prevents assignment or conversion.
	ErrTypeMismatch = errors.New("mapper: type mismatch between source and destination")

	// ErrMaxDepthExceeded indicates that the mapper exceeded the
	// maximum configured depth for nested struct traversal.
	ErrMaxDepthExceeded = errors.New("mapper: maximum depth exceeded")

	// ErrCircularReference indicates that a circular reference
	// was detected in the source object graph during deep copy.
	ErrCircularReference = errors.New("mapper: circular reference detected")
)

// MapError represents a detailed mapping failure, providing contextual
// information such as source and destination field names, types,
// operation depth, and the underlying error.
//
// MapError values are typically returned when a mapping operation
// encounters an error that cannot be handled internally.
//
// Example:
//
//	err := mapper.Map(&dst, src)
//	if e := new(mapper.MapError); errors.As(err, &e) {
//	    log.Printf("Failed to map field %s → %s: %v", e.SrcField, e.DstField, e.Err)
//	}
type MapError struct {
	// Err is the underlying cause of the failure.
	Err error

	// SrcField and DstField contain the names of the source and destination fields
	// involved in the mapping error, if available.
	SrcField string
	DstField string

	// SrcType and DstType represent the fully qualified types
	// of the source and destination values.
	SrcType string
	DstType string

	// Depth indicates the recursion level at which the error occurred.
	Depth int

	// Operation provides a short description of the failed mapping operation,
	// e.g., "mapStruct", "mapSlice", etc.
	Operation string
}

// Error implements the error interface and returns a formatted string
// describing the mapping failure in detail.
func (e *MapError) Error() string {
	if e.SrcField != "" && e.DstField != "" {
		return fmt.Sprintf(
			"mapper: failed to map %s.%s → %s.%s: %v",
			e.SrcType, e.SrcField, e.DstType, e.DstField, e.Err,
		)
	}
	return fmt.Sprintf("mapper: %s operation failed: %v", e.Operation, e.Err)
}

// Unwrap allows MapError to participate in Go's error unwrapping chain.
//
// Example:
//
//	if errors.Is(err, mapper.ErrTypeMismatch) { ... }
func (e *MapError) Unwrap() error {
	return e.Err
}

// Is enables comparison of MapError with target sentinel errors using errors.Is().
//
// This allows callers to check for underlying causes such as ErrTypeMismatch
// or ErrCircularReference directly.
func (e *MapError) Is(target error) bool {
	return errors.Is(e.Err, target)
}
