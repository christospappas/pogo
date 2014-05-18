package pogo

import (
	"fmt"
)

// The version information.
const (
	VerMajor = 0
	VerMinor = 1
	VerPatch = 0
)

// Version returns current version of the package.
func Version() string {
	return fmt.Sprintf("%d.%d.%d", VerMajor, VerMinor, VerPatch)
}
