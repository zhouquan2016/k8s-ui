package ws

import (
	"dashboard-server/client"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
)

type PtyHandler interface {
	io.Reader
	io.Writer
	remotecommand.TerminalSizeQueue
}

// TerminalSession
type TerminalSession struct {
	wsConn   *websocket.Conn
	sizeChan chan remotecommand.TerminalSize
}

type termMessage struct {
	Type    string `JSON:"type,omitempty"`
	Cols    uint16 `JSON:"cols,omitempty"`
	Rows    uint16 `JSON:"rows,omitempty"`
	Message string `JSON:"message,omitempty"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func InitWebSocket(engine *gin.Engine) {
	engine.GET("/ws", handlerWebSocket)
}

func handlerWebSocket(ctx *gin.Context) {
	namespace := ctx.Query("namespace")
	if namespace == "" {
		namespace = "default"
	}
	pod := ctx.Query("pod")
	if pod == "" {
		ctx.JSON(http.StatusOK, "pod不能为空")
		return
	}
	wsConn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.String(http.StatusOK, "升级websocket协议失败:" + err.Error())
		return
	}
	ptyHandler := &TerminalSession{wsConn: wsConn, sizeChan: make(chan remotecommand.TerminalSize, 1)}
	
	req := client.DefaultClient.RESTClient().Post().Resource("pods").Namespace(namespace).Name(pod).SubResource("exec")
	if req == nil {
		ptyHandler.WriteText("get Resource Request fail!")
		return
	}
	req.VersionedParams(&v1.PodExecOptions{
		Command: []string{"/bin/bash"},
		Stdin:   true,
		Stdout:  true,
		Stderr:  true,
		TTY:     true,
	}, scheme.ParameterCodec)

	executor, err := remotecommand.NewSPDYExecutor(&client.DefaultConfig, "POST", req.URL())
	if err != nil {
		ptyHandler.WriteText(err.Error())
		return
	}
	err = executor.Stream(remotecommand.StreamOptions{
		Stdin:             ptyHandler,
		Stdout:            ptyHandler,
		Stderr:            ptyHandler,
		TerminalSizeQueue: ptyHandler,
		Tty:               false,
	})
	if err != nil {
		ptyHandler.WriteText(err.Error())
		return
	}
}

// Next called in a loop from remotecommand as long as the process is running
func (t *TerminalSession) Next() *remotecommand.TerminalSize {
	size := <-t.sizeChan
	return &size
}

// Read called in a loop from remotecommand as long as the process is running
func (t *TerminalSession) Read(p []byte) (n int, err error) {
	defer func() {
		recover_error := recover()
		if recover_error != nil {
			n = 0
			err = fmt.Errorf("%v", recover_error)
			fmt.Printf("%v", recover_error)
			t.Close()
		}
	}()
	_, message, err := t.wsConn.ReadMessage()
	if err != nil {
		return 0, err
	}
	m := termMessage{}
	err = json.Unmarshal(message, &m)
	if err != nil {
		return 0, err
	}
	var text string
	switch m.Type {
	case "input":
		text = m.Message
	case "resize":
		t.sizeChan <- remotecommand.TerminalSize{Width: m.Rows, Height: m.Cols}
		text = ""
	default:
		text = "unsupported message type"
	}
	copy(p, []byte(text))
	return len(text), nil
}

// Write called from remotecommand whenever there is any output
func (t *TerminalSession) Write(p []byte) (n int, err error) {
	err = t.WriteText(string(p))
	if err == nil {
		n = len(p)
	}else {
		n = 0
	}
	return n, err
}

func (t *TerminalSession) WriteText(text string) (err error) {
	return t.wsConn.WriteMessage(websocket.TextMessage, []byte(text))
}

// Close close session
func (t *TerminalSession) Close() error {
	return t.wsConn.Close()
}
