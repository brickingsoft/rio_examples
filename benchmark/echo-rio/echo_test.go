package echorio_test

import (
	"fmt"
	echorio "github.com/brickingsoft/rio_examples/benchmark/echo-rio"
	"github.com/brickingsoft/rio_examples/benchmark/metric"
	"testing"
)

func BenchmarkECHO(b *testing.B) {
	cost, actions, inbounds, outbounds, failures, err := echorio.Bench(1, 500, 9000, 1024)
	if err != nil {
		b.Errorf("ECHO-RIO benching failed: %v", err)
		return
	}
	b.Log(fmt.Sprintf("ECHO-RIO benching complete(%s): %d conn/sec, %s inbounds/sec, %s outbounds/sec, %d failures",
		cost.String(), actions, metric.FormatBytes(inbounds), metric.FormatBytes(outbounds), failures))
}
