---
"chainlink": patch
---

#changed Rename the `InBackupHealthReport` to `StartUpHealthReport` and enable it for DB migrations as well. This will enable health report to be available during long start-up tasks (db backups and migrations).