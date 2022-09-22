package command

import (
	"fmt"
	cm "mydocker/common"
	cnt "mydocker/container"
	"os"
	"github.com/urfave/cli"
)

const (
	RunFlag = "run"
	ExecFlag = "Init"
	defaultCmd = "zsh"
)

func parseRunFlags(context *cli.Context) []string{  

	flaglist := []string{}
	for _,flag := range context.FlagNames() {
		flaglist = append(flaglist,fmt.Sprintf("-%v",flag))
		if flag != "cmds" {
			if str := context.String(flag); true {
				flaglist = append(flaglist,str)
			} //如果真的是空的,就相当于对cpu等资源不作出限制了
		} else {
			if cmds := context.StringSlice(flag); len(cmds) != 0 {
				flaglist = append(flaglist,cmds...)
			} else {
				defualtcmd := []string{defaultCmd}
				flaglist = append(flaglist,defualtcmd...) //
			}
		}
	}
	return flaglist
}


var RunCommand = cli.Command {
	Name: "run",
	Usage: "run a container from a image",
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
			Name: "name",
			Usage: "the name of running container",
		},
		&cli.StringSliceFlag{
			Name: "cmds",
		},
	},
	Action: func(context *cli.Context) error {

		flags := []string{}
		flags = append(flags,"runexec")
		flags = append(flags,parseRunFlags(context)...)
		containerId := RunInit()
		flags = append(flags,"-cid")
		flags = append(flags, containerId)

		//cmdargs = append(cmdargs,cmdlist...)
		//构造供exec执行的命令及其参数列表
		defer cnt.UnmountAll() 
		cmd := cnt.GetCloneContainerProc("/proc/self/exe",flags) //开始执行特定的命令,目前就先暂定为shell,也就是什么都没有的意思
		cmd.Run()
		fmt.Printf("a container quit successfully!\n")
		return nil
	},
}

var ExecCommand = cli.Command {
	Name: "exec",
	Usage: "exec a program on a running container",
	Action: func(context *cli.Context) error {

		fmt.Printf("exec a program on a running container,the pid is %v,the argsN is %v\n",cm.GetPidStr(),context.NArg())

		cmdargs := make([]string,context.NArg())
		cmd := cnt.GetCloneContainerProc("zsh",cmdargs)//目前默认开启zsh终端执行

		if err := cmd.Run(); err != nil {
			fmt.Printf("cmd run error %v in exec\n",err)
		}//退出该容器后的处理
		return nil
	},
}

var RunExecCommand = cli.Command { //该指令是从属于run的,属于run的一部分
	Name: "runexec",
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
			Name: "name",
			Usage: "the name of running container",
		},
		&cli.StringSliceFlag{
			Name: "cmds",
		},
		&cli.StringFlag{
			Name: "cid",
			Usage: "the running container's id",
		},
	},

	Action: func(context *cli.Context) error {
		//移除有关文件挂载的部分
		if context.String("cid") == "" {
			return fmt.Errorf("the ContainerId is null when runexec!")
		} //必须要有一个containerId
		cm.DPrintf("the argN is %v",context.NArg())
		//接下来根据context来构造一个cgroup的结构体
		limit := cnt.GetCgroupLimit(context.String("cpus"),context.String("mmem"),context.String("mpid"))
		
		cm.DPrintf("will RunExec\n")
		if err := RunExec(context.StringSlice("cmds"),context.String("cid"),limit); err != nil {
			return fmt.Errorf("RunExec error %v in RunExecCommand",err)
		}
		os.Exit(-1)
		return nil
	},
}

var PullCommand = cli.Command {
	Name: "pull",
	Usage: "pull an image",

	Action: func(context *cli.Context) error {


		return nil
	},

}


func InitCliApp() *cli.App {
	app := cli.NewApp()
	app.Name = "mydocker"
	app.Usage = "just for fun!" //设置一些基本的信息
	
	app.Commands = []cli.Command {
		RunCommand,
		RunExecCommand,
		ExecCommand,
	}
	
	return app
}