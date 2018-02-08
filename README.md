# TransmissionRPC
[![GoDoc](https://godoc.org/github.com/hekmon/transmissionrpc?status.svg)](https://godoc.org/github.com/hekmon/transmissionrpc)

Golang bindings to Transmission (bittorent) RPC interface

Work in Progress

Based on:
https://trac.transmissionbt.com/browser/tags/2.92/extras/rpc-spec.txt?rev=14714

## Progress

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

* torrent-remove _(to do)_

#### Moving a Torrent

* torrent-set-location _(to do)_

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
