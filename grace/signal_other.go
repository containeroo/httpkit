//go:build !unix

package grace

import "os"

func defaultShutdownSignals() []os.Signal {
	return []os.Signal{os.Interrupt}
}
