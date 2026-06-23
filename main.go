package main

import (
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/danielgtaylor/huma/v2/humacli"

	"github.com/jiehui555/gobox/admin"
	"github.com/jiehui555/gobox/handler"
	"github.com/jiehui555/gobox/pkg/database"
)

// Options CLI 的配置选项
type Options struct {
	Port   int    `help:"监听的端口号" short:"p" default:"8060"`
	DBPath string `help:"数据库文件路径" short:"d" default:"gobox.db"`
}

func main() {
	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {
		// 初始化数据库
		database.Init(options.DBPath)
		router := http.NewServeMux()
		api := humago.New(router, huma.DefaultConfig("我的 API", "1.6.2"))

		// 注册API路由（使用huma）
		handler.RegisterTopfeelSignIn(api)
		handler.RegisterTopfeelReply(api)

		// 注册管理后台路由（原生HTTP，不走huma）
		admin.Register(router)

		hooks.OnStart(func() {
			fmt.Printf("正在端口 %d 上启动服务器...\n", options.Port)
			http.ListenAndServe(fmt.Sprintf(":%d", options.Port), router)
		})
	})

	cli.Run()
}
