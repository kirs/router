package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	. "router"
)

var routers []Router

type RequestHandler struct {
	Transport *http.Transport
}

type DebugHandler struct {
}

func (dh *DebugHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<html><head></head><body>Debug page here.</body></html>")
}

func (h *RequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.RequestURI = ""
	r.URL.Scheme = "http"

	var localRequest = r.Host == "localhost"

	if localRequest {
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, "<html><head></head><body>Hello page here, %username%.</body></html>")
		return
	}

	r.URL.Host = r.Host + ":9000"

	resp, err := h.Transport.RoundTrip(r)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)

		// compiledError := CompileErrorPage(err)
		// fmt.Fprint(w, compiledError)
		fmt.Fprint(w, "<html><head></head><body>Load failed</body></html>")

		return
	}

	for k, v := range resp.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}

	w.WriteHeader(resp.StatusCode)

	io.Copy(w, resp.Body)
	resp.Body.Close()
}

func main() {
	// TODO: need `flags` here — @kavu
	routers = PopulateRouters("/Users/kir/.router/config.toml")

	PopulateEtcHosts()

	var request_handler http.Handler = &RequestHandler{Transport: &http.Transport{DisableKeepAlives: false, DisableCompression: false}}

	bind_address := "0.0.0.0:80"
	log.Printf("stating router on %s", bind_address)
	srv := &http.Server{Handler: request_handler, Addr: bind_address}

	if err := srv.ListenAndServe(); err != nil {
		log.Printf("Starting frontend failed: %v", err)
	}
}
