package common

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/urfave/cli"
)

var Debug = true

const (
	mydockerPathRoot = "/var/run/mydocker" //这个是需要检查是否已经创建的
	mydockerLibRoot = "/var/lib/mydocker"
	mydockerImagePath = mydockerLibRoot + "/image"
	mydockerTmpPath = mydockerLibRoot + "/tmp"
	procPathRoot = "/proc"
	cgroupPathRoot = "/sys/fs/cgroup"
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



func JoinPath(root string,sub string) string {
	return fmt.Sprintf("%v/%v",root,sub)
}

func GetPidStr() string {
	return strconv.Itoa(os.Getpid())
}

func IsInStrings(str string,strlist []string) bool {
	for _,s := range strlist {
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
	if _,err := os.Stat(file); err != nil && os.IsNotExist(err) { //如果err != nil && IsNotExist说明不存在
		return false
	}
	return true
}

func PrintArgs(args cli.Args) {

	for _,arg := range args.Tail() {
		DPrintf("%v ",arg)
	}
}

func GetArgsSlice(args cli.Args,argsn int) []string{

	list := make([]string,argsn)

	for _,arg := range args.Tail() {
		DPrintf("arg in GetArgsSlice is %v",arg)
		list = append(list,arg)
	}
	return list
}

