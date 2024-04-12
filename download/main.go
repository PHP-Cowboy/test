package main

import (
	"io"
	"net/http"
	"os"
)

func main() {
	// 设置路由处理程序
	http.HandleFunc("/", downloadHandler)

	// 启动服务器并监听端口
	http.ListenAndServe(":8090", nil)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	// 打开要下载的文件
	file, err := os.Open("./base.apk")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// 设置响应头，指定文件名和内容类型
	w.Header().Set("Content-Disposition", "attachment; filename=base.apk")
	w.Header().Set("Content-Type", "application/octet-stream")

	// 将文件内容写入响应主体中
	io.Copy(w, file)
}
