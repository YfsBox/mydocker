package command

import (
	"fmt"
	"github.com/urfave/cli"
	"log"
	cm "mydocker/common"
	cnt "mydocker/container"
	img "mydocker/image"
	"os"
)

const (
	RunFlag    = "run"
	ExecFlag   = "Init"
	defaultCmd = "zsh"
)

func parseRunFlags(context *cli.Context) []string {

	flaglist := []string{}
	for _, flag := range context.FlagNames() {
		if flag == "img" {
			continue
		}
		flaglist = append(flaglist, fmt.Sprintf("-%v", flag))
		if flag != "cmds" {
			if str := context.String(flag); true {
				flaglist = append(flaglist, str)
			} //如果真的是空的,就相当于对cpu等资源不作出限制了
		} else {
			if cmds := context.StringSlice(flag); len(cmds) != 0 {
				flaglist = append(flaglist, cmds...)
			} else {
				defualtcmd := []string{defaultCmd}
				flaglist = append(flaglist, defualtcmd...) //
			}
		}
	}
	return flaglist
}

var RunCommand = cli.Command{
	Name:  "run",
	Usage: "在指定镜像上运行一个容器",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "cpus",
			Usage: "对cpu的限制,限制最大cpu占用率.",
		},
		&cli.StringFlag{
			Name:  "mpid",
			Usage: "对于pid的限制,",
		},
		&cli.StringFlag{
			Name: "mmem",
		},
		&cli.StringFlag{
			Name:  "name",
			Usage: "给你要运行的容器起一个名字",
		},
		&cli.StringSliceFlag{
			Name:  "cmds",
			Usage: "给容器指定的运行命令,比如/bin/bash",
		},
		&cli.StringFlag{
			Name:  "img",
			Usage: "该容器运行所基于的镜像",
		},
	},
	Action: func(context *cli.Context) error {

		flags := []string{}
		flags = append(flags, "runexec")
		flags = append(flags, parseRunFlags(context)...)
		containerId := cnt.GetContainerId()
		flags = append(flags, "-cid")
		flags = append(flags, containerId)

		//cmdargs = append(cmdargs,cmdlist...)
		//构造供exec执行的命令及其参数列表

		hash, need, err := img.DownloadImageIfNeed(context.String("img"))
		if err != nil {
			fmt.Printf("the error is %v", err)
		}
		flags = append(flags, "-imghash")
		flags = append(flags, hash)

		imagefsList, _ := img.ProcessLayers(hash, need)
		img.RemoveTmpImage(hash)
		if err := cnt.CreateAndMountFs(imagefsList, containerId); err != nil { //难道需要将read部分
			return fmt.Errorf("CreateContainerDirs error %v", err)
		}

		cmd := cnt.GetCloneContainerProc("/proc/self/exe", flags) //开始执行特定的命令,目前就先暂定为shell,也就是什么都没有的意思
		cmd.Run()
		cmd.Wait()

		if err := cnt.RemoveContainerFs(containerId); err != nil {
			log.Fatalf("RemoveContainerFs error %v", err)
		}
		if err := cnt.RemoveCgroupForContainer(containerId); err != nil {
			log.Fatalf("RemoveCgroupForContainer %v err from Run", err)
		}
		//fmt.Printf("a container quit successfully!\n")
		return nil
	},
}

var RunExecCommand = cli.Command{ //该指令是从属于run的,属于run的一部分
	Name:  "runexec",
	Usage: "该指令属于run的一部分,用于run中,执行传入命令的阶段",

	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "cpus",
		},
		&cli.StringFlag{
			Name: "mpid",
		},
		&cli.StringFlag{
			Name: "mmem",
		},
		&cli.StringFlag{
			Name:  "name",
			Usage: "the name of running container",
		},
		&cli.StringSliceFlag{
			Name: "cmds",
		},
		&cli.StringFlag{
			Name:  "cid",
			Usage: "the running container's id",
		},
		&cli.StringFlag{
			Name: "imghash",
		},
	},

	Action: func(context *cli.Context) error {
		//移除有关文件挂载的部分
		if context.String("cid") == "" {
			return fmt.Errorf("the ContainerId is null when runexec!")
		} //必须要有一个containerId
		//接下来根据context来构造一个cgroup的结构体
		limit := cnt.GetCgroupLimit(context.String("cpus"), context.String("mmem"), context.String("mpid"))
		cm.DPrintf("will Begin RunExec!\n")
		if err := RunExec(context.StringSlice("cmds"), context.String("cid"), context.String("imghash"), limit); err != nil {
			return fmt.Errorf("RunExec error %v in RunExecCommand", err)
		}
		os.Exit(-1)
		return nil
	},
}

var PullCommand = cli.Command{
	Name:  "pull",
	Usage: "拉取一个镜像",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "img",
			Usage: "要拉取的镜像名称",
		},
	},
	Action: func(context *cli.Context) error {
		imgName := context.String("img")
		var hash string
		var err error
		var need bool

		if hash, need, err = img.DownloadImageIfNeed(imgName); err != nil {
			log.Fatalf("download the img %v error %v", imgName, err)
		}
		//cm.DPrintf("begin process")
		if _, err = img.ProcessLayers(hash, need); err != nil {
			log.Fatalf("ProcessLayers %v error: %v", hash, err)
		}
		if err = img.RemoveTmpImage(hash); err != nil {
			log.Fatalf("RemoveTmpImage %v error %v", hash, err)
		}

		return nil
	},
}

var ImagesCommand = cli.Command{
	Name:  "images",
	Usage: "展示目前已经下载好的镜像",
	Action: func(ctx *cli.Context) error {
		return img.ShowImages()
	},
}

func InitCliApp() *cli.App {
	app := cli.NewApp()
	app.Name = "mydocker"
	app.Usage = "just for fun!" //设置一些基本的信息

	app.Commands = []cli.Command{
		RunCommand,
		RunExecCommand,
		ImagesCommand,
		PullCommand,
	}
	return app
}
