package main

import (
	"flag"
	"fmt"
	"github.com/jiusanzhou/pdf2html/pkg/zhconv"
	"github.com/jiusanzhou/eagle/util"
)

func main() {

	s := flag.Bool("string", false, "transform strings.")
	source := flag.String("resource", "", "traditional=simple pairs.")
	remove := flag.Bool("remove", false, "remove the old pairs.")
	toS := flag.Bool("2s", false, "transform to simple.")
	flag.Parse()

	f := zhconv.NewFactory()

	if *source != "" {
		f.LoadResource(*source, *remove)
	}

	if *s {
		for _, s := range flag.Args() {
			if *toS {
				fmt.Println(util.B2s(f.ToSimple(s)))
			} else {
				fmt.Println(util.B2s(f.ToTraditional(s)))
			}
		}
	} else {
		args := flag.Args()
		if len(args) == 2 {
			if *toS {
				f.FileToSimple(args[0], args[1])
			} else {
				f.FileToTraditional(args[0], args[1])
			}
		} else {
			fmt.Println("should and only need 2 file name")
		}
	}

}
