#! /bin/bash


#在宿主机上挂载几个 RAM Disk（内存盘）来模拟本地磁盘。
#假设在名为k8s-node2.1的宿主机上模拟
mkdir /mnt/disks
for vol in vol1 vol2 vol3; do
    mkdir /mnt/disks/$vol
    mount -t tmpfs $vol /mnt/disks/$vol
done
