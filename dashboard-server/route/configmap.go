package route

import (
	"context"
	"dashboard-server/client"
	"dashboard-server/config"
	"sort"

	"github.com/gin-gonic/gin"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type KeyData struct {
	Key  string `json:"key"`
	Data string `json:"data"`
}


func sendMap2array(ctx *gin.Context, m map[string]string) {
	var list = make([]KeyData, 0, 10)
	for k, v := range m {
		list = append(list, KeyData{Key: k, Data: v})
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Key < list[j].Key
	})
	sendJson(ctx, config.MakeSucessResult(list))
}
func RegisterConfigMapRoutes(engine *gin.Engine) {
	engine.GET("/configmap/:namespace/:name", func(ctx *gin.Context) {
		namespace := ctx.Param("namespace")
		name := ctx.Param("name")
		configmap, err := client.DefaultClient.ConfigMaps(namespace).Get(context.TODO(), name, v1.GetOptions{})
		if err != nil {
			sendJson(ctx, config.MakeErrorResult(err.Error()))
			return
		}
		sendMap2array(ctx, configmap.Data)
	})
	engine.PUT("/configmap/:namespace/:name", func(ctx *gin.Context) {

		var keyMap map[string]string
		err := ctx.BindJSON(&keyMap)
		if err != nil {
			sendJson(ctx, config.MakeErrorResult(err.Error()))
			return
		}
		namespace := ctx.Param("namespace")
		name := ctx.Param("name")
		configmap, err := client.DefaultClient.ConfigMaps(namespace).Get(context.TODO(), name, v1.GetOptions{})
		if err != nil {
			sendJson(ctx, config.MakeErrorResult(err.Error()))
			return
		}
		if configmap.Data == nil {
			configmap.Data = keyMap
		} else {
			for k, v := range keyMap {
				configmap.Data[k] = v
			}
		}

		configmap, err = client.DefaultClient.ConfigMaps(namespace).Update(context.TODO(), configmap, v1.UpdateOptions{})
		if err != nil {
			sendJson(ctx, config.MakeErrorResult(err.Error()))
		} else {
			sendMap2array(ctx, configmap.Data)
		}
	})
	engine.DELETE("/configmap/:namespace/:name/:key", func(ctx *gin.Context) {
		namespace := ctx.Param("namespace")
		name := ctx.Param("name")
		key := ctx.Param("key")
		configmap, err := client.DefaultClient.ConfigMaps(namespace).Get(context.TODO(), name, v1.GetOptions{})
		if err != nil {
			sendJson(ctx, config.MakeErrorResult(err.Error()))
			return
		}
		delete(configmap.Data, key)
		configmap, err = client.DefaultClient.ConfigMaps(namespace).Update(context.TODO(), configmap, v1.UpdateOptions{})
		if err != nil {
			sendJson(ctx, config.MakeErrorResult(err.Error()))
		} else {
			sendMap2array(ctx, configmap.Data)
		}
	})
}
