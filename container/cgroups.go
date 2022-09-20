package container

import (
	"fmt"
	"io/ioutil"
	cm "mydocker/common"
	"os"
	"os/exec"
	//"syscall"
)

const (
	cgroupPath = "/sys/fs/cgroup"
	cpuCgroupPath = cgroupPath + "/cpu/mydocker"
	memoryCgroupPath = cgroupPath + "/memory/mydocker"
	pidsCgroupPath = cgroupPath + "/pids/mydocker"
	procsSubPath = "cgroup.procs" //记录受管理的进程pid
	norSubPath = "notify_on_release"  
)

var cgroupList = []string {
	cpuCgroupPath,
	memoryCgroupPath,
}

func isCgroupRootExist() bool { //判断cgroup根目录(mydocker)是否存在
	for _,path := range cgroupList {
		if !cm.IsFileExist(path) { //如果err != nil && IsNotExist说明不存在
			return false
		}
	}
	return true
}

func InitCgroupRootDirs() error { //创建各subsystem中关于mydocker的根目录
	for _,path := range cgroupList {
		if err := os.MkdirAll(path,0755); err != nil { //创建相关的文件夹
			return fmt.Errorf("Mkdir(%v) failed: %v",path,err)
		}
	}
	return nil
}

func RemoveCgroupRootDirs() error {
	for _,path := range cgroupList {
		if err := os.RemoveAll(path); err != nil { //创建相关的文件夹
			return fmt.Errorf("RemoveAll(%v) failed: %v",path,err)
		}
	}
	return nil
}

func CreateCgroupForContainer(containerId string) error { 
	if !isCgroupRootExist() {
		return fmt.Errorf("CreateCgroupForContainer failed,root not exist")
	}
	//在根目录下创建文件夹
	for _,rootpath := range cgroupList {

		containpath := cm.JoinPath(rootpath,containerId)
		if err := os.MkdirAll(containpath,0755); err != nil {
			return fmt.Errorf("Mkdir(%v) failed: %v from CreateCgroupForContainer %v",containpath,err,containerId)
		} //没有该Id对应的文件夹,就创建,有则什么也不干
		//在每个容器对应的文件夹之下,有两个关键的文件需要设置
		procsPath := cm.JoinPath(containpath,procsSubPath)
		norPath := cm.JoinPath(containpath,norSubPath)

		if err := ioutil.WriteFile(procsPath,[]byte(cm.GetPidStr()),0700); err != nil {
			return fmt.Errorf("WriteFile(%v) failed: %v from CreateCgroupForContainer %v",procsPath,err,containerId)
		}
		if err := ioutil.WriteFile(norPath,[]byte("1"),0700); err != nil {
			return fmt.Errorf("WriteFile(%v) failed: %v from CreateCgroupForContainer %v",norPath,err,containerId)
		}
	}
	return nil
}

func RemoveCgroupForContainer(containerId string) error { //目前关闭部分出了问题
	if !isCgroupRootExist() {
		return fmt.Errorf("CreateCgroupForContainer failed,root not exist")
	}

	for _,rootpath := range cgroupList {
		containpath := cm.JoinPath(rootpath,containerId)
		if !cm.IsFileExist(containpath) {
			continue
		}
		rmcmd := exec.Command("rmdir",containpath)
		//cm.DPrintf("%v\n",containpath)
		if err := rmcmd.Start(); err != nil {
			return fmt.Errorf("Rmdir(%v) failed: %v from CreateCgroupForContainer %v",containpath,err,containerId)
		}
	}
	return nil
}

func SetCgroupParameter(containerId string,cp float64,mp float64) error {
	return nil
}
