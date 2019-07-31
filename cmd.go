/*
Package exe provides context control and output redirections
above exec.Cmd.
*/
package exe

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"strings"
)

// A Hook provides access to the underlying Cmd.
// The returned error would indicate different
// behaviour at various process stages.
type Hook func(cmd *exec.Cmd) error

func kill(cmd *exec.Cmd) error {
	return cmd.Process.Kill()
}

// Pipe ...
func Pipe(cmd *exec.Cmd) error {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return nil
}

type cmdHooks struct {
	pre, post, exit Hook
}

// Option ...
type Option func(*cmdHooks)

// Cmd ...
type Cmd string

// const errors
var (
	ErrInvalidCmd = errors.New("invalid command")
)

// Exec ...
func (c Cmd) Exec(ctx context.Context, options ...Option) error {
	args := strings.Fields(string(c))
	if len(args) < 1 {
		return ErrInvalidCmd
	}

	// look for binary path
	path, err := exec.LookPath(args[0])
	if err != nil {
		return err
	}

	cmd := exec.Command(path, args[1:]...)

	hooks := &cmdHooks{
		exit: kill,
	}
	for _, o := range options {
		o(hooks)
	}

	if hooks.pre != nil {
		if err = hooks.pre(cmd); err != nil {
			return err
		}
	}

	if err = cmd.Start(); err != nil {
		return err
	}

	if hooks.post != nil {
		hooks.post(cmd)
	}

	// controls
	done := ctx.Done()
	cleanup := make(chan struct{})

	// exit handling
	go func() {
		select {
		case <-done:
			if hooks.exit != nil {
				hooks.exit(cmd)
			}
		case <-cleanup:
			return
		}
	}()

	err = cmd.Wait()

	// cleanup the exit handling goroutine
	close(cleanup)

	return err
}

// Pre provides a hook that runs before the cmd starts.
// A non-nil error returned would prevent starting the cmd.
func Pre(h Hook) Option {
	return func(hs *cmdHooks) {
		hs.pre = h
	}
}

// Post provides a hook that runs after the
// cmd starts, but ignores the returned error.
// The runner waits for the cmd's exit after this.
func Post(h Hook) Option {
	return func(hs *cmdHooks) {
		hs.post = h
	}
}

// Done replaces the default hook (Kill) that kills (-9)
// the process when the context cancelled
// (typically sending another signals that the binary can
// handle as graceful exit).
func Done(h Hook) Option {
	return func(hs *cmdHooks) {
		hs.exit = h
	}
}
