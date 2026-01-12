package config

import (
	"log/slog"
	"os"
)

func InitLogger() {
	var handler slog.Handler

	// 设置配置项
	opts := &slog.HandlerOptions{
		// AddSource: true 可以在日志中看到具体是哪一个文件、哪一行代码打的日志
		AddSource: true,
		// 默认级别是 Info，你可以根据环境变量动态调整
		Level: slog.LevelDebug,
	}

	if os.Getenv("ENV") == "production" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	logger := slog.New(handler)

	// 设置为全局默认，这样你在其他 package 直接调用 slog.Info 即可
	slog.SetDefault(logger)
}
