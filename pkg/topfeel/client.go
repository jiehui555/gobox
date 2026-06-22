package topfeel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

const (
	baseURL      = "https://bbs.topfeel.com"
	signInURL    = baseURL + "/api/gift/day_sign"
	replyURL     = baseURL + "/api/user/addGoodsComment"
	defaultUA    = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36"
	defaultRef   = "https://bbs.topfeel.com/h5/"
	defaultSecUA = `"Google Chrome";v="143", "Chromium";v="143", "Not A(Brand";v="24"`
)

// Client 封装了与极夜社区 API 的交互
type Client struct {
	token      string
	httpClient *http.Client
}

// NewClient 创建一个新的 Topfeel 客户端
func NewClient(token string) *Client {
	return &Client{
		token: token,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (c *Client) newRequest(method, url string, body interface{}) (*http.Request, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("序列化请求体失败: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Referer", defaultRef)
	req.Header.Set("User-Agent", defaultUA)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua", defaultSecUA)
	req.Header.Set("token", c.token)

	return req, nil
}

func (c *Client) doRequest(req *http.Request) (map[string]interface{}, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("网络请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP 状态码异常 %d: %s", resp.StatusCode, string(respBody))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("JSON 解析失败: %w", err)
	}

	return result, nil
}

// SignIn 执行签到打卡
func (c *Client) SignIn() (string, error) {
	now := time.Now().UnixMilli()
	newTime := now + int64(rand.Intn(4)+3)*1000

	body := map[string]interface{}{
		"oldtime": now,
		"newtime": newTime,
	}

	req, err := c.newRequest("POST", signInURL, body)
	if err != nil {
		return "", err
	}

	result, err := c.doRequest(req)
	if err != nil {
		return "", err
	}

	msg := ""
	if m, ok := result["msg"].(string); ok {
		msg = m
	}

	return msg, nil
}

// Reply 执行评论打卡
func (c *Client) Reply() (string, error) {
	now := time.Now()
	content := now.Format("打卡 01.02")

	body := map[string]interface{}{
		"images":   "",
		"goods_id": "1863",
		"vocdec":   0,
		"voc":      "",
		"content":  content,
		"pid":      45528,
		"to_name":  "辉HHH",
	}

	req, err := c.newRequest("POST", replyURL, body)
	if err != nil {
		return "", err
	}

	result, err := c.doRequest(req)
	if err != nil {
		return "", err
	}

	msg := ""
	if m, ok := result["msg"].(string); ok {
		msg = m
	}

	return msg, nil
}
