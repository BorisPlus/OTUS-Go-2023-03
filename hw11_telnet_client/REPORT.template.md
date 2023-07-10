https://www.telnetbbsguide.com/

```bash
nc -vvl -s localhost -p 4242 -c '
    set -x
    read request
    echo "$request"
'
```

```bash
nc localhost 4242
```

```bash
nc -l -s localhost -p 4242 -c '
set -x
echo "Hello from Server"' > server.out 2>&1
```

```bash
nc localhost 4242 -c '
set -x
echo "Hello from Client"' > client.out 2>&1
```

```bash
(echo -e "Hello\nFrom\nServer\n" && cat 2>/dev/null) | nc -ls localhost -p 4242 > ./telnet.server.out &
```

```bash
(echo -e "Hello\nFrom\nClient\n" && cat 2>/dev/null) | nc localhost 4242 > ./telnet.client.out &
```

```bash
(echo -e "Hello\nFrom\nNC\n" && cat 2>/dev/null) | nc -ls localhost -p 4242 >./nc.out.txt
NC_PID=$!

sleep 1
# /tmp/telnet.out
(echo -e "I\nam\nTELNET client\n" && cat 2>/dev/null) | nc localhost 4242 >./telnet.out
```
