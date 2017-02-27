#!/bin/bash
# unison.sh
# synchronize a local path to a vagrant host (the default host)
#
# Usage: unison.sh <PATH>

unison \
    -auto \
    -batch \
    -times \
    -ignore "Path \.vagrant" \
    -ignore "Path */node_modules" \
    -ignore "Path bazel-*" \
    -ignore "Path vendor" \
    -ignore "Name *~" \
    -ignore "Name .*.swp" \
    -ignore "Name .*~" \
    -ignore "Name ._*" \
    -ignore "Name .DS_Store" \
    -force newer \
    -terse \
    -repeat 1 \
    -confirmbigdel \
    -ui "text" \
    -prefer ./ \
    -sshargs '-i .vagrant/machines/default/virtualbox/private_key -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no' \
    ./ \
    ssh://vagrant@192.168.33.11//home/vagrant/golang/src/v2.staffjoy.com

