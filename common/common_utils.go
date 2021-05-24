// common doc

package common

import (
	"cloud-watchdog/model"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"regexp"
	"time"
)

var (
	ERROR_TAG        = "ERROR"
	DEBUG_TAG        = "DEBUG"
	WARN_TAG         = "WARN"
)
//JSONStringFormat doc
//@Description: 将json字符串转化为对应实体信息，失败返回(nil, error)
//@Author niejian
//@Date 2021-05-17 09:33:47
//@param s
//@param dto
//@return interface{}
//@return error
func JSONStringFormat(s string, dto model.K8sLogModel) (model.K8sLogModel, error) {
	err := json.Unmarshal([]byte(s), &dto)

	return dto, err
}


//IsDatePrefix doc
//@Description: 判断字符串是否是日期时间戳开头
//@Author niejian
//@Date 2021-05-17 10:17:47
//@param line
//@return bool
func IsDatePrefix(line string) bool {
	r := []rune(line)
	newLine20Prefix := string(r[0:19])
	pattern := "\\d{4}\\-\\d{2}\\-\\d{2}\\s\\d{2}:\\d{2}:\\d{2}"
	match, _ := regexp.Match(pattern, []byte(newLine20Prefix))
	return match
}

//Md5Str doc
//@Description: md5摘要
//@Author niejian
//@Date 2021-05-17 13:46:05
//@param str
//@return string
func Md5Str(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

//FormatDate doc
//@Description: 时间格式化
//@Author niejian
//@Date 2021-05-24 15:52:59
//@param pattern
//@return string
func FormatDate(pattern string) string {
	return time.Now().Format(pattern)
}