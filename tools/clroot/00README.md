# Reconstructing these files

## vrfkey.json

This is the encrypted secret key used to generate VRF proofs. Its public key is

`0x1e00cb99ccc4ab427b9c3d4b5689e4bd5ceecdcadfcb2dd6719b676e666177d2579ce2fe32ae517b91569b11335fb09606c84b96c254f70c05d16ef8461e5b4b`

Creation commands:

```
# Create key

$ ./tools/bin/cldev local vrf create -p ./tools/clroot/password.txt
Created keypair, with public key
0xe11c3b009e1977901b247e47892d2db315d5d4c6d2c6b887831ab6095018b72c9f09f489fb91d9c37cf34bd9ea9bd141aee8d257b29042949429c4e8e8b27db4

The following command will export the encrypted secret key from the db to <save_path>:

chainlink local vrf export -f <save_path> -pk 0xe11c3b009e1977901b247e47892d2db315d5d4c6d2c6b887831ab6095018b72c9f09f489fb91d9c37cf34bd9ea9bd141aee8d257b29042949429c4e8e8b27db4
```

For instance, to use the last suggested command to save the secret key to the canonical location:

```
$ ./tools/bin/cldev local vrf export -f ./tools/clroot/vrfkey.json -pk 0xe11c3b009e1977901b247e47892d2db315d5d4c6d2c6b887831ab6095018b72c9f09f489fb91d9c37cf34bd9ea9bd141aee8d257b29042949429c4e8e8b27db4
```

The secret key associated with this key is
0xaad746fe6ebbee36692d7f81a194e70db9647b98372ecc801196266a835214ca .
