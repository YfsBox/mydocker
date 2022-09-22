package command

import (
	"testing"
	"os"
	cm "mydocker/common"
)

func TestRun(t *testing.T) {
	app := InitCliApp()
	cm.RuningCliApp(app,os.Args)

}