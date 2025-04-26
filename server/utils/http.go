package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
)

// HttpRequest 用于发送HTTP请求
func HttpRequest(
	urlString string, //请求的URL字符串
	method string, // 请求方法（GET、POST）
	headers map[string]string, // 请求头（如Content-type）
	params map[string]string, // 查询参数（如?key=value&key2=value2）
	data any) (*http.Response, error) {

	// 创建url对象
	u, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}

	// 向URL添加查询参数
	query := u.Query()
	for k, v := range params {
		query.Set(k, v)
	}
	// 更新URL的查询部分
	u.RawQuery = query.Encode()

	// 将请求体数据编码成json格式
	buf := new(bytes.Buffer) // 创建一个缓冲区用于存储结构体
	if data != nil {
		b, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		buf = bytes.NewBuffer(b)
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	// 设置请求头
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	if data != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
