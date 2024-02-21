package main

import (
	"github.com/kevinanthony/collection-keep-updater/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
