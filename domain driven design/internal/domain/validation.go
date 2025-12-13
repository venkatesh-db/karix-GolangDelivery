package domain

import (
	"fmt"
	"sort"
	"strings"
)

// ValidationError aggregates field level validation failures.
type ValidationError struct {
	fields map[string]string
}

// NewValidationError constructs an empty validation error accumulator.
func NewValidationError() *ValidationError {
	return &ValidationError{fields: make(map[string]string)}
}

// Add records a field level validation issue.
func (v *ValidationError) Add(field, message string) {
	v.fields[field] = message
}

// HasErrors indicates whether validation failures were recorded.
func (v *ValidationError) HasErrors() bool {
	return len(v.fields) > 0
}

// Fields returns a copy of the validation map for safe consumption.
func (v *ValidationError) Fields() map[string]string {
	out := make(map[string]string, len(v.fields))
	for k, val := range v.fields {
		out[k] = val
	}
	return out
}

// Error implements the error interface.
func (v *ValidationError) Error() string {
	if len(v.fields) == 0 {
		return "validation error"
	}
	parts := make([]string, 0, len(v.fields))
	keys := make([]string, 0, len(v.fields))
	for k := range v.fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s: %s", k, v.fields[k]))
	}
	return strings.Join(parts, "; ")
}
