// model doc

package model

import "gorm.io/gorm"

type ErrorLogAlterConfig struct {
	gorm.Model
	// 发送工号
	ToUserIds string `json:"toUserIds" gorm:"comment: 发送工号，用|隔开"`
	// 忽略异常
	Ignores string `json:"ignores" gorm:"comment:需要忽略的异常关键字"`
	// errs 告警异常
	Errs string `json:"errs" gorm:"comment: 告警异常"`
	// 应用名称
	AppName string `json:"appName" gorm:"comment: 应用名称"`
	// 命名空间
	Namespace string `json:"namespace" gorm:"comment: 命名空间"`
}

type LogAlterConf struct {
	ToUserIds string `json:"toUserIds" yaml:"toUserIds"`
	Ignores   []string `json:"ignores" yaml:"ignores"`
	errs      []string `json:"errs" yaml:"errs"`
	AppName   string   `json:"appName" yaml:"appName"`
}
