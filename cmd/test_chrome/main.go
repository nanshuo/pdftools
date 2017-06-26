package main

import (
	"encoding/base64"
	"github.com/jiusanzhou/gcd"
	"github.com/jiusanzhou/gcd/gcdapi"
	"io/ioutil"
	"log"
	"fmt"
	"os"
)

func main() {
	chrome := gcd.NewChromeDebugger()
	defer chrome.ExitProcess()

	chrome.AddFlags([]string{"--disable-gpu", "--headless"})

	chromePath := "C:\\Program Files (x86)\\GoogleChromePortable\\App\\Chrome-bin\\chrome.exe"
	chromePath = "D:\\Program Files\\chrome-win32\\chrome_test.exe"

	chrome.StartProcess(chromePath, os.TempDir(), "8890")

	tab, err := chrome.NewTab()

	if err != nil {
		log.Fatalln(err.Error())
	}

	tab.Page.Enable()

	navigatedCh := make(chan struct{})

	tab.Subscribe("Page.domContentEventFired", func(targ *gcd.ChromeTarget, payload []byte) {
		fmt.Println("load")
		navigatedCh <- struct{}{}
	})

	_, err = tab.Page.Navigate("D:\\Zoe\\Projects\\GO\\src\\github.com\\jiusanzhou\\pdf2html\\cmd\\pdf2html\\test_data\\11湘鄂情.html", "", "")
	if err != nil {
		log.Fatalln(err.Error())
	}

	params := &gcdapi.PagePrintToPDFParams{
		MarginBottom: 0,
		MarginLeft:   0,
		MarginRight:  0,
		MarginTop:    0,

		PrintBackground: true,

		DisplayHeaderFooter:true,

		// Landscape: true,
	}

	// wait for navigation to finish
	<-navigatedCh

	data, err := tab.Page.PrintToPDFWithParams(params)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	bs, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = ioutil.WriteFile("test.pdf", bs, 0666)
	if err != nil {
		log.Fatalln(err.Error())
	}

	fmt.Println("转换完成啦~")
}
