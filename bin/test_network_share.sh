#!/bin/zsh

for cid in `docker ps -aq --filter name=busybox --last 2`; do
    echo "容器 $cid 网络配置："
    docker logs -f $cid
    echo "--------------------"
done
