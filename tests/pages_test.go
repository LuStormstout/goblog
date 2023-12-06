package tests

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test ./tests -v -count=1
// TestAllPages 测试所有页面
func TestAllPages(t *testing.T) {

	baseURL := "http://localhost:3000"

	// 声明加初始化测试数据
	var tests = []struct {
		method   string
		url      string
		expected int
	}{
		{"GET", "/", 200},
		{"GET", "/about", 200},
		{"GET", "/notfound", 404},
		{"GET", "/articles", 200},
		{"GET", "/articles/create", 200},
		{"GET", "/articles/1", 200},
		{"GET", "/articles/1/edit", 200},
		{"POST", "/articles/1", 200},
		{"POST", "/articles", 200},
		{"POST", "/articles/2/delete", 404},
	}

	// 遍历所有测试
	for _, test := range tests {
		t.Logf("当前请求 URL: %v \n", test.url)
		var (
			resp *http.Response
			err  error
		)
		// 请求以获取响应
		switch {
		case test.method == "POST":
			data := make(map[string][]string)
			resp, err = http.PostForm(baseURL+test.url, data)
		default:
			resp, err = http.Get(baseURL + test.url)
		}
		// 断言
		assert.NoError(t, err, "请求 "+test.url+" 时报错")
		assert.Equal(t, test.expected, resp.StatusCode, test.url+" 应返回状态码 "+strconv.Itoa(test.expected))
	}
}
