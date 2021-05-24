// api doc

package api

import (
	"cloud-watchdog/zapLog"
	"context"
	"flag"
	"fmt"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

var clientSet *kubernetes.Clientset

func init() {
	var kubeConfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeConfig = flag.String("kubeConfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeConfig file")
	} else {
		kubeConfig = flag.String("kubeConfig", "", "absolute path to the kubeConfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeConfig)

	clientSetInit, err := kubernetes.NewForConfig(config)

	if nil != clientSetInit {
		clientSet = clientSetInit
	} else {
		panic("k8s connect failed")
	}
	if err != nil {
		zapLog.LOGGER().Error("k8s connect failed", zap.String("err", err.Error()))
	} else {
		zapLog.LOGGER().Info("connect k8s success")
	}
}

func InitK8s() *kubernetes.Clientset {
	return clientSet
}

//DescribePod doc
//@Description: 获取pod的详细信息
//@Author niejian
//@Date 2021-05-08 11:33:27
//@param podName
//@param ns
//@return *v1.Pod
//@return error
func DescribePod(podName, ns string) (*v1.Pod, error) {
	return clientSet.CoreV1().Pods(ns).Get(context.TODO(), podName, metav1.GetOptions{})
}

//func DescribeDeploy(deployName, ns string, labels map[string]string) string {
//	deploy, _ := clientSet.AppsV1().Deployments(ns).Get(context.TODO(), deployName, metav1.GetOptions{})
//	// deploy文件中指定的label信息
//	podLabels := deploy.Spec.Template.Labels
//	for key, val := range labels {
//		// podLabels必须全部包含
//		data, isExist := podLabels[key]
//		if  !isExist {
//			break
//			return ""
//		}
//		if val != data {
//			break
//			return ""
//		}
//	}
//	return de
//}

func ListNodes()  {
	nodes, err := clientSet.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		zapLog.LOGGER().Error("获取k8s节点失败", zap.String("err", err.Error()))
	}
	for _, node := range nodes.Items {
		addresses := node.Status.Addresses
		//fmt.Printf("node addresses: %v \n", addresses)

		//isWorkerNode := false
		var addressType v1.NodeAddressType
		var addressName string

		for _, address := range addresses {
			addressType = address.Type
			addressName = address.Address

			if v1.NodeInternalIP == addressType {
				fmt.Printf("node ip: %v \n", addressName)
			}
		}
	}
}
