package web

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
)

//go:embed templates/*.html
var templatesFS embed.FS

// TemplateData 模板数据
type TemplateData struct {
	Title string
}

// RenderHome 渲染首页
func RenderHome() ([]byte, error) {
	return renderPage("home", &TemplateData{Title: "首页"})
}

// renderPage 渲染页面
func renderPage(name string, data *TemplateData) ([]byte, error) {
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