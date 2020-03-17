package addrgetter

import (
	"fmt"
	"testing"
)

func TestDefaultAddrGetter(t *testing.T) {
	var g AddrGetter
	g = NewAddrListGetter([]string{"localhost:5460", "127.0.0.1:5460"})
	fmt.Println(g.GetAddr())
	fmt.Println(g.GetAddr())
	g = NewAddrFuncGetter(func() string {
		return "[::1]:5460"
	})
	fmt.Println(g.GetAddr())
}
