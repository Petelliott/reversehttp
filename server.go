package reversehttp

import (
	"net/http"
	"errors"
	"bufio"
	"sync"
)

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

	err := req.Write(it.rw)
	if err != nil {
		return nil, err
	}
	it.rw.Flush()
	resp, err := http.ReadResponse(it.rw.Reader, req)
	return resp, err
}

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

	return &http.Client {
		Transport: newIoTripper(buf),
	}, nil
}
