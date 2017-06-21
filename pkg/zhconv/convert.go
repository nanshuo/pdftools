package zhconv

import (
	"errors"
	"github.com/jiusanzhou/pdf2html/pkg/util"
	"io/ioutil"
	"regexp"
	"strings"
)

var chineseRegex *regexp.Regexp = regexp.MustCompile("[\u4e00-\u9fa5]")

type Factory struct {
	s2t map[string]string
	t2s map[string]string
}

func NewFactory() *Factory {
	f := &Factory{
		s2t: make(map[string]string),
		t2s: make(map[string]string),
	}

	f.load(zh_db)

	return f
}

func (f *Factory) LoadResource(filename string, remove bool) error {

	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	if remove {
		f.s2t = make(map[string]string)
		f.t2s = make(map[string]string)
	}

	return f.load(util.B2s(bs))
}

func (f *Factory) load(str string) error {

	if len(str) == 0 {
		return errors.New("Resource load failed.")
	}

	pairs := strings.Split(str, "\n")
	for _, p := range pairs {
		couple := strings.SplitN(strings.TrimSpace(p), "=", 2)
		if len(couple) != 2 {
			continue
		}
		if !IsChinese(couple[0]) || !IsChinese(couple[1]) {
			continue
		}
		f.t2s[couple[0]] = couple[1]
		f.s2t[couple[1]] = couple[0]
	}

	return nil
}

// TODO: use int64 -> string

func (f *Factory) ToSimple(str string) (res string) {
	for _, s := range str {
		res += f.getSimple(string(s))
	}
	return
}

func (f *Factory) ToTraditional(str string) (res string) {
	for _, s := range str {
		res += f.getTraditional(string(s))
	}
	return
}

func (f *Factory) FileToSimple(source, dist string) (err error) {

	str, err := ioutil.ReadFile(source)
	if err != nil {
		return
	}

	res := f.ToSimple(util.B2s(str))
	err = ioutil.WriteFile(dist, util.S2b(res), 0666)

	return
}

func (f *Factory) FileToTraditional(source, dist string) (err error) {

	str, err := ioutil.ReadFile(source)
	if err != nil {
		return
	}

	res := f.ToTraditional(util.B2s(str))
	err = ioutil.WriteFile(dist, util.S2b(res), 0666)

	return
}

func (f *Factory) getTraditional(str string) string {
	if !IsChinese(str) {
		return str
	}

	v := f.s2t[str]
	if v != "" {
		return v
	} else {
		return str
	}
}

func (f *Factory) getSimple(str string) string {
	if !IsChinese(str) {
		return str
	}

	v := f.t2s[str]
	if v != "" {
		return v
	} else {
		return str
	}
}

func IsChinese(char string) bool {
	return chineseRegex.MatchString(char)
}
