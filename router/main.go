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

func (h *RequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.RequestURI = ""
	r.URL.Scheme = "http"

	targetRouter := getRouterByHost(r.Host)
	if targetRouter == nil {
		w.WriteHeader(http.StatusOK)

		compiledHello := CompileHelloPage()
		fmt.Fprint(w, compiledHello)

		return
	}

	log.Println("launched: ", targetRouter.Running)
	targetRouter.Launch()

	log.Printf("real host: %v, target: %v\n", r.Host, targetRouter.Source)
	r.URL.Host = targetRouter.Source

	resp, err := h.Transport.RoundTrip(r)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)

		compiledError := CompileErrorPage(err)
		fmt.Fprint(w, compiledError)

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

func getRouterByHost(host string) *Router {
	for i := 0; i < len(routers); i++ {
		if routers[i].Alias == host {
			return &routers[i]
		}
	}
	return nil
}

type DebugHandler struct {
}

func (dh *DebugHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	targetRouter := getRouterByHost(r.Host)

	fmt.Fprintf(w, CompileDebugPage(targetRouter))
	// fmt.Fprintf(w, "<html><head></head><body>Debug page. Running: %t</body></html>", targetRouter.Running)
}

func main() {
	// TODO: need `flags` here — @kavu
	routers = PopulateRouters("/Users/kir/.router/config.toml")

	mux := http.NewServeMux()

	var request_handler http.Handler = &RequestHandler{Transport: &http.Transport{DisableKeepAlives: false, DisableCompression: false}}

	var debug_handler http.Handler = &DebugHandler{}

	mux.Handle("/", request_handler)
	mux.Handle("/debug", debug_handler)

	bind_address := "0.0.0.0:9000"
	log.Printf("stating router on %s", bind_address)
	srv := &http.Server{Handler: mux, Addr: bind_address}

	if err := srv.ListenAndServe(); err != nil {
		log.Printf("Starting frontend failed: %v", err)
	}
}
