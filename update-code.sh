#!/bin/bash -e
# This scripts updates the go files in the current directory to use the API from
# this package.

for f in $(find . -name "*.go"); do
  sed -i -r \
    -e 's:[*]quic[.]((Conn|(Send|Receive)?Stream|Transport|EarlyListener)([,;\) ]|$)):quicapi.\1:' \
    -e 's:quic[.]((Dial|Listen)(Addr)?(Early)?):quicapi.\1:' \
    -e 's:&quic[.]((Early)?Listener)\{\}:(quicapi.\1)(nil):' \
    -e 's:(&quic[.]Transport\{[^}]*\}):quicapi.WrapTransport(\1):' \
    $f
done

for f in $(git grep -l quicapi); do
  awk '/^import/ {
        imp=1;
      }
      /^)?$/ {
        if(imp) {
          printf("\n\tquicapi \"github.com/c2FmZQ/quic-api\"\n")};
          imp=0;
        }
      !/github.com\/c2FmZQ\/quic-api/ {
        print;
      }
  ' < $f > tmpfile
  mv tmpfile $f
  gofmt -w $f
done
