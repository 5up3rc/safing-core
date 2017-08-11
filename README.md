[![Travis](https://img.shields.io/travis/Safing/safing-core.svg?style=flat-square)](https://travis-ci.org/Safing/safing-core)
[![Coveralls](https://img.shields.io/coveralls/Safing/safing-core.svg?branch=master&style=flat-square)](https://coveralls.io/github/Safing/safing-core?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/Safing/safing-core?style=flat-square)](https://goreportcard.com/report/github.com/Safing/safing-core)

# Safing Core

Safing strengthens your privacy and security online. It takes control of your network traffic and takes over handling DNS queries. Instead of prompting users if new connections shall be allowed, behavioral profiles are created and assigned to applications. Non-technical users can easily control Safing using [Security Levels](https://github.com/Safing/safing-doc/blob/master/guides/UserGuide.md). Worth mentioning is, that Safing has the ability to verify TLS connections.  
This is just the beginning. Check out our [Roadmap](https://safing.me/#roadmap).

Please note, that Safing is all about securing your network traffic and protecting your privacy, and that it will generally assume that processes that are run on your system are at least partially trusted (i.e. closed source, but not malware). While Safing will protect you from some malware, it is not designed to replace anti-virus solutions, but to work in cooperation. Especially take this into account when evaluating Safing.

The Safing Core is the core component of the Safing system. Others are the [Safing UI](https://github.com/Safing/safing-ui) and [Safing Notify](https://github.com/Safing/safing-notify).

## Current Status

Safing is now in a tech preview phase (v0.0.x), where the first features are completed and we want to open them up to the community to get feedback on the system. It is not yet ready for day to day use and should only be used to play around with the new concept.

## Download

Currently Safing is only supported on Linux. Windows and Mac will be supported by the end of 2017.

We recommend to download a packaged version of all components [here](https://github.com/Safing/safing-installer/releases).  
You can also just [download Safing Core](https://github.com/Safing/safing-core/releases).

## User Guide

To learn more about how Safing works and how you can use it protect your privacy, please read the [User Guide](https://github.com/Safing/safing-doc/blob/master/guides/UserGuide.md).

## Running

As soon as the Safing Core starts, it will intercept network connections based on a default demo profile set.

For builds that interact with the operating system, you must start Safing with operating system privileges (i.e. root on Linux). Depending on the current database configuration, Safing may use the `data` directory in the current working directory to store database data.

    Usage:
      -v	show version and exit
      -db string
          specify custom location for the database (default "data")
      -log string
          set log level [trace|debug|info|warning|error|critical] (default "trace")
      -flog string
          set log level of files: database=trace,firewall=debug

## Developing

Currently Safing is only supported on Linux.  
Builds without specific os interaction (clientapi.go) may succeed on other systems.

1. Install Go 1.8+ ([https://golang.org/dl/](https://golang.org/dl/))
2. Check additional requirements (see below)
3. Build the desired configuration:

        ./build main.go       # builds the full core
        ./build clientapi.go  # builds api service only (for debugging)
        ./build dnsserver.go  # builds nameserver only (no fw)

## Additional Requirements

#### Linux

- go dependencies
  - Some dependencies cannot be vendored, as they then would use internal packages, which is not allowed. You need to get:

        go get golang.org/x/crypto/chacha20poly1305
        go get golang.org/x/net/icmp

- nfqueue-dev
  - Ubuntu: `sudo apt-get install libnetfilter-queue-dev`
  - Arch: `sudo pacman -S libnetfilter_queue`
- Network Manager (should be installed anyway)
  - This requirement only exists until the fallbacks for detection of nameservers and connectivity are ready.
