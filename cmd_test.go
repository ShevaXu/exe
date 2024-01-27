package exe_test

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/practigo/exe"
)

func TestCmd(t *testing.T) {
	err := exe.Cmd("ls").Exec(context.TODO(), exe.Pre(exe.Pipe))
	if err != nil {
		t.Error(err)
	}

	// cannot recognize pipe `|`
	t.Log(exe.Cmd("ps -ef | grep go").Exec(context.TODO(), exe.Pre(exe.Pipe)))

	if exe.Cmd("").Exec(context.TODO()) != exe.ErrInvalidCmd {
		t.Error(err, "should be invalid cmd")
	}
}

func TestStd(t *testing.T) {
	logF, err := os.CreateTemp("", "log")
	if err != nil {
		t.Fatal(err)
	}
	err = exe.Std.Run(context.TODO(), "ls", logF.Name())
	if err != nil {
		t.Fatal(err)
	}
	os.Remove(logF.Name())

	stdout, _, err := exe.Std.Probe(context.TODO(), "pwd")
	if err != nil {
		t.Fatal(err)
	}
	pwd := strings.TrimSpace(string(stdout))
	wd, _ := os.Getwd()
	if pwd != wd {
		t.Error(pwd, "!=", wd)
	}
}

func TestStdErr(t *testing.T) {
	_, _, err := exe.Std.Probe(context.TODO(), "pwd2")
	if err == nil {
		t.Fatal("should have error on pwd2")
	}
	t.Log(err)

	_, stderr, err := exe.Std.Probe(context.TODO(), "pwd -c")
	if err == nil {
		t.Fatal("should have error on pwd -c")
	}
	if len(stderr) <= 0 {
		t.Fatal("should have stderr")
	}
	t.Log(err)
	t.Log(string(stderr))
}

func TestDummy(t *testing.T) {
	p := exe.DummyProber{
		Stdout: []byte("foo"),
		Stderr: []byte("bar"),
	}
	retout, reterr, err := p.Probe(context.TODO(), "any-cmd")
	if string(retout) != "foo" || string(reterr) != "bar" || err != nil {
		t.Error("Dummy prober should return the exact outputs")
	}

	r := exe.DummyRunner{Err: exe.ErrInvalidCmd}
	if r.Run(context.TODO(), "any-cmd", "logfile") != exe.ErrInvalidCmd {
		t.Error("Dummy runner should return the exact error")
	}
}
