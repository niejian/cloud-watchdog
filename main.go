package main

import (
	"cloud-watchdog/fsListener"
	"cloud-watchdog/initialize"
)

func init()  {
	// 初始化环境
	initialize.DoInitialize()
}

func main() {

	done := make(chan bool, 1)

	fsListener.ListenAppLog()

	<-done
}
