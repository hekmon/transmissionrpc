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
	IDs    []int    `json:"ids,omitempty"`
}

type torrentGetResults struct {
	Torrents []*Torrent `json:"torrents"`
}

// Torrent represents all the possible fields of data for a torrent
// All fields are pointers to detect if the value is nil (field not requested) or default real default value
// https://trac.transmissionbt.com/browser/tags/2.92/extras/rpc-spec.txt?rev=14714#L148
type Torrent struct {
	ActivityDate            *int               `json:"activityDate"`
	AddedDate               *int               `json:"addedDate"`
	BandwidthPriority       *int               `json:"bandwidthPriority"`
	Comment                 *string            `json:"comment"`
	CorruptEver             *int               `json:"corruptEver"`
	Creator                 *string            `json:"creator"`
	DateCreated             *int               `json:"dateCreated"`
	DesiredAvailable        *int               `json:"desiredAvailable"`
	DoneDate                *int               `json:"doneDate"`
	DownloadDir             *string            `json:"downloadDir"`
	DownloadedEver          *int               `json:"downloadedEver"`
	DownloadLimit           *int               `json:"downloadLimit"`
	DownloadLimited         *bool              `json:"downloadLimited"`
	Error                   *int               `json:"error"`
	ErrorString             *string            `json:"errorString"`
	Eta                     *int               `json:"eta"`
	EtaIdle                 *int               `json:"etaIdle"`
	Files                   []*TorrentFile     `json:"files"`
	FileStats               []*TorrentFileStat `json:"fileStats"`
	HashString              *string            `json:"hashString"`
	HaveUnchecked           *int               `json:"haveUnchecked"`
	HaveValid               *int               `json:"haveValid"`
	HonorsSessionLimits     *bool              `json:"honorsSessionLimits"`
	ID                      *int               `json:"id"`
	IsFinished              *bool              `json:"isFinished"`
	IsPrivate               *bool              `json:"isPrivate"`
	IsStalled               *bool              `json:"isStalled"`
	LeftUntilDone           *int               `json:"leftUntilDone"`
	MagnetLink              *string            `json:"magnetLink"`
	ManualAnnounceTime      *int               `json:"manualAnnounceTime"`
	MaxConnectedPeers       *int               `json:"maxConnectedPeers"`
	MetadataPercentComplete *float64           `json:"metadataPercentComplete"`
	Name                    *string            `json:"name"`
	PeerLimit               *int               `json:"peer-limit"`
	Peers                   []*Peer            `json:"peers"`
}

// TorrentFile represent one file from a Torrent
// https://trac.transmissionbt.com/browser/tags/2.92/extras/rpc-spec.txt?rev=14714#L221
type TorrentFile struct {
	BytesCompleted int    `json:"bytesCompleted"`
	Length         int    `json:"length"`
	Name           string `json:"name"`
}

// TorrentFileStat represents the metadata of a torrent's file
// https://trac.transmissionbt.com/browser/tags/2.92/extras/rpc-spec.txt?rev=14714#L227
type TorrentFileStat struct {
	BytesCompleted int  `json:"bytesCompleted"`
	Wanted         bool `json:"wanted"`
	Priority       int  `json:"priority"`
}

// Peer represent a peer metadata of a torrent's peer list
// https://trac.transmissionbt.com/browser/tags/2.92/extras/rpc-spec.txt?rev=14714#L235
type Peer struct {
	Address            string  `json:"address"`
	ClientName         string  `json:"clientName"`
	ClientIsChoked     bool    `json:"clientIsChoked"`
	ClientIsInterested bool    `json:"clientIsInterested"`
	FlagStr            string  `json:"flagStr"`
	IsDownloadingFrom  bool    `json:"isDownloadingFrom"`
	IsEncrypted        bool    `json:"isEncrypted"`
	IsIncoming         bool    `json:"isIncoming"`
	IsUploadingTo      bool    `json:"isUploadingTo"`
	IsUTP              bool    `json:"isUTP"`
	PeerIsChoked       bool    `json:"peerIsChoked"`
	PeerIsInterested   bool    `json:"peerIsInterested"`
	Port               int     `json:"port"`
	Progress           float64 `json:"progress"`
	RateToClient       int     `json:"rateToClient"`
	RateToPeer         int     `json:"rateToPeer"`
}
