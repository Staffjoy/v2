#!/bin/bash
set -e
set -u
set -x

if ! command -V vagrant >/dev/null 2>&1; then
    echo "OOPS! Run this command OUTSIDE of vagrant so that we can watch the file system better."
    exit 1
fi

if ! command -V modd >/dev/null 2>&1; then
    echo "Oops - please install `modd` on your host machine"
    echo "https://github.com/cortesi/modd"
    exit 1
fi

modd