package main

import (
	//"Me/tokens"
	"fmt"
	"os"
	//"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No input file provided.")
		fmt.Println("Usage: ./program_name <filename>")
		os.Exit(1)
	}
	if len(os.Args) > 2 {
		fmt.Println("Multiple files found.")
		fmt.Println("Usage: ./program_name <filename>")
		os.Exit(1)
	}
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("Error:",err)
		return
	}
	
}
