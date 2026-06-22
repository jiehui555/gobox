# Gobox 开发指南

## 项目概述

基于 Huma v2 框架的 Go HTTP API 服务，附带网页后台。

## 常用命令

```bash
# 开发（热重载）
air

# 构建
go build -o server .

# 格式化代码
go fmt ./...

# 验证编译
go build ./...

# Docker 构建
docker build -t gobox .
```

## 项目结构

```
main.go              # 入口，路由注册
handler/             # API 路由处理器
  greeting.go        # GET /greeting/{name}
  review.go          # POST /reviews
  topfeel.go         # POST /topfeel/sign-in, /topfeel/reply
  web.go             # GET / (网页首页)
pkg/
  topfeel/client.go  # Topfeel API 客户端封装
  web/templates.go   # 模板渲染引擎
  web/templates/     # HTML 模板文件 (embed)
```

## 技术栈

- Go 1.25.9 + Huma v2 (OpenAPI 3.1)
- 前端：Pico CSS（通过 CDN 引入）
- 模板：Go html/template + embed

## 版本号

版本号在 `main.go:22` 的 `huma.DefaultConfig` 第二个参数中维护。

## 重要提醒

- **不要自动执行 `git commit` 和 `git tag`**，这些操作应由用户手动执行
- 代码修改完成后，等待用户确认并手动提交

## 提交规范

- 提交信息使用中文，格式：`<type>: <description>`
- 类型：feat / fix / refactor / docs
- 提交前执行 `go build ./...` 和 `go fmt ./...`

## 发版流程

1. 更新 `main.go` 中的版本号
2. 提交代码（用户手动）
3. 打 tag：`git tag v1.x.x`（用户手动）
4. 推送：`git push && git push --tags`（用户手动）
5. CI 自动构建 Docker 镜像并推送到 GHCR

## CI/CD

- 推送到 main 分支触发 Docker 构建
- 推送 `v*` tag 触发版本化构建
- 镜像推送到 `ghcr.io/jiehui555/gobox`
