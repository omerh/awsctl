package main

import (
	"fmt"
	"os"
	"runtime"

	"gitlab.com/omerh/awsctl/aws/cmd"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	cmd.Execute()
}

func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "\n%v\n", err)
	os.Exit(1)
}
