// fsListener doc

package fsListener

import (
	"cloud-watchdog/api"
	"cloud-watchdog/global"
	"cloud-watchdog/model"
	"cloud-watchdog/zapLog"
	"errors"
	"go.uber.org/zap"
	"strings"
)

//WatchPodLogFile doc
//@Description: pod创建后会新建一个日志文件，可以根据这个文件找到对应的应用名从而找到配置信息
//@Author niejian
//@Date 2021-05-08 11:43:47
//@param fileName
func WatchPodLogFile(fileName string) (*model.ErrorLogAlterConfig, error) {
	var errorLogAlterConfig *model.ErrorLogAlterConfig
	names := strings.Split(fileName, "_")
	if len(names) < 2 {
		zapLog.LOGGER().Error("该文件非k8s应用相关文件", zap.String("fileName", fileName))
		return nil, errors.New("该文件非k8s应用相关文件")
	}
	podName := names[0]
	namespace := names[1]

	// 获取pod的详细信息
	pod, err := api.DescribePod(podName, namespace)
	if err != nil {
		return nil, err
	}
	// 获取标签信息， 通过模板生成的标签 app=xxxx
	labels := pod.Labels
	appName := ""
	for key, val := range labels {
		if key == "app" {
			appName = val
			break
		}
	}
	// 从数据库中获取配置信息
	global.GLOBAL_DB.Table(*global.LOG_ALTER_NAME).
		Where("namespace = ? and app_name = ?", namespace, appName).
		First(&errorLogAlterConfig)

	return errorLogAlterConfig, nil
}
