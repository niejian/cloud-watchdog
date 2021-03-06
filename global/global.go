// global doc

package global

import (
	"github.com/olivere/elastic/v7"
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
	Es_Client *elastic.Client



	GVA_INDICE_NAME_PREFIX = "watchdog_store_"
	GVA_INDICE_MAPPING     = `
	{
		"settings":{
			"number_of_shards":3,
			"number_of_replicas": 2
		},
		"mappings": {
			"properties":{
				"id":{
					"type": "long"
				},
				"year":{
					"type": "text",
					"fielddata": true
				},
				"month":{
					"type": "text",
					"fielddata": true
				},
				"day":{
					"type": "text",
					"fielddata": true
				},
				"createDate":{
					"type": "text",
					"fielddata": true
				},
				"ip":{
					"type": "text"
				},
				"exceptionTag": {
					"type": "text",
					"analyzer": "keyword",
					"fielddata": true
				},
				"createTime":{
					"type": "long"
				},
				"msg":{
					"type": "text"
				},
				"from":{
					"type": "text"
				},
				"appName": {
					"type": "text",
					"analyzer": "keyword"
				}
			}
		}
	}`
)
