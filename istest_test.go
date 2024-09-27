package istest_test

import (
	"context"
	"os"
	"os/exec"
	"testing"

	"github.com/goaux/istest"
)

func TestIs(t *testing.T) {
	// This should always return true when running tests
	if !istest.Is() {
		t.Error("istest.Is() should return true when running tests")
	}

	// This should return true for unit tests
	if !istest.Is(istest.FuncTest) {
		t.Error("istest.Is(istest.FuncTest) should return true for unit tests")
	}

	// This should return false for benchmarks in this context
	if istest.Is(istest.FuncBenchmark) {
		t.Error("istest.Is(istest.FuncBenchmark) should return false in a unit test")
	}
}

func TestIsContext(t *testing.T) {
	ctx := context.Background()

	// Should return false for an unmodified context
	if istest.IsContext(ctx) {
		t.Error("istest.IsContext should return false for an unmodified context")
	}

	// Should return true for a context with FuncTest
	testCtx := istest.Context(ctx, istest.FuncTest)
	if !istest.IsContext(testCtx, istest.FuncTest) {
		t.Error("istest.IsContext should return true for a context with FuncTest")
	}

	// Should return false for a context with a different function type
	if istest.IsContext(testCtx, istest.FuncBenchmark) {
		t.Error("istest.IsContext should return false for a context with a different function type")
	}
}

func BenchmarkIs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if !istest.Is(istest.FuncBenchmark) {
			b.Fatal("istest.Is(istest.FuncBenchmark) returns false; want true")
		}
	}
}

func FuzzIs(f *testing.F) {
	f.Fuzz(func(t *testing.T, _ int) {
		if !istest.Is() {
			t.Fatal("istest.Is() returns false; want true")
		}
		if !istest.Is(istest.FuncFuzz) {
			t.Fatal("istest.Is(istest.FuncFuzz) returns false; want true")
		}
		list := istest.Functions(
			istest.FuncTest,
			istest.FuncExample,
			istest.FuncBenchmark,
		)
		for _, f := range list {
			if istest.Is(f) {
				t.Fatalf("istest.Is(istest.%v) returns true; want false", f)
			}
		}
	})
}

func TestIs_withSubprocess(t *testing.T) {
	cmd := exec.Command("go", "run", "./testdata/main.go")
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Errorf("Expected Is() to return false in non-test environment: %v", err)
	}
}

func TestBackground(t *testing.T) {
	ctx := istest.Background(istest.FuncTest)

	if !istest.IsContext(ctx) {
		t.Fatal("istest.IsContext(ctx) returns false; want true")
	}

	if !istest.IsContext(ctx, istest.FuncTest) {
		t.Fatal("istest.IsContext(ctx, istest.FuncTest) returns false; want true")
	}

	if istest.IsContext(ctx, istest.FuncFuzz) {
		t.Fatal("istest.IsContext(ctx, istest.FuncFuzz) returns true; want false")
	}
}

func TestTODO(t *testing.T) {
	ctx := istest.TODO(istest.FuncFuzz)

	if !istest.IsContext(ctx) {
		t.Fatal("istest.IsContext(ctx) returns false; want true")
	}

	if istest.IsContext(ctx, istest.FuncTest) {
		t.Fatal("istest.IsContext(ctx, istest.FuncTest) returns true; want false")
	}

	if !istest.IsContext(ctx, istest.FuncFuzz) {
		t.Fatal("istest.IsContext(ctx, istest.FuncFuzz) returns false; want true")
	}
}
