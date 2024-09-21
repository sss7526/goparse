package arguments

import (
	"errors"
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


// Command represents a subcommand including its own arguments
type Command struct {
	Name			string
	Description		string
	Arguments		[]*Argument
	Subcommands		[]*Command
}

type ExclusiveGroup struct {
	Arguments 		[]*Argument
}

// Parser is the main type that handles argument parsing
type Parser struct {
	commands		[]*Command
	args			[]*Argument
}

// NewParser creates a new instance of the argument parser
func NewParser() *Parser {
	return &Parser{
		args:		[]*Argument{},
		commands:	[]*Command{},
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

// AddCommand adds a subcommand to the parser
func (p *Parser) AddCommand(name, description string) *Command {
	cmd := &Command {
		Name:			name,
		Description:	description,
		Arguments:		[]*Argument{},
	}
	p.commands = append(p.commands, cmd)
	return cmd
}

func (p *Parser) ValidateExclusiveGroups(groups ...*ExclusiveGroup) error {
	for _, group := range groups {
		selected := []string{}
		for _, arg := range group.Arguments {
			if value, ok := os.LookupEnv(arg.Name); ok && value != "" {
				selected = append(selected, arg.Name)
			}
		}
		if len(selected) > 1 {
			return fmt.Errorf("exclusive argument error: cannot use arguments %v together", selected)
		}
	}
	return nil
}

// Parse the CLI arguments
func (p *Parser) Parse() (map[string]interface{}, error) {
	args := os.Args[1:]

	// Handle "help" request or no arguments passed cases
	if len(args) == 0 || containsHelpArgument(args) {
		p.PrintHelp()
		return nil, errors.New("help requested or no arguments were provided")
	}

	// Parse the individual arguments based on p.args and command structure
	parsedArgs := map[string]interface{}{}

	// Parse subcommands first
	if len(args) > 0 {
		for _, cmd := range p.commands {
			if args[0] == cmd.Name {
				args = args[1:] // Remove the subcommand from args

				// Argument parsing within subcommand
				for _, def := range cmd.Arguments {
					for i := 0; i < len(args); i++ {
						arg := args[i]

						if arg == "-"+def.Short || arg == "--"+ def.Long {
							if def.DataType == "bool" {
								parsedArgs[def.Name] = true
							} else {
								if i + 1 < len(args) {
									parsedArgs[def.Name] = args[i + 1]
									i++
								} else {
									return nil, fmt.Errorf("no value provided for argument %s", arg)
								}
							}
						}
					}
				}
				return parsedArgs, nil
			}
		}
	}


	// Populate parsedArgs by parsing args based on the defined arguments
	// Iterate through the argument definitions
	for _, def := range p.args {
		found := false

		for i := 0; i < len(args); i++ {
			arg := args[i]

			// Match short or long argument form
			if arg == "-"+def.Short || arg == "--"+def.Long {

				if def.DataType == "bool" {
					parsedArgs[def.Name] = true
					found = true
					continue
				}
				// Ensure theres a value following the flag
				if i + 1 < len(args) {
					rawValue := args[i + 1]
					i++
					
					// Perform type-dependent processing
					switch def.DataType {
					case "int":
						intValue, err := strconv.Atoi(rawValue)
						if err != nil {
							return nil, fmt.Errorf("invalid value for argument '%s': expected an integer", def.Name)
						}
						parsedArgs[def.Name] = intValue
					case "string":
						parsedArgs[def.Name] = rawValue
					case "[]string":
						values := []string{rawValue}
						for i + 1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
							values = append(values, args[i+1])
							i++
						}
						parsedArgs[def.Name] = values
					default:
						return nil, fmt.Errorf("unknown data type '%s' for argument '%s'", def.DataType, def.Name)
					}
					found = true
				} else {
					return nil, fmt.Errorf("no value provided for argument %s", arg)
				}
			}
		}

		if !found {
			switch def.DataType {
			case "int":
				parsedArgs[def.Name] = 0
			case "string":
				parsedArgs[def.Name] = ""
			case "[]string":
				parsedArgs[def.Name] = nil
			case "bool":
				parsedArgs[def.Name] = false
			}
		}
	}

	return parsedArgs, nil
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

	// For subcommands
	if len(p.commands) > 0 {
		fmt.Println("\nCommands:")
		for _, command := range p.commands {
			fmt.Printf("    %s: %s\n", command.Name, command.Description)
		}
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