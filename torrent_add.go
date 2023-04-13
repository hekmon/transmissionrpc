package transmissionrpc

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
)

/*
	Adding a Torrent
    https://github.com/transmission/transmission/blob/4.0.2/docs/rpc-spec.md#34-adding-a-torrent
*/

// TorrentAddFileDownloadDir is wrapper to directly add a torrent file (it handles the base64 encoding
// and payload generation) to a DownloadDir (not the default download dir). If successful (torrent added
// or duplicate) torrent return value will only have HashString, ID and Name fields set up.
func (c *Client) TorrentAddFileDownloadDir(ctx context.Context, filepath, downloaddir string) (torrent Torrent, err error) {
	// Validate filepath
	if filepath == "" {
		err = errors.New("filepath can't be empty")
		return
	}
	// Validate downloaddir
	if downloaddir == "" {
		err = errors.New("downloaddir can't be empty")
		return
	}
	// Get base64 encoded file content
	b64, err := File2Base64(filepath)
	if err != nil {
		err = fmt.Errorf("can't encode '%s' content as base64: %w", filepath, err)
		return
	}
	// Prepare and send payload
	return c.TorrentAdd(ctx, TorrentAddPayload{MetaInfo: &b64, DownloadDir: &downloaddir})
}

// TorrentAddFile is wrapper to directly add a torrent file (it handles the base64 encoding
// and payload generation). If successful (torrent added or duplicate) torrent return value
// will only have HashString, ID and Name fields set up.
func (c *Client) TorrentAddFile(ctx context.Context, filepath string) (torrent Torrent, err error) {
	// Validate
	if filepath == "" {
		err = errors.New("filepath can't be empty")
		return
	}
	// Get base64 encoded file content
	b64, err := File2Base64(filepath)
	if err != nil {
		err = fmt.Errorf("can't encode '%s' content as base64: %w", filepath, err)
		return
	}
	// Prepare and send payload
	return c.TorrentAdd(ctx, TorrentAddPayload{MetaInfo: &b64})
}

// TorrentAdd allows to send an Add payload. If successful (torrent added or duplicate) torrent
// return value will only have HashString, ID and Name fields set up.
func (c *Client) TorrentAdd(ctx context.Context, payload TorrentAddPayload) (torrent Torrent, err error) {
	// Validate
	if payload.Filename == nil && payload.MetaInfo == nil {
		err = errors.New("fields Filename and MetaInfo can't be both nil")
		return
	}
	// Send payload
	var result torrentAddAnswer
	if err = c.rpcCall(ctx, "torrent-add", payload, &result); err != nil {
		err = fmt.Errorf("'torrent-add' rpc method failed: %w", err)
		return
	}
	// Extract results
	if result.TorrentAdded != nil {
		torrent = *result.TorrentAdded
	} else if result.TorrentDuplicate != nil {
		torrent = *result.TorrentDuplicate
	} else {
		err = errors.New("RPC call went fine but neither 'torrent-added' nor 'torrent-duplicate' result payload were found")
	}
	return
}

// TorrentAddPayload represents the data to send in order to add a torrent.
type TorrentAddPayload struct {
	Cookies           *string `json:"cookies"`           // pointer to a string of one or more cookies
	DownloadDir       *string `json:"download-dir"`      // path to download the torrent to
	Filename          *string `json:"filename"`          // filename or URL of the .torrent file
	MetaInfo          *string `json:"metainfo"`          // base64-encoded .torrent content
	Paused            *bool   `json:"paused"`            // if true, don't start the torrent
	PeerLimit         *int64  `json:"peer-limit"`        // maximum number of peers
	BandwidthPriority *int64  `json:"bandwidthPriority"` // torrent's bandwidth tr_priority_t
	FilesWanted       []int64 `json:"files-wanted"`      // indices of file(s) to download
	FilesUnwanted     []int64 `json:"files-unwanted"`    // indices of file(s) to not download
	PriorityHigh      []int64 `json:"priority-high"`     // indices of high-priority file(s)
	PriorityLow       []int64 `json:"priority-low"`      // indices of low-priority file(s)
	PriorityNormal    []int64 `json:"priority-normal"`   // indices of normal-priority file(s)
}

// MarshalJSON allows to marshall into JSON only the non nil fields.
// It differs from 'omitempty' which also skip default values
// (as 0 or false which can be valid here).
func (tap TorrentAddPayload) MarshalJSON() (data []byte, err error) {
	// Build a payload with only the non nil fields
	tspv := reflect.ValueOf(tap)
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

type torrentAddAnswer struct {
	TorrentAdded     *Torrent `json:"torrent-added"`
	TorrentDuplicate *Torrent `json:"torrent-duplicate"`
}

// File2Base64 returns the base64 encoding of the file provided by filename.
// This can then be passed as MetaInfo in TorrentAddPayload.
func File2Base64(filename string) (b64 string, err error) {
	// Try to open file
	file, err := os.Open(filename)
	if err != nil {
		err = fmt.Errorf("can't open file: %w", err)
		return
	}
	defer file.Close()
	// Prepare encoder
	buffer := new(bytes.Buffer)
	encoder := base64.NewEncoder(base64.StdEncoding, buffer)
	// Stream file to the encoder
	if _, err = io.Copy(encoder, file); err != nil {
		err = fmt.Errorf("can't copy file content into the base64 encoder: %w", err)
		return
	}
	// Flush last bytes
	if err = encoder.Close(); err != nil {
		err = fmt.Errorf("can't flush last bytes of the base64 encoder: %w", err)
		return
	}
	// Get the string form
	b64 = buffer.String()
	return
}
