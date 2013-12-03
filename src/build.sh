#!/bin/bash
# sudo launchctl unload /Library/LaunchDaemons/io.router.routerd.plist
go build router.go pages.go
# sudo launchctl load /Library/LaunchDaemons/io.router.routerd.plist
