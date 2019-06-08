package reversehttp

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"testing"
)

type benchmarkBody struct {
	b      *testing.B
	data   []byte
	reader io.Reader
}

func newBenchmarkBody(data []byte, b *testing.B) *benchmarkBody {
	return &benchmarkBody{b, data, bytes.NewReader(data)}
}

func newBenchmarkBodyReq(request *http.Request, b *testing.B) *benchmarkBody {
	buf := new(bytes.Buffer)
	request.Write(buf)
	data, _ := ioutil.ReadAll(buf)
	return newBenchmarkBody(data, b)
}

func (bb *benchmarkBody) reset() {
	bb.reader = bytes.NewReader(bb.data)
}

func (bb *benchmarkBody) Close() error {
	return nil
}

func (bb *benchmarkBody) Write(b []byte) (int, error) {
	return len(b), nil
}

func (bb *benchmarkBody) Read(b []byte) (int, error) {
	return bb.reader.Read(b)
}

func BenchmarkReverseResponse(b *testing.B) {
	b.ReportAllocs()

	var handler http.HandlerFunc
	handler = func(w http.ResponseWriter, r *http.Request) {}

	req, _ := http.NewRequest("GET", "http://example.com/path",
		ioutil.NopCloser(bytes.NewReader([]byte("hello world\n"))))

	h := http.Header{}
	h.Add("upgrade", "PTTH/1.0")
	h.Add("CONNECTION", "Upgrade")
	r := &http.Response{
		StatusCode: http.StatusSwitchingProtocols,
		Header:     h,
		Body:       newBenchmarkBodyReq(req, b),
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ReverseResponse(r, handler)
		b.StopTimer()
		r.Body.(*benchmarkBody).reset()
		b.StartTimer()
	}
}

type benchmarkRW struct {
	header http.Header
}

func (brw benchmarkRW) Header() http.Header {
	return brw.header
}

func (brw benchmarkRW) Write(b []byte) (int, error) {
	return len(b), nil
}

func (brw benchmarkRW) WriteHeader(statusCode int) {
	// do literally nothing
}

func (brw benchmarkRW) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return nil, nil, nil
}

func BenchmarkReverseRequest(b *testing.B) {
	b.ReportAllocs()

	req, _ := NewRequest("http://example.com/path")
	brw := benchmarkRW{http.Header{}}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ReverseRequest(brw, req)
	}
}
