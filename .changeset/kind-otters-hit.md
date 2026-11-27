---
"chainlink": patch
---

#bugfix When a gateway host dies without closing its sockets, workflow nodes keep writing into the dead websocket. this introduces a fix that will reduce the time it takes the node to be aware from ~17min to ~30sec
