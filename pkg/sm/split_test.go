package sm

import (
	"testing"
	"fmt"
)

func TestSplit(t *testing.T){
	files := []string{
		"D:/Zoe/Projects/GO/src/github.com/jiusanzhou/pdf2html/cmd/pdftool/test_data/13北京科锐.PDF",
	}
	outputDir := "D:/Zoe/Projects/GO/src/github.com/jiusanzhou/pdf2html/pkg/sm/test_data"

	for _, f := range files {
		fs, n ,err :=Split(f, outputDir, 10)
		if err!=nil{
			fmt.Println(err.Error())
			continue
		}
		fmt.Println("一共分了:", n, "份!")
		fmt.Println(fs)
	}
}