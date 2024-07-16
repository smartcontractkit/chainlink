/**
 * @description this script generates a master interface for interacting with the automation registry
 * @notice run this script with pnpm ts-node ./scripts/generate-automation-master-interface-v2_3.ts
 */
import { AutomationRegistry2_3__factory as Registry } from '../typechain/factories/AutomationRegistry2_3__factory'
import { AutomationRegistryLogicA2_3__factory as RegistryLogicA } from '../typechain/factories/AutomationRegistryLogicA2_3__factory'
import { AutomationRegistryLogicB2_3__factory as RegistryLogicB } from '../typechain/factories/AutomationRegistryLogicB2_3__factory'
import { AutomationRegistryLogicC2_3__factory as RegistryLogicC } from '../typechain/factories/AutomationRegistryLogicC2_3__factory'
import { utils } from 'ethers'
import fs from 'fs'
import { exec } from 'child_process'

const dest = 'src/v0.8/automation/interfaces/v2_3'
const srcDest = `${dest}/IAutomationRegistryMaster2_3.sol`
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
cat ${tmpDest} | pnpm abi-to-sol --solidity-version ^0.8.4 --license MIT > ${srcDest} IAutomationRegistryMaster2_3;
echo "// abi-checksum: ${checksum}" | cat - ${srcDest} > ${tmpDest} && mv ${tmpDest} ${srcDest};
pnpm prettier --write ${srcDest};
`

exec(cmd)

console.log('generated new master interface for automation registry')
