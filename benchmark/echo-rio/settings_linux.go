//go:build linux

package echorio

import "github.com/brickingsoft/rio"

func setting() {
	rio.UseCPUAffinity(true)
	rio.UseVortexNum(1)
}
