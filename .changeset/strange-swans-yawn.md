---
"chainlink": patch
---

#changed

Productionize transmitter for LLO

Note that some minor changes to prometheus metrics will occur in the transition to LLO. Since feed IDs no longer apply, the metrics for transmissions change as follows:

```
"mercury_transmit_*"
[]string{"feedID", ...},
```

Will change to:

```
"llo_mercury_transmit_*"
[]string{"donID", ...},
```
