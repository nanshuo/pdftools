package html2pdf

type Config struct {
	// path of chrome execute
	Chrome string

	// data dir
	TmpDir string

	// if always resident
	Permanent bool

	// port for dev tool
	Port string

	// output directory
	OutputDir string

	// other command lines
	Flags string
}
