module cloud-watchdog

go 1.14

require (
	github.com/apache/dubbo-go v1.5.6
	github.com/apache/dubbo-go-hessian2 v1.9.1
	github.com/fsnotify/fsnotify v1.4.9
	github.com/hpcloud/tail v1.0.0
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/spf13/viper v1.7.1
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.16.0
	golang.org/x/sys v0.0.0-20210514084401-e8d321eab015 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/driver/mysql v1.0.6
	gorm.io/gorm v1.21.9
	k8s.io/api v0.20.0
	k8s.io/apimachinery v0.20.0
	k8s.io/client-go v0.20.0
)
