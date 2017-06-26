# PDF 工具

PDF工具支持4中功能：
1. HTML转换至PDF
2. PDF转换至HTML
3. 繁体转换至简体
4. 简体转换至繁体

## 部署依赖

1. pdf2htmlEX 工具，这个工具在压缩包中自带，依赖库也带上了，
但是不一定会在所有机器上运行成功，因为动态库版本的问题，如果有问题，
可以在本机编译。
2. 版本在59以上Chrome 浏览器，主要用于使用headless模式进行HTML转换至PDF
    也可以使用xfvb配合低版本的Chrome。

## 使用说明

目前server模式还未支持，请使用命令模式

命令总览：
```bash
  -auth string
    	server模式下，转换服务的认证用户名和密码
  -chrome string
    	Google浏览器的路径 (default "/usr/bin/chromium-browser")
  -data-dir string
    	数据路径，一般情况下里面会包含chrome和pdf2htmlEx (default ".data")
  -http string
    	server模式下，转换服务的HTTP监听地址 (default ":8080")
  -output string
    	繁<->简转换模式下，输出的文件名
  -output-dir string
    	转换文件的输出目录
  -pdf2html-exec string
    	pdf2htmlEx工具的路径 (default "pdf2htmlEx")
  -pdf2html-tpl string
    	pdf2htmlEx工具的执行模板 (default "{{exe}} --data-dir={{data}} {{input}} {{output}}")
  -scale float
    	HTML -> PDF 的缩放 (default 1)
  -suffix string
    	转换文件的后缀，只针对于繁<->简
  -tmp-dir string
    	临时文件路径
  -wkhtmltopdf-exec string
    	wkhtmltopdf工具的路径
  -wkhtmltopdf-tpl string
    	wkhtmltopdf工具的执行模板 (default "{{exec}} --no-stop-slow-scripts -g --print-media-type --load-error-handling ignore {{input}} {{output}}")

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
```

目录结构说明：
```bash
.
├── data
│   ├── pdf2htmlEX
│   │   ├── data
│   │   │   ├── base.css
│   │   │   ├── base.min.css
│   │   │   ├── compatibility.js
│   │   │   ├── compatibility.min.js
│   │   │   ├── fancy.css
│   │   │   ├── fancy.min.css
│   │   │   ├── LICENSE
│   │   │   ├── manifest
│   │   │   ├── pdf2htmlEX-64x64.png
│   │   │   ├── pdf2htmlEX.js
│   │   │   └── pdf2htmlEX.min.js
│   │   ├── libs
│   │   │   ├── d
│   │   │   │   ├── libfontforge.so.2
│   │   │   │   ├── libgioftp.so.2
│   │   │   │   ├── libgunicode.so.4
│   │   │   │   ├── libgutils.so.2
│   │   │   │   └── libuninameslist.so.1
│   │   │   ├── libspiro.so.0
│   │   │   ├── ......
│   │   └── pdf2htmlEX
│   ├── poppler-data
│   │   ├── ......
└── pdftool
```

`pdftool`是转换的工具，`data`是`pdf2htmlEX`及其依赖的字体文件`poppler-data`。

`data/pdf2htmlEX/libs/d`是`pdf2htmlEX`依赖的动态库文件路径，与上一层目录中构成是所有的动态依赖库，
如果运行提示缺少动态库的依赖，可手动移动`libs`中的文件至`libs/d`中。

繁体PDF转简体PDF：
```bash
export LD_LIBRARY_PATH=data/pdf2htmlEX/libs/d
./pdftool -pdf2html-exec=data/pdf2htmlEX/pdf2htmlEX -chrome=[chrome所在路径] 2simple test_data/02中工国际.PDF
```

最终等待一定时间：有`test_data/02中工国际_simple.PDF`文件生成。