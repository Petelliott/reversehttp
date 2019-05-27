// reversehttp allows you to make and respond to reverse http requests
//
// Reverse Http
//
// Reverse http is a protocol for http servers to make requests to http clients it
// works by the client connecting to the server, reqesting to upgrade to reverse
// http, the server upgrades the connection, and then the connection becomes like
// a TCP connection for the server to make a request to the client on.
//
// Example
//
// Server:
//
//	http.HandleFunc("/ptth", func(w http.ResponseWriter, r *http.Request) {
//		client, _ := ReverseRequest(w, r)
//		resp, _ := c.Get("http://example.com/path")
//		// do whatever you want with the response
//	})
//
// Client:
//
//	err := ReverseFunc("http://example.com/ptth", func(w http.ResponseWriter, r *http.Request) {
//		// this could be any http.Handler
//		w.Header().Add("Content-Type", "text/plain")
//		w.Write([]byte("hello world\n"))
//	})
package reversehttp
