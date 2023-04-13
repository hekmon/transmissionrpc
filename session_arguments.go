package transmissionrpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/hekmon/cunits/v2"
)

/*
	Session Arguments
    https://github.com/transmission/transmission/blob/4.0.2/docs/rpc-spec.md#41-session-arguments
*/

var validSessionFields []string

func init() {
	torrentType := reflect.TypeOf(SessionArguments{})
	for i := 0; i < torrentType.NumField(); i++ {
		validSessionFields = append(validSessionFields, torrentType.Field(i).Tag.Get("json"))
	}
}

// RPCVersion returns true if the lib RPC version is greater or equals to the remote server rpc minimum version.
func (c *Client) RPCVersion(ctx context.Context) (ok bool, serverVersion int64, serverMinimumVersion int64, err error) {
	payload, err := c.SessionArgumentsGetAll(ctx)
	if err != nil {
		err = fmt.Errorf("can't get session values: %w", err)
		return
	}
	if payload.RPCVersion == nil {
		err = errors.New("payload RPC Version is nil")
		return
	}
	if payload.RPCVersionMinimum == nil {
		err = errors.New("payload RPC Version minimum is nil")
		return
	}
	serverVersion = *payload.RPCVersion
	serverMinimumVersion = *payload.RPCVersionMinimum
	ok = RPCVersion >= serverMinimumVersion
	return
}

// SessionArgumentsSet allows to modify global/session values.
// https://github.com/transmission/transmission/blob/4.0.2/docs/rpc-spec.md#411-mutators
func (c *Client) SessionArgumentsSet(ctx context.Context, payload SessionArguments) (err error) {
	payload.BlocklistSize = nil
	payload.ConfigDir = nil
	payload.RPCVersion = nil
	payload.RPCVersionMinimum = nil
	payload.SessionID = nil
	payload.Version = nil
	// Exec
	if err = c.rpcCall(ctx, "session-set", payload, nil); err != nil {
		err = fmt.Errorf("'session-set' rpc method failed: %w", err)
	}
	return
}

// SessionArgumentsGetAll returns global/session values.
// https://github.com/transmission/transmission/blob/4.0.2/docs/rpc-spec.md#412-accessors
func (c *Client) SessionArgumentsGetAll(ctx context.Context) (sessionArgs SessionArguments, err error) {
	if err = c.rpcCall(ctx, "session-get", nil, &sessionArgs); err != nil {
		err = fmt.Errorf("'session-get' rpc method failed: %w", err)
	}
	return
}

type sessionGetParams struct {
	Fields []string `json:"fields"`
}

// SessionArgumentsGet returns global/session values for specified fields.
// See the JSON tags of the SessionArguments struct for valid fields.
// https://github.com/transmission/transmission/blob/4.0.2/docs/rpc-spec.md#412-accessors
func (c *Client) SessionArgumentsGet(ctx context.Context, fields []string) (sessionArgs SessionArguments, err error) {
	if err = c.validateSessionFields(fields); err != nil {
		return
	}
	if err = c.rpcCall(ctx, "session-get", sessionGetParams{Fields: fields}, &sessionArgs); err != nil {
		err = fmt.Errorf("'session-get' rpc method failed: %w", err)
	}
	return
}

func (c *Client) validateSessionFields(fields []string) (err error) {
	// Validate fields
	var fieldInvalid bool
	var knownField string
	for _, inputField := range fields {
		fieldInvalid = true
		for _, knownField = range validSessionFields {
			if inputField == knownField {
				fieldInvalid = false
				break
			}
		}
		if fieldInvalid {
			err = fmt.Errorf("field '%s' is invalid", inputField)
			return
		}
	}
	return
}

// SessionArguments represents all the global/session values.
type SessionArguments struct {
	AltSpeedDown                     *int64   `json:"alt-speed-down"`                       // max global download speed (KBps)
	AltSpeedEnabled                  *bool    `json:"alt-speed-enabled"`                    // true means use the alt speeds
	AltSpeedTimeBegin                *int64   `json:"alt-speed-time-begin"`                 // when to turn on alt speeds (units: minutes after midnight)
	AltSpeedTimeEnabled              *bool    `json:"alt-speed-time-enabled"`               // true means the scheduled on/off times are used
	AltSpeedTimeEnd                  *int64   `json:"alt-speed-time-end"`                   // when to turn off alt speeds (units: same)
	AltSpeedTimeDay                  *int64   `json:"alt-speed-time-day"`                   // what day(s) to turn on alt speeds (look at tr_sched_day)
	AltSpeedUp                       *int64   `json:"alt-speed-up"`                         // max global upload speed (KBps)
	BlocklistURL                     *string  `json:"blocklist-url"`                        // location of the blocklist to use for "blocklist-update"
	BlocklistEnabled                 *bool    `json:"blocklist-enabled"`                    // true means enabled
	BlocklistSize                    *int64   `json:"blocklist-size"`                       // number of rules in the blocklist
	CacheSizeMB                      *int64   `json:"cache-size-mb"`                        // maximum size of the disk cache (MB)
	ConfigDir                        *string  `json:"config-dir"`                           // location of transmission's configuration directory
	DefaultTrackers                  *string  `json:"default-trackers"`                     // list of default trackers to use on public torrents
	DownloadDir                      *string  `json:"download-dir"`                         // default path to download torrents
	DownloadQueueSize                *int64   `json:"download-queue-size"`                  // max number of torrents to download at once (see download-queue-enabled)
	DownloadQueueEnabled             *bool    `json:"download-queue-enabled"`               // if true, limit how many torrents can be downloaded at once
	DHTEnabled                       *bool    `json:"dht-enabled"`                          // true means allow dht in public torrents
	Encryption                       *string  `json:"encryption"`                           // "required", "preferred", "tolerated"
	IdleSeedingLimit                 *int64   `json:"idle-seeding-limit"`                   // torrents we're seeding will be stopped if they're idle for this long
	IdleSeedingLimitEnabled          *bool    `json:"idle-seeding-limit-enabled"`           // true if the seeding inactivity limit is honored by default
	IncompleteDir                    *string  `json:"incomplete-dir"`                       // path for incomplete torrents, when enabled
	IncompleteDirEnabled             *bool    `json:"incomplete-dir-enabled"`               // true means keep torrents in incomplete-dir until done
	LPDEnabled                       *bool    `json:"lpd-enabled"`                          // true means allow Local Peer Discovery in public torrents
	PeerLimitGlobal                  *int64   `json:"peer-limit-global"`                    // maximum global number of peers
	PeerLimitPerTorrent              *int64   `json:"peer-limit-per-torrent"`               // maximum global number of peers
	PEXEnabled                       *bool    `json:"pex-enabled"`                          // true means allow pex in public torrents
	PeerPort                         *int64   `json:"peer-port"`                            // port number
	PeerPortRandomOnStart            *bool    `json:"peer-port-random-on-start"`            // true means pick a random peer port on launch
	PortForwardingEnabled            *bool    `json:"port-forwarding-enabled"`              // true means enabled
	QueueStalledEnabled              *bool    `json:"queue-stalled-enabled"`                // whether or not to consider idle torrents as stalled
	QueueStalledMinutes              *int64   `json:"queue-stalled-minutes"`                // torrents that are idle for N minuets aren't counted toward seed-queue-size or download-queue-size
	RenamePartialFiles               *bool    `json:"rename-partial-files"`                 // true means append ".part" to incomplete files
	RPCVersion                       *int64   `json:"rpc-version"`                          // the current RPC API version
	RPCVersionMinimum                *int64   `json:"rpc-version-minimum"`                  // the minimum RPC API version supported
	RPCVersionSemVer                 *string  `json:"rpc-version-semver"`                   // the current RPC API version in a semver-compatible string
	ScriptTorrentAddedEnabled        *bool    `json:"script-torrent-added-enabled"`         // whether or not to call the added script
	ScriptTorrentAddedFilename       *string  `json:"script-torrent-added-filename"`        //filename of the script to run
	ScriptTorrentDoneFilename        *string  `json:"script-torrent-done-filename"`         // filename of the script to run
	ScriptTorrentDoneEnabled         *bool    `json:"script-torrent-done-enabled"`          // whether or not to call the "done" script
	ScriptTorrentDoneSeedingEnabled  *bool    `json:"script-torrent-done-seeding-enabled"`  // whether or not to call the seeding-done script
	ScriptTorrentDoneSeedingFilename *string  `json:"script-torrent-done-seeding-filename"` // filename of the script to run
	SeedRatioLimit                   *float64 `json:"seedRatioLimit"`                       // the default seed ratio for torrents to use
	SeedRatioLimited                 *bool    `json:"seedRatioLimited"`                     // true if seedRatioLimit is honored by default
	SeedQueueSize                    *int64   `json:"seed-queue-size"`                      // max number of torrents to uploaded at once (see seed-queue-enabled)
	SeedQueueEnabled                 *bool    `json:"seed-queue-enabled"`                   // if true, limit how many torrents can be uploaded at once
	SessionID                        *string  `json:"session-id"`                           // the current session ID
	SpeedLimitDown                   *int64   `json:"speed-limit-down"`                     // max global download speed (KBps)
	SpeedLimitDownEnabled            *bool    `json:"speed-limit-down-enabled"`             // true means enabled
	SpeedLimitUp                     *int64   `json:"speed-limit-up"`                       // max global upload speed (KBps)
	SpeedLimitUpEnabled              *bool    `json:"speed-limit-up-enabled"`               // true means enabled
	StartAddedTorrents               *bool    `json:"start-added-torrents"`                 // true means added torrents will be started right away
	TrashOriginalTorrentFiles        *bool    `json:"trash-original-torrent-files"`         // true means the .torrent file of added torrents will be deleted
	Units                            *Units   `json:"units"`                                // see units below
	UTPEnabled                       *bool    `json:"utp-enabled"`                          // true means allow utp
	Version                          *string  `json:"version"`                              // long version string "$version ($revision)"
}

// MarshalJSON allows to marshall into JSON only the non nil fields.
// It differs from 'omitempty' which also skip default values
// (as 0 or false which can be valid here).
func (sa SessionArguments) MarshalJSON() (data []byte, err error) {
	// Build a payload with only the non nil fields
	tspv := reflect.ValueOf(sa)
	tspt := tspv.Type()
	cleanPayload := make(map[string]interface{}, tspt.NumField())
	var currentValue reflect.Value
	var currentStructField reflect.StructField
	for i := 0; i < tspv.NumField(); i++ {
		currentValue = tspv.Field(i)
		currentStructField = tspt.Field(i)
		if !currentValue.IsNil() {
			cleanPayload[currentStructField.Tag.Get("json")] = currentValue.Interface()
		}
	}
	// Marshall the clean payload
	return json.Marshal(cleanPayload)
}

// Units is subset of SessionArguments.
type Units struct {
	SpeedUnits  []string `json:"speed-units"`  // 4 strings: KB/s, MB/s, GB/s, TB/s
	SpeedBytes  int64    `json:"speed-bytes"`  // number of bytes in a KB (1000 for kB; 1024 for KiB)
	SizeUnits   []string `json:"size-units"`   // 4 strings: KB/s, MB/s, GB/s, TB/s
	SizeBytes   int64    `json:"size-bytes"`   // number of bytes in a KB (1000 for kB; 1024 for KiB)
	MemoryUnits []string `json:"memory-units"` // 4 strings: KB/s, MB/s, GB/s, TB/s
	MemoryBytes int64    `json:"memory-bytes"` // number of bytes in a KB (1000 for kB; 1024 for KiB)
}

// GetSpeed returns the speed in a handy format
func (u *Units) GetSpeed() (speed cunits.Bits) {
	return cunits.ImportInByte(float64(u.SpeedBytes))
}

// GetSize returns the size in a handy format
func (u *Units) GetSize() (size cunits.Bits) {
	return cunits.ImportInByte(float64(u.SizeBytes))
}

// GetMemory returns the memory in a handy format
func (u *Units) GetMemory() (memory cunits.Bits) {
	return cunits.ImportInByte(float64(u.MemoryBytes))
}
