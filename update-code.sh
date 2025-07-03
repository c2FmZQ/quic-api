#!/bin/bash
# This scripts updates the go files in the current directory to use the API from
# this package. It doesn't add the required import.

for f in $(find . -name "*.go"); do
  sed -i -r \
    -e 's:[*]quic[.]((Conn|(Send|Receive)?Stream|Transport)([,;\) ]|$)):quicapi.\1:' \
    -e 's:quic[.]((Dial|Listen)(Addr)?(Early)?):quicapi.\1:' $f
done
