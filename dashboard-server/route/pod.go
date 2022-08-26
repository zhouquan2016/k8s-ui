package route

import (
	"context"
	"dashboard-server/client"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	coreV1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func RegisterPodRoutes(engine *gin.Engine) {
	engine.GET("/pod/:namespace", func(ctx *gin.Context) {
		namespace := ctx.Param("namespace")
		if namespace == "" {
			namespace = "default"
		}
		podList, err := client.DefaultClient.Pods(namespace).List(context.TODO(), v1.ListOptions{})
		if err != nil {
			ctx.JSON(http.StatusOK, err.Error())
		} else {
			ctx.JSON(http.StatusOK, podList)
		}

	})

	engine.GET("/log/:namespace/:pod", func(ctx *gin.Context) {
		namespace := ctx.Param("namespace")
		if namespace == "" {
			namespace = "default"
		}
		pod := ctx.Param("pod")

		var upgrader = websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}

		wsConn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		defer wsConn.Close()
		var lines int64 = 1000
		reader, err := client.DefaultClient.Pods(namespace).GetLogs(pod, &coreV1.PodLogOptions{TailLines: &lines, Follow: true}).Stream(context.TODO())
		if err != nil {
			wsConn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			return
		}
		defer func() {
			log.Println("close reader")
			reader.Close()
		}()
		var buff [4098]byte
		for {
			size, err := reader.Read(buff[:])
			if err != nil {
				wsConn.WriteMessage(websocket.TextMessage, []byte("read log error:"+err.Error()))
				break
			}
			if size <= 0 {
				break
			}
			err = wsConn.WriteMessage(websocket.TextMessage, buff[0:size])
			if err != nil {
				log.Println("write websocket message error:", err)
				break
			}
		}
	})
}
