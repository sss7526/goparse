// package main

// import (
// 	"fmt"
// 	"os"
// 	"github.com/sss7526/goparse"
// )

// func main() {
//     parser := goparse.NewParser(
// 		goparse.WithName("MyProgram"),
// 		goparse.WithVersion("v1.0.0"),
// 		goparse.WithAuthor("crab rangoon?"),
// 		goparse.WithDescription("A simple demonstration"),
// 	)

//     // Define global arguments
//     parser.AddArgument("verbose", "v", "verbose", "Increase verbosity", "bool", false)
//     parser.AddArgument("config", "c", "config", "Path to config file", "string", true)
//     parser.AddArgument("output", "o", "output", "Output file", "string", false)
//     parser.AddArgument("log", "l", "log", "Log file", "string", false)
// 	parser.AddArgument("many", "m", "many", "many opts", "[]string", false)
    
//     // Define mutually exclusive group (e.g., only one of 'output' or 'log' can be provided)
//     parser.AddExclusiveGroup([]string{"output", "log"}, false)

//     // Parse arguments
//     parsedArgs, shouldExit, err := parser.Parse()
//     if err != nil {
//         fmt.Fprintf(os.Stderr, "Error: %v\n", err)
//         if shouldExit {
//             os.Exit(1)
//         }
//     }

//     if shouldExit {
//         os.Exit(0)
//     }

//     if v, ok := parsedArgs["verbose"]; ok && v.(bool) {
// 		fmt.Println("Verbose mode enabled")
// 	}

// 	if configPath, ok := parsedArgs["config"].(string); ok && configPath != "" {
// 		fmt.Printf("Using config file: %s\n", configPath)
// 	}

// 	if manyOptions, ok := parsedArgs["many"].([]string); ok && len(manyOptions) > 0 {
// 		fmt.Println("Many options:")
// 		for _, option := range manyOptions {
// 			fmt.Printf("%s\n", option)
// 		}
// 	}
// }

package main

import (
    "fmt"

    // Import your argparser package here
    "github.com/sss7526/goparse" // Adjust this to your actual import path
)

func main() {
    // Create a new parser
    parser := goparse.NewParser(
        goparse.WithName("Sample CLI Program"), // Set program name
        goparse.WithDescription("This is a sample program demonstrating boolean and string flags."),
        goparse.WithVersion("1.0.0"),
        goparse.WithAuthor("Your Name"),
    )

    // Define flags and options
    parser.AddArgument("pelp", "p", "pelp", "Show pelp information", "bool", false)
    parser.AddArgument("verbose", "v", "verbose", "Enable verbose output", "bool", false)
    parser.AddArgument("force", "f", "force", "Force the action", "bool", false)
    parser.AddArgument("input", "i", "input", "Input file", "string", true)
    parser.AddArgument("output", "o", "output", "Output file (default: output.log)", "string", false, "output.log")
    parser.AddArgument("labels", "l", "labels", "Labels for the process (comma-separated)", "[]string", false)

    // Simulate getting arguments from command line
    // os.Args slice by convention includes the program name at index 0

    // Parse the arguments
    parsedArgs, _, err := parser.Parse()
    if err != nil {
        fmt.Println("Error:", err)
        parser.PrintHelp()
        return
    }

    // // If the user requests help, print help and exit
    // if parsedArgs["help"].(bool) {
    //     parser.PrintHelp()
    //     return
    // }

    // Display the values of the options
    fmt.Println("Parsed Arguments:")
	fmt.Printf("Show pelp: %v\n", parsedArgs["pelp"])
    fmt.Printf("Verbose mode: %v\n", parsedArgs["verbose"])
    fmt.Printf("Force option: %v\n", parsedArgs["force"])
    fmt.Printf("Input file: %v\n", parsedArgs["input"])
    fmt.Printf("Output file: %v\n", parsedArgs["output"])

    // Labels option is a []string, handle it carefully
    if labels, ok := parsedArgs["labels"].([]string); ok {
        fmt.Printf("Labels: %v\n", labels)
    } else {
        fmt.Println("Labels: None provided")
    }

}