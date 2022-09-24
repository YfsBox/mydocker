package container

import (
	"fmt"
	"math/rand"
	cm "mydocker/common"
	"os"
	"os/exec"
	"syscall"
)

func GetContainerId() string {
	randBytes := make([]byte, 6)
	rand.Read(randBytes)
	return fmt.Sprintf("%02x%02x%02x%02x%02x%02x",
		randBytes[0], randBytes[1], randBytes[2],
		randBytes[3], randBytes[4], randBytes[5])
}

func GetCloneContainerProc(runcmd string, cmdargs []string) *exec.Cmd {

	cm.DPrintf("The runcmd is %v,the len of cmdargs is %v\n", runcmd, len(cmdargs))

	for _, arg := range cmdargs {
		cm.DPrintf("arg: %v", arg)
	}
	cmd := exec.Command(runcmd, cmdargs[0:]...)

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWIPC | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET,
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd

}
