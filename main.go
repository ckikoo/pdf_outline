package main

import (
	"bytes"
	fileUtil "file/util/file"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/unidoc/unipdf/v3/model"
)

// pdfPrintLine 递归打印PDF标签
func pdfPrintLine(item *model.OutlineItem, buf *bytes.Buffer) {
	if item.Dest.PageObj == nil {
		return
	}
	buf.WriteString(fmt.Sprintf("\t%v %v\n", item.Title, item.Dest.Page+1))
	for _, childItem := range item.Items() {
		pdfPrintLine(childItem, buf)
	}
}

// guoLvPDF 读取PDF文件，提取标签和页数
func guoLvPDF(root string, content chan string, wg *sync.WaitGroup, sem chan struct{}) {
	defer wg.Done()

	// 控制并发数量
	sem <- struct{}{}
	defer func() { <-sem }()

	var buf bytes.Buffer
	defer func() {
		content <- buf.String()
	}()

	f, err := os.Open(root)
	if err != nil {
		fmt.Printf("Error opening file %s: %s\n", root, err)
		return
	}
	defer f.Close()

	pdfReader, err := model.NewPdfReader(f)
	if err != nil {
		fmt.Printf("Error creating PDF reader for file %s: %s\n", root, err)
		return
	}

	// 获取PDF页数
	pageSize, err := pdfReader.GetNumPages()
	if err != nil {
		fmt.Printf("Error getting page size for file %s: %s\n", root, err)
		return
	}

	// 写入文件名和页数
	fmt.Fprintf(&buf, "%v %v\n", filepath.Base(root), pageSize)

	// 判断是否有无标签
	if rootNode := pdfReader.GetOutlineTree(); rootNode == nil {
		return
	}

	// 获取PDF标签
	lines, err := pdfReader.GetOutlines()
	if err != nil {
		fmt.Printf(" 无法获取pdf文件标签 file %s: %s\n", root, err)
		return
	}

	// 递归打印标签
	for _, item := range lines.Items() {
		pdfPrintLine(item, &buf)
	}

}

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
			go guoLvPDF(absPath, content, wg, sem) // 处理PDF文件
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
