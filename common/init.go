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
	dirlist := [...]string{mydockerPathRoot, mydockerLibRoot, GetMyDockerPath(), GetImagePath(), GetTmpPath(), GetContainerPath()}
	for _, path := range dirlist {
		if err := os.MkdirAll((path), 0755); err != nil {
			fmt.Printf("create %v error %v\n", path, err)
			return fmt.Errorf("Mkdir(%v) failed: %v", path, err)
		}
	}
	return nil
}

func RuningCliApp(app *cli.App, args []string) {
	if err := app.Run(args); err != nil {
		log.Fatalf("The cli.App Run %v error\n", args)
	}
}
