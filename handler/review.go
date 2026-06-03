package handler

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// ReviewInput 表示提交评价操作的请求结构
type ReviewInput struct {
	Body struct {
		Author  string `json:"author" maxLength:"10" doc:"评价作者"`
		Rating  int    `json:"rating" minimum:"1" maximum:"5" doc:"1到5的评分"`
		Message string `json:"message,omitempty" maxLength:"100" doc:"评价内容"`
	}
}

func RegisterReview(api huma.API) {
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
}
