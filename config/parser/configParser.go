// parser doc

package parser

import (
	"cloud-watchdog/config"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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

	// 判断是否是k8s环境，如果是，那么数据库，es读取k8s部署文件中的配置信息
	enableK8s := os.Getenv("ENABLE_K8S")
	if "" != enableK8s && "true" == enableK8s {
		mysqlUrl := os.Getenv("MYSQL_URLS")
		mysqlUsername := os.Getenv("MYSQL_USERNAME")
		mysqlPassword := os.Getenv("MYSQL_PASSWORD")
		// 设置库名
		mysqlDbName := os.Getenv("MYSQL_DB_NAME")

		if "" != mysqlUrl && "" != mysqlUsername && "" != mysqlPassword {
			sysConfig.Mysql.Path = mysqlUrl
			sysConfig.Mysql.Username = mysqlUsername
			sysConfig.Mysql.Password = mysqlPassword
			sysConfig.Mysql.Dbname = mysqlDbName
		}

		fmt.Println("读取k8s配置信息 ")
		esUrls := os.Getenv("ES_URLS")
		esUsername := os.Getenv("ES_USERNAME")
		esPassword := os.Getenv("ES_PASSWORD")

		if "" != esUrls && "" != esUsername && "" != esPassword {
			// 替换
			sysConfig.Es.Urls = strings.Split(esUrls, ",")
			sysConfig.Es.Username = esUsername
			sysConfig.Es.Password = esPassword
		}

	}

	return sysConfig, nil
}
