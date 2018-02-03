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
func (c *Controller) TorrentGet(fields []string, ids []int64) (torrents []*Torrent, err error) {
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

func (c *Controller) torrentGet(fields []string, ids []int64) (torrents []*Torrent, err error) {
	arguments := torrentGetParams{
		Fields: fields,
		IDs:    ids,
	}
	var result torrentGetResults
	if err = c.rpcCall("torrent-get", &arguments, &result); err != nil {
		err = fmt.Errorf("'torrent-get' rpc method failed: %v", err)
		return
	}
	torrents = result.Torrents
	return
}

type torrentGetParams struct {
	Fields []string `json:"fields"`
	IDs    []int64  `json:"ids,omitempty"`
}

type torrentGetResults struct {
	Torrents []*Torrent `json:"torrents"`
}

// Torrent represents all the possible fields of data for a torrent
// All fields are point64ers to detect if the value is nil (field not requested) or default real default value
// https://trac.transmissionbt.com/browser/tags/2.92/extras/rpc-spec.txt?rev=14714#L148
type Torrent struct {
	ActivityDate            *int64             `json:"activityDate"`
	AddedDate               *int64             `json:"addedDate"`
	BandwidthPriority       *int64             `json:"bandwidthPriority"`
	Comment                 *string            `json:"comment"`
	CorruptEver             *int64             `json:"corruptEver"`
	Creator                 *string            `json:"creator"`
	DateCreated             *int64             `json:"dateCreated"`
	DesiredAvailable        *int64             `json:"desiredAvailable"`
	DoneDate                *int64             `json:"doneDate"`
	DownloadDir             *string            `json:"downloadDir"`
	DownloadedEver          *int64             `json:"downloadedEver"`
	DownloadLimit           *int64             `json:"downloadLimit"`
	DownloadLimited         *bool              `json:"downloadLimited"`
	Error                   *int64             `json:"error"`
	ErrorString             *string            `json:"errorString"`
	Eta                     *int64             `json:"eta"`
	EtaIdle                 *int64             `json:"etaIdle"`
	Files                   []*TorrentFile     `json:"files"`
	FileStats               []*TorrentFileStat `json:"fileStats"`
	HashString              *string            `json:"hashString"`
	HaveUnchecked           *int64             `json:"haveUnchecked"`
	HaveValid               *int64             `json:"haveValid"`
	HonorsSessionLimits     *bool              `json:"honorsSessionLimits"`
	ID                      *int64             `json:"id"`
	IsFinished              *bool              `json:"isFinished"`
	IsPrivate               *bool              `json:"isPrivate"`
	IsStalled               *bool              `json:"isStalled"`
	LeftUntilDone           *int64             `json:"leftUntilDone"`
	MagnetLink              *string            `json:"magnetLink"`
	ManualAnnounceTime      *int64             `json:"manualAnnounceTime"`
	MaxConnectedPeers       *int64             `json:"maxConnectedPeers"`
	MetadataPercentComplete *float64           `json:"metadataPercentComplete"`
	Name                    *string            `json:"name"`
	PeerLimit               *int64             `json:"peer-limit"`
	Peers                   []*Peer            `json:"peers"`
}

// TorrentFile represent one file from a Torrent
// https://trac.transmissionbt.com/browser/tags/2.92/extras/rpc-spec.txt?rev=14714#L221
type TorrentFile struct {
	BytesCompleted int64  `json:"bytesCompleted"`
	Length         int64  `json:"length"`
	Name           string `json:"name"`
}

// TorrentFileStat represents the metadata of a torrent's file
// https://trac.transmissionbt.com/browser/tags/2.92/extras/rpc-spec.txt?rev=14714#L227
type TorrentFileStat struct {
	BytesCompleted int64 `json:"bytesCompleted"`
	Wanted         bool  `json:"wanted"`
	Priority       int64 `json:"priority"`
}

// Peer represent a peer metadata of a torrent's peer list
// https://trac.transmissionbt.com/browser/tags/2.92/extras/rpc-spec.txt?rev=14714#L235
type Peer struct {
	Address              string  `json:"address"`
	ClientName           string  `json:"clientName"`
	ClientIsChoked       bool    `json:"clientIsChoked"`
	ClientIsint64erested bool    `json:"clientIsint64erested"`
	FlagStr              string  `json:"flagStr"`
	IsDownloadingFrom    bool    `json:"isDownloadingFrom"`
	IsEncrypted          bool    `json:"isEncrypted"`
	IsIncoming           bool    `json:"isIncoming"`
	IsUploadingTo        bool    `json:"isUploadingTo"`
	IsUTP                bool    `json:"isUTP"`
	PeerIsChoked         bool    `json:"peerIsChoked"`
	PeerIsint64erested   bool    `json:"peerIsint64erested"`
	Port                 int64   `json:"port"`
	Progress             float64 `json:"progress"`
	RateToClient         int64   `json:"rateToClient"`
	RateToPeer           int64   `json:"rateToPeer"`
}
