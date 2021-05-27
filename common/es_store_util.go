// common doc

package common

import (
	"cloud-watchdog/model"
	"fmt"
	"time"
)

//Convert2EsStore doc
//@Description: es存储实体
//@Author niejian
//@Date 2021-05-26 16:46:46
//@param appName
//@param exceptionTag
//@param msg
//@return *model.ExceptionStore
func Convert2EsStore(appName, exceptionTag, msg string) *model.ExceptionStore {
	vo := &model.ExceptionStore{
		Id:           fmt.Sprintf("%d", time.Now().UnixNano()),
		//Ip:           GetNetIp(),
		AppName:      appName,
		CreateTime:   time.Now().UnixNano() / 1e6,
		ExceptionTag: exceptionTag,
		From:         appName,
		Msg:          msg,
		Year:         time.Now().Format("2006"),
		Month:        time.Now().Format("01"),
		Day:          time.Now().Format("02"),
		CreateDate:   time.Now().Format("20060102"),
	}

	return vo
}
