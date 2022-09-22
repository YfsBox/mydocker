package container

import (
	"fmt"
	"syscall"
	"os"
	"os/exec"
	cm "mydocker/common"
)

//该文件下,主要是包含mount和namespace相关的隔离

func GetCloneContainerProc(runcmd string,cmdargs []string) *exec.Cmd {

	cm.DPrintf("The runcmd is %v,the len of cmdargs is %v\n",runcmd,len(cmdargs))

	for _,arg := range cmdargs {
		cm.DPrintf("arg: %v",arg)
	}
	cmd := exec.Command(runcmd,cmdargs[0:]...)


	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS |  syscall.CLONE_NEWPID| syscall.CLONE_NEWIPC | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET,
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd

}

func SetUpMount() error {
    // systemd 加入linux之后, mount namespace 就变成 shared by default, 所以你必须显示
    //声明你要这个新的mount namespace独立。
    err := syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
    if err != nil {
        return err
    }
    //mount proc
    defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
    err = syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
    if err != nil {
        return fmt.Errorf("mount proc error: %v",err)
    }

    return nil
}

func unmount(file string) error {
	err := syscall.Unmount(file,0);
	return err
}


func UnmountAll() error {
	cm.DPrintf("umount proc\n")
	if err := unmount("/proc"); err != nil {
		return fmt.Errorf("Unmount proc error %v",err)
	}
	return nil
}


