package lcc_info

import (
	"fmt"
)

const (
	VERSION string = "7.0 (preview)"
)


func PrintVersionInfo() {
	fmt.Println("Luna Compiler Collection version " + VERSION)
	fmt.Println("Target: luna-l2")
}
