---
"chainlink": patch
---

Teach `utils.ToDecimal` to convert `json.Number` so pipeline tasks that consume HTTP JSON payloads (ethabiencode, fluxmonitor, OCR data sources) no longer fail with "type json.Number cannot be converted to decimal.Decimal".

#bugfix
