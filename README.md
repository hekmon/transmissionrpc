# TransmissionRPC
[![GoDoc](https://godoc.org/github.com/hekmon/transmissionrpc?status.svg)](https://godoc.org/github.com/hekmon/transmissionrpc)

Golang bindings to Transmission (bittorent) RPC interface (Work in Progress).

Even if there is some high level wrappers/helpers, the goal of this lib is to stay close to the original API in terms of methods and payloads while enhancing certain types to be more "golangish": timestamps are converted from/to time.Time, numeric durations in time.Duration, booleans in numeric form are converted to real bool, etc...

Reference:
https://trac.transmissionbt.com/browser/tags/2.92/extras/rpc-spec.txt?rev=14714

## Implementation

First the main client object must be instantiated with [New()](https://godoc.org/github.com/hekmon/transmissionrpc#New). In its basic from only host/ip, username and password must be provided. Default will apply for port (`9091`) rpc URI (`/transmission/rpc`) and others values.

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

Each rpc methods here can work with ID list, hash list or `recently-active` magic word. Therefor, there is 3 golang method variants for each of them.

```golang
transmissionbt.TorrentXXXXIDs(...)
transmissionbt.TorrentXXXXHashes(...)
transmissionbt.TorrentXXXXRecentlyActive()
```

* torrent-start
```golang
err := transmissionbt.TorrentStartIDs([]int64{55})
if err != nil {
	fmt.Fprintln(os.Stderr, err)
} else {
	fmt.Println("yay")
}
```

* torrent-start-now
```golang
err := transmissionbt.TorrentStartNowHashes([]string{"f07e0b0584745b7bcb35e98097488d34e68623d0"})
if err != nil {
	fmt.Fprintln(os.Stderr, err)
} else {
	fmt.Println("yay")
}
```

* torrent-stop
```golang
err := transmissionbt.TorrentStopIDs([]int64{55})
if err != nil {
	fmt.Fprintln(os.Stderr, err)
} else {
	fmt.Println("yay")
}
```

* torrent-verify
```golang
err := transmissionbt.TorrentVerifyHashes([]string{"f07e0b0584745b7bcb35e98097488d34e68623d0"})
if err != nil {
	fmt.Fprintln(os.Stderr, err)
} else {
	fmt.Println("yay")
}
```

* torrent-reannounce
```golang
err := transmissionbt.TorrentReannounceRecentlyActive()
if err != nil {
	fmt.Fprintln(os.Stderr, err)
} else {
	fmt.Println("yay")
}
```

#### Torrent Mutators

* torrent-set
Example: apply a 1Mo/s limit to a torrent.
```golang
uploadLimited := true
uploadLimitKbps := int64(1024)
err := transmissionbt.TorrentSet(&transmissionrpc.TorrentSetPayload{
	IDs:           []int64{55},
	UploadLimited: &uploadLimited,
	UploadLimit:   &uploadLimitKbps,
})
if err != nil {
	fmt.Fprintln(os.Stderr, err)
} else {
	fmt.Println("yay")
}
```

There is a lot more [mutators](https://godoc.org/github.com/hekmon/transmissionrpc#TorrentSetPayload).

#### Torrent Accessors

* torrent-get _(done)_

#### Adding a Torrent

* torrent-add

Adding a torrent from a file (using [TorrentAddFile](https://godoc.org/github.com/hekmon/transmissionrpc#Client.TorrentAddFile) wrapper):

```golang
filepath := "/home/hekmon/Downloads/ubuntu-17.10.1-desktop-amd64.iso.torrent"
torrent, err := transmissionbt.TorrentAddFile(filepath)
if err != nil {
	fmt.Fprintln(os.Stderr, err)
} else {
	// Only 3 fields will be returned/set in the Torrent struct
	fmt.Println(*torrent.ID)
	fmt.Println(*torrent.Name)
	fmt.Println(*torrent.HashString)
}
```

Adding a torrent from an URL (ex: a magnet) with the real [TorrentAdd](https://godoc.org/github.com/hekmon/transmissionrpc#Client.TorrentAdd) method:

```golang
magnet := "magnet:?xt=urn:btih:f07e0b0584745b7bcb35e98097488d34e68623d0&dn=ubuntu-17.10.1-desktop-amd64.iso&tr=http%3A%2F%2Ftorrent.ubuntu.com%3A6969%2Fannounce&tr=http%3A%2F%2Fipv6.torrent.ubuntu.com%3A6969%2Fannounce"
torrent, err := btserv.TorrentAdd(&transmissionrpc.TorrentAddPayload{
	Filename: &magnet,
})
if err != nil {
	fmt.Fprintln(os.Stderr, err)
} else {
	// Only 3 fields will be returned/set in the Torrent struct
	fmt.Println(*torrent.ID)
	fmt.Println(*torrent.Name)
	fmt.Println(*torrent.HashString)
}
```

Which would output:
```
55
ubuntu-17.10.1-desktop-amd64.iso
f07e0b0584745b7bcb35e98097488d34e68623d0
```

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
