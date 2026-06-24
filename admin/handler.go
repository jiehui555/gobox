package admin

import (
	"net/http"

	"github.com/jiehui555/gobox/pkg/auth"
	"github.com/jiehui555/gobox/pkg/database"
	"github.com/jiehui555/gobox/pkg/web"
)

// handleLoginRouter 登录路由（GET 返回页面，POST 处理登录）
func handleLoginRouter(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleLoginPage(w, r)
	case http.MethodPost:
		handleLogin(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "方法不允许")
	}
}

// handleLoginPage 登录页面
func handleLoginPage(w http.ResponseWriter, r *http.Request) {
	body, err := web.RenderLogin()
	if err != nil {
		http.Error(w, "渲染页面失败", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(body)
}

// handleHome 主页
func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/admin" {
		http.NotFound(w, r)
		return
	}
	body, err := web.RenderHome()
	if err != nil {
		http.Error(w, "渲染页面失败", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(body)
}

// handleUsersPage 用户管理页面
func handleUsersPage(w http.ResponseWriter, r *http.Request) {
	users, err := database.GetUsers()
	if err != nil {
		http.Error(w, "获取用户列表失败", http.StatusInternalServerError)
		return
	}
	body, err := web.RenderUsers(users)
	if err != nil {
		http.Error(w, "渲染页面失败", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(body)
}

// handleCheckAuth 检查登录状态
func handleCheckAuth(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(auth.UserIDKey).(uint)
	email, _ := r.Context().Value(auth.UserEmailKey).(string)
	username, _ := r.Context().Value(auth.UsernameKey).(string)
	role, _ := r.Context().Value(auth.UserRoleKey).(string)

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"authenticated": true,
		"user": map[string]interface{}{
			"id":       userID,
			"email":    email,
			"username": username,
			"role":     role,
		},
	})
}
