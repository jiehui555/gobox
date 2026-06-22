package handler

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jiehui555/gobox/pkg/auth"
	"github.com/jiehui555/gobox/pkg/database"
	"github.com/jiehui555/gobox/pkg/web"
	"golang.org/x/crypto/bcrypt"
)

// WebOutput 表示网页响应的结构
type WebOutput struct {
	ContentType string `header:"Content-Type"`
	Body        []byte
}

// LoginInput 登录请求输入
type LoginInput struct {
	Body struct {
		Email    string `json:"email" example:"admin@gobox.com" doc:"邮箱" required:"true"`
		Password string `json:"password" example:"admin123" doc:"密码" required:"true"`
	}
}

// LoginOutput 登录响应输出
type LoginOutput struct {
	Body struct {
		Success bool   `json:"success" example:"true" doc:"是否成功"`
		Message string `json:"message" example:"登录成功" doc:"消息"`
		Token   string `json:"token" example:"eyJhbGciOiJIUzI1NiIs..." doc:"JWT token"`
	}
}

// CheckAuthInput 检查登录状态输入
type CheckAuthInput struct {
	Token string `header:"Authorization" doc:"Bearer token"`
}

// CheckAuthOutput 检查登录状态输出
type CheckAuthOutput struct {
	Body struct {
		Authenticated bool   `json:"authenticated" example:"true" doc:"是否已登录"`
		User          *User  `json:"user,omitempty" doc:"用户信息"`
	}
}

// User 用户信息
type User struct {
	ID       uint   `json:"id" example:"1" doc:"用户ID"`
	Email    string `json:"email" example:"admin@gobox.com" doc:"邮箱"`
	Username string `json:"username" example:"管理员" doc:"用户名"`
	Role     string `json:"role" example:"admin" doc:"角色"`
}

// RegisterWeb 注册网页路由
func RegisterWeb(api huma.API) {
	// 创建Web认证中间件
	webAuthMiddleware := auth.WebAuthMiddleware(api)

	huma.Register(api, huma.Operation{
		OperationID: "web-home",
		Method:      http.MethodGet,
		Path:        "/",
		Summary:     "主页",
		Tags:        []string{"Web"},
		Middlewares: huma.Middlewares{webAuthMiddleware},
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

	huma.Register(api, huma.Operation{
		OperationID: "web-login-page",
		Method:      http.MethodGet,
		Path:        "/login",
		Summary:     "登录页面",
		Tags:        []string{"Web"},
	}, func(ctx context.Context, input *struct{}) (*WebOutput, error) {
		body, err := web.RenderLogin()
		if err != nil {
			return nil, huma.Error500InternalServerError("渲染页面失败", err)
		}

		return &WebOutput{
			ContentType: "text/html; charset=utf-8",
			Body:        body,
		}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "web-login",
		Method:      http.MethodPost,
		Path:        "/api/login",
		Summary:     "用户登录",
		Tags:        []string{"Auth"},
	}, func(ctx context.Context, input *LoginInput) (*LoginOutput, error) {
		// 根据邮箱查找用户
		user, err := database.GetUserByEmail(input.Body.Email)
		if err != nil {
			return &LoginOutput{
				Body: struct {
					Success bool   `json:"success" example:"true" doc:"是否成功"`
					Message string `json:"message" example:"登录成功" doc:"消息"`
					Token   string `json:"token" example:"eyJhbGciOiJIUzI1NiIs..." doc:"JWT token"`
				}{
					Success: false,
					Message: "邮箱或密码错误",
				},
			}, nil
		}

		// 验证密码
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Body.Password))
		if err != nil {
			return &LoginOutput{
				Body: struct {
					Success bool   `json:"success" example:"true" doc:"是否成功"`
					Message string `json:"message" example:"登录成功" doc:"消息"`
					Token   string `json:"token" example:"eyJhbGciOiJIUzI1NiIs..." doc:"JWT token"`
				}{
					Success: false,
					Message: "邮箱或密码错误",
				},
			}, nil
		}

		// 生成JWT token
		token, err := auth.GenerateToken(user.ID, user.Email, user.Username, user.Role)
		if err != nil {
			return nil, huma.Error500InternalServerError("生成token失败", err)
		}

		return &LoginOutput{
			Body: struct {
				Success bool   `json:"success" example:"true" doc:"是否成功"`
				Message string `json:"message" example:"登录成功" doc:"消息"`
				Token   string `json:"token" example:"eyJhbGciOiJIUzI1NiIs..." doc:"JWT token"`
			}{
				Success: true,
				Message: "登录成功",
				Token:   token,
			},
		}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "web-check-auth",
		Method:      http.MethodGet,
		Path:        "/api/auth/check",
		Summary:     "检查登录状态",
		Tags:        []string{"Auth"},
	}, func(ctx context.Context, input *CheckAuthInput) (*CheckAuthOutput, error) {
		// 从header中获取token
		tokenString := input.Token
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		// 解析token
		claims, err := auth.ParseToken(tokenString)
		if err != nil {
			return &CheckAuthOutput{
				Body: struct {
					Authenticated bool   `json:"authenticated" example:"true" doc:"是否已登录"`
					User          *User  `json:"user,omitempty" doc:"用户信息"`
				}{
					Authenticated: false,
				},
			}, nil
		}

		return &CheckAuthOutput{
			Body: struct {
				Authenticated bool   `json:"authenticated" example:"true" doc:"是否已登录"`
				User          *User  `json:"user,omitempty" doc:"用户信息"`
			}{
				Authenticated: true,
				User: &User{
					ID:       claims.UserID,
					Email:    claims.Email,
					Username: claims.Username,
					Role:     claims.Role,
				},
			},
		}, nil
	})
}