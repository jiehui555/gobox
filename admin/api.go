package admin

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jiehui555/gobox/pkg/auth"
	"github.com/jiehui555/gobox/pkg/database"
	"golang.org/x/crypto/bcrypt"
)

// handleLogin 用户登录
func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "方法不允许")
		return
	}

	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "请求格式错误")
		return
	}

	user, err := database.GetUserByEmail(input.Email)
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"success": false,
			"message": "邮箱或密码错误",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"success": false,
			"message": "邮箱或密码错误",
		})
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Email, user.Username, user.Role)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "生成token失败")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "登录成功",
		"token":   token,
	})
}

// handleLogout 退出登录
func handleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	writeSuccess(w, "退出成功")
}

// handleUsersAPI 用户管理API
func handleUsersAPI(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetUsers(w, r)
	case http.MethodPost:
		handleCreateUser(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "方法不允许")
	}
}

// handleUserAPI 单个用户API
func handleUserAPI(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/admin/api/users/"):]
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		writeError(w, http.StatusBadRequest, "无效的用户ID")
		return
	}

	switch r.Method {
	case http.MethodPut:
		handleUpdateUser(w, r, uint(id))
	case http.MethodDelete:
		handleDeleteUser(w, r, uint(id))
	default:
		writeError(w, http.StatusMethodNotAllowed, "方法不允许")
	}
}

// handleGetUsers 获取用户列表
func handleGetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := database.GetUsers()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "获取用户列表失败")
		return
	}
	writeJSON(w, http.StatusOK, users)
}

// handleCreateUser 创建用户
func handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Username string `json:"username"`
		Role     string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "请求格式错误")
		return
	}

	_, err := database.GetUserByEmail(input.Email)
	if err == nil {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"success": false,
			"message": "邮箱已被使用",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "密码加密失败")
		return
	}

	role := input.Role
	if role == "" {
		role = "user"
	}

	user := &database.User{
		Email:    input.Email,
		Password: string(hashedPassword),
		Username: input.Username,
		Role:     role,
	}
	if err := database.CreateUser(user); err != nil {
		writeError(w, http.StatusInternalServerError, "创建用户失败")
		return
	}

	writeSuccess(w, "创建成功")
}

// handleUpdateUser 更新用户
func handleUpdateUser(w http.ResponseWriter, r *http.Request, id uint) {
	var input struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Role     string `json:"role"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "请求格式错误")
		return
	}

	existingUser, err := database.GetUserByEmail(input.Email)
	if err == nil && existingUser.ID != id {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"success": false,
			"message": "邮箱已被其他用户使用",
		})
		return
	}

	updates := map[string]interface{}{
		"email":    input.Email,
		"username": input.Username,
		"role":     input.Role,
	}

	if input.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "密码加密失败")
			return
		}
		updates["password"] = string(hashedPassword)
	}

	if err := database.UpdateUser(id, updates); err != nil {
		writeError(w, http.StatusInternalServerError, "更新用户失败")
		return
	}

	writeSuccess(w, "更新成功")
}

// handleDeleteUser 删除用户
func handleDeleteUser(w http.ResponseWriter, r *http.Request, id uint) {
	currentUserID, _ := r.Context().Value(auth.UserIDKey).(uint)
	if currentUserID == id {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"success": false,
			"message": "不能删除自己",
		})
		return
	}

	if err := database.DeleteUser(id); err != nil {
		writeError(w, http.StatusInternalServerError, "删除用户失败")
		return
	}

	writeSuccess(w, "删除成功")
}
