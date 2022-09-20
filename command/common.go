package command

import (
	"fmt"
	"github.com/urfave/cli"
	cm "mydocker/common"
	cnt "mydocker/container"
)

const (
	RunFlag = "run"
	ExecFlag = "Init"
)

var RunCommand = cli.Command {
	Name: "run",
	Usage: "run a container from a image",
	Action: func(context *cli.Context) error {
		
		
		//defer cnt.RemoveCgroupForContainer(containerId)
		containerId := RunInit()
		cmdargs := []string{}
		cmdargs = append(cmdargs,"runexec")
		cmdargs = append(cmdargs, containerId)
		for _,arg := range context.Args().Tail() {
			cmdargs = append(cmdargs,arg)
		}
		//cmdargs = append(cmdargs,cmdlist...)
		//构造供exec执行的命令及其参数列表
		cmd := cnt.GetCloneContainerProc("/proc/self/exe",cmdargs) //开始执行特定的命令,目前就先暂定为shell,也就是什么都没有的意思
		if err := cmd.Run(); err != nil {
			fmt.Printf("cmd.Run error\n")
		}
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
			fmt.Printf("cmd run error in exec\n")
		}//退出该容器后的处理
		return nil
	},
}

var RunExecCommand = cli.Command { //该指令是从属于run的,属于run的一部分
	Name: "runexec",
	Usage: "该指令属于run的一部分,用于run中,执行传入命令的阶段",

	Action: func(context *cli.Context) error {
		defer cnt.UnmountAll() //移除有关文件挂载的部分

		cm.DPrintf("the argN is %v",context.NArg())
		cmdargs := []string{}
		cmdargs = append(cmdargs,context.Args().First())
		for _,arg := range context.Args().Tail() {
			//cm.DPrintf("the arg add to arglist is %v",arg)
			cmdargs = append(cmdargs,arg)
		}
		if err := RunExec(cmdargs[1:],cmdargs[0]); err != nil {
			return fmt.Errorf("RunExec error %v in RunExecCommand",err)
		}
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