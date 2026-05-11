package main

import (
	"os"

	"gabe565.com/ambient-weather-fusion/cmd"
	"gabe565.com/utils/cobrax"
)

func main() {
	root := cmd.New(cobrax.WithVersion(""))
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
