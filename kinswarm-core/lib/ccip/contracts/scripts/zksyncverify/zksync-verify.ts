import hre from 'hardhat'
//import big from 'ethers'
// pass in the correct constrcutor arguments for the contract &
// run with `npx hardhat run zksync-verify.ts --config ./hardhat.ccip.zksync.config.ts` & remember to change the appropriate `default network` in the config file or pass as argument
async function main() {
  await hre.run('verify:verify', {
    address: '0xb575772CD01478477eA649b82C30FD63295D7A4d',
    constructorArguments: [
      {
        voters: [
          {
            blessVoteAddr: '0x94e46574182e03cab3557ecb780643b66946970a',
            curseVoteAddr: '0x3e09450412c5765b0b91e3c7b15dd523e460aeae',
            blessWeight: 1,
            curseWeight: 1,
          },
        ],
        blessWeightThreshold: 1,
        curseWeightThreshold: 1,
      },
    ],
  })
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
