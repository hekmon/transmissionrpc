package TransmissionRPC

import (
	"fmt"
	"reflect"
)

var validTorrentFields []string

func init() {
	val := reflect.ValueOf(Torrent{})
	for i := 0; i < val.Type().NumField(); i++ {
		validTorrentFields = append(validTorrentFields, val.Type().Field(i).Tag.Get("json"))
	}
}

// TorrentGetAll returns all the known fields for all the torrents
func (c *Controller) TorrentGetAll() (torrents []*Torrent, err error) {
	// Send already validated fields to the low level fx
	return c.torrentGet(validTorrentFields, nil)
}

// TorrentGet returns the given of fields (mandatory) for each ids (optionnal)
func (c *Controller) TorrentGet(fields []string, ids []int) (torrents []*Torrent, err error) {
	// Validate fields
	var fieldInvalid bool
	var knownField string
	for _, inputField := range fields {
		fieldInvalid = true
		for _, knownField = range validTorrentFields {
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
	// Forward to real method
	return c.torrentGet(fields, ids)
}

func (c *Controller) torrentGet(fields []string, ids []int) (torrents []*Torrent, err error) {
	// Prepare
	arguments := torrentGetParams{
		Fields: fields,
		IDs:    ids,
	}
	var result torrentGetResults
	// Execute
	if err = c.rpcCall("torrent-get", &arguments, &result); err != nil {
		err = fmt.Errorf("'torrent-get' rpc method failed: %v", err)
		return
	}
	torrents = result.Torrents
	return
}

type torrentGetParams struct {
	Fields []string `json:"fields"`
	IDs    []int    `json:"ids,omitempty"`
}

type torrentGetResults struct {
	Torrents []*Torrent `json:"torrents"`
}

// Torrent represents all the possible fields of data for a torrent
// All fields are pointers to avoid unecessary memory allocation when retreiving only some fields
// but also to differentiate non present field from their default value
type Torrent struct {
	ActivityDate      *int               `json:"activityDate"`
	AddedDate         *int               `json:"addedDate"`
	BandwidthPriority *int               `json:"bandwidthPriority"`
	Comment           *string            `json:"comment"`
	CorruptEver       *int               `json:"corruptEver"`
	Creator           *string            `json:"creator"`
	DateCreated       *int               `json:"dateCreated"`
	DesiredAvailable  *int               `json:"desiredAvailable"`
	DoneDate          *int               `json:"doneDate"`
	DownloadDir       *string            `json:"downloadDir"`
	DownloadedEver    *int               `json:"downloadedEver"`
	DownloadLimit     *int               `json:"downloadLimit"`
	DownloadLimited   *bool              `json:"downloadLimited"`
	Error             *int               `json:"error"`
	ErrorString       *string            `json:"errorString"`
	Eta               *int               `json:"eta"`
	EtaIdle           *int               `json:"etaIdle"`
	Files             []*TorrentFile     `json:"files"`
	FileStats         []*TorrentFileStat `json:"fileStats"`
	HashString        *string            `json:"hashString"`
}

// TorrentFile represent one file from a Torrent
type TorrentFile struct {
	BytesCompleted int    `json:"bytesCompleted"`
	Length         int    `json:"length"`
	Name           string `json:"name"`
}

// TorrentFileStat represents the metadata of a torrent's file
type TorrentFileStat struct {
	BytesCompleted int  `json:"bytesCompleted"`
	Wanted         bool `json:"wanted"`
	Priority       int  `json:"priority"`
}
