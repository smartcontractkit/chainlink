---
"chainlink": minor
---

#added Added an auto-purge feature to the EVM TXM that identifies terminally stuck transactions either through a chain specific method or heurisitic then purges them to unblock the nonce. Includes 4 new toml configs to configure this feature.
