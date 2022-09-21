package common

import (
	"fmt"
	"log"
	"os"
	//cmd "mydocker/command"
	"github.com/urfave/cli"
)


func CheckRootUser() {
	if os.Geteuid() != 0 {
		log.Fatalf("The user is not root\n")
	}
}

func InitMyDockerDirs() error {
	mydockerPath := GetMyDockerPath()
	if err := os.MkdirAll(mydockerPath,0755);err != nil {
		return fmt.Errorf("Mkdir(%v) failed: %v",mydockerPath,err)
	}
	return nil
}


func RuningCliApp(app *cli.App,args []string ) {
	if err := app.Run(args) ; err != nil {
		log.Fatalf("The cli.App Run %v error\n",args)
	}
}

