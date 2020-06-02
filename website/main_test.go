package main

import (
	"os"
	"testing"
)

func Test_cli(t *testing.T) {
	c := &cli{
		out: os.TempDir(),
	}
	if err := c.run(); err != nil {
		t.Fatal(err)
	}
}
