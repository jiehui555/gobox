package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jiehui555/gobox/pkg/web"
)

// contextKey 类型用于 context 的 key
type contextKey string

const (
	UserIDKey    contextKey = "user_id"
	UserEmailKey contextKey = "user_email"
	UsernameKey  contextKey = "user_username"
	UserRoleKey  contextKey = "user_role"
)

// ExtractTokenFromCookie 从cookie中提取token
func ExtractTokenFromCookie(r *http.Request) string {
	cookie, err := r.Cookie("token")
	if err == nil && cookie != nil {
		return cookie.Value
	}
	return ""
}

// ExtractTokenFromHeader 从Authorization头中提取token
func ExtractTokenFromHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1]
		}
	}
	return ""
}

// AdminAuthMiddleware 管理后台认证中间件，从cookie获取token
func AdminAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := ExtractTokenFromCookie(r)

		if tokenString == "" {
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}

		claims, err := ParseToken(tokenString)
		if err != nil {
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UserEmailKey, claims.Email)
		ctx = context.WithValue(ctx, UsernameKey, claims.Username)
		ctx = context.WithValue(ctx, UserRoleKey, claims.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// extractTokenFromCookieHuma 从cookie中提取token (huma版本)
func extractTokenFromCookieHuma(ctx huma.Context) string {
	cookie, err := huma.ReadCookie(ctx, "token")
	if err == nil && cookie != nil {
		return cookie.Value
	}
	return ""
}

// extractTokenFromHeaderHuma 从Authorization头中提取token (huma版本)
func extractTokenFromHeaderHuma(ctx huma.Context) string {
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
		tokenString := extractTokenFromCookieHuma(ctx)

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
		tokenString := extractTokenFromHeaderHuma(ctx)

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
