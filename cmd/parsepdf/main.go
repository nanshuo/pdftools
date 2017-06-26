package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	unipdf "github.com/unidoc/unidoc/pdf"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("At least one PDF file.")
		// os.Exit(1)
		args = append(args, "D:/Zoe/Projects/GO/src/github.com/jiusanzhou/pdf2html/cmd/pdf2html/test_data/01双鹭药业.pdf")
	}

	f, err := os.Open(args[0])
	if err != nil {
		log.Fatalln(err.Error())
	}

	defer f.Close()

	pdfReader, _ := unipdf.NewPdfReader(f)
	pdfWriter := unipdf.NewPdfWriter()

	for _, p := range pdfReader.PageList {

		fmt.Println(p.Contents)
		pdfWriter.AddPage(p.GetPageAsIndirectObject())
	}

	n, err := os.Create("D:/Zoe/Projects/GO/src/github.com/jiusanzhou/pdf2html/cmd/parsepdf/test.pdf")
	pdfWriter.Write(n)
}
