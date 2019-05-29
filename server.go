package reversehttp

import (
	"bufio"
	"errors"
	"net/http"
	"sync"
)

// IsReverseHTTPRequest returns true if response is a valid Reverse HTTP
// upgrade Request (i.e. a valid HTTP/1.1 protocol upgrade request where the
// Upgrade Header is "PTTH/1.0).  This function will return False otherwise,
// including when req is nil.
func IsReverseHTTPRequest(req *http.Request) bool {
	if req == nil {
		return false
	}

	return req.Header.Get("Upgrade") == "PTTH/1.0" &&
		req.Header.Get("Connection") == "Upgrade"
}

type ioTripper struct {
	mu sync.Mutex
	rw *bufio.ReadWriter
}

func newIoTripper(rw *bufio.ReadWriter) *ioTripper {
	return &ioTripper{
		rw: rw,
	}
}

func (it *ioTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	it.mu.Lock()
	defer it.mu.Unlock()

	// write will usually not error, if it does flush will also error
	req.Write(it.rw)
	err := it.rw.Flush()
	if err != nil {
		return nil, err
	}
	resp, err := http.ReadResponse(it.rw.Reader, req)
	return resp, err
}

// ReverseRequest produces an http.Client from an http.ResponseWriter and http.Request.
// This Client can be used for a single request. It is possible that the
// connection is kept alive, and the client can be used more than once but this
// behavior shouldn't be relied on.
func ReverseRequest(w http.ResponseWriter, r *http.Request) (*http.Client, error) {
	if !IsReverseHTTPRequest(r) {
		return nil, errors.New("request is not a valid reverse http request")
	}
	w.Header().Add("Upgrade", "PTTH/1.0")
	w.Header().Add("Connection", "Upgrade")
	w.WriteHeader(http.StatusSwitchingProtocols)

	_, buf, err := w.(http.Hijacker).Hijack()
	if err != nil {
		return nil, err
	}

	return &http.Client{
		Transport: newIoTripper(buf),
	}, nil
}
