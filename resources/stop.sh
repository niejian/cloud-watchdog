#!/bin/bash

BinaryName=watchdog-cloud
PidFile=watchDog.pid

# 文件不存在
if [ ! -f "$PidFile" ]; then
  #dirName = "project_path=$(cd `dirname $0`; pwd)"

  #PIDS=$(ps -ef | grep $1 | grep -v grep |grep -v 'restart.sh'| awk '{print $2}' | grep -v awk)
  PIDS=`ps -ef | grep $BinaryName | grep -v grep |grep -v 'restart.sh'| awk '{print $2}' | grep -v awk`

  PIDSTR=$(echo $PIDS)
  echo "$PIDSTR" >> $PidFile

fi

# 杀进程
kill -9 `cat $PidFile`
rm -rf $PidFile
echo "【看门狗】 stop SUCCESS!"
