package handler

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jiehui555/gobox/pkg/web"
)

// WebOutput 表示网页响应的结构
type WebOutput struct {
	ContentType string `header:"Content-Type"`
	Body        []byte
}

// RegisterWeb 注册网页路由
func RegisterWeb(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "web-home",
		Method:      http.MethodGet,
		Path:        "/",
		Summary:     "主页",
		Tags:        []string{"Web"},
	}, func(ctx context.Context, input *struct{}) (*WebOutput, error) {
		body, err := web.RenderHome()
		if err != nil {
			return nil, huma.Error500InternalServerError("渲染页面失败", err)
		}

		return &WebOutput{
			ContentType: "text/html; charset=utf-8",
			Body:        body,
		}, nil
	})
}