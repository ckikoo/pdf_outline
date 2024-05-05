package pdf

import "sort"

type Label struct {
	LabelName    string //标签名
	LabelPosPage int    // 标签位于页码
	LabelLevel   int    // 几级标签
}

type Record struct {
	FileName   string // 文件名
	PageCounts int    // 页数
	Labels     []Label
}

// sortLabelsByPage 函数按照标签位于页码从小到大对标签数组进行排序
func SortLabelsByPage(record *Record) {
	sort.Slice(record.Labels, func(i, j int) bool {
		return record.Labels[i].LabelPosPage < record.Labels[j].LabelPosPage
	})
}
