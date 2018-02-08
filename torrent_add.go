package TransmissionRPC

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
)

/*
	Adding a Torrent
	https://trac.transmissionbt.com/browser/tags/2.92/extras/rpc-spec.txt?rev=14714#L356
*/

// TorrentAddFile is wrapper to directly add a torrent file (it handles the base64 encoding
// and payload generation). If successfull (torrent added or duplicate) torrent return value
// will only have HashString, ID and Name fields set up.
func (c *Controller) TorrentAddFile(filepath string) (torrent *Torrent, err error) {
	// Validate
	if filepath == "" {
		err = errors.New("filepath can't be empty")
		return
	}
	// Get base64 encoded file content
	b64, err := file2Base64(filepath)
	if err != nil {
		err = fmt.Errorf("can't encode '%s' content as base64: %v", filepath, err)
		return
	}
	// Prepare and send payload
	return c.TorrentAdd(&TorrentAddPayload{MetaInfo: &b64})
}

// TorrentAdd allows to send an Add payload. If successfull (torrent added or duplicate) torrent
// return value will only have HashString, ID and Name fields set up.
func (c *Controller) TorrentAdd(payload *TorrentAddPayload) (torrent *Torrent, err error) {
	// Validate
	if payload == nil {
		err = errors.New("payload can't be nil")
		return
	}
	if payload.Filename == nil && payload.MetaInfo == nil {
		err = errors.New("Filename and MetaInfo can't be both nil")
		return
	}
	// Send payload
	var result torrentAddAnswer
	if err = c.rpcCall("torrent-add", payload, &result); err != nil {
		err = fmt.Errorf("'torrent-add' rpc method failed: %v", err)
		return
	}
	// Extract results
	if result.TorrentAdded != nil {
		torrent = result.TorrentAdded
	} else if result.TorrentDuplicate != nil {
		torrent = result.TorrentDuplicate
	} else {
		err = errors.New("RPC call went fine but neither 'torrent-added' nor 'torrent-duplicate' result payload were found")
	}
	return
}

// TorrentAddPayload represents the data to send in order to add a torrent.
// https://trac.transmissionbt.com/browser/tags/2.92/extras/rpc-spec.txt?rev=14714#L362
type TorrentAddPayload struct {
	Cookies           *string `json:"cookies,omitempty"`
	DownloadDir       *string `json:"download-dir,omitempty"`
	Filename          *string `json:"filename,omitempty"`
	MetaInfo          *string `json:"metainfo,omitempty"`
	Paused            *bool   `json:"paused,omitempty"`
	PeerLimit         *int64  `json:"peer-limit,omitempty"`
	BandwidthPriority *int64  `json:"bandwidthPriority,omitempty"` // mm can't it be equals to 0 ? need custom marshalling with reflect
	FilesWanted       []int64 `json:"files-wanted,omitempty"`
	FilesUnwanted     []int64 `json:"files-unwanted,omitempty"`
	PriorityHigh      []int64 `json:"priority-high,omitempty"`
	PriorityLow       []int64 `json:"priority-low,omitempty"`
	PriorityNormal    []int64 `json:"priority-normal,omitempty"`
}

type torrentAddAnswer struct {
	TorrentAdded     *Torrent `json:"torrent-added"`
	TorrentDuplicate *Torrent `json:"torrent-duplicate"`
}

func file2Base64(filename string) (b64 string, err error) {
	// Try to open file
	file, err := os.Open(filename)
	if err != nil {
		err = fmt.Errorf("open error: %v", err)
		return
	}
	defer file.Close()
	// Prepare encoder
	buffer := new(bytes.Buffer)
	encoder := base64.NewEncoder(base64.StdEncoding, buffer)
	defer encoder.Close()
	// Read file & encode
	if _, err = io.Copy(encoder, file); err != nil {
		err = fmt.Errorf("can't copy file content into the base64 encoder: %v", err)
	}
	// Read it
	b64 = buffer.String()
	return
}
