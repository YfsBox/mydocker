package main

import (
	//"fmt"
	//"log"
	//cmd "mydocker/command"
	cm "mydocker/common"
	img "mydocker/image"
	//cnt "mydocker/container"
	//"os"
)

func imageTest() {
	cm.InitMyDockerDirs()
	img.DownloadImageIfNeed("busybox")

}

func main() {
	/*
	if err := cm.InitMyDockerDirs(); err != nil { //如果没有mydocker服务的总文件夹,就创建,否则什么也不做
		log.Fatalf("InitMyDockerDirs error: %v\n", err)
	}

	if err := container.InitCgroupRootDirs(); err != nil {
		log.Fatalf("InitCgroupDirs error: %v\n", err)
	} //检查Cgroups控制组中cpu,memory,pids等subsystem有没有创建好有关于mydocker的根目录

	cm.CheckRootUser()

	app := cmd.InitCliApp()
	cm.RuningCliApp(app, os.Args)

	//cm.RemoveWhenQuit()
	cm.DPrintf("the container quit,begin remove cgroup\n")
	if err := cnt.RemoveCgroupRootDirs(); err != nil {
		log.Fatalf("The RemoveWhenQuit error: %v", err)
	}

	fmt.Printf("Mydocker done(pid: %v),welcome to use!\n", cm.GetPidStr())
	*/
	imageTest()
}
