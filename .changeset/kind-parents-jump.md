---
"chainlink": patch
---

Add two new metrics for monitoring LLO transmitter health #added

`llo_mercurytransmitter_concurrent_transmit_gauge`
Gauge that measures the number of transmit threads currently waiting on a remote transmit call. You may wish to alert if this exceeds some number for a given period of time, or if it ever reaches its max.

`llo_mercurytransmitter_concurrent_delete_gauge`
Gauge that measures the number of delete threads currently waiting on a delete call to the DB. You may wish to alert if this exceeds some number for a given period of time, or if it ever reaches its max.
