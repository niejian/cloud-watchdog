// repository doc

package repository

import (
	conf "cloud-watchdog/config"
	"cloud-watchdog/global"
	"cloud-watchdog/model"
	"fmt"
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
	// appName like方式查询
	appName = fmt.Sprintf("%s%s", appName, "%")
	global.GLOBAL_DB.Table(*global.LOG_ALTER_NAME).
		Where("namespace = ? and app_name like ?", ns, appName).Find(&datas)

	if len(datas) <= 0 {
		return nil
	}

	for _, data := range datas {
		conf := &conf.AlterConf{
			ToUserIds: data.ToUserIds,
			Ignores:   nil,
			Errs:      nil,
			AppName:   data.AppName,
			Namespace: data.Namespace,
			EnableStore: data.EnableStore,
			IsEnable: data.IsEnable,
		}
		ignores := data.Ignores
		errs := data.Errs
		if "" != ignores {
			conf.Ignores = strings.Split(ignores, "|")
		}
		if "" != errs {
			conf.Errs = strings.Split(errs, "|")
		}
		confs = append(confs, conf)
	}

	return confs
}
