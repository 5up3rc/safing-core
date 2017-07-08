// Copyright Safing ICS Technologies GmbH. Use of this source code is governed by the AGPL license that can be found in the LICENSE file.

package profiles

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Safing/safing-core/log"
)

type Framework struct {
	// go hirarchy up
	FindParent      uint8 `json:",omitempty bson:",omitempty"`
	MergeWithParent bool  `json:",omitempty bson:",omitempty"`

	// go hirarchy down
	Virtual bool   `json:",omitempty bson:",omitempty"`
	Find    string `json:",omitempty bson:",omitempty"`
	Build   string `json:",omitempty bson:",omitempty"`
}

func (f *Framework) GetNewPath(command string, cwd string) (string, error) {
	// "/usr/bin/python script"
	// to
	// "/path/to/script"
	regex, err := regexp.Compile(f.Find)
	if err != nil {
		return "", fmt.Errorf("profiles(framework): failed to compile framework regex: %s", err)
	}
	matched := regex.FindAllStringSubmatch(command, -1)
	if len(matched) == 0 || len(matched[0]) < 2 {
		return "", fmt.Errorf("profiles(framework): regex \"%s\" for constructing path did not match command \"%s\"", f.Find, command)
	}

	var lastError error
	var buildPath string
	for _, buildPath = range strings.Split(f.Build, "|") {

		buildPath = strings.Replace(buildPath, "{CWD}", cwd, -1)
		for i := 1; i < len(matched[0]); i++ {
			buildPath = strings.Replace(buildPath, fmt.Sprintf("{%d}", i), matched[0][i], -1)
		}

		buildPath = filepath.Clean(buildPath)

		if !f.Virtual {
			if !strings.HasPrefix(buildPath, "~/") && !filepath.IsAbs(buildPath) {
				lastError = fmt.Errorf("constructed path \"%s\" from framework is not absolute", buildPath)
				continue
			}
			if _, err := os.Stat(buildPath); os.IsNotExist(err) {
				lastError = fmt.Errorf("constructed path \"%s\" does not exist", buildPath)
				continue
			}
		}

		lastError = nil
		break

	}

	if lastError != nil {
		return "", fmt.Errorf("profiles(framework): failed to construct valid path, last error: %s", lastError)
	}
	log.Tracef("profiles(framework): transformed \"%s\" (%s) to \"%s\"", command, cwd, buildPath)
	return buildPath, nil
}
