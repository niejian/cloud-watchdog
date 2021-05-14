// config doc

package config

//日志配置 LogConfig

type LogConfig struct {
	LogLevel   int8   `yaml:"logLevel", json:"logLevel"`     // 日志级别(debug:-1, info:0, warn: 1, error:2)
	LogName    string `yaml:"logName", json:"logName"`       // 日志全路径
	MaxSize    int    `yaml:"maxSize", json:"maxSize"`       // 日志文件最大大小
	MaxBackups int    `yaml:"maxBackups"， json:"maxBackups"` // 日志文件最多保存多少个备份
	MaxAge     int    `yaml:"maxAge", json:"maxAge"`         // 日志保存天数
	LogPath    string `yaml:"logPath", json:"logPath"`
}
