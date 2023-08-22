package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/haobinfei/ginner/config"
	"github.com/haobinfei/ginner/public/common"
	"github.com/haobinfei/ginner/routes"
)

func main() {

	// 加载配置文件
	config.InitConfig()

	// 初始化日志组件
	common.InitLogger()

	// 初始化http服务
	r := routes.Init()
	srv := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf("%s:%d", config.Conf.App.Host, config.Conf.App.Port),
	}

	// 后台启动http服务
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			common.Log.Fatalf("Listen: %s\n", err)
		}
	}()

	common.Log.Infof("Server is running at %s:%d", config.Conf.App.Host, config.Conf.App.Port)

	// 优雅退出程序
	c := make(chan os.Signal)
	signal.Notify(c)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		common.Log.Error("Server forced to shutdown:", err)
	}

	fmt.Println("Server is shutdown!")

}
