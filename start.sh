#!/bin/bash
echo "Starting MyGola with Watchexec watcher..."
watchexec -r -e go,html,tpl,yaml --shell bash "pkill -f 'go run main.go serve'; go run main.go serve"
