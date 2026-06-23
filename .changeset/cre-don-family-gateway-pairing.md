---
"chainlink": patch
---

Require explicit `don_family` on workflow and gateway nodesets for all local CRE gateway topologies (N >= 1 families). Gateway connectors, gateway worker jobs, capabilities registry families, and `env workflow deploy` are always scoped by family; deploy without local state still defaults to `DefaultDONFamily`. #internal
