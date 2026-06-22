package auth

import (
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jiehui555/gobox/pkg/web"
)

// extractToken 从请求中提取token
func extractToken(ctx huma.Context) string {
	// 首先尝试从cookie中获取token
	cookie, err := huma.ReadCookie(ctx, "token")
	if err == nil && cookie != nil {
		return cookie.Value
	}

	// 如果cookie中没有token，尝试从Authorization头获取
	authHeader := ctx.Header("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1]
		}
	}

	return ""
}

// WebAuthMiddleware Web路由认证中间件，未认证时返回登录页面
func WebAuthMiddleware(api huma.API) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		tokenString := extractToken(ctx)

		if tokenString == "" {
			body, err := web.RenderLogin()
			if err != nil {
				huma.WriteErr(api, ctx, http.StatusInternalServerError, "渲染页面失败", err)
				return
			}
			ctx.SetHeader("Content-Type", "text/html; charset=utf-8")
			ctx.SetStatus(http.StatusUnauthorized)
			ctx.BodyWriter().Write(body)
			return
		}

		claims, err := ParseToken(tokenString)
		if err != nil {
			body, err := web.RenderLogin()
			if err != nil {
				huma.WriteErr(api, ctx, http.StatusInternalServerError, "渲染页面失败", err)
				return
			}
			ctx.SetHeader("Content-Type", "text/html; charset=utf-8")
			ctx.SetStatus(http.StatusUnauthorized)
			ctx.BodyWriter().Write(body)
			return
		}

		ctx = huma.WithValue(ctx, "user_id", claims.UserID)
		ctx = huma.WithValue(ctx, "user_email", claims.Email)
		ctx = huma.WithValue(ctx, "user_username", claims.Username)
		ctx = huma.WithValue(ctx, "user_role", claims.Role)

		next(ctx)
	}
}

// APIAuthMiddleware API路由认证中间件，未认证时返回JSON错误
func APIAuthMiddleware(api huma.API) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		tokenString := extractToken(ctx)

		if tokenString == "" {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "未登录", nil)
			return
		}

		claims, err := ParseToken(tokenString)
		if err != nil {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "token无效或已过期", nil)
			return
		}

		ctx = huma.WithValue(ctx, "user_id", claims.UserID)
		ctx = huma.WithValue(ctx, "user_email", claims.Email)
		ctx = huma.WithValue(ctx, "user_username", claims.Username)
		ctx = huma.WithValue(ctx, "user_role", claims.Role)

		next(ctx)
	}
}
