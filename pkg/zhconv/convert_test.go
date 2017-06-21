package zhconv

import (
	"testing"
	"fmt"
)

func TestNewFactory(t *testing.T) {
	f := NewFactory()

	fmt.Println(f.ToTraditional("中国人牛逼的一笔"))
}
