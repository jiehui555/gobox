package admin

import (
	"net/http"

	"github.com/jiehui555/gobox/pkg/auth"
)

// Register 注册管理后台路由
func Register(mux *http.ServeMux) {
	// 公开操作
	mux.HandleFunc("/admin/login", handleLoginRouter)
	mux.HandleFunc("/admin/logout", handleLogout)

	// 需要认证的页面
	mux.Handle("/admin", auth.AdminAuthMiddleware(http.HandlerFunc(handleHome)))
	mux.Handle("/admin/users", auth.AdminAuthMiddleware(http.HandlerFunc(handleUsersPage)))

	// 需要认证的API
	mux.Handle("/admin/api/auth/check", auth.AdminAuthMiddleware(http.HandlerFunc(handleCheckAuth)))
	mux.Handle("/admin/api/users", auth.AdminAuthMiddleware(http.HandlerFunc(handleUsersAPI)))
	mux.Handle("/admin/api/users/", auth.AdminAuthMiddleware(http.HandlerFunc(handleUserAPI)))
}
