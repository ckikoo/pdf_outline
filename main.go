package main

import (
	fileUtil "file/util/file"
	pdfUtil "file/util/pdf"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// outputFile 输出文件内容
func outputFile(outFilePath string, content chan string, done chan struct{}) {

	outfile, err := os.OpenFile(outFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer outfile.Close() // 确保输出文件在函数结束时关闭

	// 从通道读取内容并写入文件
	for c := range content {
		outfile.WriteString(c)
	}

	close(done) // 关闭 done 通道，通知主函数内容已经写入完毕
}

func main() {
	rootDir := "."
	suffix := ".pdf"
	outFilePath := "find_out.txt"
	ch := make(chan string)       // 存储文件路径的缓冲通道
	content := make(chan string)  // 存储文件内容的缓冲通道
	done := make(chan struct{})   // 用于通知主函数内容已经写入完毕的通道
	sem := make(chan struct{}, 4) // 控制并发数量的信号量，限制最多同时有n个协程
	wg := &sync.WaitGroup{}       // 用于等待所有协程完成的 WaitGroup

	go fileUtil.TraverseDirectories(rootDir, ch) // 并发遍历目录

	go outputFile(outFilePath, content, done) // 启动输出文件内容的协程

	// 启动多个协程处理文件
	for filePath := range ch {
		absPath := filePath
		fileName := filepath.Base(filePath)      // 获取文件名
		if strings.HasSuffix(fileName, suffix) { // 判断是否为PDF文件
			wg.Add(1)
			go pdfUtil.GuoLvPDF(absPath, content, wg, sem) // 处理PDF文件
		}
	}

	// 等待所有协程完成
	go func() {
		wg.Wait()
		close(content) // 关闭 content 通道，通知输出文件协程没有更多内容了
	}()

	// 等待主函数写入完毕
	<-done
}
