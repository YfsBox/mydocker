package command

import (
	"fmt"
	"log"
	cm "mydocker/common"
	cnt "mydocker/container"
	"os"
	"os/exec"
	"syscall"

	"github.com/urfave/cli"
)

func checkRunExecArgslen(args *cli.Args) bool {
	return true
}


func RunInit() string { //返回的是id和error

	log.Printf("a container begin running,the pid is %v\n",cm.GetPidStr())
	containerId := cnt.GetContainerId() //得到一个新的Id
	cnt.CreateCgroupForContainer(containerId)

	return containerId
}

func RunExec(runcmd []string,containerId string) error {
	cm.DPrintf("print args\n")
	//其中应该有关于containerId的部分,暂且将第二个参数定为containerID
	containerID := containerId
	syscall.Sethostname([]byte(containerID)) 
	if err := cnt.CreateCgroupForContainer(containerID); err != nil {
		//
	}
	cm.DPrintf("the clone proc pid: %v\n",os.Getegid())
	//挂载proc
	if err := cnt.SetUpMount(); err != nil {
		return fmt.Errorf("SetUpMount error(%v) in RunExec for %v",err,containerID)
	}

	cmd := exec.Command(runcmd[0],runcmd[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("cmd.Run error(%v) in RunExec for %v",err,containerID)
	}

	return nil
}

//首先做好关于cgroup参数的修改
//然后想清楚结合clone下的cgroup是怎么进行设置的

