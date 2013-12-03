package main

import (
  "fmt"
  "strings"
  "io/ioutil"
  "github.com/BurntSushi/toml"
  "regexp"
)

func PopulateEtcHosts() {
  var result = compileHosts(readHosts())
  fmt.Println("Writing compiled /etc/hosts")

  var bytes = []byte(result)

  err := ioutil.WriteFile("/etc/hosts", bytes, 0644)
  if err != nil { panic(err) }
}

func readHosts() string {
  b, err := ioutil.ReadFile("/etc/hosts")
  if err != nil { panic(err) }

  return string(b)
}

func compileHosts(cron string) string {
  hosts := make([]string, len(routers))

  // var strings
  for i, router := range(routers) {
    hosts[i] = "127.0.0.1 " + router.Alias
  }
  var router_cron = strings.Join(hosts, "\n")

  begin_r, err := regexp.Compile(`(?m)# router start`)
  if err != nil {
      fmt.Printf("There is a problem with regexp.\n")
      return ""
  }

  end_r, err := regexp.Compile(`(?m)# router end`)
  if err != nil {
      fmt.Printf("There is a problem with regexp.\n")
      return ""
  }

  var compiled string

  if begin_r.MatchString(cron) == true && end_r.MatchString(cron) == true {
    fmt.Println("/etc/hosts already modified")

    replacer, err := regexp.Compile(`(?m)(^# router start$)[\s\w#\.-]*(^# router end)$`)
    if err != nil {
      fmt.Printf("There is a problem with your regexp.\n")
      return ""
    }

    compiled = replacer.ReplaceAllString(cron, "$1\n" + router_cron + "\n$2")
  } else {
    fmt.Println("/etc/hosts was empty")
    compiled = cron + "\n# router start\n" + router_cron + "\n# router end\n"
  }

  return compiled
}
