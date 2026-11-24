---
"chainlink": patch
---

#bugfix Vault gateway rejects secrets requests whose params the Vault nodes would re-derive differently (fields the request type does not define, a missing namespace, or a request_id other than the request ID) before authorizing them, instead of forwarding requests that no node can authorize
