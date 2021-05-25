// config doc

package config

//SysConfig doc
//@Description: 系统配置
//@Author niejian
type SysConfig struct {
	Mysql         Mysql         `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Log           LogConfig     `mapstructure:"log" json:"log" yaml:"log"`
	WxChatMsgConf WxChatMsgConf `mapstructure:"wx" json:"wx" yaml:"wx"`
	GlobalConst   GlobalConst   `mapstructure:"globalConst" json:"globalConst" yaml:"globalConst"`
	ExcludeNs []string `mapstructure:"exclude-ns" json:"exclude-ns" yaml:"exclude-ns"` // 排除解析的namespace
}

type Mysql struct {
	Username     string `mapstructure:"username" json:"username" yaml:"username"`
	Password     string `mapstructure:"password" json:"password" yaml:"password"`
	Path         string `mapstructure:"path" json:"path" yaml:"path"`
	Dbname       string `mapstructure:"db-name" json:"dbname" yaml:"db-name"`
	Config       string `mapstructure:"config" json:"config" yaml:"config"`
	MaxIdleConns int    `mapstructure:"max-idle-conns" json:"maxIdleConns" yaml:"max-idle-conns"`
	MaxOpenConns int    `mapstructure:"max-open-conns" json:"maxOpenConns" yaml:"max-open-conns"`
	LogMode      bool   `mapstructure:"log-mode" json:"logMode" yaml:"log-mode"`
}

type WxChatMsgConf struct {
	AgentId string `mapstructure:"agentid" json:"agentid" yaml:"agentid"`
	CorpId  string `mapstructure:"corpid" json:"corpid" yaml:"corpid"`
}

type AlterConf struct {
	ToUserIds string   `mapstructure:"toUserIds" json:"toUserIds" yaml:"toUserIds"`
	Ignores   []string `mapstructure:"ignores" json:"ignores" yaml:"ignores"`
	Errs      []string `mapstructure:"errs" json:"errs" yaml:"errs"`
	AppName   string   `mapstructure:"appName" json:"appName" yaml:"appName"`
	Namespace string   `mapstructure:"namespace" json:"namespace" yaml:"namespace"`
}

type GlobalConst struct {
	// 日志监听路径
	Dev          Dev    `mapstructure:"dev" json:"dev" yaml:"dev"`
	Pro          Pro    `mapstructure:"pro" json:"pro" yaml:"pro"`
	LogAlterName string `mapstructure:"logAlterName" json:"logAlterName" yaml:"logAlterName"`
	DockerLogDir string `mapstructure:"dockerLogDir" json:"dockerLogDir" yaml:"dockerLogDir"`
}

type Dev struct {
	K8sLogDir string `mapstructure:"k8sLogDir" json:"k8sLogDir" yaml:"k8sLogDir"`
}

type Pro struct {
	K8sLogDir string `mapstructure:"k8sLogDir" json:"k8sLogDir" yaml:"k8sLogDir"`
}
