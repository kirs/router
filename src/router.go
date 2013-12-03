package main

import (
  "fmt"
  "github.com/BurntSushi/toml"
  "io"
  "io/ioutil"
  "log"
  "net/http"
  "os/exec"
  "strings"
  "regexp"
)

type tomlConfig struct {
  Routers []configRouter
}

type configRouter struct {
  Alias     string
  Source    string
  LaunchCmd string
  RailsPort string
}

type RequestHandler struct {
  Transport *http.Transport
}

type DebugHandler struct {
}

func (dh *DebugHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  targetRouter := getRouterByHost(r.Host)

  fmt.Fprintf(w, CompileDebugPage(targetRouter))
  // fmt.Fprintf(w, "<html><head></head><body>Debug page. Running: %t</body></html>", targetRouter.Running)
}

type Router struct {
  Alias     string
  Source    string
  LaunchCmd string
  RailsPort string
  Running   bool
}

func (r *Router) Launch() {
  if !r.Running {

    go func(command, port string) {
      // /usr/local/Cellar/rbenv/0.4.0/bin/
      // var launchString = "rbenv exec ruby " + command +  " server -p " + port
      // fmt.Println(launchString)
      fmt.Println("starting " + r.Alias + " process in gorutine")
      cmd := exec.Command("rbenv", "exec", "ruby", command, "server", "-p", port)

      cmd.Run()

    }(r.LaunchCmd, r.RailsPort)

    r.Running = true
  }
}

var routers []Router

func getRouterByHost(host string) *Router {
  for i := 0; i < len(routers); i++ {
    if routers[i].Alias == host {
      return &routers[i]
    }
  }
  return nil
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

func NewRouterFromConfig(cr *configRouter) Router {
  return Router{Alias: cr.Alias, Source: cr.Source, LaunchCmd: cr.LaunchCmd, Running: false, RailsPort: cr.RailsPort}
}

func ReadConfig() tomlConfig {
  var config tomlConfig
  if _, err := toml.DecodeFile("/Users/kir/.router/config.toml", &config); err != nil {
    log.Fatal(err)
  }
  return config
}

func PopulateRouters() []Router {
  config := ReadConfig()

  routers := make([]Router, len(config.Routers))

  for i, router := range(config.Routers) {
    routers[i] = NewRouterFromConfig(&router)
  }

  return routers
}

func main() {
  mux := http.NewServeMux()

  routers = PopulateRouters()

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
