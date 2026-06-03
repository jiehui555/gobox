package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// GreetingOutput 表示问候操作的响应结构
type GreetingOutput struct {
	Body struct {
		Message string `json:"message" example:"Hello, world!" doc:"问候消息"`
	}
}

func RegisterGreeting(api huma.API) {
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
}
