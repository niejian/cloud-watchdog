package main

import (
	"cloud-watchdog/fsListener"
	"cloud-watchdog/global"
	"cloud-watchdog/initialize"
	"cloud-watchdog/zapLog"
	"fmt"
	"github.com/patrickmn/go-cache"
	"io/ioutil"
	"os"
	"path/filepath"
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
	linkFileName := "workorder-control-77d5fd85ff-9gtp5_workorder_workorder-control-c36c983de7cac1092bcd5c94a15e179bc7b051374f2733dbf7c6e118cba2c178.log"
	//linkFileName := "workorder-control-c36c983de7cac1092bcd5c94a15e179bc7b051374f2733db118cba2c178.log"
	//appName := "workorder-control"

	linkFileName = strings.ReplaceAll(linkFileName, string(filepath.Separator), "")
	names := strings.Split(linkFileName, "_")
	var containerId string = names[2]
	containerIds := strings.Split(containerId, "-")
	containerId = containerIds[len(containerIds)-1]
	containerId = strings.Split(containerId, ".")[0]
	//fmt.Println(strings.ReplaceAll(linkFileName, appName+"-", ""))
	fmt.Println(containerId)
}

func main() {

	// 初始化缓存
	c := cache.New(20*time.Second, 30*time.Second)

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
			zapLog.LOGGER().Debug("fileName：" + linkFileName)
			fsListener.ListenAppLogV4(linkFileName, c)
		}
		<-done
	}

}
