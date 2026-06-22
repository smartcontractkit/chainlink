---
"chainlink": patch
---

Scope local CRE gateway connectors and gateway worker jobs by explicit nodeset `don_family` for multi-family topologies. Capabilities registry `DonFamilies`, workflow registry `SetDONLimit` per family, and `env workflow deploy` (`--don-family`, `--workflow-don-name`) read TOML `don_family`; legacy single-DON topologies fall back to `DefaultDONFamily`. #internal
