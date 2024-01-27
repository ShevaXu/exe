package exe

import (
	"bytes"
	"context"
	"os"
	"os/exec"
)

// A Runner runs a command and pipes its output to a log file.
type Runner interface {
	// Run starts a process and waits for its exit.
	// The ctx is used to cancel the process while it is
	// still running.
	Run(ctx context.Context, cmd, logFile string) error
}

// A Prober runs a command and returns its output directly
// for inspection.
type Prober interface {
	// Probe starts a process and waits for its exit.
	// The ctx is used to cancel the process while it is
	// still running. The system stdout/stderr are return
	// separately.
	Probe(ctx context.Context, cmd string) (stdout, stderr []byte, err error)
}

type std struct{}

func (s *std) Run(ctx context.Context, cmd, logFile string) (err error) {
	logF, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return
	}
	defer logF.Close()

	c := Cmd(cmd)

	return c.Exec(ctx, Pre(func(cmd *exec.Cmd) error {
		// 2>&1
		cmd.Stdout = logF
		cmd.Stderr = logF
		return nil
	}))
}

func (s *std) Probe(ctx context.Context, cmd string) (stdout, stderr []byte, err error) {
	c := Cmd(cmd)
	var o, e bytes.Buffer

	err = c.Exec(ctx, Pre(func(cmd *exec.Cmd) error {
		cmd.Stdout = &o
		cmd.Stderr = &e
		return nil
	}))

	stdout = o.Bytes()
	stderr = e.Bytes()

	return
}

// Std provides a singleton.
var Std = &std{}
