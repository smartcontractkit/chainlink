# Reconstructing these files

## vrfkey.json

This is the encrypted secret key used to generate VRF proofs. Its public key is
`0xce3cc486a3aa567a2e15707cd5b874170e746f4db45ec084918e5b2b5ec6c6cf01`

Creation commands:

```
# Create key

./tools/bin/cldev local vrf \
   xxxCreateWeakKeyPeriodYesIReallyKnowWhatIAmDoingAndDoNotCareAboutThisKeyMaterialFallingIntoTheWrongHandsExclamationPointExclamationPointExclamationPointExclamationPointIAmAMasochistExclamationPointExclamationPointExclamationPointExclamationPointExclamationPoint \
   -f ./tools/clroot/vrfkey.json -p ./tools/clroot/password.txt
```

Note that this will produce a different key, though.
