package transmissionrpc

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
)

const (
	// RPCVersion indicates the exact transmission RPC version this library is build against
	RPCVersion       = 15
	defaultPort      = 9091
	defaultRPCPath   = "/transmission/rpc"
	defaultTimeout   = 30 * time.Second
	defaultUserAgent = "github.com/hekmon/transmissionrpc"
)

// Client is the base object to interract with a remote transmission rpc endpoint.
// It must be created with New().
type Client struct {
	url             string
	user            string
	password        string
	sessionID       string
	sessionIDAccess sync.RWMutex
	userAgent       string
	rnd             *rand.Rand
	httpC           *http.Client
	debug           bool
}

// AdvancedConfig handles options that are not mandatory for New().
// Default value for HTTPS is false, default port is 9091, default RPC URI is
// '/transmission/rpc', default HTTPTimeout is 30s.
type AdvancedConfig struct {
	HTTPS       bool
	Port        uint16
	RPCURI      string
	HTTPTimeout time.Duration
	UserAgent   string
	Debug       bool
}

// Config is the input data needed to make a connection to Transmission RPC.
type Config struct {
	// URL should be a full url to the /transmission/rpc endpoint.
	URL string
	// Username if authentication is required.
	Username string
	// Password if authentication is required.
	Password string
	// UserAgent is set to this package's url if not provided here.
	UserAgent string
	// Client is set to an empty http.Client{} with a sane timeout if not provided.
	Client *http.Client
}

// NewClient is the new entry into this module. Provide a Config with at minimum a URL.
// The Client returned may be used to query the Transmission RPC.
func NewClient(config Config) *Client {
	if config.Client == nil {
		config.Client = &http.Client{Timeout: defaultTimeout}
	}

	if config.UserAgent == "" {
		config.UserAgent = defaultUserAgent
	}

	return &Client{
		url:       config.URL,
		user:      config.Username,
		password:  config.Password,
		userAgent: config.UserAgent,
		rnd:       rand.New(newLockedRandomSource(time.Now().Unix())),
		httpC:     config.Client,
		// If you need debug, use a custom transport/roundtripper in your http.Client.
		// Ex: https://pkg.go.dev/golift.io/starr@main/debuglog
		debug: false,
	}
}

// New returns an initialized and ready to use Controller
func New(host, user, password string, conf *AdvancedConfig) (c *Client, err error) {
	// Config
	if conf != nil {
		// Check custom config
		if conf.Port == 0 {
			conf.Port = defaultPort
		}
		if conf.RPCURI == "" {
			conf.RPCURI = defaultRPCPath
		}
		if conf.HTTPTimeout == 0 {
			conf.HTTPTimeout = defaultTimeout
		}
		if conf.UserAgent == "" {
			conf.UserAgent = defaultUserAgent
		}
	} else {
		// Spawn default config
		conf = &AdvancedConfig{
			// HTTPS false by default
			Port:        defaultPort,
			RPCURI:      defaultRPCPath,
			HTTPTimeout: defaultTimeout,
			UserAgent:   defaultUserAgent,
		}
	}
	// Build & validate URL
	var scheme string
	if conf.HTTPS {
		scheme = "https"
	} else {
		scheme = "http"
	}
	remoteURL, err := url.Parse(fmt.Sprintf("%s://%s:%d%s", scheme, host, conf.Port, conf.RPCURI))
	if err != nil {
		err = fmt.Errorf("can't build a valid URL: %w", err)
		return
	}
	// Initialize & return ready to use client
	c = &Client{
		url:       remoteURL.String(),
		user:      user,
		password:  password,
		userAgent: conf.UserAgent,
		rnd:       rand.New(newLockedRandomSource(time.Now().Unix())),
		httpC:     cleanhttp.DefaultPooledClient(),
		debug:     conf.Debug,
	}
	c.httpC.Timeout = conf.HTTPTimeout
	return
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
