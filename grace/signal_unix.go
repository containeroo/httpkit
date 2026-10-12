//go:build unix

package grace

import (
	"os"
	"syscall"
)

func defaultShutdownSignals() []os.Signal {
	return []os.Signal{os.Interrupt, syscall.SIGTERM}
}
