## 1.简单说明

目前只完成了基本的容器和镜像部分, 初步实现利用linux namesapce的隔离和cgroup的资源限制以及利用aufs的文件系统。
接下来要做的部分是容器的exec,ps,stop等功能,以及网络部分。

## 2.预备工作

### 1)编译
```shell
$ go build main.go
```
### 2)运行时开启root权限
```shell
$ sudo ./main ..........
#或者
$ sudo su
```
## 3.使用的几个demo
目前只支持三种功能指令:run,pull,images,help.
### 1) pull
用来基于某个镜像上运行一个容器,该指令的flag如下:
```shell
OPTIONS:
   --img value  要拉取的镜像名称
```
只有一个flag,也就是要拉取镜像的名称,执行一下命令:
```shell
$ sudo ./main pull -img arm64v8/ubuntu:20.04
```
在这里以ubuntu:20.04镜像为例，由于我的机器是arm架构所以下载的是arm64v8/ubuntu:20.04,
如果是x86架构请pull`amd64/ubuntu:20.04`.
执行该指令后效果如下:
```shell
$ ./main pull -img arm64v8/ubuntu:20.04
2022/09/28 02:38:45 Not have this images,begin downloading......
2022/09/28 02:38:48 the imageHexHash is bdbe84df0b98
2022/09/28 02:39:02 SaveLegacy ok!
2022/09/28 02:39:02 tar ImageFiles......
2022/09/28 02:39:02 name:sha256:bdbe84df0b988965129e5189bae37e09e29d17d090cf4b8c8cb8dabf2443fd77
2022/09/28 02:39:02 name:7a9f619ee5e9c87f19eed59abef41d53eb0694f492da010ee069ff26e7b4ff3f.tar.gz
2022/09/28 02:39:02 name:manifest.json
2022/09/28 02:39:02 tar file reach end of file EOF!
Mydocker done(pid: 75661),welcome to use!
```
下载之后可以通过imges指令来验证是否存在.

### 2) images
该指令没有什么flag,直接如此执行就可以:
```shell
$ sudo ./main images
```
效果如下:
```shell
--------------------------------------------------
TAG              |           IMAGEHASH |
--------------------------------------------------
arm64v8/alpine          |        a6215f271958 |
arm64v8/ubuntu          |        21735dab04ba |
arm64v8/ubuntu:20.04    |        bdbe84df0b98 |
--------------------------------------------------
2022/09/28 02:05:05 the container quit,begin remove cgroup
Mydocker done(pid: 75449),welcome to use!
```
可见之前pull的`arm64v8/ubuntu:20.04`已经记录好了.
### 3) run
run这个指令其实分了两步走,第一步首先判断这个image有没有之前下过,如果没有就会下载相当于(pull)。然后生成并运行该容器，并执行指定
的命令.
执行:
```shell
$ sudo ./main run -img arm64v8/ubuntu:20.04 -cmds /bin/bash -name test
```
下过如下:
```shell
2022/09/28 02:46:55 The image has alreay exist!
2022/09/28 02:46:55 will Begin RunExec!
2022/09/28 02:46:55 cpu have no config
2022/09/28 02:46:55 mem have no config
2022/09/28 02:46:55 pid have no config
2022/09/28 02:46:55 Begin Running cmd!
root@lgQ8eGPVjmw5:/# ls                                                                                                                                                             
bin  boot  dev  etc  home  lib  media  mnt  opt  proc  root  run  sbin  srv  sys  tmp  usr  var
root@lgQ8eGPVjmw5:/# ps
    PID TTY          TIME CMD
      1 ?        00:00:00 exe
      7 ?        00:00:00 bash
     11 ?        00:00:00 ps
root@lgQ8eGPVjmw5:/# id
uid=0(root) gid=0(root) groups=0(root)
root@lgQ8eGPVjmw5:/# hostname
lgQ8eGPVjmw5
```
如果想加入cgroup的资源限制,可以执行类似于下面的命令:
```shell
$ sudo ./main run -cpus 0.5 -mmem 200 -mpid 1000 -img arm64v8/ubuntu:20.04 -cmds /bin/bash -name test
```
就cgroup资源限制相关的功能，还没有进行太多的测试.