package domain

import "errors"

// ErrNotFound signals repository lookups that miss.
var ErrNotFound = errors.New("entity not found")
