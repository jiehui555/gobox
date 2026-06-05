package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

// TopfeelSignInInput 表示极夜社区签到操作的请求结构
type TopfeelSignInInput struct {
	Body struct {
		TopfeelToken string `json:"topfeel_token" example:"your_topfeel_token" doc:"极夜社区 Token" required:"true"`
	}
}

// TopfeelSignInOutput 表示极夜社区签到操作的响应结构
type TopfeelSignInOutput struct {
	Body struct {
		Message string `json:"message" example:"签到成功" doc:"签到结果"`
	}
}

// TopfeelReplyInput 表示极夜社区回复操作的请求结构
type TopfeelReplyInput struct {
	Body struct {
		TopfeelToken string `json:"topfeel_token" example:"your_topfeel_token" doc:"极夜社区 Token" required:"true"`
	}
}

// TopfeelReplyOutput 表示极夜社区回复操作的响应结构
type TopfeelReplyOutput struct {
	Body struct {
		Message string `json:"message" example:"回复成功" doc:"回复结果"`
	}
}

func RegisterTopfeelSignIn(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "topfeel-sign-in",
		Method:      http.MethodPost,
		Path:        "/topfeel/sign-in",
		Summary:     "极夜社区-签到打卡",
		Tags:        []string{"Topfeel"},
	}, func(ctx context.Context, i *TopfeelSignInInput) (*TopfeelSignInOutput, error) {
		// 生成时间戳（毫秒）
		now := time.Now().UnixMilli()
		newTime := now + int64(rand.Intn(4)+3)*1000 // 随机 3~6 秒

		// 构建请求体
		body := map[string]interface{}{
			"oldtime": now,
			"newtime": newTime,
		}
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, huma.Error500InternalServerError("JSON 序列化失败", err)
		}

		// 创建请求
		topfeelSignInURL := "https://bbs.topfeel.com/api/gift/day_sign"
		req, err := http.NewRequest("POST", topfeelSignInURL, bytes.NewReader(bodyBytes))
		if err != nil {
			return nil, huma.Error500InternalServerError("创建请求失败", err)
		}

		// 设置请求头
		req.Header.Set("Referer", "https://bbs.topfeel.com/h5/")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("sec-ch-ua-platform", `"Windows"`)
		req.Header.Set("sec-ch-ua-mobile", "?0")
		req.Header.Set("sec-ch-ua", `"Google Chrome";v="143", "Chromium";v="143", "Not A(Brand";v="24"`)
		req.Header.Set("token", i.Body.TopfeelToken)

		// 发送请求
		client := &http.Client{Timeout: 15 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return nil, huma.Error500InternalServerError("网络请求失败", err)
		}
		defer resp.Body.Close()

		// 读取响应
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, huma.Error500InternalServerError("读取响应失败", err)
		}

		// 非 200 状态码处理
		if resp.StatusCode != http.StatusOK {
			return nil, huma.Error500InternalServerError(fmt.Sprintf("HTTP 状态码异常 %d: %s", resp.StatusCode, string(respBody)), nil)
		}

		// 解析响应
		var result map[string]interface{}
		if err := json.Unmarshal(respBody, &result); err != nil {
			return nil, huma.Error500InternalServerError("JSON 解析失败", err)
		}

		msg := ""
		if m, ok := result["msg"].(string); ok {
			msg = m
		}

		output := &TopfeelSignInOutput{}
		output.Body.Message = msg

		return output, nil
	})
}

func RegisterTopfeelReply(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "topfeel-reply",
		Method:      http.MethodPost,
		Path:        "/topfeel/reply",
		Summary:     "极夜社区-评论打卡",
		Tags:        []string{"Topfeel"},
	}, func(ctx context.Context, i *TopfeelReplyInput) (*TopfeelReplyOutput, error) {
		// 生成当天日期作为打卡内容
		now := time.Now()
		content := now.Format("打卡 01.02")

		// 构建请求体
		body := map[string]interface{}{
			"images":   "",
			"goods_id": "1863",
			"vocdec":   0,
			"voc":      "",
			"content":  content,
			"pid":      45528,
			"to_name":  "辉HHH",
		}
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, huma.Error500InternalServerError("JSON 序列化失败", err)
		}

		// 创建请求
		topfeelReplyURL := "https://bbs.topfeel.com/api/user/addGoodsComment"
		req, err := http.NewRequest("POST", topfeelReplyURL, bytes.NewReader(bodyBytes))
		if err != nil {
			return nil, huma.Error500InternalServerError("创建请求失败", err)
		}

		// 设置请求头
		req.Header.Set("Referer", "https://bbs.topfeel.com/h5/")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("sec-ch-ua-platform", `"Windows"`)
		req.Header.Set("sec-ch-ua-mobile", "?0")
		req.Header.Set("sec-ch-ua", `"Google Chrome";v="143", "Chromium";v="143", "Not A(Brand";v="24"`)
		req.Header.Set("token", i.Body.TopfeelToken)

		// 发送请求
		client := &http.Client{Timeout: 15 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return nil, huma.Error500InternalServerError("网络请求失败", err)
		}
		defer resp.Body.Close()

		// 读取响应
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, huma.Error500InternalServerError("读取响应失败", err)
		}

		// 非 200 状态码处理
		if resp.StatusCode != http.StatusOK {
			return nil, huma.Error500InternalServerError(fmt.Sprintf("HTTP 状态码异常 %d: %s", resp.StatusCode, string(respBody)), nil)
		}

		// 解析响应
		var result map[string]interface{}
		if err := json.Unmarshal(respBody, &result); err != nil {
			return nil, huma.Error500InternalServerError("JSON 解析失败", err)
		}

		msg := ""
		if m, ok := result["msg"].(string); ok {
			msg = m
		}

		output := &TopfeelReplyOutput{}
		output.Body.Message = msg

		return output, nil
	})
}
