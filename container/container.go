package container

import (
	cm "mydocker/common"
	"os"
	"os/exec"
	"syscall"
)

func GetContainerId() string {
	return cm.RandomString(12)
}

func GetCloneContainerProc(runcmd string, cmdargs []string) *exec.Cmd {

	//cm.DPrintf("The runcmd is %v,the len of cmdargs is %v\n", runcmd, len(cmdargs))
	cmd := exec.Command(runcmd, cmdargs[0:]...)

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWIPC | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET,
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd

}
