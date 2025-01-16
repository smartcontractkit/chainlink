import hre from 'hardhat'

// pass in the correct constrcutor arguments for the contract &
// Run with `npx hardhat run zksync-verify.ts --config ./hardhat.ccip.zksync.config.ts` & remember to change the appropriate `default network` in the config file or pass as argument
async function main() {
  await hre.run('verify:verify', {
    address: '0x37CbA662E9c373F2166CcA0D9c576dd089D7209a',
    constructorArguments: [
      /*{
        "linkToken": "0x52869bae3e091e36b0915941577f2d47d8d8b534",
        "chainSelector": "1562403441176082196",
        "destChainSelector": "4051577828743386545",
        "defaultTxGasLimit": 200000,
        "maxNopFeesJuels": "20000000000000000000000",
        "prevOnRamp": "0x0000000000000000000000000000000000000000",
        "rmnProxy": "0x2abb46a2d32220b8801ce96cabc32dd2da7b7b20",
        "tokenAdminRegistry": "0x100a47c9db342884e3314b91cec076bbac8e619c"
      },
      {
        "router": "0x63212825283b37528759c542ee22d70277c954d8",
        "maxNumberOfTokensPerMsg": 1,
        "destGasOverhead": 300000,
        "destGasPerPayloadByte": 16,
        "destDataAvailabilityOverheadGas": 0,
        "destGasPerDataAvailabilityByte": 16,
        "destDataAvailabilityMultiplierBps": 0,
        "priceRegistry": "0xfdae74b78045e0cdb725d57bd52cd27d92c506e7",
        "maxDataBytes": 30000,
        "maxPerMsgGasLimit": 3000000,
        "defaultTokenFeeUSDCents": 25,
        "defaultTokenDestGasOverhead": 90000,
        "enforceOutOfOrder": false
      },
      {
        "isEnabled": false,
        "capacity": 0,
        "rate": 0
      },
      [
        {
          "token": "0x52869bae3e091e36b0915941577f2d47d8d8b534",
          "networkFeeUSDCents": 10,
          "gasMultiplierWeiPerEth": "1100000000000000000",
          "premiumMultiplierWeiPerEth": "900000000000000000",
          "enabled": true
        },
        {
          "token": "0x5aea5775959fbc2557cc8789bc1bf90a239d9a91",
          "networkFeeUSDCents": 10,
          "gasMultiplierWeiPerEth": "1100000000000000000",
          "premiumMultiplierWeiPerEth": "1000000000000000000",
          "enabled": true
        }
      ],
      [],
      [] 
    */
    ],
  })
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
