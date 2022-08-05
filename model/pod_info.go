// model doc

package model

// PodInfo pod信息实体描述,如果有需要,就增加字段
type PodInfo struct {
	Ip string `json:"ip"`
	AppName string `json:"appName"`
	Namespace string `json:"namespace"`
	PodName string `json:"podName"`
}
