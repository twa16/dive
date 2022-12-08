package nerdctl

import (
	"fmt"
	"github.com/wagoodman/dive/utils"
	"io"
	"os"
	"os/exec"
)

// runNerdctlCmd runs a given Nerdctl command in the current tty
func runNerdctlCmd(cmdStr string, args ...string) error {
	if !isNerdctlClientBinaryAvailable() {
		return fmt.Errorf("cannot find nerdctl client executable")
	}

	allArgs := utils.CleanArgs(append([]string{cmdStr}, args...))

	cmd := exec.Command("nerdctl", allArgs...)
	cmd.Env = os.Environ()

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func streamNerdctlCmd(args ...string) (error, io.Reader) {
	if !isNerdctlClientBinaryAvailable() {
		return fmt.Errorf("cannot find nerdctl client executable"), nil
	}

	cmd := exec.Command("nerdctl", utils.CleanArgs(args)...)
	cmd.Env = os.Environ()

	reader, writer, err := os.Pipe()
	if err != nil {
		return err, nil
	}

	cmd.Stdout = writer
	cmd.Stderr = os.Stderr

	return cmd.Start(), reader
}

func isNerdctlClientBinaryAvailable() bool {
	_, err := exec.LookPath("nerdctl")
	return err == nil
}
