package config

import "goblog/pkg/config"

func init() {
	config.Add("app", config.StrMap{
		// 应用名称
		"name": config.Env("APP_NAME", "GoBlog"),

		// 环境，用以区分多套配置
		"env": config.Env("APP_ENV", "production"),

		// 是否进入调试模式
		"debug": config.Env("APP_DEBUG", false),

		// 应用服务端口
		"port": config.Env("APP_PORT", "3000"),

		// gorilla/sessions 在 Cookie 中加密数据时使用
		"key": config.Env("APP_KEY", "AHO99RL1ApfPUQXIgUo3jPWZMpfL5y4o"),

		// 用以生成链接
		"url": config.Env("APP_URL", "http://localhost:3000"),
	})
}
