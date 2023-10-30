package transmissionrpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/hekmon/cunits/v2"
)

/*
	Session Arguments
    https://github.com/transmission/transmission/blob/4.0.3/docs/rpc-spec.md#41-session-arguments
*/

var validSessionFields []string

func init() {
	sessionArgumentsType := reflect.TypeOf(SessionArguments{})
	validSessionFields = make([]string, sessionArgumentsType.NumField())
	for i := 0; i < sessionArgumentsType.NumField(); i++ {
		validSessionFields[i] = sessionArgumentsType.Field(i).Tag.Get("json")
	}
}

type Encryption string

const (
	EncryptionRequired  Encryption = "required"
	EncryptionPreferred Encryption = "prefered"
	EncryptionTolerated Encryption = "tolerated"
)

// SessionArguments represents all the global/session values.
type SessionArguments struct {
	AltSpeedDown                     *int64      `json:"alt-speed-down"`                       // max global download speed (KBps)
	AltSpeedEnabled                  *bool       `json:"alt-speed-enabled"`                    // true means use the alt speeds
	AltSpeedTimeBegin                *int64      `json:"alt-speed-time-begin"`                 // when to turn on alt speeds (units: minutes after midnight)
	AltSpeedTimeDay                  *int64      `json:"alt-speed-time-day"`                   // what day(s) to turn on alt speeds (look at tr_sched_day)
	AltSpeedTimeEnabled              *bool       `json:"alt-speed-time-enabled"`               // true means the scheduled on/off times are used
	AltSpeedTimeEnd                  *int64      `json:"alt-speed-time-end"`                   // when to turn off alt speeds (units: same)
	AltSpeedUp                       *int64      `json:"alt-speed-up"`                         // max global upload speed (KBps)
	BlocklistEnabled                 *bool       `json:"blocklist-enabled"`                    // true means enabled
	BlocklistSize                    *int64      `json:"blocklist-size"`                       // number of rules in the blocklist
	BlocklistURL                     *string     `json:"blocklist-url"`                        // location of the blocklist to use for "blocklist-update"
	CacheSizeMB                      *int64      `json:"cache-size-mb"`                        // maximum size of the disk cache (MB)
	ConfigDir                        *string     `json:"config-dir"`                           // location of transmission's configuration directory
	DefaultTrackers                  []string    `json:"default-trackers"`                     // list of default trackers to use on public torrents
	DHTEnabled                       *bool       `json:"dht-enabled"`                          // true means allow dht in public torrents
	DownloadDir                      *string     `json:"download-dir"`                         // default path to download torrents
	DownloadQueueEnabled             *bool       `json:"download-queue-enabled"`               // if true, limit how many torrents can be downloaded at once
	DownloadQueueSize                *int64      `json:"download-queue-size"`                  // max number of torrents to download at once (see download-queue-enabled)
	Encryption                       *Encryption `json:"encryption"`                           // "required", "preferred", "tolerated", see Encryption type constants
	IdleSeedingLimitEnabled          *bool       `json:"idle-seeding-limit-enabled"`           // true if the seeding inactivity limit is honored by default
	IdleSeedingLimit                 *int64      `json:"idle-seeding-limit"`                   // torrents we're seeding will be stopped if they're idle for this long
	IncompleteDirEnabled             *bool       `json:"incomplete-dir-enabled"`               // true means keep torrents in incomplete-dir until done
	IncompleteDir                    *string     `json:"incomplete-dir"`                       // path for incomplete torrents, when enabled
	LPDEnabled                       *bool       `json:"lpd-enabled"`                          // true means allow Local Peer Discovery in public torrents
	PeerLimitGlobal                  *int64      `json:"peer-limit-global"`                    // maximum global number of peers
	PeerLimitPerTorrent              *int64      `json:"peer-limit-per-torrent"`               // maximum global number of peers
	PeerPortRandomOnStart            *bool       `json:"peer-port-random-on-start"`            // true means pick a random peer port on launch
	PeerPort                         *int64      `json:"peer-port"`                            // port number
	PEXEnabled                       *bool       `json:"pex-enabled"`                          // true means allow pex in public torrents
	PortForwardingEnabled            *bool       `json:"port-forwarding-enabled"`              // true means enabled
	QueueStalledEnabled              *bool       `json:"queue-stalled-enabled"`                // whether or not to consider idle torrents as stalled
	QueueStalledMinutes              *int64      `json:"queue-stalled-minutes"`                // torrents that are idle for N minuets aren't counted toward seed-queue-size or download-queue-size
	RenamePartialFiles               *bool       `json:"rename-partial-files"`                 // true means append ".part" to incomplete files
	RPCVersionMinimum                *int64      `json:"rpc-version-minimum"`                  // the minimum RPC API version supported
	RPCVersionSemVer                 *string     `json:"rpc-version-semver"`                   // the current RPC API version in a semver-compatible string
	RPCVersion                       *int64      `json:"rpc-version"`                          // the current RPC API version
	ScriptTorrentAddedEnabled        *bool       `json:"script-torrent-added-enabled"`         // whether or not to call the added script
	ScriptTorrentAddedFilename       *string     `json:"script-torrent-added-filename"`        //filename of the script to run
	ScriptTorrentDoneEnabled         *bool       `json:"script-torrent-done-enabled"`          // whether or not to call the "done" script
	ScriptTorrentDoneFilename        *string     `json:"script-torrent-done-filename"`         // filename of the script to run
	ScriptTorrentDoneSeedingEnabled  *bool       `json:"script-torrent-done-seeding-enabled"`  // whether or not to call the seeding-done script
	ScriptTorrentDoneSeedingFilename *string     `json:"script-torrent-done-seeding-filename"` // filename of the script to run
	SeedQueueEnabled                 *bool       `json:"seed-queue-enabled"`                   // if true, limit how many torrents can be uploaded at once
	SeedQueueSize                    *int64      `json:"seed-queue-size"`                      // max number of torrents to uploaded at once (see seed-queue-enabled)
	SeedRatioLimit                   *float64    `json:"seedRatioLimit"`                       // the default seed ratio for torrents to use
	SeedRatioLimited                 *bool       `json:"seedRatioLimited"`                     // true if seedRatioLimit is honored by default
	SessionID                        *string     `json:"session-id"`                           // the current session ID
	SpeedLimitDownEnabled            *bool       `json:"speed-limit-down-enabled"`             // true means enabled
	SpeedLimitDown                   *int64      `json:"speed-limit-down"`                     // max global download speed (KBps)
	SpeedLimitUpEnabled              *bool       `json:"speed-limit-up-enabled"`               // true means enabled
	SpeedLimitUp                     *int64      `json:"speed-limit-up"`                       // max global upload speed (KBps)
	StartAddedTorrents               *bool       `json:"start-added-torrents"`                 // true means added torrents will be started right away
	TrashOriginalTorrentFiles        *bool       `json:"trash-original-torrent-files"`         // true means the .torrent file of added torrents will be deleted
	Units                            *Units      `json:"units"`                                // see units below
	UTPEnabled                       *bool       `json:"utp-enabled"`                          // true means allow utp
	Version                          *string     `json:"version"`                              // long version string "$version ($revision)"
}

// MarshalJSON allows to marshall into JSON only the non nil fields.
// It differs from 'omitempty' which also skip default values
// (as 0 or false which can be valid here).
func (sa SessionArguments) MarshalJSON() (data []byte, err error) {
	// Build an intermediary payload with base types
	type baseSessionArguments SessionArguments
	tmp := struct {
		DefaultTrackers *string `json:"default-trackers"` // list of default trackers to use on public torrents
		*baseSessionArguments
	}{
		baseSessionArguments: (*baseSessionArguments)(&sa),
	}
	if sa.DefaultTrackers != nil {
		oneLineDefaultTrackers := strings.Join(sa.DefaultTrackers, "\n")
		tmp.DefaultTrackers = &oneLineDefaultTrackers
	}
	// Build a payload with only the non nil fields
	sav := reflect.ValueOf(tmp)
	sat := sav.Type()
	cleanPayload := make(map[string]interface{}, sat.NumField())
	var currentValue, nestedStruct, currentNestedValue reflect.Value
	var currentStructField, currentNestedStructField reflect.StructField
	var j int
	for i := 0; i < sav.NumField(); i++ {
		currentValue = sav.Field(i)
		currentStructField = sat.Field(i)
		if !currentValue.IsNil() {
			if currentStructField.Name == "baseSessionArguments" {
				// inherited/nested struct
				nestedStruct = reflect.Indirect(currentValue)
				for j = 0; j < nestedStruct.NumField(); j++ {
					currentNestedValue = nestedStruct.Field(j)
					currentNestedStructField = nestedStruct.Type().Field(j)
					if !currentNestedValue.IsNil() {
						JSONKeyName := currentNestedStructField.Tag.Get("json")
						if JSONKeyName != "-" {
							cleanPayload[JSONKeyName] = currentNestedValue.Interface()
						}
					}
				}
			} else {
				// Overloaded field
				cleanPayload[currentStructField.Tag.Get("json")] = currentValue.Interface()
			}
		}
	}
	// Marshall the clean payload
	return json.Marshal(cleanPayload)
}

// UnmarshalJSON allows to convert timestamps to golang time.Time values.
func (sa *SessionArguments) UnmarshalJSON(data []byte) (err error) {
	// Shadow real type for regular unmarshalling
	type RawSessionArguments SessionArguments
	tmp := &struct {
		DefaultTrackers *string `json:"default-trackers"` // list of default trackers to use on public torrents
		*RawSessionArguments
	}{
		RawSessionArguments: (*RawSessionArguments)(sa),
	}
	if err = json.Unmarshal(data, &tmp); err != nil {
		return
	}
	// Custom
	if tmp.DefaultTrackers != nil {
		sa.DefaultTrackers = strings.Split(*tmp.DefaultTrackers, "\n")
	}
	return
}

// CacheSize returns the cache size in a handy format
func (sa SessionArguments) CacheSize() (size cunits.Bits) {
	if sa.CacheSizeMB != nil {
		size = cunits.ImportInMB(float64(*sa.CacheSizeMB))
	}
	return
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

// SessionArgumentsGetAll returns global/session values.
// https://github.com/transmission/transmission/blob/4.0.3/docs/rpc-spec.md#412-accessors
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
// https://github.com/transmission/transmission/blob/4.0.3/docs/rpc-spec.md#412-accessors
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

// SessionArgumentsSet allows to modify global/session values.
// https://github.com/transmission/transmission/blob/4.0.3/docs/rpc-spec.md#411-mutators
func (c *Client) SessionArgumentsSet(ctx context.Context, payload SessionArguments) (err error) {
	// Sanitize fields that can not be set
	payload.BlocklistSize = nil
	payload.ConfigDir = nil
	payload.RPCVersionMinimum = nil
	payload.RPCVersionSemVer = nil
	payload.RPCVersion = nil
	payload.SessionID = nil
	payload.Units = nil
	payload.Version = nil
	// Exec
	if err = c.rpcCall(ctx, "session-set", payload, nil); err != nil {
		err = fmt.Errorf("'session-set' rpc method failed: %w", err)
	}
	return
}
