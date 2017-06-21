package pdf2html

import (
	"testing"
	"fmt"
)

func TestNewFactory(t *testing.T) {
	c := &Config{}
	f, _:=NewFactory(c)
	fmt.Println(*f)
}

func TestFactory_NewMaterial(t *testing.T) {
	c := &Config{
		OutputDir: "data_html_output",
	}
	f, _:=NewFactory(c)

	m, _:= f.NewMaterial("testdsadasdsad.pdf", "", "")
	fmt.Println(*m)
}

func TestFactory_Convert(t *testing.T) {
	c := &Config{}
	f, _:=NewFactory(c)

	m, _:= f.NewMaterial("testdsadasdsad.pdf", "", "")


	p, err := f.Convert(m)
	fmt.Println(*p, err.Error())
}
