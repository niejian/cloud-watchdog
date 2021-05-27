package main

import (
	"cloud-watchdog/fsListener"
	"cloud-watchdog/global"
	"cloud-watchdog/initialize"
	"cloud-watchdog/zapLog"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"time"
)

func init() {

	// 获取环境变量信息
	env := os.Getenv("ENV")
	if "" == env {
		env = "dev"
	}
	// 获取命令行参数信息，默认环境dev
	//env := flag.String("env", "dev1", "get env")
	// 初始化环境
	initialize.DoInitialize(env)

}

func main() {

	// 初始化缓存
	c := cache.New(20*time.Second, 30*time.Second)
	defer func() {
		if err := recover(); err != nil {
			zapLog.LOGGER().Error("main recovery", zap.Any("err", err))
		}
	}()

	osType := runtime.GOOS
	// mac os
	if osType == "darwin" {
		done := make(chan bool, 1)
		fsListener.ListenAppLog(c)
		<-done
	}

	// linux
	if osType == "linux" {
		done := make(chan bool, 1)

		// 获取所有文件信息
		// 获取文件数
		infos, _ := ioutil.ReadDir(*global.K8S_LOG_DIR)
		for _, info := range infos {

			linkFileName := info.Name()
			// 过滤非系统命名空间
			names := strings.Split(linkFileName, "_")
			if len(names) < 2 {
				continue
			}
			namespace := names[1]
			// 判断namespace是否是过滤的namespace
			isContinue := true
			for _, ns := range *global.Exclude_Ns {
				if ns == namespace || strings.HasPrefix(namespace, ns) {
					zapLog.LOGGER().Info(linkFileName + ", ns=" + namespace + "， 为忽略的命名空间，不处理")
					isContinue = false
					break
				}
			}
			if isContinue {
				fsListener.ListenAppLogV3(linkFileName, c)
			}
		}

		// 监听连接文件的创建和删除操作
		fsListener.ListenLinkfile(c)
		<-done
	}

}
