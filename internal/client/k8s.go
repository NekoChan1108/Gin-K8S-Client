package client

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/klog/v2"
	"path/filepath"
)

//TODO修改为单例模式

func GetK8sClientSet() (*kubernetes.Clientset, error) {
	if config, err := GetK8sConfig(); err != nil {
		klog.Fatal("k8s.GetK8sClientSet.GetK8SConfig err: ", err.Error())
		return &kubernetes.Clientset{}, err
	} else {
		if clientSet, err := kubernetes.NewForConfig(config); err != nil {
			klog.Fatal("k8s.GetK8sClientSet.NewForConfig err: ", err.Error())
			return &kubernetes.Clientset{}, err
		} else {
			return clientSet, nil
		}
	}
}

// GetK8sConfig 获取config
func GetK8sConfig() (*rest.Config, error) {
	var kubeConfig, homeDir string
	if homeDir = homedir.HomeDir(); homeDir != "" {
		kubeConfig = filepath.Join(homeDir, ".kube", "config")
	}
	if config, err := clientcmd.BuildConfigFromFlags("", kubeConfig); err != nil {
		klog.Fatal("k8s.GetK8sConfig.BuildConfigFromFlags err: ", err.Error())
		return &rest.Config{}, err
	} else {
		return config, nil
	}
}
