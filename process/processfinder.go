// Copyright Safing ICS Technologies GmbH. Use of this source code is governed by the AGPL license that can be found in the LICENSE file.

package process

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"sync"
	"syscall"

	"safing/log"
)

var (
	pidsByUserLock sync.Mutex
	pidsByUser     = make(map[int][]int)
)

func GetPidOfInode(uid, inode int) (int, bool) {
	pidsByUserLock.Lock()
	defer pidsByUserLock.Unlock()

	pidsUpdated := false

	// get pids of user, update if missing
	pids, ok := pidsByUser[uid]
	if !ok {
		// log.Trace("process: no processes of user, updating table")
		updatePids()
		pidsUpdated = true
	}
	pids, ok = pidsByUser[uid]
	if ok {
		// if user has pids, start checking them first
		var checkedUserPids []int
		for _, possiblePID := range pids {
			if findSocketFromPid(possiblePID, inode) {
				return possiblePID, true
			}
			checkedUserPids = append(checkedUserPids, possiblePID)
		}
		// if we fail on the first run and have not updated, update and check the ones we haven't tried so far.
		if !pidsUpdated {
			// log.Trace("process: socket not found in any process of user, updating table")
			// update
			updatePids()
			// sort for faster search
			for i, j := 0, len(checkedUserPids)-1; i < j; i, j = i+1, j-1 {
				checkedUserPids[i], checkedUserPids[j] = checkedUserPids[j], checkedUserPids[i]
			}
			len := len(checkedUserPids)
			// check unchecked pids
			for _, possiblePID := range pids {
				// only check if not already checked
				if sort.SearchInts(checkedUserPids, possiblePID) == len {
					if findSocketFromPid(possiblePID, inode) {
						return possiblePID, true
					}
				}
			}
		}
	}

	// check all other pids
	// log.Trace("process: socket not found in any process of user, checking all pids")
	// TODO: find best order for pidsByUser for best performance
	for possibleUID, pids := range pidsByUser {
		if possibleUID != uid {
			for _, possiblePID := range pids {
				if findSocketFromPid(possiblePID, inode) {
					return possiblePID, true
				}
			}
		}
	}

	return -1, false
}

func findSocketFromPid(pid, inode int) bool {
	socketName := fmt.Sprintf("socket:[%d]", inode)
	entries := readDirNames(fmt.Sprintf("/proc/%d/fd", pid))
	if len(entries) == 0 {
		return false
	}

	for _, entry := range entries {
		link, err := os.Readlink(fmt.Sprintf("/proc/%d/fd/%s", pid, entry))
		if err != nil {
			if !os.IsNotExist(err) {
				log.Warningf("process: failed to read link /proc/%d/fd/%s: %s", pid, entry, err)
			}
			continue
		}
		if link == socketName {
			return true
		}
	}

	return false
}

func updatePids() {
	pidsByUser = make(map[int][]int)

	entries := readDirNames("/proc")
	if len(entries) == 0 {
		return
	}

entryLoop:
	for _, entry := range entries {
		pid, err := strconv.ParseInt(entry, 10, 32)
		if err != nil {
			continue entryLoop
		}

		statData, err := os.Stat(fmt.Sprintf("/proc/%d", pid))
		if err != nil {
			log.Warningf("process: could not stat /proc/%d: %s", pid, err)
			continue entryLoop
		}
		sys, ok := statData.Sys().(*syscall.Stat_t)
		if !ok {
			log.Warningf("process: unable to parse /proc/%d: wrong type", pid)
			continue entryLoop
		}

		pids, ok := pidsByUser[int(sys.Uid)]
		if ok {
			pidsByUser[int(sys.Uid)] = append(pids, int(pid))
		} else {
			pidsByUser[int(sys.Uid)] = []int{int(pid)}
		}

	}

	for _, slice := range pidsByUser {
		for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
			slice[i], slice[j] = slice[j], slice[i]
		}
	}

}

func readDirNames(dir string) (names []string) {
	file, err := os.Open(dir)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Warningf("process: could not open directory %s: %s", dir, err)
		}
		return
	}
	defer file.Close()
	names, err = file.Readdirnames(0)
	if err != nil {
		log.Warningf("process: could not get entries from direcotry %s: %s", dir, err)
		return []string{}
	}
	return
}
