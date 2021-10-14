.PHONY: all build run gotool clean zip done help

BINARY="watchdog-cloud"
TRA_NAME="watchdog-cloud.tar.gz"
TAR_DIR="./watchdog-cloud-bin/"
TAR_RESOURCES="./watchdog-cloud-bin/resources/"

all: clean gotool build copy zip done

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BINARY}

run:
	@go run ./

gotool:
	go mod tidy
	go fmt ./
	go vet ./

clean:
	@echo "clean 清理工作空间"
	@if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
	@if [ -d ${TAR_DIR} ]; then rm -rf ${TAR_DIR}; fi
	@if [ -d ${TAR_TAR} ]; then rm -rf ${TAR_TAR}; fi
copy:
	@echo "复制脚本"
	mkdir -p ${TAR_RESOURCES} && cp ./resources/*.sh ${TAR_DIR}
zip:
	@echo "打压缩包"
	cp -r ./resources/*.y*ml ${TAR_RESOURCES} && cp ${BINARY} ${TAR_DIR} && tar zcvf ${TRA_NAME} ${TAR_DIR}
done:
	@if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
	@if [ -d ${TAR_DIR} ]; then rm -rf ${TAR_DIR}; fi

help:
	@echo "make - 格式化 Go 代码, 并编译生成二进制文件"
	@echo "make build - 编译 Go 代码, 生成二进制文件"
	@echo "make run - 直接运行 Go 代码"
	@echo "make clean - 移除二进制文件和 vim swap files"
	@echo "make gotool - 运行 Go 工具 'fmt' and 'vet'"
	@echo "make zip - 打压缩包"
