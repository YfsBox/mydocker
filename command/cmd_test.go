package command

import (
	cm "mydocker/common"
	"os"
	"testing"
)

func TestRun(t *testing.T) {
	app := InitCliApp()
	cm.RuningCliApp(app, os.Args)

}
