#!/bin/bash
set -e

if ! command -V vagrant >/dev/null 2>&1; then
    echo "OOPS! Run this command OUTSIDE of vagrant so that we can watch the file system better."
    exit 1
fi

if ! command -V modd >/dev/null 2>&1; then
    echo "Oops - please install `modd` on your host machine"
    echo "https://github.com/cortesi/modd"
    exit 1
fi

# Boot the vm
vagrant up

# Clear unison caches. This is surprisingly important after rebuilds and such.
rm -rf ~/Library/Application\ Support/Unison/
vagrant ssh -c "rm -rf /home/vagrant/.unison"

# Catch shutdown signal and kill both
trap 'kill %1;' SIGINT

# Run both in parallel
./vagrant/unison.sh | sed -e 's/^/[rsync] /' & (sleep 5; echo "booting modd"; modd) | sed -e 's/^/[modd] /'


# shut down the VM
vagrant halt
