package main

import (
  // "fmt"
  "os/exec"
  // "log"
  // "bytes"
  // "time"
)

func main() {
  var port string
  port = "9000"
  var command = "/Users/kir/Projects/idinaidi/bin/rails"
  cmd := exec.Command("rbenv", "exec", "ruby", command, "server", "-p", port)

  cmd.Run()

}
