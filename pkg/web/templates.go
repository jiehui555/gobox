package web

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"

	"github.com/jiehui555/gobox/pkg/database"
)

//go:embed templates/*.html
var templatesFS embed.FS

// TemplateData 模板数据
type TemplateData struct {
	Title string
}

// UsersPageData 用户列表页面数据
type UsersPageData struct {
	TemplateData
	Users []database.User
}

// RenderHome 渲染首页
func RenderHome() ([]byte, error) {
	return renderPage("home", &TemplateData{Title: "首页"})
}

// RenderUsers 渲染用户列表页面
func RenderUsers(users []database.User) ([]byte, error) {
	return renderPage("users", &UsersPageData{
		TemplateData: TemplateData{Title: "用户管理"},
		Users:        users,
	})
}

// RenderLogin 渲染登录页面
func RenderLogin() ([]byte, error) {
	return renderPage("login", &TemplateData{Title: "登录"})
}

// renderPage 渲染页面
func renderPage(name string, data any) ([]byte, error) {
	// 读取页面模板（先解析，定义"content"模板）
	pageBytes, err := templatesFS.ReadFile("templates/" + name + ".html")
	if err != nil {
		return nil, fmt.Errorf("读取页面模板失败: %w", err)
	}

	// 先解析页面模板，定义"content"模板
	tmpl, err := template.New(name).Parse(string(pageBytes))
	if err != nil {
		return nil, fmt.Errorf("解析页面模板失败: %w", err)
	}

	// 读取布局模板
	layoutBytes, err := templatesFS.ReadFile("templates/layout.html")
	if err != nil {
		return nil, fmt.Errorf("读取布局模板失败: %w", err)
	}

	// 解析布局模板
	tmpl, err = tmpl.Parse(string(layoutBytes))
	if err != nil {
		return nil, fmt.Errorf("解析布局模板失败: %w", err)
	}

	// 执行模板
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "layout", data); err != nil {
		return nil, fmt.Errorf("执行模板失败: %w", err)
	}

	return buf.Bytes(), nil
}
