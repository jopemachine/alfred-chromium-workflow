#!/bin/bash
rm -rf *alfred-chromium-workflow*
go build .
dir_name=$(basename "$PWD")
zip -r "${dir_name}.zip" *
mv "${dir_name}.zip" "${dir_name}.alfredworkflow"
