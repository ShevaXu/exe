package exe_test

import (
	"context"
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
