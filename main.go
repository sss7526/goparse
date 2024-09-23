package main

import (
	"fmt"
	"goparse/arguments"
	"os"
)

func main() {
    parser := arguments.NewParser("My Program","","","")

    // Define global arguments
    parser.AddArgument("verbose", "v", "verbose", "Increase verbosity", "bool", false)
    parser.AddArgument("config", "c", "config", "Path to config file", "string", true)
    parser.AddArgument("output", "o", "output", "Output file", "string", false)
    parser.AddArgument("log", "l", "log", "Log file", "string", false)
	parser.AddArgument("many", "m", "many", "many opts", "[]string", false)
    
    // Define mutually exclusive group (e.g., only one of 'output' or 'log' can be provided)
    parser.AddExclusiveGroup([]string{"output", "log"}, false)

    // Parse arguments
    parsedArgs, shouldExit, err := parser.Parse()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        if shouldExit {
            os.Exit(1)
        }
    }

    if shouldExit {
        os.Exit(0)
    }

    if v, ok := parsedArgs["verbose"]; ok && v.(bool) {
		fmt.Println("Verbose mode enabled")
	}

	if configPath, ok := parsedArgs["config"]; ok && configPath != "" {
		fmt.Printf("Using config file: %s\n", configPath)
	}
}