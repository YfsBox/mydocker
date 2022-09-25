package container

import (
	"fmt"
	"golang.org/x/sys/unix"
	"io/ioutil"
	cm "mydocker/common"
	img "mydocker/image"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func getFsMntPath(containerId string) string {
	return fmt.Sprintf("%v/%v/fs/mnt", cm.GetContainerPath(), containerId)
}

func aufsMount(imgfslist []string, containerId string, containerFsPath string, mntPath string) error {
	mntOptions := "lowerdir=" + strings.Join(imgfslist, ":") + ",upperdir=" + containerFsPath + "/writelayer,workdir=" + containerFsPath + "/worklayer"
	cm.DPrintf("The mnt options are %v", mntOptions)

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

func SetUpMount() error {
	// systemd 加入linux之后, mount namespace 就变成 shared by default, 所以你必须显示
	//声明你要这个新的mount namespace独立。
	err := syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
	if err != nil {
		return err
	}
	//mount proc
	//var err error
	//defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV

	err = unix.Mount("proc", "/proc", "proc", 0, "")
	if err != nil {
		fmt.Printf("mount proc error: %v\n", err)
		return fmt.Errorf("mount proc error: %v", err)
	}
	err = unix.Mount("tmpfs", "/dev", "tmpfs", 0, "mode=755")
	if err != nil {
		fmt.Printf("mount tmpfs error\n")
		return fmt.Errorf("mount tmpfs error: %v", err)
	}

	return nil
}

func Unmount(file string) error {
	err := syscall.Unmount(file, 0)
	return err
}

func UnmountAll() error {

	unmountList := [...]string{"/dev", "/proc", "/"}

	for _, fs := range unmountList {
		if err := unix.Unmount(fs, 0); err != nil {
			return fmt.Errorf("Unmount %v error %v", fs, err)
		}
	}

	return nil
}

func ChangeRootDir(containerHash string) error {

	mntPath := getFsMntPath(containerHash)

	fmt.Printf("The mnt Path is %v\n", mntPath)
	defer fmt.Printf("ChangeRootDir ok\n")

	if err := unix.Chroot(mntPath); err != nil {
		fmt.Printf("chroot error")
		return fmt.Errorf("Chroot %v error %v", mntPath, err)
	}
	pwd, _ := os.Getwd()
	cm.DPrintf("the current dir is %v", pwd)

	if err := os.Chdir("/"); err != nil {
		fmt.Printf("chroot error")
		return fmt.Errorf("Chdir to / error %v", err)
	}

	return nil
}

func RemoveContainerFs(containerId string) error {

	containerPath := img.GetContainerPath(containerId)
	mntfsPath := fmt.Sprintf("%v/fs/mnt", containerPath)

	cm.DPrintf("the mntfsPath is %v", mntfsPath)

	//lsofcmd := exec.Command("lsof | grep /var/lib/mydocker/container/BpLnfgDsc2WD/fs/mnt")
	opts := fmt.Sprintf("lsof | grep /var/lib/mydocker/container/%v/fs/mnt", containerId)
	c := exec.Command("bash", "-c", opts)
	out, _ := c.Output()
	cm.DPrintf("the out is:\n%v", out)

	if err := Unmount(mntfsPath); err != nil {
		return fmt.Errorf("Unmount %v error %v", mntfsPath, err)
	}
	if err := os.RemoveAll(containerPath); err != nil {
		cm.DPrintf("Remove %v error %v", containerPath, err)
		return fmt.Errorf("Remove %v error %v", containerPath, err)
	}

	return nil
}

func listAll(path string) {
	files, _ := ioutil.ReadDir(path)
	for _, fi := range files {
		if fi.IsDir() {
			//listAll(path + "/" + fi.Name())
			println(path + "/" + fi.Name())
		} else {
			println(path + "/" + fi.Name())
		}
	}
}
