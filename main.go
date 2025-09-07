package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/joho/godotenv"
	"github.com/omerh/awsctl/cmd/awsctl/cmd"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		// .env file is optional, so we don't exit on error
		log.Println("No .env file found or error loading .env file")
	}

	cmd.Execute()
}

func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "\n%v\n", err)
	os.Exit(1)
}
