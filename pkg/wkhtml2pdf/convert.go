package wkhtml2pdf

import (
	"bytes"
	"fmt"
	"github.com/jiusanzhou/pdf2html/pkg/util"
	"io/ioutil"
	"os"
	"path"
	"time"
)

// material of factory input
type Material struct {
	// TODO: support http get file
	Url string

	// file path
	FilePath string

	// output file path
	OutputFilePath string
}

// product of factory output
type Product struct {
	// status of this convert
	// 0: normal
	// 1: error
	// ...
	Status int

	// file path of output file
	FilePath string

	// size of out file
	Size int64

	// coast time of born
	Coast time.Duration

	// Material
	Material *Material
}

// factory for transforming pdf to html
type Factory struct {
	// configuration
	config *Config

	// execute path of pdf2htmlEx
	cmdTpl string

	// material channel for put
	in chan *Material

	// product channel for get
	out chan *Product
}

var (
	defaultData = ".data"

	defaultExec = path.Join(defaultData, "wkhtmltopdf")

	defaultExecTpl = "{{exec}} --no-stop-slow-scripts -g --print-media-type --load-error-handling ignore {{input}} {{output}}"
)

func NewFactory(c *Config) (f *Factory, err error) {

	// exec := c.exec
	// execTpl := c.execTpl

	var exec, execTpl string

	if c.Exec != "" {
		exec = c.Exec
	} else {
		exec = defaultExec
	}

	if c.ExecTpl != "" {
		execTpl = c.ExecTpl
	} else {
		execTpl = defaultExecTpl
	}

	if c.OutputDir != "" {
		// TODO: make sure it writable
	}

	f = &Factory{
		config: c,
		cmdTpl: util.ExecTpl(execTpl, map[string]string{"exec": exec}),

		in:  make(chan *Material),
		out: make(chan *Product),
	}

	go f.Start()

	return
}

func (f *Factory) NewMaterial(filePath, outputDir, outputFileName string) (m *Material, err error) {

	// TODO: check if file path is url

	// TODO: check if file exits and is legal pdf file

	// TODO: check if output directory is writable

	// linux and windows is different

	var name, basePath, outputPath string
	if outputFileName == "" {
		basePath, name = path.Split(filePath)
	} else {
		name = outputFileName
	}

	name = name[:len(name)-len(path.Ext(name))] + ".pdf"

	if !path.IsAbs(name) {
		// replace suffix
		if outputDir != "" {
			outputPath = path.Join(outputDir, name)
		} else if f.config.OutputDir != "" {
			outputPath = path.Join(f.config.OutputDir, name)
		} else {
			outputPath = path.Join(basePath, name)
		}
	}

	m = &Material{
		FilePath:       filePath,
		OutputFilePath: outputPath,
	}

	return
}

func fixFontsBug(f string) {
	data, err := ioutil.ReadFile(f)
	if err == nil {
		ioutil.WriteFile(f, bytes.Replace(data, []byte("@font-face{font-family:"), []byte("@font-face{font-family:_remove_"), -1), 0666)
	}
}

func (f *Factory) Convert(m *Material) (p *Product, err error) {

	cmd := util.ExecTpl(f.cmdTpl, map[string]string{"input": m.FilePath, "output": m.OutputFilePath})

	fmt.Println(cmd)

	fixFontsBug(m.FilePath)

	p = &Product{
		Material: m,
	}

	startTime := time.Now()
	err = util.DoCmd(cmd)
	coast := time.Now().Sub(startTime)

	p.Coast = coast

	if err != nil {
		fmt.Println("转换HTML有错误,", err.Error())
		p.Status = 1
		return
	}

	fi, err := os.Stat(m.OutputFilePath)
	if err != nil {
		fmt.Println("HTML->PDF输出的文件发生错误,", err.Error())
		p.Status = 1
		return
	}

	p.Status = 0
	p.FilePath = m.OutputFilePath
	p.Size = fi.Size()
	return
}

func (f *Factory) Put(m *Material) (err error) {

	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()

	// log

	f.in <- m

	return
}

func (f *Factory) Get() (*Product, error) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()

	return <-f.out, nil
}

func (f *Factory) Start() {
	defer func() {
		if err := recover(); err != nil {

		}
	}()

	for m := range f.in {

		// get a material

		// convert
		p, _ := f.Convert(m)

		// log error

		// put product
		f.out <- p
	}
}

func (f *Factory) Close() {
	// TODO: wait finished all

	defer func() {
		if err := recover(); err != nil {
		}
	}()
	close(f.in)
	close(f.out)
}
