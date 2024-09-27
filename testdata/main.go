package main

import (
	"fmt"
	"os"

	"github.com/goaux/istest"
)

func main() {
	if istest.Is() {
		fmt.Fprintln(os.Stderr, "error: istest.Is() returns true")
		os.Exit(1)
	}
	if istest.Is(
		istest.FuncTest,
		istest.FuncExample,
		istest.FuncBenchmark,
		istest.FuncFuzz,
	) {
		fmt.Fprintln(os.Stderr, "error: istest.Is(FuncTest,FuncExample,FuncBenchmark,FuncFuzz) returns true")
		os.Exit(1)
	}
	fmt.Println("ok")
}
