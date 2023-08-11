package transmissionrpc

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	t.Parallel()

	expect := &Client{
		httpC:     &http.Client{Timeout: defaultTimeout},
		user:      "localuser",
		password:  "localpass",
		url:       fmt.Sprint("http://localhost:", defaultPort, defaultRPCPath),
		userAgent: defaultUserAgent,
	}

	client, err := New("localhost", expect.user, expect.password, nil)
	if err != nil {
		t.Fatalf("Requesting a simple rpc client must not produce an error: %v", err)
	}

	testCheckClient(t, client, expect)
}

func TestNewClient(t *testing.T) {
	t.Parallel()

	expect := &Client{
		httpC:     &http.Client{Timeout: defaultTimeout},
		user:      "localuser",
		password:  "localpass",
		url:       "http://localhost:999/rpc",
		debug:     false,
		userAgent: "test agent",
	}
	config := Config{
		URL:       expect.url,
		Username:  expect.user,
		Password:  expect.password,
		UserAgent: expect.userAgent,
	}

	testCheckClient(t, NewClient(config), expect)

	expect.httpC.Timeout = time.Hour
	config.Client = &http.Client{Timeout: time.Hour}
	testCheckClient(t, NewClient(config), expect)
}

func testCheckClient(t *testing.T, received *Client, expected *Client) {
	t.Helper()

	switch {
	case received == nil:
		t.Fatal("Requesting a rpc client must not produce a nil client.")
	case received.httpC == nil:
		t.Fatal("Requesting a rpc client must not produce a nil http client.")
	}

	if received.password != expected.password {
		t.Error("Provided client was returned with the wrong password.")
	}

	if received.user != expected.user {
		t.Error("Provided client was returned with the wrong username.")
	}

	if received.url != expected.url {
		t.Error("Provided client was returned with the wrong URL.")
	}

	if received.userAgent != expected.userAgent {
		t.Error("Provided client was returned with the wrong User Agent.")
	}

	if received.httpC.Timeout != expected.httpC.Timeout {
		t.Error("Provided client was returned with the wrong http timeout.")
	}

	if received.debug != expected.debug {
		t.Error("Provided client was returned with the wrong debug setting.")
	}
}
