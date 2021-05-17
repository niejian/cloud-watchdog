#!/bin/bash
BinaryName=watchdog-cloud
echo "开始重启...."
./stop.sh $BinaryName
./start.sh $BinaryName
echo "启动成功"

#/bin/sh $COMMOND
