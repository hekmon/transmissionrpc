package transmissionrpc

/*
	Torrent Accessors
    https://github.com/transmission/transmission/blob/4.0.3/docs/rpc-spec.md#33-torrent-accessor-torrent-get
*/

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/hekmon/cunits/v2"
)

var validTorrentFields []string

func init() {
	torrentType := reflect.TypeOf(Torrent{})
	for i := 0; i < torrentType.NumField(); i++ {
		validTorrentFields = append(validTorrentFields, torrentType.Field(i).Tag.Get("json"))
	}
}

// TorrentGetAll returns all the known fields for all the torrents.
func (c *Client) TorrentGetAll(ctx context.Context) (torrents []Torrent, err error) {
	// Send already validated fields to the low level fx
	return c.torrentGet(ctx, validTorrentFields, nil)
}

// TorrentGetAllFor returns all known fields for the given torrent's ids.
func (c *Client) TorrentGetAllFor(ctx context.Context, ids []int64) (torrents []Torrent, err error) {
	return c.torrentGet(ctx, validTorrentFields, ids)
}

// TorrentGetAllForHashes returns all known fields for the given torrent's ids by string (usually hash).
func (c *Client) TorrentGetAllForHashes(ctx context.Context, hashes []string) (torrents []Torrent, err error) {
	return c.torrentGetHash(ctx, validTorrentFields, hashes)
}

// TorrentGet returns the given of fields (mandatory) for each ids (optionnal).
func (c *Client) TorrentGet(ctx context.Context, fields []string, ids []int64) (torrents []Torrent, err error) {
	if err = c.validateTorrentFields(fields); err != nil {
		return
	}
	return c.torrentGet(ctx, fields, ids)
}

// TorrentGetHashes returns the given of fields (mandatory) for each ids (optionnal).
func (c *Client) TorrentGetHashes(ctx context.Context, fields []string, hashes []string) (torrents []Torrent, err error) {
	if err = c.validateTorrentFields(fields); err != nil {
		return
	}
	return c.torrentGetHash(ctx, fields, hashes)
}

func (c *Client) validateTorrentFields(fields []string) (err error) {
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
	return
}

func (c *Client) torrentGet(ctx context.Context, fields []string, ids []int64) (torrents []Torrent, err error) {
	var result torrentGetResults
	if err = c.rpcCall(ctx, "torrent-get", &torrentGetParams{
		Fields: fields,
		IDs:    ids,
	}, &result); err != nil {
		err = fmt.Errorf("'torrent-get' rpc method failed: %w", err)
		return
	}
	torrents = result.Torrents
	return
}

func (c *Client) torrentGetHash(ctx context.Context, fields []string, hashes []string) (torrents []Torrent, err error) {
	var result torrentGetResults
	if err = c.rpcCall(ctx, "torrent-get", &torrentGetHashParams{
		Fields: fields,
		Hashes: hashes,
	}, &result); err != nil {
		err = fmt.Errorf("'torrent-get' rpc method failed: %w", err)
		return
	}
	torrents = result.Torrents
	return
}

type torrentGetParams struct {
	Fields []string `json:"fields"`
	IDs    []int64  `json:"ids,omitempty"`
}

type torrentGetHashParams struct {
	Fields []string `json:"fields"`
	Hashes []string `json:"ids,omitempty"`
}

type torrentGetResults struct {
	Torrents []Torrent `json:"torrents"`
}

// Torrent represents all the possible fields of data for a torrent.
// All fields are pointers to detect if the value is nil (field not requested) or real default value.
type Torrent struct {
	ActivityDate            *time.Time        `json:"activityDate"`
	AddedDate               *time.Time        `json:"addedDate"`
	Availability            []int64           `json:"availability"` // RPC v17
	BandwidthPriority       *int64            `json:"bandwidthPriority"`
	Comment                 *string           `json:"comment"`
	CorruptEver             *int64            `json:"corruptEver"`
	Creator                 *string           `json:"creator"`
	DateCreated             *time.Time        `json:"dateCreated"`
	DesiredAvailable        *int64            `json:"desiredAvailable"`
	DoneDate                *time.Time        `json:"doneDate"`
	DownloadDir             *string           `json:"downloadDir"`
	DownloadedEver          *int64            `json:"downloadedEver"`
	DownloadLimit           *int64            `json:"downloadLimit"`
	DownloadLimited         *bool             `json:"downloadLimited"`
	EditDate                *time.Time        `json:"editDate"`
	Error                   *int64            `json:"error"`
	ErrorString             *string           `json:"errorString"`
	ETA                     *int64            `json:"eta"`
	ETAIdle                 *int64            `json:"etaIdle"`
	FileCount               *int64            `json:"file-count"` // RPC v17
	Files                   []TorrentFile     `json:"files"`
	FileStats               []TorrentFileStat `json:"fileStats"`
	Group                   *string           `json:"group"` // RPC v17
	HashString              *string           `json:"hashString"`
	HaveUnchecked           *int64            `json:"haveUnchecked"`
	HaveValid               *int64            `json:"haveValid"`
	HonorsSessionLimits     *bool             `json:"honorsSessionLimits"`
	ID                      *int64            `json:"id"`
	IsFinished              *bool             `json:"isFinished"`
	IsPrivate               *bool             `json:"isPrivate"`
	IsStalled               *bool             `json:"isStalled"`
	Labels                  []string          `json:"labels"` // RPC v16
	LeftUntilDone           *int64            `json:"leftUntilDone"`
	MagnetLink              *string           `json:"magnetLink"`
	ManualAnnounceTime      *int64            `json:"manualAnnounceTime"`
	MaxConnectedPeers       *int64            `json:"maxConnectedPeers"`
	MetadataPercentComplete *float64          `json:"metadataPercentComplete"`
	Name                    *string           `json:"name"`
	PeerLimit               *int64            `json:"peer-limit"`
	Peers                   []Peer            `json:"peers"`
	PeersConnected          *int64            `json:"peersConnected"`
	PeersFrom               *TorrentPeersFrom `json:"peersFrom"`
	PeersGettingFromUs      *int64            `json:"peersGettingFromUs"`
	PeersSendingToUs        *int64            `json:"peersSendingToUs"`
	PercentComplete         *float64          `json:"percentComplete"` // RPC v17
	PercentDone             *float64          `json:"percentDone"`
	Pieces                  *string           `json:"pieces"`
	PieceCount              *int64            `json:"pieceCount"`
	PieceSize               *cunits.Bits      `json:"PieceSize"`
	Priorities              []int64           `json:"priorities"`
	PrimaryMimeType         *string           `json:"primary-mime-type"` // RPC v17
	QueuePosition           *int64            `json:"queuePosition"`
	RateDownload            *int64            `json:"rateDownload"` // B/s
	RateUpload              *int64            `json:"rateUpload"`   // B/s
	RecheckProgress         *float64          `json:"recheckProgress"`
	TimeDownloading         *time.Duration    `json:"secondsDownloading"`
	TimeSeeding             *time.Duration    `json:"secondsSeeding"`
	SeedIdleLimit           *time.Duration    `json:"seedIdleLimit"`
	SeedIdleMode            *int64            `json:"seedIdleMode"`
	SeedRatioLimit          *float64          `json:"seedRatioLimit"`
	SeedRatioMode           *SeedRatioMode    `json:"seedRatioMode"`
	SizeWhenDone            *cunits.Bits      `json:"sizeWhenDone"`
	StartDate               *time.Time        `json:"startDate"`
	Status                  *TorrentStatus    `json:"status"`
	Trackers                []Tracker         `json:"trackers"`
	TrackerList             *string           `json:"trackerList"`
	TrackerStats            []TrackerStats    `json:"trackerStats"`
	TotalSize               *cunits.Bits      `json:"totalSize"`
	TorrentFile             *string           `json:"torrentFile"`
	UploadedEver            *int64            `json:"uploadedEver"`
	UploadLimit             *int64            `json:"uploadLimit"`
	UploadLimited           *bool             `json:"uploadLimited"`
	UploadRatio             *float64          `json:"uploadRatio"`
	Wanted                  []bool            `json:"wanted"`
	WebSeeds                []string          `json:"webseeds"`
	WebSeedsSendingToUs     *int64            `json:"webseedsSendingToUs"`
}

// ConvertDownloadSpeed will return the download speed as cunits.Bits/second
func (t *Torrent) ConvertDownloadSpeed() (speed cunits.Bits) {
	if t.RateDownload != nil {
		speed = cunits.ImportInByte(float64(*t.RateDownload))
	}
	return
}

// ConvertUploadSpeed will return the upload speed as cunits.Bits/second
func (t *Torrent) ConvertUploadSpeed() (speed cunits.Bits) {
	if t.RateUpload != nil {
		speed = cunits.ImportInByte(float64(*t.RateUpload))
	}
	return
}

// UnmarshalJSON allows to convert timestamps to golang time.Time values.
func (t *Torrent) UnmarshalJSON(data []byte) (err error) {
	// Shadow real type for regular unmarshalling
	type RawTorrent Torrent
	tmp := &struct {
		ActivityDate       *int64  `json:"activityDate"`
		AddedDate          *int64  `json:"addedDate"`
		DateCreated        *int64  `json:"dateCreated"`
		DoneDate           *int64  `json:"doneDate"`
		EditDate           *int64  `json:"editDate"`
		PieceSize          *int64  `json:"pieceSize"`
		SecondsDownloading *int64  `json:"secondsDownloading"`
		SecondsSeeding     *int64  `json:"secondsSeeding"`
		SeedIdleLimit      *int64  `json:"seedIdleLimit"`
		SizeWhenDone       *int64  `json:"sizeWhenDone"`
		StartDate          *int64  `json:"startDate"`
		TotalSize          *int64  `json:"totalSize"`
		Wanted             []int64 `json:"wanted"` // boolean in number form
		*RawTorrent
	}{
		RawTorrent: (*RawTorrent)(t),
	}
	// Unmarshal (with timestamps as number)
	if err = json.Unmarshal(data, &tmp); err != nil {
		return
	}
	// Create the real time & duration from timsteamps and seconds
	if tmp.ActivityDate != nil {
		ad := time.Unix(*tmp.ActivityDate, 0)
		t.ActivityDate = &ad
	}
	if tmp.AddedDate != nil {
		ad := time.Unix(*tmp.AddedDate, 0)
		t.AddedDate = &ad
	}
	if tmp.DateCreated != nil {
		dc := time.Unix(*tmp.DateCreated, 0)
		t.DateCreated = &dc
	}
	if tmp.DoneDate != nil {
		dd := time.Unix(*tmp.DoneDate, 0)
		t.DoneDate = &dd
	}
	if tmp.EditDate != nil {
		dd := time.Unix(*tmp.EditDate, 0)
		t.EditDate = &dd
	}
	if tmp.PieceSize != nil {
		ps := cunits.ImportInByte(float64(*tmp.PieceSize))
		t.PieceSize = &ps
	}
	if tmp.SecondsDownloading != nil {
		dur := time.Duration(*tmp.SecondsDownloading) * time.Second
		t.TimeDownloading = &dur
	}
	if tmp.SecondsSeeding != nil {
		dur := time.Duration(*tmp.SecondsSeeding) * time.Second
		t.TimeSeeding = &dur
	}
	if tmp.SeedIdleLimit != nil {
		dur := time.Duration(*tmp.SeedIdleLimit) * time.Minute
		t.SeedIdleLimit = &dur
	}
	if tmp.SizeWhenDone != nil {
		swd := cunits.ImportInByte(float64(*tmp.SizeWhenDone))
		t.SizeWhenDone = &swd
	}
	if tmp.StartDate != nil {
		st := time.Unix(*tmp.StartDate, 0)
		t.StartDate = &st
	}
	if tmp.TotalSize != nil {
		ts := cunits.ImportInByte(float64(*tmp.TotalSize))
		t.TotalSize = &ts
	}
	// Boolean slice in decimal form
	if tmp.Wanted != nil {
		t.Wanted = make([]bool, len(tmp.Wanted))
		for index, value := range tmp.Wanted {
			if value == 1 {
				t.Wanted[index] = true
			} else if value != 0 {
				return fmt.Errorf("can't convert wanted index %d value '%d' as boolean", index, value)
			}
		}
	}
	return
}

// MarshalJSON allows to convert back golang values to original payload values.
func (t Torrent) MarshalJSON() (data []byte, err error) {
	// Shadow real type for regular unmarshalling
	type RawTorrent Torrent
	tmp := &struct {
		ActivityDate       *int64  `json:"activityDate"`
		AddedDate          *int64  `json:"addedDate"`
		DateCreated        *int64  `json:"dateCreated"`
		DoneDate           *int64  `json:"doneDate"`
		SecondsDownloading *int64  `json:"secondsDownloading"`
		SecondsSeeding     *int64  `json:"secondsSeeding"`
		SeedIdleLimit      *int64  `json:"seedIdleLimit"`
		StartDate          *int64  `json:"startDate"`
		Wanted             []int64 `json:"wanted"` // boolean in number form
		*RawTorrent
	}{
		RawTorrent: (*RawTorrent)(&t),
	}
	// Timestamps & Duration
	if t.ActivityDate != nil {
		ad := t.ActivityDate.Unix()
		tmp.ActivityDate = &ad
	}
	if t.AddedDate != nil {
		ad := t.AddedDate.Unix()
		tmp.AddedDate = &ad
	}
	if t.DateCreated != nil {
		dc := t.DateCreated.Unix()
		tmp.DateCreated = &dc
	}
	if t.DoneDate != nil {
		dd := t.DoneDate.Unix()
		tmp.DoneDate = &dd
	}
	if t.TimeDownloading != nil {
		sd := int64(*t.TimeDownloading / time.Second)
		tmp.SecondsDownloading = &sd
	}
	if t.TimeSeeding != nil {
		ss := int64(*t.TimeSeeding / time.Second)
		tmp.SecondsSeeding = &ss
	}
	if t.SeedIdleLimit != nil {
		sil := int64(*t.SeedIdleLimit / time.Minute)
		tmp.SeedIdleLimit = &sil
	}
	if t.StartDate != nil {
		st := t.StartDate.Unix()
		tmp.StartDate = &st
	}
	// Boolean as number
	if t.Wanted != nil {
		tmp.Wanted = make([]int64, len(t.Wanted))
		for index, value := range t.Wanted {
			if value {
				tmp.Wanted[index] = 1
			}
		}
	}
	// Marshall original values within the tmp payload
	return json.Marshal(&tmp)
}

// TorrentFile represent one file from a Torrent.
type TorrentFile struct {
	BytesCompleted int64  `json:"bytesCompleted"`
	Length         int64  `json:"length"`
	Name           string `json:"name"`
}

// TorrentFileStat represents the metadata of a torrent's file.
type TorrentFileStat struct {
	BytesCompleted int64 `json:"bytesCompleted"`
	Wanted         bool  `json:"wanted"`
	Priority       int64 `json:"priority"`
}

// Peer represent a peer metadata of a torrent's peer list.
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
	Port               int64   `json:"port"`
	Progress           float64 `json:"progress"`
	RateToClient       int64   `json:"rateToClient"` // B/s
	RateToPeer         int64   `json:"rateToPeer"`   // B/s
}

// ConvertDownloadSpeed will return the download speed from peer as cunits.Bits/second
func (p *Peer) ConvertDownloadSpeed() (speed cunits.Bits) {
	return cunits.ImportInByte(float64(p.RateToClient))
}

// ConvertUploadSpeed will return the upload speed to peer as cunits.Bits/second
func (p *Peer) ConvertUploadSpeed() (speed cunits.Bits) {
	return cunits.ImportInByte(float64(p.RateToPeer))
}

// TorrentPeersFrom represents the peers statistics of a torrent.
type TorrentPeersFrom struct {
	FromCache    int64 `json:"fromCache"`
	FromDHT      int64 `json:"fromDht"`
	FromIncoming int64 `json:"fromIncoming"`
	FromLPD      int64 `json:"fromLpd"`
	FromLTEP     int64 `json:"fromLtep"`
	FromPEX      int64 `json:"fromPex"`
	FromTracker  int64 `json:"fromTracker"`
}

// SeedRatioMode represents a torrent current seeding mode
type SeedRatioMode int64

const (
	// SeedRatioModeGlobal represents the use of the global ratio for a torrent
	SeedRatioModeGlobal SeedRatioMode = 0
	// SeedRatioModeCustom represents the use of a custom ratio for a torrent
	SeedRatioModeCustom SeedRatioMode = 1
	// SeedRatioModeNoRatio represents the absence of ratio for a torrent
	SeedRatioModeNoRatio SeedRatioMode = 2
)

func (srm SeedRatioMode) String() string {
	switch srm {
	case SeedRatioModeGlobal:
		return "global"
	case SeedRatioModeCustom:
		return "custom"
	case SeedRatioModeNoRatio:
		return "no ratio"
	default:
		return "<unknown>"
	}
}

// GoString implements the GoStringer interface from the stdlib fmt package
func (srm SeedRatioMode) GoString() string {
	switch srm {
	case SeedRatioModeGlobal:
		return fmt.Sprintf("global (%d)", srm)
	case SeedRatioModeCustom:
		return fmt.Sprintf("custom (%d)", srm)
	case SeedRatioModeNoRatio:
		return fmt.Sprintf("no ratio (%d)", srm)
	default:
		return fmt.Sprintf("<unknown> (%d)", srm)
	}
}

// TorrentStatus binds torrent status to a status code
type TorrentStatus int64

const (
	// TorrentStatusStopped represents a stopped torrent
	TorrentStatusStopped TorrentStatus = 0
	// TorrentStatusCheckWait represents a torrent queued for files checking
	TorrentStatusCheckWait TorrentStatus = 1
	// TorrentStatusCheck represents a torrent which files are currently checked
	TorrentStatusCheck TorrentStatus = 2
	// TorrentStatusDownloadWait represents a torrent queue to download
	TorrentStatusDownloadWait TorrentStatus = 3
	// TorrentStatusDownload represents a torrent currently downloading
	TorrentStatusDownload TorrentStatus = 4
	// TorrentStatusSeedWait represents a torrent queued to seed
	TorrentStatusSeedWait TorrentStatus = 5
	// TorrentStatusSeed represents a torrent currently seeding
	TorrentStatusSeed TorrentStatus = 6
	// TorrentStatusIsolated represents a torrent which can't find peers
	TorrentStatusIsolated TorrentStatus = 7
)

func (status TorrentStatus) String() string {
	switch status {
	case TorrentStatusStopped:
		return "stopped"
	case TorrentStatusCheckWait:
		return "waiting to check files"
	case TorrentStatusCheck:
		return "checking files"
	case TorrentStatusDownloadWait:
		return "waiting to download"
	case TorrentStatusDownload:
		return "downloading"
	case TorrentStatusSeedWait:
		return "waiting to seed"
	case TorrentStatusSeed:
		return "seeding"
	case TorrentStatusIsolated:
		return "can't find peers"
	default:
		return "<unknown>"
	}
}

// GoString implements the GoStringer interface from the stdlib fmt package
func (status TorrentStatus) GoString() string {
	switch status {
	case TorrentStatusStopped:
		return fmt.Sprintf("stopped (%d)", status)
	case TorrentStatusCheckWait:
		return fmt.Sprintf("waiting to check files (%d)", status)
	case TorrentStatusCheck:
		return fmt.Sprintf("checking files (%d)", status)
	case TorrentStatusDownloadWait:
		return fmt.Sprintf("waiting to download (%d)", status)
	case TorrentStatusDownload:
		return fmt.Sprintf("downloading (%d)", status)
	case TorrentStatusSeedWait:
		return fmt.Sprintf("waiting to seed (%d)", status)
	case TorrentStatusSeed:
		return fmt.Sprintf("seeding (%d)", status)
	case TorrentStatusIsolated:
		return fmt.Sprintf("can't find peers (%d)", status)
	default:
		return fmt.Sprintf("<unknown> (%d)", status)
	}
}

// Tracker represent the base data of a torrent's tracker.
type Tracker struct {
	Announce string `json:"announce"`
	ID       int64  `json:"id"`
	Scrape   string `json:"scrape"`
	SiteName string `json:"sitename"`
	Tier     int64  `json:"tier"`
}

// TrackerStats represent the extended data of a torrent's tracker.
type TrackerStats struct {
	Announce              string    `json:"announce"`
	AnnounceState         int64     `json:"announceState"`
	DownloadCount         int64     `json:"downloadCount"`
	HasAnnounced          bool      `json:"hasAnnounced"`
	HasScraped            bool      `json:"hasScraped"`
	Host                  string    `json:"host"`
	ID                    int64     `json:"id"`
	IsBackup              bool      `json:"isBackup"`
	LastAnnouncePeerCount int64     `json:"lastAnnouncePeerCount"`
	LastAnnounceResult    string    `json:"lastAnnounceResult"`
	LastAnnounceStartTime time.Time `json:"-"`
	LastAnnounceSucceeded bool      `json:"lastAnnounceSucceeded"`
	LastAnnounceTime      time.Time `json:"-"`
	LastAnnounceTimedOut  bool      `json:"lastAnnounceTimedOut"`
	LastScrapeResult      string    `json:"lastScrapeResult"`
	LastScrapeStartTime   time.Time `json:"-"`
	LastScrapeSucceeded   bool      `json:"lastScrapeSucceeded"`
	LastScrapeTime        time.Time `json:"-"`
	LastScrapeTimedOut    bool      `json:"-"` // should be boolean but number. Will be converter in UnmarshalJSON
	LeecherCount          int64     `json:"leecherCount"`
	NextAnnounceTime      time.Time `json:"-"`
	NextScrapeTime        time.Time `json:"-"`
	Scrape                string    `json:"scrape"`
	ScrapeState           int64     `json:"scrapeState"`
	SiteName              string    `json:"sitename"`
	SeederCount           int64     `json:"seederCount"`
	Tier                  int64     `json:"tier"`
}

// UnmarshalJSON allows to convert timestamps to golang time.Time values.
func (ts *TrackerStats) UnmarshalJSON(data []byte) (err error) {
	// Shadow real type for regular unmarshalling
	type RawTrackerStats TrackerStats
	tmp := struct {
		LastAnnounceStartTime int64       `json:"lastAnnounceStartTime"`
		LastAnnounceTime      int64       `json:"lastAnnounceTime"`
		LastScrapeStartTime   int64       `json:"lastScrapeStartTime"`
		LastScrapeTime        int64       `json:"lastScrapeTime"`
		LastScrapeTimedOut    interface{} `json:"lastScrapeTimedOut"`
		NextAnnounceTime      int64       `json:"nextAnnounceTime"`
		NextScrapeTime        int64       `json:"nextScrapeTime"`
		*RawTrackerStats
	}{
		RawTrackerStats: (*RawTrackerStats)(ts),
	}
	// Unmarshal (with timestamps as number)
	if err = json.Unmarshal(data, &tmp); err != nil {
		return
	}
	// Convert to real boolean
	// transmission rpc version 15 returns 1 or 0
	// transmission rpc version 16 returns true or false
	switch v := tmp.LastScrapeTimedOut.(type) {
	case bool:
		ts.LastScrapeTimedOut = v
	case float64:
		if v == 1. {
			ts.LastScrapeTimedOut = true
		} else if v == 0. {
			ts.LastScrapeTimedOut = false
		} else {
			return fmt.Errorf("can't convert 'lastScrapeTimedOut' value '%v' into boolean", tmp.LastScrapeTimedOut)
		}
	default:
		return fmt.Errorf("can't convert 'lastScrapeTimedOut' value '%v' into boolean", tmp.LastScrapeTimedOut)
	}
	// Create the real time value from the timestamps
	ts.LastAnnounceStartTime = time.Unix(tmp.LastAnnounceStartTime, 0)
	ts.LastAnnounceTime = time.Unix(tmp.LastAnnounceTime, 0)
	ts.LastScrapeStartTime = time.Unix(tmp.LastScrapeStartTime, 0)
	ts.LastScrapeTime = time.Unix(tmp.LastScrapeTime, 0)
	ts.NextAnnounceTime = time.Unix(tmp.NextAnnounceTime, 0)
	ts.NextScrapeTime = time.Unix(tmp.NextScrapeTime, 0)
	return
}

// MarshalJSON allows to convert back golang values to original payload values.
func (ts TrackerStats) MarshalJSON() (data []byte, err error) {
	// Shadow real type for regular unmarshalling
	type RawTrackerStats TrackerStats
	tmp := struct {
		LastAnnounceStartTime int64 `json:"lastAnnounceStartTime"`
		LastAnnounceTime      int64 `json:"lastAnnounceTime"`
		LastScrapeStartTime   int64 `json:"lastScrapeStartTime"`
		LastScrapeTime        int64 `json:"lastScrapeTime"`
		LastScrapeTimedOut    int64 `json:"lastScrapeTimedOut"`
		NextAnnounceTime      int64 `json:"nextAnnounceTime"`
		NextScrapeTime        int64 `json:"nextScrapeTime"`
		*RawTrackerStats
	}{
		LastAnnounceStartTime: ts.LastAnnounceStartTime.Unix(),
		LastAnnounceTime:      ts.LastAnnounceTime.Unix(),
		LastScrapeStartTime:   ts.LastScrapeStartTime.Unix(),
		LastScrapeTime:        ts.LastScrapeTime.Unix(),
		NextAnnounceTime:      ts.NextAnnounceTime.Unix(),
		NextScrapeTime:        ts.NextScrapeTime.Unix(),
		RawTrackerStats:       (*RawTrackerStats)(&ts),
	}
	// Convert real bool to its number form
	if ts.LastScrapeTimedOut {
		tmp.LastScrapeTimedOut = 1
	}
	// MarshalJSON allows to convert back golang values to original payload values
	return json.Marshal(&tmp)
}
