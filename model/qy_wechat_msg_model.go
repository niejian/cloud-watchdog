// 企业微信消息实体 doc

package model


type MsgText struct {
	Content string `json:"content"`
}

type MsgData struct {
	Touser string `json:"touser"`
	MsgType string `json:"msgtype"`
	Agentid int32 `json:"agentid"`
	Text interface{} `json:"text"`
}

type Msg struct {
	CorpId string `json:"corpId"`
	Agentid int32 `json:"agentId"`
	Data interface{} `json:"data"`
}
