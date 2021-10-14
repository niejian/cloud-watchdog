
```shell

docker build -f Dockerfile -t k8s-watchdog:v1 .
docker rmi -f $(docker images|grep none |awk '{print $3}')
docker run -it -v /Users/a/.kube:/root/.kube ede5cd8b4083

docker login harbor.bluemoon.com.cn

docker tag k8s-watchdog:v1 harbor.bluemoon.com.cn/cloud-platform/k8s-watchdog:v1

docker push harbor.bluemoon.com.cn/cloud-platform/k8s-watchdog:v1

v4
```