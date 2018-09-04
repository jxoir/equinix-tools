# Equinix CLI Tools for ECX & ECP [Unofficial]

<!-- toc -->
- [Overview](#overview)
- [Installation](#installation)
<!-- tocstop -->

## Overview

An **UNOFFICIAL** GO CLI for ECX and ECP Tested with Go 1.10+

:warning: WARNING: This CLI is **NOT official**, What does this mean?

* There is no formal Equinix [support] (https://developer.equinix.com) for this SDK at this point
* Bugs may or may not get fixed
* Not all SDK features may be implemented and implemented features may be buggy or incorrect
* Only implements Buyer API _for now_

- [ ] ECX CLI
   - [ ] Buyer API
   - Metros
   - [x] List metros
   - Connections
   - [x] List connections
   - [x] Get connection by uuid
   - [ ] Validate authorization key
   - [ ] Create a connection
   - [ ] Delete a connection
   - [ ] Modify a connection
   - Routing Instance
   - Connector
   - Subscription
   - Bundle Offering
   - Public IPBlock
   - Buyer Preferences

## Installation

Make sure you have a working Go environment.  Go version 1.10+ is supported.  [See
the install instructions for Go](http://golang.org/doc/install.html).

To install cli, simply run:
```
$ go get github.com/jxoir/equinix-tools
```

Make sure your `PATH` includes the `$GOPATH/bin` directory so your commands can
be easily used:
```sh
export PATH=$PATH:$GOPATH/bin
````

Supported env vars

```sh
export ECX_API_HOST="api.equinix.com"
export ECX_API_USER="yourapiuser@yourdomain.com"
export ECX_API_USER_PASSWORD="yourapipassword"
export EQUINIX_API_ID="yourAppId"
export EQUINIX_API_SECRET="yourSecret"
```