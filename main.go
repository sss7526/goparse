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
	parser.AddArgument("urls", "u", "urls", "List of URLs", "[]string", false)
	parser.AddArgument("count", "c", "count", "A count of something", "int", false)
    
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

	if countVal, ok := args["count"].(int); ok {
		fmt.Printf("Count: %d\n", countVal)
	} else {
		fmt.Println("Could not retrieve count")
	}

	if urls, ok := args["urls"].([]string); ok && len(urls) > 0 {
		fmt.Printf("URSl: %v\n", urls)
	} else {
		fmt.Println("No URLs provided")
	}

    fmt.Printf("Parsed arguments: %v\n", args)
}