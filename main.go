package main

import (
	"cloud-watchdog/fsListener"
	"cloud-watchdog/initialize"
	"github.com/patrickmn/go-cache"
	"os"
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
	done := make(chan bool, 1)

	fsListener.ListenAppLog(c)

	<-done
}
