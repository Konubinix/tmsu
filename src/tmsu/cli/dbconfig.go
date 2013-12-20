package cli

import (
	"fmt"
	"log"
)

var DBConfigCommand = Command{
	Name:     "dbconfig",
	Synopsis: "Modify database configuration",
	Description: `tmsu dbconfig <name> [<value>]

Put the value <value> in the configuration of name <name>. If the value is not specified, print the current value.`,
	Options: Options{},
	Exec: DBConfigExec,
}

func DBConfigExec(options Options, args []string) error {
	if len(args) > 2 {
		return fmt.Errorf("Must provide atmost two arguments.")
	}
	if len(args) == 0 {
		return fmt.Errorf("Must provide at least one argument.")
	}
	if len(args) == 1 {
		DBConfigExternalProgramGet(args[0])
	}
	if len(args) == 2 {
		DBConfigExternalProgramSet(args[0], args[1])
	}
	return nil
}

func DBConfigExternalProgramSet(name string, value string) error {
	log.Fatal(name, " = ", value, ", not implemented yet")
	return nil
}

func DBConfigExternalProgramGet(name string) error {
	log.Fatal("Get ", name, ", not implemented yet")
	return nil
}
