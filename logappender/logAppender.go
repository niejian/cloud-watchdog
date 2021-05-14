// logappender doc

package logappender

import (
	"cloud-watchdog/common"
	"cloud-watchdog/zapLog"
	"fmt"
	"github.com/hpcloud/tail"
	"go.uber.org/zap"
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

func LogAppender(namespace, appName, logFileName string) {
	// 获取告警配置信息

	zapLog.LOGGER().Debug("文件写操作：", zap.String("fileName", logFileName))


	tailLog, err := tailLogFile(logFileName)
	if nil != err {
		zapLog.LOGGER().Error("构造tail对象失败，", zap.String("err", err.Error()))
		return
	}

	// tail -f
	for line := range tailLog.Lines {
		// 获取该文件的配置信息
		conf, err := common.GetLogAlterConfByFileName(namespace, appName)
		if nil != err {
			zapLog.LOGGER().Error("err", zap.String("err", err.Error()))
			continue
		}
		if nil == conf {
			zapLog.LOGGER().Error("配置为空")
			continue
		}
		text := line.Text
		fmt.Println("========》追踪到的日志信息：", text)
		zapLog.LOGGER().Info("追踪到的日志信息: " + text, zap.String("logFile", logFileName))
		zapLog.LOGGER().Debug("追踪到的日志信息: " + text)
		common.SendMsgUtil(text, conf)

		// 测试主动触发panic
		//if time.Now().Nanosecond() % 2 == 0 {
		//	panic("主动触发panic")
		//}
	}

	// 等待上游发送消息至此管道，如果有值，就停止阻塞
	//<-quitChan
}

