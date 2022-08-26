package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRoutes(ctx *gin.Engine) {
	
	RegisterPodRoutes(ctx)
	RegisterConfigMapRoutes(ctx)
}

func sendJson(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, data)
}