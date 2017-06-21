package pdf2html

type Config struct {
	// output directory
	OutputDir string

	// auto rename if has the same file name
	// if not replace the old file with the newer
	AutoRename bool

	// max concurrent
	Concurrent int

	// main execute file
	// default: pdf2htmlEX
	Exec string

	// template of command lines.
	ExecTpl string

	// data dir
	Pdf2htmlDataDir string
}
