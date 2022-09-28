#!/bin/bash
mkdir aufs
mkdir aufs/mnt
mkdir aufs/container-layer
echo "I am container layer" > aufs/container-layer/container-layer.txt
mkdir aufs/{image-layer1,image-layer2,image-layer3}
echo "I am image layer 1" > aufs/image-layer1/image-layer1.txt
echo "I am image layer 2" > aufs/image-layer2/image-layer2.txt
echo "I am image layer 3" > aufs/image-layer3/image-layer3.txt

cd aufs
sudo mount -t aufs -o dirs=./container-layer:./image-layer1:./image-layer2:./image-layer3 none ./mnt
