// initialize doc

package initialize

//DoInitialize doc
//@Description: 初始化
//@Author niejian
//@Date 2021-05-08 14:50:56
func DoInitialize(env string)  {
	// mysql 初始化
	MysqlInitialize()
	GlobalConstInitialize(env)
	
}
