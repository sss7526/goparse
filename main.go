package main

import (
    "fmt"
    "log"
    "goparse/arguments" // Assuming the argument parsing code is in a package named `arguments`
)

func main() {
    // Create a new parser
    parser := arguments.NewParser()

    // Add global arguments
    parser.AddArgument("verbose", "v", "verbose", "Enable verbose output", "bool", false).Required = false
    parser.AddArgument("config", "c", "config", "Path to configuration file", "string", true).Required = true

    // Add a subcommand for "build"
    buildCmd := parser.AddCommand("build", "Build the project")
    buildCmd.Arguments = append(buildCmd.Arguments,
        parser.AddArgument("optimize", "o", "optimize", "Enable optimizations", "bool", false),                   // Boolean flag
        parser.AddArgument("input", "i", "input", "Input files to build", "[]string", false),                    // List of strings
        parser.AddArgument("threads", "t", "threads", "Number of threads to use during build", "int", false),     // Integer argument
    )

    // // Add a subcommand for "deploy"
    deployCmd := parser.AddCommand("deploy", "Deploy the project")
    deployCmd.Arguments = append(deployCmd.Arguments,
        parser.AddArgument("target", "t", "target", "Deployment target", "string", false),   // Required argument
        parser.AddArgument("dryRun", "d", "dryrun", "Simulate deployment without changes", "bool", false),  // Boolean flag
    )

    // Parse the arguments from the command line
    parsedArgs, err := parser.Parse()
    if err != nil {
        log.Fatalf("Error parsing arguments: %v", err)
    }

    // Safely extract and print the results
    fmt.Printf("Global arguments:\n")
    fmt.Printf("  Verbose: %v\n", parsedArgs["verbose"])
    fmt.Printf("  Config: %v\n", parsedArgs["config"])

    // Handle subcommand-specific arguments
    if subcommandArgs, ok := parsedArgs["build"]; ok {
        fmt.Printf("Subcommand: build\n")
        fmt.Printf("  Optimize: %v\n", subcommandArgs.(map[string]interface{})["optimize"])
        fmt.Printf("  Input files: %v\n", subcommandArgs.(map[string]interface{})["input"])
        fmt.Printf("  Threads: %v\n", subcommandArgs.(map[string]interface{})["threads"])
    } else if subcommandArgs, ok := parsedArgs["deploy"]; ok {
        fmt.Printf("Subcommand: deploy\n")
        fmt.Printf("  Target: %v\n", subcommandArgs.(map[string]interface{})["target"])
        fmt.Printf("  Dry run: %v\n", subcommandArgs.(map[string]interface{})["dryRun"])
    } else {
        fmt.Printf("No subcommand recognized.\n")
    }
}