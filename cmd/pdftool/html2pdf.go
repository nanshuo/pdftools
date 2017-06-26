package main

import (
	"fmt"
	"github.com/jiusanzhou/pdf2html/pkg/html2pdf"
	"github.com/jiusanzhou/pdf2html/pkg/wkhtml2pdf"
	"log"
)

var html2pdfFactory *html2pdf.Factory
var wkhtml2pdfFactory *wkhtml2pdf.Factory

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
