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

// WebLoginInput 登录请求输入
type WebLoginInput struct {
	Body struct {
		Email    string `json:"email" example:"admin@gobox.com" doc:"邮箱" required:"true"`
		Password string `json:"password" example:"admin123" doc:"密码" required:"true"`
	}
}

// WebLoginOutput 登录响应输出
type WebLoginOutput struct {
	Body struct {
		Success bool   `json:"success" example:"true" doc:"是否成功"`
		Message string `json:"message" example:"登录成功" doc:"消息"`
		Token   string `json:"token" example:"eyJhbGciOiJIUzI1NiIs..." doc:"JWT token"`
	}
}

// WebCheckAuthInput 检查登录状态输入
type WebCheckAuthInput struct {
	Token string `header:"Authorization" doc:"Bearer token"`
}

// WebCheckAuthOutput 检查登录状态输出
type WebCheckAuthOutput struct {
	Body struct {
		Authenticated bool     `json:"authenticated" example:"true" doc:"是否已登录"`
		User          *WebUser `json:"user,omitempty" doc:"用户信息"`
	}
}

// WebUser 用户信息
type WebUser struct {
	ID       uint   `json:"id" example:"1" doc:"用户ID"`
	Email    string `json:"email" example:"admin@gobox.com" doc:"邮箱"`
	Username string `json:"username" example:"管理员" doc:"用户名"`
	Role     string `json:"role" example:"admin" doc:"角色"`
}

// WebCreateUserInput 创建用户请求输入
type WebCreateUserInput struct {
	Body struct {
		Email    string `json:"email" example:"user@example.com" doc:"邮箱" required:"true"`
		Password string `json:"password" example:"password123" doc:"密码" required:"true"`
		Username string `json:"username" example:"张三" doc:"用户名" required:"true"`
		Role     string `json:"role" example:"user" doc:"角色"`
	}
}

// WebDeleteUserInput 删除用户请求输入
type WebDeleteUserInput struct {
	ID uint `path:"id" doc:"用户ID"`
}

// WebUpdateUserInput 更新用户请求输入
type WebUpdateUserInput struct {
	ID   uint `path:"id" doc:"用户ID"`
	Body struct {
		Email    string `json:"email" example:"user@example.com" doc:"邮箱" required:"true"`
		Username string `json:"username" example:"张三" doc:"用户名" required:"true"`
		Role     string `json:"role" example:"user" doc:"角色"`
		Password string `json:"password" example:"" doc:"密码（留空则不修改）"`
	}
}

// WebResultOutput 通用操作结果输出
type WebResultOutput struct {
	Body struct {
		Success bool   `json:"success" example:"true" doc:"是否成功"`
		Message string `json:"message" doc:"消息"`
	}
}

// RegisterWeb 注册网页路由
func RegisterWeb(api huma.API) {
	// 创建Web认证中间件
	webAuthMiddleware := auth.WebAuthMiddleware(api)
	// 创建API认证中间件
	apiAuthMiddleware := auth.APIAuthMiddleware(api)

	huma.Register(api, huma.Operation{
		OperationID: "web-home",
		Method:      http.MethodGet,
		Path:        "/",
		Summary:     "页面-主页",
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
		OperationID: "web-users",
		Method:      http.MethodGet,
		Path:        "/users",
		Summary:     "页面-用户管理",
		Tags:        []string{"Web"},
		Middlewares: huma.Middlewares{webAuthMiddleware},
	}, func(ctx context.Context, input *struct{}) (*WebOutput, error) {
		users, err := database.GetUsers()
		if err != nil {
			return nil, huma.Error500InternalServerError("获取用户列表失败", err)
		}

		body, err := web.RenderUsers(users)
		if err != nil {
			return nil, huma.Error500InternalServerError("渲染页面失败", err)
		}

		return &WebOutput{
			ContentType: "text/html; charset=utf-8",
			Body:        body,
		}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "web-api-create-user",
		Method:      http.MethodPost,
		Path:        "/web/users",
		Summary:     "用户管理-创建用户",
		Tags:        []string{"Web"},
		Middlewares: huma.Middlewares{apiAuthMiddleware},
	}, func(ctx context.Context, input *WebCreateUserInput) (*WebResultOutput, error) {
		// 检查邮箱是否已存在
		_, err := database.GetUserByEmail(input.Body.Email)
		if err == nil {
			output := &WebResultOutput{}
			output.Body.Success = false
			output.Body.Message = "邮箱已被使用"
			return output, nil
		}

		// 密码哈希
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Body.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, huma.Error500InternalServerError("密码加密失败", err)
		}

		// 设置默认角色
		role := input.Body.Role
		if role == "" {
			role = "user"
		}

		// 创建用户
		user := &database.User{
			Email:    input.Body.Email,
			Password: string(hashedPassword),
			Username: input.Body.Username,
			Role:     role,
		}
		if err := database.CreateUser(user); err != nil {
			return nil, huma.Error500InternalServerError("创建用户失败", err)
		}

		output := &WebResultOutput{}
		output.Body.Success = true
		output.Body.Message = "创建成功"
		return output, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "web-api-delete-user",
		Method:      http.MethodDelete,
		Path:        "/web/users/{id}",
		Summary:     "用户管理-删除用户",
		Tags:        []string{"Web"},
		Middlewares: huma.Middlewares{apiAuthMiddleware},
	}, func(ctx context.Context, input *WebDeleteUserInput) (*WebResultOutput, error) {
		// 获取当前用户ID，不能删除自己
		currentUserID, _ := ctx.Value("user_id").(uint)
		if currentUserID == input.ID {
			output := &WebResultOutput{}
			output.Body.Success = false
			output.Body.Message = "不能删除自己"
			return output, nil
		}

		if err := database.DeleteUser(input.ID); err != nil {
			return nil, huma.Error500InternalServerError("删除用户失败", err)
		}

		output := &WebResultOutput{}
		output.Body.Success = true
		output.Body.Message = "删除成功"
		return output, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "web-api-update-user",
		Method:      http.MethodPut,
		Path:        "/web/users/{id}",
		Summary:     "用户管理-更新用户",
		Tags:        []string{"Web"},
		Middlewares: huma.Middlewares{apiAuthMiddleware},
	}, func(ctx context.Context, input *WebUpdateUserInput) (*WebResultOutput, error) {
		// 检查邮箱是否被其他用户使用
		existingUser, err := database.GetUserByEmail(input.Body.Email)
		if err == nil && existingUser.ID != input.ID {
			output := &WebResultOutput{}
			output.Body.Success = false
			output.Body.Message = "邮箱已被其他用户使用"
			return output, nil
		}

		updates := map[string]interface{}{
			"email":    input.Body.Email,
			"username": input.Body.Username,
			"role":     input.Body.Role,
		}

		// 如果提供了新密码则更新
		if input.Body.Password != "" {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Body.Password), bcrypt.DefaultCost)
			if err != nil {
				return nil, huma.Error500InternalServerError("密码加密失败", err)
			}
			updates["password"] = string(hashedPassword)
		}

		if err := database.UpdateUser(input.ID, updates); err != nil {
			return nil, huma.Error500InternalServerError("更新用户失败", err)
		}

		output := &WebResultOutput{}
		output.Body.Success = true
		output.Body.Message = "更新成功"
		return output, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "web-login-page",
		Method:      http.MethodGet,
		Path:        "/login",
		Summary:     "页面-登录",
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
		OperationID: "web-api-login",
		Method:      http.MethodPost,
		Path:        "/web/login",
		Summary:     "认证-用户登录",
		Tags:        []string{"Web"},
	}, func(ctx context.Context, input *WebLoginInput) (*WebLoginOutput, error) {
		// 根据邮箱查找用户
		user, err := database.GetUserByEmail(input.Body.Email)
		if err != nil {
			output := &WebLoginOutput{}
			output.Body.Success = false
			output.Body.Message = "邮箱或密码错误"
			return output, nil
		}

		// 验证密码
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Body.Password))
		if err != nil {
			output := &WebLoginOutput{}
			output.Body.Success = false
			output.Body.Message = "邮箱或密码错误"
			return output, nil
		}

		// 生成JWT token
		token, err := auth.GenerateToken(user.ID, user.Email, user.Username, user.Role)
		if err != nil {
			return nil, huma.Error500InternalServerError("生成token失败", err)
		}

		output := &WebLoginOutput{}
		output.Body.Success = true
		output.Body.Message = "登录成功"
		output.Body.Token = token
		return output, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "web-api-logout",
		Method:      http.MethodPost,
		Path:        "/web/logout",
		Summary:     "认证-退出登录",
		Tags:        []string{"Web"},
	}, func(ctx context.Context, input *struct{}) (*WebResultOutput, error) {
		output := &WebResultOutput{}
		output.Body.Success = true
		output.Body.Message = "退出成功"
		return output, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "web-api-check-auth",
		Method:      http.MethodGet,
		Path:        "/web/auth/check",
		Summary:     "认证-检查登录状态",
		Tags:        []string{"Web"},
	}, func(ctx context.Context, input *WebCheckAuthInput) (*WebCheckAuthOutput, error) {
		// 从header中获取token
		tokenString := input.Token
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		// 解析token
		claims, err := auth.ParseToken(tokenString)
		if err != nil {
			output := &WebCheckAuthOutput{}
			output.Body.Authenticated = false
			return output, nil
		}

		output := &WebCheckAuthOutput{}
		output.Body.Authenticated = true
		output.Body.User = &WebUser{
			ID:       claims.UserID,
			Email:    claims.Email,
			Username: claims.Username,
			Role:     claims.Role,
		}
		return output, nil
	})
}
