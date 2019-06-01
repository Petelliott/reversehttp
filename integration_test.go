package reversehttp

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReverseHTTPGet(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := ReverseRequest(w, r)
		expect(t, nil, err)
		resp, err := c.Get("http://example.com/path2")
		expect(t, nil, err)

		b, err := ioutil.ReadAll(resp.Body)
		expect(t, nil, err)
		expect(t, []byte("hello world\n"), b)
	}))
	defer srv.Close()

	http.DefaultClient = srv.Client()

	err := ReverseFunc(srv.URL, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plan")
		_, err := w.Write([]byte("hello world\n"))
		expect(t, nil, err)
	})
	expect(t, nil, err)
}

func TestReverseReverseHTTP(t *testing.T) {

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := ReverseRequest(w, r)
		expect(t, nil, err)
		http.DefaultClient = c

		err = ReverseFunc("http://whatever.co.uk/blah", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "text/plan")
			_, err := w.Write([]byte("hello world\n"))
			expect(t, nil, err)
		})
		expect(t, nil, err)
	}))
	defer srv.Close()

	http.DefaultClient = srv.Client()

	err := ReverseFunc(srv.URL, func(w http.ResponseWriter, r *http.Request) {
		c, err := ReverseRequest(w, r)
		expect(t, nil, err)
		resp, err := c.Get("http://example.com/path2")
		expect(t, nil, err)

		b, err := ioutil.ReadAll(resp.Body)
		expect(t, nil, err)
		expect(t, []byte("hello world\n"), b)
	})
	expect(t, nil, err)
}
