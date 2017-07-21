package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	unipdf "github.com/jiusanzhou/unidoc/pdf"
	// "github.com/jung-kurt/gofpdf"
	"bytes"
	"github.com/jiusanzhou/unidoc/common"
	"regexp"
	"github.com/signintech/gopdf"
	"github.com/jiusanzhou/pdf2html/pkg/zhconv"
	"strings"
)

func pFont(page *unipdf.PdfPage) {

	source := page.Resources
	if fonts, ok := source.Font.(*unipdf.PdfObjectDictionary); ok {
		for fontName, fontItem := range *fonts {
			if fontItemObj, ok := fontItem.(*unipdf.PdfIndirectObject); ok {
				fmt.Println(fontName, fontItemObj)
				if fontItemDict, ok := fontItemObj.PdfObject.(*unipdf.PdfObjectDictionary); ok {
					for k, v := range *fontItemDict {
						if k.String() == "ToUnicode" {
							//if oldStream := v.(*unipdf.PdfObjectStream); ok {
							//	bs, err := unipdf.DecodeStream(oldStream)
							//	fmt.Println("========", fontName, "=======")
							//	fmt.Println(string(bs))
							//	if err == nil {
							//		ds := reUnicode.ReplaceAllFunc(bs, func(s []byte) []byte {
							//			return []byte(" <4E8B>\n")
							//		})
							//		stream, _ := unipdf.EncodeStream(ds, "FlateDecode")
							//		fontItemDict.Set(k, stream)
							//	}
							//}
						} else if k.String() == "FontDescriptor" {
							if fontDesc, ok := v.(*unipdf.PdfIndirectObject); ok {
								if fontDescDict, ok := fontDesc.PdfObject.(*unipdf.PdfObjectDictionary); ok {
									fmt.Println(fontDescDict)
									if fontFile2, ok := (*fontDescDict)["FontFile2"]; ok {
										fmt.Println("=====Font File 2====")
										if oldStream := fontFile2.(*unipdf.PdfObjectStream); ok {
											bs, err := unipdf.DecodeStream(oldStream)
											if err == nil {
												f, _ := os.Create("font.ttf")
												f.Write(bs)
											}else{
												fmt.Println(err.Error())
											}
											// oldStream.Stream = []byte{}
										}else{
											fmt.Println("Not ok")
										}
									}
								}
							}
						} else {
							fmt.Println(k, v)
						}
					}
				}
			}
		}
	}
}

func pStream(stream *unipdf.PdfObjectStream) {
	bs, err := unipdf.DecodeStream(stream)
	if err == nil {
		fmt.Println(string(bs))
		// 02d6 05D40099
		// 05D40099
		// 02950A3703C10D360ABA05030DAB0FA003DE082202AE0321
		ds := bytes.Replace(bs, []byte("02950C820C1C0453"), []byte("02950A3703C10D360ABA05030DAB0FA003DE082202AE0321"), -1)
		en, err := unipdf.EncodeStream(ds, "FlateDecode")
		if err == nil {
			stream.Stream = en.Stream
		}
	} else {
		fmt.Println(err.Error())
	}
}

func pContent(page *unipdf.PdfPage) {
	fmt.Println("===Page===", page.Contents)
	switch v := page.Contents.(type) {
	case *unipdf.PdfObjectStream:
		pStream(v)
	case *unipdf.PdfObjectArray:
		for _, vv := range *v {
			if vvv, ok := vv.(*unipdf.PdfObjectStream); ok {
				pStream(vvv)
			}
		}
	}
}

var reUnicode *regexp.Regexp = regexp.MustCompile(" <[0-9A-F]{4}>\n")

func main() {


	text := "微信抓取账号与手机等其他信sssssssssss;.~@#$$%^&*)(*&^}:\"><?息㟡㒂˪ᶲⶪ夷⇯˫䫔Ḵ⋩"

	fmt.Println("检查是否是乱码！")

	sss := zhconv.DoubleCharRegex.FindAllString(text, -1)
	fmt.Println(sss)
	fmt.Println("第一步：扣出双字节字符")

	ss := zhconv.ChineseRegex.FindAllString(strings.Join(sss, ""), -1)
	fmt.Println(ss)
	fmt.Println("第二步：扣出汉字")

	fmt.Println(float32(len(ss))/float32(len(sss)))
	fmt.Println("第三步：计算汉字在字表中的比例")
	fmt.Println("小于0.5可能就是乱码了")
	fmt.Println("退出程序！")
	//CreatePdf("你好啊中文")
	return
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("At least one PDF file.")
		// os.Exit(1)
		// args = append(args, "D:/Zoe/Projects/GO/src/github.com/jiusanzhou/pdf2html/cmd/parsepdf/_.pdf")
		args = append(args, "D:/Zoe/Projects/GO/src/github.com/jiusanzhou/pdf2html/cmd/parsepdf/1.pdf")
		args = append(args, "D:/Zoe/Projects/GO/src/github.com/jiusanzhou/pdf2html/cmd/pdftool/test_data/13北京科锐.PDF")
	}

	f, err := os.Open(args[0])
	if err != nil {
		log.Fatalln(err.Error())
	}

	defer f.Close()

	pdfReader, _ := unipdf.NewPdfReader(f)
	pdfWriter := unipdf.NewPdfWriter()

	common.SetLogger(common.ConsoleLogger{})
	fmt.Println(reUnicode)

	i := 0
	for _, page := range pdfReader.PageList {
		//pdfWriter.AddPage(page.GetPageAsIndirectObject())
		//break
		pFont(page)
		pContent(page)
		pdfWriter.AddPage(page.GetPageAsIndirectObject())
		i++
		if i > 2 {
			break
		}
	}

	n, err := os.Create("D:/Zoe/Projects/GO/src/github.com/jiusanzhou/pdf2html/cmd/parsepdf/_1.pdf")
	pdfWriter.Write(n)
}

func CreatePdf(content string) {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{ PageSize: gopdf.Rect{W: 595.28, H: 841.89}}) //595.28, 841.89 = A4
	pdf.AddPage()
	err := pdf.AddTTFFont("t", "D:\\Zoe\\Projects\\GO\\src\\github.com\\jiusanzhou\\pdf2html\\cmd\\parsepdf\\wts11.ttf")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = pdf.SetFont("t", "", 14)
	if err != nil {
		log.Print(err.Error())
		return
	}
	pdf.Cell(nil, content)
	pdf.WritePdf("D:\\Zoe\\Projects\\GO\\src\\github.com\\jiusanzhou\\pdf2html\\cmd\\parsepdf\\hello.pdf")
}
