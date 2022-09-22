package common

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"io"
	"archive/tar"
	"compress/gzip"
	"github.com/urfave/cli"
	"path/filepath"
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

func Untar(dst, src string) (err error) {
    // 打开准备解压的 tar 包
    fr, err := os.Open(src)
    if err != nil {
        return
    }
    defer fr.Close()

    // 将打开的文件先解压
    gr, err := gzip.NewReader(fr)
    if err != nil {
        return
    }
    defer gr.Close()

    // 通过 gr 创建 tar.Reader
    tr := tar.NewReader(gr)

    // 现在已经获得了 tar.Reader 结构了，只需要循环里面的数据写入文件就可以了
    for {
        hdr, err := tr.Next()

        switch {
        case err == io.EOF:
            return nil
        case err != nil:
            return err
        case hdr == nil:
            continue
        }

        // 处理下保存路径，将要保存的目录加上 header 中的 Name
        // 这个变量保存的有可能是目录，有可能是文件，所以就叫 FileDir 了……
        dstFileDir := filepath.Join(dst, hdr.Name)

        // 根据 header 的 Typeflag 字段，判断文件的类型
        switch hdr.Typeflag {
        case tar.TypeDir: // 如果是目录时候，创建目录
            // 判断下目录是否存在，不存在就创建
            if b := ExistDir(dstFileDir); !b {
                // 使用 MkdirAll 不使用 Mkdir ，就类似 Linux 终端下的 mkdir -p，
                // 可以递归创建每一级目录
                if err := os.MkdirAll(dstFileDir, 0775); err != nil {
                    return err
                }
            }
        case tar.TypeReg: // 如果是文件就写入到磁盘
            // 创建一个可以读写的文件，权限就使用 header 中记录的权限
            // 因为操作系统的 FileMode 是 int32 类型的，hdr 中的是 int64，所以转换下
            file, err := os.OpenFile(dstFileDir, os.O_CREATE|os.O_RDWR, os.FileMode(hdr.Mode))
            if err != nil {
                return err
            }
            n, err := io.Copy(file, tr)
            if err != nil {
                return err
            }
            // 将解压结果输出显示
            fmt.Printf("成功解压： %s , 共处理了 %d 个字符\n", dstFileDir, n)

            // 不要忘记关闭打开的文件，因为它是在 for 循环中，不能使用 defer
            // 如果想使用 defer 就放在一个单独的函数中
            file.Close()
        }
    }

    return nil
}

func ExistDir(dirname string) bool {
    fi, err := os.Stat(dirname)
    return (err == nil || os.IsExist(err)) && fi.IsDir()
}
