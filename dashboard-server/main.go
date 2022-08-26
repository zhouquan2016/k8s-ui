package main

import (
	"dashboard-server/config"
	"dashboard-server/middleware"
	"dashboard-server/route"
	"dashboard-server/ws"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	engine.Use(middleware.Cors())
	engine.Use(gin.Recovery())
	route.InitRoutes(engine)
	ws.InitWebSocket(engine);
	engine.Run(fmt.Sprintf(":%d", config.GetConfig().Server.Port))
}
