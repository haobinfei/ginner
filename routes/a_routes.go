package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/haobinfei/ginner/config"
)

func Init() *gin.Engine {
	// 设置模式
	gin.SetMode(config.Conf.App.Mode)

	// 创建gin路由, 增加恢复中间件
	r := gin.New()
	r.Use(gin.Recovery())

	return r
}
