---
"chainlink": patch
---

Performance improvements and bugfixes for LLO transmission queue #internal #fixed

- Reduces required number of DB connections and reduce overall DB transaction load
- Remove duplicate deletion code from server.go and manage everything in the persistence manager
- Introduce an application-wide global reaper for last-ditch cleanup effort
- Implement delete batching for more reliable and incremental deletion
- Ensure that records are properly removed on exit
