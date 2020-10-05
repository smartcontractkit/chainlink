package config

const setConfigEncodedComponentsABI = `[
  {
    "name": "setConfigEncodedComponents",
    "type": "tuple",
    "components": [
      {
        "name": "deltaProgress",
        "type": "int64"
      },
      {
        "name": "deltaResend",
        "type": "int64"
      },
      {
        "name": "deltaRound",
        "type": "int64"
      },
      {
        "name": "deltaGrace",
        "type": "int64"
      },
      {
        "name": "deltaC",
        "type": "int64"
      },
      {
        "name": "alpha",
        "type": "uint64"
      },
      {
        "name": "deltaStage",
        "type": "int64"
      },
      {
        "name": "rMax",
        "type": "uint8"
      },
      {
        "name": "s",
        "type": "uint8[]"
      },
      {
        "name": "offchainPublicKeys",
        "type": "bytes32[]"
      },
      {
        "name": "peerIDs",
        "type": "string"
      },
      {
        "name": "sharedSecretEncryptions",
        "type": "tuple",
        "components": [
          {
            "name": "diffieHellmanPoint",
            "type": "bytes32"
          },
          {
            "name": "sharedSecretHash",
            "type": "bytes32"
          },
          {
            "name": "encryptions",
            "type": "bytes16[]"
          }
        ]
      }
    ]
  }
]`
