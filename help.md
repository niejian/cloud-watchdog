
```shell

docker build -f Dockerfile -t k8s-watchdog:v8 .

docker login harbor.bluemoon.com.cn

docker tag k8s-watchdog:v8 harbor.bluemoon.com.cn/cloud-platform/k8s-watchdog:v8

docker push harbor.bluemoon.com.cn/cloud-platform/k8s-watchdog:v8

```