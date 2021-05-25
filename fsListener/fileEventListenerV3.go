// fsListener doc

package fsListener

import (
	"cloud-watchdog/common"
	"cloud-watchdog/global"
	"cloud-watchdog/logappender"
	"cloud-watchdog/zapLog"
	"github.com/fsnotify/fsnotify"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
	"path/filepath"
	"strings"
	"sync"
)

var fileMap sync.Map
var namespaceMap sync.Map

func ListenAppLogV3(linkFileName string, c *cache.Cache)  {
	// 通过 连接文件名找找到实际文件
	getAppNameAndNamespace(linkFileName, c)

}

func getAppNameAndNamespace(linkFileName string, c *cache.Cache)  {
	var namespace  = ""
	var appName = ""
	var key = common.Md5Str(linkFileName)
	// 缓存namespace等信息
	if val, ok := namespaceMap.Load(key); ok {
		valStr, _ := val.(string)
		vals := strings.Split(valStr, ",")
		namespace = vals[0]
		appName = vals[1]

	} else {
		namespace, appName, _ = common.GetAppNameByLogFileName(linkFileName)
		if namespace == "" || appName == "" {
			zapLog.LOGGER().Error("解析失败: " + linkFileName)
		}
		if "" != namespace && "" != appName {
			namespaceMap.Store(key, namespace + "," + appName)
		}
	}

	// 判断文件是否是我要监听的（podName_namespace_xxxx.log）
	_, loaded := fileMap.LoadOrStore(linkFileName, 1)

	// 该文件没被监听
	if !loaded {
		doLogAppender(namespace, appName, linkFileName, c)
	}
}

func doLogAppender(namespace, appName, linkFileName string, c *cache.Cache)  {

	// 通过连接文件名获取到真正的文件信息 /var/lib/docker/containers/cid/cid-json.log
	// 获取真实的docker日志文件信息
	dockerLogFileName := getDockerLogFilePath(linkFileName)
	if "" == dockerLogFileName {
		zapLog.LOGGER().Error("linkfile：" + linkFileName + ", 无法获取到实际文件信息")
		return
	}
	zapLog.LOGGER().Debug("dockerLogFileName: "+ dockerLogFileName)
	logappender.LogAppender(namespace, appName, dockerLogFileName, c)

}

//ListenLinkfile doc
//@Description: 在linux环境中，fsnotify只能监听到连接文件的创建和删除操作，当监听到文件创建时，找到实际文件
//@Author niejian
//@Date 2021-04-28 14:17:12
//@param c
func ListenLinkfile(c *cache.Cache) {
	// 缓存
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		zapLog.LOGGER().Error("创建文件监听失败", zap.Any("err", err))
		return
	}

	defer watcher.Close()

	/*

	https://github.com/fsnotify/fsnotify/issues?q=Symlinks
	https://blog.csdn.net/u013536232/article/details/104123861
	在linux环境中，fsnotify无法追踪到连接文件
	这里需要做的是需要将这个文件夹中所有文件全部拿出来，每个文件都设置一个watcher
	 */

	// 监听整个目录
	//watcher.a


	err = watcher.Add(*global.K8S_LOG_DIR)
	zapLog.LOGGER().Info("监听文件：" + *global.K8S_LOG_DIR)
	if err != nil {
		zapLog.LOGGER().Error("添加目录监听失败，", zap.String("err", err.Error()))
	}

	done := make(chan bool)

	// 获取操作系统类型

	// 建立监听
	go func() {
		// 接收panic操作
		defer func() {
			if err := recover(); err != nil {
				zapLog.LOGGER().Error("watch panic->recovery", zap.Any("err", err))
			}
		}()

		for {
			select {
			case event := <-watcher.Events:
				fileName := event.Name
				op := event.Op

				// 写操作 查看v2
				// 文件创建操作
				if op&fsnotify.Create == fsnotify.Create {
					zapLog.LOGGER().Info("文件创建操作：", zap.String("fileName", fileName))
					// 通过连接文件找到实际文件信息，并执行tail -f操作
					getAppNameAndNamespace(fileName, c)
				}
				// 文件删除操作
				if op&fsnotify.Remove == fsnotify.Remove {
					dockerLogFilePath := getDockerLogFilePath(fileName)
					if "" == dockerLogFilePath {
						break
					}
					zapLog.LOGGER().Info("文件删除操作， 取消tail -f：", zap.String("source", fileName), zap.String("target", dockerLogFilePath))
					// 取消tail -f
					logappender.StopTail(dockerLogFilePath)
				}

			}
		}
	}()


	<-done
}

//getDockerLogFilePath doc
//@Description: 获取docker日志文件实际目录信息
//@Author niejian
//@Date 2021-05-25 10:11:08
//@param linkFileName
//@return string
func getDockerLogFilePath(linkFileName string) string {

	containerId := common.GetRealFileName(linkFileName)
	if "" == containerId {
		return ""
	}
	// 获取真实的docker日志文件信息
	dockerLogFileName := *global.Docker_Log_Dir + containerId + string(filepath.Separator) +containerId + "-json.log"
	return dockerLogFileName

}


