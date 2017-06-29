package sm

import (
	"fmt"
	"github.com/CodisLabs/codis/pkg/utils/errors"
	"github.com/jiusanzhou/pdf2html/pkg/util"
	unipdf "github.com/jiusanzhou/unidoc/pdf"
	"os"
	"path/filepath"
)

func MergetFromDir(dir, target string) (n int, err error) {

	if !util.Exists(dir) {
		err = errors.New("Not exits file.")
		return
	}

	// list all file from dir
	files, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		return
	}

	// check whether pdf file
	// for _, f := range files {}

	return Merge(target, files)
}

func Merge(f string, fs []string) (n int, err error) {

	if !filepath.IsAbs(f) {
		f, _ = filepath.Abs(f)
	}

	writer := unipdf.NewPdfWriter()

	for _, fi := range fs {
		if !filepath.IsAbs(fi) {
			fi, _ = filepath.Abs(fi)
		}

		pages, err := readPdf(fi)
		if err != nil {
			// add error message?
		} else {
			for _, p := range pages {
				_, _, r, _ := getSize(p)
				if r == int64(90) {
					fmt.Println(p.PieceInfo)
					fmt.Println(p.Contents)
					fmt.Println(*p.Resources)
				}

				n++
				writer.AddPage(p.GetPageAsIndirectObject())
			}
		}
	}

	fio, err := os.Create(f)
	if err != nil {
		return
	}
	writer.Write(fio)
	fio.Close()
	return
}
