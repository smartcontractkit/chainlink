/**
 * @description this script generates a master interface for interacting with the automation registry
 * @notice run this script with pnpm ts-node ./scripts/generate-zksync-automation-master-interface-v2_3.ts
 */
import { ZKSyncAutomationRegistry2_3__factory as Registry } from '../typechain/factories/ZKSyncAutomationRegistry2_3__factory'
import { ZKSyncAutomationRegistryLogicA2_3__factory as RegistryLogicA } from '../typechain/factories/ZKSyncAutomationRegistryLogicA2_3__factory'
import { ZKSyncAutomationRegistryLogicB2_3__factory as RegistryLogicB } from '../typechain/factories/ZKSyncAutomationRegistryLogicB2_3__factory'
import { ZKSyncAutomationRegistryLogicC2_3__factory as RegistryLogicC } from '../typechain/factories/ZKSyncAutomationRegistryLogicC2_3__factory'
import { utils } from 'ethers'
import fs from 'fs'
import { exec } from 'child_process'

const dest = 'src/v0.8/automation/interfaces/zksync'
const srcDest = `${dest}/IZKSyncAutomationRegistryMaster2_3.sol`
const tmpDest = `${dest}/tmp.txt`

const combinedABI = []
const abiSet = new Set()
const abis = [
  Registry.abi,
  RegistryLogicA.abi,
  RegistryLogicB.abi,
  RegistryLogicC.abi,
]

for (const abi of abis) {
  for (const entry of abi) {
    const id = utils.id(JSON.stringify(entry))
    if (!abiSet.has(id)) {
      abiSet.add(id)
      if (
        entry.type === 'function' &&
        (entry.name === 'checkUpkeep' ||
          entry.name === 'checkCallback' ||
          entry.name === 'simulatePerformUpkeep')
      ) {
        entry.stateMutability = 'view' // override stateMutability for check / callback / simulate functions
      }
      combinedABI.push(entry)
    }
  }
}

const checksum = utils.id(abis.join(''))

fs.writeFileSync(`${tmpDest}`, JSON.stringify(combinedABI))

const cmd = `
cat ${tmpDest} | pnpm abi-to-sol --solidity-version ^0.8.19 --license MIT > ${srcDest} IZKSyncAutomationRegistryMaster2_3;
echo "// solhint-disable \n// abi-checksum: ${checksum}" | cat - ${srcDest} > ${tmpDest} && mv ${tmpDest} ${srcDest};
pnpm prettier --write ${srcDest};
`

exec(cmd)

console.log(
  'generated new master interface for zksync automation registry v2_3',
)
