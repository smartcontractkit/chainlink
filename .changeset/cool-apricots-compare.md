---
"chainlink": patch
---

Increase default config for postgres max open conns from 20 to 100.

Also, add autoscaling for mercury jobs. The max open conns limit will be
automatically increased to the number of mercury jobs if this exceeds the
configured value.
