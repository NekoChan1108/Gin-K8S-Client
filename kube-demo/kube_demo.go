package main

import (
	"context"
	"flag"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/klog/v2"
	"path/filepath"
)

func main() {
	var kubeConfig *string
	ctx := context.Background()
	//获取./kube/config文件
	if homeDir := homedir.HomeDir(); homeDir != "" {
		kubeConfig = flag.String("kubeConfig", filepath.Join(homeDir, ".kube", "config"), "absolute path to the kube config file")
	} else {
		kubeConfig = flag.String("kubeConfig", "", "no kube config file found")
	}
	//解析参数
	flag.Parse()
	//通过kubeConfig参数创建client
	config, err := clientcmd.BuildConfigFromFlags("", *kubeConfig)
	if err != nil {
		klog.Fatal(err)
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatal(err)
	}
	//拿到命名空间数组
	namespaceList, err := clientSet.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	//拿到pod数组
	podList, err := clientSet.CoreV1().Pods("default").List(ctx, metav1.ListOptions{})
	if err != nil {
		klog.Fatal(err)
	}
	namespaces := namespaceList.Items
	pods := podList.Items
	//遍历打印命名空间相关
	for _, item := range namespaces {
		fmt.Printf("name ====> %v\n  staus ===> %v\n", item.Name, item.Status)
	}
	//遍历打印pod相关
	for _, item := range pods {
		fmt.Printf("pod ====> %v\n", item.Name)
	}
}
