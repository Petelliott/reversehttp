package reversehttp

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"bytes"
	"io/ioutil"
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
	r.Body = ioutil.NopCloser(bytes.NewReader([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nhello world\n")))
	if err != nil {
		t.Error()
	}

	c, err := ReverseRequest(w, r)
	if err != nil {
		t.Error()
	} else {
		resp, err := c.Get("http://example.com/path2")
		if expect(t, nil, err) {
			expect(t, "text/plain", resp.Header.Get("Content-Type"))

			b, err := ioutil.ReadAll(resp.Body)
			expect(t, nil, err)
			expect(t, []byte("hello world\n"), b)
		}
	}
}
