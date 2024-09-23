# GoParse - Command-line Argument Parsing Library for Go

GoParse is a simple customizable command-line argument parsing library for Go. You can define options flags/arguments with various properties (name, description, required/optional status, data type, default values, etc.) and enforce options that are mutually exclusive. I made this more or less as a learning project and because I didn't like Go's default arg parsing syntax or the other, far more developed arg parsing libs out there (I'm just too lazy to read the documentation).

## Features

- **Simple Argument Definitions**: Support for short/long flags, description, defaults, and required flags.
- **Mutually Exclusive Argument Groups**: Ensures only one option from a group is passed.
- **Type-Safe Argument Parsing**: Automatically parses types such as `int`, `string`, and `[]string`.
- **Help and Version Output**: Provides automatic help (`--help`) and version (`--version`) support.
- **Graceful Error Handling**: Return value includes a `shouldExit` flag, leaving the program exit handling to the programmer.

## Installation

To include GoParse in your project, run:

```bash
go get -u github.com/yourusername/goparse
```

## Getting Started

### Basic Example (with Safe Argument Access and Exit Handling)

The following example shows how to define arguments, safely access them, and handle the `shouldExit` flag, allowing for graceful error handling:

```go
package main

import (
	"fmt"
	"os"
	"github.com/sss7526/goparse"
)

func main() {
	// Create a new parser with optional program metadata
	parser := goparse.NewParser(
		goparse.WithName("My CLI Tool"),
		goparse.WithDescription("A description of my CLI tool"),
		goparse.WithAuthor("The program author"),
		goparse.WithVersion("1.0.0"),
	)
	
	// Add a required argument.
	parser.AddArgument("input", "i", "input", "Input file path", "string", true)
	
	// Add an optional argument with default value.
	parser.AddArgument("verbose", "v", "verbose", "Enable verbose mode", "bool", false, false)

    // Add an option that takes one or more arguments.
    parser.AddArgument("manythings", "m", "manythings", "Takes space separated list of one or more strings", "[]string", false)



	// Parse the provided arguments.
	parsedArgs, shouldExit, err := parser.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing arguments:", err)
	}
    
	// `shouldExit` flag checks if we need to terminate the program due to help flag or error
	if shouldExit {
		if err != nil {
			// Error case, exit with failure code
			os.Exit(1)
		}
		// If help was printed, exit successfully
		os.Exit(0)
	}

	// Safely access/validate parsed arguments
	input := parsedArgs["input"].(string)
	verbose := parsedArgs["verbose"].(bool)

	fmt.Println("Input file:", input)
	fmt.Println("Verbose mode:", verbose)
}
```

### Running the Example

Assuming the program is compiled to `mycli`:

```bash
./mycli --input myfile.txt --verbose
```

Output:
```bash
Input file: myfile.txt
Verbose mode: true
```

### Handling Errors and `shouldExit` Properly

- **When the required `input` argument is missing**, the error is handled, the `shouldExit` flag is set, and the library returns control to your program:

```bash
$ ./mycli
Error parsing arguments: missing required global argument: input
```

- **When help is requested** (`-h` or `--help`), the library prints the help information, and the `shouldExit` flag is set to gracefully exit after displaying help:

```bash
$ ./mycli --help
```
Output:
```bash
My CLI Tool
Author: The program author
Version: 1.0.0
A description of my CLI tool

Usage:
    -i, --input: Input file path (string, required)
    -v, --verbose: Enable verbose mode (bool, optional)
    -m, --manythings: Takes space separated list of one or more strings
```

In both cases, returning `shouldExit` allows the user to manage the flow of the program without the library forcing a premature exit.

### Other Advanced Features

#### Adding Mutually Exclusive Arguments

GoParse allows mutually exclusive argument groups (only one option from a group may be provided):

```go
parser.AddArgument("foo", "f", "foo", "Foo option", "bool", false)
parser.AddArgument("bar", "b", "bar", "Bar option", "bool", false)

// Define that "foo" and "bar" are mutually exclusive.
parser.AddExclusiveGroup([]string{"foo", "bar"}, false)
```

This ensures that the user cannot pass both `--foo` and `--bar` at the same time.

#### Handling Different Data Types

GoParse manages various data types like `int`, `string`, `[]string`, `bool`:

```go
// String
parser.AddArgument("config", "c", "config", "Configuration file path", "string", false, "/default/config/path")

// Integer
parser.AddArgument("threads", "t", "threads", "Number of threads", "int", false)
```

Values are type-validated during parsing, ensuring robust error checking.

## API Reference

### `NewParser(options ...Option) *Parser`
Creates a new argument parser. You can pass optional configurations such as `WithName`, `WithDescription`, `WithAuthor`, and `WithVersion` for program metadata.

### `AddArgument(name, short, long, description, dataType string, required bool, defaultValue ...interface{}) *Argument`
Adds an argument to the parser:
- `name`: The internal argument name (used in the code).
- `short`: Short flag version (`-x`).
- `long`: Long flag version (`--example`).
- `description`: Description of the argument for help output.
- `dataType`: Argument type (`string`, `int`, `bool`, etc.).
- `required`: Set to `true` if the argument is required, otherwise `false`.
- `defaultValue`: (optional) Value used by default when the argument isn't provided.

### `AddExclusiveGroup(options []string, mustHave bool)`
Defines a group of mutually exclusive arguments:
- `options`: List of argument names in the mutual exclusion group.
- `mustHave`: Set to `true` if at least one option in the group must be provided.

### `PrintHelp()`
Prints the help message showing program metadata (name, version, description) and the usage instructions for all available arguments.

### `Parse(args []string) (map[string]interface{}, bool, error)`
Parses the provided command-line arguments:
- Returns a map of parsed arguments with their values.
- The `bool` flag (`shouldExit`) is set to `true` if the help flag was passed or an error occurred (indicating the program should exit).
- Returns an `error` if an issue is encountered (like missing required arguments or invalid types).

## Example Scenarios

### Run with Required Arguments:

```bash
./mycli -i input.txt
```

Output:
```bash
Input file: input.txt
Verbose mode: false
```

### Run with Help:

```bash
./mycli --help
```
Output:
```bash
My CLI Tool
A description of my CLI tool

Usage:
    -i, --input: Input file path (string, required)
    -v, --verbose: Enable verbose mode (bool, optional)
```

### Error Handling:

Missing required argument:
```bash
$ ./mycli
Error parsing arguments: missing required global argument: input
```

## Contributing

Contributions are welcome! Please follow these steps:
1. Fork the repository and make a new branch.
2. Make your changes in your branch.
3. Submit a pull request describing your changes.

## License

This project is licensed under the [MIT License](LICENSE).

---