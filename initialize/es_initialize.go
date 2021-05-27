// initialize doc

package initialize

import (
	"cloud-watchdog/config/parser"
	esService "cloud-watchdog/es"
	"cloud-watchdog/global"
	"log"
)

//EslInitialize doc
//@Description: es初始化
//@Author niejian
//@Date 2021-05-26 16:49:37
func EslInitialize()  {
	config, _ := parser.SysConfigParser()
	es := config.Es
	client, err := esService.InitEs(&es)
	if err != nil {
		log.Printf("初始化es失败: %v \n", err)
	} else {
		global.Es_Client = client
	}


}
