package goparse

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

)

// Argument represents an argument (flag) definition
type Argument struct {
	Name			string
	Short			string
	Long			string
	Description		string
	DataType 		string 		// e.g., string, []string, int, bool, etc.
	DefaultValue 	interface{}
	Required		bool
}

type ExclusiveGroup struct {
	Options 		[]string	// Names of mutually exclusive options
	MustHave 		bool		// If true, exactly one option must be provided
}

// Option represents a functional option for configuring the parser with program metadata
type Option func(*Parser)

// Parser is the main type that handles argument parsing
type Parser struct {
	Name 			string // (Optional) program name
	Description		string // (Optional) program description
	Author			string // (Optional) program's author
	Version			string // (Optional) Program version
	args			[]*Argument
	exclusiveGroups	[]*ExclusiveGroup	
}


// Withname optionally sets the program name.
func WithName(name string) Option {
	return func(p *Parser) {
		p.Name = name
	}
}

// WithDescription optionally sets the program description
func WithDescription(desc string) Option {
	return func(p *Parser) {
		p.Description = desc
	}
}

// WithAuthor optionally sets the program author.
func WithAuthor(author string) Option {
	return func(p *Parser) {
		p.Author = author
	}
}

// WithVersion optionally sets the program version.
func WithVersion(version string) Option {
	return func(p *Parser) {
		p.Version = version
	}
}
// NewParser creates a new instance of the argument parser
func NewParser(options ...Option) *Parser {
	p := &Parser {
		args:				[]*Argument{},
		exclusiveGroups:	[]*ExclusiveGroup{},
	}

	// Apply all optionally provided function options
	for _, option := range options {
		option(p)
	}

	return p
}

// AddArgument adds a positional or optional argument to the parser
func (p *Parser) AddArgument(name, short, long, description string, dataType string, required bool, defaultValue ...interface{}) *Argument {
	arg := &Argument {
		Name:			name,
		Short:			short,
		Long:			long,
		Description:	description,
		DataType:		dataType,
		Required:		required,
	}

	if len(defaultValue) > 0 {
		arg.DefaultValue = defaultValue[0]
	}
	p.args = append(p.args, arg)
	return arg
}

func (p *Parser) AddExclusiveGroup(optionNames []string, mustHave bool) {
	p.exclusiveGroups = append(p.exclusiveGroups, &ExclusiveGroup{
		Options: 		optionNames,
		MustHave:		mustHave,
	})
}

func (p *Parser) validateExclusiveGroups(parsedArgs map[string]interface{}) error {
	for _, group := range p.exclusiveGroups {
		foundCount := 0

		// Count how many mutually exclusive options are passed
		for _, optionName := range group.Options {
			if _, exists := parsedArgs[optionName]; exists {
				foundCount++
			}
		}

		// If more than one option in the group is passed, it's an error
		if foundCount > 1 {
			return fmt.Errorf("mutually exclusive options passed: %v", group.Options)
		}

		// If 'mustHave' is true but none were provided
		if group.MustHave && foundCount == 0 {
			return fmt.Errorf("one of the mutually exlusive options must be provided: %v", group.Options)
		}
	}
	return nil
}

func parseArguments(defs []*Argument, args []string, parsedArgs map[string]interface{}) error {
    for i := 0; i < len(args); i++ {
        arg := args[i]

        // Handle stacked short form flags (e.g., -abc => -a -b -c)
        if strings.HasPrefix(arg, "-") && !strings.HasPrefix(arg, "--") && len(arg) > 2 {
            for j := 1; j < len(arg); j++ {
                shortFlag := string(arg[j])
                found := false
                
                // Look for the short flag definition
                for _, def := range defs {
                    if def.Short == shortFlag && def.DataType == "bool" {
                        parsedArgs[def.Name] = true
                        found = true
                        break
                    }
                }

                if !found {
                    return fmt.Errorf("unknown argument: -%s", shortFlag)
                }
            }
            continue // Move to the next argument since a stacked group was processed
        }

        // Handle normal (non-stacked) flags
        found := false
        for _, def := range defs {
            if arg == "-"+def.Short || arg == "--"+def.Long {
                found = true

                if def.DataType == "bool" {
                    parsedArgs[def.Name] = true
                    break
                }

                // Ensure non-boolean flags have a value following them
                if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
                    rawValue := args[i+1]
                    i++

                    switch def.DataType {
                    case "int":
                        intValue, err := strconv.Atoi(rawValue)
                        if err != nil {
                            return fmt.Errorf("invalid value for argument '%s': expected an integer", def.Name)
                        }
                        parsedArgs[def.Name] = intValue
                    case "string":
                        parsedArgs[def.Name] = rawValue
                    case "[]string":
                        values := []string{rawValue}
                        for i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
                            values = append(values, args[i+1])
                            i++
                        }
                        parsedArgs[def.Name] = values
                    default:
                        return fmt.Errorf("unknown data type '%s' for argument '%s'", def.DataType, def.Name)
                    }
                } else {
                    return fmt.Errorf("no value provided for argument %s", arg)
                }
            }
        }

        if !found {
            return fmt.Errorf("unknown argument: %s", arg)
        }
    }

    // Handle defaults after parsing
    for _, def := range defs {
        if _, ok := parsedArgs[def.Name]; !ok {
            if def.DefaultValue != nil {
                parsedArgs[def.Name] = def.DefaultValue
            } else if def.DataType == "bool" {
                parsedArgs[def.Name] = false
            }
        }
    }

    return nil
}

// Parse the CLI arguments
func (p *Parser) Parse() (map[string]interface{}, bool, error) {
	args := os.Args[1:]

	// Handle "help" request or no arguments passed cases
	if len(args) == 0 || containsHelpArgument(args) {
		p.PrintHelp()
		return nil, true, nil
	}

	if requestedVersion(args) {
		p.PrintVersion()
		return nil, true, nil
	}

	// Parse the individual arguments based on p.args and command structure
	parsedArgs := map[string]interface{}{}

	// Parse global arguments using helper parseArguments func
	err := parseArguments(p.args, args, parsedArgs)
	if err != nil {
		if strings.HasPrefix(err.Error(), "unknown argument") {
			return nil, true, fmt.Errorf("unknown argument: %s", args[0])
		}
		return nil, true, err
	}

	// Validate global required args after parsing all subcommands
	for _, arg := range p.args {
		if arg.Required {
			if _, ok := parsedArgs[arg.Name]; !ok {
				return nil, true, fmt.Errorf("missing required global argument: %s", arg.Name)
			}
		}
	}

	// Validate mutual exclusivity
	err = p.validateExclusiveGroups(parsedArgs)
	if err != nil {
		return nil, true, err
	}

	return parsedArgs, false, nil
}

// PrintVersion does the obvious
func(p *Parser) PrintVersion() {
	if p.Version != "" {
		fmt.Printf("%s Version: %s\n", p.Name, p.Version)
	} else {
		fmt.Println("No version information provided by program.")
	}
}


// PrintHelp does the obvious
func (p *Parser) PrintHelp() {
	// Optional program metadata
	if p.Name != "" {
		fmt.Printf("%s\n", p.Name)
	}
	if p.Author != "" {
		fmt.Printf("Author: %s\n", p.Author)
	}
	if p.Version != "" {
		fmt.Printf("Version: %s\n", p.Version)
	}
	if p.Description != "" {
		fmt.Printf("%s\n", p.Description)
	}



	fmt.Println("Usage:")

	// Sort arguments by name (or long form if available)
	sort.Slice(p.args, func(i, j int) bool {
		return p.args[i].Name < p.args[j].Name
	})

	for _, arg := range p.args {
		fmt.Printf("    -%s, --%s: %s\n", arg.Short, arg.Long, arg.Description)
	}
}

// Helper function to check for help request
func containsHelpArgument(args []string) bool {
	for _, arg := range args {
		if arg == "-h" || arg == "--help" {
			return true
		}
	}
	return false
}

func requestedVersion(args []string) bool {
	for _, arg := range args {
		if arg == "--version" {
			return true
		}
	}
	return false
}