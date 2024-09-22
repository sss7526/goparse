package arguments

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
	DataType 		string // e.g., string, []string, int, bool, etc.
	DefaultValue 	interface{}
	Required		bool
}

// Parser is the main type that handles argument parsing
type Parser struct {
	args			[]*Argument
}

// NewParser creates a new instance of the argument parser
func NewParser() *Parser {
	return &Parser{
		args:		[]*Argument{},
	}
}

// AddArgument adds a positional or optional argument to the parser
func (p *Parser) AddArgument(name, short, long, description string, dataType string, required bool) *Argument {
	arg := &Argument {
		Name:			name,
		Short:			short,
		Long:			long,
		Description:	description,
		DataType:		dataType,
		Required:		required,
	}
	p.args = append(p.args, arg)
	return arg
}

func parseArguments(defs []*Argument, args []string, parsedArgs map[string]interface{}) error {
	for _, def := range defs {
		found := false

		for i := 0; i < len(args); i++ {
			arg := args[i]

			// Match short or long argument form
			if arg == "-" + def.Short || arg == "--" + def.Long {

				if def.DataType == "bool" {
					parsedArgs[def.Name] = true
					found = true
					continue
				}

				// Ensure there's a value following non boolean flags
				if i + 1 < len(args) {
					rawValue := args[i + 1]
					i++

					// Perform type-dependent processing
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
						for i + 1 < len(args) && !strings.HasPrefix(args[i + 1], "-") {
							values = append(values, args[i + 1])
							i++
						}
						parsedArgs[def.Name] = values
					default:
						return fmt.Errorf("unknown data type '%s' for argument '%s'", def.DataType, def.Name)
					}
					found = true
				} else {
					return fmt.Errorf("no value provided for argument %s", arg)
				}
			}
		}

		// Check for required arguments
		if def.Required && !found {
			return fmt.Errorf("missing required argument: %s", def.Name)
		}

		// Assign default values for non-found optional arguments
		if !found {
			switch def.DataType {
			case "int":
				parsedArgs[def.Name] = 0
			case "string":
				parsedArgs[def.Name] = ""
			case "[]string":
				parsedArgs[def.Name] = []string{}
			case "bool":
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

	// Parse the individual arguments based on p.args and command structure
	parsedArgs := map[string]interface{}{}

	// Parse global arguments using helper parseArguments func
	err := parseArguments(p.args, args, parsedArgs)
	if err != nil {
		if strings.HasPrefix(err.Error(), "unknown argument") {
			return nil, true, fmt.Errorf("unknown argument: %s", args[0])
		}
		return nil, false, err
	}

	// Validate global required args after parsing all subcommands
	for _, arg := range p.args {
		if arg.Required {
			if _, ok := parsedArgs[arg.Name]; !ok {
				return nil, true, fmt.Errorf("missing required global argument: %s", arg.Name)
			}
		}
	}

	return parsedArgs, false, nil
}

// PrintHelp does the obvious
func (p *Parser) PrintHelp() {
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