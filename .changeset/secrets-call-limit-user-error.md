---
"chainlink": patch
---

Classify the per-execution secrets call limit as a user LimitExceeded error, so the confidential relay answers it with invalid params instead of an internal error. #bugfix
