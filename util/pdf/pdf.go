package pdfUtil

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
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
func GuoLvPDF(root string, content chan string, wg *sync.WaitGroup, sem chan struct{}) {
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
