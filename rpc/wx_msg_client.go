// rpc doc

package rpc

import (
	conf "cloud-watchdog/config"
	"cloud-watchdog/rpc/pkg"
	"cloud-watchdog/zapLog"
	"encoding/json"
	"github.com/apache/dubbo-go/config"
	"go.uber.org/zap"
)

import (
	hessian "github.com/apache/dubbo-go-hessian2"
	_ "github.com/apache/dubbo-go/cluster/cluster_impl"
	_ "github.com/apache/dubbo-go/cluster/loadbalance"
	_ "github.com/apache/dubbo-go/common/proxy/proxy_factory"
	_ "github.com/apache/dubbo-go/filter/filter_impl"
	_ "github.com/apache/dubbo-go/protocol/dubbo"
	_ "github.com/apache/dubbo-go/registry/protocol"
	_ "github.com/apache/dubbo-go/registry/zookeeper"
)

var (
	wechatMsgProvider *pkg.WechatMsgProvider
	Msg_Type = "text"
)
func init()  {
	wechatMsgProvider = new(pkg.WechatMsgProvider)
	hessian.RegisterPOJO(&pkg.AgentMsgVo{})
	hessian.RegisterPOJO(&pkg.PostWechatMsg{})
	hessian.RegisterPOJO(&pkg.WechatMsgSendInfo{})
	hessian.RegisterPOJO(&pkg.ResponseBean{})
	// 注册引用信息
	config.SetConsumerService(wechatMsgProvider)
	config.Load()
}

//SendMsg doc
//@Description: 调用dubbo接口发送微信消息
//@Author niejian
//@Date 2021-05-13 09:07:39
//@param msg
//@param wxChatMsgConf
//@return *pkg.ResponseBean
//@return error
func SendMsg(msg, toUserIds string, wxChatMsgConf *conf.WxChatMsgConf) (*pkg.ResponseBean, error) {
	zapLog.LOGGER().Info("开始调用dubbo接口，", zap.String("method", "SendAgentMsg"))
	// 组装消息信息
	agentId := wxChatMsgConf.AgentId
	corpId := wxChatMsgConf.CorpId

	postWechatMsg := &pkg.PostWechatMsg{
		Content: msg,
	}

	wechatMsgSendInfo := pkg.WechatMsgSendInfo{
		ToUser:  toUserIds,
		MsgType: Msg_Type,
		AgentId: agentId,
		Text: postWechatMsg,
	}
	// 转json字符串
	wechatMsgSendInfoJSON, _ := json.Marshal(wechatMsgSendInfo)

	responseBean, err := wechatMsgProvider.SendAgentMsg(corpId, agentId, string(wechatMsgSendInfoJSON))
	zapLog.LOGGER().Info("调用dubbo接口结束", zap.String("method", "SendAgentMsg"))
	return responseBean, err
}
