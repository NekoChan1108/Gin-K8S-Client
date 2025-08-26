package service

import (
	"Gin-K8S-Client/internal/client"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

type NamespaceSvc struct {
}

var (
	namespaceSvcInstance *NamespaceSvc
	//once                 sync.Once
	//ctx                  = context.Background()
)

// NewNamespaceSvc 单例加载
func NewNamespaceSvc() *NamespaceSvc {
	once.Do(func() {
		namespaceSvcInstance = &NamespaceSvc{}
	})
	return namespaceSvcInstance
}

func (ns *NamespaceSvc) GetNamespace() ([]v1.Namespace, error) {
	clientSet, err := client.GetK8sClientSet()
	if err != nil {
		klog.Error("service.GetNamespace.GetK8sClientSet err: ", err.Error())
		return []v1.Namespace{}, err
	}
	namespaceList, err := clientSet.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		klog.Error("service.GetNamespace.Namespaces().List err: ", err.Error())
		return []v1.Namespace{}, err
	}
	return namespaceList.Items, nil
}
