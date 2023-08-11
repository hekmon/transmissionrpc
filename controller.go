package transmissionrpc

import (
	"errors"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
)

const (
	// RPCVersion indicates the exact transmission RPC version this library is build against
	RPCVersion       = 17
	defaultUserAgent = "github.com/hekmon/transmissionrpc"
)

// Config is the input data needed to make a connection to Transmission RPC.
type Config struct {
	// UserAgent is set to this package's url if not provided here.
	UserAgent string
	// Client is set to a clean and isolated client if not provided
	CustomClient *http.Client
}

// New returns an initialized and ready to use Controller
func New(transmissionRPCendpoint *url.URL, extra *Config) (c *Client, err error) {
	// handle user input
	if transmissionRPCendpoint == nil {
		err = errors.New("please provide an Transmission RPC endpoint URL")
		return
	}
	if extra != nil {
		if extra.UserAgent == "" {
			extra.UserAgent = defaultUserAgent
		}
		if extra.CustomClient == nil {
			extra.CustomClient = cleanhttp.DefaultPooledClient()
		}
	} else {
		extra = &Config{
			UserAgent:    defaultUserAgent,
			CustomClient: cleanhttp.DefaultPooledClient(),
		}
	}
	// Initialize & return ready to use client
	c = &Client{
		endpoint:     *transmissionRPCendpoint,
		http:         extra.CustomClient,
		userAgent:    extra.UserAgent,
		tagGenerator: rand.New(newLockedRandomSource(time.Now().Unix())),
	}
	return
}

// Client is the base object to interract with a remote transmission rpc endpoint.
// It must be created with New().
type Client struct {
	// HTTP Client
	endpoint  url.URL
	http      *http.Client
	userAgent string
	// Transmission RPC protections
	tagGenerator    *rand.Rand
	sessionID       string
	sessionIDAccess sync.RWMutex
}

func (c *Client) getRandomTag() int {
	return c.tagGenerator.Int()
}

func (c *Client) getSessionID() string {
	defer c.sessionIDAccess.RUnlock()
	c.sessionIDAccess.RLock()
	return c.sessionID
}

func (c *Client) updateSessionID(newID string) {
	defer c.sessionIDAccess.Unlock()
	c.sessionIDAccess.Lock()
	c.sessionID = newID
}

// rand.NewSource is not thread-safe, so access should be serialized
type lockedRandomSource struct {
	mut sync.Mutex
	src rand.Source
}

func newLockedRandomSource(seed int64) rand.Source {
	return &lockedRandomSource{
		src: rand.NewSource(seed),
	}
}

func (lrs *lockedRandomSource) Int63() (rnd int64) {
	lrs.mut.Lock()
	rnd = lrs.src.Int63()
	lrs.mut.Unlock()
	return
}

func (lrs *lockedRandomSource) Seed(seed int64) {
	lrs.mut.Lock()
	lrs.src.Seed(seed)
	lrs.mut.Unlock()
}
