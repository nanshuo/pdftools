package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
)

const usage string = `
高级使用说明: pdftool <command> [OPTIONS]
	pdftool server				启动HTTP转换服务，通过HTTP接口接收转换请求
	pdftool 2pdf				转换文件为PDF文档
	pdftool 2html				转换文件为HTML文档
	pdftool 2simple				转换文档为简体中文
	pdftool 2tradition			转换文档为繁体中文
	pdftool help				打印帮助文档
	pdftool version				查看版本信息

示例：
	pdftool server -http=:8080 -auth=username:password
	pdftool 2pdf test.html test.pdf
	pdftool 2pdf test-1.html test-2.html
	pdftool 2pdf -chrome=./chrome/chrome test.html test.pdf
	pdftool 2pdf -output-dir=output_test test.html
	pdftool 2html test.pdf test.html
	pdftool 2html -pdf2htmlEx=data/pdf2htmlEx/pdf2htmlEx test.pdf test.html
	pdftool 2html test-1.pdf test-2.pdf
	pdftool 2html -output-dir=output_test test.pdf
	pdftool 2simple -suffix=_simple zh_tradition.txt zh_tradition.html
	pdftool 2tradition -output=zh_tradition_simple.txt zh_tradition.txt
	pdftool version
`

type Options struct {
	// command type
	// enum: server, 2pdf, 2html, 2simple, 2tradition
	cmd string

	// http serve address for server mode only
	http string

	// http authorization for server mode only
	auth string

	// data directory, maybe contains all files we need
	dataDir string

	// directory for output files
	outputDir string

	// execute file of pdf2htmlEx
	// if we assign this, we will not find it in data dir
	pdf2htmlEx string
	// command line tpl
	pdf2htmlExTpl string

	// wkhtml2pdf
	wkhtml2pdf    string
	wkhtml2pdfTpl string

	// execute file of chrome
	// if we assign this, we will not find it in data dir
	chrome string

	// input files
	inputs []string

	// output files
	// if the length of outputs don't equal inputs'
	// give up the output, and auto named the output files with suffix.
	outputs []string

	// suffix for outputs
	// will append after file-name and before ext-name
	// like: test.html -> test_suffix.html
	suffix string

	// tmp dir
	tmpDir string

	scale float64

	maxPage int

	debug bool
}

var opts *Options

func init() {

	var (
		chrome string

		wkhtml2pdf    string
		wkhtml2pdfTpl string

		pdf2htmlEx    string
		pdf2htmlExTpl string

		dataDir   string
		outputDir string

		suffix string

		http string
		auth string

		tmpDir string

		output string

		scale float64

		maxPage int

		debug bool
	)

	flag.Usage = func() {
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, usage)
	}

	flag.StringVar(&wkhtml2pdf, "wkhtmltopdf-exec", "", "wkhtmltopdf工具的路径")
	flag.StringVar(&wkhtml2pdfTpl,
		"wkhtmltopdf-tpl",
		"{{exec}} --no-stop-slow-scripts -g --print-media-type --load-error-handling ignore {{input}} {{output}}",
		"wkhtmltopdf工具的执行模板")
	flag.StringVar(&pdf2htmlEx, "pdf2html-exec", "pdf2htmlEx", "pdf2htmlEx工具的路径")
	flag.StringVar(&pdf2htmlExTpl, "pdf2html-tpl", "{{exe}} --data-dir={{data}} {{input}} {{output}}", "pdf2htmlEx工具的执行模板")
	flag.StringVar(&dataDir, "data-dir", ".data", "数据路径，一般情况下里面会包含chrome和pdf2htmlEx")
	flag.StringVar(&outputDir, "output-dir", "", "转换文件的输出目录")
	flag.StringVar(&suffix, "suffix", "", "转换文件的后缀，只针对于繁<->简")
	flag.StringVar(&http, "http", ":8080", "server模式下，转换服务的HTTP监听地址")
	flag.StringVar(&auth, "auth", "", "server模式下，转换服务的认证用户名和密码")
	flag.StringVar(&output, "output", "", "繁<->简转换模式下，输出的文件名")
	flag.Float64Var(&scale, "scale", 1, "HTML -> PDF 的缩放")
	flag.IntVar(&maxPage, "max-page", 50, "一次转换最多的页数，用此参数可控制并发。不是越小越好啊，越好越占CPU哦")
	flag.BoolVar(&debug, "debug", false,"")

	switch runtime.GOOS {
	case "windows":
		flag.StringVar(&chrome, "chrome", "C:\\Program Files (x86)\\Google\\Chrome\\Application\\chrome.exe", "Google浏览器的路径")
		flag.StringVar(&tmpDir, "tmp-dir", "", "临时文件路径")
	case "darwin":
		flag.StringVar(&chrome, "chrome", "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", "Google浏览器的路径")
		flag.StringVar(&tmpDir, "tmp-dir", "", "临时文件路径")
	case "linux":
		flag.StringVar(&chrome, "chrome", "/usr/bin/chromium-browser", "Google浏览器的路径")
		flag.StringVar(&tmpDir, "tmp-dir", "", "临时文件路径")
	}

	flag.Parse()

	if tmpDir == "" {
		tmpDir = os.TempDir()
	}

	opts = &Options{
		cmd:       flag.Arg(0),
		http:      http,
		auth:      auth,
		outputDir: outputDir,
		dataDir:   dataDir,
		suffix:    suffix,
		chrome:    chrome,

		wkhtml2pdf:    wkhtml2pdf,
		wkhtml2pdfTpl: wkhtml2pdfTpl,

		tmpDir:        tmpDir,
		pdf2htmlEx:    pdf2htmlEx,
		pdf2htmlExTpl: pdf2htmlExTpl,

		scale:   scale,
		maxPage: maxPage,

		debug: debug,
	}

	if len(flag.Args()) == 0 {
		fmt.Println("你要干嘛，不知道就 -h 看帮助")
		os.Exit(1)
	}

	args := flag.Args()[1:]
	switch len(args) {
	case 0:
		switch opts.cmd {
		case "2pdf", "2html", "2simple", "2tradition":
			fmt.Println("要转换文件，你起码要输入一个文件名的吧!")
			os.Exit(1)
		}
	case 1:
		switch opts.cmd {
		case "server":
			fmt.Println("[警告]你准备在HTTP服务模式下，不需要带命令!")
		case "2pdf":
			ext := strings.ToUpper(path.Ext(args[0]))
			if ext != "" && ext != ".HTML" {
				fmt.Println("不好意思，目前只支持从HTML文档转PDF，其他格式的你可以贡献代码!")
				fmt.Println("你非得强行转的话，你可以把后缀名去掉试试!")
				os.Exit(1)
			}
		case "2html":
			ext := strings.ToUpper(path.Ext(args[0]))
			if ext != "" && ext != ".PDF" {
				fmt.Println("不好意思，目前只支持从PDF文档转HTML，其他格式的你可以贡献代码!")
				fmt.Println("你非得强行转的话，你可以把后缀名去掉试试!")
				os.Exit(1)
			}
		case "2simple", "2tradition":
			if output != "" {
				opts.outputs = append(opts.outputs, output)
			}
		}
		opts.inputs = append(opts.inputs, args[0])
	case 2:
		switch opts.cmd {
		case "server":
			fmt.Println("[警告]你准备在HTTP服务模式下，不需要带命令!")
		case "2pdf":
			ext := strings.ToUpper(path.Ext(args[0]))
			if ext != "" && ext != ".HTML" {
				fmt.Println("不好意思，目前只支持从HTML文档转PDF，其他格式的你可以贡献代码!")
				fmt.Println("你非得强行转的话，你可以把后缀名去掉试试!")
				os.Exit(1)
			}
			opts.inputs = append(opts.inputs, args[0])
			opts.outputs = append(opts.outputs, args[1])
		case "2html":
			ext := strings.ToUpper(path.Ext(args[0]))
			if ext != "" && ext != ".PDF" {
				fmt.Println("不好意思，目前只支持从PDF文档转HTML，其他格式的你可以贡献代码!")
				fmt.Println("你非得强行转的话，你可以把后缀名去掉试试!")
				os.Exit(1)
			}
			opts.inputs = append(opts.inputs, args[0])
			opts.outputs = append(opts.outputs, args[1])
		case "2simple", "2tradition":
			opts.inputs = append(opts.inputs, args[0])
			opts.outputs = append(opts.outputs, args[1])
		}
	default:
		switch opts.cmd {
		case "server":
			fmt.Println("[警告]你准备在HTTP服务模式下，不需要带命令!")
		default:
			fmt.Println("这么多我就都当做你是要进行转换的哦，没有检查类型，自行确认文件正确呀!")
			opts.inputs = append(opts.inputs, args...)
			opts.outputs = []string{}
		}
	}

	switch opts.cmd {
	case "version":
		fmt.Println("暂时还不想去定版本，有疑问请找我")
	case "help":
		flag.Usage()
		os.Exit(0)
	case "server":
	case "2pdf":
	case "2html":
	case "2simple":
	case "2tradition":
	default:
		fmt.Printf("未知命令: %s，看下你是不是打错了!\n", opts.cmd)
		os.Exit(1)
	}
}
