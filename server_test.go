package reversehttp

import (
	"bufio"
	"bytes"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"reflect"
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

type errorWriter struct{}

func (ew errorWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("error writers always fail, this is expected")
}

func (ew errorWriter) Read(p []byte) (n int, err error) {
	return 0, errors.New("error writers always fail, this is expected")
}

func TestUpgradeBodyRead(t *testing.T) {
	var buf bytes.Buffer
	buf.Write([]byte("hello world\n"))
	ub := newUpgradeBody(bufio.NewReadWriter(bufio.NewReader(&buf), nil), nil)

	b, err := ioutil.ReadAll(ub)
	expect(t, nil, err)
	expect(t, "hello world\n", string(b))
}

func TestUpgradeBodyWrite(t *testing.T) {
	var buf bytes.Buffer
	ub := newUpgradeBody(bufio.NewReadWriter(nil, bufio.NewWriter(&buf)), nil)

	n, err := ub.Write([]byte("hello world\n"))
	expect(t, 12, n)
	expect(t, nil, err)

	// verify that the write is not buffered
	b, err := ioutil.ReadAll(&buf)
	expect(t, nil, err)
	expect(t, "hello world\n", string(b))

	ub = newUpgradeBody(bufio.NewReadWriter(nil, bufio.NewWriter(errorWriter{})), nil)
	_, err = ub.Write([]byte("hello world\n"))
	if err == nil {
		t.Error("write did not fail")
	}

	ub = newUpgradeBody(bufio.NewReadWriter(nil, bufio.NewWriterSize(errorWriter{}, 1)), nil)
	_, err = ub.Write([]byte("hello world\n"))
	if err == nil {
		t.Error("write did not fail")
	}
}

type closeChecker struct {
	hasClosed bool
}

func (cc *closeChecker) Close() error {
	cc.hasClosed = true
	return nil
}

func TestUpgradeBodyClose(t *testing.T) {
	cc := closeChecker{false}
	ub := newUpgradeBody(nil, &cc)

	err := ub.Close()
	expect(t, nil, err)
	expect(t, true, cc.hasClosed)
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

	// test 101 switching protocols case
	rbuf := new(bytes.Buffer)
	wbuf := new(bytes.Buffer)
	rbuf.Write([]byte("HTTP/1.1 101 Switching Protocols\r\nUpgrade: test\r\nConnection: Upgrade\r\n\r\nhello world\n"))
	it = newIoTripper(bufio.NewReadWriter(bufio.NewReader(rbuf), bufio.NewWriter(wbuf)))

	resp, err := it.RoundTrip(r)
	expect(t, nil, err)
	expect(t, 101, resp.StatusCode)
	expect(t, "reversehttp.upgradeBody", reflect.TypeOf(resp.Body).String())
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
	defer srv.Close()

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
