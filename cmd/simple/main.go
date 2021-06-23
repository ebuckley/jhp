package main

import (
	"bufio"
	"fmt"
	"github.com/dop251/goja"
	"os"
)

func main() {
	fmt.Println("This is a very simple REPL. write some javascript here! use special command .exit to quit")

	vm := goja.New()

	reader := bufio.NewReader(os.Stdin)
	for  {
		fmt.Print("👽️")
		text, _ := reader.ReadString('\n')
		if text==".exit\n" {
			os.Exit(0)
		}
		val, err := vm.RunString(text)
		if err != nil {
			_, _ = fmt.Fprint(os.Stderr, "ERROR:", err, "\n")
		}
		fmt.Println(val)
	}
}