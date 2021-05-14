// pkg doc

package pkg

import (
	"fmt"
)

//AgentMsgVo doc
//@Description: 企业微信消息结构体
//@Author niejian
type AgentMsgVo struct {
	ErrCode     int    `json:"errcode"` // 0 成功
	ErrMsg      string `json:"errmsg"`  // ok 成功
	InvalidUser string `json:"invaliduser"`
}

//name doc
//@Description:  重写tostring
//@Author niejian
//@Date 2021-05-12 17:09:07
//@receiver w
//@return string
func (w AgentMsgVo) String() string {
	return fmt.Sprintf(
		"AgentMsgVo{ErrCode:%s, ErrMsg:%s, InvalidUser:%s}",
		w.ErrCode, w.ErrMsg, w.InvalidUser,
	)
}

// JavaClassName 声明java包名
func (AgentMsgVo) JavaClassName() string  {
	return "cn.com.bluemoon.qy.dubbo.api.domain.vo.AgentMsgVo"
}

type WechatMsgSendInfo struct {
	ToUser string `json:"touser"`
	MsgType string `json:"msgtype"`
	AgentId string `json:"agentid"`
	Text *PostWechatMsg `json:"text"`
}

func (w WechatMsgSendInfo) JavaClassName() string {
	return fmt.Sprintf(
		"WechatMsgSendInfo{MsgType:%s, Text:%s, ToUser:%s, AgentId:%s}",
		w.MsgType,
		w.Text,
		w.ToUser,
		w.AgentId,
	)
}

//PostWechatMsg doc
//@Description: 消息结构体
//@Author niejian
type PostWechatMsg struct {
	Content string `json:"content"`
}

func (p PostWechatMsg) JavaClassName() string {
	return fmt.Sprintf(
		"PostWechatMsg{content:%s}", p.Content,
	)
}


//WechatMsgProvider doc
//@Description: dubbo接口定义方法
//@Author niejian
type WechatMsgProvider struct {
	// 发送企业微信信息
	SendAgentMsg func(corpid, agentid ,postWechatMsg string) (*ResponseBean, error) `dubbo:"sendAgentMsg"`
}

func (wechatMsgProvider *WechatMsgProvider) Reference() string {
	return "WechatMsgProvider"
}

type ResponseBean struct {
	IsSuccess bool `json:"isSuccess"`
	ResponseCode int `json:"responseCode"`
	ResponseMsg string `json:"responseMsg"`
	Data interface{} `json:"data"`
}

func (r ResponseBean) String() string {
	return fmt.Sprintf(
		"ResponseBean{isSuccess:%v, responseCode:%s, responseMsg:%s, data:%v}",
		r.IsSuccess, r.ResponseCode, r.ResponseMsg, r.Data,
	)
}

func (ResponseBean) JavaClassName() string {
	return "com.bluemoon.pf.standard.bean.ResponseBean"
}
