package k6_test

import (
	"github.com/brickingsoft/rio_examples/benchmark/cli/k6"
	"os"
	"testing"
)

func TestScript(t *testing.T) {
	file := k6.ScriptFile("test", "127.0.0.1", 9000, 100, "10s")

	b, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(b))

	os.Remove(file)
}
