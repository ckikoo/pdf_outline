package fileUtil

import (
	"log"
	"os"
	"path/filepath"
)

// 寻找文件
func TraverseDirectories(root string, ch chan<- string) {
	defer close(ch) // Close the channel when traversal is complete

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			ch <- path
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
