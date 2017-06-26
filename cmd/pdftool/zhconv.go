package main

import (
	"errors"
	"fmt"
	"github.com/jiusanzhou/pdf2html/pkg/zhconv"
	"os"
	"path/filepath"
	"strings"
)

var zhconvFactory *zhconv.Factory

func zhconvInit() {
	zhconvFactory = zhconv.NewFactory()
}

func zhconvClose() {

}

type T struct {
	i, m, fileType string
}

func prepare(filename string) (*T, error) {
	ext := strings.ToUpper(filepath.Ext(filename))

	var t *T
	switch ext {
	case ".PDF":
		// convert to html first
		fmt.Println("先将PDF转换为HTML,这会花一些时间，请耐心等候")

		if pdf2htmlFactory == nil {
			pdf2htmlInit()
		}

		m, err := pdf2htmlFactory.NewMaterial(filename, opts.tmpDir, "")
		if err != nil {
			return nil, err
		}

		// convert html
		p, err := pdf2htmlFactory.Convert(m)
		if err != nil {
			return nil, err
		}

		fmt.Println("PDF->HTML 耗时:", p.Coast.String())

		t = &T{filename, p.FilePath, "PDF"}
	case "", ".HTML", ".TXT", ".JSON", ".XML", ".HTM":
		// convert directly
		t = &T{filename, filename, "TEXT"}
	default:
		return nil, errors.New("不支持的文件类型")
	}

	return t, nil
}

func toSimple() {

	var err error
	var t *T
	var outputFile string
	var simpleFile string

	if zhconvFactory == nil {
		zhconvInit()
	}

	for _, i := range opts.inputs {

		t, err = prepare(i)

		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		ext := filepath.Ext(t.m)
		base, nameWithExt := filepath.Split(t.m)

		// name := strings.Replace(nameWithExt, ext, "", -1)
		name := base + nameWithExt[:len(nameWithExt)-len(ext)]

		if opts.suffix != "" {
			simpleFile = name + opts.suffix + ext
		} else {
			simpleFile = name + "_simple" + ext
		}

		err = zhconvFactory.FileToSimple(t.m, simpleFile)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		if len(opts.inputs) == 1 && len(opts.outputs) == 1 {
			outputFile = opts.outputs[0]
		} else {
			ext := filepath.Ext(t.i)
			base, nameWithExt := filepath.Split(t.i)
			// name := strings.Replace(nameWithExt, ext, "", -1)
			name = base + nameWithExt[:len(nameWithExt)-len(ext)]

			if opts.suffix != "" {
				outputFile = name + opts.suffix + ext
			} else {
				outputFile = name + "_simple" + ext
			}
		}

		switch t.fileType {
		case "PDF":
			fmt.Println("再将HTML转换为PDF,这会花一些时间,请耐心等候")

			if opts.wkhtml2pdf != "" {

				if wkhtml2pdfFactory == nil {
					wkhtml2pdfInit()
				}

				m, err := wkhtml2pdfFactory.NewMaterial(simpleFile, opts.outputDir, outputFile)
				if err != nil {
					fmt.Println(err.Error())
					continue
				}

				p, err := wkhtml2pdfFactory.Convert(m)
				if err != nil {
					fmt.Println(err.Error())
					continue
				}

				fmt.Println("[wkhtmltopdf]HTML->PDF 耗时:", p.Coast.String())

				// remove old file
				os.Remove(t.m)
				os.Remove(simpleFile)

				fmt.Println("转换完成:", p.FilePath)

			} else {

				// convert back to pdf
				if html2pdfFactory == nil {
					html2pdfInit()
				}

				m, err := html2pdfFactory.NewMaterial(simpleFile, opts.outputDir, outputFile, opts.scale)
				if err != nil {
					fmt.Println(err.Error())
					continue
				}

				p, err := html2pdfFactory.Convert(m)
				if err != nil {
					fmt.Println(err.Error())
					continue
				}

				fmt.Println("[chrome]HTML->PDF 耗时:", p.Coast.String())

				// remove old file
				os.Remove(t.m)
				os.Remove(simpleFile)

				fmt.Println("转换完成:", p.FilePath)
			}
		}
	}
}

func toTradition() {

	var err error
	var t *T
	var outputFile string

	if zhconvFactory == nil {
		zhconvInit()
	}

	for _, i := range opts.inputs {
		t, err = prepare(i)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		err = zhconvFactory.FileToTraditional(t.i, "")
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		if len(opts.inputs) == 1 && len(opts.outputs) == 1 {
			outputFile = opts.outputs[0]
		}

		switch t.fileType {
		case "PDF":
			// convert back to pdf
			m, err := html2pdfFactory.NewMaterial(t.m, opts.outputDir, outputFile, opts.scale)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			p, err := html2pdfFactory.Convert(m)
			if err != nil {
				fmt.Println(err.Error())
			}

			// remove old file
			os.Remove(m.FilePath)

			fmt.Println("转换完成：", p.FilePath)
		}
	}
}
