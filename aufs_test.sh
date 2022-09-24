#!/bin/zsh

mkdir aufs_test && cd aufs_test
mkdir dir0 dir1 root

echo dir0 > dir0/001.txt
echo dir0 > dir0/002.txt
echo dir1 > dir1/002.txt
echo dir1 > dir1/003.txt

sudo mount -t aufs -o br=./dir0=ro:./dir1=ro none ./root
#ls root/