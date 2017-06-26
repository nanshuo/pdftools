package wkhtml2pdf

type Config struct {
	// output directory
	OutputDir string

	// auto rename if has the same file name
	// if not replace the old file with the newer
	AutoRename bool

	// max concurrent
	Concurrent int

	// main execute file
	// default: wkhtml2pdf
	Exec string

	// template of command lines.
	ExecTpl string
}
