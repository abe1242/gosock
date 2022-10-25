#!/usr/bin/env bash

set -e

if [ $OSTYPE == "linux-gnu" ]; then
    file="gso_linux_x86_64.tar.gz"
    set -x
    if ! curl -sL "https://github.com/abe1242/gosock/releases/latest/download/$file" -o $file; then
        wget -q "https://github.com/abe1242/gosock/releases/latest/download/$file"
    fi
    
    tar xzf $file
    rm $file
    sudo mkdir -p /usr/local/bin
    sudo mv gso /usr/local/bin/
elif [ $OSTYPE == "linux-android" ]; then
    file="gso_linux_arm64.tar.gz"
    set -x
    curl -sL "https://github.com/abe1242/gosock/releases/latest/download/$file" -o $file
    tar xzf $file
    rm $file

    mkdir -p ~/.local/bin
    set +x
    if ! echo $PATH | grep -q $HOME/.local/bin; then
        set -x
        echo 'PATH=$PATH:$HOME/.local/bin' >> ~/.profile
        source ~/.profile
    fi
    set -x

    mkdir -p ~/.local/bin
    mv gso ~/.local/bin
fi
