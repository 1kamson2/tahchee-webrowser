#!/bin/bash

readonly AGENT="agent.json"
readonly DOM="dom.json"
readonly RESOURCES="/home/kums0nd/Dev/scrapper/go/resources/"

mkdir -p "$RESOURCES$1"
go run main.go $DOM $AGENT $RESOURCES
