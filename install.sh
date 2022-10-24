#!/usr/bin/env bash

set -e

if [ $OSTYPE == "linux-gnu" ]; then
    file="gso_linux_x86_64.tar.gz"
    curl -L "https://github.com/abe1242/gosock/releases/latest/download/$file" -o $file
    tar xzf $file
    rm $file
    sudo mkdir -p /usr/local/bin
    sudo mv gso /usr/local/bin/
elif [ $OSTYPE == "linux-android" ]; then
    file="gso_linux_arm64.tar.gz"
    curl -L "https://github.com/abe1242/gosock/releases/latest/download/$file" -o $file
    tar xzf $file
    rm $file

    mkdir -p ~/.local/bin
    if ! echo $PATH | grep -q $HOME/.local/bin; then
        echo 'PATH=$PATH:$HOME/.local/bin' >> ~/.profile
        source ~/.profile
    fi

    mkdir -p ~/.local/bin
    mv gso ~/.local/bin
fi
