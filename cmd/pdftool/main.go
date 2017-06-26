package main

import "fmt"

func startHttp() {
	fmt.Println("Start HTTP service on:", opts.http)
}

func close(){
	html2pdfClose()
	wkhtml2pdfClose()
	pdf2htmlClose()
	zhconvClose()
}

func main() {

	// pdf2htmlInit()
	// html2pdfInit()

	switch opts.cmd {
	case "server":
		startHttp()
	case "2pdf":
		toPdf()
	case "2html":
		toHtml()
	case "2simple":
		toSimple()
	case "2tradition":
		toTradition()
	}

	close()
}
