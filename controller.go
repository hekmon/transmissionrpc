package TransmissionRPC

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

// Controller is the base object to interract with a remote transmission rpc endpoint
// It must be created with New()
type Controller struct {
	url       string
	user      string
	password  string
	sessionID string
	rnd       *rand.Rand
	httpC     *http.Client
}

// AdvancedConfig handles options that are not mandatory for New()
// If Port is not specified it will be set at 80 if HTTPS is false, 443 otherwise
type AdvancedConfig struct {
	HTTPS       bool
	Port        uint16
	RPCURI      string
	HTTPTimeout time.Duration
}

// New returns an initialized and ready to use Controller
func New(host, user, password string, conf *AdvancedConfig) (c *Controller, err error) {
	// Check extra config if any
	if conf == nil {
		conf = &AdvancedConfig{
			// HTTPS false by default
			RPCURI:      "/transmission/rpc",
			HTTPTimeout: 30 * time.Second,
		}
	}
	if conf.HTTPTimeout == 0 {
		err = errors.New("HTTPTimeout can't be 0")
		return
	}
	// Compute missing data
	var scheme string
	if conf.HTTPS {
		scheme = "https"
		if conf.Port == 0 {
			conf.Port = 443
		}
	} else {
		scheme = "http"
		if conf.Port == 0 {
			conf.Port = 80
		}
	}
	// Initialize & return
	c = &Controller{
		url:      fmt.Sprintf("%s://%s:%d%s", scheme, host, conf.Port, conf.RPCURI),
		user:     user,
		password: password,
		rnd:      rand.New(rand.NewSource(time.Now().Unix())),
		httpC: &http.Client{
			Timeout: time.Minute,
		},
	}
	return
}
