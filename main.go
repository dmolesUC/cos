package main

import (
	"fmt"
	"os"

	"github.com/dmolesUC3/cos/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
