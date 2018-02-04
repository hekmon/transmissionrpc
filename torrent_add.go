package TransmissionRPC

/*
	Adding a Torrent
	https://trac.transmissionbt.com/browser/tags/2.92/extras/rpc-spec.txt?rev=14714#L356
*/

// TorrentAddPayload represents the data to send in order to add a torrent.
// https://trac.transmissionbt.com/browser/tags/2.92/extras/rpc-spec.txt?rev=14714#L362
type TorrentAddPayload struct {
	Cookies           *string `json:"cookies"`
	DownloadDir       *string `json:"download-dir"`
	Filename          *string `json:"filename"`
	MetaInfo          *string `json:"metainfo"`
	Paused            *bool   `json:"paused"`
	PeerLimit         *int64  `json:"peer-limit"`
	BandwidthPriority *int64  `json:"bandwidthPriority"`
	FilesWanted       []int64 `json:"files-wanted"`
	FilesUnwanted     []int64 `json:"files-unwanted"`
	PriorityHigh      []int64 `json:"priority-high"`
	PriorityLow       []int64 `json:"priority-low"`
	PriorityNormal    []int64 `json:"priority-normal"`
}
