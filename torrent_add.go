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

// TorrentAddPayload represents the data to send in order to add a torrent.
// https://trac.transmissionbt.com/browser/tags/2.92/extras/rpc-spec.txt?rev=14714#L362
type TorrentAddPayload struct {
	Cookies           *string `json:"cookies,omitempty"`
	DownloadDir       *string `json:"download-dir,omitempty"`
	Filename          *string `json:"filename,omitempty"`
	MetaInfo          *string `json:"metainfo,omitempty"`
	Paused            *bool   `json:"paused,omitempty"`
	PeerLimit         *int64  `json:"peer-limit,omitempty"`
	BandwidthPriority *int64  `json:"bandwidthPriority,omitempty"`
	FilesWanted       []int64 `json:"files-wanted,omitempty"`
	FilesUnwanted     []int64 `json:"files-unwanted,omitempty"`
	PriorityHigh      []int64 `json:"priority-high,omitempty"`
	PriorityLow       []int64 `json:"priority-low,omitempty"`
	PriorityNormal    []int64 `json:"priority-normal,omitempty"`
}

// TorrentAdd allows to send an Add payload.
// If successfull (torrent added or duplicate) torrent return value will have only
// HashString, ID and Name set up.
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

type torrentAddAnswer struct {
	TorrentAdded     *Torrent `json:"torrent-added"`
	TorrentDuplicate *Torrent `json:"torrent-duplicate"`
}

// TorrentAddFile is wrapper to directly add a torrent file (it handles the base64 encoding).
// If successfull (torrent added or duplicate) torrent return value will have only
// HashString, ID and Name set up.
func (c *Controller) TorrentAddFile(filename string) (torrent *Torrent, err error) {
	// Validate
	if filename == "" {
		err = errors.New("filename can't be empty")
		return
	}
	// Get base64 encoded file content
	b64, err := file2Base64(filename)
	if err != nil {
		err = fmt.Errorf("can't encode '%s' content as base64: %v", filename, err)
		return
	}
	// Prepare and send payload
	return c.TorrentAdd(&TorrentAddPayload{MetaInfo: &b64})
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
