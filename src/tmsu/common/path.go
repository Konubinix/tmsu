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

package common

import (
	"os"
	"path/filepath"
	"strings"
)

func IsDir(path string) bool {
    info, err := os.Stat(path)
    if err != nil {
        switch {
        case os.IsPermission(err):
            Warnf("'%v': Permission denied", path)
        case os.IsNotExist(err):
            Warnf("'%v': No such file", path)
        default:
            Warnf("'%v': Error: %v", err)
        }

        return false
    }

    return info.IsDir()
}

func MakeRelative(path string) string {
	workingDirectory, err := os.Getwd()
	if err != nil {
		return path
	}

	if path == workingDirectory {
		return "."
	}

	if strings.HasPrefix(path, workingDirectory+string(filepath.Separator)) {
		return path[len(workingDirectory)+1:]
	}

	return path
}

func Join(dir, path string) string {
    if filepath.IsAbs(path) {
        return path
    }

    return filepath.Join(dir, path)
}
