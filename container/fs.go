package container

import (
	"fmt"
	"golang.org/x/sys/unix"
	cm "mydocker/common"
	"os"
	"strings"
)

func aufsMount(imgfslist []string, containerId string, containerFsPath string, mntPath string) error {
	mntOptions := "lowerdir=" + strings.Join(imgfslist, ":") + ",upperdir=" + containerFsPath + "/writelayer,workdir=" + containerFsPath + "/worklayer"
	cm.DPrintf("The mnt options are %v", mntOptions)
	/*cmd := exec.Command("mount", "-t", "aufs", "-o", mntOptions, "none", mntPath)

	fmt.Printf("the cmd is %v\n", cmd.String())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("cmd run err %v", err)
	}*/
	if err := unix.Mount("none", fmt.Sprintf("%v/mnt", containerFsPath), "overlay", 0, mntOptions); err != nil {
		fmt.Printf("mount error\n")
		return fmt.Errorf("mount %v error %v", mntOptions, err)
	}

	return nil
}

func CreateAndMountFs(imgfslist []string, containerId string) error { //创建有关于镜像
	containerFsPath := fmt.Sprintf("%v/%v/fs", cm.GetContainerPath(), containerId)
	mntFsPath := fmt.Sprintf("%v/mnt", containerFsPath)
	paths := [...]string{
		containerFsPath,
		fmt.Sprintf("%v/writelayer", containerFsPath),
		fmt.Sprintf("%v/worklayer", containerFsPath), //这个主要是为了dockerfile中的RUN等构件指令服务的地方
		mntFsPath,
	}

	for _, path := range paths {
		if err := os.MkdirAll(path, 0777); err != nil {
			return fmt.Errorf("the error %v when mkdir %v", err, path)
		}
	}
	if err := aufsMount(imgfslist, containerId, containerFsPath, mntFsPath); err != nil {
		return fmt.Errorf("the mountfs error %v", err)
	}
	return nil
}
