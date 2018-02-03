package TransmissionRPC

import (
	"errors"
	"net/http"
	"time"
)

// Controller is the base object to interract with a remote transmission rpc endpoint
// It must be created with New()
type Controller struct {
	host     string
	user     string
	password string
	rpcURI   string
	httpC    *http.Client
}

// AdvancedConfig handles options that are not mandatory for New()
type AdvancedConfig struct {
	RPCURI      string
	HTTPTimeout time.Duration
}

// New returns an initialized and ready to use Controller
func New(host, user, password string, conf *AdvancedConfig) (c *Controller, err error) {
	// Check extra config if any
	if conf == nil {
		conf = &AdvancedConfig{
			RPCURI:      "/transmission/rpc",
			HTTPTimeout: 30 * time.Second,
		}
	}
	if conf.HTTPTimeout == 0 {
		err = errors.New("HTTPTimeout can't be 0")
		return
	}
	// Initialize & return
	c = &Controller{
		host:     host,
		user:     user,
		password: password,
		rpcURI:   conf.RPCURI,
		httpC: &http.Client{
			Timeout: time.Minute,
		},
	}
	return
}
