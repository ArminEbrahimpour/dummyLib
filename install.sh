#!/bin/bash


sudo apt install libayatana-appindicator-glib-dev

go mod download

go build -o dummylib ./cmd

ln -s /usr/bin/dummylib dummylib



