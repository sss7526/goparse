// parser_test.go
package main

import (
    "os"
    "testing"
	"goparse/arguments"
)

func TestParser_NoArgs(t *testing.T) {
    // Simulate no arguments
    os.Args = []string{"app"}

    parser := arguments.NewParser()
    _, err := parser.Parse()
    if err == nil {
        t.Errorf("Expected help text error due to missing arguments")
    }
}

func TestParser_BasicArgs(t *testing.T) {
    // Simulate user-provided arguments
    os.Args = []string{"app", "--verbose", "--output", "file.txt"}

    parser := arguments.NewParser()
    parser.AddArgument("verbose", "v", "verbose", "Enable verbose mode", "bool", false)
    parser.AddArgument("output", "o", "output", "Output file", "string", false)

    args, err := parser.Parse()
    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }

    // Check parsed values
    if args["output"] != "file.txt" {
        t.Errorf("Expected output to be 'file.txt', got %v", args["output"])
    }
    if !args["verbose"].(bool) {
        t.Errorf("Expected verbose argument to be true")
    }
}