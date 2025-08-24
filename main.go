package main

import (
	"Gin-K8S-Client/internal/config"
	"Gin-K8S-Client/router"
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
)

func main() {
	r := gin.Default()
	gin.SetMode(gin.DebugMode)
	router.InitRouter(r)
	if err := r.Run(fmt.Sprintf("%s:%d",
		config.GetString(config.ServerHost),
		config.GetInt(config.ServerPort))); err != nil {
		klog.Fatal("main.Run failed: ", err.Error())
		return
	}
}
