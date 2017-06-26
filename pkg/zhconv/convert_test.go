package zhconv

import (
	"testing"
	"fmt"
)

func TestNewFactory(t *testing.T) {
	f := NewFactory()

	f.load("钟=中")

	fmt.Printf("%s\n", f.ToTraditional("中国人牛逼的一笔"))

	fmt.Printf("%s\n", f.ToSimple("中國人牛逼的一筆"))
}
