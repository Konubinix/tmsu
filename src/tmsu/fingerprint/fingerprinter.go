/*
Copyright 2011-2013 Paul Ruane.

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

package fingerprint

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path"
	"strings"
	"regexp"
)

const sparseFingerprintThreshold = 5 * 1024 * 1024
const sparseFingerprintSize = 512 * 1024

func Create(_path string) (Fingerprint, error) {
	// This is a quick hack to test git annex files tagging based on the
	// name of the file the symlink points to
	link, _ := os.Readlink(_path)
	_, file := path.Split(link)
	fg := strings.TrimSuffix(file, path.Ext(file))
	sha_regexp := regexp.MustCompile(`--([a-z0-9]+)`)
	sha := sha_regexp.FindStringSubmatch(fg)[1]

	return Fingerprint(sha), nil
}

func createSparseFingerprint(path string, fileSize int64) (Fingerprint, error) {
	buffer := make([]byte, sparseFingerprintSize)
	hash := sha256.New()

	file, err := os.Open(path)
	if err != nil {
		return EMPTY, err
	}
	defer file.Close()

	// start
	count, err := file.Read(buffer)
	if err != nil {
		return EMPTY, err
	}
	hash.Write(buffer[:count])

	// middle
	_, err = file.Seek((fileSize-sparseFingerprintSize)/2, 0)
	if err != nil {
		return EMPTY, err
	}

	count, err = file.Read(buffer)
	if err != nil {
		return EMPTY, err
	}
	hash.Write(buffer[:count])

	// end
	_, err = file.Seek(-sparseFingerprintSize, 2)
	if err != nil {
		return EMPTY, err
	}

	count, err = file.Read(buffer)
	if err != nil {
		return EMPTY, err
	}
	hash.Write(buffer[:count])

	sum := hash.Sum(make([]byte, 0, 64))
	fingerprint := hex.EncodeToString(sum)

	return Fingerprint(fingerprint), nil
}

func createFullFingerprint(path string) (Fingerprint, error) {
	hash := sha256.New()

	file, err := os.Open(path)
	if err != nil {
		return EMPTY, err
	}
	defer file.Close()

	buffer := make([]byte, 1024)
	for count := 0; err == nil; count, err = file.Read(buffer) {
		hash.Write(buffer[:count])
	}

	sum := hash.Sum(make([]byte, 0, 64))
	fingerprint := hex.EncodeToString(sum)

	return Fingerprint(fingerprint), nil
}

type FileInfoSlice []os.FileInfo

func (infos FileInfoSlice) Len() int {
	return len(infos)
}

func (infos FileInfoSlice) Less(i, j int) bool {
	return infos[i].Name() < infos[j].Name()
}

func (infos FileInfoSlice) Swap(i, j int) {
	infos[j], infos[i] = infos[i], infos[j]
}
