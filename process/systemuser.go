// Copyright Safing ICS Technologies GmbH. Use of this source code is governed by the AGPL license that can be found in the LICENSE file.

package process

import (
	"bufio"
	"fmt"
	"os"
	"github.com/safing/safing-core/log"
	"strconv"
	"strings"
)

var userNames map[int]string
var userHomes map[int]string

func (process *Process) getUserName() {
	name, ok := userNames[process.UserID]
	if !ok {
		// TODO: rescan for new users (and make this goroutine safe)
		process.UserName = fmt.Sprintf("%d", process.UserID)
	} else {
		process.UserName = name
	}
}

func (process *Process) getUserHome() {
	home, ok := userHomes[process.UserID]
	if ok {
		process.UserHome = home
	}
}

func init() {

	userNames = make(map[int]string)
	userHomes = make(map[int]string)

	file, err := os.Open("/etc/passwd")
	if err != nil {
		log.Errorf("process: failed to open /etc/passwd, will use IDs instead of names: %s", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ":")
		if len(parts) < 6 {
			log.Warningf("process: invalid line in /etc/passwd")
		}
		id, err := strconv.Atoi(parts[2])
		if err != nil {
			log.Warningf("process: failed to parse user id %s in /etc/passwd", parts[2])
		}
		userNames[int(id)] = parts[0]
		if parts[5] != "" {
			userHomes[int(id)] = parts[5]
		}
	}

	if err := scanner.Err(); err != nil {
		log.Errorf("process: error while reading /etc/passwd, some usernames might be missing, will use IDs in that case: %s", err)
	}
}
