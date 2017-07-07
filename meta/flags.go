// Copyright Safing ICS Technologies GmbH. Use of this source code is governed by the AGPL license that can be found in the LICENSE file.

package meta

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var (
	databaseDir = flag.String("db", "data", "specify custom location for the database")
	logLevel    = flag.String("log", "trace", "set log level [trace|debug|info|warning|error|critical]")
	fileLevels  = flag.String("flog", "", "set log level of files: database=trace,firewall=debug")
	showVersion = flag.Bool("v", false, "show version and exit")
)

func init() {
	flag.Parse()

	if *showVersion {
		fmt.Println(FullVersion())
		os.Exit(0)
	}
}

func DatabaseDir() string {
	cleanedPath := filepath.Clean(*databaseDir)
	if _, err := os.Stat(filepath.Dir(cleanedPath)); err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "FATAL ERROR: path for database does not exist: %s\n", filepath.Dir(cleanedPath))
		} else {
			fmt.Fprintf(os.Stderr, "FATAL ERROR: error accessing database path (%s): %s\n", filepath.Dir(cleanedPath), err)
		}
		os.Exit(1)
	}
	return cleanedPath
}

func LogLevel() string {
	return *logLevel
}

func FileLogLevels() string {
	return *fileLevels
}
