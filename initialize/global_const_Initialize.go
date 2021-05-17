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
func GlobalConstInitialize()  {
	config, err := parser.SysConfigParser()
	if err != nil {

	}
	global.K8S_LOG_DIR = &config.GlobalConst.K8sLogDir
	global.LOG_ALTER_NAME = &config.GlobalConst.LogAlterName
}
