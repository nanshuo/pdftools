package main

import (
	"fmt"
	"github.com/CodisLabs/codis/pkg/utils/errors"
	"github.com/jiusanzhou/pdf2html/pkg/html2pdf"
	"github.com/jiusanzhou/pdf2html/pkg/wkhtml2pdf"
	"log"
	"path/filepath"
	"strings"
	"strconv"
)

var html2pdfFactory *html2pdf.Factory
var wkhtml2pdfFactory *wkhtml2pdf.Factory

const (
	point2inch = 72.5
)

func getHtml2PdfMetial(filePath, outputDir, outputFileName string, scale float64, landScape bool) (m *html2pdf.Material, err error) {
	if wkhtml2pdfFactory != nil {
		return wkhtml2pdfFactory.NewMaterial(filePath, outputDir, outputFileName)
	} else if html2pdfFactory != nil {
		return html2pdfFactory.NewMaterial(filePath, outputDir, outputFileName, scale, landScape)
	} else {
		err = errors.New("You should init html2pdf factory at first.")
		return
	}
}

func getFromHtml2PdfFactory() (*html2pdf.Product, error) {
	if wkhtml2pdfFactory != nil {
		return wkhtml2pdfFactory.Get()
	} else if html2pdfFactory != nil {
		return html2pdfFactory.Get()
	} else {
		return nil, errors.New("You should init html2pdf factory at first.")
	}
}

func addToHtml2PdfFactory(simpleFile, outputDir, simplePdfFile string, scale float64) error {

	if wkhtml2pdfFactory != nil {
		m, err := wkhtml2pdfFactory.NewMaterial(simpleFile, outputDir, simplePdfFile)
		if err != nil {
			return err
		}
		return wkhtml2pdfFactory.Put(m)
	} else if html2pdfFactory != nil {
		var landScape bool
		_, name := filepath.Split(simpleFile)
		pureName := name[:len(name)-len(filepath.Ext(name))]
		xx := strings.Split(pureName, "_")
		if xx[len(xx)-1] == "90" {
			landScape = true
		}

		m, err := html2pdfFactory.NewMaterial(simpleFile, outputDir, simplePdfFile, scale, landScape)

		iY := 11.6
		iX := 8.2

		if len(xx)>3{
			if aX, err := strconv.Atoi(xx[len(xx)-3]); err==nil{
				if aY, err:= strconv.Atoi(xx[len(xx)-2]); err==nil{

					// set page size
					// from page size table
					iY = float64(aY) / point2inch
					iX = float64(aX) / point2inch

					// if x and y reserve
					// set landscape to true
					if aX > aY {
						m.Params.Landscape = true
					}
				}
			}
		}
		// set width and high
		m.Params.PaperHeight = iY
		m.Params.PaperWidth = iX

		if err != nil {
			return err
		}
		return html2pdfFactory.Put(m)
	} else {
		return errors.New("You should init html2pdf factory at first.")
	}
}

func html2pdfInit() {

	var err error
	html2pdfConfig := &html2pdf.Config{
		Chrome:    opts.chrome,
		OutputDir: opts.outputDir,
	}
	html2pdfFactory, err = html2pdf.NewFactory(html2pdfConfig)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func wkhtml2pdfInit() {
	var err error
	wkhtml2pdfConfig := &wkhtml2pdf.Config{
		Exec:      opts.wkhtml2pdf,
		ExecTpl:   opts.wkhtml2pdfTpl,
		OutputDir: opts.outputDir,
	}
	wkhtml2pdfFactory, err = wkhtml2pdf.NewFactory(wkhtml2pdfConfig)
	if err != nil {
		log.Fatalln(err.Error())
	}

}

func html2pdfClose() {
	if html2pdfFactory != nil {
		html2pdfFactory.Close()
	}
}

func wkhtml2pdfClose() {
	if wkhtml2pdfFactory != nil {
		wkhtml2pdfFactory.Close()
	}
}

func toPdf() {
	fmt.Println("Convert to PDF")
}
