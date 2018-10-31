# Equinix CLI Tools for ECX & ECP [Unofficial]

<!-- toc -->
- [Overview](#overview)
- [Installation](#installation)
- [Playground](#Playground)
<!-- tocstop -->

## Overview

An **UNOFFICIAL** GO CLI for ECX and ECP Tested with Go 1.10+

:warning: WARNING: This CLI is **NOT official**, What does this mean?

* There is no formal Equinix [support] for this CLI at this point
* Bugs may or may not get fixed
* Not all API features may be implemented and implemented features may be buggy or incorrect
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
   - [x] Delete a connection
   - [ ] Modify a connection
   - [x] Seller services list/fetch
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
$ go get github.com/jxoir/equinix-tools/...
```

Example use:
```
$ ecxctl connections list
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

## Playground

In order to use playground endpoint you should use the "playground-token" flag with the token.

```
ecxctl connections list --playground-token=xxxxxxxxxxxx
````

# Filtering

Basic filtering options available (connections initially)

Key/Value filtering

To filter a connection list by connection name with value "something"

Filtering *only* works with one filter and *doesn't* traverse nested structures

```
ecxctl connections list --filter=Key=name,Value=something
```

## Connections

List available connections
```
ecxctl connections list
```
List connections filtered 
```
ecxctl connections list --filter=Key=name,Value=something
```


Retrieve connection details (uuid as argument, no need to flag --uuid)
```
ecxctl connections get xxxxxxxxx-xxxxxxxx-xxxxxxx-xxxxxxx
```

Delete (use --uuid flag, security measure)
```
ecxctl connections delete --uuid=xxxxxxxxx-xxxxxxxx-xxxxxxx-xxxxxxx
```
