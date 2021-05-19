// initialize doc

package initialize

import (
	"cloud-watchdog/config/parser"
	"cloud-watchdog/global"
)

//GlobalConstInitialize doc
//@Description: 常量初始化
//@Author niejian
//@Date 2021-05-17 15:05:20
func GlobalConstInitialize(env string)  {
	config, _ := parser.SysConfigParser()

	if env == "dev" {
		global.K8S_LOG_DIR = &config.GlobalConst.Dev.K8sLogDir
	}else if env == "pro"{
		global.K8S_LOG_DIR = &config.GlobalConst.Pro.K8sLogDir
	}

	global.LOG_ALTER_NAME = &config.GlobalConst.LogAlterName
}
