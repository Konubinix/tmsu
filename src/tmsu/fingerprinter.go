/*
Copyright 2011 Paul Ruane.

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
package main
*/

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
)

func Fingerprint(path string) (string, error) {
	file, error := os.Open(path)
	if error != nil { return "", error }
	defer file.Close()

	hash := sha256.New()

	buffer := make([]byte, 1024)
	for count := 0; error == nil; count, error = file.Read(buffer) {
		hash.Write(buffer[:count])
	}

	sum := hash.Sum()
	fingerprint := hex.EncodeToString(sum)

	return fingerprint, nil
}
