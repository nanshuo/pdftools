package util

import "strings"

func ExecTpl(tpl string, kvs map[string]string) string {
	rs := []string{}
	for k, v := range kvs {
		rs = append(rs, "{{"+k+"}}", v)
	}
	r := strings.NewReplacer(rs...)
	return r.Replace(tpl)
}
