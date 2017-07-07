[![Travis](https://img.shields.io/travis/Safing/safing-core.svg?style=flat-square)](https://travis-ci.org/Safing/safing-core)
[![Coveralls](https://img.shields.io/coveralls/Safing/safing-core.svg?branch=master&style=flat-square)](https://coveralls.io/github/Safing/safing-core?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/Safing/safing-core?style=flat-square)](https://goreportcard.com/report/github.com/Safing/safing-core)

# Safing Core

Safing aims to protect your privacy online. Safing is a software that you install on your device and will control the network traffic to stop communication of apps that hurt your privacy.  
You can think of it like Browser Add-Ons such as Adblock, Ghostery or Disconnect, but instead of only protecting your browser, it protects your whole system. But this is just the beginning. Check out our [Roadmap](https://safing.me/#roadmap).

_Technically speaking_ Safing is designed to protect you from possible unwanted connections made by applications you use. It, however, does assume that every application that is executed is at least partially trusted. Wherever Safing is in the position to provide a clean solution to also increase the security of the system, it will do so.

The Safing Core is the core component of the Safing system. Others are the [Safing UI](#) and [Safing Notify](#).

## Current Status

Safing is now in a tech preview phase (v0.0.x), where the first features are completed and we want to open them up to the community to get feedback on the system. It is not yet ready for day to day use and should only be used to play around with the new concept.

## Download

Currently Safing is only supported on Linux. Windows and Mac will be supported by the end of 2017.
There is no installer yet, you will have to execute the components you need directly:
- Core: TODO
- UI: TODO
- Notify (Tray Icon): TODO

## User Guide

To learn more about how Safing works and how you can use it protect your privacy, please read the [User Guide](#).

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

## Building

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
