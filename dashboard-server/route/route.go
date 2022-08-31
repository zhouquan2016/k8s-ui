package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRoutes(ctx *gin.Engine) {
	RegisterIndexRoute(ctx)
	RegisterPodRoutes(ctx)
	RegisterConfigMapRoutes(ctx)
}

func RegisterIndexRoute(engine *gin.Engine) {
	engine.GET("/", func(ctx *gin.Context) {
		sendJson(ctx, "ok")
	})
}
func sendJson(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, data)
}
