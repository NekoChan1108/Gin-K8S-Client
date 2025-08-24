package router

import (
	"Gin-K8S-Client/api"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	mainGroup := r.Group("/api")
	{
		mainGroup.GET("/", api.Hello)
		k8sGroup := mainGroup.Group("/kubernetes")
		{
			k8sGroup.GET("/namespace", api.GetNamespace)
			k8sGroup.GET("/pod", api.GetPod)
		}
	}
}
