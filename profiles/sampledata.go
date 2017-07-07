// Copyright Safing ICS Technologies GmbH. Use of this source code is governed by the AGPL license that can be found in the LICENSE file.

package profiles

import (
	ds "github.com/ipfs/go-datastore"

	"safing/database"
)

func init() {

	// Data here is for demo purposes, Profiles will be served over network soon™.

	(&Profile{
		Name:         "Chromium",
		Description:  "Browser by Google",
		Path:         "/usr/lib/chromium-browser/chromium-browser",
		Flags:        []int8{User, Internet, LocalNet, Browser},
		ConnectPorts: []uint16{80, 443},
	}).CreateInDist()

	(&Profile{
		Name:          "Evolution",
		Description:   "PIM solution by GNOME",
		Path:          "/usr/bin/evolution",
		Flags:         []int8{User, Internet, Gateway},
		ConnectPorts:  []uint16{25, 80, 143, 443, 465, 587, 993, 995},
		SecurityLevel: 2,
	}).CreateInDist()

	(&Profile{
		Name:          "Evolution Calendar",
		Description:   "PIM solution by GNOME - Calendar",
		Path:          "/usr/lib/evolution/evolution-calendar-factory-subprocess",
		Flags:         []int8{User, Internet, Gateway},
		ConnectPorts:  []uint16{80, 443},
		SecurityLevel: 2,
	}).CreateInDist()

	(&Profile{
		Name:         "Spotify",
		Description:  "Music streaming",
		Path:         "/usr/share/spotify/spotify",
		ConnectPorts: []uint16{80, 443, 4070},
		Flags:        []int8{User, Internet, Strict},
	}).CreateInDist()

	(&Profile{
		// flatpak edition
		Name:         "Spotify",
		Description:  "Music streaming",
		Path:         "/newroot/app/extra/share/spotify/spotify",
		ConnectPorts: []uint16{80, 443, 4070},
		Flags:        []int8{User, Internet, Strict},
	}).CreateInDist()

	(&Profile{
		Name:          "Evince",
		Description:   "PDF Document Reader",
		Path:          "/usr/bin/evince",
		Flags:         []int8{},
		SecurityLevel: 2,
	}).CreateInDist()

	(&Profile{
		Name:        "Ahavi",
		Description: "mDNS service",
		Path:        "/usr/bin/avahi-daemon",
		Flags:       []int8{System, LocalNet, Service, Directconnect},
	}).CreateInDist()

	(&Profile{
		Name:        "Python 2.7 Framework",
		Description: "Correctly handle python scripts",
		Path:        "/usr/bin/python2.7",
		Framework: &Framework{
			Find:  "^[^ ]+ ([^ ]+)",
			Build: "{1}|{CWD}/{1}",
		},
	}).CreateInDist()

	(&Profile{
		Name:        "Python 3.5 Framework",
		Description: "Correctly handle python scripts",
		Path:        "/usr/bin/python3.5",
		Framework: &Framework{
			Find:  "^[^ ]+ ([^ ]+)",
			Build: "{1}|{CWD}/{1}",
		},
	}).CreateInDist()

	(&Profile{
		Name:        "DHCP Client",
		Description: "Client software for the DHCP protocol",
		Path:        "/sbin/dhclient",
		Framework: &Framework{
			FindParent:      1,
			MergeWithParent: true,
		},
	}).CreateInDist()

	// Default Profiles
	// Until Profiles are distributed over the network, default profiles are activated when the Default Profile for "/" is missing.

	if ok, err := database.Has(ds.NewKey("/Data/Profiles/Profile:d-2f")); !ok || err != nil {

		(&Profile{
			Name:        "Default Base",
			Description: "Default Profile for /",
			Path:        "/",
			Flags:       []int8{Internet, LocalNet, Strict},
			Default:     true,
		}).Create()

		(&Profile{
			Name:        "Installed Applications",
			Description: "Default Profile for /usr/bin",
			Path:        "/usr/bin/",
			Flags:       []int8{Internet, LocalNet, Gateway},
			Default:     true,
		}).Create()

		(&Profile{
			Name:        "System Binaries (/sbin)",
			Description: "Default Profile for ~/Downloads",
			Path:        "/sbin/",
			Flags:       []int8{Internet, LocalNet, Directconnect, Service, System},
			Default:     true,
		}).Create()

		(&Profile{
			Name:        "System Binaries (/usr/sbin)",
			Description: "Default Profile for ~/Downloads",
			Path:        "/usr/sbin/",
			Flags:       []int8{Internet, LocalNet, Directconnect, Service, System},
			Default:     true,
		}).Create()

		(&Profile{
			Name:        "System Tmp folder",
			Description: "Default Profile for /tmp",
			Path:        "/tmp/",
			Flags:       []int8{}, // deny all
			Default:     true,
		}).Create()

		(&Profile{
			Name:        "User Home",
			Description: "Default Profile for ~/",
			Path:        "~/",
			Flags:       []int8{Internet, LocalNet, Gateway},
			Default:     true,
		}).Create()

		(&Profile{
			Name:        "User Downloads",
			Description: "Default Profile for ~/Downloads",
			Path:        "~/Downloads/",
			Flags:       []int8{}, // deny all
			Default:     true,
		}).Create()

		(&Profile{
			Name:        "User Cache",
			Description: "Default Profile for ~/.cache",
			Path:        "~/.cache/",
			Flags:       []int8{}, // deny all
			Default:     true,
		}).Create()

	}

}
