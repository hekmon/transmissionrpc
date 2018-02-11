# TransmissionRPC
[![GoDoc](https://godoc.org/github.com/hekmon/transmissionrpc?status.svg)](https://godoc.org/github.com/hekmon/transmissionrpc)

Golang bindings to Transmission (bittorent) RPC interface

Work in Progress

Even if there is some high level wrappers/helpers, the goal of this lib is to stay close to the original API in terms of methods and payloads while enhancing certain types to be more "golangish": timestamps are converted from/to time.Time, numeric durations in time.Duration, booleans in numeric form are converted to real bool, etc...

Reference:
https://trac.transmissionbt.com/browser/tags/2.92/extras/rpc-spec.txt?rev=14714

## Implementation

First the main client object must be [instancied](https://godoc.org/github.com/hekmon/transmissionrpc#New). In its basic from only host/ip, username and password must be provided. Default will apply for port (`9091`) rpc URI ('`/transmission/rpc`) and others values.

```golang
    transmissionbt := transmissionrpc.New("127.0.0.1", "rpcuser", "rpcpass", nil)
```

But advanced values can also be configured to your liking using [AdvancedConfig](https://godoc.org/github.com/hekmon/transmissionrpc#AdvancedConfig). 
Each value of `AdvancedConfig` with a type default value will be replaced by the lib default value, so you can set only the ones you want:

```golang
	transmissionbt := transmissionrpc.New("bt.mydomain.net", "rpcuser", "rpcpass",
		&transmissionrpc.AdvancedConfig{
			HTTPS: true,
			Port:  443,
		})
```

### Torrent Requests

#### Torrent Action Requests

* torrent-start _(done)_
* torrent-start-now _(done)_
* torrent-stop _(done)_
* torrent-verify _(done)_
* torrent-reannounce _(done)_

#### Torrent Mutators

* torrent-set _(done)_

#### Torrent Accessors

* torrent-get _(done)_

#### Adding a Torrent

* torrent-add _(done)_

#### Removing a Torrent

* torrent-remove _(done)_

#### Moving a Torrent

* torrent-set-location _(done)_

#### Renaming a Torrent's Path

* torrent-rename-path _(to do)_

### Session Requests

#### Mutators

* session-set _(to do)_

#### Accessors

* session-get _(to do)_

#### Session Statistics

* session-stats _(to do)_

#### Blocklist

* blocklist-update _(to do)_

#### Port Checking

* port-test _(to do)_

#### Session shutdown

* session-close _(to do)_

#### Queue Movement Requests

* queue-move-top _(to do)_
* queue-move-up _(to do)_
* queue-move-down _(to do)_
* queue-move-bottom _(to do)_

#### Free Space

* free-space _(to do)_
