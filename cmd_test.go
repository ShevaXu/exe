package exe_test

import (
	"context"
	"io/ioutil"
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
}

func TestStd(t *testing.T) {
	logF, err := ioutil.TempFile("", "log")
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
