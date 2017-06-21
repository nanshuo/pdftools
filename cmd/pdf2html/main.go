package main

import (
	"github.com/jiusanzhou/pdf2html/pkg/pdf2html"
	"log"
)

func main() {

	opts := ParseArgs()

	c := &pdf2html.Config{
		Exec:       opts.exec,
		ExecTpl:    opts.execTpl,
		Concurrent: opts.parallel,
	}

	factory, err := pdf2html.NewFactory(c)
	if err != nil {
		log.Fatalln("Init factory error, ", err.Error())
	}

	if opts.input == "" {
		log.Fatalln("You must offer at least one file to convert.")
	}

	log.Printf("准备转换文件: %s\n", opts.input)

	m, err := factory.NewMaterial(opts.input, opts.outputDir, opts.output)
	if err != nil {
		log.Fatalln(err.Error())
	}

	p, err := factory.Convert(m)
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Printf(p.FilePath, "转换完成，耗时:", p.Coast.String(), "，文件大小:", p.Size)
}
