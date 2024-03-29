#!/usr/bin/env bash
set -xeuo pipefail

go build -o go-telnet.goc

(echo -e "Hello\nFrom\nNC\n" && cat 2>/dev/null) | nc -ls localhost -p 4242 > ./nc.out &
NC_PID=$!

sleep 1

(echo -e "I\nam\nTELNET client\n" && cat 2>/dev/null) | ./go-telnet.goc localhost 4242 > ./telnet.out &
TL_PID=$!

sleep 5
kill ${TL_PID} 2>/dev/null || true
kill ${NC_PID} 2>/dev/null || true

function fileEquals() {
  local fileData
  fileData=$(cat "$1")
  [ "${fileData}" = "${2}" ] || (echo -e "FAIL: unexpected output, $1:\n${fileData}" && exit 1)
}

expected_nc_out='I
am
TELNET client'
fileEquals ./nc.out "${expected_nc_out}"

expected_telnet_out='...Try connect to localhost:4242
...Connected to localhost:4242
Hello
From
NC

EOF
...Force disconnect from localhost:4242'
fileEquals ./telnet.out "${expected_telnet_out}"

rm -f go-telnet.goc
echo "PASS"
