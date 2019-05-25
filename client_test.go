package reversehttp

import (
	"testing"
	"fmt"
	"net/http"
	"reflect"
	"bytes"
	"io/ioutil"
)

func expect(t *testing.T, expected interface{}, got interface{}) bool {
	if !reflect.DeepEqual(expected, got) {
		fmt.Printf("expected: %v, got: %v\n", expected, got)
		t.Error()
		return false
	}
	return true
}

func TestNewRequest(t *testing.T) {
	req, err := NewRequest("http://example.com/path")
	if err != nil {
		t.Error()
	} else {
		if !IsReverseHTTPRequest(req) {
			fmt.Printf("req '%v' is not a valid reverse http request\n", req)
			t.Error()
		}
	}

	req, err = NewRequest("asdkjfklvqnvnon  idga %%2")
	if err == nil {
		t.Error()
	}
}

func TestIsReverseHTTPResponse(t *testing.T) {
	h := http.Header{}
	h.Add("upgrade", "PTTH/1.0")
	h.Add("CONNECTION", "Upgrade")
	expect(t, true, IsReverseHTTPResponse(&http.Response{
		StatusCode: 101,
		Header: h,
	}))

	expect(t, false, IsReverseHTTPResponse(&http.Response{
		StatusCode: 200,
		Header: h,
	}))

	h.Del("upgrade")
	expect(t, false, IsReverseHTTPResponse(&http.Response{
		StatusCode: 101,
		Header: h,
	}))

	expect(t, false, IsReverseHTTPResponse(nil))
}

func TestInternalResponse(t *testing.T) {
	req, err := NewRequest("http://example.com/path")
	expect(t, nil, err)

	buf := bytes.NewBuffer(make([]byte, 0))

	resp := newResponse(req, buf)
	resp.Header().Add("Content-Type", "application/x-testtype")

	resp.Write([]byte("hello world\n"))

	b, err := ioutil.ReadAll(buf)
	expect(t, nil, err)
	expected := []byte("HTTP/1.1 200 OK\r\nContent-Length: 0\r\nContent-Type: application/x-testtype\r\n\r\nhello world\n")
	expect(t, expected, b)

	// test with writeheader
	buf = bytes.NewBuffer(make([]byte, 0))

	resp = newResponse(req, buf)
	resp.Header().Add("Content-Type", "application/x-testtype")

	resp.WriteHeader(http.StatusOK)
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte("hello world\n"))

	b, err = ioutil.ReadAll(buf)
	expect(t, nil, err)
	expect(t, expected, b)

}

