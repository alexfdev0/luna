//go:build linux || darwin

package component

import (
	"plugin"
	"os"
	"fmt"
)

type Component *plugin.Plugin

func ReturnComponentFunction(_Component *plugin.Plugin, Name string) any {
	function, err := _Component.Lookup(Name)
	if err != nil {
		fmt.Println("luna-l2: failed sending '" + Name + "' to component:", err)
		return function
	}
	return function
}

func InitializeComponent(Path string) *plugin.Plugin {
	_Component, err := plugin.Open(Path)
	if err != nil {
		fmt.Println("luna-l2: failed to initialize component with path '" + Path + "':", err)
		os.Exit(1)
	}
	return Component(_Component)
}
