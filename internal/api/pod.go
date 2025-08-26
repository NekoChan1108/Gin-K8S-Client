package api

import (
	"Gin-K8S-Client/internal/service"
	"Gin-K8S-Client/internal/websocket"
	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
	"net/http"
)

func GetPod(c *gin.Context) {
	podSvc := service.NewPodSvc()
	pods, err := podSvc.GetPod()
	if err != nil {
		klog.Error("api.GetPod.GetPod err: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"data":    nil,
			"message": "获取Pod列表失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    pods,
		"message": "ok",
	})
}

func ExecPod(c *gin.Context) {
	namespace := c.Param("namespace")
	podName := c.Param("podName")
	containerName := c.Param("containerName")
	method := c.DefaultQuery("action", "sh")
	if err := websocket.WebSSH(namespace, podName, containerName, method, c.Writer, c.Request); err != nil {
		klog.Error("api.ExecPod.WebSSH err: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"data":    nil,
			"message": "WebSSH失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    nil,
		"message": "ok",
	})
}
