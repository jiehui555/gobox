package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

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

// TopfeelSignInInput 表示极夜社区签到操作的请求结构
type TopfeelSignInInput struct {
	Body struct {
		TopfeelToken string `json:"topfeel_token" example:"your_topfeel_token" doc:"极夜社区 Token" required:"true"`
	}
}

// TopfeelSignInOutput 表示极夜社区签到操作的响应结构
type TopfeelSignInOutput struct {
	Body struct {
		Message string `json:"message" example:"签到成功" doc:"签到结果"`
	}
}

func main() {
	// 创建一个支持端口配置的 CLI 应用
	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {
		// 创建原生 HTTP ServeMux 和 API
		router := http.NewServeMux()
		api := humago.New(router, huma.DefaultConfig("我的 API", "1.0.0"))

		// 注册 POST /topfeel/sign-in 接口
		huma.Register(api, huma.Operation{
			OperationID: "topfeel-sign-in",
			Method:      http.MethodPost,
			Path:        "/topfeel/sign-in",
			Summary:     "极夜社区签到",
			Tags:        []string{"Topfeel"},
		}, func(ctx context.Context, i *TopfeelSignInInput) (*TopfeelSignInOutput, error) {
			// 生成时间戳（毫秒）
			now := time.Now().UnixMilli()
			newTime := now + int64(rand.Intn(4)+3)*1000 // 随机 3~6 秒

			// 构建请求体
			body := map[string]interface{}{
				"oldtime": now,
				"newtime": newTime,
			}
			bodyBytes, err := json.Marshal(body)
			if err != nil {
				return nil, huma.Error500InternalServerError("JSON 序列化失败", err)
			}

			// 创建请求
			topfeelSignInURL := "https://bbs.topfeel.com/api/gift/day_sign"
			req, err := http.NewRequest("POST", topfeelSignInURL, bytes.NewReader(bodyBytes))
			if err != nil {
				return nil, huma.Error500InternalServerError("创建请求失败", err)
			}

			// 设置请求头
			req.Header.Set("Referer", "https://bbs.topfeel.com/h5/")
			req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36")
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("sec-ch-ua-platform", `"Windows"`)
			req.Header.Set("sec-ch-ua-mobile", "?0")
			req.Header.Set("sec-ch-ua", `"Google Chrome";v="143", "Chromium";v="143", "Not A(Brand";v="24"`)
			req.Header.Set("token", i.Body.TopfeelToken)

			// 发送请求
			client := &http.Client{Timeout: 15 * time.Second}
			resp, err := client.Do(req)
			if err != nil {
				return nil, huma.Error500InternalServerError("网络请求失败", err)
			}
			defer resp.Body.Close()

			// 读取响应
			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, huma.Error500InternalServerError("读取响应失败", err)
			}

			// 非 200 状态码处理
			if resp.StatusCode != http.StatusOK {
				return nil, huma.Error500InternalServerError(fmt.Sprintf("HTTP 状态码异常 %d: %s", resp.StatusCode, string(respBody)), nil)
			}

			// 解析响应
			var result map[string]interface{}
			if err := json.Unmarshal(respBody, &result); err != nil {
				return nil, huma.Error500InternalServerError("JSON 解析失败", err)
			}

			msg := ""
			if m, ok := result["msg"].(string); ok {
				msg = m
			}

			output := &TopfeelSignInOutput{}
			output.Body.Message = msg

			return output, nil
		})

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
