package exe

import "context"

// DummyProber returns the stdout/stderr and error
// for testing purpose.
type DummyProber struct {
	Stdout []byte
	Stderr []byte
	Err    error
}

// Probe implements Prober.
func (p *DummyProber) Probe(_ context.Context, cmd string) (stdout, stderr []byte, err error) {
	return p.Stdout, p.Stderr, p.Err
}

// DummyRunner returns the error provided
// for testing purpose.
type DummyRunner struct {
	Err error
}

// Run implements Runner.
func (r *DummyRunner) Run(_ context.Context, cmd, logFile string) error {
	return r.Err
}
