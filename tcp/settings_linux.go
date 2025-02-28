//go:build linux

package main

import "github.com/brickingsoft/rio"

func setting() {
	rio.UseZeroCopy(true)
}
