---
"chainlink": patch
---

Remove finality depth as the default value for minConfirmation for tx jobs. 
Update the sql query for fetching pending callback transactions:
if minConfirmation is not null, we check difference if the current block - tx block > minConfirmation
else we check if the tx block is <= finalizedBlock 
#updated
