package zhconv

import (
	"errors"
	"github.com/jiusanzhou/pdf2html/pkg/util"
	"io/ioutil"
	"regexp"
	"strings"
)

var chineseRegex *regexp.Regexp = regexp.MustCompile("[\u4e00-\u9fa5]")
var ChineseRegex = chineseRegex
var DoubleCharRegex *regexp.Regexp = regexp.MustCompile("[\u00A1-\u9fa5]")

type Factory struct {
	s2t map[int32][]byte
	t2s map[int32][]byte
}

func NewFactory() *Factory {
	f := &Factory{
		s2t: make(map[int32][]byte),
		t2s: make(map[int32][]byte),
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
		f.s2t = make(map[int32][]byte)
		f.t2s = make(map[int32][]byte)
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
		f.t2s[int32([]rune(couple[0])[0])] = util.S2b(couple[1])
		f.s2t[int32([]rune(couple[1])[0])] = util.S2b(couple[0])
	}

	return nil
}

func (f *Factory) ToSimple(str string) (res []byte) {
	for _, s := range str {
		res = append(res, f.getSimple(s)...)
	}
	return
}

func (f *Factory) ToTraditional(str string) (res []byte) {
	for _, s := range str {
		res = append(res, f.getTraditional(s)...)
	}
	return
}

func (f *Factory) FileToSimple(source, dist string) (err error, precent float64) {

	str, err := ioutil.ReadFile(source)
	if err != nil {
		return
	}

	s := util.B2s(str)
	// 这里的逻辑处理得不够好
	// 但是为了业务职能加在这里了

	// 获取非字母数据和符号字符数据
	ds := DoubleCharRegex.FindAllString(s, -1)

	// 获取中文数据
	cs := ChineseRegex.FindAllString(strings.Join(ds, ""), -1)

	if len(ds) > 0 {
		precent = float64(len(cs)) / float64(len(ds))
	}

	res := f.ToSimple(s)
	err = ioutil.WriteFile(dist, res, 0666)

	return
}

func (f *Factory) FileToTraditional(source, dist string) (err error) {

	str, err := ioutil.ReadFile(source)
	if err != nil {
		return
	}

	res := f.ToTraditional(util.B2s(str))
	err = ioutil.WriteFile(dist, res, 0666)

	return
}

func (f *Factory) getTraditional(str int32) []byte {
	if !IsChineseInt(str) {
		return util.S2b(string(str))
	}

	v := f.s2t[str]
	if v != nil {
		return v
	} else {
		return util.S2b(string(str))
	}
}

func (f *Factory) getSimple(str int32) []byte {
	if !IsChineseInt(str) {
		return util.S2b(string(str))
	}

	v := f.t2s[str]
	if v != nil {
		return v
	} else {
		return util.S2b(string(str))
	}
}

func IsChinese(char string) bool {
	return chineseRegex.MatchString(char)
}

func IsChineseInt(str int32) bool {
	if 19968 <= str && str <= 40869 {
		return true
	} else {
		return false
	}
}
