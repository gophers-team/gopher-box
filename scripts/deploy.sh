#!/bin/bash

mkdir -p /var/lib/gopher-box
touch /var/lib/gopher-box/db

cp /tmp/gopher-box /srv/gopher-box
kill $(pgrep gopher-box) &> /dev/null
/srv/gopher-box
