// repository doc

package repository

import (
	conf "cloud-watchdog/config"
	"cloud-watchdog/global"
	"cloud-watchdog/model"
	"strings"
)

//ListAllLogAlterConf doc
//@Description: 获取所有的配置信息
//@Author niejian
//@Date 2021-05-13 10:49:31
//@return []*conf.AlterConf
//@return error
func ListAllLogAlterConf() ([]*conf.AlterConf, error) {
	var datas []*conf.AlterConf
	global.GLOBAL_DB.Table(*global.LOG_ALTER_NAME).Find(datas)
	return datas, nil
}

func ListLogAlterConfByAppNameAndNamespace(ns, appName string) []*conf.AlterConf  {
	var datas []model.ErrorLogAlterConfig
	var confs []*conf.AlterConf
	global.GLOBAL_DB.Table(*global.LOG_ALTER_NAME).
		Where("namespace = ? and app_name = ?", ns, appName).Find(&datas)

	if len(datas) <= 0 {
		return nil
	}

	for _, data := range datas {
		conf := &conf.AlterConf{
			ToUserIds: data.ToUserIds,
			Ignores:   strings.Split(data.Ignores, "|"),
			Errs:      strings.Split(data.Errs, "|"),
			AppName:   data.AppName,
			Namespace: data.Namespace,
		}
		confs = append(confs, conf)
	}

	return confs
}
