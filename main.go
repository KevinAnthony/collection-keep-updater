package main

import (
	"github.com/kevinanthony/collection-keep-updater/cmd"
)

func main() {
	_ = cmd.GetCmd().Execute()
}
