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
	"strings"
	"sync"
)

//ListenAppLog doc
//@Description: 监听k8s中应用容器的日志文件变化信息。当有追加的指令时，就做相关操作
//@Author niejian
//@Date 2021-04-28 14:17:12
//@param logDir
//@param appName
func ListenAppLog(c *cache.Cache) {
	// 缓存
	var fileMap sync.Map
	var namespaceMap sync.Map
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		zapLog.LOGGER().Error("创建文件监听失败", zap.Any("err", err))
		return
	}

	defer watcher.Close()

	// 监听整个目录
	err = watcher.Add(*global.K8S_LOG_DIR)
	zapLog.LOGGER().Info("监听文件：" + *global.K8S_LOG_DIR)
	if err != nil {
		zapLog.LOGGER().Error("添加目录监听失败，", zap.String("err", err.Error()))
	}

	done := make(chan bool)

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
				zapLog.LOGGER().Info("op: " + string(op))

				// 写操作
				if op&fsnotify.Write == fsnotify.Write {
					var namespace string
					var appName string
					// 缓存namespace等信息
					if val, ok := namespaceMap.Load(fileName); ok {
						valStr, _ := val.(string)
						vals := strings.Split(valStr, ",")
						namespace = vals[0]
						appName = vals[1]

					} else {
						namespace, appName, err = common.GetAppNameByLogFileName(fileName)
						if err != nil {
							zapLog.LOGGER().Error("解析失败"+ err.Error())
						}
						if "" != namespace && "" != appName {
							namespaceMap.Store(fileName, namespace + "," + appName)
						}
					}

					// 判断文件是否是我要监听的（podName_namespace_xxxx.log）
					// 判断是否已经监听了这个文件
					isMonitor := false

					if _, ok := fileMap.Load(fileName); ok {
						isMonitor = true
					}
					fileMap.LoadOrStore(fileName, 1)

					// 追加日志
					if !isMonitor {
						go func() {
							// panic处理
							defer func() {
								if err := recover(); err != nil {
									zapLog.LOGGER().Error("tail log panic, prepare recovery", zap.Any("err", err))
									// panic之后要重新赋值，状态复位
									isMonitor = false
									fileMap.Delete(fileName)
								}
							}()

							logappender.LogAppender(namespace, appName, fileName, c)

						}()
					}
				}
				// 文件创建操作
				if op&fsnotify.Create == fsnotify.Create {
					// 有新的pod创建，从文件名获取到配置信息
					//zapLog.LOGGER().Info("新增pod日志文件", zap.String("pod_log_file", fileName))
					//namespace, appName, err := common.GetAppNameByLogFileName(fileName)
					//if err != nil {
					//	zapLog.LOGGER().Error("解析失败"+ err.Error())
					//}
					//conf, err := common.GetLogAlterConfByFileName(namespace, appName)
					//if nil != err {
					//	zapLog.LOGGER().Error("err", zap.String("err", err.Error()))
					//	continue
					//}
					//if nil == conf {
					//	zapLog.LOGGER().Error("配置为空")
					//	continue
					//}
					//ListenAppLogV2()
				}
				// 文件删除操作
				if op&fsnotify.Remove == fsnotify.Remove {
					zapLog.LOGGER().Info("文件删除操作：", zap.String("fileName", fileName))
					// 取消监听
					//watcher.Remove(fileName)
					//
					panic(fileName + "已经删除，主动断开操作")
					break
				}

			}
		}
	}()


	<-done
}


