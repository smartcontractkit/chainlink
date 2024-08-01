package config

// setConfigEncodedComponentsABI specifies the serialization schema for the
// encoded config, in a form which can be parsed by abigen. The "name" of each
// component must match the name of the corresponding field in
// setConfigSerializationTypes, and the "name" of each component in
// "sharedSecretEncryptions" must match the name of the corresponding field in
// sseSerializationTypes.
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
        "name": "alphaPPB",
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
