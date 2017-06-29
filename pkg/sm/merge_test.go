package sm

import "testing"

func TestMerge(t *testing.T) {

	dirs := []string{
		"D:/Zoe/Projects/GO/src/github.com/jiusanzhou/pdf2html/pkg/sm/test_data",
	}

	for _, dir := range dirs {
		MergetFromDir(dir)
	}
}