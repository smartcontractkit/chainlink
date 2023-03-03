## go-sockaddr - `{Raw,}Sockaddr` conversions

See https://groups.google.com/d/msg/golang-nuts/B-meiFfkmH0/-TxP1r6zvk8J
This package extracts unexported code from `golang.org/x/unix` to help in converting
between:

```Go
${platform}.Sockaddr
${platform}.RawSockaddrAny
C.struct_sockaddr_any
net.*Addr
```

Godoc:

- sockaddr - http://godoc.org/github.com/libp2p/go-sockaddr
- sockaddr/net - http://godoc.org/github.com/libp2p/go-sockaddr/net

---

The last gx published version of this module was: 1.0.3: QmNzEyX7vjWiqinyLeavcAF1oegav6dZ1aQpAkYvG9m5Ze
