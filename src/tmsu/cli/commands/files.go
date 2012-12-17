/*
Copyright 2011-2012 Paul Ruane.

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package commands

import (
	"errors"
	"fmt"
	"sort"
	"tmsu/cli"
	"tmsu/common"
	"tmsu/log"
	"tmsu/storage"
)

type FilesCommand struct{}

func (FilesCommand) Name() cli.CommandName {
	return "files"
}

func (FilesCommand) Synopsis() string {
	return "List files with particular tags"
}

func (FilesCommand) Description() string {
	return `tmsu files [--explicit] [-]TAG...
tmsu files --all

Lists the files, if any, that have all of the TAGs specified. Tags can be excluded by prefixing them with a minus (-).`
}

func (FilesCommand) Options() cli.Options {
	return cli.Options{{"-a", "--all", "show the complete set of tagged files"},
		{"-e", "--explicit", "show only the explicitly tagged files"}}
}

func (command FilesCommand) Exec(options cli.Options, args []string) error {
	if cli.HasOption(options, "--all") {
		return command.listAllFiles()
	}

	explicitOnly := cli.HasOption(options, "--explicit")

	return command.listFiles(args, explicitOnly)
}

func (FilesCommand) listAllFiles() error {
	store, err := storage.Open()
	if err != nil {
		return err
	}
	defer store.Close()

	files, err := store.Files()
	if err != nil {
		return err
	}

	for _, file := range files {
		fmt.Println(file.Path())
	}

	return nil
}

func (FilesCommand) listFiles(args []string, explicitOnly bool) error {
	if len(args) == 0 {
		return errors.New("At least one tag must be specified. Use --all to show all files.")
	}

	store, err := storage.Open()
	if err != nil {
		return err
	}
	defer store.Close()

	includeTagIds := make([]uint, 0)
	excludeTagIds := make([]uint, 0)
	for _, arg := range args {
		var tagName string
		var include bool
		if arg[0] == '-' {
			tagName = arg[1:]
			include = false
		} else {
			tagName = arg
			include = true
		}

		tag, err := store.TagByName(tagName)
		if err != nil {
			return err
		}
		if tag == nil {
			log.Fatalf("No such tag '%v'.", tagName)
		}

		if include {
			includeTagIds = append(includeTagIds, tag.Id)
		} else {
			excludeTagIds = append(excludeTagIds, tag.Id)
		}
	}

	files, err := store.FilesWithTags(includeTagIds, excludeTagIds, explicitOnly)
	if err != nil {
		return err
	}

	paths := make([]string, len(files))
	for _, file := range files {
		relPath := common.RelPath(file.Path())
		paths = append(paths, relPath)
	}

	sort.Strings(paths)

	previousPath := ""
	for _, path := range paths {
		if path != previousPath {
			fmt.Println(path)
		}

		previousPath = path
	}

	return nil
}
