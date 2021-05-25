// global doc

package global

import (
	"gorm.io/gorm"
)

var (
	GLOBAL_DB      *gorm.DB
	//LOG_ALTER_NAME = "log_alter_conf"
	LOG_ALTER_NAME *string
	//K8S_LOG_DIR    = "/var/log/containers"
	//K8S_LOG_DIR    = "/Users/a/logs/"
	K8S_LOG_DIR  *string
	Docker_Log_Dir *string
	Exclude_Ns *[]string
)
