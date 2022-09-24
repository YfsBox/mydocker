package command

import (
	"fmt"
	"log"
	cm "mydocker/common"
	cnt "mydocker/container"
	img "mydocker/image"
	"os"
	"os/exec"
	"syscall"

	"github.com/urfave/cli"
)

func checkRunExecArgslen(args *cli.Args) bool {
	return true
}

func RunInit() string { //返回的是id和error

	log.Printf("a container begin running,the pid is %v\n", cm.GetPidStr())
	containerId := cnt.GetContainerId() //得到一个新的Id
	cnt.CreateCgroupForContainer(containerId)

	return containerId
}

func RunExec(runcmd []string, containerId string, limit *cnt.CgroupLimit) error {
	cm.DPrintf("print args\n")
	//其中应该有关于containerId的部分,暂且将第二个参数定为containerID
	containerID := containerId
	syscall.Sethostname([]byte(containerID))

	imagefsList, _ := img.ProcessLayers(containerID)

	if err := cnt.CreateAndMountFs(imagefsList, containerID); err != nil { //难道需要将read部分
		return fmt.Errorf("CreateContainerDirs error %v", err)
	}
	if err := cnt.CreateCgroupForContainer(containerID); err != nil {
		return fmt.Errorf("CreateCgroupForContainer %v err from RunExec", err)
	}
	if err := cnt.ConfigCgroupParameter(containerID, limit); err != nil {
		return fmt.Errorf("ConfigCgroupParameter %v err from RunExec", err)
	}

	cm.DPrintf("the clone proc pid: %v\n", os.Getegid())
	//挂载proc
	if err := cnt.SetUpMount(); err != nil {
		return fmt.Errorf("SetUpMount error(%v) in RunExec for %v", err, containerID)
	}

	cmd := exec.Command(runcmd[0], runcmd[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("cmd.Run error(%v) in RunExec for %v", err, containerID)
	}
	cm.DPrintf("Remove the container's files\n")
	if err := cnt.RemoveCgroupForContainer(containerID); err != nil {
		return fmt.Errorf("ConfigCgroupParameter %v err from RunExec", err)
	}

	return nil
}

//首先做好关于cgroup参数的修改
//然后想清楚结合clone下的cgroup是怎么进行设置的
