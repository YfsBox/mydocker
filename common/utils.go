package common

import (
	"fmt"
	"github.com/urfave/cli"
	"io"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"time"
)

var Debug = true

const (
	mydockerPathRoot      = "/var/run/mydocker" //这个是需要检查是否已经创建的
	mydockerLibRoot       = "/var/lib/mydocker"
	mydockerImagePath     = mydockerLibRoot + "/image"
	mydockerTmpPath       = mydockerLibRoot + "/tmp"
	mydockerContainerPath = mydockerLibRoot + "/container"
	procPathRoot          = "/proc"
	cgroupPathRoot        = "/sys/fs/cgroup"
)

func GetMyDockerPath() string {
	return mydockerPathRoot
}

func GetProcPath() string {
	return procPathRoot
}

func GetCgroupPath() string {
	return cgroupPathRoot
}

func GetTmpPath() string {
	return mydockerTmpPath
}

func GetImagePath() string {
	return mydockerImagePath
}

func GetContainerPath() string {
	return mydockerContainerPath
}

func JoinPath(root string, sub string) string {
	return fmt.Sprintf("%v/%v", root, sub)
}

func GetPidStr() string {
	return strconv.Itoa(os.Getpid())
}

func IsInStrings(str string, strlist []string) bool {
	for _, s := range strlist {
		if s == str {
			return true
		}
	}
	return false
}

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug {
		log.Printf(format, a...)
	}
	return
}

func IsFileExist(file string) bool {
	if _, err := os.Stat(file); err != nil && os.IsNotExist(err) { //如果err != nil && IsNotExist说明不存在
		return false
	}
	return true
}

func PrintArgs(args cli.Args) {

	for _, arg := range args.Tail() {
		DPrintf("%v ", arg)
	}
}

func GetArgsSlice(args cli.Args, argsn int) []string {

	list := make([]string, argsn)

	for _, arg := range args.Tail() {
		DPrintf("arg in GetArgsSlice is %v", arg)
		list = append(list, arg)
	}
	return list
}

func Untar(destDirPath, srcFilePath string) error {
	cmd := exec.Command("tar", "-zxvf", srcFilePath, "-C", destDirPath)
	if err := cmd.Run(); err != nil {
		log.Fatalf("tar %v error %v", srcFilePath, err)
		return err
	}
	return nil
}

func ExistDir(dirname string) bool {
	fi, err := os.Stat(dirname)
	return (err == nil || os.IsExist(err)) && fi.IsDir()
}

func CopyFile(src, dst string) error {
	//1.open src file
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func(in *os.File) {
		err := in.Close()
		if err != nil {
			log.Println("in file close error")
		}
	}(in)
	//2.create dst file
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			log.Println("out file close error")
		}
	}(out)
	//3.copy file content
	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return nil
}

var defaultLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// RandomString returns a random string with a fixed length
func randomString(n int, allowedChars ...[]rune) string {
	var letters []rune

	if len(allowedChars) == 0 {
		letters = defaultLetters
	} else {
		letters = allowedChars[0]
	}
	rand.Seed(time.Now().Unix())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

func RandomString(n int) string {
	return randomString(n, defaultLetters)
}
