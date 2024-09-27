// Package istest provides utilities for detecting and working with test environments in Go programs.

// Use with care as excessive use can ruin your tests.
//
// The determination is made based on whether `os.Args[0]` ends with ".test".
//
// The determination of which of the four test functions, namely TestXxx, ExampleXxx, BenchmarkXxx and FuzzXxx, is based on the function name in the call stack and the file name ending with "_test.go".
//
// Alternatively, you can set it explicitly in the context instead of determining it automatically.
//
// Usage
//
//	import "github.com/goaux/istest"
//
//	// Check if running in any test environment
//	if istest.Is() {
//	    // Execute test-specific code
//	}
//
//	// Check if running a specific type of test
//	if istest.Is(istest.FuncTest) {
//	    // Execute code specific to unit tests
//	}
//
//	// Check for multiple test types
//	if istest.Is(istest.FuncBenchmark, istest.FuncFuzz) {
//	    // Execute code specific to benchmarks or fuzz tests
//	}
//
//	// Use with context
//	ctx := istest.Context(context.Background(), istest.FuncTest)
//	if istest.IsContext(ctx, istest.FuncTest) {
//	    // Execute test-specific code
//	}
package istest

import (
	"context"
	"os"
	"runtime"
	"strings"

	"github.com/goaux/funcname"
)

// function is a custom string type representing a test function type.
type function string

// Match returns true if the name starts with the function prefix, false otherwise.
func (f function) Match(name string) bool { return strings.HasPrefix(name, string(f)) }

// String returns a string representation of the function type.
func (f function) String() string { return "Func" + string(f) }

var (
	FuncBenchmark = function("Benchmark") // FuncBenchmark represents a benchmark function.
	FuncFuzz      = function("Fuzz")      // FuncFuzz represents a fuzz test function.
	FuncExample   = function("Example")   // FuncExample represents an example function.
	FuncTest      = function("Test")      // FuncTest represents a test function.
)

// Functions returns a slice of specified function types.
func Functions(functions ...function) []function {
	return functions
}

// isTest is true if os.Args[0] has the suffix ".test".
var isTest = len(os.Args) > 0 && strings.HasSuffix(os.Args[0], ".test")

// Is determines if the current execution context is within a test environment.
//
// It returns true if any of the following conditions are met:
// 1. The program is running as a test (os.Args[0] ends with ".test") and no specific test functions are specified.
// 2. The program is running as a test and there's a function in the call stack that:
//   - Has a name prefixed with any of the provided function types
//   - Is defined in a file with a "_test.go" suffix
//
// Parameters:
//   - functions: Optional variadic parameter specifying the types of tests to check for.
//     Valid values are FuncTest, FuncExample, FuncBenchmark, and FuncFuzz.
//
// Returns:
//   - bool: True if the execution is within a test environment matching the specified criteria, false otherwise.
func Is(functions ...function) bool {
	if !isTest {
		return false
	}
	if len(functions) == 0 {
		return true
	}
	pc := make([]uintptr, 16)
	n := runtime.Callers(2, pc)
	if n == 0 {
		return false
	}
	iter := runtime.CallersFrames(pc[:n])
	for {
		frame, more := iter.Next()
		if strings.HasSuffix(frame.File, "_test.go") {
			_, name := funcname.Split(frame.Function)
			for _, f := range functions {
				if f.Match(name) {
					return true
				}
			}
		}
		if !more {
			break
		}
	}
	return false
}

type key struct{}

// Context creates a new context with the specified function type.
//
// Parameters:
//   - parent: The parent context.
//   - fn: The function type to store in the context.
//
// Returns:
//   - context.Context: A new context containing the specified function type.
//
// See [IsContext].
func Context(parent context.Context, fn function) context.Context {
	return context.WithValue(parent, key{}, fn)
}

// IsContext checks if the context contains any of the specified function types.
//
// Parameters:
//   - ctx: The context to check.
//   - functions: Optional variadic parameter specifying the function types to check for.
//
// Returns:
//   - bool: True if the context contains any of the specified function types (or any function type if none specified), false otherwise.
//
// See [Context].
func IsContext(ctx context.Context, functions ...function) bool {
	fn, ok := ctx.Value(key{}).(function)
	if !ok {
		return false
	}
	if len(functions) == 0 {
		return true
	}
	for _, f := range functions {
		if f == fn {
			return true
		}
	}
	return false
}

// Background returns the result of calling [Context]([context.Background](), fn).
//
// See [TODO].
func Background(fn function) context.Context {
	return Context(context.Background(), fn)
}

// TODO returns the result of calling [Context]([context.TODO](), fn).
//
// See [Background].
func TODO(fn function) context.Context {
	return Context(context.TODO(), fn)
}
