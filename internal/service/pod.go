package service

import (
	"Gin-K8S-Client/internal/client"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

type PodSvc struct {
}

var (
	podSvcInstance *PodSvc
)

// NewPodSvc 单例加载
func NewPodSvc() *PodSvc {
	once.Do(func() {
		podSvcInstance = &PodSvc{}
	})
	return podSvcInstance
}

func (ps *PodSvc) GetPod() ([]v1.Pod, error) {
	clientSet, err := client.GetK8sClientSet()
	if err != nil {
		klog.Error("service.GetPod.GetK8sClientSet err: ", err.Error())
		return []v1.Pod{}, err
	}
	podList, err := clientSet.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		klog.Error("service.GetPod.Namespaces().List err: ", err.Error())
		return []v1.Pod{}, err
	}
	return podList.Items, nil
}
