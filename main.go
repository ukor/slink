package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
)

type FlagDescription struct {
	flagType    string
	description string
}

type Flag struct {
	hasValue bool
	value    FlagDescription
	command  string
}

type FlagValue struct {
	flag  string
	value string
}

var supportedFlags = map[string]Flag{
	"source": {
		hasValue: true,
		value:    FlagDescription{flagType: "directory", description: ""},
		command:  "link_to",
	},
	"target": {
		hasValue: true,
		value:    FlagDescription{flagType: "directory", description: ""},
		command:  "link_from",
	},
	"destination": {
		hasValue: true,
		value:    FlagDescription{flagType: "directory", description: ""},
	},
	"dry-run": {
		hasValue: false,
		value:    FlagDescription{flagType: "boolean", description: ""},
	},
	"dotfiles": {
		hasValue: false,
		value:    FlagDescription{flagType: "boolean", description: ""},
	},
	"help": {
		hasValue: false,
		value:    FlagDescription{flagType: "boolean", description: ""},
	},
}

var subcommands = [2]string{"list", "ls"}

func inArray(term string, list ...string) bool {
	// ---
	for _, b := range list {
		if term == b {
			return true
		}
	}

	return false
}

func cleanFlag(flag string) string {
	if after, ok := strings.CutPrefix(flag, "--"); ok {
		return after
	}

	panic("Not a flag")
}

func parseEqualTo(list []string) []string {
	withoutEqualTo := []string{}

	for _, v := range list {

		trimmed := strings.TrimSpace(v)

		if trimmed == "=" {
			continue
		} else if strings.Contains(trimmed, "=") {
			c := strings.Split(trimmed, "=")
			word := strings.Join(c, " ")
			word = strings.TrimSpace(word)
			withoutEqualTo = append(withoutEqualTo, strings.Split(word, " ")...)
		} else {
			withoutEqualTo = append(withoutEqualTo, trimmed)
		}

	}

	return withoutEqualTo
}

func parseFlagsAndArgument(args []string) map[string]any {
	flagToValues := map[string]any{
		"dry-run":  false,
		"dotfiles": false,
	}

	argsSize := len(args)

	for i, arg := range args {
		if flagName, ok := strings.CutPrefix(arg, "--"); ok {

			f := supportedFlags[flagName]

			if f.hasValue {

				if i+1 >= argsSize {
					log.Fatalf("Error: flag --%s is expecting a value but got none", flagName)
				}

				flagValue := args[i+1]
				fmt.Println(flagName, "is expecting a value. Next =>", flagValue)

				if strings.HasPrefix(flagValue, "--") {
					log.Fatalf("Error : %v is expecting a value but got none", flagName)
				}

				if p, ok := strings.CutPrefix(flagValue, "~/"); ok {
					home, err := os.UserHomeDir()
					if err != nil {
						log.Fatalf("Error : %v ", err)
					}

					flagValue = filepath.Join(home, p)
				}

				err := checkArgumentType(flagValue, f.value.flagType)
				if err != nil {
					log.Fatalf("Error : %v", err)
				}

				flagToValues[flagName] = flagValue

				i++
			} else {
				flagToValues[flagName] = true
			}
			fmt.Println("flag ->", f, "flag values ->", flagToValues, "after = ", flagName)

		}
	}

	return flagToValues
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		printUsage()
		os.Exit(1)
	}

	subcommand := args[0]
	fmt.Println("args before parsing =>", args)
	args = parseEqualTo(args)
	fmt.Println("args after parsing => ", args)

	if slices.Contains(args, "--target") && slices.Contains(args, "--source") {
		log.Fatalf("Command can not contain --target and --source \n")
	}

	if inArray("=", subcommands[:]...) {
		message := fmt.Sprintf("%s is not supported ", subcommand)
		fmt.Println(message)
		os.Exit(2)
	}

	flagValuesMap := parseFlagsAndArgument(args)
	fmt.Println()
	fmt.Println("flagValuesMap === ", flagValuesMap)
	fmt.Println()

	source, isSourceOk := flagValuesMap["source"]
	destination, isDestinationOk := flagValuesMap["destination"]
	if isSourceOk && !isDestinationOk {
		log.Fatalf("Error: Destination is required \n")
	}

	fmt.Printf("source => %v destination => %v \n", source, destination)
}

// printUsage displays how to use the CLI application.
func printUsage() {
	fmt.Println("Simple Directory Lister CLI")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  linker [path] [--all]")
	fmt.Println("")
	// fmt.Println("Commands:")
	// fmt.Println("  list   Lists directories in the specified path.")
	// fmt.Println("         If no path is provided, it defaults to the current directory.")
	fmt.Println("")
	fmt.Println("Flags:")
	fmt.Println("  --source  .")
	fmt.Println("  --target  .")
	fmt.Println("  --destination .")
	fmt.Println("  --dotfiles   .")
	fmt.Println("  --dry-run .")
	fmt.Println("  --help .")
}

func checkArgumentType(argument interface{}, argType string) error {
	switch argType {

	case "directory":
		argString, ok := argument.(string)
		if !ok {
			return fmt.Errorf("expected a string for directory path, got %T", argument)
		}

		info, err := os.Stat(argString)
		if err != nil {

			if os.IsNotExist(err) {
				return fmt.Errorf("%s :Path does not exist", argument)
			}
			return fmt.Errorf("could not access path %s: %w", argument, err)

		}

		if info.IsDir() {
			return nil
		}
		return fmt.Errorf("path is not a directory: %s", argument)

	case "boolean":
		if reflect.ValueOf(argument).Kind() != reflect.Bool {
			return fmt.Errorf("expected a boolean for: %v", argument)
		}

		return nil
	default:

		return fmt.Errorf("unspecified type received for: %s", argument)
	}
}
