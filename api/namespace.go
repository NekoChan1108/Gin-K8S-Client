package api

import (
	"Gin-K8S-Client/internal/service"
	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
	"net/http"
)

func GetNamespace(c *gin.Context) {
	namespaceSvc := service.NewNamespaceSvc()
	namespace, err := namespaceSvc.GetNamespace()
	if err != nil {
		klog.Error("api.GetNamespace.GetNamespace err: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"data":    nil,
			"message": "获取命名空间列表失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    namespace,
		"message": "ok",
	})
}
