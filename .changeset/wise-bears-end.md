---
"chainlink": minor
---

#updated mercury plugin to consider PluginConfig as optional if EnableTriggerCapability relay config is true. Then if PluginConfig is nil, skip fetching latestPrice for linkFeedId and nativeFeedId.
