#!/bin/bash

set -e

if [ -z "$1" ]; then
    echo "required patch/minor/major" 1>&2
    exit 1
fi

# gofmt
gofmt -s -w .
git add ./*.go
git commit -m "gofmt"

# gobump
new_version=$(gobump "$1" -w -v cmd | jq -r '.[]')
git add ./*.go
git commit -m "Bump version $new_version"
git push origin master

# build release files
goxc -tasks='xc archive' -bc 'linux,!arm windows darwin' -d .
cp -p $(PWD)/snapshot/linux_amd64/grabeni $(PWD)/snapshot/grabeni_linux_amd64
cp -p $(PWD)/snapshot/linux_386/grabeni $(PWD)/snapshot/grabeni_linux_386
cp -p $(PWD)/snapshot/darwin_amd64/grabeni $(PWD)/snapshot/grabeni_darwin_amd64
cp -p $(PWD)/snapshot/darwin_386/grabeni $(PWD)/snapshot/grabeni_darwin_386
cp -p $(PWD)/snapshot/windows_amd64/grabeni.exe $(PWD)/snapshot/grabeni_windows_amd64.exe
cp -p $(PWD)/snapshot/windows_386/grabeni.exe $(PWD)/snapshot/grabeni_windows_386.exe

# make github release draft
ghr --username yuuki1 --replace --draft "v$new_version" snapshot/
