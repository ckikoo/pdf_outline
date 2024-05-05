package main

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

func main() {
	f := excelize.NewFile()

	// 创建一个新的工作表
	sheetName := "pdf 标签 总计"
	f.NewSheet(sheetName)

	// 向工作表中写入数据
	data := [][]interface{}{{"Name", "Age"}, {"John", 30}, {"Jane", 25}}
	for r, row := range data {
		for c, value := range row {
			columnName, err := excelize.ColumnNumberToName(c + 1)
			if err != nil {
				fmt.Println(err)
				return
			}
			cell := columnName + fmt.Sprintf("%d", r+1)
			f.SetCellValue(sheetName, cell, value)
		}
	}

	f.SetColWidth(sheetName, "A", "A", 40)
	f.SetColWidth(sheetName, "B", "B", 40)
	f.SetColWidth(sheetName, "C", "C", 20)

	// 保存Excel文件
	err := f.SaveAs("example.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
}
