#! /bin/bash
#在模拟的宿主机上执行
for vol in vol1 vol2 vol3; do
    umount $vol /mnt/disks/$vol
done
