package main

import (
	"cloud-watchdog/fsListener"
	"cloud-watchdog/initialize"
	"github.com/patrickmn/go-cache"
	"time"
)

func init() {
	// 初始化环境
	initialize.DoInitialize()
}

func main() {

	// 初始化缓存
	c := cache.New(20*time.Second, 30*time.Second)
	done := make(chan bool, 1)

	fsListener.ListenAppLog(c)

	<-done
}
