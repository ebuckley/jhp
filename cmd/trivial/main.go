package main

import (
	"fmt"
	"github.com/dop251/goja"
)

func main() {
	vm := goja.New()
	val, _ := vm.RunString(`
	var adder = function(a,b) {
		return a + b;
    }
	
		adder(98, 1);
	`)
	fmt.Println(val)
}
