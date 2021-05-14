// common doc

package common

import (
	"cloud-watchdog/api"
	conf "cloud-watchdog/config"
	"cloud-watchdog/global"
	"cloud-watchdog/repository"
	"cloud-watchdog/zapLog"
	"errors"
	"go.uber.org/zap"
	"path/filepath"
	"strings"
)

//GetLogAlterConfByFileName doc
//@Description: 通过namespace、appName获取到配置信息
//@Author niejian
//@Date 2021-05-13 14:15:03
//@param namespace
//@param appName
//@return *conf.AlterConf
func GetLogAlterConfByFileName(namespace, appName string) (*conf.AlterConf, error) {

	if "" == namespace || "" == appName {
		return nil, errors.New("获取ns，appName失败")
	}

	// 根据ns， appName 获取配置信息
	confs := repository.ListLogAlterConfByAppNameAndNamespace(namespace, appName)
	if len(confs) == 0 {
		return nil, errors.New("not find")
	}
	return confs[0], nil

	//confs := repository.ListLogAlterConfByAppNameAndNamespace("apollo", "portal-apollo-portal")
	return confs[0], nil
}


//GetAppNameByLogFileName doc
//@Description: 通过日志文件名获取appName
//@Author niejian
//@Date 2021-05-14 11:34:53
//@param fileName
//@return string
//@return error
func GetAppNameByLogFileName(fileName string) (string, string, error)  {
	if strings.Contains(fileName, global.K8S_LOG_DIR) {
		fileName = strings.ReplaceAll(fileName,  global.K8S_LOG_DIR, "")
	}

	fileName = strings.ReplaceAll(fileName, string(filepath.Separator), "")
	names := strings.Split(fileName, "_")
	zapLog.LOGGER().Info("解析日志文件", zap.String("fileName", fileName))
	podName := ""
	namespace := ""

	if len(names) > 1 {
		podName = names[0]
		namespace = names[1]
	}else{
		return "", "", errors.New(fileName + ", not a format pod log file")
	}
	// 根据podName，namespace获取到相关pod详细信息，label包含appName
	pod, err := api.DescribePod(podName, namespace)
	if nil != err {
		zapLog.LOGGER().Error("获取pod信息失败",
			zap.String("podName", podName),
			zap.String("namespace", namespace),
			zap.String("err", err.Error()))
		return "", "", err
	}
	labels := pod.Labels
	// 获取 app标签
	appName := labels["app"]

	return namespace, appName, nil
}

