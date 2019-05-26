package reversehttp

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"bytes"
	"io/ioutil"
	"errors"
	"bufio"
	"net"
)

func TestIsReverseHTTPRequest(t *testing.T) {
	h := http.Header{}
	h.Add("upgrade", "PTTH/1.0")
	h.Add("CONNECTION", "Upgrade")
	expect(t, true, IsReverseHTTPRequest(&http.Request{
		Header: h,
	}))


	h.Del("upgrade")
	expect(t, false, IsReverseHTTPRequest(&http.Request{
		Header: h,
	}))

	expect(t, false, IsReverseHTTPRequest(nil))
}

type errorWriter struct {}

func (ew errorWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("error writers always fail, this is expected")
}

func (ew errorWriter) Read(p []byte) (n int, err error) {
	return 0, errors.New("error writers always fail, this is expected")
}


func TestIoTripper(t *testing.T) {
	it := newIoTripper(bufio.NewReadWriter(bufio.NewReader(errorWriter{}), bufio.NewWriter(errorWriter{})))

	r, err := http.NewRequest("GET", "http://example.com/path", nil)
	if err != nil {
		t.Error()
	}

	_, err = it.RoundTrip(r)
	if err == nil {
		t.Error()
	}
}

type ResponseHijackFailer struct {
	http.ResponseWriter
}

func (rhf *ResponseHijackFailer) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return nil, nil, errors.New("ResponseHijackFailer will always fail to hijack")
}

func TestReverseRequest(t *testing.T) {
	w := httptest.NewRecorder()

	r, err := http.NewRequest("GET", "http://example.com/path", nil)
	if err != nil {
		t.Error()
	}

	_, err = ReverseRequest(w, r)
	if err == nil {
		t.Error()
	}

	r, err = NewRequest("http://example.com/path")
	expect(t, nil, err)
	_, err = ReverseRequest(&ResponseHijackFailer{w}, r)
	if err == nil {
		t.Error()
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := ReverseRequest(w, r)
		expect(t, nil, err)
		resp, err := c.Get("http://example.com/path2")
		expect(t, nil, err)

		expect(t, "text/plain", resp.Header.Get("Content-Type"))

		b, err := ioutil.ReadAll(resp.Body)
		expect(t, nil, err)
		expect(t, []byte("hello world\n"), b)
	}))

	r, err = NewRequest(srv.URL)
	r.Body = ioutil.NopCloser(bytes.NewReader([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nhello world\n")))
	expect(t, nil, err)

	client := srv.Client()
	resp, err := client.Do(r)
	expect(t, nil, err)

	innerreq, err := http.ReadRequest(bufio.NewReader(resp.Body))
	expect(t, nil, err)

	expect(t, "GET", innerreq.Method)
	expect(t, "/path2", innerreq.URL.String())
}
