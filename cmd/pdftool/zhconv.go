package main

import (
	"errors"
	"fmt"
	"github.com/jiusanzhou/pdf2html/pkg/sm"
	"github.com/jiusanzhou/pdf2html/pkg/util"
	"github.com/jiusanzhou/pdf2html/pkg/zhconv"
	"os"
	"path/filepath"
	"strings"
	"sync"
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
		// split split

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

func old() {

	//var err error
	//var t *T
	//var outputFile string
	//var simpleFile string
	//t, err = prepare(i)
	//
	//if err != nil {
	//	fmt.Println(err.Error())
	//	continue
	//}
	//
	//ext := filepath.Ext(t.m)
	//base, nameWithExt := filepath.Split(t.m)
	//
	//// name := strings.Replace(nameWithExt, ext, "", -1)
	//name := base + nameWithExt[:len(nameWithExt)-len(ext)]
	//
	//if opts.suffix != "" {
	//	simpleFile = name + opts.suffix + ext
	//} else {
	//	simpleFile = name + "_simple" + ext
	//}
	//
	//err = zhconvFactory.FileToSimple(t.m, simpleFile)
	//if err != nil {
	//	fmt.Println(err.Error())
	//	continue
	//}
	//
	//if len(opts.inputs) == 1 && len(opts.outputs) == 1 {
	//	outputFile = opts.outputs[0]
	//} else {
	//	ext := filepath.Ext(t.i)
	//	base, nameWithExt := filepath.Split(t.i)
	//	// name := strings.Replace(nameWithExt, ext, "", -1)
	//	name = base + nameWithExt[:len(nameWithExt)-len(ext)]
	//
	//	if opts.suffix != "" {
	//		outputFile = name + opts.suffix + ext
	//	} else {
	//		outputFile = name + "_simple" + ext
	//	}
	//}
	//
	//switch t.fileType {
	//case "PDF":
	//	fmt.Println("再将HTML转换为PDF,这会花一些时间,请耐心等候")
	//
	//	if opts.wkhtml2pdf != "" {
	//
	//		if wkhtml2pdfFactory == nil {
	//			wkhtml2pdfInit()
	//		}
	//
	//		m, err := wkhtml2pdfFactory.NewMaterial(simpleFile, opts.outputDir, outputFile)
	//		if err != nil {
	//			fmt.Println(err.Error())
	//			continue
	//		}
	//
	//		p, err := wkhtml2pdfFactory.Convert(m)
	//		if err != nil {
	//			fmt.Println(err.Error())
	//			continue
	//		}
	//
	//		fmt.Println("[wkhtmltopdf]HTML->PDF 耗时:", p.Coast.String())
	//
	//		// remove old file
	//		os.Remove(t.m)
	//		os.Remove(simpleFile)
	//
	//		// merge pdf files
	//
	//		fmt.Println("转换完成:", p.FilePath)
	//
	//	} else {
	//
	//		// convert back to pdf
	//		if html2pdfFactory == nil {
	//			html2pdfInit()
	//		}
	//
	//		m, err := html2pdfFactory.NewMaterial(simpleFile, opts.outputDir, outputFile, opts.scale)
	//		if err != nil {
	//			fmt.Println(err.Error())
	//			continue
	//		}
	//
	//		p, err := html2pdfFactory.Convert(m)
	//		if err != nil {
	//			fmt.Println(err.Error())
	//			continue
	//		}
	//
	//		fmt.Println("[chrome]HTML->PDF 耗时:", p.Coast.String())
	//
	//		// remove old file
	//		os.Remove(t.m)
	//		os.Remove(simpleFile)
	//
	//		// merge pdf files
	//
	//		fmt.Println("转换完成:", p.FilePath)
	//	}
	//}
}

func toSimple() {

	if zhconvFactory == nil {
		zhconvInit()
	}

	for _, i := range opts.inputs {
		convertToSimple(i)
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
			m, err := html2pdfFactory.NewMaterial(t.m, opts.outputDir, outputFile, opts.scale, true)
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

func convertToSimple(f string) error {
	fPath, fName := filepath.Split(f)
	ext := filepath.Ext(fName)

	var outputFile string
	if len(opts.inputs) == 1 && len(opts.outputs) == 1 {
		outputFile = opts.outputs[0]
	} else {
		fNameNoneExt := fPath + fName[:len(fName)-len(ext)]
		if opts.suffix != "" {
			outputFile = fNameNoneExt + opts.suffix + ext
		} else {
			outputFile = fNameNoneExt + "_simple" + ext
		}
	}
	switch strings.ToUpper(ext) {
	case ".PDF":

		if pdf2htmlFactory == nil {
			pdf2htmlInit()
		}

		if opts.wkhtml2pdf != "" {
			if wkhtml2pdfFactory == nil {
				wkhtml2pdfInit()
			}
		} else {
			if html2pdfFactory == nil {
				html2pdfInit()
			}
		}
		pdfToSimple(f, outputFile)
	case "", ".HTML", ".TXT", ".JSON", ".XML", ".HTM":
		zhconvFactory.FileToSimple(f, outputFile)
	default:
		fmt.Println("不支持的文件类型")
	}

	return nil
}

func pdfToSimple(f, t string) error {

	// f 是输入的PDF
	// t 是输出的PDF

	// 检查文件是否存在
	if !util.Exists(f) {
		return errors.New("No such file to trans.")
	}

	// 切分出路径和文件名
	fpath, fname := filepath.Split(f)

	var suffix string
	if opts.suffix != "" {
		suffix = opts.suffix
	} else {
		suffix = "_simple"
	}

	// 后缀
	ext := filepath.Ext(fname)

	// 去掉文件后缀
	pureName := fname[:len(fname)-len(ext)]

	if t == "" {
		// 如果输出文件为空，在的当前目录下pureName+suffix+ext
		t = filepath.Join(fpath, pureName+opts.suffix+ext)
	}

	tmpDir := os.TempDir()
	if opts.tmpDir != "" {
		if p, err := filepath.Abs(opts.tmpDir); err == nil {
			tmpDir = p
		}
	}
	// 1. 建立存放切分文件和HTML文件以及转换后的PDF文件的临时目录
	// - tmp+pureName_pieces 		  # 切分后的PDF存放目录
	// - tmp+pureName_html			  # 转换至HTML的存放目录
	// - tmp+pureName_{suffix}_html	  # 繁体转简体的HTML存放目录
	// - tmp+pureName_{suffix}_pieces # HTML转PDF的存放目录
	// 不知道为什么会建立不成功文件夹，所以建两次
	tmpPiecesDir := filepath.Join(tmpDir, pureName+"_pieces")
	if util.Exists(tmpPiecesDir) {
		os.RemoveAll(tmpPiecesDir)
	}
	for {
		if err := os.MkdirAll(tmpPiecesDir, 0711); err == nil {
			break
		}
	}

	tmpHtmlDir := filepath.Join(tmpDir, pureName+"_html")
	if util.Exists(tmpHtmlDir) {
		os.RemoveAll(tmpHtmlDir)
	}
	for {
		if err := os.MkdirAll(tmpHtmlDir, 0711); err == nil {
			break
		}
	}

	tmpSimpleHtmlDir := filepath.Join(tmpDir, pureName+suffix+"_html")
	if util.Exists(tmpSimpleHtmlDir) {
		os.RemoveAll(tmpSimpleHtmlDir)
	}
	for {
		if err := os.MkdirAll(tmpSimpleHtmlDir, 0711); err == nil {
			break
		}
	}

	tmpSimplePiecesDir := filepath.Join(tmpDir, pureName+suffix+"_pieces")
	if util.Exists(tmpSimplePiecesDir) {
		os.RemoveAll(tmpSimplePiecesDir)
	}
	for {
		if err := os.MkdirAll(tmpSimplePiecesDir, 0711); err == nil {
			break
		}
	}

	// 6. 确保删除所有的临时目录和文件
	if !opts.debug {
		defer func() {
			os.RemoveAll(tmpPiecesDir)
			os.RemoveAll(tmpHtmlDir)
			os.RemoveAll(tmpSimpleHtmlDir)
			os.RemoveAll(tmpSimplePiecesDir)
		}()
	}

	// 2. 切分PDF文件
	pieces, piecesCount, err := sm.Split(f, tmpPiecesDir, opts.maxPage)
	if err != nil {
		return errors.New("Split pdf error")
	}
	fmt.Println("切为PDF文件的数量：", piecesCount)

	var wg1 sync.WaitGroup = sync.WaitGroup{}
	var wg2 sync.WaitGroup = sync.WaitGroup{}
	var ct1, ct2 int

	// 先启动HTML2PDF工厂
	go func() {
		for {
			p, _ := getFromHtml2PdfFactory()
			if p == nil {
				// must be channel been closed
				break
			}
			wg2.Done()
		}
	}()

	// 先启动PDF占HTML工厂
	go func() {
		// 去消费pdf2html的产品
		// 转为简体->送到HTML转PDF工厂
		for {
			// 这个方法是HANG方式
			pHtml, err := pdf2htmlFactory.Get()
			fmt.Println(" -> HTML ", pHtml.FilePath)

			if err == nil && pHtml != nil {
				// 数据转换正常

				// 取得HTML的文件名
				_, name := filepath.Split(pHtml.FilePath)

				// 转换为简体
				simpleHtmlFile := filepath.Join(tmpSimpleHtmlDir, name)
				err := zhconvFactory.FileToSimple(pHtml.FilePath, simpleHtmlFile)
				if err != nil {
					fmt.Println(err.Error())
				}

				// 加入HTML转PDF工厂
				addToHtml2PdfFactory(simpleHtmlFile, tmpSimplePiecesDir, name, opts.scale)

				// 任务计数+1
				wg2.Add(1)
				ct2++
				fmt.Printf("HTML -> PDF %d\n", ct2)
			} else {
				fmt.Println(err.Error())
			}
			wg1.Done()
		}
	}()

	// 3. 将所有PDF全部送到PDF转HTML工厂
	for _, p := range pieces {
		// 这里使用工厂模式，非阻塞

		_, name := filepath.Split(p)
		m, err := pdf2htmlFactory.NewMaterial(p, tmpHtmlDir, name)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		pdf2htmlFactory.Put(m)
		wg1.Add(1)
		ct1++
		fmt.Printf("===[%s]===\n", p)
		fmt.Printf("PDF -> HTML %d\n", ct1)
	}

	// 等待所有的HTML转换成PDF
	fmt.Println("等待所有转换完成")
	wg1.Wait()
	wg2.Wait()

	// 转换完马上关闭Tab有问题，在这里关闭Chrome
	if html2pdfFactory != nil {
		html2pdfFactory.Close()
	}

	// 5. 将PDF目录下的所有文件合并（注意顺序）
	n, err := sm.MergetFromDir(tmpSimplePiecesDir, t)
	fmt.Printf("将%d页合并: %s\n", n, t)

	return err
}
