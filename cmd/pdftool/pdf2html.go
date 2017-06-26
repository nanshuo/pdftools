package main

import (
	"fmt"
	"github.com/jiusanzhou/pdf2html/pkg/pdf2html"
	"log"
)

var pdf2htmlFactory *pdf2html.Factory

func pdf2htmlInit() {

	var err error
	pdf2htmlConfig := &pdf2html.Config{
		OutputDir: opts.outputDir,
		Exec:      opts.pdf2htmlEx,
		ExecTpl:   opts.pdf2htmlExTpl,
	}
	pdf2htmlFactory, err = pdf2html.NewFactory(pdf2htmlConfig)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func pdf2htmlClose() {
	if pdf2htmlFactory != nil {
		pdf2htmlFactory.Close()
	}
}

func toHtml() {
	fmt.Println("Convert to HTML")
}
