package main

import (
	"flag"
	"log"
)

type Options struct {
	outputDir string
	parallel  int
	exec      string
	execTpl   string

	input  string
	output string
}

var (
	defaultParallel = 5
)

func ParseArgs() *Options {

	outputDir := flag.String("output-dir", "", "Output directory path")
	parallel := flag.Int("parallel", defaultParallel, "Parallel to convert files")

	exe := flag.String("exec", "", "pdf2htmlEx execute file")
	exetpl := flag.String("exec-tpl", "", "Execute cmd")

	flag.Parse()

	var input, output string

	switch len(flag.Args()) {
	case 0:
		log.Fatalln("You must offer one file to convert.")
	case 1:
		input = flag.Arg(0)
	case 2:
		input = flag.Arg(0)
		output = flag.Arg(1)
	default:
		log.Fatalln("Too more arguments.")

	}

	return &Options{
		outputDir: *outputDir,
		parallel:  *parallel,

		exec:    *exe,
		execTpl: *exetpl,

		input:  input,
		output: output,
	}
}
