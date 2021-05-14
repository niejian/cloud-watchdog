// global doc

package global

import "gorm.io/gorm"

var (
	GLOBAL_DB      *gorm.DB
	LOG_ALTER_NAME = "log_alter_conf"
	//K8S_LOG_DIR    = "/var/log/containers"
	K8S_LOG_DIR    = "/Users/a/logs/"
)
