//go:build linux

package http_rio

import "github.com/brickingsoft/rio"

func setting() {
	rio.UseCPUAffinity(true)
	rio.UseVortexNum(1)
}
