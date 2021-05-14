// common doc

package common

import (
	conf "cloud-watchdog/config"
	"cloud-watchdog/config/parser"
	"cloud-watchdog/rpc"
	"cloud-watchdog/zapLog"
	"go.uber.org/zap"
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
	}


}
