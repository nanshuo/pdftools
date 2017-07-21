package sm

import (
	"fmt"
	"github.com/jiusanzhou/pdf2html/pkg/util"
	unipdf "github.com/jiusanzhou/unidoc/pdf"
	"os"
	"path/filepath"
	"math"
)

func Split(f string, outputDir string, maxPage int) (fs []string, n int, err error) {
	// pureName.index_Ux_Uy_rotate.pdf

	fmt.Println("Max page:", maxPage)

	// make sure f/outputDir is ths abs path
	if !filepath.IsAbs(f) {
		f, _ = filepath.Abs(f)
	}

	if !util.Exists(outputDir) {
		os.MkdirAll(outputDir, 0666)
	}

	if !filepath.IsAbs(outputDir) {
		outputDir, _ = filepath.Abs(outputDir)
	}

	_, name := filepath.Split(f)
	ext := filepath.Ext(name)
	pureName := name[:len(name)-len(ext)]

	pages, err := readPdf(f)
	if err != nil {
		return
	}

	var x, y float64
	var r int64
	var pageContainer []*unipdf.PdfPage = []*unipdf.PdfPage{}
	for _, p := range pages {

		// check page content


		Px, Py, Pr, err := getSize(p)
		if err != nil {
			fmt.Println(err.Error())
			pageContainer = append(pageContainer, p)
			continue
		}

		// if pages' length greater than max page
		// write them all
		l := len(pageContainer)
		if l == 0 {
		} else if l >= maxPage {
			// over fill the container
			// flush them all
			flush(&pageContainer, &fs, &n, outputDir, pureName, ext, x, y, r)
		} else {
			// check x, y, r
			if Px != x || Py != y || Pr != r {
				// different with last one
				// flush to file
				flush(&pageContainer, &fs, &n, outputDir, pureName, ext, x, y, r)
			}
		}
		// add page to container
		pageContainer = append(pageContainer, p)
		//update x, y, r with the new value
		x = Px
		y = Py
		r = Pr

	}
	if len(pageContainer) != 0 {
		flush(&pageContainer, &fs, &n, outputDir, pureName, ext, x, y, r)
	}
	return
}

func flush(pageContainer *[]*unipdf.PdfPage, fileList *[]string, counter *int, outputDir, pureName, ext string, x, y float64, rotate int64) {

	iX := int(math.Floor(x))
	iY := int(math.Floor(y))
	of := fmt.Sprintf("%s.%d_%d_%d_%d%s", filepath.Join(outputDir, pureName), *counter, iX, iY, rotate, ext)
	if err := writePdf(of, *pageContainer); err == nil {
		*pageContainer = []*unipdf.PdfPage{}
		*fileList = append(*fileList, of)
		*counter++
	} else {
		fmt.Println(err.Error())
	}
}
