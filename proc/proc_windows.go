package proc

import (
	"fmt"
	"os"
)

func Setuid(b bool) {
}

func SetSignalExit(exitFunc func()) {
}

func SendExitSignal() {
	fmt.Fprintf(os.Stderr, "Not yet supported.\n")
}
