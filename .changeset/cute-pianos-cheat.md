---
"chainlink": minor
---

**DirectRequest and FluxMonitor job types have been removed.** Creating new jobs of these types is no longer supported and will return an error. Any existing jobs of these types that are still present in the database will surface an error in the job UI on node startup rather than running. The underlying database tables (`direct_request_specs`, `flux_monitor_specs`, `flux_monitor_round_stats_v2`) are **unchanged in this release** and will be cleaned up in a future migration. The `[FluxMonitor]` TOML config section is now a no-op but is still accepted to avoid breaking existing config files during the transition. #breaking_change #nops
