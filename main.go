package main

import (
    "fmt"
    "goparse/arguments"
)

func main() {
    parser := arguments.NewParser()

    // Adding some arguments
    parser.AddArgument("verbose", "v", "verbose", "Enable verbose mode", "bool", false)
    parser.AddArgument("output", "o", "output", "Output file", "string", false)
    
    // Adding a subcommand
    convertCmd := parser.AddCommand("convert", "Convert files between formats")
    convertCmd.Arguments = append(convertCmd.Arguments, &arguments.Argument{
        Name:        "input",
        Short:       "i",
        Long:        "input",
        DataType:    "string",
        Description: "Input file to convert",
    })

    // Parsing
    args, err := parser.Parse()
    if err != nil {
        fmt.Println(err)
        return
    }

    // Handle parsed arguments
    if verbose, ok := args["verbose"].(bool); ok && verbose {
        fmt.Println("Verbose mode enabled.")
    }
    fmt.Printf("Parsed arguments: %v\n", args)
}