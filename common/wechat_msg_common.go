// common doc

package common

import (
	"bytes"
	conf "cloud-watchdog/config"
	"cloud-watchdog/config/parser"
	"cloud-watchdog/rpc"
	"cloud-watchdog/zapLog"
	"go.uber.org/zap"
	"os"
)

func SendMsgUtil(msg string, conf *conf.AlterConf)  {
	sysConfigParser, err := parser.SysConfigParser()
	if err != nil {
		zapLog.LOGGER().Error("配置转换失败", zap.String("err", err.Error()))
		return
	}
	toUserIds := conf.ToUserIds
	// 发送微信消息
	responseBean, err := rpc.SendMsg(msg, toUserIds, &sysConfigParser.WxChatMsgConf)
	if err != nil {
		zapLog.LOGGER().Error("调用dubbo接口失败" + err.Error())
	} else {
		zapLog.LOGGER().Info("返回结果:" +  responseBean.String())
		// 发送失败原因
		if !responseBean.IsSuccess {
			// 通知管理员检查为什么没有将消息发出去
			checkUserIds := os.Getenv("msg.checkers")
			if "" == checkUserIds {
				checkUserIds = "80487083|80468295"
			}

			var buffer bytes.Buffer
			buffer.WriteString("消息发送失败,请检查!\n")
			buffer.WriteString("responseMsg：" + responseBean.ResponseMsg + "\n")
			buffer.WriteString("userIds：" + toUserIds + "\n")
			rpc.SendMsg(buffer.String(), checkUserIds, &sysConfigParser.WxChatMsgConf)

		}

	}


}
