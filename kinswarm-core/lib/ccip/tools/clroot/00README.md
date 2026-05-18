# Reconstructing these files

## vrfkey.json

This is the encrypted secret key used to generate VRF proofs. Its public key is

``

Creation commands:

```
# Create key

./tools/bin/cldev local vrf \
   createWeakKeyPeriodYesIReallyKnowWhatIAmDoingAndDoNotCareAboutThisKeyMaterialFallingIntoTheWrongHandsExclamationPointExclamationPointExclamationPointExclamationPointIAmAMasochistExclamationPointExclamationPointExclamationPointExclamationPointExclamationPoint \
   -f ./tools/clroot/vrfkey.json -p ./tools/clroot/password.txt
```
