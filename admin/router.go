package admin

import (
	"net/http"

	"github.com/jiehui555/gobox/pkg/auth"
)

// Register 注册管理后台路由
func Register(mux *http.ServeMux) {
	// 公开页面
	mux.HandleFunc("/login", handleLoginPage)
	mux.HandleFunc("/login/submit", handleLogin)
	mux.HandleFunc("/logout/submit", handleLogout)

	// 需要认证的页面
	mux.Handle("/", auth.AdminAuthMiddleware(http.HandlerFunc(handleHome)))
	mux.Handle("/users", auth.AdminAuthMiddleware(http.HandlerFunc(handleUsersPage)))

	// 需要认证的API
	mux.Handle("/admin/auth/check", auth.AdminAuthMiddleware(http.HandlerFunc(handleCheckAuth)))
	mux.Handle("/admin/users", auth.AdminAuthMiddleware(http.HandlerFunc(handleUsersAPI)))
	mux.Handle("/admin/users/", auth.AdminAuthMiddleware(http.HandlerFunc(handleUserAPI)))
}
