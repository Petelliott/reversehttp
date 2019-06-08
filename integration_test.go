package reversehttp

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReverseHTTPGet(t *testing.T) {
	endserver := make(chan struct{})

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := ReverseRequest(w, r)
		expect(t, nil, err)
		resp, err := c.Get("http://example.com/path2")
		expect(t, nil, err)

		b, err := ioutil.ReadAll(resp.Body)
		expect(t, nil, err)
		expect(t, []byte("hello world\n"), b)

		close(endserver)
	}))
	defer srv.Close()

	http.DefaultClient = srv.Client()

	err := ReverseFunc(srv.URL, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plan")
		_, err := w.Write([]byte("hello world\n"))
		expect(t, nil, err)
	})
	expect(t, nil, err)
	<-endserver
}

func TestReverseReverseHTTP(t *testing.T) {
	endserver := make(chan struct{})

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

		close(endserver)
	}))
	defer srv.Close()

	http.DefaultClient = srv.Client()

	endclient := make(chan struct{})
	err := ReverseFunc(srv.URL, func(w http.ResponseWriter, r *http.Request) {
		c, err := ReverseRequest(w, r)
		expect(t, nil, err)
		resp, err := c.Get("http://example.com/path2")
		expect(t, nil, err)

		b, err := ioutil.ReadAll(resp.Body)
		expect(t, nil, err)
		expect(t, []byte("hello world\n"), b)
		resp.Body.Close()
		close(endclient)
	})
	expect(t, nil, err)
	<-endclient
	<-endserver
}
