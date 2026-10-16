---
"chainlink": patch
---

#bugfix Registry-launched OCR2 plugins (DON Time) now match multichain signers covering several chain families (e.g. EVM + Solana + Aptos), signing with every matched family's OCR key bundle and using the transmitter paired with the node's signer. Launch fails when there is no registry OCR config or no signer in it matches the node's keys.
