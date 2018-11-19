/*
 * Copyright 1999-2018 Alibaba Group Holding Ltd.
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *      http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package nacos

import (
	"strings"
	"os"
	"bytes"
	"path/filepath"
	"os/user"
	"runtime"
	"os/exec"
	"errors"
)

var DNSDomains = make(map[string]string)
var DNSTTL uint32 = 1

func Exist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {

	}
	return dir
}

func Home() string {
	user2, err := user.Current()
	if nil == err {
		return user2.HomeDir
	}

	if "windows" == runtime.GOOS {
		return homeWindows()
	}

	// Unix-like system, so just assume Unix
	return homeUnix()
}

func homeUnix() string {
	// First prefer the HOME environmental variable
	if home := os.Getenv("HOME"); home != "" {
		return home
	}

	// If that fails, try the shell
	var stdout bytes.Buffer
	cmd := exec.Command("sh", "-c", "eval echo ~$USER")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	result := strings.TrimSpace(stdout.String())
	if result == "" {
		panic(errors.New("blank output when reading home directory"))
	}

	return result
}

func homeWindows() string{
	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	home := drive + path
	if drive == "" || path == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home == "" {
		panic(errors.New("HOMEDRIVE, HOMEPATH, and USERPROFILE are blank"))
	}

	return home
}
