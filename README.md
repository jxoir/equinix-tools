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
   - [x] Create a L2 connection to a Seller profile (specific use case for AWS/Azure,Others)
   - [ ] Create a generic L2 connection
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

Create L2 connection to seller service (shortcut to establish a simple connection to AWS initially)

- Required flags
  - name - user specified name for the new connection
  - port-uuid - user port to establish the connection to
  - vlan - user side VLAN for the specific connection (primary connection)
  - auth-key - specific seller authorization key, in the case of AWS is the Account ID, Azure auth key must be retrieved from Azure portal
  - seller-uuid - seller specific uuid (can be retrieved with seller command)
  - seller-region - specific seller param, in the case of AWS is the destination region ex.: eu-west-1
  - seller-metro - seller metro to connect to, some sellers allows to use a "remote" connection (incurring in extra charges)
  - speed - speed for the connection, must be allowed by the platform and seller (can be retrieved with seller command)
  - speed-unit - MB / GB, must be allowed by the platform and the seller (can be retrieved with seller command)
  - notifications-email - email for notifications

### Create Connection Flowchart

```
+------------------------+               +---------------------------+                           +------------------------------------+
|                        | Get Service   |                           |  Get Service              |                                    |
|   User Initiates the   | Profile       |  Seller                   |  Profiles for Customer    |  Buyer                             |
|   Creation Process     | +-----------> |  /serviceprofiles/{uuid}  | +-----------------------> |  /serviceprofiles/services/{uuid}  |
|                        |               |                           |                           |                                    |
+------------------------+               +---------------------------+                           +------------------------------------+

                                                                                                                          +         +
                                                                                                                          |         |
                                                                                                                          |         |
                                                                                          Get Available Connection        |         |
                                                                                          Tiers for the Profile           |         |
                                                                                                                          |         |
                                                                                        +-------------------------+       |         |
                                                                                        |                         |       |         |
                                                                                        |  /common/billingTiers/  | <-----+         |
         +-------------------------------+                                              |                         |                 |
         | NO                            |                                              +-------------------------+                 |
         |                               |                                                                                          |
         +                               v                                                                                          |
                                                                                                                                    |
+-----------------+        +-----------------------------------------+   Validate Keys  +---------------------------+               |
|                 |        |                                         |   with Provider  |                           | Get Service   |
| Are Keys Valid? | <----+ |  Buyer                                  |                  |  Seller                   | Profile       |
|                 |        |  /connections/validateAuthorizationKey  | <--------------+ |  /serviceprofiles/{uuid}  |               |
+-----------------+        |                                         |                  |                           | <-------------+
                           +-----------------------------------------+                  +---------------------------+
   YES  +
        |                                                   ^
        |                                                   |
        v                                                   |
                                                            |
 +---------------+                                          |
 | Buyer         |                                          |
 | /connections/ |                                          |
 |               |                                          |
 +---------------+                                          |
                                                            |
        +                                                   |
        |                                                   |
        |                                                   |
        v                                                   |
                                                            |
  +-------------+  YES                                      |
  | Any Errors? |  +----------------------------------------+
  +-------------+

    NO  +
        |            +--------------+
        |            |  Connection  |
        +--------->  |  Created     |
                     +--------------+
```

### Create connection to AWS

```sh
ecxctl connections create --name=EQUINIX_DEMO_CONN --port-uuid=2813d8f6-4623-4a5c-9c71-34de7e100933 --seller-metro=LD --seller-region=eu-west-1 --seller-uuid=9b460b5a-5461-4186-a3d5-2e8d8fb4c91b 
 --speed=50 --speed-unit=MB --vlan=3022 --auth-key=12345678912 --notifications-email=some@email.com
```

### Create connection to Azure

```sh
ecxctl connections create --name=EQUINIX_DEMO_CONN_AZ --name-sec=EQUINIX_DEMO_CONN_AZ_SEC --port-uuid=66284add-49a3-9a30-b4e0-30ac094f8af1 --port-uuid-sec=66284add-49a5-9a50-b4e0-30ac094f8af1 --seller-metro=LD --seller-region=westeurope --seller-uuid=a1390b22-bbe0-4e93-ad37-85beef9d254d --speed=50 --named-tag=Microsoft --speed-unit=MB --vlan=3143 --vlan-sec=3143 --auth-key=12345678912 --notifications-email=some@email.com
```

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
