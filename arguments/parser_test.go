package arguments_test

import (
    "testing"
    "os"
    "goparse/arguments" // Adjust this import path accordingly
)

// setupParser for tests not requiring required arguments
func setupParser() *arguments.Parser {
    parser := arguments.NewParser()

    // Adding some basic arguments (none of these are required)
    parser.AddArgument("verbose", "v", "verbose", "Increase verbosity", "bool", false)
    parser.AddArgument("config", "c", "config", "Path to config file", "string", false)
    parser.AddArgument("retry-count", "r", "retry", "Retry count", "int", false)

    return parser
}

// Test string argument parsing
func TestStringArgument(t *testing.T) {
    parser := setupParser()

    // Mocking os.Args
    os.Args = []string{"program", "--config", "config.yaml"}

    // Parse the arguments
    parsedArgs, shouldExit, err := parser.Parse()

    // Check for error
    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }

    // Check exit condition
    if shouldExit {
        t.Fatalf("The program should not have signaled an exit.")
    }

    // Check parsed value
    if parsedArgs["config"] != "config.yaml" {
        t.Errorf("Expected config.yaml but got %v", parsedArgs["config"])
    }
}

// Test integer argument parsing
func TestIntArgument(t *testing.T) {
    parser := setupParser()

    // Mock os.Args with an integer flag
    os.Args = []string{"program", "--retry", "5"}

    // Parse the arguments
    parsedArgs, shouldExit, err := parser.Parse()

    // Check for error
    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }

    // Check exit condition
    if shouldExit {
        t.Fatalf("The program should not have signaled an exit.")
    }

    // Check parsed value
    if val, ok := parsedArgs["retry-count"].(int); !ok || val != 5 {
        t.Errorf("Expected retry to be 5 but got %v", parsedArgs["retry-count"])
    }
}

// Test boolean flag parsing
func TestBoolArgument(t *testing.T) {
    parser := setupParser()

    // Mock os.Args for boolean flag (verbose)
    os.Args = []string{"program", "--verbose"}

    // Parse the arguments
    parsedArgs, shouldExit, err := parser.Parse()

    // Check for error
    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }

    // Check exit condition
    if shouldExit {
        t.Fatalf("The program should not have signaled an exit.")
    }

    // Check parsed boolean flag
    if val, ok := parsedArgs["verbose"].(bool); !ok || !val {
        t.Errorf("Expected verbose to be true but got %v", parsedArgs["verbose"])
    }
}

// Test handling of missing required argument
func TestMissingRequiredArgument(t *testing.T) {
    parser := arguments.NewParser()

    // Add a required argument for this test
    parser.AddArgument("output", "o", "output", "Output file path", "string", true)  // Required argument

    // Mock os.Args without the required 'output' argument
    os.Args = []string{"program"}

    // Parse the arguments
    _, shouldExit, err := parser.Parse()

    // Check if error is returned due to missing required argument
    if err == nil {
        t.Fatalf("Expected an error for missing required argument, but got none")
    }

    // Check exit signal
    if !shouldExit {
        t.Fatalf("Expected 'shouldExit' to be true but got false")
    }

    // Verify the error message
    expectedError := "missing required global argument: output"
    if err.Error() != expectedError {
        t.Errorf("Expected error message %v but got %v", expectedError, err.Error())
    }
}

// Test parsing of multiple arguments
func TestMultipleArguments(t *testing.T) {
    parser := setupParser()

    // Mock os.Args with multiple valid flags
    os.Args = []string{"program", "--retry", "3", "--config", "myapp.yaml"}

    // Parse args
    parsedArgs, shouldExit, err := parser.Parse()

    // Check for error
    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }

    // Check exit condition
    if shouldExit {
        t.Fatalf("Program should not have signaled an exit.")
    }

    // Check multiple parsed values
    if parsedArgs["retry-count"] != 3 {
        t.Errorf("Expected retry-count to be 3 but got %v", parsedArgs["retry-count"])
    }
    if parsedArgs["config"] != "myapp.yaml" {
        t.Errorf("Expected config file to be 'myapp.yaml' but got %v", parsedArgs["config"])
    }
}

// Test handling missing argument value
func TestMissingArgumentValue(t *testing.T) {
    parser := setupParser()

    // Mock os.Args to simulate missing argument value
    os.Args = []string{"program", "--config"}  // Missing value for '--config'

    // Parse args
    _, shouldExit, err := parser.Parse()

    // Check if error due to missing value
    if err == nil {
        t.Fatalf("Expected an error due to missing value for argument, but got none")
    }

    // Verify exit signal
    if !shouldExit {
        t.Fatalf("Expected 'shouldExit' to be true but got false")
    }

    // Verify error message
    expectedError := "no value provided for argument --config"
    if err.Error() != expectedError {
        t.Errorf("Expected error message %v but got %v", expectedError, err.Error())
    }
}

// Test duplicate string flags (last one wins)
func TestRepeatedFlag(t *testing.T) {
    parser := setupParser()

    // Mock os.Args with repeated flag
    os.Args = []string{"program", "--config", "config1.yaml", "--config", "config2.yaml"}

    // Parse
    parsedArgs, shouldExit, err := parser.Parse()

    // Ensure no error
    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }

    // Check exit signal
    if shouldExit {
        t.Fatalf("Program should not have signaled an exit.")
    }

    // Again, the last flag's value should take precedence
    if parsedArgs["config"] != "config2.yaml" {
        t.Errorf("Expected config to be 'config2.yaml' but got %v", parsedArgs["config"])
    }
}

// Test duplicate boolean flag (should still be true)
func TestDuplicateBooleanFlag(t *testing.T) {
    parser := setupParser()

    // Mock os.Args -- pass '--verbose' twice
    os.Args = []string{"program", "--verbose", "--verbose"}

    // Parse args
    parsedArgs, shouldExit, err := parser.Parse()

    // Ensure no error
    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }

    // Check exit signal
    if shouldExit {
        t.Fatalf("Program should not have signaled an exit.")
    }

    // Duplicate boolean flag should still result in true
    if verbose, ok := parsedArgs["verbose"].(bool); !ok || !verbose {
        t.Errorf("Expected verbose to be true but got %v", parsedArgs["verbose"])
    }
}

// Test checking the default values for arguments not provided
func TestArgumentDefaults(t *testing.T) {
    parser := setupParser()

    // Mock os.Args without any flags provided
    os.Args = []string{"program"}

    // Parse the arguments
    parsedArgs, shouldExit, err := parser.Parse()

    // Ensure no error
    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }

    // Ensure program does not exit
    if shouldExit {
        t.Fatalf("Program should not have signaled an exit.")
    }

    // Verify defaults
    if retryCount := parsedArgs["retry-count"].(int); retryCount != 0 {
        t.Errorf("Expected retry-count to default to 0, but got %v", retryCount)
    }
    if config := parsedArgs["config"].(string); config != "" {
        t.Errorf("Expected config to default to empty string, but got %v", config)
    }
    if verbose := parsedArgs["verbose"].(bool); verbose != false {
        t.Errorf("Expected verbose to default to false, but got %v", verbose)
    }
}

// Test help flag triggers exit without error
func TestHelpFlag(t *testing.T) {
    parser := setupParser()

    // Mock os.Args with --help flag
    os.Args = []string{"program", "--help"}

    // Parse
    _, shouldExit, err := parser.Parse()

    // Ensure no error
    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }

    // Ensure program signals an exit
    if !shouldExit {
        t.Fatalf("Expected 'shouldExit' to be true due to help flag, but got false")
    }
}