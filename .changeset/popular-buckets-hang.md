---
"chainlink": patch
---

Add new config option Pipeline.VerboseLogging

VerboseLogging enables detailed logging of pipeline execution steps. This is
disabled by default because it increases log volume for pipeline runs, but can
be useful for debugging failed runs without relying on the UI or database.
Consider enabling this if you disabled run saving by setting MaxSuccessfulRuns
to zero.

Set it like the following example:

```
[Pipeline]
VerboseLogging = true
```
