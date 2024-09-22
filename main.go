package main

import (
	"fmt"
	"goparse/arguments"
	"os"
)

func main() {
	// Initialize the parser
	parser := arguments.NewParser()

	// Define global arguments
	parser.AddArgument("verbose", "v", "verbose", "Increase verbosity", "bool", false)
	parser.AddArgument("config", "c", "config", "Path to config file", "string", false, "/etc/config/yaml")
	parser.AddArgument("retry-count", "r", "retry", "Number of retries", "int", false)
	parser.AddArgument("required", "l", "required", "required arg", "string", true)

	// Parse arguments
	parsedArgs, shouldExit, err := parser.Parse()

	// Handle parsing errors
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		if shouldExit {
			os.Exit(1)
		}
	}

	// Handle help or exit flag
	if shouldExit {
		os.Exit(0)
	}

	// Process verbose mode if enabled
	if v, ok := parsedArgs["verbose"]; ok && v.(bool) {
		fmt.Println("Verbose mode enabled")
	}

	// Process config file if provided
	if configPath, ok := parsedArgs["config"].(string); ok && configPath != "" {
		fmt.Printf("Using config file: %s\n", configPath)
	}

	// Process retry count if provided
	if retryCount, ok := parsedArgs["retry-count"].(int); ok && retryCount > 0 {
		fmt.Printf("Retrying %d times\n", retryCount)
	}

	if required, ok := parsedArgs["required"].(string); ok && required != "" {
		fmt.Printf("Required flag passed\n")
	}

	fmt.Println("Program running...")
}
