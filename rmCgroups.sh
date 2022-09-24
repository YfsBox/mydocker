#!/bin/zsh
rootpath=/sys/fs/cgroup
list=(cpu memory pids)
containerId=$1
rmAll=$2

for element in ${list[@]}
do
    echo rmdir ${rootpath}/${element}/mydocker/${containerId}
    rmdir ${rootpath}/${element}/mydocker/${containerId}
done



