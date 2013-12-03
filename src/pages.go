package main

import (
  "io/ioutil"
  "log"
  "strings"
  "path/filepath"
  "os"
  "fmt"
)

func CompileErrorPage(error_message error) string {
  var errorTemplatePath = GetProjectDir() + "/views/error.html"
  content, err := ioutil.ReadFile(errorTemplatePath)
  if err != nil {
    log.Fatal(err)
  }

  reason := fmt.Sprintf("%v", error_message)

  return strings.Replace(string(content), "{reason}", reason, 1)
}

func CompileDebugPage(router *Router) string {
  var debugTemplatePath = GetProjectDir() + "/views/debug.html"
  content, err := ioutil.ReadFile(debugTemplatePath)
  if err != nil {
    log.Fatal(err)
  }

  return string(content)
}

func CompileHelloPage() string {
  var helloTemplatePath = GetProjectDir() + "/views/hello.html"
  content, err := ioutil.ReadFile(helloTemplatePath)
  if err != nil {
    log.Fatal(err)
  }

  return string(content)
}

func GetProjectDir() string {
  dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
  if err != nil {
    log.Fatal(err)
  }

  return dir
}
