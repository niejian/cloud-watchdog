// fsListener doc

package fsListener

import (
	"cloud-watchdog/common"
	"cloud-watchdog/global"
	"cloud-watchdog/logappender"
	"cloud-watchdog/zapLog"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
	"path/filepath"
	"strings"
	"sync"
)

var fileMap sync.Map
var namespaceMap sync.Map

func ListenAppLogV4(linkFileName string, c *cache.Cache)  {
	// 通过 连接文件名找找到实际文件
	getAppNameAndNamespace(linkFileName, c)

}

func getAppNameAndNamespace(linkFileName string, c *cache.Cache)  {
	namespaceChan := make(chan string, 1)
	appNameChan := make(chan string, 1)
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
	// 判断是否已经监听了这个文件
	isMonitor := false

	if _, ok := fileMap.Load(linkFileName); ok {
		isMonitor = true
	}
	fileMap.LoadOrStore(linkFileName, 1)
	namespaceChan <- namespace
	appNameChan <- appName
	doLogAppender(namespace, appName, linkFileName, c)

	//
	//select {
	//	case appNameTmp := <- appNameChan:
	//		namespaceTmp := <-namespaceChan
	//		doLogAppender(namespaceTmp, appNameTmp, linkFileName, c)
	//
	//}
	// 追加日志
	if !isMonitor {

	}
}

func doLogAppender(namespace, appName, linkFileName string, c *cache.Cache)  {

	// 通过连接文件名获取到真正的文件信息 /var/lib/docker/containers/cid/cid-json.log
	containerId := common.GetRealFileName(linkFileName)
	if "" == containerId {
		return
	}
	// containerId workorder-web-9076b2e17db484c00468866f7a44274a58fe8b2d46be5ea66c6d78d95188a169
	//containerId = strings.ReplaceAll(containerId, appName+"-", "")

	zapLog.LOGGER().Debug("namespace: "+ namespace + ", appName:" + appName)

	// 获取真实的docker日志文件信息
	dockerLogFileName := *global.Docker_Log_Dir + containerId + string(filepath.Separator) +containerId + "-json.log"
	zapLog.LOGGER().Debug("dockerLogFileName: "+ dockerLogFileName)
	logappender.LogAppender(namespace, appName, dockerLogFileName, c)

}

//ListenAppLogV3 doc
//@Description: 监听k8s中应用容器的日志文件变化信息。当有追加的指令时，就做相关操作
//@Author niejian
//@Date 2021-04-28 14:17:12
//@param logDir
//@param appName
func ListenAppLogV3(c *cache.Cache) {
	// 缓存
	var fileMap sync.Map
	var namespaceMap sync.Map
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
				//zapLog.LOGGER().Info("fileName:" + fileName)
				fmt.Printf("============ \n")
				fmt.Printf("fileName: %v, Evnet: %v \n", fileName, event )
				fmt.Printf("============ \n")



				symlinks, _ := filepath.EvalSymlinks(fileName)
				fmt.Println("---->" + symlinks)

				// 写操作
				if op&fsnotify.Write == fsnotify.Write {
					var namespace string
					var appName string
					var key = common.Md5Str(fileName)
					// 缓存namespace等信息
					if val, ok := namespaceMap.Load(key); ok {
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
							namespaceMap.Store(key, namespace + "," + appName)
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
					zapLog.LOGGER().Info("文件创建操作：", zap.String("fileName", fileName))
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


