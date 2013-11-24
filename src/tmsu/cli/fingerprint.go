package cli

import (
	"fmt"
	"log"
)

var FingerprintCommand = Command{
	Name:     "fingerprint",
	Synopsis: "Modify fingerprint computation facilities",
	Description: `tmsu fingerprint <external-program>

Use external-program to compute the fingerprint of files.`,
	Options: Options{},
	Exec: fingerprintExec,
}

func fingerprintExec(options Options, args []string) error {
	if len(args) > 1 {
		return fmt.Errorf("Must provide atmost one argument")
	}
	if len(args) == 0 {
		fingerprintExternalProgramUnset()
	}
	if len(args) == 1 {
		fingerprintExternalProgramSet(args[0])
	}
	return nil
}

func fingerprintExternalProgramSet(program string) error {
	log.Fatal("Not implemented yet, would add ", program, " as external program")
	return nil
}

func fingerprintExternalProgramUnset() error {
	log.Fatal("Not implemented yet")
	return nil
}
