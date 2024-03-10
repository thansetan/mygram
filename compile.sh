#!/bin/bash
file="."
name=$1
output_dir="./bin"
mkdir -p $output_dir

# Compile for Windows
GOOS=windows GOARCH=amd64 go build -o ${output_dir}/${name}_windows.exe ${file}
# Compile for Linux
GOOS=linux GOARCH=amd64 go build -o ${output_dir}/${name}_linux ${file}
# Compile for Darwin
GOOS=darwin GOARCH=amd64 go build -o ${output_dir}/${name}_darwin ${file}

echo "Done!"