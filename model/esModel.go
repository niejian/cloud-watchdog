// model doc

package model

type EsModel struct {
	Urls     []string `yaml:"urls" json:"urls" mapstructure:"urls"`
	Username string   `yaml:"username" json:"username" mapstructure:"username"`
	Password string   `yaml:"password" json:"password" mapstructure:"password"`
}

type ExceptionStore struct {
	Id           string `json:"id"`
	Ip           string `json:"ip"`
	Year         string `json:"year"`
	Month        string `json:"month"`
	Day          string `json:"day"`
	CreateDate   string `json:"createDate"`
	AppName      string `json:"appName"`
	CreateTime   int64  `json:"createTime"`
	ExceptionTag string `json:"exceptionTag"`
	From         string `json:"from"`
	Msg          string `json:"msg"`
}
