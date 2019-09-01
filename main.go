package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/omerh/awsctl/cmd/awsctl/cmd"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	cmd.Execute()
}

func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "\n%v\n", err)
	os.Exit(1)
}
