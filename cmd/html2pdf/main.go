package main

import (
	"flag"
	"github.com/jiusanzhou/pdf2html/pkg/html2pdf"
	"log"
	"runtime"
	"os"
	"os/signal"
	"syscall"
	"fmt"
)

var path string
var dir string
var port string
var outputDir string

func init() {
	switch runtime.GOOS {
	case "windows":
		flag.StringVar(&path, "chrome", "C:\\Program Files (x86)\\Google\\Chrome\\Application\\chrome.exe", "path to chrome")
		flag.StringVar(&dir, "dir", "C:\\temp\\", "user directory")
	case "darwin":
		flag.StringVar(&path, "chrome", "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", "path to chrome")
		flag.StringVar(&dir, "dir", "/tmp/", "user directory")
	case "linux":
		flag.StringVar(&path, "chrome", "/usr/bin/chromium-browser", "path to chrome")
		flag.StringVar(&dir, "dir", "/tmp/", "user directory")
	}

	flag.StringVar(&port, "port", "9999", "Debugger port")
	flag.StringVar(&outputDir, "output-dir", "", "Output directory")

	flag.Parse()
}

func main() {
	c := &html2pdf.Config{
		Port:   port,
		TmpDir: dir,
		Chrome: path,
	}

	args := flag.Args()
	var i, o string
	switch len(args) {
	case 1:
		i = args[0]
	case 2:
		i = args[0]
		o = args[1]
	default:
		log.Fatalln("You only need input and output.")
	}

	f, err := html2pdf.NewFactory(c)

	if err != nil {
		log.Fatalln("Init factory error, ", err.Error())
	}

	defer f.Close()


	// Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal)

	go func(){
		m, err := f.NewMaterial(i, o, outputDir)
		if err != nil {
			log.Fatalln(err.Error())
		}
		p, err := f.Convert(m)

		if err != nil {
			log.Fatalln(err.Error())
		}
		log.Printf("%s 转换完成，耗时: %s，文件大小: %d\n", p.FilePath, p.Coast.String(), p.Size)
		ch <-syscall.SIGINT
	}()

	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	<-ch
	fmt.Println("关闭Chrome")
	f.Close()
}
