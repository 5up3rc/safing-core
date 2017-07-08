// Copyright Safing ICS Technologies GmbH. Use of this source code is governed by the AGPL license that can be found in the LICENSE file.

package process

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"syscall"

	datastore "github.com/ipfs/go-datastore"

	"github.com/Safing/safing-core/database"
	"github.com/Safing/safing-core/log"
	"github.com/Safing/safing-core/profiles"
)

// A Process represents a process running on the operating system
type Process struct {
	database.Base
	UserID     int
	UserName   string
	UserHome   string
	Pid        int
	ParentPid  int
	Path       string
	Cwd        string
	FileInfo   *FileInfo
	CmdLine    string
	FirstArg   string
	ProfileKey string
	Profile    *profiles.Profile
	Name       string
	Icon       string
	// Icon is a path to the icon and is either prefixed "f:" for filepath, "d:" for database cache path or "c:"/"a:" for a the icon key to fetch it from a company / authoritative node and cache it in its own cache.
}

var processModel *Process // only use this as parameter for database.EnsureModel-like functions

func init() {
	database.RegisterModel(processModel, func() database.Model { return new(Process) })
}

// Create saves Process with the provided name in the default namespace.
func (m *Process) Create(name string) error {
	return m.CreateObject(&database.Processes, name, m)
}

// CreateInNamespace saves Process with the provided name in the provided namespace.
func (m *Process) CreateInNamespace(namespace *datastore.Key, name string) error {
	return m.CreateObject(namespace, name, m)
}

// Save saves Process.
func (m *Process) Save() error {
	return m.SaveObject(m)
}

// GetProcess fetches Process with the provided name from the default namespace.
func GetProcess(name string) (*Process, error) {
	return GetProcessFromNamespace(&database.Processes, name)
}

// GetProcessFromNamespace fetches Process with the provided name from the provided namespace.
func GetProcessFromNamespace(namespace *datastore.Key, name string) (*Process, error) {
	object, err := database.GetAndEnsureModel(namespace, name, processModel)
	if err != nil {
		return nil, err
	}
	model, ok := object.(*Process)
	if !ok {
		return nil, database.NewMismatchError(object, processModel)
	}
	return model, nil
}

func (m *Process) String() string {
	if m == nil {
		return "?"
	}
	if m.Profile != nil && !m.Profile.Default {
		return fmt.Sprintf("%s:%s:%d", m.UserName, m.Profile, m.Pid)
	}
	return fmt.Sprintf("%s:%s:%d", m.UserName, m.Path, m.Pid)
}

func GetOrFindProcess(pid int) (*Process, error) {
	process, err := GetProcess(strconv.Itoa(pid))
	if err == nil {
		return process, nil
	}

	new := &Process{
		Pid: pid,
	}

	new.Path, err = os.Readlink(fmt.Sprintf("/proc/%d/exe", pid))
	if err != nil {
		return nil, fmt.Errorf("could not read /proc/%d/exe: %s", pid, err)
	}

	new.Cwd, err = os.Readlink(fmt.Sprintf("/proc/%d/cwd", pid))
	if err != nil {
		return nil, fmt.Errorf("could not read /proc/%d/cwd: %s", pid, err)
	}

	cmdLine, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid))
	if err != nil {
		return nil, fmt.Errorf("could not read /proc/%d/cmdline: %s", pid, err)
	}
	// convert null bytes to spaces before converting to string
	new.CmdLine = strings.Trim(string(bytes.Replace(cmdLine, []byte{0x00}, []byte{0x20}, -1)), " ")
	// log.Tracef("loaded cmdline: %v", new.CmdLine)

	procStatData, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
	if err != nil {
		return nil, fmt.Errorf("could not read /proc/%d/stat: %s", pid, err)
	}
	fields := strings.SplitN(string(procStatData), " ", 5)
	if len(fields) < 5 {
		return nil, fmt.Errorf("could not parse %d", pid)
	}
	parentPid, err := strconv.ParseInt(fields[3], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("could not parse parent pid: %s", fields[3])
	}
	new.ParentPid = int(parentPid)

	statData, err := os.Stat(fmt.Sprintf("/proc/%d", pid))
	if err != nil {
		return nil, fmt.Errorf("could not stat /proc/%d: %s", pid, err)
	}
	sys, ok := statData.Sys().(*syscall.Stat_t)
	if !ok {
		return nil, fmt.Errorf("unable to parse /proc/%d: wrong type", pid)
	}
	new.UserID = int(sys.Uid)

	// get username and home
	new.getUserName()
	new.getUserHome()

	// try to get Process name and icon
	new.GetHumanInfo()

	// get Profile
	processPath := new.Path
	var applyProfile *profiles.Profile
	iterations := 0
	for applyProfile == nil {

		iterations++
		if iterations > 10 {
			log.Warningf("process: got into loop while getting profile for %s", new)
			break
		}

		applyProfile, err = profiles.GetActiveProfileByPath(processPath)
		if err == database.ErrNotFound {
			applyProfile, err = profiles.FindProfileByPath(processPath, new.UserHome)
		}
		if err != nil {
			log.Warningf("process: could not get profile for %s: %s", new, err)
		} else if applyProfile == nil {
			log.Warningf("process: no default profile found for %s", new)
		} else {

			// TODO: there is a lot of undefined behaviour if chaining framework profiles

			// process framework
			if applyProfile.Framework != nil {
				if applyProfile.Framework.FindParent > 0 {
					ppid := new.ParentPid
					for i := uint8(1); i < applyProfile.Framework.FindParent; i++ {
						ppid, err = GetParentPid(ppid)
						if err != nil {
							return nil, err
						}
					}
					if applyProfile.Framework.MergeWithParent {
						return GetOrFindProcess(ppid)
					}
					processPath, err = os.Readlink(fmt.Sprintf("/proc/%d/exe", pid))
					if err != nil {
						return nil, fmt.Errorf("could not read /proc/%d/exe: %s", pid, err)
					}
					continue
				}

				newCommand, err := applyProfile.Framework.GetNewPath(new.CmdLine, new.Cwd)
				if err != nil {
					return nil, err
				}

				// assign
				new.CmdLine = newCommand
				new.Path = strings.SplitN(newCommand, " ", 2)[0]
				processPath = new.Path

				// make sure we loop
				applyProfile = nil
				continue
			}

			// apply profile to process
			log.Debugf("process: applied profile to %s: %s", new, applyProfile)
			new.Profile = applyProfile
			new.ProfileKey = applyProfile.GetKey().String()

			// update Profile with Process icon if Profile does not have one
			if !new.Profile.Default && new.Icon != "" && new.Profile.Icon == "" {
				new.Profile.Icon = new.Icon
				new.Profile.Save()
			}
		}
	}

	// get FileInfo
	new.FileInfo = GetFileInfo(new.Path)

	// save to DB
	new.Create(strconv.Itoa(new.Pid))

	return new, nil
}

func GetParentPid(pid int) (int, error) {
	procStatData, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
	if err != nil {
		return 0, fmt.Errorf("could not read /proc/%d/stat: %s", pid, err)
	}
	fields := strings.SplitN(string(procStatData), " ", 5)
	if len(fields) < 5 {
		return 0, fmt.Errorf("could not parse %d", pid)
	}
	parentPid, err := strconv.ParseInt(fields[3], 10, 32)
	if err != nil {
		return 0, fmt.Errorf("could not parse parent pid: %s", fields[3])
	}
	return int(parentPid), nil
}
