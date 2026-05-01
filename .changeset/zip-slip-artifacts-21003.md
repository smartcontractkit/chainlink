---
"chainlink": patch
---

#bugfix

`downloadProgramArtifacts` in `deployment/utils/solutils/artifacts.go` now rejects archive entries whose names would resolve outside the target directory and validates the absolute extraction path against the configured target directory. This hardens the Solana artifact downloader against Zip Slip / path-traversal attacks where a malicious archive could otherwise overwrite arbitrary files. Closes #21003.
