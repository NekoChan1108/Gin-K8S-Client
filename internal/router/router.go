package router

import (
	api2 "Gin-K8S-Client/internal/api"
	"Gin-K8S-Client/internal/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	middleware.InitMiddleware(r)
	mainGroup := r.Group("/api")
	{
		mainGroup.GET("/", api2.Hello)
		k8sGroup := mainGroup.Group("/kubernetes")
		{
			k8sGroup.GET("/namespace", api2.GetNamespace)
			k8sGroup.GET("/pod", api2.GetPod)
		}
	}
}
