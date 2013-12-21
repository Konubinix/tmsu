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
	//"path"
	//"strings"
	//"regexp"
	"fmt"
	// for CreateExternal
	"os/exec"
	"log"
	"bufio"
	"tmsu/common"
	"tmsu/storage"
)

const EMPTY common.Fingerprint = common.Fingerprint("")

const sparseFingerprintThreshold = 5 * 1024 * 1024
const sparseFingerprintSize = 512 * 1024

func Create(_path string) (common.Fingerprint, error) {
	var config common.DBConfig
	var res common.Fingerprint
	store, err := storage.Open()
	if err != nil {
 		log.Fatalf("Something wrong happened: %v", err)
	}
	config, err = store.Db.DBConfigGetConfig()
	if err != nil {
 		log.Fatalf("Something wrong happened: %v", err)
	}
	if config.FingerPrintCommand != "" {
		res, err = CreateExternal(_path, config.FingerPrintCommand)
	} else {
		res, err = CreateInternal(_path)
	}
	return res, err
}

func CreateExternal(path string, command string) (common.Fingerprint, error) {
	cmd := exec.Command(command)
	in, _ := cmd.StdinPipe()
	out, _ := cmd.StdoutPipe()
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(out)
	in.Write([]byte(path + "\n"))
	in.Close()
	res := ""
	for scanner.Scan() {
		res += scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	// interprete the result
	var result common.Fingerprint
	if res == "internal" {
		result, err = CreateInternal(path)
	} else {
		if res == "nil" {
			log.Fatal("Could not compute the fingerprint for %s", path)
		}
		result = common.Fingerprint(res)
	}
	return result, err
}

func CreateInternal(path string) (common.Fingerprint, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return EMPTY, fmt.Errorf("'%v': could not determine if path is a directory: %v", path, err)
	}
	if stat.IsDir() {
		return EMPTY, nil
	}

	fileSize := stat.Size()

	if fileSize > sparseFingerprintThreshold {
		return createSparseFingerprint(path, fileSize)
	}

	return createFullFingerprint(path)
}

func createSparseFingerprint(path string, fileSize int64) (common.Fingerprint, error) {
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

	return common.Fingerprint(fingerprint), nil
}

func createFullFingerprint(path string) (common.Fingerprint, error) {
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

	return common.Fingerprint(fingerprint), nil
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
