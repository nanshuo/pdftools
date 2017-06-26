package html2pdf

import (
	"encoding/base64"
	"fmt"
	"github.com/CodisLabs/codis/pkg/utils/errors"
	"github.com/jiusanzhou/gcd"
	"github.com/jiusanzhou/gcd/gcdapi"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	minVersion = 60
)

// material of factory input
type Material struct {
	// TODO: support http get file
	Url string

	// file path
	FilePath string

	// output file path
	OutputFilePath string

	Params *gcdapi.PagePrintToPDFParams
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
	// chrome debugger
	// use: github.com/wirepair/gcd
	chrome *gcd.Gcd

	// configuration
	config *Config

	// material channel for put
	in chan *Material

	// product channel for get
	out chan *Product
}

func NewFactory(c *Config) (f *Factory, err error) {

	var chromePath, tmpDir, port string
	if c.Chrome == "" {
		switch runtime.GOOS {
		case "windows":
			chromePath = "C:\\Program Files (x86)\\Google\\Chrome\\Application\\chrome.exe"
		case "linux":
			chromePath = "/usr/bin/chromium-browser"
		case "darwin":
			chromePath = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
		}
	} else {
		chromePath = c.Chrome
	}

	if c.TmpDir == "" {
		switch runtime.GOOS {
		case "windows":
			tmpDir = "C:\\temp\\"
		case "linux", "darwin":
			tmpDir = "/tmp/"
		}
	} else {
		tmpDir = c.TmpDir
	}

	f = &Factory{
		chrome: gcd.NewChromeDebugger(),
		config: c,

		in:  make(chan *Material),
		out: make(chan *Product),
	}

	if c.Port == "" {
		port = "9999"
	} else {
		port = c.Port
	}

	// check version of chrome
	manifests, _ := filepath.Glob(path.Join(filepath.Dir(chromePath), "*.manifest"))
	if len(manifests) > 0 {
		// get the biggest one of manifest
		var v int
		_, name := filepath.Split(manifests[0])
		v, err = strconv.Atoi(strings.Split(name, ".")[0])
		if err == nil && v < minVersion {
			fmt.Println("Chrome version is too old,", manifests[0], v)
			err = errors.New("Chrome version is too old.")
			return
		}
	}

	f.chrome.StartProcess(chromePath, tmpDir, port, "--headless", "--disable-gpu")
	// f.chrome.AddFlags([]string{"--headless", "--disable-gpu", c.Flags})

	go f.Start()

	return
}

func (f *Factory) NewMaterial(filePath, outputDir, outputFileName string, scale float64) (m *Material, err error) {

	// TODO: check if file path is url

	// TODO: check if file exits and is legal pdf file

	// TODO: check if output directory is writable

	var name, outputPath string
	if outputFileName == "" {
		_, name = path.Split(filePath)
	} else {
		name = outputFileName
	}

	// replace suffix
	name = name[:len(name)-len(path.Ext(name))] + ".pdf"

	if outputDir != "" {
		outputPath = path.Join(outputDir, name)
	} else {
		outputPath = path.Join(f.config.OutputDir, name)
	}

	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return
	}
	m = &Material{
		FilePath:       absPath,
		OutputFilePath: outputPath,
		Params:         &gcdapi.PagePrintToPDFParams{Scale: scale},
	}

	return
}

func (f *Factory) Convert(m *Material) (p *Product, err error) {

	p = &Product{
		Material: m,
	}

	startTime := time.Now()

	// TODO: use channel to deal with tab
	tab, err := f.chrome.NewTab()
	if err != nil {
		p.Status = 1
		return
	}

	defer f.chrome.CloseTab(tab)

	tab.Page.Enable()

	navigatedCh := make(chan struct{})

	// subscribe to page loaded event
	tab.Subscribe("Page.loadEventFired", func(targ *gcd.ChromeTarget, payload []byte) {
		navigatedCh <- struct{}{}
	})

	_, err = tab.Page.Navigate(m.FilePath, "", "")
	if err != nil {
		p.Status = 1
		return
	}

	if m.Params == nil {
		m.Params = &gcdapi.PagePrintToPDFParams{
			MarginBottom: 0,
			MarginLeft:   0,
			MarginRight:  0,
			MarginTop:    0,
			Scale:        1,
		}
	}

	// wait for navigation to finish
	<-navigatedCh

	data, err := tab.Page.PrintToPDFWithParams(m.Params)
	if err != nil {
		return
	}

	bs, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		p.Status = 1
		return
	}

	err = ioutil.WriteFile(m.OutputFilePath, bs, 0666)
	if err != nil {
		p.Status = 1
		return
	}

	fi, err := os.Stat(m.OutputFilePath)
	if err != nil {
		p.Status = 1
		return
	}

	p.Coast = time.Now().Sub(startTime)
	p.Status = 0
	p.Size = fi.Size()
	p.FilePath = m.OutputFilePath

	return
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
	// f.chrome.ExitProcess()

	defer func() {
		if err := recover(); err != nil {
		}
	}()
	close(f.in)
	close(f.out)
}
