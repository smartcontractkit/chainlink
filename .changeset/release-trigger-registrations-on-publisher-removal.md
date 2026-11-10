---
"chainlink": patch
---

Unregister a trigger publisher's registrations from the underlying capability when the launcher removes the publisher, so the same triggers can register again when a capability moves to another DON on the same node. The unregister runs in the background, and new registrations for that capability wait until it finishes. #bugfix
