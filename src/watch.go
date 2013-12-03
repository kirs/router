package main

import (
    "log"
    "io/ioutil"
    "github.com/howeyc/fsnotify"
    "os/exec"
)

func PrepareCmd() *exec.Cmd {
    return exec.Command("/Users/kir/Projects/idinaidi/bin/rails", "server", "-p", "3000")
}

func main() {
    var cmd *exec.Cmd

    cmd = PrepareCmd()
    err := cmd.Start()
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Waiting for command to finish...")

    // done := make(chan error)
    go func() {
        cmd.Wait()
    }()

    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        log.Fatal(err)
    }

    done := make(chan bool)

    // Process events
    go func() {
        for {
            select {
            case ev := <-watcher.Event:
                if err := cmd.Process.Kill(); err != nil {
                    log.Printf("failed to kill: ", err)
                }

                cmd = PrepareCmd()
                err := cmd.Start()
                if err != nil {
                    log.Fatal(err)
                }

                log.Println("event:", ev)
            case err := <-watcher.Error:
                log.Println("error:", err)
            }
        }
    }()

    empty := []byte{}

    var restart_file = "/Users/kir/Projects/idinaidi/tmp/restart.txt"

    err = ioutil.WriteFile(restart_file, empty, 0644)
    if err != nil { panic(err) }

    err = watcher.Watch(restart_file)
    if err != nil {
        log.Fatal(err)
    }

    <-done

    /* ... do stuff ... */
    log.Println("doing more")
    watcher.Close()
}
