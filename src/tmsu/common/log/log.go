/*
Copyright 2011-2014 Paul Ruane.

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

package log

import (
	"fmt"
	"os"
)

var Verbosity uint = 1

func Fatal(values ...interface{}) {
	Warn(values...)
	os.Exit(1)
}

func Fatalf(format string, values ...interface{}) {
	Warnf(format, values...)
	os.Exit(1)
}

func Warn(values ...interface{}) {
	fmt.Fprint(os.Stderr, "tmsu: ")
	fmt.Fprintln(os.Stderr, values...)
}

func Warnf(format string, values ...interface{}) {
	format = "tmsu: " + format + "\n"
	fmt.Fprintf(os.Stderr, format, values...)
}

func Info(verbosity uint, values ...interface{}) {
	if verbosity > Verbosity {
		return
	}

	fmt.Printf("tmsu: ")
	fmt.Println(values...)
}

func Infof(verbosity uint, format string, values ...interface{}) {
	if verbosity > Verbosity {
		return
	}

	format = "tmsu: " + format + "\n"
	fmt.Printf(format, values...)
}
