package reversehttp

import (
	"net/http"
	"errors"
	"io"
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
	writer io.Writer
	reader io.Reader
}

func newIoTripper(writer io.Writer, reader io.Reader) *ioTripper {
	return &ioTripper{
		writer: writer,
		reader: reader,
	}
}

func (it *ioTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	it.mu.Lock()
	defer it.mu.Unlock()

	err := req.Write(it.writer)
	if err != nil {
		return nil, err
	}
	return http.ReadResponse(bufio.NewReader(it.reader), req)
}

func ReverseRequest(w http.ResponseWriter, r *http.Request) (*http.Client, error) {
	if !IsReverseHTTPRequest(r) {
		return nil, errors.New("request is not a valid reverse http request")
	}
	w.Header().Add("Upgrade", "PTTH/1.0")
	w.Header().Add("Connection", "Upgrade")
	w.WriteHeader(http.StatusSwitchingProtocols)

	return &http.Client {
		Transport: newIoTripper(w, r.Body),
	}, nil
}
