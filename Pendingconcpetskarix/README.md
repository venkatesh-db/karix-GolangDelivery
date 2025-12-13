# Go Generics Examples

This repository contains comprehensive examples of Go generics, including generic functions, type constraints, and advanced data structures.

## Files

### 1. `01_generic_functions.go`
Basic generic functions demonstrating:
- Generic print functions
- Pair functions with multiple type parameters
- Slice operations (First, Last, Reverse)
- Functional programming (Filter, Map, Reduce)
- Utility functions (Contains, Unique)

**Run:**
```bash
go run 01_generic_functions.go
```

### 2. `02_type_constraints.go`
Type constraints including:
- Built-in constraints (any, comparable)
- Custom constraints (Number)
- Using `constraints` package (Ordered, Integer, Float)
- Interface-based constraints
- Method constraints
- Combining constraints
- Approximate element constraints (~)

**Run:**
```bash
go run 02_type_constraints.go
```

### 3. `03_generic_data_structures.go`
Advanced generic data structures:
- Stack (LIFO)
- Queue (FIFO)
- Linked List
- Binary Search Tree
- Pair/Tuple
- Optional/Maybe type
- Result type (for error handling)
- Cache (key-value store)
- Set with union and intersection

**Run:**
```bash
go run 03_generic_data_structures.go
```

## Setup

1. Initialize the module and download dependencies:
```bash
go mod download
```

2. Run any example:
```bash
go run <filename>.go
```

## Key Concepts

### Type Parameters
```go
func Print[T any](value T) {
    fmt.Println(value)
}
```

### Type Constraints
```go
func Max[T constraints.Ordered](a, b T) T {
    if a > b {
        return a
    }
    return b
}
```

### Generic Types
```go
type Stack[T any] struct {
    items []T
}
```

### Custom Constraints
```go
type Number interface {
    int | int64 | float64
}

func Add[T Number](a, b T) T {
    return a + b
}
```

## Requirements
- Go 1.18 or higher (generics support)
- `golang.org/x/exp` package for constraints

## Benefits of Generics
- **Type Safety**: Compile-time type checking
- **Code Reusability**: Write once, use with multiple types
- **Performance**: No runtime overhead like interface{}
- **Cleaner Code**: Eliminate type assertions and reflection
