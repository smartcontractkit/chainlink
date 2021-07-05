package testspecs

var (
	RandomnessJob = `
{
  "initiators": [
    {
      "type": "randomnesslog",
      "params": {
        "address": "0xaba5edc1a551e55b1a570c0e1f1055e5be11eca7",
        "_comment": "should be the address of the VRF Coordinator"
      }
    }
  ],
  "tasks": [
    {
      "type": "random",
      "confirmations": 30,
      "params": {
        "_comment": "Note: the following key is ONLY AN EXAMPLE, and not secure.",
        "_comment2": "Use the public key reported when you ran chainlink local vrf create, instead",
        "publicKey":"0x79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F8179800",
        "_comment3": "Corresponds to a secret key of 1. (So not secure at all!)"
      }
    },
    {
      "type": "ethtx",
      "params": {
        "format": "preformatted",
        "comment": "ethereum address of the VRF coordinator contract goes in 'address' field:",
        "address": "0x5e1f1e555ca1ab1eb01dfaceca11ab1eba5eba11",
        "comment2": "functionSelector from VRFCoordinator.sol, and javascript call:",
        "comment3": "web3.eth.abi.encodeFunctionSignature('fulfillRandomnessRequest(bytes)')",
        "functionSelector": "0x5e1c1059"
      }
    }
  ]
}
`
)
