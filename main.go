package main

import (
	"os"

	"github.com/tlinden/epuppy/cmd"
)

func main() {
	os.Exit(Main())
}

func Main() int {
	return cmd.Execute(os.Stdout)
}
