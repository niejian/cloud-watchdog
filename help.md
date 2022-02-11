
```shell

docker build -f Dockerfile -t k8s-watchdog:v7 .

docker login harbor.bluemoon.com.cn

docker tag k8s-watchdog:v7 harbor.bluemoon.com.cn/cloud-platform/k8s-watchdog:v7

docker push harbor.bluemoon.com.cn/cloud-platform/k8s-watchdog:v7

```