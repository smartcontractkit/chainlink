#bugfix

---
"chainlink": patch
---

Fix EVMService.SubmitTransaction to return TxReverted when receipt status is 0 (reverted), instead of always returning TxSuccess. This resolves OCR consensus failures on Polygon WriteReport where DON nodes returned inconsistent results for reverted transactions.
