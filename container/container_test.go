package container

import (
	"fmt"
	cm "mydocker/common"
	//"os/exec"
	img "mydocker/image"
	"testing"
)

func initAndCreateFs(imageName string) string {
	cm.InitMyDockerDirs()
	hash, err := img.DownloadImageIfNeed(imageName)
	if err != nil {
		fmt.Printf("the error is %v", err)
	}
	cm.DPrintf("hash is %v\n", hash)
	imglist, _ := img.ProcessLayers(hash)
	CreateAndMountFs(imglist, hash)
	img.RemoveTmpImage(hash)

	return hash
}

func shownewContainer() string {
	ctnId := GetContainerId()
	cm.DPrintf("the new containerId is %v\n", ctnId) //对于这个ContainId暂时还有一定的疑问
	return ctnId
}

func TestCgroupsInitAndRemoveRoot(t *testing.T) { //删除某个Root还是可以的,但是独居mydocker底下的一个目录目前有问题
	//ctnId := shownewContainer()
	//反复调用,一定不要忘了开root
	for i := 0; i <= 2; i++ {
		if err := InitCgroupRootDirs(); err != nil {
			cm.DPrintf("InitCgroupRootDirs error: %v\n", err)
		}
	}
	for i := 0; i <= 2; i++ { //多次调用Remove尝试
		if err := RemoveCgroupRootDirs(); err != nil {
			cm.DPrintf("RemoveCgroupRootDirs error: %v\n", err)
		}
	}

}

func TestCgroupCreateAndRemove(t *testing.T) { //这里有个问题,当两个终端同时操作某个cgroup文件夹的时候,会导致无法删除(rmdir)

	ctnId := shownewContainer()

	InitCgroupRootDirs()

	for i := 0; i <= 2; i++ { //反复删除调用
		if err := CreateCgroupForContainer(ctnId); err != nil {
			cm.DPrintf("CreateCgroupForContainer error: %v\n", err)
		}
		if err := RemoveCgroupForContainer(ctnId); err != nil {
			cm.DPrintf("RemoveCgroupForContainer error: %v\n", err)
		}
	}

	RemoveCgroupRootDirs()

}

func TestCreateAndMountFs(t *testing.T) {
	hash := initAndCreateFs("ubuntu:16.04")
	config := img.ParseContainerConfig(hash)
	cm.DPrintf("the config of %v is %v,the env is %v", hash, config, config.Config.Env)

	RemoveContainerFs(hash)
}

func TestChangeRootDir(t *testing.T) {
	hash := initAndCreateFs("busybox")
	ChangeRootDir(hash)
}

/*
func TestIsolation(t *testing.T) {
	cmd := GetCloneContainerProc("")
	cmd.Run()
}
*/
