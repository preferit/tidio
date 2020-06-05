package perm

import (
	"fmt"
	"testing"
)

func ExampleOperation_String() {
	fmt.Println(OpRead, OpWrite, OpExec)
	// output: read write exec
}

func ExampleOperation_Short() {
	fmt.Printf("%s%s%s", OpRead.Short(), OpWrite.Short(), OpExec.Short())
	// output: rwx
}

func Test_unknown_operations(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		defer catchPanic(t)
		Operation(99).String()
	})
	t.Run("Short", func(t *testing.T) {
		defer catchPanic(t)
		Operation(99).Short()
	})
}

func catchPanic(t *testing.T) {
	t.Helper()
	if recover() == nil {
		t.Fatal("no panic")
	}
}
