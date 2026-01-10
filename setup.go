package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
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
		// 1. 生成原始字节
		rawToken := generateRandomKey()

		// 2. 转换为人类可读的十六进制字符串 (不再是乱码)
		tokenString := hex.EncodeToString(rawToken)

		// 3. 写入文件
		err := os.WriteFile(".superusertoken", []byte(tokenString), 0400)
		if err != nil {
			panic("Create Superuser token failed")
		}

		// 4. 打印明文 (使用 \033)
		fmt.Println("Created superuser token: " + tokenString)

		// 修正颜色代码 \033[31m
		warning := "Please keep your Super User Token file safe and do not disclose or delete it."
		fmt.Printf("\033[31m%s\033[0m\n", strings.ToUpper(warning))
	} else {
		// 进入常规逻辑
		fmt.Println("\n\n\n\n\n\n\nYour superuser key has been generated: Skip.\n\nIf you haven't it, please delete the .superusertoken file, and restart.\n\nDo not delete the .superusertoken file. If you confirm that this token has been leaked, please delete it immediately.")
		fmt.Printf("\033[31m")
		fmt.Println(strings.ToUpper("\n **Please note that deleting the .superusertoken file means all users will be unable to log in.** \n\n"))
		fmt.Printf("\033[0m")

	}
}
