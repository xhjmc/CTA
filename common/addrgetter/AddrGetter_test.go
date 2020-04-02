package addrgetter_test

import (
	"cta/common/addrgetter"
	"fmt"
	"testing"
)

func TestDefaultAddrGetter(t *testing.T) {
	var g addrgetter.AddrGetter
	g = addrgetter.NewAddrListGetter([]string{"localhost:5460", "127.0.0.1:5460"})
	fmt.Println(g.GetAddr())
	fmt.Println(g.GetAddr())
	g = addrgetter.NewAddrFuncGetter(func() string {
		return "[::1]:5460"
	})
	fmt.Println(g.GetAddr())
}
