package main

import (
	"bufio"
	"minik8s/pkg/config"
	"os"
	"strings"
)

func main() {
	configPath := config.LocalConfigPath
	configPath = strings.Replace(configPath, ":namespace", "default", -1)
	configPath = strings.Replace(configPath, ":name", "my-dns", -1)

	// 如果不存在该文件，则创建
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		file, err := os.Create(configPath)
		if err != nil {
			panic(err)
		}
		defer file.Close()
	}

	// 读取文件
	file, err := os.Open(configPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 一行一行读取
	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// 直接在文件末尾追加

	lines = append(lines, "server {")
	lines = append(lines, "    listen 80;")
	lines = append(lines, "    server_name test.com ;")

	// 子路径
	lines = append(lines, "    location /service/example {")
	lines = append(lines, "        proxy_pass http://86.140.7.16:7080/;")
	lines = append(lines, "    }")

	// 子路径
	lines = append(lines, "    location /service/example2 {")
	lines = append(lines, "        proxy_pass http://109.51.139.197:9876/;")
	lines = append(lines, "    }")

	lines = append(lines, "}")

	// 写入文件
	err = os.WriteFile(configPath, []byte(strings.Join(lines, "\n")), 0644)
	if err != nil {
		panic(err)
	}

}
