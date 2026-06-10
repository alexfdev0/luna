package lcc_info

import (
	"fmt"
)

const (
	VERSION string = "6.3"
)


func PrintVersionInfo() {
	fmt.Println("Luna Compiler Collection version " + VERSION)
	fmt.Println("Target: luna-l2")
}
