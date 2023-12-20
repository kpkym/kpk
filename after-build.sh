#!/bin/bash

os_arch="darwin_amd64"
bin_dir="$HOME/config/bin"

if [[ $1 == *$os_arch* ]]
then
    # Copy the file to the destination directory
    cp $1 $bin_dir
    echo "Copied: $1"
fi