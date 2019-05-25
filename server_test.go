package reversehttp

import (
	"testing"
	"net/http"
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
