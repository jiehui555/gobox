package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
	"github.com/danielgtaylor/huma/v2/humacli"
)

// CLI 的配置选项
type Options struct {
	Port int `help:"监听的端口号" short:"p" default:"8060"`
}

// GreetingOutput 表示问候操作的响应结构
type GreetingOutput struct {
	Body struct {
		Message string `json:"message" example:"Hello, world!" doc:"问候消息"`
	}
}

// ReviewInput 表示提交评价操作的请求结构
type ReviewInput struct {
	Body struct {
		Author  string `json:"author" maxLength:"10" doc:"评价作者"`
		Rating  int    `json:"rating" minimum:"1" maximum:"5" doc:"1到5的评分"`
		Message string `json:"message,omitempty" maxLength:"100" doc:"评价内容"`
	}
}

func main() {
	// 创建一个支持端口配置的 CLI 应用
	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {
		// 创建原生 HTTP ServeMux 和 API
		router := http.NewServeMux()
		api := humago.New(router, huma.DefaultConfig("我的 API", "1.0.0"))

		// 注册 GET /greeting/{name} 接口
		huma.Register(api, huma.Operation{
			OperationID: "get-greeting",
			Method:      http.MethodGet,
			Path:        "/greeting/{name}",
			Summary:     "获取问候语",
			Description: "根据姓名获取对应的问候语",
			Tags:        []string{"Greetings"},
		}, func(ctx context.Context, input *struct {
			Name string `path:"name" maxLength:"30" example:"world" doc:"要问候的姓名"`
		}) (*GreetingOutput, error) {
			resp := &GreetingOutput{}
			resp.Body.Message = fmt.Sprintf("你好，%s！", input.Name)
			return resp, nil
		})

		// 注册 POST /reviews 接口
		huma.Register(api, huma.Operation{
			OperationID:   "post-review",
			Method:        http.MethodPost,
			Path:          "/reviews",
			Summary:       "提交评价",
			Tags:          []string{"Reviews"},
			DefaultStatus: http.StatusCreated,
		}, func(ctx context.Context, i *ReviewInput) (*struct{}, error) {
			// TODO: 将评价保存到数据存储中
			return nil, nil
		})

		// 配置 CLI 如何启动服务器
		hooks.OnStart(func() {
			fmt.Printf("正在端口 %d 上启动服务器...\n", options.Port)
			http.ListenAndServe(fmt.Sprintf(":%d", options.Port), router)
		})
	})

	// 运行 CLI。当未传入命令时，默认启动服务器
	cli.Run()
}
