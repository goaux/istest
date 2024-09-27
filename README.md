# istest

`istest` is a Go package that provides utilities for detecting and working with test environments in Go programs.

Use with care as excessive use can ruin your tests.

The determination is made based on whether `os.Args[0]` ends with ".test".

The determination of which of the four test functions, namely TestXxx, ExampleXxx, BenchmarkXxx and FuzzXxx, is based on the function name in the call stack and the file name ending with "_test.go".

Alternatively, you can set it explicitly in the context instead of determining it automatically.

## Features

- Detect if the current execution is within a test environment
- Filter for specific types of tests (unit tests, benchmarks, examples, fuzz tests)
- Store and retrieve test function types in a context

## Usage

Here are some examples of how to use the `istest` package:

```go
import "github.com/goaux/istest"

// Check if running in any test environment
if istest.Is() {
    // Execute test-specific code
}

// Check if running a specific type of test
if istest.Is(istest.FuncTest) {
    // Execute code specific to unit tests
}

// Check for multiple test types
if istest.Is(istest.FuncBenchmark, istest.FuncFuzz) {
    // Execute code specific to benchmarks or fuzz tests
}

// Use with context
ctx := istest.Context(context.Background(), istest.FuncTest)
if istest.IsContext(ctx, istest.FuncTest) {
    // Execute test-specific code
}
```
