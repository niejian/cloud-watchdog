FROM alpine:latest
WORKDIR /data/watchdog/
ADD watchdog-cloud.tar.gz ./
CMD ["/bin/sh", "/data/watchdog/watchdog-cloud-bin/startDocker.sh"]
#  docker build -f Dockerfile -t k8s-watchdog:v1 .
#  docker run -it -v /Users/a/.kube:/root/.kube ede5cd8b4083