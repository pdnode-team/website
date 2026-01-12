package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

func generateRandomKey() []byte {
	bytes := make([]byte, 32) // 256位随机性
	if _, err := rand.Read(bytes); err != nil {
		panic(err) // 这种级别出错通常是操作系统出问题了
	}

	return bytes
}
func SetUpSuperuser() {
	if _, err := os.Stat(".superusertoken"); os.IsNotExist(err) {
		// --- 这里保持你原有的生成逻辑 ---
		rawToken := generateRandomKey()
		tokenString := hex.EncodeToString(rawToken)
		err := os.WriteFile(".superusertoken", []byte(tokenString), 0400)

		if os.Getenv("ENV") == "production" {
			if err != nil {
				slog.Error("Create Superuser token failed", "err", err.Error())
				panic("Create Superuser token failed")
			}
			slog.Info("Superuser token created successfully")
			slog.Warn("SECURITY ALERT", "msg", "Keep .superusertoken file safe. Leaking or deleting it will impact all logins.")
		} else {
			if err != nil {
				panic(err)
			}
			fmt.Printf("\n\033[32m[INIT]\033[0m Superuser token created: \033[33m%s\033[0m\n", tokenString)
			fmt.Printf("\033[31m%s\033[0m\n\n", "WARNING: DO NOT DELETE OR DISCLOSE .SUPERUSERTOKEN")
		}

	} else {
		// --- 优化后的 Else 分支 ---
		if os.Getenv("ENV") == "production" {
			// 生产环境：保持专业且简洁的结构化输出
			slog.Info("Superuser token check", "status", "exists", "action", "skip_generation")
			// 只有在生产环境，这种重要警告才用 Warn 等级，方便监控系统抓取
			slog.Warn("Critical file protection",
				"file", ".superusertoken",
				"warning", "Deleting this file will break user authentication")
		} else {
			// 开发环境：保持显眼，提醒开发者
			fmt.Println("\n------------------------------------------------------------------")
			fmt.Println("[\033[34mINFO\033[0m] Superuser token is already initialized. (Skip)")
			fmt.Println("[\033[33mTIP\033[0m] If you lost it, delete '.superusertoken' and restart.")

			warning := "** DANGER: Deleting this file will prevent all users from logging in **"
			fmt.Printf("\033[31m%s\033[0m\n", strings.ToUpper(warning))
			fmt.Println("------------------------------------------------------------------")
		}
	}
}
