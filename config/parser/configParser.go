// parser doc

package parser

import (
	"cloud-watchdog/config"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

var logConfigFileName = "log.yaml"

//LogConfigParser doc
//@Description: 日志配置解析
//@Author niejian
//@Date 2021-04-26 17:26:53
//@return *config.LogConfig
func LogConfigParser() *config.LogConfig {
	logConfig := &config.LogConfig{}
	projectPath, _ := os.Getwd()
	logConfigPath := fmt.Sprintf("%s%s%s%s%s", projectPath, string(filepath.Separator), "resources",
		string(filepath.Separator), logConfigFileName)
	file, err := ioutil.ReadFile(logConfigPath)
	if err != nil {
		fmt.Printf("初始化日志配置失败，err:%v \n", err)
	}

	err = yaml.Unmarshal(file, logConfig)
	if err != nil {
		fmt.Printf("日志配置转化yaml失败，err:%v \n", err)
	}
	return logConfig
}

//SysConfigParser doc
//@Description: 解析系统配置
//@Author niejian
//@Date 2021-05-08 14:36:12
//@return *config.SysConfig
//@return error
func SysConfigParser() (*config.SysConfig, error) {
	sysConfig := &config.SysConfig{}

	projectPath, _ := os.Getwd()
	logConfigPath := fmt.Sprintf("%s%s%s%s%s", projectPath, string(filepath.Separator), "resources",
		string(filepath.Separator), "conf.yaml")
	file, err := ioutil.ReadFile(logConfigPath)
	if err != nil {
		fmt.Printf("初始化日志配置失败，err:%v \n", err)
	}

	err = yaml.Unmarshal(file, sysConfig)
	if err != nil {
		fmt.Printf("日志配置转化yaml失败，err:%v \n", err)
	}

	return sysConfig, nil
}
