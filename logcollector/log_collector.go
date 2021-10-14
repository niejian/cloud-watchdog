// logcollector doc

package logcollector

import (
	"bytes"
	"cloud-watchdog/common"
	"cloud-watchdog/es"
	"cloud-watchdog/global"
	"cloud-watchdog/model"
	"cloud-watchdog/zapLog"
	"github.com/hpcloud/tail"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	lock sync.Mutex
	logCollectorTailMap sync.Map
)

//CollectorLog doc
//@Description: 日志收集
//@Author niejian
//@Date 2021-10-09 11:06:53
//@param namespace
//@param appName
//@param linkFileName
//@param c
func CollectorLog(namespace, appName, linkFileName string, c *cache.Cache)  {

	podName := getPodNameByLinkFile(linkFileName)
	if "" == podName {
		zapLog.LOGGER().Error("获取不了podName", zap.String("file", linkFileName))
		return
	}
	// 通过连接文件名获取到真正的文件信息 /var/lib/docker/containers/cid/cid-json.log
	// 获取真实的docker日志文件信息
	zapLog.LOGGER().Info("开始收集日志", zap.String("appName", appName))
	dockerLogFileName := getDockerLogFilePath(linkFileName)
	if "" == dockerLogFileName {
		zapLog.LOGGER().Error("linkfile：" + linkFileName + ", 无法获取到实际文件信息")
		return
	}
	zapLog.LOGGER().Info("dockerLogFileName: "+ dockerLogFileName)

	// 开始tail log, 并收集log信息
	go func() {
		defer func() {
			if err := recover(); err != nil {
				zapLog.LOGGER().Error("tail log panic, prepare recovery", zap.Any("err", err))
				time.Sleep(1 * time.Second)
				doLogCollector(dockerLogFileName, appName, namespace, linkFileName)
			}
		}()

		doLogCollector(dockerLogFileName, appName, namespace, linkFileName)
	}()

}


//doLogCollector doc
//@Description: 收集日志动作
//@Author niejian
//@Date 2021-10-09 11:26:55
//@param logFileName
//@param appName
//@param namespace
//@param c
func doLogCollector(logFileName, appName, namespace , linkFileName string) {
	podName := getPodNameByLinkFile(linkFileName)
	tailLog, err := tailLogFile(logFileName)
	zapLog.LOGGER().Debug("doLogCollector 开始收集日志", zap.String("appName", appName), zap.String("logFileName", logFileName))
	if nil != err {
		zapLog.LOGGER().Error("构造tail对象失败，", zap.String("err", err.Error()))
		return
	}

	logCollectorTailMap.LoadOrStore(logFileName, tailLog)

	var buffer bytes.Buffer
	//var nextLine bytes.Buffer
	for line := range tailLog.Lines {
		// 查询数据库,是否配置日志收集
		conf, err := common.GetLogAlterConfByFileName(namespace, appName)
		if nil != err {
			zapLog.LOGGER().Debug("err", zap.String("namespace", namespace),
				zap.String("appname", appName), zap.String("err", err.Error()))
			continue
		}
		if nil == conf {
			zapLog.LOGGER().Debug("配置为空")
			continue
		}

		// 是否开启日志收集
		if conf.IsCollectLog != 1 {

			zapLog.LOGGER().Info("未开启日志收集", zap.String("namespace", namespace), zap.String("appName", appName), zap.Any("查询实体", conf))
			continue
		}

		text := line.Text
		//k8s 日志
		var k8sLogModel model.K8sLogModel
		data, err := common.JSONStringFormat(text, k8sLogModel)
		if err != nil {
			zapLog.LOGGER().Error("json转换失败", zap.String("err", err.Error()))
			continue
		}

		content := data.Log
		if "" == content || " " == content{
			continue
		}
		//判断日志是否是日期开头
		match := common.IsDatePrefix(content)
		buffer.WriteString(content)
		if !match {
			// 如果不是,就将他们拼成一行
			//buffer.WriteString("\n")
			continue
		}
		logContent := buffer.String()
		// 发送数据至elastic search
		//go func() {
		//	defer func() {
		//		if err := recover(); err != nil {
		//			zapLog.LOGGER().Error("saveLogContent panic, prepare recovery", zap.Any("err", err))
		//			time.Sleep(1 * time.Second)
		//			saveLogContent(logContent, appName, podName)
		//		}
		//	}()
		//	saveLogContent(logContent, appName, podName)
		//}()
		saveLogContent(logContent, appName, podName)
		// buffer置空
		buffer.Reset()

	}
}

//name doc
//@Description:
//@Author niejian
//@Date 2021-10-09 13:56:23
//@param content
//@param appName
func saveLogContent(content, appName, podName string) {
	if "" == content {
		zapLog.LOGGER().Error("日志数据为空", zap.String("podName", podName))
		return
	}
	vo := &model.LogContent{
		Message:    content,
		Timestamp: time.Now(),
		Version: "1",
		PodName: podName,
	}
	err := es.InsertLogContent(global.Es_Client, vo, appName)
	if err != nil {
		zapLog.LOGGER().Error("日志添加失败", zap.String("err", err.Error()))
		es.InsertLogContent(global.Es_Client, vo, appName)
	}
}

//tailLogFile doc
//@Description: 构造tail对象
//@Author niejian
//@Date 2021-04-28 16:13:30
//@param logFileName
//@return *tail.Tail
//@return error
func tailLogFile(logFileName string) (*tail.Tail, error) {
	return tail.TailFile(logFileName, tail.Config{
		ReOpen:    true,                                 // 文件被移除或被打包，需要重新打开
		Follow:    true,                                 // 实时跟踪
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // 如果程序出现异常，保存上次读取的位置，避免重新读取。 whence，从哪开始：0从头，1当前，2末尾
		MustExist: false,                                // 如果文件不存在，是否推出程序，false是不退出
		Poll:      true,
	})
}

func StopTail(logFileName string)  {
	// 从map中获取
	tailLog, ok := logCollectorTailMap.Load(logFileName)
	if ok {
		re, _ := tailLog.(*tail.Tail)
		re.Stop()
		// map 中删除
		logCollectorTailMap.Delete(logFileName)
	}
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

//getPodNameByLinkFile doc
//@Description: 通过连接文件获取pod名称,日志显示用
//@Author niejian
//@Date 2021-10-12 15:25:53
//@param linkFileName
//@return string
func getPodNameByLinkFile(linkFileName string) string {
	names := strings.Split(linkFileName, "_")
	var podName string
	if len(names) > 2 {
		podName = names[0]
		if strings.Contains(*global.K8S_LOG_DIR, podName) {
			podName = strings.ReplaceAll(*global.K8S_LOG_DIR, podName, "")
		}
	}
	return podName
}

