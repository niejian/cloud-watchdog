// config doc

package zapLog

import (
	"cloud-watchdog/config"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

//WatchAndReloadConfig doc
//@Description: 监听和热更新配置文件
//@Author niejian
//@Date 2021-04-29 09:36:39
//@param configFilePath
//@param configObj
//@return *interface{}
func watchAndReloadConfig(configFilePath string) *config.LogConfig {

	configObj := &config.LogConfig{}
	viper.SetConfigFile(configFilePath)
	err := viper.ReadInConfig() // 读取配置信息
	if err != nil {
		fmt.Printf("读取配置失败，请重试！%v \n", err)
	}

	// 将读取的配置信息保存到全局变量conf中
	if err = viper.Unmarshal(configObj); err != nil {
		fmt.Printf("序列化配置失败，请重试！%v \n", err)
	}

	// 监控文件变化
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Printf("配置文件更新: %v \n", configFilePath)
		// 判断结构体类型，然后再重新初始化使配置更新
		reloadLogger()
	})

	if err = viper.Unmarshal(configObj); err != nil {
		fmt.Printf("序列化配置失败，请重试！%v \n", err)
	}

	return configObj

}

