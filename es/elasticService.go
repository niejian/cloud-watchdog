// es doc

package es

import (
	"cloud-watchdog/common"
	"cloud-watchdog/config"
	"cloud-watchdog/global"
	"cloud-watchdog/model"
	"cloud-watchdog/zapLog"
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func InitEs(es *config.Es) (*elastic.Client, error) {
	urls := es.Urls
	if nil == urls || len(urls) == 0 {
		zapLog.LOGGER().Error("es 地址为空")
		return nil, nil
	}
	username := es.Username
	password := es.Password

	var client *elastic.Client

	if "" != username && "" != password {
		client1, err := elastic.NewClient(
			elastic.SetURL(urls...),
			elastic.SetBasicAuth(username, password),
			elastic.SetSniff(true),
			elastic.SetHealthcheckInterval(10*time.Second),
			elastic.SetGzip(true),
			elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC-ERR-", log.LstdFlags)),
			elastic.SetInfoLog(log.New(os.Stdout, "ELASTIC-INFO-", log.LstdFlags)),
			elastic.SetHeaders(http.Header{
				"X-Caller-Id": []string{"..."},
			}),
		)

		if nil != err {
			zapLog.LOGGER().Error("es连接失败，退出....err:  ", zap.String("err", err.Error()))
			return nil, err
		}
		client = client1

	} else {
		client1, err := elastic.NewClient(
			elastic.SetURL(urls...),
			elastic.SetSniff(true),
			elastic.SetHealthcheckInterval(10*time.Second),
			elastic.SetGzip(true),
			elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC-ERR-", log.LstdFlags)),
			elastic.SetInfoLog(log.New(os.Stdout, "ELASTIC-INFO-", log.LstdFlags)),
			elastic.SetHeaders(http.Header{
				"X-Caller-Id": []string{"..."},
			}),
		)

		if nil != err {
			zapLog.LOGGER().Error("es连接失败，退出....err:  ", zap.String("err", err.Error()))
			return nil, err
		}
		client = client1
	}
	return client, nil
}


// 判断当前索引是否存在，不存在则创建之
// client esclient 客户端
// appName 应用名称
func getIndice(client *elastic.Client, appName string) string {
	appName = strings.ReplaceAll(appName, "-", "_")
	indexName := fmt.Sprintf("%s%s", global.GVA_INDICE_NAME_PREFIX, appName)

	exists, err := client.IndexExists(indexName).Do(context.Background())
	if err != nil {
		zapLog.LOGGER().Error("es 连接失败", zap.String("err", err.Error()))

	}

	if !exists {
		zapLog.LOGGER().Info("索引：[" + indexName + "], 不存在")
		// 创建索引
		createIndice(client, indexName, "")
	}

	return indexName
}

// 创建索引
func createIndice(client *elastic.Client, indexName, mappingName string) {
	var indicesCreateResult *elastic.IndicesCreateResult
	var err error

	if "" == mappingName {
		indicesCreateResult, err = client.CreateIndex(indexName).BodyString(global.GVA_INDICE_MAPPING).Do(context.Background())
	} else {
		indicesCreateResult, err = client.CreateIndex(indexName).BodyString(mappingName).Do(context.Background())

	}
	if nil != err {
		zapLog.LOGGER().Error("索引创建失败", zap.String("indexName", indexName), zap.String("err", err.Error()))

		//log.Fatalf("创建索引：%v 失败：%v", indexName, err)
	}

	if !indicesCreateResult.Acknowledged {
		zapLog.LOGGER().Error("索引创建失败", zap.String("indexName", indexName), zap.String("err", err.Error()))
	} else {
		zapLog.LOGGER().Info("索引创建成功", zap.String("indexName", indexName))
	}
}

//InsertDocument doc
//@Description: 插入数据
//@Author niejian
//@Date 2021-05-26 16:38:51
//@param client
//@param vo
//@return error
func InsertDocument(client *elastic.Client, vo *model.ExceptionStore) error {
	indexName := getIndice(client, vo.AppName)
	ctx := context.Background()
	create, err := client.Index().Index(indexName).
		Id(fmt.Sprintf("%d%s", time.Now().UnixNano(), "")).
		BodyJson(vo).Do(ctx)
	if err != nil {
		zapLog.LOGGER().Error("添加数据失败 %v ", zap.String("err", err.Error()))
		return err
	}
	zapLog.LOGGER().Info("数据添加成功", zap.String("Indexed", indexName + " " +create.Id),
	zap.String("index", create.Index), zap.String("type", create.Type))
	log.Printf("Indexed %v %s to index %s, type %s\n", indexName, create.Id, create.Index, create.Type)
	return nil
}

//createIndiceByDate doc
//@Description: 按天创建索引
//@Author niejian
//@Date 2021-10-09 13:37:42
//@param client
//@param indexName
func createIndiceByDate(client *elastic.Client, indexName string)  {
	date := common.FormatDate(common.TIME_DATE_FORMAT)
	indexName = fmt.Sprintf("%s_%s", indexName, date)
	indicesCreateResult, err := client.CreateIndex(indexName).BodyString(global.LOG_COLLECT_INDICE_MAPPING).Do(context.Background())
	if nil != err {
		zapLog.LOGGER().Error("索引创建失败", zap.String("indexName", indexName), zap.String("err", err.Error()))

		//log.Fatalf("创建索引：%v 失败：%v", indexName, err)
	}

	if !indicesCreateResult.Acknowledged {
		zapLog.LOGGER().Error("索引创建失败", zap.String("indexName", indexName), zap.String("err", err.Error()))
	} else {
		zapLog.LOGGER().Info("索引创建成功", zap.String("indexName", indexName))
	}
}

//InsertLogContent doc
//@Description: 插入日志
//@Author niejian
//@Date 2021-10-09 13:50:01
//@param client
//@param vo
//@param appName
//@return error
func InsertLogContent(client *elastic.Client, vo *model.LogContent, appName string) error {
	indexName := getDateIndice(client, appName)
	ctx := context.Background()
	create, err := client.Index().
		Index(indexName).
		BodyJson(vo).
		Do(ctx)
	if err != nil {
		zapLog.LOGGER().Error("添加数据失败", zap.String("err", err.Error()))
		return err
	}
	zapLog.LOGGER().Info("数据添加成功", zap.String("Indexed", indexName + " " +create.Id),
		zap.String("index", create.Index), zap.String("type", create.Type))
	//log.Printf("Indexed %v %s to index %s, type %s\n", indexName, create.Id, create.Index, create.Type)
	return nil
}

//getDateIndice doc
//@Description: 判断日志索引是否存在,如果不存在则创建
//@Author niejian
//@Date 2021-10-09 13:50:13
//@param client
//@param appName
//@return string
func getDateIndice(client *elastic.Client, appName string) string {
	appName = strings.ReplaceAll(appName, "-", "_")
	// appName 添加日期
	date := common.FormatDate(common.TIME_DATE_FORMAT)
	indexName := fmt.Sprintf("%s%s%s%s", global.LOG_COLLECT_INDICE_NAME_PREFIX, appName, "_", date)

	exists, err := client.IndexExists(indexName).Do(context.Background())
	if err != nil {
		zapLog.LOGGER().Error("es 连接失败", zap.String("err", err.Error()))
	}

	if !exists {
		zapLog.LOGGER().Info("索引：[" + indexName + "], 不存在")
		// 创建索引
		createIndice(client, indexName, global.LOG_COLLECT_INDICE_MAPPING)
	}
	return indexName
}

