package config

import "goblog/pkg/config"

func init() {
	config.Add("session", config.StrMap{

		// session 驱动，可选: cookie、file、memory、redis、mysql、postgresql、sqlite3、mssql，目前只支持 cookie
		"default": config.Env("SESSION_DRIVER", "cookie"),

		// 会话的 Cookie 名称
		"session_name": config.Env("SESSION_NAME", "goblog-session"),
	})
}
