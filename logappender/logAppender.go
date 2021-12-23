// logappender doc

package logappender

import (
	"bytes"
	"cloud-watchdog/common"
	"cloud-watchdog/config"
	"cloud-watchdog/es"
	"cloud-watchdog/global"
	"cloud-watchdog/model"
	"cloud-watchdog/zapLog"
	"github.com/hpcloud/tail"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
	"strings"
	"sync"
	"time"
)

var (
	lock sync.Mutex
	tailMap sync.Map
)



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

func LogAppender(namespace, appName, logFileName string, c *cache.Cache) {
	zapLog.LOGGER().Info("文件写操作：", zap.String("fileName", logFileName))
	//quitChan := make(chan bool, 1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				zapLog.LOGGER().Error("tail log panic, prepare recovery", zap.Any("err", err))
				time.Sleep(1 * time.Second)
				tailLog(logFileName, namespace, appName, c)
			}
		}()

		tailLog(logFileName, namespace, appName, c)
	}()

	// 等待上游发送消息至此管道，如果有值，就停止阻塞
	//<-quitChan
}

func tailLog(logFileName, namespace, appName string, c *cache.Cache)  {

	tailLog, err := tailLogFile(logFileName)
	if nil != err {
		zapLog.LOGGER().Error("构造tail对象失败，", zap.String("err", err.Error()))
		return
	}

	tailMap.LoadOrStore(logFileName, tailLog)

	msg := ""
	// tail -f
	for line := range tailLog.Lines {
		custErr := ""
		// 获取该文件的配置信息
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

		// 是否开启告警
		if conf.IsEnable != 1 {
			zapLog.LOGGER().Info("未开启告警", zap.String("namespace", namespace), zap.String("appName", appName))
			continue
		}

		text := line.Text
		//fmt.Println("========》追踪到的日志信息：", text)
		zapLog.LOGGER().Debug("追踪到的日志信息: " + text, zap.String("logFile", logFileName))
		//k8s 日志
		var k8sLogModel model.K8sLogModel
		data, err := common.JSONStringFormat(text, k8sLogModel)
		if err != nil {
			zapLog.LOGGER().Error("json转换失败", zap.String("err", err.Error()))
			continue
		}


		//zapLog.LOGGER().Debug("追踪到的日志信息: " + text)

		msg = convertErrMsg(data.Log, msg, conf)

		ignores := conf.Ignores
		errs := conf.Errs
		hasExp := false
		//custErr := ""
		zapLog.LOGGER().Debug("", zap.String("msg", msg))



		// 1秒没操作，判断是需要发送消息
		time.AfterFunc(1*time.Second, func() {
			if "" == msg {
				return
			}
			//fmt.Println("时间静止500MS")
			msg = strings.TrimSpace(msg)
			if "" != msg && len(ignores) > 0 {
				// 判断是否包含忽略异常
				for _, ignoreEx := range ignores {

					if "" != ignoreEx && strings.Contains(msg, ignoreEx) {
						// 当前msg包含忽略异常关键字，将msg置空
						msg = ""
						zapLog.LOGGER().Info("当前消息包含忽略异常, 不予发送 ",  zap.String("ignore", ignoreEx))
						break
					}
				}
			}
			if "" != msg {
				for _, errTag := range errs {
					//fmt.Printf("errTag：%v, newLine: %v \n", errTag, newLine)
					// 含有异常关键字，发送提示告警
					if "" != errTag && strings.Contains(msg, errTag) {
						if errTag == "Exception"{
							// 拿到具体的异常信息
							index := strings.Index(msg, errTag+":")
							if index > 0 && len(msg) > 0{
								s := msg[0:index]
								if len(s) > 0 {
									split := strings.Split(s, ".")
									length := len(split)
									custErr = split[length - 1] + errTag
								}
							}

						}else {
							custErr = errTag
						}

						zapLog.LOGGER().Debug("has error", zap.String("err", errTag))
						hasExp = true
						break
					}
				}
				if hasExp && "" != msg &&  "" != custErr{
					lock.Lock()
					zapLog.LOGGER().Debug("", zap.String("msg", msg))
					isIgnore := isIgnoreMsg(msg, conf)
					if isIgnore {
						msg = ""
					}
					md5Str := ""
					isExist := true
					if "" != msg {
						md5Str = common.Md5Str(msg)
						_, isExist = c.Get(md5Str)
					}

					if !isExist && "" != msg &&  "" != custErr{
						c.Set(md5Str, "a", cache.DefaultExpiration)
						alarmMsg := convertWxchatMsg(custErr, appName, msg)
						common.SendMsgUtil(alarmMsg, conf)
						if conf.EnableStore == 1 {
							// 存储至es中
							vo := common.Convert2EsStore(appName, custErr, msg)
							// 插入数据
							es.InsertDocument(global.Es_Client, vo)
						}

					}
					lock.Unlock()
				}else {
					zapLog.LOGGER().Info("已发送该条告警消息")
				}
			}

			msg = ""

		})

	}
}

/*是否为忽略消息*/
func isIgnoreMsg(msg string, c *config.AlterConf) bool {
	ignores := c.Ignores
	isIgnore := false
	if "" != msg && len(ignores) > 0 {
		// 判断是否包含忽略异常
		for _, ignoreEx := range ignores {

			if "" != ignoreEx && strings.Contains(msg, ignoreEx) {
				// 当前msg包含忽略异常关键字，将msg置空
				msg = ""
				zapLog.LOGGER().Info("当前消息包含忽略异常, 不予发送 ",  zap.String("ignore", ignoreEx))
				return true
			}
		}
	}

	return  isIgnore
}

func convertWxchatMsg(custErr, appName, msg string) string {
	var buffer bytes.Buffer
	buffer.WriteString("应用名称：" + appName + "\n")
	buffer.WriteString("时间：" + common.FormatDate("2006-01-02 15:04:05") + "\n")
	buffer.WriteString("异常信息：" + custErr + "\n")
	buffer.WriteString(msg)
	alarmMsg := buffer.String()

	return alarmMsg

}

//convertErrMsg doc
//@Description: 组合告警消息
//@Author niejian
//@Date 2021-05-17 09:12:28
//@param text
//@param msg
//@param conf
func convertErrMsg(text, msg string, conf *config.AlterConf) string  {
	errs := conf.Errs
	newLine := text

	match := common.IsDatePrefix(text)
	hasExp := false
	for _, errTag := range errs {
		//fmt.Printf("errTag：%v, newLine: %v \n", errTag, newLine)
		// 含有异常关键字, 发送提示告警
		if strings.Contains(newLine, errTag) {
			hasExp = true
			break
		}

	}

	// log.error 输出方式，
	if !hasExp && match && strings.Contains(newLine, common.ERROR_TAG) {
		msg += newLine
	}

	if hasExp && match {
		if !strings.Contains(newLine, common.DEBUG_TAG) && !strings.Contains(newLine, common.WARN_TAG) {
			msg += newLine
		}
	}

	if hasExp && !match {
		//errContentChan <- newLine + "\n"
		//writing <- true
		msg += newLine
	}

	if !hasExp && !match && !strings.Contains(newLine, common.DEBUG_TAG) {
		//errContentChan <- newLine + "\n"
		//writing <- true
		msg += newLine
	}

	return msg
}

func StopTail(logFileName string)  {
	// 从map中获取
	tailLog, ok := tailMap.Load(logFileName)
	if ok {
		re, _ := tailLog.(*tail.Tail)
		re.Stop()
	}
}

