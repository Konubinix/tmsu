package cli

import (
	"fmt"
	"log"
	"tmsu/common"
	"tmsu/storage"
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
	store, err := storage.Open()
	if err != nil {
		log.Fatalf("Something wrong happened: %v", err)
	}
	defer store.Close()

	if len(args) > 2 {
		return fmt.Errorf("Must provide atmost two arguments.")
	}
	if len(args) == 0 {
		return fmt.Errorf("Must provide at least one argument.")
	}
	var name string
	name = args[0]
	if name != "fingerprint_command" {
 		log.Fatalf("Config name not known: %v", name)
	}
	if len(args) == 1 {
		DBConfigExternalProgramGet(store, name)
	}
	if len(args) == 2 {
		DBConfigExternalProgramSet(store, name, args[1])
	}
	return nil
}

func DBConfigExternalProgramSet(store *storage.Storage, name string, value string) error {
	store.Db.DBConfigSetConfig(common.DBConfig{value})
	return nil
}

func DBConfigExternalProgramGet(store *storage.Storage, name string) error {
	config, err := store.Db.DBConfigGetConfig()
	if err != nil {
 		log.Fatalf("Something wrong happened: %v", err)
	}
	fmt.Printf("'%v' = '%v'\n", name, config.FingerPrintCommand)
	return nil
}
