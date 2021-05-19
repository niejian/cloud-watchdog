!/bin/sh
BinaryName=watchdog-cloud
PidFile=watchDog.pid
DirName=$(cd $(dirname $0); pwd)
cd $DirName
# 设置环境变量信息
export APP_LOG_CONF_FILE=$DirName/resources/dubbo-log.yml
export CONF_CONSUMER_FILE_PATH=$DirName/resources/client.yml
export ENV=pro

if [ -f "$PIDFILE" ]; then
    echo "【看门狗】已启动 ..., 运行stop.sh"
    ./stop.sh
else
  echo "【看门狗】开始启动..."
  nohup $DirName/$BinaryName 2>&1 &
  printf '%d' $! > $PidFile
  echo "【看门狗】启动成功"
fi
