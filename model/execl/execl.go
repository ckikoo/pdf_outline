package execl

import (
	"file/model/pdf"
	"fmt"

	"github.com/xuri/excelize/v2"
)

type pair struct {
	first  string // 左上角
	second string //右下角
}

type Execl struct {
	fileName  string
	file      *excelize.File
	sheetName string
}

func NewExecl(path string, sheetName string) *Execl {
	f := excelize.NewFile()
	f.SetSheetName("Sheet1", sheetName)
	f.SetColWidth(sheetName, "A", "D", 60)
	return &Execl{
		fileName:  path,
		file:      f,
		sheetName: sheetName,
	}
}

// 多次调用这个接口
func (e *Execl) WriteData(r *pdf.Record) {
	rows, err := e.file.GetRows(e.sheetName)
	if err != nil {
		fmt.Println(err)
		return
	}
	rowCount := len(rows) + 1 // 加1是因为索引从1开始

	// 写入文件名和页数
	e.file.SetCellValue(e.sheetName, "A"+fmt.Sprintf("%d", rowCount), r.FileName)
	e.file.SetCellValue(e.sheetName, "B"+fmt.Sprintf("%d", rowCount), r.PageCounts)

	e.fillColor(
		pair{first: "A" + fmt.Sprintf("%d", rowCount),
			second: "B" + fmt.Sprintf("%d", rowCount)}, "FFFF00") // 填充颜色

	rowCount++
	// 写入标签数据
	for _, label := range r.Labels {
		e.file.SetCellValue(e.sheetName, "B"+fmt.Sprintf("%d", rowCount), label.LabelName)
		e.file.SetCellValue(e.sheetName, "C"+fmt.Sprintf("%d", rowCount), label.LabelPosPage)
		rowCount++
	}
	debug := false
	if debug {
		rows, err := e.file.GetRows(e.sheetName)
		if err != nil {
			panic(err)
		}

		fmt.Printf("rows: %v\n", rows)
	}

}

func (e *Execl) fillColor(pos pair, color string) error {
	style, _ := e.file.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{color},
			Pattern: 1,
		},
	})

	return e.file.SetCellStyle(e.sheetName, pos.first, pos.second, style)
}

func (e *Execl) Close() {
	err := e.file.SaveAs(e.fileName)
	if err != nil {
		panic(err)
	}
	err = e.file.Close()
	if err != nil {
		panic(err)
	}
}
