package websocket

import (
	"errors"
	"github.com/gorilla/websocket"
	"k8s.io/klog/v2"
	"sync"
)

/**
TODO
1. 封装websocket结构 ✅
2. 封装websocket方法(读 写) ✅
3. 封装websocket读循环、写循环 ✅
4. 封装web终端(WebSSH)
5. 封装终端读、写
*/

// WsMessage websocket消息结构体
type WsMessage struct {
	// 消息类型
	MsgType int
	// 消息内容
	Data []byte
}

// WsConnection websocket连接结构体
type WsConnection struct {
	// websocket连接
	wsConn *websocket.Conn
	//读消息管道
	inChan chan *WsMessage
	//写消息管道
	outChan chan *WsMessage
	//关闭连接的管道
	closeChan chan byte
	//标记连接是否关闭
	isClosed bool
	//连接上锁
	lock sync.Mutex
}

// ReadLoop 读循环
func (ws *WsConnection) ReadLoop() {
	var (
		msgType int
		data    []byte
		msg     *WsMessage
		err     error
	)
	for {
		//将从连接读取的消息写入inChan给ReadMsg()调用
		if msgType, data, err = ws.wsConn.ReadMessage(); err != nil {
			klog.Fatal("websocket.ReadLoop err: ", err.Error())
			//报错就关闭
			if err = ws.Close(); err != nil {
				klog.Fatal("websocket.ReadLoop.Close err: ", err.Error())
				return
			}
		}
		msg = &WsMessage{
			MsgType: msgType,
			Data:    data,
		}
		select {
		case ws.inChan <- msg:
		case <-ws.inChan:
			if ws.isClosed {
				klog.Fatal("websocket.ReadLoop err: websocket has been closed")
				return
			}
			klog.Fatal("websocket.ReadLoop err: websocket closed")
			//关闭函数
			goto ENDLOOP
		}
	}
ENDLOOP:
}

// WriteLoop 写循环
func (ws *WsConnection) WriteLoop() {
	var (
		msg *WsMessage
		err error
	)
	for {
		//从outChan中取出消息发送给连接
		select {
		case msg = <-ws.outChan:
			if err = ws.wsConn.WriteMessage(msg.MsgType, msg.Data); err != nil {
				klog.Fatal("websocket.WriteLoop err: ", err.Error())
				if err = ws.Close(); err != nil {
					klog.Fatal("websocket.WriteLoop.Close err: ", err.Error())
					return
				}
			}
		case <-ws.closeChan:
			if ws.isClosed {
				klog.Fatal("websocket.WriteLoop err: websocket has been closed")
				return
			}
			klog.Fatal("websocket.WriteLoop err: websocket closed")
			//关闭函数
			goto ENDLOOP
		}
	}
ENDLOOP:
}

// ReadMsg 读取消息
func (ws *WsConnection) ReadMsg() (*WsMessage, error) {
	select {
	case msg := <-ws.inChan:
		return msg, nil
	case <-ws.closeChan:
		//先判断是否已经关闭
		if ws.isClosed {
			klog.Fatal("websocket.ReadMsg err: websocket has been closed")
			return nil, errors.New("websocket has been closed")
		}
		return nil, errors.New("websocket.ReadMsg err: websocket closed")
	}
}

// WriteMsg 写入消息
func (ws *WsConnection) WriteMsg(msgType int, data []byte) error {
	msg := &WsMessage{
		MsgType: msgType,
		Data:    data,
	}
	select {
	case ws.outChan <- msg:
		return nil
	case <-ws.closeChan:
		//先判断是否已经关闭
		if ws.isClosed {
			klog.Fatal("websocket.WriteMsg err: websocket has been closed")
			return errors.New("websocket has been closed")
		}
		return errors.New("websocket.WriteMsg err: websocket closed")
	}
}

// Close 关闭连接
func (ws *WsConnection) Close() error {
	//上锁保证安全不被其他协程关闭
	ws.lock.Lock()
	defer ws.lock.Unlock()
	if !ws.isClosed {
		ws.isClosed = true
		close(ws.closeChan)
		if err := ws.wsConn.Close(); err != nil {
			klog.Fatal("websocket.Close err: ", err.Error())
			return errors.New("websocket close error")
		}
	}
	return nil
}
