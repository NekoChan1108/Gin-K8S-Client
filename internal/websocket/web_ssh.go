package websocket

import (
	"Gin-K8S-Client/internal/client"
	"context"
	"github.com/gorilla/websocket"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/klog/v2"
	"net/http"
)

/** 前端xterm的数据类型
{
  type: 'input' | 'resize',
  data?: string,
  cols?: number,
  rows?: number
}
*/

// xtermMessage xterm数据结构
type xtermMessage struct {
	// type: input | resize
	Type string `json:"type"`
	Data string `json:"data"`
	//行列用于调整大小
	Cols uint16 `json:"cols"`
	Rows uint16 `json:"rows"`
}

// streamHandler 流式处理websocket数据
type streamHandler struct {
	wsConn *WsConnection
	// resizeEvent 终端resize事件用于调整终端大小
	resizeEvent chan remotecommand.TerminalSize
}

// Read 流式读取websocket数据
func (handler *streamHandler) Read(data []byte) (int, error) {
	var (
		xtermMsg xtermMessage
		msg      *WsMessage
		err      error
		size     int
	)
	//读取websocket数据
	if msg, err = handler.wsConn.ReadMsg(); err != nil {
		klog.Error("websocket.Read err: ", err.Error())
		return -1, err
	}
	//反序列化前端xterm传来的消息数据到xtermMsg
	if err = json.Unmarshal(msg.Data, &xtermMsg); err != nil {
		klog.Error("websocket.streamHandler.Read.Unmarshal err: ", err.Error())
		return -1, err
	}
	//判断消息类型
	switch xtermMsg.Type {
	case "input":
		size = len(xtermMsg.Data)
		copy(data, xtermMsg.Data)
		//如果是调整大小就向管道写入对应的终端大小
	case "resize":
		handler.resizeEvent <- remotecommand.TerminalSize{
			Width:  xtermMsg.Cols,
			Height: xtermMsg.Rows,
		}
	}
	return size, nil
}

// Write 流式写入websocket数据
func (handler *streamHandler) Write(data []byte) (int, error) {
	//拷贝一份
	copyData := make([]byte, len(data))
	copy(copyData, data)
	size := len(copyData)
	//发送文本数据
	if err := handler.wsConn.WriteMsg(websocket.TextMessage, copyData); err != nil {
		klog.Error("websocket.WriteMsg err: ", err.Error())
		return -1, err
	}
	return size, nil
}

// Next 获取终端大小 实现TerminalSize的接口方法
func (handler *streamHandler) Next() *remotecommand.TerminalSize {
	terminalSize := <-handler.resizeEvent
	return &terminalSize
}

// WebSSH 开启Web终端
func WebSSH(namespace, podName, containerName, method string, resp http.ResponseWriter, req *http.Request) error {
	var (
		wsConnection *WsConnection
		err          error
		k8sConfig    *rest.Config
		k8sClientSet *kubernetes.Clientset
		executor     remotecommand.Executor
		handler      *streamHandler
		ctx          = context.Background()
	)

	if k8sConfig, err = client.GetK8sConfig(); err != nil {
		klog.Error("websocket.WebSSH.GetK8sConfig err: ", err.Error())
		return err
	}
	if k8sClientSet, err = client.GetK8sClientSet(); err != nil {
		klog.Error("websocket.WebSSH.GetK8sClientSet err: ", err.Error())
		return err
	}
	// 创建RESTClient
	requestSSH := k8sClientSet.CoreV1().RESTClient().Post().Resource("pods").Name(podName).
		Namespace(namespace).SubResource("exec").VersionedParams(&v1.PodExecOptions{
		Container: containerName,
		Command:   []string{method},
		Stdin:     true,
		Stdout:    true,
		Stderr:    true,
		TTY:       true,
	}, scheme.ParameterCodec)
	//执行命令
	if executor, err = remotecommand.NewSPDYExecutor(k8sConfig, "POST", requestSSH.URL()); err != nil {
		klog.Error("websocket.WebSSH.NewSPDYExecutor err: ", err.Error())
		return err
	}
	if wsConnection, err = NewWsConnection(resp, req); err != nil {
		klog.Error("websocket.WebSSH.NewWsConnection err: ", err.Error())
		return err
	}
	//构建流式处理websocket数据
	handler = &streamHandler{
		wsConn:      wsConnection,
		resizeEvent: make(chan remotecommand.TerminalSize),
	}
	//流式执行命令
	if err = executor.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:             handler,
		Stdout:            handler,
		Stderr:            handler,
		Tty:               true,
		TerminalSizeQueue: handler,
	}); err != nil {
		klog.Error("websocket.WebSSH.Stream err: ", err.Error())
		return wsConnection.Close()
	}
	return err
}
