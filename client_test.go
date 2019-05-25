package reversehttp

import (
	"testing"
	"fmt"
)

func expect(t *testing.T, expected interface{}, got interface{}) {
	if expected != got {
		fmt.Printf("expected: %v, got: %v", expected, got)
		t.Error()
	}
}

func TestNewRequest(t *testing.T) {
	req, err := NewRequest("http://example.com/path")
	if err != nil {
		t.Error()
	} else {
		if !IsReverseHTTPRequest(req) {
			fmt.Printf("req '%v' is not a valid reverse http request", req)
			t.Error()
		}
	}

	req, err = NewRequest("asdkjfklvqnvnon  idga %%2")
	if err == nil {
		t.Error()
	}
}
