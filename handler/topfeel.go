package handler

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jiehui555/gobox/pkg/topfeel"
)

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

// TopfeelReplyInput 表示极夜社区回复操作的请求结构
type TopfeelReplyInput struct {
	Body struct {
		TopfeelToken string `json:"topfeel_token" example:"your_topfeel_token" doc:"极夜社区 Token" required:"true"`
	}
}

// TopfeelReplyOutput 表示极夜社区回复操作的响应结构
type TopfeelReplyOutput struct {
	Body struct {
		Message string `json:"message" example:"回复成功" doc:"回复结果"`
	}
}

func RegisterTopfeelSignIn(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "topfeel-sign-in",
		Method:      http.MethodPost,
		Path:        "/topfeel/sign-in",
		Summary:     "极夜社区-签到打卡",
		Tags:        []string{"Topfeel"},
	}, func(ctx context.Context, i *TopfeelSignInInput) (*TopfeelSignInOutput, error) {
		client := topfeel.NewClient(i.Body.TopfeelToken)
		msg, err := client.SignIn()
		if err != nil {
			return nil, huma.Error500InternalServerError("签到失败", err)
		}

		output := &TopfeelSignInOutput{}
		output.Body.Message = msg

		return output, nil
	})
}

func RegisterTopfeelReply(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "topfeel-reply",
		Method:      http.MethodPost,
		Path:        "/topfeel/reply",
		Summary:     "极夜社区-评论打卡",
		Tags:        []string{"Topfeel"},
	}, func(ctx context.Context, i *TopfeelReplyInput) (*TopfeelReplyOutput, error) {
		client := topfeel.NewClient(i.Body.TopfeelToken)
		msg, err := client.Reply()
		if err != nil {
			return nil, huma.Error500InternalServerError("回复失败", err)
		}

		output := &TopfeelReplyOutput{}
		output.Body.Message = msg

		return output, nil
	})
}
